package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/quynhtruong/backend-master-class/util"
	"github.com/stretchr/testify/require"
)

func creatRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	creatRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := creatRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := creatRandomUser(t)
	newFullName := util.RandomOwner()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.Equal(t, updatedUser.FullName, newFullName)
	require.Equal(t, updatedUser.Username, oldUser.Username)
	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.Email, oldUser.Email)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := creatRandomUser(t)
	newEmail := util.RandomEmail()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.Equal(t, updatedUser.Email, newEmail)
	require.Equal(t, updatedUser.Username, oldUser.Username)
	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.FullName, oldUser.FullName)
}
