.PHONY: openapi
openapi: openapi_http

.PHONY: openapi_http
openapi_http:
	oapi-codegen -generate types -o internal/categories/ports/openapi_types.gen.go -package ports api/openapi/categories.yml
	oapi-codegen -generate chi-server -o internal/categories/ports/openapi_api.gen.go -package ports api/openapi/categories.yml
