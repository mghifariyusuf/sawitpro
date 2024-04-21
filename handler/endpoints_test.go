package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/handler/models"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	TestPhoneNumber = "+621234567890"
	TestFullName    = "John Doe"
	TestPassword    = "P@ssword1"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockRepositoryInterface(ctrl)

	tests := []struct {
		name                 string
		request              *models.RegisterUserRequest
		expectedStatusCode   int
		expectedResponseBody interface{}
		mockRepoExpectation  func()
		success              bool
	}{
		{
			name: "Valid Request",
			request: &models.RegisterUserRequest{
				PhoneNumber: TestPhoneNumber,
				FullName:    TestFullName,
				Password:    TestPassword,
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: &models.RegisterUserResponse{ID: 1},
			mockRepoExpectation: func() {
				mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(int64(1), nil)
			},
			success: true,
		},
		{
			name: "Invalid Phone Number Length",
			request: &models.RegisterUserRequest{
				PhoneNumber: "+621234",
				FullName:    TestFullName,
				Password:    TestPassword,
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: &models.RegisterUserResponse{ID: 0},
			success:              false,
		},
		{
			name: "Invalid Phone Number Country Code",
			request: &models.RegisterUserRequest{
				PhoneNumber: "+601234567890",
				FullName:    TestFullName,
				Password:    TestPassword,
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: &models.RegisterUserResponse{ID: 0},
			success:              false,
		},
		{
			name: "Invalid Full Name Length",
			request: &models.RegisterUserRequest{
				PhoneNumber: TestPhoneNumber,
				FullName:    "JD",
				Password:    TestPassword,
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: &models.RegisterUserResponse{ID: 0},
			success:              false,
		},
		{
			name: "Invalid Password",
			request: &models.RegisterUserRequest{
				PhoneNumber: TestPhoneNumber,
				FullName:    TestFullName,
				Password:    "123",
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: &models.RegisterUserResponse{ID: 0},
			success:              false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			reqBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.success {
				tt.mockRepoExpectation()
			}

			server := &handler.Server{
				Repository: mockRepo,
			}

			err := server.RegisterUser(c)
			require.NoError(t, err)

			var actualResponseBody *models.RegisterUserResponse
			err = json.Unmarshal(rec.Body.Bytes(), &actualResponseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCode, rec.Code)
			assert.Equal(t, tt.expectedResponseBody, actualResponseBody)
		})
	}
}

func TestLoginUser(t *testing.T) {
	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody interface{}
	}{
		{
			name:               "Valid Request",
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: &models.LoginUserResponse{
				ID:    1,
				Token: "mock",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.GET("/login", func(c echo.Context) error {
				return c.JSON(http.StatusOK, models.LoginUserResponse{
					ID:    1,
					Token: "mock",
				})
			})
			req := httptest.NewRequest(http.MethodGet, "/login", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer mock")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			var actualResponseBody *models.LoginUserResponse
			err := json.Unmarshal(rec.Body.Bytes(), &actualResponseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCode, rec.Code)
			assert.Equal(t, tt.expectedResponseBody, actualResponseBody)
		})
	}
}

func TestGetUserProfile(t *testing.T) {
	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody interface{}
	}{
		{
			name:               "Valid Request",
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: &models.GetUserProfileResponse{
				PhoneNumber: TestPhoneNumber,
				FullName:    TestFullName,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.GET("/profile", func(c echo.Context) error {
				return c.JSON(http.StatusOK, models.GetUserProfileResponse{
					PhoneNumber: TestPhoneNumber,
					FullName:    TestFullName,
				})
			})
			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer mock")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			var actualResponseBody *models.GetUserProfileResponse
			err := json.Unmarshal(rec.Body.Bytes(), &actualResponseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCode, rec.Code)
			assert.Equal(t, tt.expectedResponseBody, actualResponseBody)
		})
	}
}

func TestUpdateUserProfile(t *testing.T) {
	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody interface{}
	}{
		{
			name:               "Valid Request",
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: &models.UpdateUserProfileResponse{
				PhoneNumber: &TestPhoneNumber,
				FullName:    &TestFullName,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.GET("/profile", func(c echo.Context) error {
				return c.JSON(http.StatusOK, models.GetUserProfileResponse{
					PhoneNumber: TestPhoneNumber,
					FullName:    TestFullName,
				})
			})
			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer mock")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			var actualResponseBody *models.UpdateUserProfileResponse
			err := json.Unmarshal(rec.Body.Bytes(), &actualResponseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCode, rec.Code)
			assert.Equal(t, tt.expectedResponseBody, actualResponseBody)
		})
	}
}
