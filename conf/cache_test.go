package conf

import (
	"fmt"
	"testing"
)

func TestNdbPutArray(t *testing.T) {
	//v, _ := NdbGet("ceshi2")
	//array := []string{"123", "456", "789"}
	//NdbPutAny("test1", &array)

	var newa []string

	NdbGetAny("test1", &newa)

	fmt.Println(newa)
}
