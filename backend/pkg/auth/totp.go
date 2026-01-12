package auth

import (
	"encoding/base64"
	"fmt"

	"github.com/pquerna/otp/totp"
)

func GenerateTOTPSecret(email, issuer string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: email,
	})
	if err != nil {
		return "", "", err
	}

	// Convert to PNG QR code
	var buf []byte
	_, err := key.Image(200, 200)
	if err != nil {
		return "", "", err
	}

	// Convert image to base64
	qrCode := base64.StdEncoding.EncodeToString(buf)

	return key.Secret(), fmt.Sprintf("data:image/png;base64,%s", qrCode), nil
}

func ValidateTOTP(code, secret string) bool {
	return totp.Validate(code, secret)
}

func GenerateTOTPURL(secret, email, issuer string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s",
		issuer, email, secret, issuer)
}
