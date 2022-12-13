package go_aptos_sdk

import (
	"fmt"
	"testing"
)

var err error
var client *AleoClient

func init() {
	endpoint := "https://vm.aleo.org/api"
	client, err = NewAleoClient(endpoint, true)
	if err != nil {
		panic(err)
	}
}

func TestAleoClient_GetLatestHeight(t *testing.T) {
	height, err := client.GetLatestHeight()
	if err != nil {
		panic(err)
	}
	fmt.Println(height)
}

func TestAleoClient_GetTransactionsByHeight(t *testing.T) {
	transactions, err := client.GetTransactionsByHeight(194728)
	if err != nil {
		panic(err)
	}
	fmt.Println(transactions)
}

func TestAleoClient_Transaction(t *testing.T) {
	tx, err := client.Transaction("")
	if err != nil {
		panic(err)
	}
	fmt.Println(tx)
}
