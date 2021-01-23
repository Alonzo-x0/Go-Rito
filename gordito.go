package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"

	boosted "./isHeBoosted/lib"
	database "./userdb/lib"
	weather "./weather/lib"

	//foo "./youtube/lib"
	//"github.com/bwmarrin/dgvoice"
	"database/sql"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	//"reflect"
)

func deleteEmpty(s []string) []string {
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

func setCommand(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB, args []string) {
	if len(args) == 2 {
		zip := args[1]
		lokey, err := weather.PostalKey(zip, WeatherKey)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Hol' up, that command didn't work too well. Try doing it correctly.")
			log.Println(err)
			return
		}

		err = database.Insert(db, "weather", "userid", "lkey", m.Author.ID, lokey)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error saving your zipcode")
			log.Println(err)
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Your zipcode has been saved, very cool!")

		err = database.DeleteDupes(db, "weather", "userid", "id")
		if err != nil {
			log.Println("Error deleting duplicates in DB")
		}
	} else if len(args) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Hey baka, usage is !set {zipcode}")
	}
}

func conditions(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB, args []string) string {
	//SelectRows(db *sql.DB, tarCol string, table string, desCol string, value string )
	message := ""
	if len(args) == 1 {

		loKey, err := database.SelectRows(db, "lkey", "weather", "userid", m.Author.ID)

		if err != nil {
			log.Println(err)
			return ""
		}

		if loKey != 0 {
			zip := strconv.Itoa(loKey)
			message, err := weather.CurrConditions(zip, WeatherKey)
			if err != nil {
				fmt.Println(err)
				return ""
			}
			s.ChannelMessageSend(m.ChannelID, message)
		}
	}
	if len(args) == 2 {
		loKey, err := weather.PostalKey(args[1], WeatherKey)

		if err != nil {
			fmt.Println(err)

		}
		message, err := weather.CurrConditions(loKey, WeatherKey)
		if err != nil {
			fmt.Println(err)
		}
		log.Println(message)
		s.ChannelMessageSend(m.ChannelID, message)
	} else if len(args) != 1 && len(args) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Hey baka, this command takes 1 argument, your zipcode, or leave blank if you already set your zipcode with the !set xxxxx command!")
	}
	return message

}

func stampCheck(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	s.ChannelMessageSend(m.ChannelID, "Hol' up")
	var embedded discordgo.MessageEmbed

	if len(args) == 2 {
		args := strings.SplitAfter(m.Content, " ")
		t, err := discordgo.SnowflakeTimestamp(args[1])

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error in getting info on that timestamp, whoops.")
			return
		}
		stamp := t.Format("2006-01-02 15:04:05")

		mValue, err := s.ChannelMessage(m.ChannelID, args[1])

		if err != nil {
			//log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "Unknown Message, unable to be parsed")
			return
		}

		//log.Println(mValue.Content)

		embedded.Title = mValue.Content

		message := fmt.Sprintf("was sent at %s EST", stamp)
		embedded.Description = message

		s.ChannelMessageSendEmbed(m.ChannelID, &embedded)
	} else if len(args) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Hey baka, usage is !time {messageID}")

	}

}

func searchBoosted(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}
	args = deleteEmpty(args)

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

	} else if len(args) != 4 {
		s.ChannelMessageSend(m.ChannelID, "Hey baka, usage is !search \"booster\" \"boostee\" range")
	}
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//log.Println(m.Content)

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
		setCommand(s, m, db, args)
	}

	if strings.HasPrefix(m.Content, "!conditions") {
		s.ChannelMessageSend(m.ChannelID, "Hol' up")
		args := strings.SplitAfter(m.Content, " ")

		conditions(s, m, db, args)
	}

	//!time messageID
	if strings.HasPrefix(m.Content, "!time") == true {
		s.ChannelMessageSend(m.ChannelID, "Hol' up")
		args := strings.SplitAfter(m.Content, " ")

		stampCheck(s, m, args)

	}

	//!search "booster" "boostee" range
	if strings.HasPrefix(m.Content, "!search") == true {
		s.ChannelMessageSend(m.ChannelID, "Hol' up")
		args := strings.Split(m.Content, "\"")

		searchBoosted(s, m, args)

	}
	//!spect "player"
	if strings.HasPrefix(m.Content, "!spect") == true {
		buf := new(bytes.Buffer)
		w := tabwriter.NewWriter(buf, 5, 0, 3, ' ', tabwriter.Debug)

		s.ChannelMessageSend(m.ChannelID, "Hol' up, Sir.")

		args := strings.SplitAfter(m.Content, "\"")

		for i := range args {
			args[i] = strings.TrimRight(args[i], "\"")
			args[i] = strings.TrimLeft(args[i], " ")
		}
		args = deleteEmpty(args)
		if len(args) == 2 {
			//log.Println(args)
			teamA, teamB, err := boosted.SpectGame(args[1], LeagueKey)

			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error raised, double check those arguments are valid buddy")
				log.Println(err)
				return
			} else if err == nil || teamA != nil {
				keys := make([]string, len(teamA))
				i := 0
				for k := range teamA {
					keys[i] = k
					i++
				}
				i = 0
				for k, z := range teamB {
					kA := keys[i]
					vA := teamA[kA]
					//log.Println()
					fmt.Fprintln(w, k+"\t"+z+"\t"+vA+"\t"+kA)
					i++
				}

				w.Flush()
				y := string(string(buf.Bytes()))
				s.ChannelMessageSend(m.ChannelID, "```"+y+"```")

			} else if len(args) != 2 {
				s.ChannelMessageSend(m.ChannelID, "Hey baka, usage is !spect \"player\"")

				//s.ChannelMessageSend(m.ChannelID, boosted.SpectGame(args[1], key)[1])
			}
		}
	}
}

var (
	DiscordKey string
	LeagueKey  string
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
