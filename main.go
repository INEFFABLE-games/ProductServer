package main

import (
	"context"
	"main/src/server"
)

func main(){
	server.Start(context.Background())
}
