GO_PACKAGE_NAME=shopifyoauth
GO_PACKAGE_PATH=github.com/tkeech1/${GO_PACKAGE_NAME}

build-docker:
	docker build -t ${GO_PACKAGE_NAME}:latest .

build-docker-nocache:
	docker build --no-cache -t ${GO_PACKAGE_NAME}:latest .

deploy-dev: build-docker
	docker run -it --rm -e "AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}" -e "AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}" -e "AWS_REGION=${AWS_REGION}" -e "JWT_KEY=${JWT_KEY}" -e "SHOPIFY_SHARED_SECRET=${SHOPIFY_SHARED_SECRET}" -e "SHOPIFY_API_KEY=${SHOPIFY_API_KEY}" ${GO_PACKAGE_NAME}:latest serverless deploy --stage=dev -v

undeploy-dev:
	docker run -it --rm -e "AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}" -e "AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}" -e "AWS_REGION=${AWS_REGION}" ${GO_PACKAGE_NAME}:latest serverless remove --stage=dev -v

updatecode-dev-oauth_install: build-docker
	docker run -it --rm -e "AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}" -e "AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}" -e "AWS_REGION=${AWS_REGION}" -e "JWT_KEY=${JWT_KEY}" -e "SHOPIFY_SHARED_SECRET=${SHOPIFY_SHARED_SECRET}" -e "SHOPIFY_API_KEY=${SHOPIFY_API_KEY}" ${GO_PACKAGE_NAME}:latest serverless deploy function --function oauth_install --stage=dev -v

updatecode-dev-oauth_callback: build-docker 
	docker run -it --rm -e "AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}" -e "AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}" -e "AWS_REGION=${AWS_REGION}" -e "JWT_KEY=${JWT_KEY}" -e "SHOPIFY_SHARED_SECRET=${SHOPIFY_SHARED_SECRET}" -e "SHOPIFY_API_KEY=${SHOPIFY_API_KEY}" ${GO_PACKAGE_NAME}:latest serverless deploy function --function oauth_callback --stage=dev -v
