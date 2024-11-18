package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GenerateXSignature(timestamp, bodyRequest, apiToken, privateKey string) (string, error) {
	// Concatenate the payload
	var payload string

	if bodyRequest == "" {
		payload = fmt.Sprintf("%s%s%s", timestamp, bodyRequest, apiToken)
	} else {
		payload = fmt.Sprintf("%s%s", timestamp, apiToken)
	}

	// Create HMAC-SHA256 hash
	h := hmac.New(sha256.New, []byte(privateKey))
	_, err := h.Write([]byte(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create HMAC: %v", err)
	}

	// Encode the HMAC hash to hexadecimal
	signature := hex.EncodeToString(h.Sum(nil))
	return signature, nil
}

func ValidateXSignature(receivedSignature, timestamp, bodyRequest, apiToken, privateKey string) bool {
	// Generate signature dari data request
	expectedSignature, err := GenerateXSignature(timestamp, bodyRequest, apiToken, privateKey)
	if err != nil {
		return false
	}

	fmt.Println("Expected Signature:", expectedSignature)
	fmt.Println("Received Signature:", receivedSignature)
	return hmac.Equal([]byte(receivedSignature), []byte(expectedSignature))
}
