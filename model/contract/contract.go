package contract

import (
	"github.com/bitxx/evm-utils/util"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func Signer(privateKey string, chainId *big.Int) (bind.SignerFn, error) {
	priData, err := util.HexDecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	privateKeyECDSA, err := crypto.ToECDSA(priData)
	if err != nil {
		return nil, err
	}

	return func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		signer := types.LatestSignerForChainID(chainId)
		signature, err := crypto.Sign(signer.Hash(tx).Bytes(), privateKeyECDSA)
		if err != nil {
			return nil, err
		}
		return tx.WithSignature(signer, signature)
	}, nil
}
