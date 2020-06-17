package main

import (
	"fmt"
	"os"
	"net/http"
	"log"
	//"strconv"
	"io/ioutil"
	"github.com/joho/godotenv"
	"encoding/json"

)

func urlRequest(url string) []byte{
	//simple request function returns response body
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}

//
//postalKey("zipcode")
/* used to get the location key, necessary for accuweather endpoints*/
func postalKey(zip string, key string) string {
	var postal Postal
	var url = "http://dataservice.accuweather.com/locations/v1/postalcodes/search?apikey=" + key + "&q=" + zip + "&details=false"//%20HTTP/1.1"
	fmt.Println(url)
	err := json.Unmarshal(urlRequest(url), &postal)

	if err != nil {
		log.Println(err)
	}
	//fmt.Println(postal[0].ParentCity.Key)
	return postal[0].ParentCity.Key

}

func currConditions(loKey string, key string) string{
	var weather curWeather
	//var url = "http://dataservice.accuweather.com/currentconditions/v1/" + loKey + "?apikey=" + key// + "%20HTTP/1.1"

	jsonFile, err := os.Open("conditions.json")

	if err != nil {
		log.Println(err)
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(byteValue, &weather)

	if err != nil {
		log.Println(err)
	}

	s := fmt.Sprintf("%s %v", weather[0].WeatherText, weather[0].Temperature.Imperial.Value)
	fmt.Println(s)
	
	return ""
}



func main() {
	err := godotenv.Load("C:/Users/Alonzo/Programming/Go-Rito/killerkeys.env")
	if err != nil {
		log.Fatal(err)
	} 
	WeatherKey := os.Getenv("WeatherKey")

	loKey := ""
	currConditions(loKey, WeatherKey)
	}



	
type curWeather []struct {
	LocalObservationDateTime string      `json:"LocalObservationDateTime"`
	EpochTime                int         `json:"EpochTime"`
	WeatherText              string      `json:"WeatherText"`
	WeatherIcon              int         `json:"WeatherIcon"`
	HasPrecipitation         bool        `json:"HasPrecipitation"`
	PrecipitationType        interface{} `json:"PrecipitationType"`
	IsDayTime                bool        `json:"IsDayTime"`
	Temperature              struct {
		Metric struct {
			Value    float64 `json:"Value"`
			Unit     string  `json:"Unit"`
			UnitType int     `json:"UnitType"`
		} `json:"Metric"`
		Imperial struct {
			Value    float64    `json:"Value"`
			Unit     string `json:"Unit"`
			UnitType int    `json:"UnitType"`
		} `json:"Imperial"`
	} `json:"Temperature"`
	MobileLink string `json:"MobileLink"`
	Link       string `json:"Link"`
}


type Postal []struct {
	Version           int    `json:"Version"`
	Key               string `json:"Key"`
	Type              string `json:"Type"`
	Rank              int    `json:"Rank"`
	LocalizedName     string `json:"LocalizedName"`
	EnglishName       string `json:"EnglishName"`
	PrimaryPostalCode string `json:"PrimaryPostalCode"`
	Region            struct {
		ID            string `json:"ID"`
		LocalizedName string `json:"LocalizedName"`
		EnglishName   string `json:"EnglishName"`
	} `json:"Region"`
	Country struct {
		ID            string `json:"ID"`
		LocalizedName string `json:"LocalizedName"`
		EnglishName   string `json:"EnglishName"`
	} `json:"Country"`
	AdministrativeArea struct {
		ID            string `json:"ID"`
		LocalizedName string `json:"LocalizedName"`
		EnglishName   string `json:"EnglishName"`
		Level         int    `json:"Level"`
		LocalizedType string `json:"LocalizedType"`
		EnglishType   string `json:"EnglishType"`
		CountryID     string `json:"CountryID"`
	} `json:"AdministrativeArea"`
	GeoPosition struct {
		Latitude  float64 `json:"Latitude"`
		Longitude float64 `json:"Longitude"`
		Elevation struct {
			Metric struct {
				Value    float64 `json:"Value"`
				Unit     string  `json:"Unit"`
				UnitType int     `json:"UnitType"`
			} `json:"Metric"`
			Imperial struct {
				Value    float64 `json:"Value"`
				Unit     string  `json:"Unit"`
				UnitType int     `json:"UnitType"`
			} `json:"Imperial"`
		} `json:"Elevation"`
	} `json:"GeoPosition"`
	IsAlias    bool `json:"IsAlias"`
	ParentCity struct {
		Key           string `json:"Key"`
		LocalizedName string `json:"LocalizedName"`
		EnglishName   string `json:"EnglishName"`
	} `json:"ParentCity"`
	SupplementalAdminAreas []struct {
		Level         int    `json:"Level"`
		LocalizedName string `json:"LocalizedName"`
		EnglishName   string `json:"EnglishName"`
	} `json:"SupplementalAdminAreas"`
	DataSets []string `json:"DataSets"`
}