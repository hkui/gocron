package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	client := http.Client{}
	data := ""
	url := ""
	request, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8") //设置Content-Type

	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36") //设置User-Agent
	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body)) //打印返回文本

}
