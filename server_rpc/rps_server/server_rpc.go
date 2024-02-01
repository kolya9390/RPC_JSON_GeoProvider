package rpcserver

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/kolya9390/RPC_JSON_GeoProvider/server_rpc/app"
	"github.com/kolya9390/RPC_JSON_GeoProvider/server_rpc/config"
	servisgeo "github.com/kolya9390/RPC_JSON_GeoProvider/server_rpc/servis_geo"
	"github.com/kolya9390/RPC_JSON_GeoProvider/server_rpc/storage"
	_ "github.com/lib/pq"
)



type GeoService struct {
	geoProvider app.GeoProvider
}

func NewGeoServis() *GeoService{
	return &GeoService{}
}

func (gs *GeoService) StartServer(port string) error {

	config := config.NewAppConf("server_app/.env")
		log.Println(config)
	// Инициализация подключения к базе данных
	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name)

	db, err := sqlx.Open("postgres", connstr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}
	time.Sleep(time.Second * 3)
	// Проверка соединения с базой данных
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %s", err)
	}

	defer db.Close()

	postgresDB := storage.NewGeoRepositoryDB(db)

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",config.Cache.Address,config.Cache.Port),
	})

	defer redisClient.Close()

	cache := storage.NewGeoRedis(redisClient)
	storagDB := storage.NewGeoRepositoryProxy(*postgresDB, cache)
	sevisDAdata := servisgeo.NewDadataService(config.AuthorizationDADATA)
	gs.geoProvider = app.NewGeoProvider(storagDB, sevisDAdata)

	err = postgresDB.ConnectToDB()

	if err != nil {
		log.Printf("Error conect DB %s", err)
	}

	if err := rpc.Register(gs); err != nil {
		log.Printf("Error Registretions rpc %v", err)
		return err
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s",config.RPCServer.Port))
	if err != nil {
		log.Printf("Eroor Listen %v", err)
		return err
	}
	defer listen.Close()

	log.Printf("RPC типа %s сервер запущен и прослушивает порт :%s",config.RPCServer.Type ,config.RPCServer.Port)

	switch true {
	case config.RPCServer.Type == "json-rpc":
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal("Accept error:", err)
			}
	
			go jsonrpc.ServeConn(conn)
		}

	case config.RPCServer.Type == "rpc":
		rpc.Accept(listen)

	}

	return nil
}

func (gs *GeoService) AddressSearchRPC(query RequestAddressSearch, reply *[]*Address) error {
    addresses, err := gs.geoProvider.AddressSearch(query.Query)
    if err != nil {
        log.Printf("Error AddressSearch: %v", err)
        return err
    }

	for _,adres := range addresses{
		*reply = append(*reply, &Address{
			GeoLat: adres.GeoLat,
			GeoLon: adres.GeoLon,
			Result: adres.Result,
		})

	}

    return nil
}


func (gs *GeoService) AddressGeoCodeRPC(geocode RequestAddressGeocode, reply *[]*Address) error {
		addresses, err := gs.geoProvider.GeoCode(geocode.Lat,geocode.Lng)
		if err != nil {
			log.Printf("Error AddressGeoCode: %v", err)
			return err
		}
		// Просто присваиваем новое значение reply через косвенное разыменование
		for _,adres := range addresses{
			*reply = append(*reply, &Address{
				GeoLat: adres.GeoLat,
				GeoLon: adres.GeoLon,
				Result: adres.Result,
			})
	
		}
	
		return nil
	}