package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTManager_Generate(t *testing.T) {
	manager := NewJWTManager("test-secret-key-min-32-characters", 24)

	token, err := manager.Generate("user-123", "test@example.com", "user")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTManager_Validate_Success(t *testing.T) {
	manager := NewJWTManager("test-secret-key-min-32-characters", 24)

	token, _ := manager.Generate("user-123", "test@example.com", "admin")

	claims, err := manager.Validate(token)

	assert.NoError(t, err)
	assert.Equal(t, "user-123", claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "admin", claims.Role)
}

func TestJWTManager_Validate_InvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key-min-32-characters", 24)

	claims, err := manager.Validate("invalid-token")

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestJWTManager_Validate_WrongSecret(t *testing.T) {
	manager1 := NewJWTManager("secret-key-one-min-32-characters", 24)
	manager2 := NewJWTManager("secret-key-two-min-32-characters", 24)

	token, _ := manager1.Generate("user-123", "test@example.com", "user")

	claims, err := manager2.Validate(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTManager_Validate_ExpiredToken(t *testing.T) {
	manager := &JWTManager{
		secret:      "test-secret-key-min-32-characters",
		expireHours: 0,
	}

	token, _ := manager.Generate("user-123", "test@example.com", "user")

	time.Sleep(time.Second * 2)

	claims, err := manager.Validate(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}