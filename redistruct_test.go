package redistruct

import (
	"bytes"
	"testing"

	"github.com/go-redis/redis/v7"
)

type testObj struct {
	Nome  string
	Idade int
}

func Test_serialize(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"Teste serialização", args{&testObj{"Carlos", 19}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := serialize(tt.args.obj)
			t.Log(got)
		})
	}
}

func Test_deserialize(t *testing.T) {
	data, _ := serialize(&testObj{"Carlos", 19})

	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"Teste deserialize", args{data}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got testObj
			deserialize(&got, tt.args.data)
			t.Log(got)

			gotData, _ := serialize(got)
			t.Log(bytes.Compare(gotData, data))

			t.Log(data)
			t.Log(gotData)
		})
	}
}

func TestRedistruct_SaveStruct(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()

	conn := NewRedistruct(client, "obj")

	objToken, err := conn.SaveStruct(&testObj{
		Nome:  "Carlos",
		Idade: 20,
	})

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(objToken)
}

func TestRedistruct_LoadStruct(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()

	conn := NewRedistruct(client, "obj")

	var obj testObj
	err := conn.LoadStruct("obj:bqou9hmi6ljddi3p1umg", &obj)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(obj.Nome)
}
