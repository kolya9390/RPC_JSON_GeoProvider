package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	reverproxy "github.com/kolya9390/RPC_JSON_GeoProvider/client_Proxy/infrastructure/reverProxy"
	swaggerui "github.com/kolya9390/RPC_JSON_GeoProvider/client_Proxy/infrastructure/swaggerUI"
	rpcclient "github.com/kolya9390/RPC_JSON_GeoProvider/client_Proxy/servis/rpc_client"
)

func NewApiRouter( /*RPC CLIENT*/ ) http.Handler {
	r := chi.NewRouter()

	proxy := reverproxy.NewReverseProxy("hugo", "1313")

	r.Use(middleware.Logger)
	r.Use(proxy.ReverseProxy)

	//SwaggerUI
	r.Get("/swagger", swaggerui.SwaggerUI)

	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("./client_app/public"))).ServeHTTP(w, r)
	})

	// API
	r.Route("/api", func(r chi.Router) {
		//	geo := controllers.GeoController

		r.Post("/address/search", func(w http.ResponseWriter, r *http.Request) {
			var requestBody RequestAddressSearch

			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				log.Println("Decoder Body")
				return
			}

			client := rpcclient.NewGeoClient()

			addresses := client.SearchSer(rpcclient.RequestAddressSearch(requestBody))
			var adreses_resp []Address
			for _,adres := range addresses{
				adreses_resp = append(adreses_resp, Address(adres))
			}
			response := ResponseAddress{
				Addresses: adreses_resp,
			}

			// Конвертируйте объект ResponseAddress в JSON
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(jsonResponse)
			if err != nil {
				log.Println("Error writing JSON response:", err)
				return
			}

		})

		r.Post("/address/geocode", func(w http.ResponseWriter, r *http.Request) {
			var requestBody RequestAddressGeocode

			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				log.Println("Decoder Body")
				return
			}

			client := rpcclient.NewGeoClient()

			addresses := client.GeoCoder(rpcclient.RequestAddressGeocode(requestBody))
			var adreses_resp []Address
			for _,adres := range addresses{
				adreses_resp = append(adreses_resp, Address(adres))
			}
			response := ResponseAddress{
				Addresses: adreses_resp,
			}

			// Конвертируйте объект ResponseAddress в JSON
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(jsonResponse)
			if err != nil {
				log.Println("Error writing JSON response:", err)
				return
			}

		})

		// Group Adress
		r.Route("/address", func(r chi.Router) {

		})
	})

	return r
}
