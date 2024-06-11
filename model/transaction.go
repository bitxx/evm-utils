package model

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	"math/big"
	"strconv"
	"time"
)

type Transaction struct {
	Hash              string
	Protected         bool
	Nonce             uint64
	Data              []byte
	Size              uint64
	Value             decimal.Decimal
	GasPrice          decimal.Decimal
	Type              string
	ChainId           decimal.Decimal
	Gas               uint64
	Cost              decimal.Decimal
	GasFeeCap         decimal.Decimal
	GasTipCap         decimal.Decimal
	To                string
	From              string
	Time              uint64
	GasUsed           uint64
	CumulativeGasUsed uint64
	ReceiptStatus     uint64
	EffectiveGasPrice decimal.Decimal
	BlobGasUsed       uint64
	BlobGasPrice      decimal.Decimal
	TransactionIndex  uint
	ContractAddress   string

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
//	@Description: 根据编号获取块
//	@receiver t
//	@param number 如果number<=0，则读取最新块
//	@return *types.Block
//	@return error
func (t *Transaction) BlockByNumber(number uint64) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.chain.Timeout)*time.Second)
	defer cancel()
	if number <= 0 {
		return t.chain.RemoteRpcClient.BlockByNumber(ctx, nil)
	}
	return t.chain.RemoteRpcClient.BlockByNumber(ctx, new(big.Int).SetUint64(number))
}

// BlockReceiptsByNumber
//
//	@Description: 读取一个块交易的回执
//	@receiver t
//	@param number
//	@return []*types.Receipt
//	@return error
/*func (t *Transaction) BlockReceiptsByNumber(number uint64) ([]*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.chain.Timeout)*time.Second)
	defer cancel()
	if number <= 0 {
		return t.chain.RemoteRpcClient.BlockReceipts(ctx, rpc.BlockNumberOrHash{})
	}
	n := rpc.BlockNumber(1)
	return t.chain.RemoteRpcClient.BlockReceipts(ctx, rpc.BlockNumberOrHash{BlockNumber: &n, RequireCanonical: false})
}*/

// TxByHash
//
//	@Description: 根据hash获取交易记录回执
//	@receiver t
//	@param hash
//	@return *types.Receipt
//	@return error
func (t *Transaction) TxByHash(hash string) (*Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.chain.Timeout)*time.Second)
	defer cancel()

	tx, _, err := t.chain.RemoteRpcClient.TransactionByHash(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, err
	}
	return t.parseTx(ctx, tx, nil)

}

func (t *Transaction) parseTx(ctx context.Context, tx *types.Transaction, block *types.Block) (*Transaction, error) {

	receipt, err := t.chain.RemoteRpcClient.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return nil, err
	}

	if block == nil {
		block, err = t.BlockByNumber(receipt.BlockNumber.Uint64())
		if err != nil {
			return nil, err
		}
	}

	var signer types.Signer
	switch {
	case tx.Type() == types.AccessListTxType:
		signer = types.NewEIP2930Signer(tx.ChainId())
	case tx.Type() == types.DynamicFeeTxType:
		signer = types.NewLondonSigner(tx.ChainId())
	case tx.Type() == types.BlobTxType:
		signer = types.NewCancunSigner(tx.ChainId())
	default:
		signer = types.NewEIP155Signer(tx.ChainId())
	}
	from, err := types.Sender(signer, tx)
	if err != nil {
		return nil, err
	}

	to := ""
	if tx.To() != nil {
		to = tx.To().String()
	}

	effectiveGasPrice := receipt.EffectiveGasPrice
	if effectiveGasPrice == nil {
		effectiveGasPrice = big.NewInt(0)
	}

	blobGasPrice := receipt.BlobGasPrice
	if blobGasPrice == nil {
		blobGasPrice = big.NewInt(0)
	}

	return &Transaction{
		Hash:              tx.Hash().String(),
		Protected:         tx.Protected(),
		Nonce:             tx.Nonce(),
		Data:              tx.Data(),
		Size:              tx.Size(),
		Gas:               tx.Gas(),
		Value:             decimal.NewFromBigInt(tx.Value(), 0),
		GasPrice:          decimal.NewFromBigInt(tx.GasPrice(), 0),
		Type:              strconv.Itoa(int(tx.Type())),
		ChainId:           decimal.NewFromBigInt(tx.ChainId(), 0),
		Cost:              decimal.NewFromBigInt(tx.Cost(), 0),
		GasFeeCap:         decimal.NewFromBigInt(tx.GasFeeCap(), 0),
		GasTipCap:         decimal.NewFromBigInt(tx.GasTipCap(), 0),
		To:                to,
		From:              from.String(),
		Time:              block.Time(),
		GasUsed:           receipt.GasUsed,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		ReceiptStatus:     receipt.Status,
		EffectiveGasPrice: decimal.NewFromBigInt(effectiveGasPrice, 0),
		BlobGasUsed:       receipt.BlobGasUsed,
		BlobGasPrice:      decimal.NewFromBigInt(blobGasPrice, 0),
		TransactionIndex:  receipt.TransactionIndex,
		ContractAddress:   receipt.ContractAddress.Hex(),
	}, nil
}

// TxByBlockNumber
//
//	@Description: 获取一个块的交易
//	@receiver t
//	@param number
//	@return []Transaction
//	@return error
func (t *Transaction) TxByBlockNumber(number uint64) ([]Transaction, error) {
	block, err := t.BlockByNumber(number)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.chain.Timeout)*time.Second)
	defer cancel()
	var transactions []Transaction
	for _, tx := range block.Transactions() {
		if tx.Hash().Hex() == "" {
			continue
		}
		transaction, err := t.parseTx(ctx, tx, block)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *transaction)
	}
	return transactions, nil
}

// TxIsPending
//
//	@Description: is pendding
//	@receiver t
//	@param hash
//	@return bool
//	@return error
func (t *Transaction) TxIsPending(hash string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.chain.Timeout)*time.Second)
	defer cancel()
	_, isPending, err := t.chain.RemoteRpcClient.TransactionByHash(ctx, common.HexToHash(hash))
	return isPending, err
}

// LatestBlockNumber
//
//	@Description: 获取最新块高度
//	@receiver t
//	@return uint64
//	@return error
func (t *Transaction) LatestBlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.chain.Timeout)*time.Second)
	defer cancel()
	return t.chain.RemoteRpcClient.BlockNumber(ctx)
}
