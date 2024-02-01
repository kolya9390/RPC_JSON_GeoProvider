package rpcclient

type GeoSearchServicer interface {
	SearchSer(query string) []Address
	GeoCoder(lat,lng string) []Address
}

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}