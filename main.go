package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"math/rand"
)

func main() {
	key := "a300b6e80af0a7c85b7a6713b1ffc652b70aaa77bec2150b76c3a2e99100921c"
	SendBundle(key, buildSignedTxs)
}

func GetSigner(pk string) (*ecdsa.PrivateKey, common.Address) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, common.Address{}
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println(fromAddress.String())
	return privateKey, common.HexToAddress(fromAddress.String())
}

func buildSignedTxs(pk *ecdsa.PrivateKey, addr common.Address, nonce uint64) []*types.Transaction {
	signedTxs := make([]*types.Transaction, 0)

	to := common.HexToAddress("0x0000000000000000000000000000000000001000")
	for i := 0; i < 1; i++ {
		fmt.Println("nonce:", nonce)
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: big.NewInt(6000000000),
			Gas:      7000000,
			To:       &to,
			Value:    big.NewInt(int64(rand.Intn(10))),
		})

		fmt.Println(tx.Hash().Hex())

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(97)), pk)
		if err != nil {
			panic(err)
		}

		signedTxs = append(signedTxs, signedTx)
		nonce++
	}

	return signedTxs
}

type buildTxs func(pk *ecdsa.PrivateKey, addr common.Address, nonce uint64) []*types.Transaction

func SendBundle(key string, fn buildTxs) {
	pk, addr := GetSigner(key)

	client, _ := ethclient.Dial("https://bsc-testnet-builder.bnbchain.org")
	nonce, err := client.NonceAt(context.TODO(), addr, nil)
	if err != nil {
		fmt.Println(err)
	}

	signedTxs := fn(pk, addr, nonce)

	txBytes := make([]hexutil.Bytes, 0)
	for _, signedTx := range signedTxs {
		txByte, err := signedTx.MarshalBinary()
		if err != nil {
			panic(err)
		}
		txBytes = append(txBytes, txByte)
	}

	bundleArgs := &types.SendBundleArgs{
		Txs: txBytes,
	}

	err = client.SendBundle(context.Background(), bundleArgs)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Bundle sent successfully")
	}
}
