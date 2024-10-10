package util

import (
	"RandomsServer/blockchain/abi"
	"crypto/rand"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type RandomsData struct {
	ContractAddress common.Address
	RequestId       common.Hash
	Randoms         *big.Int
}

var eventData struct {
	SubContract common.Address
	RequestId   common.Hash
}

func Randoms(oracleEvents <-chan types.Log, randomsDatas chan<- RandomsData) {
	simpleOracleAbi := abi.Read("SimpleOracle")

	for event := range oracleEvents {
		err := simpleOracleAbi.UnpackIntoInterface(&eventData, "RequestedData", event.Data)
		if err != nil {
			log.Fatalln("Failed to unpack event data:", err)
		}

		newRandoms, err := rand.Int(rand.Reader, big.NewInt(1000))
		if err != nil {
			log.Fatalln("生成真亂數時發生錯誤:", err)
			return
		}
		log.Println("新真亂數：", newRandoms)

		injectionData := &RandomsData{
			ContractAddress: eventData.SubContract,
			RequestId:       eventData.RequestId,
			Randoms:         newRandoms,
		}

		randomsDatas <- *injectionData
	}
}
