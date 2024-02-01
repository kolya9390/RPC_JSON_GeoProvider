package app

import (
	"fmt"
	"log"

	servisgeo "github.com/kolya9390/RPC_JSON_GeoProvider/server_rpc/servis_geo"
	"github.com/kolya9390/RPC_JSON_GeoProvider/server_rpc/storage"
)

type GeoProviderService struct {
	storeg storage.GeoRepository
	servis servisgeo.DadataService

}

func NewGeoProvider(storagDB storage.GeoRepository,servis servisgeo.DadataService)*GeoProviderService{
	return &GeoProviderService{storeg: storagDB,servis: servis}
}

func (gp *GeoProviderService) AddressSearch(input string) ([]*servisgeo.Address, error) {
    var result []*servisgeo.Address

    // Проверка кэша
    if ok, err := gp.storeg.CheckAvailability(input); ok {
		if err != nil{
			return nil,err
		}
        addresses, err := gp.storeg.Get(input)
        if err != nil {
            log.Printf("Get error: %s", err)
            return nil, err
        }

        for _, address := range addresses {
            result = append(result, &servisgeo.Address{
                GeoLat: address.GeoLat,
                GeoLon: address.GeoLon,
                Result: address.Region,
            })
        }
        return result, nil
    }

    // Если данные в кэше устарели или их нет, обращаемся к сервису Dadata
    respData, err := gp.servis.AddressSearch(input)
    if err != nil {
        log.Printf("AddressSearch error: %s", err)
        return nil, err
    }

    if len(respData) < 1{
        return nil,fmt.Errorf("address was not found")
    }
    // Сохранение данных в кэше
    err = gp.storeg.Add(input, respData[0].Result, respData[0].GeoLat, respData[0].GeoLon)
    if err != nil {
        log.Printf("Add error: %s", err)
    }

    for _, address := range respData {
        result = append(result, &address)
    }

    return result, nil
}


func (gp *GeoProviderService) GeoCode(lat, lng string) ([]*servisgeo.Address, error) {
	
    var result []*servisgeo.Address

    geocood := fmt.Sprintf("%s %s",lat,lng)

     // Проверка кэша
     if ok, err := gp.storeg.CheckAvailability(geocood); ok {
		if err != nil{
			return nil,err
		}
        addresses, err := gp.storeg.Get(geocood)
        if err != nil {
            log.Printf("Get error: %s", err)
            return nil, err
        }


        for _, address := range addresses {
            result = append(result, &servisgeo.Address{
                GeoLat: address.GeoLat,
                GeoLon: address.GeoLon,
                Result: address.Region,
            })
        }
        return result, nil
    }

     // Если данные в кэше устарели или их нет, обращаемся к сервису Dadata
     respData, err := gp.servis.GeoCode(lat,lng)
     if err != nil {
         log.Printf("AddressSearch error: %s", err)
         return nil, err
     }

     if len(respData) < 1{
        return nil,fmt.Errorf("address was not found")
    }
 
     // Сохранение данных в кэше
     err = gp.storeg.Add(geocood, respData[0].Result, respData[0].GeoLat, respData[0].GeoLon)
     if err != nil {
         log.Printf("Add error: %s", err)
     }
 
     for _, address := range respData {
         result = append(result, &address)
     }
 
     return result, nil
}