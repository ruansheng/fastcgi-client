package gofastcgi

import (
	"fmt"
	"os"
	"testing"
	"time"
)

type Request struct {
	Rid     string                 `json:"rid"`
	Service string                 `json:"service"`
	Method  string                 `json:"method"`
	Args    map[string]interface{} `json:"args"`
}

func TestGet(t *testing.T) {
	host := "127.0.0.1"
	//host := "172.16.11.117"
	port := 9000

	client, err := NewClient(host, port, true)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	for i := 0; i < 2; i++ {
		reqParams := "name=zhangsan"
		env := make(map[string]string)
		env["REQUEST_METHOD"] = "GET"
		//env["SCRIPT_FILENAME"] = "/usr/local/php/test/index.php"
		//env["SCRIPT_FILENAME"] = "/Users/ruansheng/PhpstormProjects/sample/index_web.php"
		//env["SCRIPT_FILENAME"] = "/var/www/code/php-rpc/index_rpc.php"
		env["QUERY_STRING"] = reqParams

		reponse, err := client.Request(env, "")
		if err != nil {
			fmt.Printf("err: %v\n", err)
			fmt.Println("-------------------------------------------")
			continue
		}

		fmt.Println(reponse.GetContent())
		fmt.Println("-------------------------------------------")
		//client.Reset()
		//fmt.Printf("%+v\n", client)
		//time.Sleep(time.Second * 10)
	}

	time.Sleep(time.Second * 30)
}

func TestPost(t *testing.T) {
	/*
		//host := "127.0.0.1"
		host := "172.16.11.117"
		port := 9000
		client, err := NewClient(host, port, false)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}

		request := &Request{
			Rid:     "111",
			Service: "test",
			Method:  "demo",
			Args:    map[string]interface{}{"a": 1, "b": 2},
		}

		jsonStr, err := json.Marshal(request)
		if err != nil {
			fmt.Printf("err: %v", err)
		}

		reqParams := string(jsonStr)
		fmt.Println(reqParams)
		env := make(map[string]string)
		env["REQUEST_METHOD"] = "POST"
		//env["SCRIPT_FILENAME"] = "/Users/ruansheng/PhpstormProjects/php-rpc/index_rpc.php"
		env["SCRIPT_FILENAME"] = "/var/www/code/php-rpc/index_rpc.php"
		env["CONTENT_TYPE"] = "application/x-www-form-urlencoded"
		env["CONTENT_LENGTH"] = strconv.Itoa(strings.Count(reqParams, "") - 1)

		reponse, err := client.Request(env, reqParams)
		if err != nil {
			fmt.Printf("err: %v", err)
			os.Exit(-1)
		}

		fmt.Println(reponse.GetContent())
	*/
}
