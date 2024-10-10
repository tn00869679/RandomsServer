package blockchain

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func waitForConfirmations(client *ethclient.Client, signedTx *types.Transaction, confirmations uint64) (*types.Receipt, error) {
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		return nil, fmt.Errorf("transaction mining error: %v", err)
	}

	blockNumber := receipt.BlockNumber.Uint64()

	for {
		currentBlock, err := client.BlockByNumber(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get current block: %v", err)
		}

		currentBlockNumber := currentBlock.NumberU64()

		if currentBlockNumber >= blockNumber+confirmations {
			log.Printf("Transaction confirmed with %d confirmations\n", confirmations)

			latestReceipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
			if err != nil {
				return nil, fmt.Errorf("failed to get latest transaction receipt: %v", err)
			}

			return latestReceipt, nil
		}

		time.Sleep(15 * time.Second)
	}
}
