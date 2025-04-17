# AWS Lambda API 서버

이 프로젝트는 AWS Lambda와 API Gateway를 사용하여 다음 기능을 제공하는 서버리스 API 서버입니다:

1. S3 Presigned URL 생성 기능
2. 매장 요청 관리 기능 (관리자용)

## 기능 개요

### 1. S3 Presigned URL 생성 API

클라이언트가 AWS S3에 직접 파일을 업로드하거나 다운로드할 수 있도록 하는 Presigned URL을 생성합니다.

- GET/PUT 메서드에 대한 Presigned URL 생성
- 콘텐츠 타입 자동 감지 및 설정
- URL 만료 시간 설정 기능
- CORS 지원

### 2. 관리자 API

매장 생성 요청을 관리하는 기능을 제공합니다.

- 매장 생성 요청 목록 조회
- 매장 생성 요청 승인/거절 처리

## API 엔드포인트

### S3 Presigned URL API

#### `GET /s3/presigned-url`

S3 업로드/다운로드용 Presigned URL을 생성합니다.

**요청 예시:**

```json
{
  "bucket": "my-s3-bucket",
  "key": "uploads/image.jpg",
  "method": "PUT",
  "duration": 3600,
  "contentType": "image/jpeg",
  "maxFileSize": 5242880
}
```

**응답 예시:**

```json
{
  "url": "https://my-s3-bucket.s3.amazonaws.com/uploads/image.jpg?X-Amz-Algorithm=...",
  "contentType": "image/jpeg",
  "maxFileSize": 5242880,
  "expiresAt": 1623456789
}
```

### 관리자 API

#### `GET /admin/restaurant/request`

매장 생성 요청 목록을 조회합니다.

**응답 예시:**

```json
{
  "requests": [
    {
      "id": 1,
      "restaurantId": "550e8400-e29b-41d4-a716-446655440000",
      "userId": "123456789",
      "status": "PENDING",
      "createdAt": "2023-04-01T12:00:00Z",
      "updatedAt": "2023-04-01T12:00:00Z"
    }
  ],
  "total": 1
}
```

#### `POST /admin/restaurant/request/{id}/process`

매장 생성 요청을 승인하거나 거절합니다.

**요청 예시:**

```json
{
  "status": "APPROVED",
  "rejectReason": null
}
```

**응답 예시:**

```json
{
  "id": 1,
  "restaurantId": "550e8400-e29b-41d4-a716-446655440000",
  "userId": "123456789",
  "status": "APPROVED",
  "rejectReason": null,
  "createdAt": "2023-04-01T12:00:00Z",
  "updatedAt": "2023-04-01T12:30:00Z"
}
```

## 인증

- `/admin/` 경로의 API는 세션 인증(SessionAuth)이 필요합니다.
- `/s3/presigned-url` API는 인증이 필요하지 않습니다.

## 프로젝트 구조

```
lambda-go/
├── main.go                 # 주 진입점 및 핸들러
├── pkg/
│   ├── configs/            # 환경 설정 관련 코드
│   ├── contexts/           # 컨텍스트 관련 유틸리티
│   ├── handlers/           # API 핸들러
│   │   ├── admin/          # 관리자 핸들러
│   │   └── public/         # 공개 핸들러
│   ├── middlewares/        # 인증 및 요청 처리 미들웨어
│   ├── models/             # 데이터 모델
│   ├── repositories/       # 데이터베이스 액세스 레이어
│   ├── routes/             # 라우팅 정의
│   ├── services/           # 비즈니스 로직
│   └── utils/              # 유틸리티 함수
└── template.yaml           # AWS SAM 템플릿
```

## 배포 방법

### 사전 요구사항

- AWS CLI 설치 및 구성
- AWS SAM CLI 설치
- Go 개발 환경 설정
- PostgreSQL 데이터베이스 접근 정보

### 배포 단계

1. 의존성 설치:

```bash
go mod tidy
```

2. 빌드:

```bash
GOARCH=amd64 GOOS=linux go build -o main
```

3. AWS SAM을 사용한 배포:

```bash
sam deploy --template-file template.yaml --stack-name multi-purpose-lambda --capabilities CAPABILITY_IAM --parameter-overrides \
DefaultBucketName=your-s3-bucket \
Environment=prod \
DefaultMaxFileSize=10485760 \
DBHost=your-db-host \
DBPort=5432 \
DBUser=your-db-user \
DBPassword=your-db-password \
DBName=your-db-name \
DBSSLMode=require
```

## 클라이언트 사용 예시 (JavaScript)

### Presigned URL을 사용한 파일 업로드

```javascript
async function getUploadURL(filename) {
  const response = await fetch("https://your-api-endpoint/s3/presigned-url", {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      key: `uploads/${filename}`,
      method: "PUT",
    }),
  });

  return await response.json();
}

async function uploadFile(file) {
  const { url, contentType } = await getUploadURL(file.name);

  await fetch(url, {
    method: "PUT",
    headers: {
      "Content-Type": contentType,
    },
    body: file,
  });

  return url.split("?")[0]; // 순수 S3 URL 반환
}
```

### 관리자 API 사용 예시

```javascript
// 매장 생성 요청 목록 조회
async function getRestaurantRequests() {
  const response = await fetch(
    "https://your-api-endpoint/admin/restaurant/request",
    {
      method: "GET",
      headers: {
        Authorization: "Bearer your-auth-token",
      },
    }
  );

  return await response.json();
}

// 매장 생성 요청 처리
async function processRestaurantRequest(
  requestId,
  approved,
  rejectReason = null
) {
  const response = await fetch(
    `https://your-api-endpoint/admin/restaurant/request/${requestId}/process`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: "Bearer your-auth-token",
      },
      body: JSON.stringify({
        status: approved ? "APPROVED" : "REJECTED",
        rejectReason: approved ? null : rejectReason,
      }),
    }
  );

  return await response.json();
}
```
