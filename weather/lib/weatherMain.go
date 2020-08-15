package lib

import (
    "fmt"
    "os"
    "net/http"
    "log"
    "io/ioutil"
    "github.com/joho/godotenv"
    "encoding/json"
    "time"
    "strconv"
    "errors"
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
func PostalKey(zip string, key string) (string, error) {
    //var postal Postal
    //var url = "http://dataservice.accuweather.com/locations/v1/postalcodes/search?apikey=" + key + "&q=" + zip + "&details=false"//%20HTTP/1.1"
    err := errors.New("")
    //err := json.Unmarshal(urlRequest(url), &postal)
//
    //if err != nil {
        //return "", err
    //}
    //fmt.Println(postal[0].ParentCity.Key)
   // return postal[0].ParentCity.Key, err
   return "335315", err
}

func CurrConditions(loKey string, key string) (string, error){
    var weather curWeather
    //var url = "http://dataservice.accuweather.com/currentconditions/v1/" + loKey + "?apikey=" + key// + "%20HTTP/1.1"

    jsonFile, err := os.Open("C:/Users/Alonzo/Programming/Go-Rito/weather/lib/conditions.json")

    if err != nil {
        return "", err
    }

    defer jsonFile.Close()

    byteValue, err := ioutil.ReadAll(jsonFile)

    if err != nil {
        return "", err
    }

    err = json.Unmarshal(byteValue, &weather)

    if err != nil {
        return "", err
    }
    x := fmt.Sprintf("%.1f", weather[0].Temperature.Imperial.Value)

    s := weather[0].WeatherText + " " + x + "F"
    
    return s, err
}

func fiveDay(lokey string, key string) (string, error) {
    var forecast Forcast
    var blank string
    var lead [2]string
    var highs []int
    var lows []int
    var days []string
    var descript []string

    jsonFile, err := os.Open("forecast.json")

    if err != nil {
        return blank, err
    }

    defer jsonFile.Close()

    byteValue, err := ioutil.ReadAll(jsonFile)

    if err != nil {
        return blank, err
    }

    err = json.Unmarshal(byteValue, &forecast)

    if err != nil {
        return blank, err
    }
    t, err := time.Parse("2006-01-02T15:04:05Z07:00", forecast.Headline.EffectiveDate)
    if err != nil {
        return blank, err
    }
    //fmt.Println(t.Weekday())
    lead[0] = t.Format("2006-01-02 15:04:05")
    lead[1] = forecast.Headline.Text 

    for x, _ := range forecast.DailyForecasts{
        highs = append(highs, forecast.DailyForecasts[x].Temperature.Maximum.Value)
        lows = append(lows, forecast.DailyForecasts[x].Temperature.Minimum.Value)
        descript = append(descript, forecast.DailyForecasts[x].Day.IconPhrase)
        t, err := time.Parse("2006-01-02T15:04:05Z07:00", forecast.DailyForecasts[x].Date)
        if err != nil {
            return blank, err
        }
        days = append(days, t.Weekday().String())


    }
    
    var f string
    for i := 0; i < 5; i++ {
       f = f + fmt.Sprintf("%v %d \\ %d %v\n", days[i], highs[i], lows[i], descript[i])
        
    }

    return f, err
}

func hourly(lokey string, key string) (string, error) {
    //url :=http://dataservice.accuweather.com/forecasts/v1/hourly/12hour/+ lokey + ?apikey= + weatherkey
    jsonFile, err := os.Open("hourly.json")
    var hour Hour
    if err != nil {
        return "", err
    }

    defer jsonFile.Close()

    byteValue, err := ioutil.ReadAll(jsonFile)

    if err != nil {
        return "", err
    }

    err = json.Unmarshal(byteValue, &hour)

    if err != nil {
        return "", err
    }
    var f string
    for i := 0; i < 5; i++ {
        t, err := time.Parse("2006-01-02T15:04:05Z07:00", hour[i].DateTime)
        if err != nil {
            return "", err
        }
        f = f + "@" + strconv.Itoa(t.Hour()) + " it will be: " + strconv.Itoa(hour[i].Temperature.Value) + "F\n"
    }
    fmt.Println(f)
    return "", err
}
func main() {
    err := godotenv.Load("C:/Users/Alonzo/Programming/Go-Rito/killerkeys.env")
    if err != nil {
        log.Println("Error pulling the keys: ", err)
    } 
    WeatherKey := os.Getenv("WeatherKey")
    //335315
    loKey := ""
    hourly(loKey, WeatherKey)
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


type Forcast struct {
    Headline       Headline         `json:"Headline"`
    DailyForecasts []DailyForecasts `json:"DailyForecasts"`
}
type Headline struct {
    EffectiveDate      string `json:"EffectiveDate"`
    EffectiveEpochDate int    `json:"EffectiveEpochDate"`
    Severity           int    `json:"Severity"`
    Text               string `json:"Text"`
    Category           string `json:"Category"`
    EndDate            string `json:"EndDate"`
    EndEpochDate       int    `json:"EndEpochDate"`
    MobileLink         string `json:"MobileLink"`
    Link               string `json:"Link"`
}
type Minimum struct {
    Value    int    `json:"Value"`
    Unit     string `json:"Unit"`
    UnitType int    `json:"UnitType"`
}
type Maximum struct {
    Value    int    `json:"Value"`
    Unit     string `json:"Unit"`
    UnitType int    `json:"UnitType"`
}
type Temperature struct {
    Minimum Minimum `json:"Minimum"`
    Maximum Maximum `json:"Maximum"`
}

type Night struct {
    Icon                   int    `json:"Icon"`
    IconPhrase             string `json:"IconPhrase"`
    HasPrecipitation       bool   `json:"HasPrecipitation"`
    PrecipitationType      string `json:"PrecipitationType"`
    PrecipitationIntensity string `json:"PrecipitationIntensity"`
}
type Day struct {
    Icon                   int    `json:"Icon"`
    IconPhrase             string `json:"IconPhrase"`
    HasPrecipitation       bool   `json:"HasPrecipitation"`
    PrecipitationType      string `json:"PrecipitationType"`
    PrecipitationIntensity string `json:"PrecipitationIntensity"`
}

type DailyForecasts struct {
    Date        string      `json:"Date"`
    EpochDate   int         `json:"EpochDate"`
    Temperature Temperature `json:"Temperature"`
    Sources     []string    `json:"Sources"`
    MobileLink  string      `json:"MobileLink"`
    Link        string      `json:"Link"`
    Day         Day         `json:"Day,omitempty"`
    Night       Night       `json:"Night,omitempty"`

}

    
type Hour []struct {
    DateTime         string `json:"DateTime"`
    EpochDateTime    int    `json:"EpochDateTime"`
    WeatherIcon      int    `json:"WeatherIcon"`
    IconPhrase       string `json:"IconPhrase"`
    HasPrecipitation bool   `json:"HasPrecipitation"`
    IsDaylight       bool   `json:"IsDaylight"`
    Temperature      struct {
        Value    int    `json:"Value"`
        Unit     string `json:"Unit"`
        UnitType int    `json:"UnitType"`
    } `json:"Temperature"`
    PrecipitationProbability int    `json:"PrecipitationProbability"`
    MobileLink               string `json:"MobileLink"`
    Link                     string `json:"Link"`
    PrecipitationType        string `json:"PrecipitationType,omitempty"`
    PrecipitationIntensity   string `json:"PrecipitationIntensity,omitempty"`
}