package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/ethereum/go-ethereum/log"
	"myCode/public_blockchain/part7-network/util"
)

type Transaction struct {
	TxHash []byte
	//UTXO输入
	Vint []txInput
	//UTXO输出
	Vout []txOutput
}

func (t *Transaction) hash() {
	tBytes := t.Serialize()
	//加入随机数byte
	randomNumber := util.GenerateRealRandom()
	randomByte := util.Int64ToBytes(randomNumber)
	sumByte := bytes.Join([][]byte{tBytes, randomByte}, []byte(""))
	hashByte := sha256.Sum256(sumByte)
	t.TxHash = hashByte[:]
}

func (t *Transaction) hashSign() []byte {
	t.TxHash = nil
	tBytes := t.Serialize()
	//加入随机数byte
	hashByte := sha256.Sum256(tBytes)
	return hashByte[:]
}

// 将transaction序列化成[]byte
func (t *Transaction) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(t)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}


func (t *Transaction) getTransBytes() []byte {
	if t.TxHash == nil || t.Vint == nil || t.Vout == nil{
		log.Error("交易信息不完整，无法拼接成字节数组")
		return nil
	}
	transBytes:=[]byte{}
	transBytes = append(transBytes,t.TxHash...)
	for _,v := range t.Vint {
		transBytes = append(transBytes,v.TxHash...)
		transBytes = append(transBytes,util.Int64ToBytes(int64(v.Index))...)
		transBytes = append(transBytes,v.Signature...)
		transBytes = append(transBytes,v.PublicKey...)
	}
	for _,v := range t.Vout {
		transBytes = append(transBytes,util.Int64ToBytes(int64(v.Value))...)
		transBytes = append(transBytes,v.PublicKeyHash...)
	}
	return transBytes
}

func (t *Transaction) customCopy() Transaction {
	newVin := []txInput{}
	newVout := []txOutput{}
	for _, vin := range t.Vint {
		newVin = append(newVin, txInput{vin.TxHash, vin.Index, nil, nil})
	}
	for _, vout := range t.Vout {
		newVout = append(newVout, txOutput{vout.Value, vout.PublicKeyHash})
	}
	return Transaction{t.TxHash, newVin, newVout}
}

func isGenesisTransaction(tss []Transaction) bool {
	if tss != nil {
		if tss[0].Vint[0].Index == -1 {
			return true
		}
	}
	return false
}
