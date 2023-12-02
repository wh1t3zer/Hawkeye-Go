package lb

import (
	"fmt"
	"testing"
)

func TestRoundRobinBalance(t *testing.T) {
	fmt.Println("[-]TestRoundRobinBalance start")
	rb := &RoundRobinBalance{}
	for i := 0; i < 5; i++ {
		rb.Add(fmt.Sprintf("127.0.0.1:808%v", i))
	}

	for i := 0; i < 10; i++ {
		fmt.Println(rb.Next())
	}
	fmt.Println("[+]TestRoundRobinBalance finished")
}
