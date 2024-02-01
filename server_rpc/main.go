package main

import (
	"log"

	rpcserver "github.com/kolya9390/RPC_JSON_GeoProvider/server_rpc/rps_server"
)

func main() {

	rpcServer := rpcserver.NewGeoServis()
	err := rpcServer.StartServer("1234")

	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

}
