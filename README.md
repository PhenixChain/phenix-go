# Go library for easy usage of PhenixChain

# API
* [创建账户](#创建账户)  
* [查询](#查询)  
    * [账户](#账户)  
    * [交易](#交易)  
* [广播](#广播)  
    * [交易](#交易)  

## 创建账户
```
公钥前缀: pub
地址前缀: adr
```
参考: [bip39](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki)
[bip44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)
[bech32](https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki#Bech32)

## 查询
### 账户
```
http://120.132.120.245/abci_query?path="/store/acc/key"&data=<Address HexString>
Address HexString 如: 0x6163636f756e743a5ae01bd40c2c6ffff0f00fdead41e34edd74fcad
```
### 交易
```
http://120.132.120.245/tx?hash=<TX HASH>
TX HASH 如: 0xBB83B9A3A0D41CF0FAB1933F08CD6FD7000F28CB04AAEAD30FDF70BE466D3714
```

## 广播
### 交易
```
http://120.132.120.245/broadcast_tx_sync?tx=<Hex SignData>
```