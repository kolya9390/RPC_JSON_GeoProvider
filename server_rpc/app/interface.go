package app

import servisgeo "github.com/kolya9390/RPC_JSON_GeoProvider/server_rpc/servis_geo"

type GeoProvider interface {
	AddressSearch(input string) ([]*servisgeo.Address, error)
	GeoCode(lat, lng string) ([]*servisgeo.Address, error)
}