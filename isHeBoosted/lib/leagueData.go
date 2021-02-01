package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func champFind(ID int) string {
	welcome, err := UnmarshalWelcome(urlRequest("http://ddragon.leagueoflegends.com/cdn/10.8.1/data/en_US/champion.json"))
	if err != nil {
		return ""
	}
	for k, _ := range welcome.Data {
		if welcome.Data[k].Key == strconv.Itoa(ID) {
			return k
		}
	}

	return ""
}

func SpectGame(username string, key string) (map[string]string, map[string]string, error) {
	var enID, _, err = getSummoner(username, key)

	if err != nil {
		return nil, nil, err
	}

	var url = "https://na1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/" + enID + "?api_key=" + key
	var spect specMode
	//var whoops errSpec
	//var teamT map[string]int
	var teamA, teamB []string
	var matchupA map[string]string
	var matchupB map[string]string

	matchupA = make(map[string]string)
	matchupB = make(map[string]string)

	err = json.Unmarshal(urlRequest(url), &spect)

	if err != nil {
		return nil, nil, err
	}

	if spect.GameMode == "" {
		return nil, nil, err
	}

	for index, _ := range spect.Participants {
		if spect.Participants[index].TeamID == 100 {
			teamA = append(teamA, spect.Participants[index].SummonerName)
		} else if spect.Participants[index].TeamID == 200 {
			teamB = append(teamB, spect.Participants[index].SummonerName)
		}
	}
	log.Println(url)
	for index, names := range teamA {
		if spect.Participants[index].SummonerName == names {
			//fmt.Fprintln(w, spect.Participants[index].SummonerName+" \t "+champFind(spect.Participants[index].ChampionID))
			matchupA[spect.Participants[index].SummonerName] = champFind(spect.Participants[index].ChampionID)
		}
	}

	for index, names := range teamB {
		if spect.Participants[index+5].SummonerName == names {

			//fmt.Fprintln(w, spect.Participants[index+5].SummonerName+" \t "+champFind(spect.Participants[index+5].ChampionID))

			matchupB[spect.Participants[index+5].SummonerName] = champFind(spect.Participants[index+5].ChampionID)
		}
	}

	return matchupA, matchupB, err

	//start
	//jsonFile, err := os.Open("C:/Users/Alonzo/Programming/Go-Rito/sampleData.json")
	//
	//if err != nil {
	//return nil, nil, err
	//}

	//defer jsonFile.Close()

	//byteValue, err := ioutil.ReadAll(jsonFile)

	//if err != nil {
	//return nil, nil, err
	//}
	//
	//err = json.Unmarshal(byteValue, &spect)
	//end
	//buf := new(bytes.Buffer)

	//w := tabwriter.NewWriter(buf, 5, 0, 3, '-', 0)

}

func checkTeam(gameList int64, key string) ([]string, []string) {
	var url = "https://na1.api.riotgames.com/lol/match/v4/matches/" + strconv.FormatInt(gameList, 10) + "?api_key=" + key
	var matchinfo matchINFO
	var loserID []int
	var winnerID []int
	var winnerName []string
	var loserName []string

	err := json.Unmarshal(urlRequest(url), &matchinfo)
	if err != nil {
		log.Println(err)
	}

	for x, _ := range matchinfo.Participants {
		if matchinfo.Participants[x].Stats.Win == false {
			loserID = append(loserID, matchinfo.ParticipantIdentities[x].ParticipantID)
		} else {
			winnerID = append(winnerID, matchinfo.ParticipantIdentities[x].ParticipantID)

		}
	}

	for i := 0; i < 10; i++ {
		for _, y := range winnerID {
			if matchinfo.ParticipantIdentities[i].ParticipantID == y {
				winnerName = append(winnerName, matchinfo.ParticipantIdentities[i].Player.SummonerName)
			} else {
				loserName = append(loserName, matchinfo.ParticipantIdentities[i].Player.SummonerName)
			}
		}
	}

	return winnerName, loserName

}

