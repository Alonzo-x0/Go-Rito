package lib

import (
	"net/http"
	"strings"
	"fmt"  
	"log"
	"io/ioutil"
	"encoding/json"
	"strconv"
	//"github.com/joho/godotenv"
	//"os"
	"sort"
	"time"
	//"github.com/pkg/profile"
)


func champFind(ID int) string{
	welcome, err := UnmarshalWelcome(urlRequest("http://ddragon.leagueoflegends.com/cdn/10.8.1/data/en_US/champion.json"))
	if err != nil {
			log.Println(err) 
		}	
	for k, _ := range welcome.Data {
		if welcome.Data[k].Key == strconv.Itoa(ID) {
			return k
		}
	}

	return "FUCK UP"
}

func SpectGame(username string, key string) ([]string, []string) {
	var enID, v = getSummoner(username, key)
	var errors []string
	errors = append(errors,  "Error in SpectGame")
	if v == "" {
		errors = append(errors, "")
		return errors, errors
	}
	var url = "https://na1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/" + enID + "?api_key=" + key
	var spect specMode
	var whoops errSpec
	var teamA []string
	var teamB []string
	var champA map[string]int
	var champB map[string]int
	var matchup map[string]string
	var matchupB map[string]string
	var error []string

	err := json.Unmarshal(urlRequest(url), &whoops)

	
	if err != nil {
		log.Println(err)
		error := append(error, "Error in SpectGame")
		return error, error

	}
	fmt.Println(whoops.Status.StatusCode == 0)
	if whoops.Status.StatusCode == 0 {

		err := json.Unmarshal(urlRequest(url), &spect)
	
		if err != nil {
			log.Println(err)
			error := append(error, "Error in SpectGame")
			return error, error
		}
	
	
		matchup = make(map[string]string)
		champA = make(map[string]int)
		champB = make(map[string]int)
		matchupB = make(map[string]string)
	
		for index, _ := range spect.Participants {
			if spect.Participants[index].TeamID == 100 {
				teamA = append(teamA, spect.Participants[index].SummonerName)
			}else if spect.Participants[index].TeamID == 200 {
				teamB = append(teamB, spect.Participants[index].SummonerName)
		} }
	
		for index, names := range teamA {
			if spect.Participants[index].SummonerName == names {
				champA[spect.Participants[index].SummonerName] = spect.Participants[index].ChampionID
			}
		}
	
		for i, x := range champA {
			matchup[i] = champFind(x)
		}
		var names []string
		for k, _ := range matchup {
			names = append(names, k)
		} 
		var summLen []int
		for _, y := range names {
			summLen = append(summLen, len(y))
		}
	
		sort.Sort(sort.Reverse(sort.IntSlice(summLen)))
	
		var aTeam []string
	 
		for _, y := range names {
			if len(y) != summLen[0] {
				ghost := summLen[0]-len(y)
				for k, v := range matchup {
					if y == k {
						aTeam = append(aTeam, y + strings.Repeat(" ", ghost) + " | " + v)
	
					}
				}
			} else if len(y) == summLen[0] {
				for k, v := range matchup {
					if y == k {
						aTeam = append(aTeam, y + " | " + v)
					}
				}
			}
		}
	
	
		for index, names := range teamB {
	
			if spect.Participants[index+5].SummonerName == names {
				////fmt.Println(names)
				champB[spect.Participants[index+5].SummonerName] = spect.Participants[index+5].ChampionID
			}
		}
	
		var namesB []string
	
		for i, x := range champB {
			////fmt.Println(i)
			matchupB[i] = champFind(x)
		}
		for k, _ := range matchupB {
			namesB = append(namesB, k)
		} 
		var fuckerB []int
		for _, y := range namesB {
			fuckerB = append(fuckerB, len(y))
		}
	
		sort.Sort(sort.Reverse(sort.IntSlice(fuckerB)))
	
		var bTeam []string
		for _, y := range namesB {
	
			if len(y) != fuckerB[0] {
				ghost := fuckerB[0]-len(y)
				for k, v := range matchupB {
					if y == k {
						bTeam = append(bTeam, y + strings.Repeat(" ", ghost) + " | " + v)
					}
				}
			} else if len(y) == fuckerB[0] {
				for k, v := range matchupB {
					if y == k {
						bTeam = append(bTeam, y + " | " + v)
					}
				}
			}
		}
	 
		return aTeam, bTeam
	} else if whoops.Status.StatusCode != 0{
		log.Println(whoops.Status.Message)
		return  errors, errors
	}
	
	return errors, errors
}

