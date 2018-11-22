package main

import (
	"bytes"
	"encoding/gob"
)

type rpcData struct {
	funcName string
	args []interface{}
}

func encodeData(funcName string, args []interface{}) ([]byte, error) {
	data := rpcData{}
	data.funcName = funcName

	data.args = args

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(data); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func decodeData(data []byte) (funcName string, args []interface{}, err error) {
	var m_data rpcData

	var buf = bytes.NewBuffer(data)

	dec := gob.NewDecoder(buf)

	if err = dec.Decode(&m_data); err != nil {
		return
	}
	funcName = m_data.funcName
	args = m_data.args

	return
}
