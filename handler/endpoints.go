package handler

import (
	"net/http"

	"github.com/SawitProRecruitment/UserService/handler/models"
	"github.com/SawitProRecruitment/UserService/repository/entity"
	"github.com/labstack/echo/v4"
)

func (server *Server) RegisterUser(c echo.Context) error {
	ctx := c.Request().Context()
	registerRequest := &models.RegisterUserRequest{}

	err := c.Bind(registerRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Failed to read request",
			Error:   err.Error(),
		})
	}

	errs := registerRequest.Validate()
	if len(errs) > 0 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request",
			Error:   errs,
		})
	}

	hashedPassword, err := HashPassword(registerRequest.Password, registerRequest.PhoneNumber)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Failed to generate hashed password",
			Error:   err.Error(),
		})
	}
	registerRequest.Password = hashedPassword

	id, err := server.Repository.CreateUser(ctx, registerRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Failed to register user",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.RegisterUserResponse{
		ID: id,
	})
}

func (server *Server) LoginUser(c echo.Context) error {
	ctx := c.Request().Context()
	loginRequest := &models.LoginUserRequest{}

	err := c.Bind(loginRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Failed to read request",
			Error:   err.Error(),
		})
	}

	user, err := server.Repository.GetUser(ctx, &entity.UserFilter{
		PhoneNumber: &loginRequest.PhoneNumber,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Failed to get user",
			Error:   err.Error(),
		})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Message: "User not found",
		})
	}

	// Compare password from request and db
	err = ValidatePassword(loginRequest.Password, user.PhoneNumber, user.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Authentication failed",
			Error:   err.Error(),
		})
	}

	// Generate JWT token
	token, err := server.JWT.GenerateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
	}

	// Increase successful login number
	err = server.Repository.IncLogin(ctx, user.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Failed to login",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.LoginUserResponse{
		ID:    user.ID,
		Token: token,
	})
}

func (server *Server) GetUserProfile(c echo.Context) error {
	ctx := c.Request().Context()

	// Extract JWT token from Authorization header
	headerAuthorization := c.Request().Header.Get("Authorization")
	if headerAuthorization == "" {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Message: "Missing authorization token",
		})
	}

	// Parse JWT to get token claims
	claims, err := server.JWT.ValidateToken(headerAuthorization)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Message: "Failed to validate token",
			Error:   err.Error(),
		})
	}

	// Extract user ID from the token claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Message: "Invalid user ID in token claims",
		})
	}
	id := int64(userID)

	user, err := server.Repository.GetUser(ctx, &entity.UserFilter{
		ID: &id,
	})
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Message: "Failed to get user",
			Error:   err.Error(),
		})
	}
	if user == nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Message: "User not found",
		})
	}

	return c.JSON(http.StatusOK, models.GetUserProfileResponse{
		PhoneNumber: user.PhoneNumber,
		FullName:    user.FullName,
	})
}

func (server *Server) UpdateUserProfile(c echo.Context) error {
	ctx := c.Request().Context()

	// Extract JWT token from Authorization header
	headerAuthorization := c.Request().Header.Get("Authorization")
	if headerAuthorization == "" {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Message: "Missing authorization token",
		})
	}

	// Parse JWT to get token claims
	claims, err := server.JWT.ValidateToken(headerAuthorization)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Message: "Failed to validate token",
			Error:   err.Error(),
		})
	}

	// Extract user ID from the token claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Message: "Invalid user ID in token claims",
		})
	}
	id := int64(userID)

	updateRequest := &models.UpdateUserProfileRequest{}

	err = c.Bind(updateRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Failed to read request",
			Error:   err.Error(),
		})
	}

	errs := updateRequest.Validate()
	if len(errs) > 0 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request",
			Error:   errs,
		})
	}

	// Check phone number
	if updateRequest.PhoneNumber != nil {
		user, err := server.Repository.GetUser(ctx, &entity.UserFilter{
			PhoneNumber: updateRequest.PhoneNumber,
		})
		if err != nil {
			return c.JSON(http.StatusForbidden, models.ErrorResponse{
				Message: "Failed to get user",
				Error:   err.Error(),
			})
		}

		if user != nil {
			return c.JSON(http.StatusConflict, models.ErrorResponse{
				Message: "Phone number already registered",
			})
		}
	}

	err = server.Repository.UpdateProfile(ctx, id, updateRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Failed to update user",
			Error:   err,
		})
	}

	return c.JSON(http.StatusOK, models.UpdateUserProfileResponse{
		PhoneNumber: updateRequest.PhoneNumber,
		FullName:    updateRequest.FullName,
	})
}
