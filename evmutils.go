package evmutils

import (
	"errors"
	"github.com/bitxx/evm-utils/model"
	"github.com/bitxx/evm-utils/model/contract/erc20"
	"github.com/bitxx/evm-utils/util/signutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type EvmClient struct {
	RpcUrl  string
	timeout int64
}

// NewEthClient
//
//	@Description: if rpcUrl and timeout is empty,you can‘t connect the node,but you can use the function about wallet
//	@param rpcUrl
//	@param timeout
//	@return *EvmClient
func NewEthClient(rpcUrl string, timeout int64) *EvmClient {
	return &EvmClient{
		RpcUrl:  rpcUrl,
		timeout: timeout,
	}
}

// NewSimpleEthClient
//
//	@Description: not support connect to the node
//	@return *EvmClient
func NewSimpleEthClient() *EvmClient {
	return &EvmClient{}
}

func (o *EvmClient) AccountByMnemonic() (account *model.Account, err error) {
	return model.NewAccount().AccountByMnemonic()
}

// AccountInfoByMnemonic
//
//	@Description:
//	@receiver o
//	@param mnemonic
//	@return account
//	@return err
func (o *EvmClient) AccountInfoByMnemonic(mnemonic string) (account *model.Account, err error) {
	return model.NewAccount().AccountInfoByMnemonic(mnemonic)
}

func (o *EvmClient) AccountWithPrivateKey(privateKey string) (account *model.Account, err error) {
	return model.NewAccount().AccountWithPrivateKey(privateKey)
}

func (o *EvmClient) AccountGenKeystore(privateKey, pwd, path string) (address string, err error) {
	return model.NewAccount().AccountGenKeystore(privateKey, pwd, path)
}

func (o *EvmClient) TokenBalanceOf(address string) (balance string, err error) {
	chain, err := o.Chain()
	if err != nil {
		return "", err
	}
	token := model.NewToken(chain)
	return token.BalanceOf(address)
}

// TokenEstimateGasLimit
//
//	@Description: 估算gas ，如果是合约地址，data肯定不得为空
//	@receiver o
//	@param fromAddress
//	@param receiverAddress
//	@param gasPrice
//	@param amount
//	@param data
//	@return balance
//	@return err
func (o *EvmClient) TokenEstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount string, data []byte) (balance string, err error) {
	chain, err := o.Chain()
	if err != nil {
		return "", err
	}
	token := model.NewToken(chain)
	return token.EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount, data)
}

func (o *EvmClient) Chain() (*model.Chain, error) {
	return model.GetChain(o.RpcUrl, o.timeout)
}

func (o *EvmClient) Nonce(address string) (nonce uint64, err error) {
	chain, err := o.Chain()
	if err != nil {
		return 0, err
	}
	return chain.Nonce(address)
}

func (o *EvmClient) TokenTransfer(privateKey, nonce, gasPrice, gasLimit, maxPriorityFeePerGas, value, to, data string) (hash string, err error) {
	chain, err := o.Chain()
	if err != nil {
		return "", err
	}
	token := model.NewToken(chain)
	return token.Transfer(privateKey, nonce, gasPrice, gasLimit, maxPriorityFeePerGas, value, to, data)
}

// TxByBlockNumber
//
//	@Description: get all tx by block number
//	@receiver o
//	@param number
//	@return []model.Transaction
//	@return error
func (o *EvmClient) TxByBlockNumber(number uint64) ([]model.Transaction, error) {
	chain, err := o.Chain()
	if err != nil {
		return nil, err
	}
	transaction := model.NewTransaction(chain)
	return transaction.TxByBlockNumber(number)
}

// BlockByNumber
//
//	@Description: 读取一个块
//	@receiver o
//	@param number 如果number<=0，则读取最新块
//	@return *types.Block
//	@return error
/*func (o *EvmClient) BlockByNumber(number uint64) (*types.Block, error) {
	chain, err := o.Chain()
	if err != nil {
		return nil, err
	}
	transaction := model.NewTransaction(chain)
	return transaction.BlockByNumber(number)
}*/

