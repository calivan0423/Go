package main

import (
	"fmt"
	//"log"

	"github.com/calivan0423/learngo/accounts"
)

func main() {
	account := accounts.NewAccount("calivan")
	account.Deposit(10)
	fmt.Println(account.Balance())
	err := account.Withdraw(20)
	if err != nil {
		//log.Fatalln(err)
		fmt.Println(err)
	}
	fmt.Println(account)
}
