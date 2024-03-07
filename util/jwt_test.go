package util

import (
	"goflet/config"
	"testing"
)

func TestNone(t *testing.T) {
	tokenString := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ."

	config.GofletCfg.JWTConfig.Algorithm = "none"
	JwtInit()

	_, err := ParseJwtToken(tokenString)
	if err == nil {
		t.Fatal("None Algorithm should not be supported")
	}
}

func doTest(t *testing.T, tokenString string) {
	_, err := ParseJwtToken(tokenString)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHS256(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	key := "your-256-bit-secret"

	config.GofletCfg.JWTConfig.Algorithm = "HS256"
	config.GofletCfg.JWTConfig.Security.SigningKey = key
	JwtInit()

	doTest(t, tokenString)
}

func TestRS256(t *testing.T) {
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ"
	publicKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1SU1LfVLPHCozMxH2Mo
4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0/IzW7yWR7QkrmBL7jTKEn5u
+qKhbwKfBstIs+bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyeh
kd3qqGElvW/VDL5AaWTg0nLVkjRo9z+40RQzuVaE8AkAFmxZzow3x+VJYKdjykkJ
0iT9wCS0DRTXu269V264Vf/3jvredZiKRkgwlL9xNAwxXFg0x/XFw005UWVRIkdg
cKWTjpBP2dPwVZ4WWC+9aGVd+Gyn1o0CLelf4rEjGoXbAAEgAqeGUxrcIlbjXfbc
mwIDAQAB
-----END PUBLIC KEY-----`

	config.GofletCfg.JWTConfig.Algorithm = "RS256"
	config.GofletCfg.JWTConfig.Security.PublicKey = publicKey
	JwtInit()

	doTest(t, tokenString)
}

func TestES256(t *testing.T) {
	tokenString := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.tyh-VfuzIxCyGYDlkBA7DfyjrqmSHu6pQ2hoZuFqUSLPNY2N0mpHb3nk5K17HWP_3cYHBw7AhHale5wky6-sVA"
	publicKey := `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEEVs/o5+uQbTjL3chynL4wXgUg2R9
q9UU8I5mEovUf86QZ7kOBIjJwqnzD1omageEHWwHdBO6B+dFabmdT9POxg==
-----END PUBLIC KEY-----`

	config.GofletCfg.JWTConfig.Algorithm = "ES256"
	config.GofletCfg.JWTConfig.Security.PublicKey = publicKey
	JwtInit()

	doTest(t, tokenString)
}

func TestPS256(t *testing.T) {
	tokenString := "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.iOeNU4dAFFeBwNj6qdhdvm-IvDQrTa6R22lQVJVuWJxorJfeQww5Nwsra0PjaOYhAMj9jNMO5YLmud8U7iQ5gJK2zYyepeSuXhfSi8yjFZfRiSkelqSkU19I-Ja8aQBDbqXf2SAWA8mHF8VS3F08rgEaLCyv98fLLH4vSvsJGf6ueZSLKDVXz24rZRXGWtYYk_OYYTVgR1cg0BLCsuCvqZvHleImJKiWmtS0-CymMO4MMjCy_FIl6I56NqLE9C87tUVpo1mT-kbg5cHDD8I7MjCW5Iii5dethB4Vid3mZ6emKjVYgXrtkOQ-JyGMh6fnQxEFN1ft33GX2eRHluK9eg"
	publicKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1SU1LfVLPHCozMxH2Mo
4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0/IzW7yWR7QkrmBL7jTKEn5u
+qKhbwKfBstIs+bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyeh
kd3qqGElvW/VDL5AaWTg0nLVkjRo9z+40RQzuVaE8AkAFmxZzow3x+VJYKdjykkJ
0iT9wCS0DRTXu269V264Vf/3jvredZiKRkgwlL9xNAwxXFg0x/XFw005UWVRIkdg
cKWTjpBP2dPwVZ4WWC+9aGVd+Gyn1o0CLelf4rEjGoXbAAEgAqeGUxrcIlbjXfbc
mwIDAQAB
-----END PUBLIC KEY-----`

	config.GofletCfg.JWTConfig.Algorithm = "PS256"
	config.GofletCfg.JWTConfig.Security.PublicKey = publicKey
	JwtInit()

	doTest(t, tokenString)
}

func prepareForHS256() {
	config.GofletCfg.JWTConfig.Algorithm = "HS256"
	config.GofletCfg.JWTConfig.Security.SigningKey = "your-256-bit-secret"
	JwtInit()
}

func TestInvalidAlgorithm(t *testing.T) {
	prepareForHS256()
	tokenString := "eyJhbGciOiJIUzExNDUxNCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

	_, err := ParseJwtToken(tokenString)
	if err == nil {
		t.Fatal("Invalid algorithm should not be supported")
	}
}

func TestInvalidToken(t *testing.T) {
	prepareForHS256()
	tokenString := "123456"

	_, err := ParseJwtToken(tokenString)
	if err == nil {
		t.Fatal("Invalid token should not be supported")
	}
}

func TestInvalidKey(t *testing.T) {
	prepareForHS256()
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.EbVePZ7UuIHAYoyyH5KNBXVMnezJl8ut9Scx5XA42vc"

	_, err := ParseJwtToken(tokenString)
	if err == nil {
		t.Fatal("Invalid key should not be supported")
	}
}

func TestBeforeNbf(t *testing.T) {
	prepareForHS256()
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwibmJmIjo5OTk5OTk5OTk5fQ.HCxJ2E1Km6BM5EHMptURfFGqDLh4BbymYSfem-mdqvo"

	_, err := ParseJwtToken(tokenString)
	if err == nil {
		t.Fatal("Before nbf should not be supported")
	}
}

func TestAfterExp(t *testing.T) {
	prepareForHS256()
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZXhwIjoxMTQ1MTQxOTE5fQ.unr6-gTptTt4HskE_NWK4vvd8jZqFZU6SMmScOuTFt4"

	_, err := ParseJwtToken(tokenString)
	if err == nil {
		t.Fatal("After exp should not be supported")
	}
}

func TestValidIssuer(t *testing.T) {
	config.GofletCfg.JWTConfig.TrustedIssuers = []string{"issuer"}
	prepareForHS256()

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaXNzIjoiaXNzdWVyIn0.RY_M5zJMCCj1r4u6KV9fK3NCY5ubIctKON9fhFG63K8"

	_, err := ParseJwtToken(tokenString)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInvalidIssuer(t *testing.T) {
	config.GofletCfg.JWTConfig.TrustedIssuers = []string{"issuer"}
	prepareForHS256()

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaXNzIjoiVGFkb2tvcm8ifQ.GCRVtomcBEahfueDXIIn51U8rnwI86VgIQadQz2Od1c"

	_, err := ParseJwtToken(tokenString)
	if err == nil {
		t.Fatal("Invalid issuer should not be supported")
	}
}
