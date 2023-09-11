package ethutil

import (
	"encoding/json"
	"fmt"
	"github.com/bitxx/ethutil/util/dateutil"
	"github.com/bitxx/ethutil/util/httputil"
	"github.com/bitxx/ethutil/util/idgenutil"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/status-im/keycard-go/hexutils"
	"math/rand"
	"strconv"
	"time"

	"github.com/bitxx/ethutil/config"

	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

const (
	rpcUrl  = "https://alpha-us-http-geth.opside.network/"
	timeout = 60 //second

	testAccountFromAddress           = "0x7a547A149A79A03F4dd441B6806ffCBb1b63F383"
	testAccountFromAddressPrivateKey = ""
	testAccountToAddress             = "0x8B63293748e058F47a31c0D2Af0B1b3FeDdc4D4C"
	accountFile                      = "./account.txt"
	addressFile                      = "./address.txt"
	privateKeyFile                   = "./privateKey.txt"
)

func MyClient() *EthClient {
	return NewEthClient(rpcUrl, timeout)
}

func TestAccountByMnemonic(t *testing.T) {
	num := 150
	result := "account：\n"
	addresses := "\n\n\nSummary of the above account addresses:"
	privateKeys := "\n\n\nSummary of the above account privateKeys:"
	menmonics := "\n\n\nSummary of the above account menmonics:"
	for i := 0; i < num; i++ {
		account, _ := MyClient().AccountByMnemonic()
		addresses = addresses + "\n" + account.Address
		privateKeys = privateKeys + "\n" + account.PrivateKey
		menmonics = menmonics + "\n" + account.Mnemonic
		result = result + fmt.Sprintf("NO. %d group account：\nmnemonic：%s\naddress：%s\nprivateKey：%s\npublicKey：%s\n\n", i+1, account.Mnemonic, account.Address, account.PrivateKey, account.PublicKey)
	}
	result = result + addresses + privateKeys + menmonics
	_ = os.WriteFile(accountFile, []byte(result), 0666)
}

func TestAccountWithPrivateKey(t *testing.T) {
	account1, _ := MyClient().AccountByMnemonic()
	account2, _ := MyClient().AccountWithPrivateKey(account1.PrivateKey)
	require.Equal(t, account1.PublicKey, account2.PublicKey)
}

func TestTokenBalance(t *testing.T) {
	balance, err := MyClient().TokenBalanceOf(testAccountToAddress)
	require.Nil(t, err)
	t.Log("balance：", balance)
}

func TestBatchTokenBalance(t *testing.T) {
	bytes, err := os.ReadFile(addressFile)
	require.Nil(t, err)
	addresses := strings.Split(string(bytes), "\n")
	client := MyClient()
	for i, address := range addresses {
		address = strings.TrimSpace(address)
		if len(address) <= 0 {
			continue
		}
		balance, err := client.TokenBalanceOf(address)
		if err != nil {
			t.Error(fmt.Sprintf("index: %d ,address:%s, request err: %s", i, address, err.Error()))
			break
		}
		/*value, err := decimal.NewFromString(balance)
		if err != nil {
			t.Error(fmt.Sprintf("index: %d ,address:%s, value err: %s", i, address, err.Error()))
			break
		}
		if value.Cmp(decimal.Zero) > 0 {
			println("the", i+1, "group,address: ", address)
		}*/

		t.Log(fmt.Sprintf("address: %s,balance: %s", address, balance))
	}
}

func TestNonce(t *testing.T) {
	nonce, err := MyClient().Nonce(testAccountToAddress)
	require.Nil(t, err)
	t.Log("nonce: ", nonce)
}

func TestTokenEstimateGasLimit(t *testing.T) {
	value := "13000000000000000000"
	data := hexutils.HexToBytes("73c45c98000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000540000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000000000000000014f9251c4ef894355d34d2a98a26e9d60b9d56d6bc0000000000000000000000000000000000000000000000000000000000000000000000000000000000000008a688906bd8b00000000000000000000000000000000000000000000000000000")
	gasLimit, err := MyClient().TokenEstimateGasLimit("0xf9251c4ef894355d34d2a98a26e9d60b9d56d6bc", "0xbd927011759b2c4f2602c3008f8ef3407db53473", config.DefaultEthGasPrice, value, data)
	require.Nil(t, err)
	t.Log("estimate gas limit: ", gasLimit)
}

func TestTokenTransfer(t *testing.T) {
	value := "13000000000000000000"
	data := hexutils.HexToBytes("73c45c98000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000540000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000000000000000000000000000000000000000000014f9251c4ef894355d34d2a98a26e9d60b9d56d6bc0000000000000000000000000000000000000000000000000000000000000000000000000000000000000008a688906bd8b00000000000000000000000000000000000000000000000000000")

	gasLimit, err := MyClient().TokenEstimateGasLimit("0xf9251c4ef894355d34d2a98a26e9d60b9d56d6bc", "0xbd927011759b2c4f2602c3008f8ef3407db53473", config.DefaultEthGasPrice, value, data)
	require.Nil(t, err)
	hash, err := MyClient().TokenTransfer("cfe211ce0489dce12e6cf204e720c25e01d654ada9d839df1727a59043287cbb", config.DefaultEthGasPrice, gasLimit, value, "0xbd927011759b2c4f2602c3008f8ef3407db53473", data)
	require.Nil(t, err)
	t.Log("hash:", hash)
}

func TestBatchTokenTransferToManyAddress(t *testing.T) {
	privateKey := ""
	gasLimit := "21000"

	bytes, err := os.ReadFile(addressFile)
	require.Nil(t, err)
	toAddresses := strings.Split(string(bytes), "\n")
	client := MyClient()
	total := 0
	for i, toAddress := range toAddresses {
		time.Sleep(1 * time.Second)
		//transfer random token to each address, token number from 2 to 45
		rValue := rand.Intn(45-2) + 2
		total = total + rValue
		value := strconv.Itoa(rValue) + "000000000000000000"
		hash, err := client.TokenTransfer(privateKey, config.DefaultEthGasPrice, gasLimit, value, toAddress, nil)
		if err != nil {
			t.Error(fmt.Sprintf("index:%d,toAddress: %s,error: %s", i, toAddress, err.Error()))
			continue
		}
		t.Log(fmt.Sprintf("index: %d,toAddress: %s,hash: %s,value: %s", i, toAddress, hash, strings.Replace(value, "000000000000000000", "", 1)))
	}
	t.Log(fmt.Sprintf("transfer over,address count:%d,all token: %d", len(toAddresses), total))

}

func TestBatchTokenTransferToOneAddress(t *testing.T) {
	toAddress := ""
	value := "25000000000000000000000"
	gasLimit := "21000"
	bytes, err := os.ReadFile(privateKeyFile)
	require.Nil(t, err)
	privateKeys := strings.Split(string(bytes), "\n")
	client := MyClient()
	for i, privateKey := range privateKeys {
		privateKey = strings.TrimSpace(privateKey)
		if len(privateKey) <= 0 {
			continue
		}
		hash, err := client.TokenTransfer(privateKey, config.DefaultEthGasPrice, gasLimit, value, toAddress, nil)
		if err != nil {
			t.Error(fmt.Sprintf("index:%d,privateKey: %s,error: %s", i, privateKey, err.Error()))
			continue
		}
		t.Log(fmt.Sprintf("index: %d,privateKey: %s,hash: %s", i, privateKey, hash))
	}
}

func TestMetamaskLoginSign(t *testing.T) {
	url := "https://graphigo.prd.galaxy.eco/query"
	privateKey := ""
	//1. 获取账户
	ethClient := NewSimpleEthClient()
	account, err := ethClient.AccountWithPrivateKey(privateKey)
	require.Nil(t, err)

	//2. 生成未签名消息
	version := "1"
	chainId := "1"
	nonce := idgenutil.ID()
	//nonce := "kxjFHNusQHg9vbEGl"
	now := time.Now()
	startTime := dateutil.ConvertToStr(time.Now(), 4)
	endTime := dateutil.ConvertToStr(now.AddDate(0, 0, 7), 4)
	/*startTime := "2023-09-05T09:23:27.197Z"
	endTime := "2023-09-12T09:23:27.173Z"*/

	msg := fmt.Sprintf("galxe.com wants you to sign in with your Ethereum account:\n%s\n\nSign in with Ethereum to the app.\n\nURI: https://galxe.com\nVersion: %s\nChain ID: %s\nNonce: %s\nIssued At: %s\nExpiration Time: %s", account.Address, version, chainId, nonce, startTime, endTime)

	//3. metamask消息签名
	sign, err := ethClient.MetamaskSignLogin(msg, privateKey)
	require.Nil(t, err)

	//4. 请求提交galxe登录信息
	param := map[string]interface{}{
		"operationName": "SignIn",
		"query":         "mutation SignIn($input: Auth) {\n  signin(input: $input)\n}\n",
		"variables": map[string]interface{}{
			"input": map[string]interface{}{
				"address":   account.Address,
				"message":   msg,
				"signature": sign,
			},
		},
	}
	data, err := httputil.JasonSend(url, httputil.PostMethod, param)
	require.Nil(t, err)

	//5. 返回结果
	resp := map[string]interface{}{}
	err = json.Unmarshal(data, &resp)
	require.Nil(t, err)

	r1, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("result is error")
	}
	token, ok := r1["signin"].(string)
	if !ok {
		t.Fatal("result signin is error")
	}
	t.Log("token: ", token)
}

