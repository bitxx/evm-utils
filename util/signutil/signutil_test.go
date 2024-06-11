package signutil

import (
	"encoding/json"
	"fmt"
	evmutils "github.com/bitxx/evm-utils"
	"github.com/bitxx/evm-utils/util/dateutil"
	"github.com/bitxx/evm-utils/util/httputil"
	"github.com/bitxx/evm-utils/util/idgenutil"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func TestSignEip721(t *testing.T) {
	loginUrl := "https://opside.network/api/user/custom/login"
	evmClient := evmutils.NewEthClient("", 0)
	privateKey := ""
	account, err := evmClient.AccountWithPrivateKey(privateKey)
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
	signature, err := SignEip721(privateKey, &typedData)
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

func TestMetamaskLoginSign(t *testing.T) {
	url := "https://graphigo.prd.galaxy.eco/query"
	privateKey := ""
	//1. 获取账户
	evmClient := evmutils.NewEthClient("", 0)
	account, err := evmClient.AccountWithPrivateKey(privateKey)
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
	sign, err := MetamaskSignLogin(msg, privateKey)
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
