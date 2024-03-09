package util

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/vvbbnn00/goflet/config"
)

var (
	secretKey      string
	publicKey      string
	alg            string
	trustedIssuers []string
)

func init() {
	JwtInit()
}

func JwtInit() {
	conf := config.GofletCfg.JWTConfig
	alg = conf.Algorithm // The only supported algorithm defined in the configuration
	secretKey = conf.Security.SigningKey
	publicKey = conf.Security.PublicKey
	trustedIssuers = conf.TrustedIssuers
}

var ErrInvalidAlgorithm = errors.New("invalid algorithm")
var ErrUnsafeNoneAlgorithm = errors.New("none algorithm is not supported for security reasons")

// Permission The permission of the token
type Permission struct {
	Path    string            `json:"path"`    // The path that the token is allowed to access, supports wildcards
	Methods []string          `json:"methods"` // The methods that the token is allowed to access
	Query   map[string]string `json:"query"`   // The query parameters, if set in the map, the query should match the map
}

// JwtClaims The body of the JWT token
type JwtClaims struct {
	*jwt.StandardClaims
	Permissions []Permission `json:"permissions"` // The permissions of the token
}

func (c *JwtClaims) Valid() error {
	if c.StandardClaims == nil { // Check StandardClaims in case of nil pointer dereference
		return errors.New("missing required fields")
	}
	return c.StandardClaims.Valid()
}

// ParseJwtToken Parse the JWT token
func ParseJwtToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, selectSecretKey)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token") // Should make an error instance here.
	}

	// If there is no trusted issuer, trust any issuer
	if len(trustedIssuers) == 0 {
		return claims, nil
	}

	// Check if the issuer is trusted
	found := false
	issuer := claims.Issuer
	for _, trustedIssuer := range trustedIssuers {
		if issuer == trustedIssuer {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("untrusted issuer")
	}

	return claims, nil
}

// selectSecretKey The function to get the key for the JWT token
func selectSecretKey(token *jwt.Token) (interface{}, error) {
	if token.Method.Alg() != alg {
		return nil, ErrInvalidAlgorithm
	}

	switch token.Method.Alg() {
	case jwt.SigningMethodHS256.Name, jwt.SigningMethodHS384.Name, jwt.SigningMethodHS512.Name:
		return []byte(secretKey), nil
	case jwt.SigningMethodRS256.Name, jwt.SigningMethodRS384.Name, jwt.SigningMethodRS512.Name:
		return jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	case jwt.SigningMethodES256.Name, jwt.SigningMethodES384.Name, jwt.SigningMethodES512.Name:
		return jwt.ParseECPublicKeyFromPEM([]byte(publicKey))
	case jwt.SigningMethodPS256.Name, jwt.SigningMethodPS384.Name, jwt.SigningMethodPS512.Name:
		return jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	case jwt.SigningMethodNone.Alg():
		return nil, ErrUnsafeNoneAlgorithm
	default:
		return nil, ErrInvalidAlgorithm
	}
}
