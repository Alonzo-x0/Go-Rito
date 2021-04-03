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
	tube "./youtube/lib"

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
			log.Println(err)
			return
		}
		stamp := t.Format("2006-01-02 15:04:05")

		mValue, err := s.ChannelMessage(m.ChannelID, args[1])

		if err != nil {
			//log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "Unknown Message, unable to be parsed")
			return
		}
		embedded.Title = "Time Stamp"
		//log.Println(mValue.Content)
		var field []*discordgo.MessageEmbedField

		embedded.Fields = field

		message := fmt.Sprintf("\"%s\" was sent at %s EST", mValue.Content, stamp)
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
	voiceInstances[m.GuildID] = vi
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

		args := strings.SplitAfter(m.Content, "\"")

		for i := range args {
			args[i] = strings.TrimRight(args[i], "\"")
			args[i] = strings.TrimLeft(args[i], " ")
		}
		args = deleteEmpty(args)
		if len(args) == 2 {
			//log.Println(args)
			teamA, teamB, err := boosted.SpectGame(args[1], LeagueKey)

			if teamA == nil {
				s.ChannelMessageSend(m.ChannelID, "Uh oh, cant seem to find that player in a current game")
				log.Println(err)
				return
			} else if err == nil && teamA != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, embedMatchup(s, m, teamA, teamB))
			}

		} else if len(args) != 2 {
			s.ChannelMessageSend(m.ChannelID, "Hey baka, usage is !spect \"player\"")

			//s.ChannelMessageSend(m.ChannelID, boosted.SpectGame(args[1], key)[1])
		}
	}

	if strings.HasPrefix(m.Content, "!play") {
		s.ChannelMessageSend(m.ChannelID, "Hol' up, Sir.")

		//if anything exists after !play is typed, itll save it as a arg
		args := strings.SplitAfter(m.Content, "!play")[1]

		if vi.trackPlaying {
			//maybe add currPlay on channelmessagesend below
			s.ChannelMessageSend(m.ChannelID, "A song is currently playing!")
			return

		} else if vi.trackPlaying == false {
			//find a way to iterate thru vi.queue
			//title := strings.SplitAfter(m.Content, "!play")[1]
			log.Println(len(vi.queue))

			//if nothing is in queue, play the title after !play command
			if len(vi.queue) == 0 {
				tube.Zoop(s, m, args)
				return
				//if vi queue isnt empty, play it before anything else
			} else if len(vi.queue) != 0 {
				for _, song := range vi.queue {

					tube.Zoop(s, m, song)
					//queuePlay(s, m, title)

				}
			}
			//now that everything that could be played, is played, play whatever was in the !play command
			if args != "" {
				tube.Zoop(s, m, args)
				return
			}
		}
	}

	if strings.HasPrefix(m.Content, "!test") {
		x := len(vi.queue)
		for i := 0; i < x; i++ {
			s.ChannelMessageSend(m.ChannelID, vi.queue[i])
			tube.Zoop(s, m, vi.queue[i])
			//log.Println(vi.queue[i])

		}
	}

	if strings.HasPrefix(m.Content, "!queue") {
		s.ChannelMessageSend(m.ChannelID, "Hol' up, Sir.")
		args := strings.SplitAfter(m.Content, "!queue")[1]
		vi.queue = append(vi.queue, args)
		log.Println(vi.queue)
		s.ChannelMessageSend(m.ChannelID, "Queued up!")
	}

	if strings.HasPrefix(m.Content, "!stop") {
		voiceInstances[m.GuildID].StopAudio()
	}

	if strings.HasPrefix(m.Content, "!clear") {
		vi.queue = nil
	}
}

func queuePlay(s *discordgo.Session, m *discordgo.MessageCreate, title string) {
	log.Println(len(""))
	if len(title) > len("") {
		tube.Zoop(s, m, title)
	} else if len(title) <= len("") {
		for index, song := range vi.queue {
			log.Println(index, song)
		}
	}
}

//serverID := "690961298384486410"
func embedMatchup(s *discordgo.Session, m *discordgo.MessageCreate, blue map[string]string, red map[string]string) *discordgo.MessageEmbed {

	var embedded discordgo.MessageEmbed
	var field []*discordgo.MessageEmbedField
	var lCol *discordgo.MessageEmbedField = new(discordgo.MessageEmbedField)
	var rCol *discordgo.MessageEmbedField = new(discordgo.MessageEmbedField)

	embedded.Title = "MATCHUP"

	//tabwriter to maintain proper space formatting
	buf := new(bytes.Buffer)
	bufNew := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 5, 0, 3, ' ', tabwriter.Debug)
	b := tabwriter.NewWriter(bufNew, 5, 0, 3, ' ', tabwriter.Debug)

	//lCol is left header section, not the title
	lCol.Name = "BLUE"
	lCol.Inline = true
	field = append(field, lCol)

	//rCol is right header section
	rCol.Name = "RED"
	rCol.Inline = true
	field = append(field, rCol)

	//since we got 2 maps, we gotta make one of them into a list to iterate thru both
	keys := make([]string, len(blue))
	i := 0
	for k := range blue {
		keys[i] = k
		i++
	}

	embedded.Fields = field

	i = 0
	//indexs through that list we made earlier, this allows us to iterate through both maps letting us get the variables on the same line.
	for playerA, championA := range red {
		kA := keys[i]
		vA := blue[kA]
		fmt.Fprintln(w, playerA+"\t "+championA) //+"\t"+vA+"\t"+kA)
		fmt.Fprintln(b, kA+"\t "+vA)
		i++
	}

	//flush writes to tabwriter "sealing the deal" to the io
	//formats bytes from tabwriter
	//TODO: fix these dumb ass variable names
	w.Flush()
	b.Flush()
	y := string(buf.Bytes())
	foo := string(bufNew.Bytes())
	//keeps block formatting and preserves tabwriter format
	lCol.Value = "```" + y + "```"
	rCol.Value = "```" + foo + "```"

	//

	//s.ChannelMessageSend(m.ChannelID, y)
	return &embedded
	//s.ChannelMessageSendEmbed(m.ChannelID, &embedded)
}

//StopAudio stops audio playing in the server
func (vi *VoiceInstance) StopAudio() {
	vi.stop = true
	vi.trackPlaying = false
	vi.stop = false
}

var (
	//DiscordKey discord api key
	DiscordKey string
	//LeagueKey riot api key
	LeagueKey string
	//WeatherKey accuweather api key
	WeatherKey string
	//vi creates a voice instance global instance
	vi *VoiceInstance
	//GoogleKey is for using youtube api
	GoogleKey string
	//voiceInstances is a map of instances
	voiceInstances = map[string]*VoiceInstance{}
)

//InitApp initialize variables in global state
func InitApp() (*VoiceInstance, string, string, string, string) {
	err := godotenv.Load("C:/Users/Alonzo/Programming/Go-Rito/isHeBoosted/killerkeys.env")
	if err != nil {
		log.Fatal(err)
	}
	var vi *VoiceInstance = new(VoiceInstance)
	dkey := os.Getenv("DisKey")
	rkey := os.Getenv("APIkey")
	wkey := os.Getenv("WeatherKey")
	gkey := os.Getenv("youtubeKey")
	return vi, dkey, rkey, wkey, gkey
}

func init() {
	vi, DiscordKey, LeagueKey, WeatherKey, GoogleKey = InitApp()

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

//VoiceInstance is going to hold live values of current voicce instance by bot
type VoiceInstance struct {
	serverID     string
	skip         bool
	stop         bool
	trackPlaying bool
	queue        []string
	curPlay      string
}