func checkTeam(gameList int64, key string) ([]string, []string){
	var url = "https://na1.api.riotgames.com/lol/match/v4/matches/" + strconv.FormatInt(gameList, 10) + "?api_key=" + key
	var matchinfo matchINFO
	var loserID []int
	var winnerID []int
	var winnerName []string
	var loserName [] string
	//var loseID  []int
	err := json.Unmarshal(urlRequest(url), &matchinfo)
	if err != nil {
		log.Println(err)
	}
	//fmt.Println(matchinfo.Participants[0].Stats.Win)

	for x, _ := range matchinfo.Participants {
		if matchinfo.Participants[x].Stats.Win == false{
			//fmt.Println(matchinfo.ParticipantIdentities[x].Player.SummonerName, "  is a loser")
			loserID = append(loserID, matchinfo.ParticipantIdentities[x].ParticipantID)
		}else {
			//fmt.Println(matchinfo.ParticipantIdentities[x].Player.SummonerName, "  is a winner")
			winnerID = append(winnerID, matchinfo.ParticipantIdentities[x].ParticipantID)

		}
	}

	for i := 0; i < 10; i++ {
		for _, y := range winnerID {
			if matchinfo.ParticipantIdentities[i].ParticipantID == y{
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


func getSummoner(name string, key string) (string, string){
	//id=encryptedID accountid=encryptedaccountid
	var url = "https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + name + "?api_key=" + key
	var summoner Summoner
	var error errSpec
	fmt.Println(url)




	err := json.Unmarshal(urlRequest(url), &error)
	if err != nil {
		log.Println(err)
	}
	var erMessage = error.Status.Message

	if erMessage == ""{
		err := json.Unmarshal(urlRequest(url), &summoner)
		
		if err != nil {
			log.Println(err)
			return "Error, check logs", ""
		}

		enID := summoner.ID
		accID := summoner.AccountID
		return enID, accID
	}

	log.Println("Error in getSummoner: ", error.Status.StatusCode)
	return erMessage, ""
}


func freqCount(list []string) map[string]int{

	freqBind := make(map[string]int)

	for _, item := range list {

		_, exist := freqBind[item]

		if exist {
			freqBind[item] += 1
		}else {
			freqBind[item] = 1
		}
	}
	return freqBind
}



func gameCount(accID string, max int, key string) []int64{
	var IDs []int64
	//var ratio []string
	
	var url = "https://na1.api.riotgames.com/lol/match/v4/matchlists/by-account/" + accID + "?queue=420&endIndex="+ strconv.Itoa(max) + "&api_key=" + key
	var history matchHistory
	var index int
	

	err := json.Unmarshal(urlRequest(url), &history)
	if err != nil {
		log.Println(err)
	}
	//finds out max amount of recorded matches and uses the the max if the number provided is too high
	for x, _ := range history.Matches {
		index = x
	}

	for i:=0; i <=index; i++ {
		IDs = append(IDs, history.Matches[i].GameID)		
	}


	
	
	return IDs
} 

func checkWin(gameId int64, key string) string{
	var url = "https://na1.api.riotgames.com/lol/match/v4/matches/" + strconv.FormatInt(gameId, 10) + "?api_key=" + key
	var matchinfo matchINFO
	err := json.Unmarshal(urlRequest(url), &matchinfo)
	if err != nil {
		log.Println(err)
	}

	return matchinfo.Teams[0].Win
}

func getMatchID(enID string, key string) (string) {
	var url = "https://na1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/" + enID + "?api_key=" + key
	
	response := make(map[string]interface{})
	err := json.Unmarshal(urlRequest(url), &response)
	if err != nil {
		log.Println(err)
	}

	if response["status"] != nil{
		
		return response["status"].(map[string]interface{})["message"].(string)
	} else {
		
		return strconv.FormatFloat(response["gameId"].(float64), 'f', -1, 64)
	
	return "you should not be seeing this."
	}
}
 
func UsrSearch(booster string, boostee string, indexMax int, key string) string{
	var url string
	var matchinfo matchINFO
	recent := 0
	start := time.Now()

	var whoops, accID = getSummoner(boostee, key)
	

	if accID == "" {
		log.Println("fuckup in usrsearch")
		log.Println(whoops)
		return whoops
	}
	log.Println("fuckup in usrsearch")
	var gameList = gameCount(accID, indexMax, key)

	for _, p := range gameList {
		url = "https://na1.api.riotgames.com/lol/match/v4/matches/" + strconv.FormatInt(p, 10) + "?api_key=" + key
		err := json.Unmarshal(urlRequest(url), &matchinfo)
		if err != nil {
			log.Println(err)
			//return err
		}
		for i := 0; i < 10; i++ {
			if strings.ToLower(matchinfo.ParticipantIdentities[i].Player.SummonerName) == strings.ToLower(booster) {
				//fmt.Println(booster, " is in match: ", p)
				recent = recent+1	
				//fmt.Println(recent)
			}
		}
	}
	elapsed := time.Since(start)

	output := "Looking through " + boostee + " match history and found " + booster + " in " + strconv.Itoa(recent) + " out of " + strconv.Itoa(indexMax) + " games"
	log.Println(output)
	log.Println("usrsearch took ", elapsed)
	return output
} 

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
		Win                  string `json:"win"`
		TeamID               int    `json:"teamId"`
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
			Kills                           int  `json:"kills"`
			ParticipantID                   int  `json:"participantId"`
			Win                             bool `json:"win"`
			Deaths                          int  `json:"deaths"`

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
	Type    Type             `json:"type"`   
	Format  string           `json:"format"` 
	Version Version          `json:"version"`
	Data    map[string]Datum `json:"data"`   
}

type Datum struct {
	Version Version            `json:"version"`
	ID      string             `json:"id"`     
	Key     string             `json:"key"`    
	Name    string             `json:"name"`   
	Title   string             `json:"title"`  
	Blurb   string             `json:"blurb"`  
	Info    Info               `json:"info"`   
	Image   Image              `json:"image"`  
	Tags    []Tag              `json:"tags"`   
	Partype string             `json:"partype"`
	Stats   map[string]float64 `json:"stats"`  
}

type Image struct {
	Full   string `json:"full"`  
	Sprite Sprite `json:"sprite"`
	Group  Type   `json:"group"` 
	X      int64  `json:"x"`     
	Y      int64  `json:"y"`     
	W      int64  `json:"w"`     
	H      int64  `json:"h"`     
}

type Info struct {
	Attack     int64 `json:"attack"`    
	Defense    int64 `json:"defense"`   
	Magic      int64 `json:"magic"`     
	Difficulty int64 `json:"difficulty"`
}

type Type string
const (
	Champion Type = "champion"
)

type Sprite string
const (
	Champion0PNG Sprite = "champion0.png"
	Champion1PNG Sprite = "champion1.png"
	Champion2PNG Sprite = "champion2.png"
	Champion3PNG Sprite = "champion3.png"
	Champion4PNG Sprite = "champion4.png"
)

type Tag string
const (
	Assassin Tag = "Assassin"
	Fighter Tag = "Fighter"
	Mage Tag = "Mage"
	Marksman Tag = "Marksman"
	Support Tag = "Support"
	Tank Tag = "Tank"
)

type Version string
const (
	The1081 Version = "10.8.1"
)

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