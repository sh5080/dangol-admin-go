.PHONY: build clean deploy

# 변수 정의
BINARY_NAME=main
STACK_NAME=admin-lambda
S3_BUCKET=jdg-admin-lambda
PARAMETER_OVERRIDES=
AWS_PROFILE?=sh5080

# 기본 타겟
all: clean build

# 빌드: Go Lambda 함수 컴파일
build:
	@echo "Building Go Lambda function..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) main.go
	chmod +x $(BINARY_NAME)
	echo '#!/bin/sh\n./$(BINARY_NAME)' > bootstrap
	chmod +x bootstrap
	@echo "Build completed"

# Go 모듈 초기화 및 의존성 다운로드
init:
	@echo "Initializing Go modules..."
	go mod tidy
	@echo "Initialization completed"

# 클린: 빌드 산출물 제거
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	@echo "Clean completed"

# 패키지: SAM 패키지 생성
package:
	@echo "Packaging SAM application..."
	AWS_PROFILE=$(AWS_PROFILE) sam package \
		--template-file template.yaml \
		--output-template-file packaged.yaml \
		--s3-bucket $(S3_BUCKET)
	@echo "Package completed"

# 배포: SAM 배포
deploy:
	@echo "Deploying SAM application..."
	AWS_PROFILE=$(AWS_PROFILE) sam deploy \
		--template-file packaged.yaml \
		--stack-name $(STACK_NAME) \
		--capabilities CAPABILITY_IAM \
		$(if $(PARAMETER_OVERRIDES),--parameter-overrides $(PARAMETER_OVERRIDES),)
	@echo "Deployment completed"

# 로컬 테스트: SAM 로컬로 API 실행
local-api:
	@echo "Starting local API Gateway..."
	sam local start-api
	@echo "Local API stopped"

# 로컬 호출: Lambda 함수 로컬 호출
local-invoke: build
	@echo "Invoking Lambda function locally..."
	sam local invoke PresignedURLFunction --event events/event.json
	@echo "Local invocation completed"

# 테스트 이벤트 생성
create-event:
	@echo "Creating test event..."
	mkdir -p events
	echo '{ "body": "{ \"bucket\": \"test-bucket\", \"key\": \"test/file.jpg\", \"method\": \"PUT\", \"duration\": 3600 }" }' > events/event.json
	@echo "Test event created at events/event.json"

# AWS 리소스 삭제
delete:
	@echo "Deleting AWS resources..."
	aws cloudformation delete-stack --stack-name $(STACK_NAME)
	@echo "Deletion request sent"

# 배포된 API URL 얻기
get-api-url:
	@echo "Getting deployed API URL..."
	aws cloudformation describe-stacks \
		--stack-name $(STACK_NAME) \
		--query 'Stacks[0].Outputs[?OutputKey==`APIEndpoint`].OutputValue' \
		--output text
	@echo "API URL retrieved"

# 도움말
help:
	@echo "사용법: make [타겟]"
	@echo ""
	@echo "타겟:"
	@echo "  build         Lambda 함수 빌드"
	@echo "  init          Go 모듈 초기화 및 의존성 다운로드"
	@echo "  clean         빌드 산출물 제거"
	@echo "  package       SAM 패키지 생성 (S3_BUCKET을 지정해야 함)"
	@echo "  deploy        SAM 배포 (S3_BUCKET을 지정해야 함)"
	@echo "  local-api     SAM 로컬로 API 실행"
	@echo "  local-invoke  Lambda 함수 로컬 호출"
	@echo "  create-event  테스트 이벤트 생성"
	@echo "  delete        AWS 리소스 삭제"
	@echo "  get-api-url   배포된 API URL 얻기"
	@echo ""
	@echo "예시:"
	@echo "  make package S3_BUCKET=your-deployment-bucket"
	@echo "  make deploy S3_BUCKET=your-deployment-bucket PARAMETER_OVERRIDES=\"DefaultBucketName=your-s3-bucket Environment=prod\""
