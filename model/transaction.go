package model

import (
	"context"
	"encoding/json"
	"fmt"
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

func (t *Transaction) Transactions(number uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.chain.Timeout)*time.Second)
	defer cancel()
	block, err := t.chain.RemoteRpcClient.BlockByNumber(ctx, new(big.Int).SetUint64(number))
	if err != nil {
		return err
	}

	//fmt.Println(block)
	b, _ := json.Marshal(block.Transactions())
	fmt.Println(string(b))
	return nil
}
