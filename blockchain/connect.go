package blockchain

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

const maxAttempts = 3
const retryDelay = 5 * time.Second

func Connect(rpcUrl string) (*ethclient.Client, error) {
	var client *ethclient.Client
	var err error

	attempts := 0
	for attempts < maxAttempts {
		client, err = ethclient.Dial(rpcUrl)
		if err == nil {
			return client, nil
		}

		attempts++
		log.Printf("Failed to connect, attempt %d out of %d", attempts, maxAttempts)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect after %d attempts: %v", maxAttempts, err)
}
