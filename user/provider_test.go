package user

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"main/model"
	"testing"
)

func TestProvider_GerUserByName(t *testing.T) {
	provider, err := NewProvider("test_users.json")
	assert.NoError(t, err)

	usr := provider.GerUserByName("Michael")
	assert.Nil(t, usr)

	id, _ := uuid.Parse("6f2d2ec6-a89a-11ed-afa1-0242ac120002")

	expectedBob := &model.User{
		Name:           "Bob",
		ID:             id,
		HashedPassword: "$2a$14$tSnuH38X0AZ3J2cyS.OYAubMYBeXtCzexI/M1lCYx7rlbNkEq.frC",
	}
	bob := provider.GerUserByName("Bob")

	assert.Equal(t, expectedBob, bob)
}
