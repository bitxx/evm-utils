package model

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"
)

type Transaction struct {
	chain *Chain
}

func NewTransaction(chain *Chain) *Transaction {
	return &Transaction{
		chain: chain,
	}
}

func (t *Transaction) TransactionsByBlockNum(number uint64) (types.Transactions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.chain.Timeout)*time.Second)
	defer cancel()
	block, err := t.chain.RemoteRpcClient.BlockByNumber(ctx, new(big.Int).SetUint64(number))
	if err != nil {
		return nil, err
	}
	return block.Transactions(), nil
}
