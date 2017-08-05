package tail

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func runLogger(filename string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0700)
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(100 * time.Millisecond)
		_, err := f.WriteString("testing\n")
		if err != nil {
			panic(err)
		}
	}
}

func TestTail_Tail(t *testing.T) {
	go runLogger("test.log")
	go runLogger("test2.log")
	go runLogger("test3.log")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	feed := make(chan Message)

	go func() {
		fmt.Println("tailing")
		Tail(ctx, "*.log", feed)
		close(feed)
		fmt.Println("done")
	}()

	for line := range feed {
		fmt.Printf("%+v\n", line)
	}
}
