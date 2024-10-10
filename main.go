package main

import (
	"RandomsServer/blockchain"
	"RandomsServer/script"
	"RandomsServer/util"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	log.Println("Start!!")
	client, err := blockchain.Connect("wss://bsc-testnet-rpc.publicnode.com")
	if err != nil {
		panic("Fuck you")
	}

	eventLogs := make(chan types.Log)
	randomsDatas := make(chan util.RandomsData)
	go script.MonitorOracleEvent(client, eventLogs)
	go util.Randoms(eventLogs, randomsDatas)
	go blockchain.InjectRandoms(client, randomsDatas)

	select {}
}
