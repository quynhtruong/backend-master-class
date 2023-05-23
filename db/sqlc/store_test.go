package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	store := NewStore(testDB)
	account1 := creatRandomAccount(t)
	account2 := creatRandomAccount(t)

	fmt.Println(">>before: ", account1.Balance, account2.Balance)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	n := 5
	amount := int64(10)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		//check account entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)

		// TODO: check account balance
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)
		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
	}

	//check thefinal updated balace
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">>after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

func TestTransferDeadlock(t *testing.T) {
	store := NewStore(testDB)
	account1 := creatRandomAccount(t)
	account2 := creatRandomAccount(t)

	fmt.Println(">>before: ", account1.Balance, account2.Balance)
	errs := make(chan error)

	n := 10
	amount := int64(10)
	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: fromAccountID,
				ToAccountId:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	//check thefinal updated balace
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	fmt.Println(">>after: ", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
