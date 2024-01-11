package db

import (
	"context"
	"github.com/StanislavIvanovQA/simple_bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser() (User, CreateUserParams, error) {
	hashedPassword, _ := util.HashPassword(util.RandomPassword())

	arg := CreateUserParams{
		Username:       util.RandomName(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	return user, arg, err
}

func TestCreateUser(t *testing.T) {
	user, params, err := createRandomUser()

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, params.Username, user.Username)
	require.Equal(t, params.HashedPassword, user.HashedPassword)
	require.Equal(t, params.FullName, user.FullName)
	require.Equal(t, params.Email, user.Email)

	require.NotZero(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
}

func TestGetUser(t *testing.T) {
	account1, _, err := createRandomUser()
	require.NoError(t, err)

	account2, err := testQueries.GetUser(context.Background(), account1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.FullName, account2.FullName)
	require.Equal(t, account1.Username, account2.Username)
	require.Equal(t, account1.Email, account2.Email)
	require.Equal(t, account1.HashedPassword, account2.HashedPassword)
	require.WithinDuration(t, account1.PasswordChangedAt, account2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}