// BlockReceiptsByNumber
//
//	@Description: 读取一个块所有交易的回执
//	@receiver o
//	@param number
//	@return []*types.Receipt
//	@return error
/*func (o *EvmClient) BlockReceiptsByNumber(number uint64) ([]*types.Receipt, error) {
	chain, err := o.Chain()
	if err != nil {
		return nil, err
	}
	transaction := model.NewTransaction(chain)
	return transaction.BlockReceiptsByNumber(number)
}*/

// TxByHash
//
//	@Description: 根据hash获取交易回执
//	@receiver o
//	@param hash
//	@return *types.Receipt
//	@return error
func (o *EvmClient) TxByHash(hash string) (*model.Transaction, error) {
	chain, err := o.Chain()
	if err != nil {
		return nil, err
	}
	transaction := model.NewTransaction(chain)
	return transaction.TxByHash(hash)
}

// TxIsPending
//
//	@Description: is pending
//	@receiver o
//	@param hash
//	@return bool
//	@return error
func (o *EvmClient) TxIsPending(hash string) (bool, error) {
	chain, err := o.Chain()
	if err != nil {
		return false, err
	}
	transaction := model.NewTransaction(chain)
	return transaction.TxIsPending(hash)
}

// LatestBlockNumber
//
//	@Description: 获取最新块
//	@receiver o
//	@return uint64
//	@return error
func (o *EvmClient) LatestBlockNumber() (uint64, error) {
	chain, err := o.Chain()
	if err != nil {
		return 0, err
	}
	transaction := model.NewTransaction(chain)
	return transaction.LatestBlockNumber()
}

// MetamaskSignLogin
//
//	@Description: metamask sign login
//	@receiver o
//	@param message
//	@param privateKey
//	@return string
//	@return error
func (o *EvmClient) MetamaskSignLogin(message, privateKey string) (string, error) {
	if message == "" || privateKey == "" {
		return "", errors.New("param is empty")
	}
	return signutil.MetamaskSignLogin(message, privateKey)
}

// SignEip721
//
//	@Description: eip721 sign
//	@receiver o
//	@param privateKey
//	@param typedData
//	@return string
//	@return error
func (o *EvmClient) SignEip721(privateKey string, typedData *apitypes.TypedData) (string, error) {
	if typedData == nil || privateKey == "" {
		return "", errors.New("param is empty")
	}
	return signutil.SignEip721(privateKey, typedData)
}

// TokenErc20BalanceOf
//
//	@Description: erc20 balance
//	@receiver o
//	@param address user's account address
//	@param contractAddress erc20 address
//	@opts options
//	@return balance
//	@return err
func (o *EvmClient) TokenErc20BalanceOf(address, contractAddress string, opts *bind.CallOpts) (balance string, err error) {
	chain, err := o.Chain()
	if err != nil {
		return "", err
	}

	link, err := erc20.NewERC20(common.HexToAddress(contractAddress), chain.RemoteRpcClient)
	b, err := link.BalanceOf(opts, common.HexToAddress(address))
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// TokenErc20Approve
//
//	@Description: approve erc20
//	@receiver o
//	@param spenderAddress
//	@param approveValue
//	@param contractAddress
//	@param opts
//	@return tx
//	@return err
/*func (o *EvmClient) TokenErc20Approve(privateKey, spenderAddress, approveValue, contractAddress string, opts *bind.TransactOpts) (tx string, err error) {
	chain, err := o.Chain()
	if err != nil {
		return "", err
	}
	var value *big.Int
	var valid bool
	if value, valid = big.NewInt(0).SetString(approveValue, 10); !valid {
		return "", errors.New("invalid approve value")
	}
	if value.Cmp(big.NewInt(0)) <= 0 {
		return "", errors.New("value need bigger than 0")
	}

	if opts == nil {
		opts = &bind.TransactOpts{}
	}

	singer, err := contract.Signer(privateKey, chain.ChainId)
	if err != nil {
		return "", err
	}
	opts.Signer = singer

	link, err := erc20.NewERC20(common.HexToAddress(contractAddress), chain.RemoteRpcClient)
	b, err := link.Approve(opts, common.HexToAddress(spenderAddress), value)
	if err != nil {
		return "", err
	}
	return b.Hash().String(), nil
}*/
