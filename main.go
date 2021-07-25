package main

import "time"
import "fmt"
import "syscall"
import "os"
import "io"
import "io/ioutil"
import "log"

func main() {
		err := syscall.Mkfifo("/tmp/the-pipe", 0666)
		if(nil != err) {
			fmt.Println("%v", err)
		}
		dataChannel := make(chan []byte)
		go func() {
			for {
				dataChannel <- pollChannelForData()
			}
		}()

		//--------------------------------------------------
		var a [4]int
		a[0] = 77
    fmt.Println("hello world")
		fmt.Println("bytes: %o", a)
		fmt.Println("dataChannel: %o", dataChannel)

		for{
			select {
			case theBytes := <-dataChannel:
				// We received telemetry
				fmt.Println("bytes", theBytes)
			default:
			}
			time.Sleep(100 * time.Millisecond)
		}

}

func pollChannelForData() []byte {
	dataPipe, err := os.OpenFile("/tmp/the-pipe", os.O_RDONLY, 0)
	if err != nil {
		log.Panic("failed to open pipe", err)
	}

	defer close(dataPipe)

	// When the write side closes, we get an EOF.
	bytes, err := ioutil.ReadAll(dataPipe)
	if err != nil {
		log.Panic("failed to read telemetry pipe", err)
	}

	return bytes
}

func close(thing io.Closer) {
	err := thing.Close()
	if err != nil {
		log.Println(err)
	}
}
