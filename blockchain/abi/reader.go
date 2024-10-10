package abi

import (
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func Read(fileName string) abi.ABI {
	abiData, err := os.ReadFile("./blockchain/abi/" + fileName + ".json")
	if err != nil {
		log.Fatalln("Failed to read ABI file:", err)
	}
	contractAbi, err := abi.JSON(strings.NewReader(string(abiData)))
	if err != nil {
		log.Fatalln("Failed to parse ABI:", err)
	}

	return contractAbi
}
