package rpcclient

import (
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/joho/godotenv"
)

type GeoClient struct {
	client *rpc.Client
}

func NewGeoClient() *GeoClient {
	var client *rpc.Client
	env, err := godotenv.Read("client_app/.env")
	if err != nil {
		log.Fatal(err)
	}
	typeRPC:= env["RPC_PROTOCOL"]
	switch true{
	case  typeRPC == "rpc":
		client, err = rpc.Dial("tcp", "server_rpc:1234")
	if err != nil {
		log.Fatal("type-rpc",err)
	}
case typeRPC == "json-rpc":
	client, err = jsonrpc.Dial("tcp", "server_rpc:1234")
	if err != nil {
		log.Fatal("type-json",err)
	}

}
    return &GeoClient{client: client}
}

type Address struct {
    GeoLat string `json:"lat"`
    GeoLon string `json:"lon"`
    Result string `json:"result"`
}

func (gss *GeoClient) SearchSer(query RequestAddressSearch) []Address {
	
	var result []Address                                             
	err := gss.client.Call("GeoService.AddressSearchRPC", query, &result)


	if err != nil {
		log.Println(err)
	}

	log.Println("Результат поиска:", result)

	return result
}

func (gss *GeoClient) GeoCoder(geocode RequestAddressGeocode) []Address {

	var result []Address                                              // Инициализируйте переменную для результата
	err := gss.client.Call("GeoService.AddressGeoCodeRPC",geocode, &result) // Измените имя метода и передайте &result
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Результат поиска:", result)

	return result
}