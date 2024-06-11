# evm-utils
ethereum生态的常用工具，如钱包生成、账户余额查询、转账等等


## 使用方式
具体使用方式，先看根目录的`evmutils_test.go`吧，后面有空了再写


## ** 关于model/contract/中的使用说明
合约交互不比简单转账，它是有参数的，而且参数的编码也要严格遵守规范，`go-ethereum`自带的`abigen`可以根据合约的`ABI`文件，自动生成对应的Go文件。  
[go-ethereum相关程序官方下载入口](https://geth.ethereum.org/downloads)  

举例：  
假设要获取erc20的基本操作流程  
从[OpenZeppelin库](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v5.0.1/contracts/token/ERC20/IERC20.sol) 下载接口文件，然后放到Remix里编译，拷贝出`ABI`文件  
保存的`ABI`文件名是`IERC20.json`。  
运行命令  
```go
abigen --abi IERC20.json --type ERC20 --pkg erc20 --out erc20.go
```  
之后进一步操作erc20的golang代码即可。
