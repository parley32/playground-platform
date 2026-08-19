package main

import "testing"

func TestAccountDeposit(t *testing.T) {
	tests := []struct {
		name        string
		initial     float64
		amount      float64
		wantBalance float64
		wantErr     bool
	}{
		{
			name:        "deposits a positive amount",
			initial:     0,
			amount:      100,
			wantBalance: 100,
		},
		{
			name:        "adds to existing balance",
			initial:     50,
			amount:      25.50,
			wantBalance: 75.50,
		},
		{
			name:        "rejects zero",
			initial:     50,
			amount:      0,
			wantBalance: 50,
			wantErr:     true,
		},
		{
			name:        "rejects a negative amount",
			initial:     50,
			amount:      -10,
			wantBalance: 50,
			wantErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account := Account{accountBalance: test.initial}

			err := account.deposit(test.amount)
			if (err != nil) != test.wantErr {
				t.Fatalf("deposit() error = %v, wantErr %v", err, test.wantErr)
			}

			if account.accountBalance != test.wantBalance {
				t.Errorf("account balance = %.2f, want %.2f", account.accountBalance, test.wantBalance)
			}
		})
	}
}

func TestAccountWithdrawal(t *testing.T) {
	tests := []struct {
		name        string
		initial     float64
		amount      float64
		wantBalance float64
		wantErr     bool
	}{
		{
			name:        "withdraw a positive amount",
			initial:     100,
			amount:      50,
			wantBalance: 50,
		},
		{
			name:        "rejects an amount larger than account balance",
			initial:     50,
			amount:      100,
			wantBalance: 50,
			wantErr:     true,
		},
		{
			name:        "rejects zero",
			initial:     50,
			amount:      0,
			wantBalance: 50,
			wantErr:     true,
		},
		{
			name:        "rejects a negative amount",
			initial:     50,
			amount:      -10,
			wantBalance: 50,
			wantErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			account := Account{accountBalance: test.initial}

			err := account.withdraw(test.amount)
			if (err != nil) != test.wantErr {
				t.Fatalf("withdraw() error = %v, wantErr %v", err, test.wantErr)
			}

			if account.accountBalance != test.wantBalance {
				t.Errorf("account balance = %.2f, want %.2f", account.accountBalance, test.wantBalance)
			}
		})
	}
}

func TestAccountOwner(t *testing.T) {
	want := "lachlan"
	account := Account{accountName: want}
	got := account.getAccountOwner()
	if got != want {
		t.Errorf("account owner function returned = %v, want %v", got, want)
	}
}

func TestAccountBalanceCheck(t *testing.T) {
	want := float64(100)
	account := Account{accountBalance: want}
	got := account.checkBalance()
	if got != want {
		t.Errorf("account balance function returned = %0.2f, want %0.2f", got, want)
	}
}
