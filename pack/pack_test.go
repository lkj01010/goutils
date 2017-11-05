package pack

import (
	"net"
	"testing"
	"fmt"
	"io"
)

func sender(writer io.Writer) {
	for i := 0; i < 50; i++ {
		words := "{\"Id\":1,\"Name\":\"golang\",\"Message\":\"message\"}"
		writer.Write(Pack([]byte(words)))
	}
	fmt.Println("send over")
}

func handleReader(reader io.Reader) {
	tmpBuffer := make([]byte, 0)

	readCh := make(chan []byte, 16)
	go func(readCh chan []byte) {
		for {
			select {
			case data := <-readCh:
				fmt.Println("read data: ", string(data))
			}
		}
	}(readCh)

	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			fmt.Println("read error: ", err)
		}
		fmt.Println("reader bytes: ", buffer[:n])
		tmpBuffer = Unpack(append(tmpBuffer, buffer[:n] ...), readCh)
	}
}

func TestMalformedInput(t *testing.T) {
	cli, srv := net.Pipe()
	//go cli.Write([]byte(`{id:1}`)) // invalid json
	//jsonrpc.ServeConn(srv)                 // must return, not loop

	go sender(cli);

	handleReader(srv);

	endless := make(chan interface{}, 0)

	<- endless
}