func TestSignEip721(t *testing.T) {
	loginUrl := "https://opside.network/api/user/custom/login"
	ethClient := NewSimpleEthClient()
	privateKey := ""
	account, err := ethClient.AccountWithPrivateKey(privateKey)
	require.Nil(t, err)
	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
			},
			"Action": {
				{Name: "actionType", Type: "string"},
				{Name: "timestamp", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
			},
		},
		PrimaryType: "Action",
		Domain: apitypes.TypedDataDomain{
			Name:    "Opside",
			Version: "1",
		},
		Message: map[string]interface{}{
			"actionType": "LOGIN",
			"timestamp":  strconv.FormatInt(time.Now().UnixMilli(), 10),
			"nonce":      idgenutil.IDNum(),
		},
	}
	signature, err := ethClient.SignEip721(privateKey, &typedData)
	require.Nil(t, err)

	data, err := json.Marshal(typedData.Map())
	require.Nil(t, err)

	param := map[string]interface{}{
		"payload":   string(data),
		"signature": signature,
		"address":   account.Address,
	}
	header := map[string]string{
		"Content-Type": "application/json",
	}
	c := httputil.NewHttpSend(loginUrl)
	c.SetSendType(httputil.SendtypeJson)
	c.SetBody(param)
	c.SetHeader(header)
	result, err := c.Post()
	require.Nil(t, err)

	resp := map[string]interface{}{}
	err = json.Unmarshal(result, &resp)
	require.Nil(t, err)

	token, ok := resp["result"].(string)
	if !ok {
		t.Error("result is error")
	}
	t.Log("token: ", token)
}
