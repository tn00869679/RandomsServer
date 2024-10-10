package blockchain

import (
	"RandomsServer/blockchain/abi"
	"RandomsServer/util"
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func InjectRandoms(client *ethclient.Client, randomsDatas <-chan util.RandomsData) {
	applicationContractAbi := abi.Read("ApplicationContract")

	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatalln("Failed to load private key:", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalln("Failed to cast public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	for data := range randomsDatas {
		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatalln("Failed to get nonce:", err)
		}

		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatalln("Failed to suggest gas price:", err)
		}

		inputData, err := applicationContractAbi.Pack("updateRequestIdAndRandoms", data.RequestId, data.Randoms)
		if err != nil {
			log.Fatalln("Failed to pack input data:", err)
		}

		gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
			From:     fromAddress,
			To:       &data.ContractAddress,
			GasPrice: gasPrice,
			Value:    big.NewInt(0),
			Data:     inputData,
		})
		if err != nil {
			log.Fatalln("Failed to estimate gas:", err)
		}

		tx := types.NewTransaction(nonce, data.ContractAddress, big.NewInt(0), gasLimit, gasPrice, inputData)

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatalln("Failed to get chain ID:", err)
		}

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			log.Fatalln("Failed to sign transaction:", err)
		}

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			log.Fatalln("Failed to send transaction:", err)
		}

		log.Printf("Transaction sent! TX Hash: %s\n", signedTx.Hash().Hex())

		confirmations := uint64(5)
		receipt, err := waitForConfirmations(client, signedTx, confirmations)
		if err != nil && err != rpc.ErrNoResult {
			log.Println(err)
		}
		if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
			log.Printf("Success transaction, BlockNumber: %d\n", receipt.BlockNumber.Uint64())
		}
	}
}
