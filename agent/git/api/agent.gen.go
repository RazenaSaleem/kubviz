// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

// PostAzureJSONBody defines parameters for PostAzure.
type PostAzureJSONBody = map[string]interface{}

// PostBitbucketJSONBody defines parameters for PostBitbucket.
type PostBitbucketJSONBody = map[string]interface{}

// PostGiteaJSONBody defines parameters for PostGitea.
type PostGiteaJSONBody = map[string]interface{}

// PostGithubJSONBody defines parameters for PostGithub.
type PostGithubJSONBody = map[string]interface{}

// PostGitlabJSONBody defines parameters for PostGitlab.
type PostGitlabJSONBody = map[string]interface{}

// PostAzureJSONRequestBody defines body for PostAzure for application/json ContentType.
type PostAzureJSONRequestBody = PostAzureJSONBody

// PostBitbucketJSONRequestBody defines body for PostBitbucket for application/json ContentType.
type PostBitbucketJSONRequestBody = PostBitbucketJSONBody

// PostGiteaJSONRequestBody defines body for PostGitea for application/json ContentType.
type PostGiteaJSONRequestBody = PostGiteaJSONBody

// PostGithubJSONRequestBody defines body for PostGithub for application/json ContentType.
type PostGithubJSONRequestBody = PostGithubJSONBody

// PostGitlabJSONRequestBody defines body for PostGitlab for application/json ContentType.
type PostGitlabJSONRequestBody = PostGitlabJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List of APIs provided by the service
	// (GET /api-docs)
	GetApiDocs(c *gin.Context)
	// Handle Azure webhooks post requests
	// (POST /azure)
	PostAzure(c *gin.Context)
	// Handle Bitbucket webhooks post requests
	// (POST /bitbucket)
	PostBitbucket(c *gin.Context)
	// Handle Gitea webhooks post requests
	// (POST /gitea)
	PostGitea(c *gin.Context)
	// Handle Github webhooks post requests
	// (POST /github)
	PostGithub(c *gin.Context)
	// Handle Gitlab webhooks post requests
	// (POST /gitlab)
	PostGitlab(c *gin.Context)
	// Kubernetes readiness and liveness probe endpoint
	// (GET /liveness)
	GetLiveness(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetApiDocs operation middleware
func (siw *ServerInterfaceWrapper) GetApiDocs(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetApiDocs(c)
}

// PostAzure operation middleware
func (siw *ServerInterfaceWrapper) PostAzure(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostAzure(c)
}

// PostBitbucket operation middleware
func (siw *ServerInterfaceWrapper) PostBitbucket(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostBitbucket(c)
}

// PostGitea operation middleware
func (siw *ServerInterfaceWrapper) PostGitea(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostGitea(c)
}

// PostGithub operation middleware
func (siw *ServerInterfaceWrapper) PostGithub(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostGithub(c)
}

// PostGitlab operation middleware
func (siw *ServerInterfaceWrapper) PostGitlab(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.PostGitlab(c)
}

// GetLiveness operation middleware
func (siw *ServerInterfaceWrapper) GetLiveness(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
	}

	siw.Handler.GetLiveness(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router *gin.Engine, si ServerInterface) *gin.Engine {
	return RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router *gin.Engine, si ServerInterface, options GinServerOptions) *gin.Engine {

	errorHandler := options.ErrorHandler

	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/api-docs", wrapper.GetApiDocs)

	router.POST(options.BaseURL+"/azure", wrapper.PostAzure)

	router.POST(options.BaseURL+"/bitbucket", wrapper.PostBitbucket)

	router.POST(options.BaseURL+"/gitea", wrapper.PostGitea)

	router.POST(options.BaseURL+"/github", wrapper.PostGithub)

	router.POST(options.BaseURL+"/gitlab", wrapper.PostGitlab)

	router.GET(options.BaseURL+"/liveness", wrapper.GetLiveness)

	return router
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8yUP4/bPAyHv4rA2W+Se7t5y6HANbgDmq1DcYMs0THvHFEVpRRp4O9eUNekQP/ES4dM",
	"toAfqeehAJ7A8T5ywJAF2tPUAIWeoT2BR3GJYiYO0MI6mPV2Y3pOZrDBjxR25it2A/OrmD7x3njqe0wY",
	"skkYWShzIhRoIFMeEVr4dI6vtxto4IBJ3nrfLVaLFUwNcMRgI0EL7xarxR00EG0elAuWNtJ/nl097DDr",
	"hyMmq3wbDy08YF5Heq+RBhJK5CBY4/+vVr8LfXyEaWpAyn5v0xFaeCLJhnulExMTH8ijN93R5AGNYDqQ",
	"Q7WxO4H2M8TSjeTgWZss7beSUO+ILH9A27LkdY0o2ZeCku/ZHzXoOGQMtcbGOJKrVcsXUcgTiBtwb/Uv",
	"H6MOkbsXdBmmSe+dt5TiHIr0ZTQXpF+8P+hzoql8P59URcwPVqkVy45yV9zr2+z/Lnp/id2u7IXxqvCO",
	"Mtrrsg81cruilW9OcijdrKVmblpzKN2c52jnPTVz056jve450gEDytU9+XTO/CO6x9JhCphRTELrSXsb",
	"G7w5w+g+7dBg8JEp1Mk0oDsVky7TE5Q0QgtLmJ6n7wEAAP//s7retZIGAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
