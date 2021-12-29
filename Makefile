include .env
export

.PHONY: openapi
openapi: openapi_http openapi_js

.PHONY: openapi_http
openapi_http:
	oapi-codegen -generate types -o internal/checkouts/ports/openapi_types.gen.go -package ports api/openapi/checkouts.yml
	oapi-codegen -generate chi-server -o internal/checkouts/ports/openapi_api.gen.go -package ports api/openapi/checkouts.yml
	oapi-codegen -generate types -o internal/common/client/checkouts/openapi_types.gen.go -package checkouts api/openapi/checkouts.yml
	oapi-codegen -generate client -o internal/common/client/checkouts/openapi_client_gen.go -package checkouts api/openapi/checkouts.yml

	oapi-codegen -generate types -o internal/orders/ports/openapi_types.gen.go -package ports api/openapi/orders.yml
	oapi-codegen -generate chi-server -o internal/orders/ports/openapi_api.gen.go -package ports api/openapi/orders.yml
	oapi-codegen -generate types -o internal/common/client/orders/openapi_types.gen.go -package orders api/openapi/orders.yml
	oapi-codegen -generate client -o internal/common/client/orders/openapi_client_gen.go -package orders api/openapi/orders.yml

	oapi-codegen -generate types -o internal/products/ports/openapi_types.gen.go -package ports api/openapi/products.yml
	oapi-codegen -generate chi-server -o internal/products/ports/openapi_api.gen.go -package ports api/openapi/products.yml
	oapi-codegen -generate types -o internal/common/client/products/openapi_types.gen.go -package products api/openapi/products.yml
	oapi-codegen -generate client -o internal/common/client/products/openapi_client_gen.go -package products api/openapi/products.yml

	oapi-codegen -generate types -o internal/users/ports/openapi_types.gen.go -package ports api/openapi/users.yml
	oapi-codegen -generate chi-server -o internal/users/ports/openapi_api.gen.go -package ports api/openapi/users.yml
	oapi-codegen -generate types -o internal/common/client/users/openapi_types.gen.go -package users api/openapi/users.yml
	oapi-codegen -generate client -o internal/common/client/users/openapi_client_gen.go -package users api/openapi/users.yml

.PHONY: openapi_js
openapi_js:
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.0 generate \
        -i /local/api/openapi/checkouts.yml \
        -g javascript \
        -o /local/web/src/repositories/clients/checkouts


	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.0 generate \
        -i /local/api/openapi/orders.yml \
        -g javascript \
        -o /local/web/src/repositories/clients/orders

	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.0 generate \
		-i /local/api/openapi/products.yml \
		-g javascript \
		-o /local/web/src/repositories/clients/products

	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli:v4.3.0 generate \
		-i /local/api/openapi/users.yml \
		-g javascript \
		-o /local/web/src/repositories/clients/users

.PHONY: proto
proto:
	# protoc --go_out=plugins=grpc:internal/common/genproto/orders -I api/protobuf api/protobuf/orders.proto
	# protoc --go_out=plugins=grpc:internal/common/genproto/products -I api/protobuf api/protobuf/products.proto
	# protoc --go_out=plugins=grpc:internal/common/genproto/users -I api/protobuf api/protobuf/users.proto
	
	protoc --go_out=plugins=grpc:internal/common/genproto/orders --go_opt=paths=source_relative -I api/protobuf api/protobuf/orders.proto
	protoc --go_out=plugins=grpc:internal/common/genproto/products --go_opt=paths=source_relative -I api/protobuf api/protobuf/products.proto
	protoc --go_out=plugins=grpc:internal/common/genproto/users --go_opt=paths=source_relative -I api/protobuf api/protobuf/users.proto

.PHONY: lint
lint:
	@./scripts/lint.sh checkouts
	@./scripts/lint.sh orders
	@./scripts/lint.sh products
	@./scripts/lint.sh users

.PHONY: fmt
fmt:
	goimports -l -w internal/

.PHONY: mycli
mycli:
	mycli -u ${MYSQL_USER} -p ${MYSQL_PASSWORD} ${MYSQL_DATABASE}

test:
	@./scripts/test.sh common .e2e.env
	@./scripts/test.sh checkouts .test.env
	@./scripts/test.sh orders .test.env
	@./scripts/test.sh products .test.env
	@./scripts/test.sh users .test.env
