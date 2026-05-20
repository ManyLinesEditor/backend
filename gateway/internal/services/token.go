package services

import (
	"crypto/sha256"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwe"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type TokenService struct {
	signKey []byte
	encKey  []byte
	ttl     time.Duration
}

func NewTokenService(secret string, ttl time.Duration) *TokenService {
	raw := []byte(secret)
	s := sha256.Sum256(append([]byte("sign:"), raw...))
	e := sha256.Sum256(append([]byte("enc:"), raw...))
	return &TokenService{signKey: s[:], encKey: e[:], ttl: ttl}
}

func (s *TokenService) Issue(userID uuid.UUID) (string, error) {
	tok, err := jwt.NewBuilder().
		Subject(userID.String()).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(s.ttl)).
		Build()
	if err != nil {
		return "", err
	}
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.HS256, s.signKey))
	if err != nil {
		return "", err
	}
	encrypted, err := jwe.Encrypt(signed,
		jwe.WithKey(jwa.A256KW, s.encKey),
		jwe.WithContentEncryption(jwa.A256GCM),
	)
	return string(encrypted), err
}

func (s *TokenService) Verify(tokenStr string) (uuid.UUID, error) {
	signed, err := jwe.Decrypt([]byte(tokenStr), jwe.WithKey(jwa.A256KW, s.encKey))
	if err != nil {
		return uuid.Nil, err
	}
	tok, err := jwt.Parse(signed,
		jwt.WithKey(jwa.HS256, s.signKey),
		jwt.WithValidate(true),
	)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(tok.Subject())
}
