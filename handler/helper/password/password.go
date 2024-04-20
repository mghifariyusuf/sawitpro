package password

import "golang.org/x/crypto/bcrypt"

func Validate(password, phoneNumber, hashedPassword string) error {
	// Concatenate the password and phone number to add uniqueness
	passwordWithSalt := password + phoneNumber

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(passwordWithSalt))
}

func Hash(password, phoneNumber string) (string, error) {
	// Concatenate the password and phone number to add uniqueness
	passwordWithSalt := password + phoneNumber

	// Generate salted hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
