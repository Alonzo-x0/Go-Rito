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
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"time"
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

	
	err := godotenv.Load("killerkeys.env")
	if err != nil {
		log.Fatal(err)
	} 
	key := os.Getenv("APIkey")
	w := os.Getenv("WeatherKey")

	if m.Author.ID == s.State.User.ID {

		return
	}

	if strings.HasPrefix(m.Content, "!conditions") {
		s.ChannelMessageSend(m.ChannelID, "Hol' up")

		args := strings.SplitAfter(m.Content, " ")

		fmt.Println(args[0], args[1])
		time.Sleep(4 * time.Second)
		zipcode, err := weather.PostalKey(args[1], w)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		message, err := weather.CurrConditions(zipcode, w)
		if err != nil {
			fmt.Println(err)
			return
		}
		s.ChannelMessageSend(m.ChannelID, message)
	}

	killmepls := func(s *discordgo.Session, y *discordgo.GuildMembersChunk) {

	count := 0

	for p, _ := range y.Presences{
		if y.Presences[p].Status == "online" {
			count++
		}
	}

	time.Sleep(1 * time.Second)
	testing := strconv.Itoa(count) + " Homies are online"
	
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

			message, err := boosted.UsrSearch(args[1], args[2], index, key)
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
			teamA, teamB, err := boosted.SpectGame(args[1], key)

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

	



func main() {
	err := godotenv.Load("killerkeys.env")
	if err != nil {
		log.Fatal(err)
	} 
	discToken := os.Getenv("DisKey")

	dg, err := discordgo.New("Bot " + discToken)
	
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

