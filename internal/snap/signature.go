package snap

import (
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

// AccessTokenToString builds "clientId|timestamp"
func AccessTokenToString(clientId string, timestamp string) string {
	return clientId + "|" + timestamp
}

// SignAccessToken signs the access-token stringToSign with SHA256withRSA and returns lowecase hex
func SignAccessToken(privKey *rsa.PrivateKey, clientId, timestamp string) (string, error) {
	msg := []byte(AccessTokenToString(clientId, timestamp))
	digest := sha256.Sum256(msg)
	sig, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, digest[:])
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sig), nil
}

// VerifyAccessToken checks an access-token signature against the partner's public key.
func VerifyAccessToken(pub *rsa.PublicKey, clientId, timestamp, signatureHex string) error {
	sig, err := hex.DecodeString(signatureHex)
	if err != nil {
		return err
	}

	digest := sha256.Sum256([]byte(AccessTokenToString(clientId, timestamp)))

	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sig)
}

// TransactionStringToSign builds the canonical transaction string.
// withAccessToken selects between the symmetric (true) and asymmetric (false) variants.
func TransactionStringToSign(method, endpointURL, accessToken string, body []byte, timestamp string, withAccessToken bool) string {
	bodyHash := SHA256HexLower(body)
	if withAccessToken {
		return fmt.Sprintf("%s:%s:%s:%s:%s", method, endpointURL, accessToken, bodyHash, timestamp)
	}
	return fmt.Sprintf("%s:%s:%s:%s", method, endpointURL, bodyHash, timestamp)
}

// SignTransactionSymmetric computes HMAC_SHA512(clientSecret, stringToSign) as lowercase hex.
func SignTransactionSymmetric(clientSecret, method, endpointURL, accessToken string, body []byte, timestamp string) (string, error) {
	sts := TransactionStringToSign(method, endpointURL, accessToken, body, timestamp, true)
	mac := hmac.New(sha512.New, []byte(clientSecret))
	if _, err := mac.Write([]byte(sts)); err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// VerifyTransactionSymmetric recomputes the MAC and compares it in constant time.
func VerifyTransactionSymmetric(clientSecret, method, endpointURL, accessToken string, body []byte, timestamp, signatureHex string) error {
	expected, err := SignTransactionSymmetric(clientSecret, method, endpointURL, accessToken, body, timestamp)
	if err != nil {
		return err
	}

	// hmac.Equal is constant-time; never compare signatures with ==.
	if !hmac.Equal([]byte(expected), []byte(signatureHex)) {
		return fmt.Errorf("transaction signature mismatch")
	}
	return nil
}

// SignTransactionAsymmetric signs the asymmetric stringToSign (no access token) with the private key.
func SignTransactionAsymmetric(priv *rsa.PrivateKey, method, endpointURL string, body []byte, timestamp string) (string, error) {
	digest := sha256.Sum256([]byte(TransactionStringToSign(method, endpointURL, "", body, timestamp, false)))

	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, digest[:])
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig), nil
}

// VerifyTransactionAsymmetric checks an asymmetric transaction signature.
func VerifyTransactionAsymmetric(pub *rsa.PublicKey, method, endpointURL string, body []byte, timestamp, signatureHex string) error {
	sig, err := hex.DecodeString(signatureHex)
	if err != nil {
		return fmt.Errorf("signature is not valid hex: %w", err)
	}

	sts := TransactionStringToSign(method, endpointURL, "", body, timestamp, false)
	digest := sha256.Sum256([]byte(sts))

	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sig)
}
