package utils

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"hash/fnv"
	"log"
)

// Hash as
func Hash(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}

// DelSliceByStr 删除指定切片元素
func DelSliceByStr(array []string, elem string) []string {
	if array != nil && len(array) > 0 {
		target := array[:0]
		for _, v := range array {
			if v != elem {
				target = append(target, v)
			}
		}
		return target
	}
	return nil
}

func request(webUrl string) {
	client := resty.New()

	resp, err := client.R().Get(webUrl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response Info:", string(resp.Body()))
	fmt.Println("Status Code:", resp.StatusCode())
}
