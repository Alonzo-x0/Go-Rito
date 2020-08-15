package main

import (
	"log"
	"strings"
	"strconv"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	boosted "./isHeBoosted/lib"
	weather "./weather/lib"
	database "./userdb/lib"
	//foo "./youtube/lib"
	//"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"time"
	"database/sql"
	//"reflect"
)


func delete_empty (s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}


func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	//
	disect := m.BeforeDelete
	
	if disect != nil {
		time.Sleep(1 * time.Second)	

		if disect.Author.Bot == false {
			message := disect.Content + " sent by @" + disect.Author.String() + " was deleted."
			//log.Println(message)
			s.ChannelMessageSend(m.ChannelID, message)
		}
	}
}
	




func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var embedded discordgo.MessageEmbed
	//log.Println(m.Content)
	embedded.URL = "https://github.com/Alonzo-x0/Go-Rito"

	db, err := sql.Open("mysql", "killer:toor@tcp(127.0.0.1:3306)/discord")

	if err != nil {
		log.Println(err)
	}

	if m.Author.ID == s.State.User.ID {

		return
	}


	//Insert(db, "weather", "userid", "lkey", "157959951501885440", "335315")
	if strings.HasPrefix(m.Content, "!set") {
		

		s.ChannelMessageSend(m.ChannelID, "Hol' up")

		args := strings.SplitAfter(m.Content, " ")

		//len = 2 for !set zipcode

		if len(args) == 2 {
			//convert zipcode from argument into lkey
			if len(args[1]) != 5 {
				s.ChannelMessageSend(m.ChannelID, "Hey baka, zipcodes are 5 digits!")
				return
			}

			loKey, err := weather.PostalKey(args[1], WeatherKey)
			//TODO: add if err 


			err = database.Insert(db, "weather", "userid", "lkey", m.Author.ID, loKey)
			if err != nil {
				return
			}
	
			err = database.DeleteDupes(db, "weather", "userid", "id")
			if err != nil {
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Your location has been set!")
		}else if len(args) != 2 {
			s.ChannelMessageSend(m.ChannelID, "Hey baka, usage is !set zipcode")
		}
	}



	if strings.HasPrefix(m.Content, "!conditions") {
		s.ChannelMessageSend(m.ChannelID, "Hol' up")

		//SelectRows(db *sql.DB, tarCol string, table string, desCol string, value string )
		args := strings.SplitAfter(m.Content, " ")
		log.Println(len(args))
		if len(args) == 1 {

			loKey, err := database.SelectRows(db, "lkey", "weather", "userid", m.Author.ID)

			if err != nil {
				log.Println(err)
				return
			}
	
			if loKey != 0 {
				zip := strconv.Itoa(loKey)
				message, err := weather.CurrConditions(zip, WeatherKey)
				if err != nil {
					fmt.Println(err)
					return
				}
				s.ChannelMessageSend(m.ChannelID, message)
			}
		}
		if len(args) == 2 {
			log.Println("1")
			loKey, err := weather.PostalKey(args[1], WeatherKey)
	
			if err != nil {
				fmt.Println(err)
				
			}
			log.Println("2")
			message, err := weather.CurrConditions(loKey, WeatherKey)
			if err != nil {
				fmt.Println(err)
			}
			log.Println("3")
			log.Println(message)
			s.ChannelMessageSend(m.ChannelID, message)
		}else if len(args) != 1 && len(args) != 2 {
			s.ChannelMessageSend(m.ChannelID, "Hey baka, this command takes 1 argument, your zipcode, or leave blank if you already set your zipcode with the !set xxxxx command!")
		}
	}

	killmepls := func(s *discordgo.Session, y *discordgo.GuildMembersChunk) {

	count := 0

	for p, _ := range y.Presences{
		if y.Presences[p].Status == "online" {
			count++
		}
	}

	time.Sleep(1 * time.Second)
	testing := strconv.Itoa(count) + " homies are online"
	
	s.ChannelMessageSend(m.ChannelID, testing)
	
	}

	//!online
	if strings.HasPrefix(m.Content, "!online") == true {
		s.AddHandler(killmepls)

		f := s.RequestGuildMembers(m.GuildID, "", 0, true)
		
		if f != nil {
			//log.Println(err)
			
		}
	}

	//!time messageID
	if strings.HasPrefix(m.Content, "!time") == true {
		s.ChannelMessageSend(m.ChannelID, "Hol' up")

				

		args := strings.SplitAfter(m.Content, " ")
		t, err := discordgo.SnowflakeTimestamp(args[1])
		
		if err != nil {
			//log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "Error in getting time stamp")
			return
		}


		//t.String()
		stamp := t.Format("2006-01-02 15:04:05")


		mValue, err := s.ChannelMessage(m.ChannelID, args[1])

		if err != nil {
			//log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "Unknown Message, unable to be parsed")
			return
		}

		//log.Println(mValue.Content)
		

		embedded.Title = "Time Stamp"

		message := fmt.Sprintf("%s | was sent at %s EST", mValue.Content, stamp)
		embedded.Description = message
		
		s.ChannelMessageSendEmbed(m.ChannelID, &embedded)

		
	}

	//!search "booster" "boostee" range
	if strings.HasPrefix(m.Content, "!search") == true{
		s.ChannelMessageSend(m.ChannelID, "Hol' up")
		
		
		args := strings.SplitAfter(m.Content, "\"")


		for i := range args{
			args[i] = strings.TrimRight(args[i], "\"")
			args[i] = strings.TrimLeft(args[i], " ")
		}
		args = delete_empty(args)
	


		if len(args) == 4 {
			index, err := strconv.Atoi(args[3])
			if index == 0 {
				index = 10
			}
			if err != nil {
				log.Println(err)
				return
			}

			message, err := boosted.UsrSearch(args[1], args[2], index, LeagueKey)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error raised, double check those arguments are valid buddy")
			}
			s.ChannelMessageSend(m.ChannelID, message)

		}else if len(args) != 4 { 
			s.ChannelMessageSend(m.ChannelID, "Hey baka, usage is !search \"booster\" \"boostee\" range")
		}

	}
	//!spect "player"
	if strings.HasPrefix(m.Content, "!spect") == true{
		s.ChannelMessageSend(m.ChannelID, "Hol' up, Sir.")

		

		args := strings.SplitAfter(m.Content, "\"")


		for i := range args{
			args[i] = strings.TrimRight(args[i], "\"")
			args[i] = strings.TrimLeft(args[i], " ")
		}
		args = delete_empty(args)
		log.Println(len(args))
		time.Sleep(2 * time.Second)
		if len(args) == 2 {
			//log.Println(args)
			teamA, teamB, err := boosted.SpectGame(args[1], LeagueKey)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error raised, double check those arguments are valid buddy")
				log.Println(err)
				return
			}else if err == nil {
				message := "```" + teamA[0]+ "\t\t\t\t" + teamB[0] + "\n" + teamA[1] + "\t" + teamB[1] + "\n" + teamA[2] + "\t" + teamB[2] + "\n" + teamA[3] + "\t" + teamB[3] + "\n" + teamA[4]  +"\t" + teamB[4] + "```"
				s.ChannelMessageSend(m.ChannelID, message)
			}

		} else if len(args) != 2{
			s.ChannelMessageSend(m.ChannelID, "Hey baka, usage is !spect \"player\"")

			//s.ChannelMessageSend(m.ChannelID, boosted.SpectGame(args[1], key)[1])
		}
	}
}

	
var (
	DiscordKey string
	LeagueKey string
	WeatherKey string
)

func InitApp() (string, string, string) {
	err := godotenv.Load("killerkeys.env")
	if err != nil {
		log.Fatal(err)
	} 
	dkey := os.Getenv("DisKey")
	rkey := os.Getenv("APIkey")
	wkey := os.Getenv("WeatherKey")
	return dkey, rkey, wkey
}

func init() {
	DiscordKey, LeagueKey, WeatherKey = InitApp()
}



func main() {
	
	
	dg, err := discordgo.New("Bot " + DiscordKey)
	
	//log.Println(reflect.TypeOf(dg))
	if err != nil {
		fmt.Println(err)
		return 
	}

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	dg.AddHandler(messageCreate)
	
	

	dg.State.MaxMessageCount = 50
	discordgo.NewState()


	err1 := dg.Open()

	if err1 != nil {
		fmt.Println(err1)
		return 
	}

	fmt.Println("CTRL-C to exit")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
	//messageCreate()
}

