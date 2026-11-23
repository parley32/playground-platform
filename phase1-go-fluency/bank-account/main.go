package main

import "fmt"

func main() {
	account := Account{accountName: "Lachlan"}

	if err := account.deposit(40); err != nil {
		fmt.Println(err)
		return
	}

	if err := account.withdraw(50); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Account balance is %.2f\n", account.checkBalance())
}
