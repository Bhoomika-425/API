package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"project/internal/middleware"
	mock_files "project/internal/mock-files"
	"project/internal/models"

	service "project/internal/service"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func Test_handler_Login(t *testing.T) {
	// type args struct {
	// 	c *gin.Context
	// }
	tests := []struct {
		name string
		// h    *handler
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
			name: "error in validating json",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", strings.NewReader(`{"hii"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"please provide valid email and password"}`,
		},
		{
			name: "validating",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", strings.NewReader(`{"email":"bhoomi@gmail.com","password":""}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"please provide valid  email and password"}`,
		},
		{
			name: "failure",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", strings.NewReader(`{"username":"asdfghjk","email":"bhoomi@gmail.com","password":"12345678"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				mc := gomock.NewController(t)
				mu := mock_files.NewMockUserService(mc)

				mu.EXPECT().UserLogin(gomock.Any(), gomock.Any()).Return("", errors.New("error")).AnyTimes()

				return c, rr, mu
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Bad Request"}`,
		},
		{
			name: "success",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", strings.NewReader(`{"username":"asdfghjk","email":"bhoomi@gmail.com","password":"12345678"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				mc := gomock.NewController(t)
				mu := mock_files.NewMockUserService(mc)

				mu.EXPECT().UserLogin(gomock.Any(), gomock.Any()).Return("", nil).AnyTimes()

				return c, rr, mu
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"token":""}`,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()
			h := &handler{
				service: ms,
			}
			h.Login(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())

			// tt.h.Login(tt.args.c)
		})
	}
}

func Test_handler_SignUp(t *testing.T) {
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
			name: "error in validating json",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", strings.NewReader(`{"hii"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"please provide valid username, email and password"}`,
		},
		{
			name: "validating",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", strings.NewReader(`{"username":"anii","email":"bhoomi@gmail.com","password":""}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest
				return c, rr, nil
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"please provide valid username, email and password"}`,
		},
		{
			name: "failure",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", strings.NewReader(`{"username":"asdfghjk","email":"bhoomi@gmail.com","password":"12345678"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				mc := gomock.NewController(t)
				mu := mock_files.NewMockUserService(mc)

				mu.EXPECT().UserSignup(gomock.Any(), gomock.Any()).Return(models.User{}, errors.New("error")).AnyTimes()

				return c, rr, mu
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"error"}`,
		},
		{
			name: "success",
			setup: func() (*gin.Context, *httptest.ResponseRecorder, service.UserService) {
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				httpRequest, _ := http.NewRequest(http.MethodGet, "http://test.com", strings.NewReader(`{"username":"ani","email":"bhoomi@gmail.com","password":"1234"}`))
				ctx := httpRequest.Context()
				ctx = context.WithValue(ctx, middleware.TraceIDKey, "1")
				httpRequest = httpRequest.WithContext(ctx)
				c.Request = httpRequest

				mc := gomock.NewController(t)
				mu := mock_files.NewMockUserService(mc)

				mu.EXPECT().UserSignup(gomock.Any(), gomock.Any()).Return(models.User{}, nil).AnyTimes()

				return c, rr, mu
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"username":"","email":""}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, rr, ms := tt.setup()
			h := &handler{
				service: ms,
			}
			h.SignUp(c)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())

		})
	}
}
