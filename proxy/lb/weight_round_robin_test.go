package lb

import (
	"fmt"
	"testing"
)

func TestWeightRoundRobinBalance(t *testing.T) {
	fmt.Println("[-]TestWeightRoundRobinBalance start")
	rb := &WeightRoundRobinBalance{}
	rb.Add("127.0.0.1:8084", "4")
	rb.Add("127.0.0.1:8083", "3")
	rb.Add("127.0.0.1:8082", "2")

	for i := 0; i < 9; i++ {
		fmt.Println(rb.Next())
	}
	fmt.Println("[+]TestWeightRoundRobinBalance finished")
}
