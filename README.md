# S3 Presigned URL 생성 Lambda 함수

이 프로젝트는 AWS Lambda를 사용하여 S3 Presigned URL을 생성하는 API를 제공합니다. 이를 통해 클라이언트는 S3 버킷에 직접 파일을 업로드하거나, S3 버킷에서 파일을 다운로드할 수 있습니다.

## 기능

- GET 또는 PUT 메서드에 대한 Presigned URL 생성
- 콘텐츠 타입 자동 감지 및 설정 지원
- 파일 크기 제한 설정 기능
- CORS 지원
- URL 만료 시간 설정 기능

## 요청 예시

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

- `bucket`: S3 버킷 이름 (선택적, 환경 변수에서 기본값 사용 가능)
- `key`: 파일 경로/이름 (필수)
- `method`: HTTP 메서드, `GET` 또는 `PUT` (선택적, 기본값: `GET`)
- `duration`: URL 유효 기간(초) (선택적, 기본값: 3600초)
- `contentType`: 파일 콘텐츠 타입, PUT 요청 시에만 유효 (선택적, 자동 감지)
- `maxFileSize`: 최대 파일 크기(바이트), PUT 요청 시에만 유효 (선택적, 환경 변수에서 기본값 사용)

## 응답 예시

```json
{
  "url": "https://my-s3-bucket.s3.amazonaws.com/uploads/image.jpg?X-Amz-Algorithm=...",
  "contentType": "image/jpeg",
  "maxFileSize": 5242880,
  "expiresAt": 1623456789
}
```

- `url`: 생성된 presigned URL
- `contentType`: 콘텐츠 타입 (PUT 요청 시에만 반환)
- `maxFileSize`: 최대 파일 크기 (PUT 요청 시에만 반환)
- `expiresAt`: URL 만료 시간 (Unix 타임스탬프)

## 배포 방법

### 사전 요구사항

- AWS CLI 설치 및 구성
- AWS SAM CLI 설치
- Go 개발 환경 설정

### 배포 단계

1. 의존성 설치:

```bash
make init
```

2. 빌드:

```bash
make build
```

3. 패키징:

```bash
make package S3_BUCKET=deployment-bucket-name
```

4. 배포:

```bash
make deploy S3_BUCKET=deployment-bucket-name PARAMETER_OVERRIDES="DefaultBucketName=your-s3-bucket Environment=prod DefaultMaxFileSize=10485760"
```

## 로컬에서 테스트

1. 테스트 이벤트 생성:

```bash
make create-event
```

2. 로컬에서 함수 실행:

```bash
make local-invoke
```

3. 로컬 API 시작:

```bash
make local-api
```

## 클라이언트 사용 예시 (JavaScript)

```javascript
// PUT(업로드) URL 요청
async function getUploadURL(filename) {
  const response = await fetch("https://your-api-endpoint/presigned-url", {
    method: "POST",
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

// 파일 업로드
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

// GET(다운로드) URL 요청
async function getDownloadURL(fileKey) {
  const response = await fetch("https://your-api-endpoint/presigned-url", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      key: fileKey,
      method: "GET",
      duration: 7200, // 2시간
    }),
  });

  return await response.json();
}
```
