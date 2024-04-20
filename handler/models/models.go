package models

import (
	"regexp"
)

const (
	RegexUppercase            = "[A-Z]+"
	RegexNumber               = "[0-9]+"
	RegexSpecialChars         = "[^a-zA-Z0-9 ]+"
	RegexIndonesiaPhoneNumber = `^\+62[0-9]*$`
)

type ErrorResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

type RegisterUserRequest struct {
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
	Password    string `json:"password"`
}

func (registerRequest *RegisterUserRequest) Validate() map[string][]string {
	errs := make(map[string][]string)

	if len(registerRequest.PhoneNumber) < 10 ||
		len(registerRequest.PhoneNumber) > 13 {
		errs["phone_number"] = append(errs["phone_number"], "Phone number must be at minimum 10 characters and maximum 13 characters")
	}

	validCountryCode := regexp.MustCompile(RegexIndonesiaPhoneNumber).MatchString(registerRequest.PhoneNumber)
	if !validCountryCode {
		errs["phone_number"] = append(errs["phone_number"], "Phone number must start with the Indonesia country code")
	}

	if len(registerRequest.FullName) < 3 ||
		len(registerRequest.FullName) > 60 {
		errs["full_name"] = append(errs["full_name"], "Full name must be at minimum 3 characters and maximum 60 characters")
	}

	if len(registerRequest.Password) < 6 || len(registerRequest.Password) > 64 {
		errs["password"] = append(errs["password"], "Password must be minimum 6 characters and maximum 64 characters")
	}

	passwordContainUppercase := regexp.MustCompile(RegexUppercase).MatchString(registerRequest.Password)
	passwordContainNumber := regexp.MustCompile(RegexNumber).MatchString(registerRequest.Password)
	passwordContainSpecialChars := regexp.MustCompile(RegexSpecialChars).MatchString(registerRequest.Password)
	if !passwordContainUppercase || !passwordContainNumber || !passwordContainSpecialChars {
		errs["password"] = append(errs["password"], "Password must contain at least 1 capital letter, 1 number, and 1 special characters")
	}

	return errs
}

type RegisterUserResponse struct {
	ID int64 `json:"id"`
}

type LoginUserRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginUserResponse struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}

type GetUserProfileResponse struct {
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
}

type UpdateUserProfileRequest struct {
	PhoneNumber *string `json:"phone_number,omitempty"`
	FullName    *string `json:"full_name,omitempty"`
}

func (updateRequest *UpdateUserProfileRequest) Validate() map[string][]string {
	errs := make(map[string][]string)

	if updateRequest.PhoneNumber != nil {
		if len(*updateRequest.PhoneNumber) < 10 ||
			len(*updateRequest.PhoneNumber) > 13 {
			errs["phone_number"] = append(errs["phone_number"], "Phone number must be at minimum 10 characters and maximum 13 characters")
		}

		validCountryCode := regexp.MustCompile(RegexIndonesiaPhoneNumber).MatchString(*updateRequest.PhoneNumber)
		if !validCountryCode {
			errs["phone_number"] = append(errs["phone_number"], "Phone number must start with the Indonesia country code")
		}
	}

	if updateRequest.FullName != nil {
		if len(*updateRequest.FullName) < 3 ||
			len(*updateRequest.FullName) > 60 {
			errs["full_name"] = append(errs["full_name"], "Full name must be at minimum 3 characters and maximum 60 characters")
		}
	}

	return errs
}

type UpdateUserProfileResponse struct {
	PhoneNumber *string `json:"phone_number,omitempty"`
	FullName    *string `json:"full_name,omitempty"`
}
