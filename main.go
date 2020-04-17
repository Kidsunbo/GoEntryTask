package main

import (
	"awesomeProject/HTTPServer"
	"awesomeProject/TCPServer"
	"sync"
)

func main(){

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		HTTPServer.Run()
		wg.Done()
	}()
	go func() {
		TCPServer.NewServer(12345).Run()
		wg.Done()
	}()


	wg.Wait()

}
