package main

import (
	"fastcgi-client/fastcgi"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main()  {
	host := "127.0.0.1"
	port := 9000
	client, err := fastcgi.New(host, port)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	reqParams := "name=zhangsan"

	env := make(map[string]string)
	env["REQUEST_METHOD"] = "POST"
	env["SCRIPT_FILENAME"] = "/usr/local/php7/test/index.php"
	env["CONTENT_TYPE"] = "application/x-www-form-urlencoded"
	env["CONTENT_LENGTH"] = strconv.Itoa(strings.Count(reqParams,"")-1)

	reponse, err := client.Request(env, reqParams)
	if err != nil {
		fmt.Printf("err: %v", err)
	}

	fmt.Println(reponse.GetContent())
}
