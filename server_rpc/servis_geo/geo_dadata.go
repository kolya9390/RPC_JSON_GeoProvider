package servisgeo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/kolya9390/RPC_JSON_GeoProvider/server_rpc/config"
)

type DadataService interface {
	AddressSearch(input string) ([]Address, error)
	GeoCode(lat, lng string) ([]Address, error)
}



type DadataServiceImpl struct {
	client http.Client
	AuthorizationDADATA	config.AuthorizationDADATA
}



func NewDadataService(configTokenDa config.AuthorizationDADATA) DadataService {

	return &DadataServiceImpl{
		client: http.Client{},
		AuthorizationDADATA: configTokenDa,
		}
}


// Реализация методов Запроса к ДаДата
func (d *DadataServiceImpl) makeRequest(ctx context.Context, url, method, contentType string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", d.AuthorizationDADATA.ApiKeyValue))
	req.Header.Set("X-Secret", d.AuthorizationDADATA.SecretKeyValue)


	return d.client.Do(req)

}

// Реализация методов DadataService
// SearchAddress и GeocodeAddress
func(d *DadataServiceImpl) AddressSearch(input string) ([]Address, error){

	ctx := context.Background()
	// москва сухонская 11

	data := strings.NewReader(fmt.Sprintf(`[ "%s" ]`, input))


	url := "https://cleaner.dadata.ru/api/v1/clean/address"

	resp, err := d.makeRequest(ctx, url, "POST", "application/json", data)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к Dadata API: %w", err)
	}

	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respData respDadataAdres

	//log.Printf("%s",bodyText)
	err = json.Unmarshal(bodyText, &respData)
	if err != nil {
		log.Println("respData unmarshal",err)
		return nil, err
	}

	var result []Address

	for _,adres := range respData{
		address := Address{
			GeoLat: adres.GeoLat,
			GeoLon: adres.GeoLon,
			Result: adres.Region,
		}
		result = append(result, address)
	}

	return result,nil

}

func(d *DadataServiceImpl)	GeoCode(lat, lng string) ([]Address, error){

	ctx := context.Background()
	// москва сухонская 11

	data := strings.NewReader(fmt.Sprintf(`{ "lat": %s, "lon": %s }`, lat, lng))

	url := "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address"

	resp, err := d.makeRequest(ctx, url, "POST", "application/json", data)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к Dadata API: %w", err)
	}

	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respData responseDadataGeo

	err = json.Unmarshal(bodyText, &respData)
	if err != nil {
		return nil,err
	}

	var adresses []Address

	for _, suggestion := range respData.Suggestions {
		address := Address{
			GeoLat: suggestion.Data.GeoLat,
			GeoLon: suggestion.Data.GeoLon,
			Result: suggestion.Data.Result,
		}
		adresses = append(adresses, address)
	
	}

	return adresses,nil
}