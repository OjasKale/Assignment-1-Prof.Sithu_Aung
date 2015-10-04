package main

import (
	"fmt"
	"net/rpc"
	"bufio"
	"log"
	"os"
)


type Args struct {
	StockSymbolAndPercentage string
	UserBudget               float64
}

type Quote struct {
	Stocks         string  `json:"stocksymbol"`
	UnvestedAmount float64 `json:"stockprice"`
	TrdId        int     `json:"id"`
}

type Id struct {
	TrdId int `json:"id"`
}

type UpdQuote struct {
	Stocks         string  `json:"stocksymbol"`
	UnvestedAmount float64 `json:"stockprice"`
}


func main() {


	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1528")

	var ch int
fmt.Println("1. Buying Stocks")
fmt.Println("2. Portfolio")
fmt.Print("Enter you choice: ")
fmt.Scanf("%d",&ch)
switch ch {
case 1:
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Stock Symbol: ")
	StockSymbolAndPercentage, _ := reader.ReadString('\n')
	fmt.Println(StockSymbolAndPercentage)
	fmt.Print("Enter budget: ")
	var UserBudget float64
	fmt.Scan(&UserBudget)

	args := Args{StockSymbolAndPercentage, UserBudget}

	var quo Quote
	err = client.Call("StkC.StockPrice", args, &quo)

	if err != nil {
		log.Fatal("arith error:", err)
	}

	fmt.Print("Stocks : ")
	fmt.Println(quo.Stocks)
	fmt.Print("Amount: ")
	fmt.Println(quo.UnvestedAmount)
	fmt.Print("ID: ")
	fmt.Println(quo.TrdId)
	break
case 2:

	fmt.Println("Enter your trade id")
	var id int
	fmt.Scanf("%d",&id)
	Tid := Id{id}
	var quo Quote
	err = client.Call("StkC.Portfo",Tid,&quo)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Print("Stock : Current market price : bought ")
	fmt.Println(quo.Stocks)
	fmt.Print("Unvested Amount: ")
	fmt.Println(quo.UnvestedAmount)
}


	}
