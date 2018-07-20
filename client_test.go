package gofastcgi

import (
	"fmt"
	"os"
	"testing"
	"strconv"
	"strings"
)

func TestGet(t *testing.T) {
	for i := 0; i < 100 ; i++ {
		host := "127.0.0.1"
		port := 9000
		client, err := NewClient(host, port)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
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
		fmt.Println("-------------------------------------------")
	}
}

func TestPost(t *testing.T) {
	host := "127.0.0.1"
	port := 9000
	client, err := NewClient(host, port)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	reqParams := "name=zhangsan"

	env := make(map[string]string)
	env["REQUEST_METHOD"] = "POST"
	env["SCRIPT_FILENAME"] = "/usr/local/php/test/index.php"
	env["CONTENT_TYPE"] = "application/x-www-form-urlencoded"
	env["CONTENT_LENGTH"] = strconv.Itoa(strings.Count(reqParams,"")-1)

	reponse, err := client.Request(env, reqParams)
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	fmt.Println(reponse.GetContent())
}