func dumpMap(space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			fmt.Printf("{ \"%v\": \n", k)
			dumpMap(space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}

func getSummoner(name string, key string) (string, string, error) {
	//id=encryptedID accountid=encryptedaccountid
	var url = "https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + name + "?api_key=" + key
	var summoner Summoner
	var error errSpec
	log.Println(url)

	err := json.Unmarshal(urlRequest(url), &error)
	if err != nil {
		return "Error", "", err
	}

	var erMessage = error.Status.Message

	if erMessage != "" {
		return "", "", err
	}

	err = json.Unmarshal(urlRequest(url), &summoner)

	if err != nil {
		return "", "", err
	}

	enID := summoner.ID
	accID := summoner.AccountID
	return enID, accID, err

}

func freqCount(list []string) map[string]int {

	freqBind := make(map[string]int)

	for _, item := range list {

		_, exist := freqBind[item]

		if exist {
			freqBind[item] += 1
		} else {
			freqBind[item] = 1
		}
	}
	return freqBind
}

func gameCount(accID string, max int, key string) []int64 {
	var IDs []int64
	var error errSpec
	var history matchHistory
	var index int
	var url = "https://na1.api.riotgames.com/lol/match/v4/matchlists/by-account/" + accID + "?queue=420&endIndex=" + strconv.Itoa(max) + "&api_key=" + key

	log.Println(url)

	err := json.Unmarshal(urlRequest(url), &error)
	if err != nil {
		log.Println("error unmarshalling to error")
		log.Println(err)
	}

	if error.Status.Message != "" {
		log.Println(error.Status.Message)

	} else if error.Status.Message == "" {

		err = json.Unmarshal(urlRequest(url), &history)
		if err != nil {
			log.Println("error parsing gameCount")
			log.Println(err)

		}
		//finds out max amount of recorded matches and uses the the max if the number provided is too high
		for x, _ := range history.Matches {
			index = x
		}
		for i := 0; i <= index; i++ {

			IDs = append(IDs, history.Matches[i].GameID)
		}

		return IDs
	}
	return IDs
}

func checkWin(gameId int64, key string) string {
	var url = "https://na1.api.riotgames.com/lol/match/v4/matches/" + strconv.FormatInt(gameId, 10) + "?api_key=" + key
	var matchinfo matchINFO
	err := json.Unmarshal(urlRequest(url), &matchinfo)
	if err != nil {
		log.Println(err)
	}

	return matchinfo.Teams[0].Win
}

func getMatchID(enID string, key string) string {
	var url = "https://na1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/" + enID + "?api_key=" + key

	response := make(map[string]interface{})
	err := json.Unmarshal(urlRequest(url), &response)
	if err != nil {
		log.Println(err)
	}

	if response["status"] != nil {

		return response["status"].(map[string]interface{})["message"].(string)
	} else {

		return strconv.FormatFloat(response["gameId"].(float64), 'f', -1, 64)

	}
}

func UsrSearch(booster string, boostee string, indexMax int, key string) (string, error) {
	var url string
	var matchinfo matchINFO
	recent := 0
	start := time.Now()

	_, accID, err := getSummoner(boostee, key)

	if err != nil {
		return "", err
	}
	var gameList = gameCount(accID, indexMax, key)

	for _, p := range gameList {
		url = "https://na1.api.riotgames.com/lol/match/v4/matches/" + strconv.FormatInt(p, 10) + "?api_key=" + key
		err := json.Unmarshal(urlRequest(url), &matchinfo)

		if err != nil {
			return "", err
		}
		for i := 0; i < 10; i++ {
			if strings.ToLower(matchinfo.ParticipantIdentities[i].Player.SummonerName) == strings.ToLower(booster) {
				recent = recent + 1
			}
		}
	}
	elapsed := time.Since(start)

	output := "Looking through " + boostee + " match history and found " + booster + " in " + strconv.Itoa(recent) + " out of " + strconv.Itoa(indexMax) + " games"
	log.Println(output)
	log.Println("usrsearch took ", elapsed)
	return output, err
}

func urlRequest(url string) []byte {
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

type matchINFO struct {
	SeasonID              int   `json:"seasonId"`
	QueueID               int   `json:"queueId"`
	GameID                int64 `json:"gameId"`
	ParticipantIdentities []struct {
		Player struct {
			CurrentPlatformID string `json:"currentPlatformId"`
			SummonerName      string `json:"summonerName"`
			MatchHistoryURI   string `json:"matchHistoryUri"`
			PlatformID        string `json:"platformId"`
			CurrentAccountID  string `json:"currentAccountId"`
			ProfileIcon       int    `json:"profileIcon"`
			SummonerID        string `json:"summonerId"`
			AccountID         string `json:"accountId"`
		} `json:"player"`
		ParticipantID int `json:"participantId"`
	} `json:"participantIdentities"`
	GameVersion string `json:"gameVersion"`
	PlatformID  string `json:"platformId"`
	GameMode    string `json:"gameMode"`
	MapID       int    `json:"mapId"`
	GameType    string `json:"gameType"`
	Teams       []struct {
		FirstDragon bool `json:"firstDragon"`
		Bans        []struct {
			PickTurn   int `json:"pickTurn"`
			ChampionID int `json:"championId"`
		} `json:"bans"`
		Win    string `json:"win"`
		TeamID int    `json:"teamId"`
	} `json:"teams"`
	Participants []struct {
		Spell1ID      int `json:"spell1Id"`
		ParticipantID int `json:"participantId"`
		Timeline      struct {
			Lane          string `json:"lane"`
			Role          string `json:"role"`
			ParticipantID int    `json:"participantId"`
		} `json:"timeline"`
		Spell2ID int `json:"spell2Id"`
		TeamID   int `json:"teamId"`
		Stats    struct {
			Kills         int  `json:"kills"`
			ParticipantID int  `json:"participantId"`
			Win           bool `json:"win"`
			Deaths        int  `json:"deaths"`
		} `json:"stats"`
		ChampionID int `json:"championId"`
	} `json:"participants"`
	GameDuration int   `json:"gameDuration"`
	GameCreation int64 `json:"gameCreation"`
}

type Summoner struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

type matchHistory struct {
	Matches []struct {
		Lane       string `json:"lane"`
		GameID     int64  `json:"gameId"`
		Champion   int    `json:"champion"`
		PlatformID string `json:"platformId"`
		Timestamp  int64  `json:"timestamp"`
		Queue      int    `json:"queue"`
		Role       string `json:"role"`
		Season     int    `json:"season"`
	} `json:"matches"`
	EndIndex   int `json:"endIndex"`
	StartIndex int `json:"startIndex"`
	TotalGames int `json:"totalGames"`
}

type errSpec struct {
	Status struct {
		Message    string `json:"message"`
		StatusCode int    `json:"status_code"`
	} `json:"status"`
}

func UnmarshalWelcome(data []byte) (Welcome, error) {
	var r Welcome
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Welcome) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Welcome struct {
	Format string           `json:"format"`
	Data   map[string]Datum `json:"data"`
}

type Datum struct {
	ID      string             `json:"id"`
	Key     string             `json:"key"`
	Name    string             `json:"name"`
	Title   string             `json:"title"`
	Blurb   string             `json:"blurb"`
	Partype string             `json:"partype"`
	Stats   map[string]float64 `json:"stats"`
}

type specMode struct {
	GameID            int64  `json:"gameId"`
	MapID             int    `json:"mapId"`
	GameMode          string `json:"gameMode"`
	GameType          string `json:"gameType"`
	GameQueueConfigID int    `json:"gameQueueConfigId"`
	Participants      []struct {
		TeamID                   int           `json:"teamId"`
		Spell1ID                 int           `json:"spell1Id"`
		Spell2ID                 int           `json:"spell2Id"`
		ChampionID               int           `json:"championId"`
		ProfileIconID            int           `json:"profileIconId"`
		SummonerName             string        `json:"summonerName"`
		Bot                      bool          `json:"bot"`
		SummonerID               string        `json:"summonerId"`
		GameCustomizationObjects []interface{} `json:"gameCustomizationObjects"`
		Perks                    struct {
			PerkIds      []int `json:"perkIds"`
			PerkStyle    int   `json:"perkStyle"`
			PerkSubStyle int   `json:"perkSubStyle"`
		} `json:"perks"`
	} `json:"participants"`
	Observers struct {
		EncryptionKey string `json:"encryptionKey"`
	} `json:"observers"`
	PlatformID      string `json:"platformId"`
	BannedChampions []struct {
		ChampionID int `json:"championId"`
		TeamID     int `json:"teamId"`
		PickTurn   int `json:"pickTurn"`
	} `json:"bannedChampions"`
	GameStartTime int64 `json:"gameStartTime"`
	GameLength    int   `json:"gameLength"`
}

type sampleData struct {
	GameID            int64  `json:"gameId"`
	MapID             int    `json:"mapId"`
	GameMode          string `json:"gameMode"`
	GameType          string `json:"gameType"`
	GameQueueConfigID int    `json:"gameQueueConfigId"`
	Participants      []struct {
		TeamID                   int           `json:"teamId"`
		Spell1ID                 int           `json:"spell1Id"`
		Spell2ID                 int           `json:"spell2Id"`
		ChampionID               int           `json:"championId"`
		ProfileIconID            int           `json:"profileIconId"`
		SummonerName             string        `json:"summonerName"`
		Bot                      bool          `json:"bot"`
		SummonerID               string        `json:"summonerId"`
		GameCustomizationObjects []interface{} `json:"gameCustomizationObjects"`
		Perks                    struct {
			PerkIds      []int `json:"perkIds"`
			PerkStyle    int   `json:"perkStyle"`
			PerkSubStyle int   `json:"perkSubStyle"`
		} `json:"perks"`
	} `json:"participants"`
	Observers struct {
		EncryptionKey string `json:"encryptionKey"`
	} `json:"observers"`
	PlatformID      string `json:"platformId"`
	BannedChampions []struct {
		ChampionID int `json:"championId"`
		TeamID     int `json:"teamId"`
		PickTurn   int `json:"pickTurn"`
	} `json:"bannedChampions"`
	GameStartTime int64 `json:"gameStartTime"`
	GameLength    int   `json:"gameLength"`
}
