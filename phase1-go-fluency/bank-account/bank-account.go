package main

import "errors"

type Account struct {
	accountName    string
	accountBalance float64
}

func (account *Account) deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("must deposit a value greater than 0")
	}
	account.accountBalance += amount
	return nil
}

func (account *Account) withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("must withdraw a value greater than 0")
	}
	if amount > account.accountBalance {
		return errors.New("must withdraw a value smaller than or equal to account balance")
	}
	account.accountBalance -= amount
	return nil
}

func (account *Account) checkBalance() float64 {
	return account.accountBalance
}

func (account *Account) getAccountOwner() string {
	return account.accountName
}
