package redistruct

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/rs/xid"
)

// Redistruct manage connection and data persistence
type Redistruct struct {
	connection *redis.Client
	Prefix     string
}

// NewRedistruct create Redistruct object
func NewRedistruct(conn *redis.Client, prefix string) *Redistruct {
	return &Redistruct{
		Prefix:     prefix,
		connection: conn,
	}
}

// SaveStruct persist struct into redis and return a unique ID and a error
func (r *Redistruct) SaveStruct(obj interface{}) (string, error) {
	uid := xid.New()
	data, err := serialize(obj)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%s:%s", r.Prefix, uid.String())

	err = r.connection.Set(key, data, 0).Err()
	if err != nil {
		return "", err
	}

	return key, nil
}

// LoadStruct get struct from redis
func (r *Redistruct) LoadStruct(key string, obj interface{}) error {
	val, err := r.connection.Get(key).Result()
	if err != nil {
		return err
	}

	return deserialize(obj, []byte(val))
}

func serialize(obj interface{}) ([]byte, error) {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(obj)
	return data.Bytes(), err
}

func deserialize(obj interface{}, data []byte) error {
	dataBuffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(dataBuffer)
	err := dec.Decode(obj)
	return err
}
