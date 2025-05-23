package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct{
	*Queries
	db * sql.DB
}


func NewStore(db *sql.DB) Store{
	return Store{
		Queries: New(db),
		db: db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error{


	tx,err:= store.db.BeginTx(ctx, nil)

	if err != nil{
		return err
	}

	q:= New(tx)
	err= fn(q)
	if err != nil{

		rbErr:=tx.Rollback();
		if rbErr!= nil{
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()

	
}

type TransferTxParams struct{
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 	`json:"to_account_id"` 
	Amount int64 `json:"amount"`

}

type TransferTxResult struct{

	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

var txKey= struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error){
	
	var result TransferTxResult

	err:= store.execTx(ctx, func(q *Queries ) error{
		
		var err error

		txName:= ctx.Value(txKey)

		fmt.Println(txName,"create transfer")
		result.Transfer,err= q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err!= nil {
			return err
		}

		fmt.Println("id transfer", result.Transfer.ID)



		fmt.Println(txName,"createEntry From")
		result.FromEntry, err= q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err!=nil {

			return err
		}

		fmt.Println(txName,"createEntry To")
		result.ToEntry, err= q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err!=nil {

			return err
		}

		//TODO : update account balance

		fmt.Println(txName,"get acount 1")
		/*
		account1,err:= q.GetAccountForUpdate(ctx,arg.FromAccountID)
		if err!=nil {

			return err
		}

		fmt.Println(txName,"update acount 1")
		result.FromAccount, err= q.UpdateAccount(ctx, UpdateAccountParams{
			ID: arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})

		if err!=nil {

			return err
		}*/

		if arg.FromAccountID < arg.ToAccountID{

		
		result.FromAccount, err= q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.FromAccountID,
			Amount: -arg.Amount,
		})

		if err!=nil {

			return err
		}


		fmt.Println(txName,"get acount 2")
		/*
		account2,err:= q.GetAccountForUpdate(ctx,arg.ToAccountID)
		if err!=nil {

			return err
		}

		fmt.Println(txName,"update acount 2")
		result.ToAccount, err= q.UpdateAccount(ctx, UpdateAccountParams{
			ID: arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})

		if err!=nil {

			return err
		}
		*/


		result.ToAccount, err= q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err!=nil {

			return err
		}
	}else {

		result.ToAccount, err= q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.ToAccountID,
			Amount: arg.Amount,
		})

		if err!=nil {

			return err
		}


		result.FromAccount, err= q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: arg.FromAccountID,
			Amount: -arg.Amount,
		})

		if err!=nil {

			return err
		}


		fmt.Println(txName,"get acount 2")

	}
		
		
		return nil
	})

	return result, err


}