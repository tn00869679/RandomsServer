package script

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func MonitorOracleEvent(client *ethclient.Client, eventLogs chan types.Log) {
	contractAddress := common.HexToAddress("0x3a138Fd97C3df2D02B891822b9a3eA980CB1Cd0f")

	eventSignature := []byte("RequestedData(address,bytes32)")
	eventSignatureHash := crypto.Keccak256Hash(eventSignature)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{eventSignatureHash}},
	}

	sub, err := client.SubscribeFilterLogs(context.Background(), query, eventLogs)
	if err != nil {
		log.Fatalln("Failed to subscribe to events:", err)
	}

	for err := range sub.Err() {
		log.Fatalln("<-sub.Err (Oracle):", err)
		return
	}
}
