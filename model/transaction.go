package model

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"math/big"
	"strconv"
	"time"
)

type Transaction struct {
	Hash      string
	Protected bool
	Nonce     uint64
	Data      []byte
	Size      uint64
	Value     decimal.Decimal
	GasPrice  decimal.Decimal
	Type      string
	ChainId   decimal.Decimal
	Cost      decimal.Decimal
	GasFeeCap decimal.Decimal
	GasTipCap decimal.Decimal
	To        string
	From      string
	Time      time.Time
	//TODO blobs待定

	chain *Chain
}

func NewTransaction(chain *Chain) *Transaction {
	return &Transaction{
		chain: chain,
	}
}

// BlockByNumber
//
//	@Description: 获取块
//	@receiver t
//	@param number
//	@return *types.Block
//	@return error
func (t *Transaction) BlockByNumber(number uint64) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.chain.Timeout)*time.Second)
	defer cancel()
	return t.chain.RemoteRpcClient.BlockByNumber(ctx, new(big.Int).SetUint64(number))
}

// TransactionsByBlockNumber
//
//	@Description: 获取一个块的交易
//	@receiver t
//	@param number
//	@return []Transaction
//	@return error
func (t *Transaction) TransactionsByBlockNumber(number uint64) ([]Transaction, error) {
	block, err := t.BlockByNumber(number)
	if err != nil {
		return nil, err
	}
	var transactions []Transaction
	for _, tx := range block.Transactions() {
		tx.MarshalJSON()
		var signer types.Signer
		switch {
		case tx.Type() == types.AccessListTxType:
			signer = types.NewEIP2930Signer(tx.ChainId())
		case tx.Type() == types.DynamicFeeTxType:
			signer = types.NewLondonSigner(tx.ChainId())
		default:
			signer = types.NewEIP155Signer(tx.ChainId())
		}
		from, err := types.Sender(signer, tx)
		if err != nil {
			return nil, err
		}

		transaction := Transaction{
			Hash:      tx.Hash().String(),
			Protected: tx.Protected(),
			Nonce:     tx.Nonce(),
			Data:      tx.Data(),
			Size:      tx.Size(),
			Value:     decimal.NewFromBigInt(tx.Value(), 0),
			GasPrice:  decimal.NewFromBigInt(tx.GasPrice(), 0),
			Type:      strconv.Itoa(int(tx.Type())),
			ChainId:   decimal.NewFromBigInt(tx.ChainId(), 0),
			Cost:      decimal.NewFromBigInt(tx.Cost(), 0),
			GasFeeCap: decimal.NewFromBigInt(tx.GasFeeCap(), 0),
			GasTipCap: decimal.NewFromBigInt(tx.GasTipCap(), 0),
			To:        tx.To().String(),
			From:      from.String(),
			Time:      tx.Time(),
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
