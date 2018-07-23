package gofastcgi

import (
	"fmt"
	"testing"
	"time"
)

func TestNewClientPool(t *testing.T) {
	poolOptions := PoolOptions{
		host:        "127.0.0.1",
		port:        9000,
		maxOpen:     5,
		startOpen:   2,
		maxLifetime: time.Second * 6,
	}
	pool, err := NewClientPool(poolOptions)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("%+v\n", pool)
	for i := 0; i < 10; i++ {
		client, err := pool.Acquire()
		if err != nil {
			t.Fatal(err.Error())
		}
		fmt.Println(client)
		pool.Release(client)
		time.Sleep(time.Second * 3)
		fmt.Printf("%+v\n", pool)
	}

	c := make(chan int)
	<-c
}

func TestNewClientPool1(t *testing.T) {
	poolOptions := PoolOptions{
		host:        "127.0.0.1",
		port:        9000,
		maxOpen:     5,
		startOpen:   2,
		maxLifetime: time.Second * 6,
	}
	pool, err := NewClientPool(poolOptions)
	if err != nil {
		t.Fatal(err.Error())
	}

	for i := 0; i < 10; i++ {
		client, err := pool.Acquire()
		if err != nil {
			t.Fatal(err.Error())
		}

		reqParams := "name=zhangsan"
		env := make(map[string]string)
		env["REQUEST_METHOD"] = "GET"
		env["SCRIPT_FILENAME"] = "/usr/local/php/test/index.php"
		env["QUERY_STRING"] = reqParams

		reponse, err := client.Request(env, "")
		if err != nil {
			fmt.Printf("err: %v\n", err)
			fmt.Println("-------------------------------------------")
			continue
		}

		fmt.Println(reponse.GetContent())
		fmt.Println(client)
		fmt.Printf("%+v\n", pool)
		fmt.Println("-------------------------------------------")
		pool.Release(client)
		time.Sleep(time.Second * 3)
	}

	c := make(chan int)
	<-c
}
