package main

import (
	"microservice-with-grpc/entity"
	"strconv"
)

func TransferBalance(sourceAccount, destinationAccount *entity.Account, amount string) {
	sourceBalance, _ := strconv.Atoi(sourceAccount.Balance)
	destinationBalance, _ := strconv.Atoi(destinationAccount.Balance)
	amountToTransfer, _ := strconv.Atoi(amount)

	sourceBalance -= amountToTransfer
	destinationBalance += amountToTransfer

	sourceAccount.Balance = strconv.Itoa(sourceBalance)
	destinationAccount.Balance = strconv.Itoa(destinationBalance)
}
