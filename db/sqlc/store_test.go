package db

import (
	"context"
	"testing"
	"log"
	"github.com/stretchr/testify/require"
	"fmt"
)


func TestTransferTx(t *testing.T){


	store:= NewStore(testDB)

	account1:= CreateAccountRandom(t)
	account2:= CreateAccountRandom(t)


	log.Printf("Account1 id: %v Balance1 %v  Account2 Id: %v Balance2 %v",account1.ID,account1.Balance,account2.ID,account2.Balance)


	//run n transaction 
	n:=5
	amount:= int64(10)
	errsChan:= make(chan error)
	resultsChan:= make(chan TransferTxResult)



		for i:=0; i<n; i++{

			txName:= fmt.Sprintf("tx %d",i+1)
			go func(){
				ctx:= context.WithValue(context.Background(),txKey,txName)
				result,err:= store.TransferTx(ctx, TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID: account2.ID,
					Amount: amount,
				})

				if err!=nil {
					log.Println("Error :",err)
				} else {
					log.Println("Result :", result.Transfer.ID)
	
					
				}

				errsChan <- err
				resultsChan <- result

			}()
		}

		//check results
		for i:= 0; i< n ; i++ {

		err:= <- errsChan

		require.NoError(t, err)

		result:= <-resultsChan
		require.NotEmpty(t,result)
		transfer:= result.Transfer
		require.NotEmpty(t,transfer)
		require.Equal(t,account1.ID,transfer.FromAccountID)
		require.Equal(t,account2.ID,transfer.ToAccountID)
		require.Equal(t, amount,transfer.Amount)
		require.NotZero(t,transfer.ID)
		require.NotZero(t,transfer.CreatedAt)

		_,err= store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t,err)

		//check entries
		fromEntry:=result.FromEntry
		require.NotEmpty(t,fromEntry)
		require.Equal(t,account1.ID,fromEntry.AccountID)
		require.Equal(t,-amount,fromEntry.Amount)
		require.NotZero(t,fromEntry.ID)
		require.NotZero(t,fromEntry.CreatedAt)

		_,err= store.GetEntry(context.Background(),fromEntry.ID)
		require.NoError(t,err)

		toEntry:=result.ToEntry
		require.NotEmpty(t,toEntry)
		require.Equal(t,account2.ID,toEntry.AccountID)
		require.Equal(t,amount,toEntry.Amount)
		require.NotZero(t,toEntry.ID)
		require.NotZero(t,toEntry.CreatedAt)

		_,err= store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t,err)

		// chec cuentas
		fromAccount:= result.FromAccount
		require.NotEmpty(t,fromAccount)
		require.Equal(t,account1.ID, fromAccount.ID)

		toAccount:= result.ToAccount
		require.NotEmpty(t,toAccount)
		require.Equal(t,account2.ID,toAccount.ID)

		//check account balance
		diff1:= account1.Balance - fromAccount.Balance
		diff2:= toAccount.Balance - account2.Balance
		log.Printf("Dif1 %v Diff2 %v",diff1,diff2)
		require.Equal(t,diff1,diff2)
		require.True(t,diff1>0)
		require.True(t,diff1%amount==0)

		log.Printf("Post - Account1 id: %v Balance1 %v  Account2 Id: %v Balance2 %v",account1.ID,account1.Balance,account2.ID,account2.Balance)



	}
}


	func TestTransferDeadLockTx(t *testing.T){


		store:= NewStore(testDB)
	
		account1:= CreateAccountRandom(t)
		account2:= CreateAccountRandom(t)
	
	
		log.Printf("Account1 id: %v Balance1 %v  Account2 Id: %v Balance2 %v",account1.ID,account1.Balance,account2.ID,account2.Balance)
	
	
		//run n transaction 
		n:=5
		amount:= int64(10)
		errsChan:= make(chan error)

	
	
	
			for i:=0; i<n; i++{

				fromAccountID:=account1.ID
				toAccountID:=account2.ID

				if i%2 ==1{

					fromAccountID=account2.ID
					toAccountID=account1.ID
				}

	
				txName:= fmt.Sprintf("tx %d",i+1)
				go func(){
					ctx:= context.WithValue(context.Background(),txKey,txName)
					result,err:= store.TransferTx(ctx, TransferTxParams{
						FromAccountID: fromAccountID,
						ToAccountID: toAccountID,
						Amount: amount,
					})
	
					if err!=nil {
						log.Println("Error :",err)
					} else {
						log.Println("Result :", result.Transfer.ID)
		
						
					}
	
					errsChan <- err
					
	
				}()
			}
	
			//check results
			for i:= 0; i< n ; i++ {
	
			err:= <- errsChan
	
			require.NoError(t, err)

	
			
	
			log.Printf("Post - Account1 id: %v Balance1 %v  Account2 Id: %v Balance2 %v",account1.ID,account1.Balance,account2.ID,account2.Balance)
	
	
	
		}
	}

