package storage

type GeoRepository interface {
	Add(query, region, geoLat, geoLon string) error // Вставка в DB's
	Get(query string) ([]AddressData, error) // Получаем данные из базы или из Редиса
	CheckAvailability(query string) (bool, error) // Проверка наличая в базе and Cache
}

type GeoRepositoryDB interface {
	InsertSearchHistory(query string) (int, error) // Вставка строки поиска в Таблицу search_history и надо б кэш
	InsertAddress(region, geoLat, geoLon string) (int, error) // Вставка адреса, и координат в Таблицу address и кэш
	InsertHistorySearchAddress(searchHistoryID, addressID int) error // Вставка айд в Таблицу history_search_address
	SearchInHistory(query string) (bool, error) // Проверка наличая в базе 
	FindAddressByQueryAndHistory(query string) ([]AddressData, error) // Селекс по двум таблицам 

}

type Cacher interface {
    Set(key string, value []AddressData) error // Устанавливает запись в редис
    Get(key string) ([]AddressData, error) // получаем данные из редиса
	Check(query string) (bool, error)
}