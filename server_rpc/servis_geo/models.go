package servisgeo

type respDadataAdres []struct {
	Region string `json:"region"`
	GeoLat string `json:"geo_lat"`
	GeoLon string `json:"geo_lon"`
}

type Address struct {
	GeoLat string `json:"lat"`
	GeoLon string `json:"lon"`
	Result string `json:"result"`
}

type responseDadataGeo struct {
	Suggestions []struct {
		Value string `json:"value"`
		Data  struct {
			GeoLat string `json:"geo_lat"`
			GeoLon string `json:"geo_lon"`
			Result string `json:"region_with_type"`
		} `json:"data"`
	} `json:"suggestions"`
}