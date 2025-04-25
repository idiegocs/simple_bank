package db

import (
	"context"
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require" //go mod tidy -e
)


func CreateAccountRandom(t *testing.T) Account{

	arg:= CreateAccountParams{
		Owner:    util.RandomOwner(), //random generate
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),

	}

	account, err := testQueries.CreateAccount(context.Background(), arg)


	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {

	CreateAccountRandom(t)


}

func TestGetAccount(t *testing.T) {

	ac:=CreateAccountRandom(t)

	account,err:=testQueries.GetAccount(context.Background(), ac.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, ac.ID, account.ID)


}

func TestUpdateAccount(t *testing.T) {

	ac:=CreateAccountRandom(t)

	arg:= UpdateAccountParams{
		ID: ac.ID,
		Balance: util.RandomMoney(),
	}

	account,err:=testQueries.UpdateAccount(context.Background(), arg)


	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, ac.ID, account.ID)
	require.Equal(t, arg.Balance, account.Balance)


}


func TestDeleteAccount(t *testing.T) {

	ac:=CreateAccountRandom(t)

	err:=testQueries.DeleteAccount(context.Background(), ac.ID)

	require.NoError(t, err)
	

}


func TestListAccount(t *testing.T) {

  arg:= ListAccountsParams{
		Limit: 5,
		Offset: 0,

  }

	lacc,err:=testQueries.ListAccounts(context.Background(),arg)

	require.NoError(t, err)
	require.NotEmpty(t, lacc)

	require.Len(t, lacc, 5)

	

}


