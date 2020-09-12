.PHONY: openapi
openapi: openapi_http

.PHONY: openapi_http
openapi_http:
	oapi-codegen -generate types -o pkg/api/rest/openapi_types.gen.go -package rest api/openapi/crud.yml
	oapi-codegen -generate chi-server -o pkg/api/rest/openapi_api.gen.go -package rest api/openapi/crud.yml
