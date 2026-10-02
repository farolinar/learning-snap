package snap

import (
	"crypto/rand"
	"crypto/rsa"
	"strings"
	"testing"
)

const (
	testClientID     = "962489e9-de5d-4eb7-92a4-b07d44d64bf4"
	testClientSecret = "xS3vNQQgJRemFF0SZfXkZOq3r7kQ9n5YJgK4Wg0tVCQ="
	testTimestamp    = "2020-01-01T00:00:00+07:00"
)

func TestAccessTokenSignatureRoundTrip(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	sig, err := SignAccessToken(priv, testClientID, testTimestamp)
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyAccessToken(&priv.PublicKey, testClientID, testTimestamp, sig); err != nil {
		t.Fatalf("expected valid signature, got %v", err)
	}
	// A different timestamp must fail verification.
	if err := VerifyAccessToken(&priv.PublicKey, testClientID, "2020-01-01T00:00:01+07:00", sig); err == nil {
		t.Fatal("expected verification to fail for a different timestamp")
	}
}

func TestTransactionStringToSignFormat(t *testing.T) {
	body := []byte(`{"accountNo":"115471119"}`)
	bodyHash := SHA256HexLower(body)

	gotSymmetric := TransactionStringToSign("POST", "/api/v1.0/balance-inquiry", "TOKEN", body, testTimestamp, true)
	wantSymmetric := "POST:/api/v1.0/balance-inquiry:TOKEN:" + bodyHash + ":" + testTimestamp
	if gotSymmetric != wantSymmetric {
		t.Fatalf("symmetric stringToSign mismatch:\n got %q\nwant %q", gotSymmetric, wantSymmetric)
	}

	gotAsymmetric := TransactionStringToSign("POST", "/api/v1.0/balance-inquiry", "TOKEN", body, testTimestamp, false)
	wantAsymmetric := "POST:/api/v1.0/balance-inquiry:" + bodyHash + ":" + testTimestamp
	if gotAsymmetric != wantAsymmetric {
		t.Fatalf("asymmetric stringToSign mismatch:\n got %q\nwant %q", gotAsymmetric, wantAsymmetric)
	}
	if strings.Contains(gotAsymmetric, "TOKEN") {
		t.Fatal("asymmetric stringToSign must not contain the access token")
	}
}

func TestSymmetricTransactionRoundTrip(t *testing.T) {
	body := []byte(`{"accountNo":"115471119","balanceTypes":["Cash"]}`)
	sig, err := SignTransactionSymmetric(testClientSecret, "POST", "/api/v1.0/balance-inquiry", "TOKEN", body, testTimestamp)
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyTransactionSymmetric(testClientSecret, "POST", "/api/v1.0/balance-inquiry", "TOKEN", body, testTimestamp, sig); err != nil {
		t.Fatalf("expected valid signature, got %v", err)
	}
	if err := VerifyTransactionSymmetric("wrong-secret", "POST", "/api/v1.0/balance-inquiry", "TOKEN", body, testTimestamp, sig); err == nil {
		t.Fatal("expected verification to fail with the wrong secret")
	}
}

func TestAsymmetricTransactionRoundTrip(t *testing.T) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	body := []byte(`{"accountNo":"115471119"}`)
	sig, err := SignTransactionAsymmetric(priv, "POST", "/api/v1.0/balance-inquiry", body, testTimestamp)
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyTransactionAsymmetric(&priv.PublicKey, "POST", "/api/v1.0/balance-inquiry", body, testTimestamp, sig); err != nil {
		t.Fatalf("expected valid signature, got %v", err)
	}
}

func TestSHA256HexLowerIsLowercase(t *testing.T) {
	got := SHA256HexLower([]byte(""))
	const want = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if got != want {
		t.Fatalf("sha256 of empty string: got %q want %q", got, want)
	}
	if strings.ToLower(got) != got {
		t.Fatal("hash must be lowercase hex")
	}
}
