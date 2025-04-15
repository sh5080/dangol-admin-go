package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"lambda-go/pkg/models"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// RestaurantRepository는 매장 관련 데이터 액세스를 처리합니다.
type RestaurantRepository struct {
	dbPool *pgxpool.Pool
}

// NewRestaurantRepository는 새 RestaurantRepository 인스턴스를 생성합니다.
func NewRestaurantRepository(dbPool *pgxpool.Pool) *RestaurantRepository {
	return &RestaurantRepository{
		dbPool: dbPool,
	}
}

// GetRestaurantRequests는 매장 생성 요청 목록을 조회합니다.
func (r *RestaurantRepository) GetRestaurantRequests(ctx context.Context) ([]models.RestaurantRequest, int, error) {
	// 트랜잭션 시작
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("트랜잭션 시작 오류: %w", err)
	}
	defer tx.Rollback(ctx)
	
	// 전체 개수 조회
	var total int
	countQuery := `SELECT COUNT(*) FROM "RestaurantRequest" WHERE "deletedAt" IS NULL`
	err = tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("요청 개수 조회 오류: %w", err)
	}
	
	// 요청 목록 조회
	query := `
		SELECT r."id", r."restaurantId", r."userId", r."rejectReason", 
			r."createdAt", r."updatedAt", r."deletedAt", r."status"
		FROM "RestaurantRequest" r
		WHERE r."deletedAt" IS NULL
		ORDER BY r."createdAt" DESC
	`
	
	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("요청 목록 조회 오류: %w", err)
	}
	defer rows.Close()
	
	// 모델 변환
	var result []models.RestaurantRequest
	for rows.Next() {
		var req models.RestaurantRequest
		var rejectReason pgtype.Text
		var deletedAt pgtype.Timestamp
		
		err := rows.Scan(
			&req.ID, &req.RestaurantID, &req.UserID, &rejectReason,
			&req.CreatedAt, &req.UpdatedAt, &deletedAt, &req.Status,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("행 스캔 오류: %w", err)
		}
		
		if rejectReason.Status == pgtype.Present {
			reason := rejectReason.String
			req.RejectReason = &reason
		}
		
		if deletedAt.Status == pgtype.Present {
			deleteTime := deletedAt.Time
			req.DeletedAt = &deleteTime
		}
		
		result = append(result, req)
	}
	
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("행 반복 오류: %w", err)
	}
	
	// 트랜잭션 커밋
	if err := tx.Commit(ctx); err != nil {
		return nil, 0, fmt.Errorf("트랜잭션 커밋 오류: %w", err)
	}
	
	return result, total, nil
}

// GetRestaurantRequestByID는 ID로 매장 생성 요청을 조회합니다.
func (r *RestaurantRepository) GetRestaurantRequestByID(ctx context.Context, requestID int) (models.RequestStatus, error) {
	var status string
	query := `SELECT "status" FROM "RestaurantRequest" WHERE "id" = $1 AND "deletedAt" IS NULL`
	
	err := r.dbPool.QueryRow(ctx, query, requestID).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("요청 ID %d를 찾을 수 없습니다", requestID)
		}
		return "", fmt.Errorf("요청 조회 오류: %w", err)
	}
	
	return models.RequestStatus(status), nil
}

// ProcessRestaurantRequest는 매장 생성 요청을 처리합니다.
func (r *RestaurantRepository) ProcessRestaurantRequest(ctx context.Context, requestID int, payload *models.ProcessRestaurantRequest) (*models.RestaurantRequest, error) {
	// 트랜잭션 시작
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("트랜잭션 시작 오류: %w", err)
	}
	defer tx.Rollback(ctx)
	
	// 현재 시간
	now := time.Now()
	
	// 요청 상태 업데이트
	updateQuery := `
		UPDATE "RestaurantRequest"
		SET "status" = $1, "updatedAt" = $2, "rejectReason" = $3
		WHERE "id" = $4 AND "deletedAt" IS NULL
		RETURNING "id", "restaurantId", "userId", "rejectReason", "createdAt", "updatedAt", "deletedAt", "status"
	`
	
	// 결과 저장 변수
	var request models.RestaurantRequest
	var rejectReasonSQL pgtype.Text
	var deletedAtSQL pgtype.Timestamp
	
	// 값 설정
	var rejectReasonVal pgtype.Text
	if payload.RejectReason != nil {
		rejectReasonVal.String = *payload.RejectReason
		rejectReasonVal.Status = pgtype.Present
	} else {
		rejectReasonVal.Status = pgtype.Null
	}
	
	// 업데이트 실행 및 결과 스캔
	err = tx.QueryRow(ctx, updateQuery, 
		payload.Status, now, rejectReasonVal, requestID,
	).Scan(
		&request.ID, &request.RestaurantID, &request.UserID, &rejectReasonSQL,
		&request.CreatedAt, &request.UpdatedAt, &deletedAtSQL, &request.Status,
	)
	
	if err != nil {
		return nil, fmt.Errorf("요청 업데이트 오류: %w", err)
	}
	
	// NULL 값 처리
	if rejectReasonSQL.Status == pgtype.Present {
		reason := rejectReasonSQL.String
		request.RejectReason = &reason
	}
	
	if deletedAtSQL.Status == pgtype.Present {
		deleteTime := deletedAtSQL.Time
		request.DeletedAt = &deleteTime
	}
	
	// 사용자 정보 조회
	userQuery := `
		SELECT "id", "email", "name"
		FROM "user"
		WHERE "id" = $1
	`
	
	var user models.User
	err = tx.QueryRow(ctx, userQuery, request.UserID).Scan(
		&user.ID, &user.Email, &user.Name,
	)
	
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("사용자 조회 오류: %w", err)
		}
		// 사용자가 없는 경우에는 계속 진행
	} else {
		// 사용자 정보 설정
		request.User = &user
	}
	
	// 트랜잭션 커밋
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("트랜잭션 커밋 오류: %w", err)
	}
	
	return &request, nil
} 