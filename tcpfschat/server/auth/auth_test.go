package auth

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuth(t *testing.T) {
	a := NewAuthenticator()

	username := "test"
	pass := "password"
	err := a.Auth(username, pass)
	assert.NoError(t, err)

	err = a.Auth(username, "sfsdfgsdfg")
	assert.Equal(t, errors.New("wrong password"), err)

	err = a.Auth(username, pass)
	assert.NoError(t, err)
}
