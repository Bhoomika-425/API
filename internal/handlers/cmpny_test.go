package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"project/internal/auth"
	"project/internal/middleware"
	mock_files "project/internal/mock-files"
	"project/internal/models"
	service "project/internal/service"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
)

func Test_handler_ViewCompany(t *testing.T) {
	// type args struct {
	// 	c *gin.Context
	// }
	tests := []struct {
		name               string
		setup              func() (*gin.Context, *httptest.ResponseRecorder, service.UserService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "missing trace id",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
		{
			name: "jwtclamis not there",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://read.com", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "456")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				return c, rr, nil
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Unauthorized"}`,
		},
		{
			name: "id not found",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://read.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "456")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Params = append(c.Params, gin.Param{Key: "id", Value: "abc"})
				c.Request = httpRequest

				mc := gomock.NewController(t)
				ms := mock_files.NewMockUserService(mc)

				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Bad Request"}`,
		},
		{
			name: "no companies to view",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://read.com:8080", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "456")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Params = append(c.Params, gin.Param{Key: "id", Value: "123"})
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := mock_files.NewMockUserService(mc)

				ms.EXPECT().ViewCompanyDetails(c.Request.Context(), gomock.Any()).Return(models.Company{}, errors.New("error while fetching")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"error while fetching"}`,
		},
		{
			name: " companies to view",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://read.com:8080", strings.NewReader(`{"name":"accen,"location":"bng","field":"sw"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "456")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Params = append(c.Params, gin.Param{Key: "id", Value: "123"})
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := mock_files.NewMockUserService(mc)

				ms.EXPECT().ViewCompanyDetails(c.Request.Context(), gomock.Any()).Return(models.Company{}, nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"name":"","location":"","field":""}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()
			h := &handler{
				service: ms,
			}
			h.ViewCompany(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
			// tt.h.ViewCompany(tt.args.c)
		})
	}
}

func Test_handler_ViewAllCompanies(t *testing.T) {
	// type args struct {
	// 	c *gin.Context
	// }
	tests := []struct {
		name               string
		setup              func() (*gin.Context, *httptest.ResponseRecorder, service.UserService)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "missing trace id",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
				c.Request = httpRequest

				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
		{
			name: "jwtclamis not there",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://read.com", nil)
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "456")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				return c, rr, nil
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `{"error":"Unauthorized"}`,
		},

		{
			name: "no companies to view",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://read.com:8080", strings.NewReader(
					`{"name":"accen,"location":"bng","field":"sw"}`,
				))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "456")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Params = append(c.Params, gin.Param{Key: "id", Value: "123"})
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := mock_files.NewMockUserService(mc)

				ms.EXPECT().ViewAllCompanies(c.Request.Context()).Return([]models.Company{}, errors.New("error while fetching")).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"error while fetching"}`,
		},
		{
			name: "companies to view success",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://read.com:8080", strings.NewReader(
					`{"name":"accen,"location":"bng","field":"sw"}`,
				))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "456")
				ctx = context.WithValue(ctx, auth.Key, jwt.RegisteredClaims{})
				httpRequest = httpRequest.WithContext(ctx)
				c.Params = append(c.Params, gin.Param{Key: "id", Value: "123"})
				c.Request = httpRequest
				mc := gomock.NewController(t)
				ms := mock_files.NewMockUserService(mc)

				ms.EXPECT().ViewAllCompanies(c.Request.Context()).Return([]models.Company{}, nil).AnyTimes()

				return c, rr, ms
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()
			h := &handler{
				service: ms,
			}
			h.ViewAllCompanies(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
			// tt.h.ViewAllCompanies(tt.args.c)
		})
	}
}

func Test_handler_AddCompany(t *testing.T) {
	// type args struct {
	// 	c *gin.Context
	// }
	tests := []struct {
		name               string
		setup              func() (*gin.Context, *httptest.ResponseRecorder, service.UserService)
		expectedStatusCode int
		expectedResponse   string
	}{
		// {
		// 	name: "missing trace id",
		// 	setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
		// 		rr := httptest.NewRecorder()
		// 		c, _ := gin.CreateTestContext(rr)
		// 		httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
		// 		c.Request = httpRequest

		// 		return c, rr, nil
		// 	},
		// 	expectedStatusCode: http.StatusInternalServerError,
		// 	expectedResponse:   `{"error":"Internal Server Error"}`,
		// },
		// {
		// 	name: "jwtclamis not there",
		// 	setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
		// 		rr := httptest.NewRecorder()
		// 		c, _ := gin.CreateTestContext(rr)
		// 		httpRequest, _ := http.NewRequest(http.MethodGet, "http://read.com", nil)
		// 		ctx := httpRequest.Context()
		// 		ctx = context.WithValue(ctx, middleware.TraceIDKey, "456")
		// 		httpRequest = httpRequest.WithContext(ctx)
		// 		c.Request = httpRequest
		// 		return c, rr, nil
		// 	},
		// 	expectedStatusCode: http.StatusUnauthorized,
		// 	expectedResponse:   `{"error":"Unauthorized"}`,
		// },
		{
			name: "error in validating json",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", strings.NewReader(`{"name":"info","location":"bng","field":"sw}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "1")

				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				return c, rr, nil
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"msg":"Internal Server Error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()
			h := &handler{
				service: ms,
			}
			h.AddCompany(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())
		})
	}
}
