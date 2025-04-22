package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"lambda-go/pkg/models"
	dto "lambda-go/pkg/models/dtos"

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
func (r *RestaurantRepository) GetRestaurantRequests(ctx context.Context, query dto.RestaurantRequestQuery) ([]models.RestaurantRequest, int, error) {
	// 쿼리 빌더 패턴 적용
	whereClause := `WHERE "deletedAt" IS NULL`
	params := []interface{}{}
	paramIndex := 1

	// 상태 필터 적용
	if query.Status != nil {
		whereClause += fmt.Sprintf(` AND "status" = $%d`, paramIndex)
		params = append(params, string(*query.Status))
		paramIndex++
	}

	// 전체 개수 조회
	var total int
	countQuery := `SELECT COUNT(*) FROM "RestaurantRequest" ` + whereClause
	err := r.dbPool.QueryRow(ctx, countQuery, params...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("요청 개수 조회 오류: %w", err)
	}

	// 페이지네이션 적용
	offset := (query.Page - 1) * query.PageSize
	limit := query.PageSize

	// 요청 목록 조회
	queryStr := fmt.Sprintf(`
		SELECT r."id", r."restaurantId", r."userId", r."rejectReason", 
			r."createdAt", r."updatedAt", r."deletedAt", r."status", r."type",
			r."businessLicenseImageUrl", r."businessLicenseNumber"
		FROM "RestaurantRequest" r
		%s
		ORDER BY r."createdAt" DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, paramIndex, paramIndex+1)

	params = append(params, limit, offset)

	rows, err := r.dbPool.Query(ctx, queryStr, params...)
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
			&req.CreatedAt, &req.UpdatedAt, &deletedAt, &req.Status, &req.Type,
			&req.BusinessLicenseImageUrl, &req.BusinessLicenseNumber,
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

	return result, total, nil
}

// GetRestaurantRequestByID는 ID로 매장 요청을 조회합니다.
func (r *RestaurantRepository) GetRestaurantRequestByID(ctx context.Context, requestID string) (models.RestaurantRequestStatus, error) {
	var status string
	query := `SELECT "status" FROM "RestaurantRequest" WHERE "id" = $1 AND "deletedAt" IS NULL`

	err := r.dbPool.QueryRow(ctx, query, requestID).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("요청 ID %s를 찾을 수 없습니다", requestID)
		}
		return "", fmt.Errorf("요청 조회 오류: %w", err)
	}

	return models.RestaurantRequestStatus(status), nil
}

// ProcessRestaurantRequest는 매장 생성 요청을 처리합니다.
func (r *RestaurantRepository) ProcessRestaurantRequest(ctx context.Context, requestID string, payload *models.ProcessRestaurantRequest) (*models.RestaurantRequest, error) {
	// 트랜잭션 시작
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("트랜잭션 시작 오류: %w", err)
	}
	defer tx.Rollback(ctx)

	// 현재 시간
	now := time.Now()

	// 요청 상태 업데이트
	updateRequestQuery := `
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
	err = tx.QueryRow(ctx, updateRequestQuery,
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

	// 승인된 경우에만 Restaurant 상태 업데이트
	if payload.Status == models.APPROVED {
		// Restaurant 상태를 HIDDEN으로 업데이트
		updateRestaurantQuery := `
			UPDATE "Restaurant"
			SET "status" = $1, "updatedAt" = $2
			WHERE "id" = $3
		`

		_, err = tx.Exec(ctx, updateRestaurantQuery,
			models.HIDDEN, // HIDDEN 상태로 설정
			now,
			request.RestaurantID,
		)

		if err != nil {
			return nil, fmt.Errorf("매장 상태 업데이트 오류: %w", err)
		}
	}

	// 트랜잭션 커밋
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("트랜잭션 커밋 오류: %w", err)
	}

	return &request, nil
}
