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
	//weather "./weather/lib"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"text/tabwriter"
	"bytes"
	"time"
	"reflect"
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
			log.Println(disect.Content)
			log.Println(message)
			s.ChannelMessageSend(m.ChannelID, message)
		}
	}
	
}



func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 1, 1, 1, ' ', 0)
	
	err := godotenv.Load("killerkeys.env")
	if err != nil {
		log.Fatal(err)
	} 
	key := os.Getenv("APIkey")

	if m.Author.ID == s.State.User.ID {
		time.Sleep(3 * time.Minute)
		s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		log.Println("deletion here")
		return
	}

	var embedded discordgo.MessageEmbed
	embedded.URL = "https://github.com/Alonzo-x0/Go-Rito"

	killmepls := func(s *discordgo.Session, y *discordgo.GuildMembersChunk) {
	count := 0
	log.Println(y.Presences[0].Status)
	for p, _ := range y.Presences{
		if y.Presences[p].Status == "online" {
			count++
		}
	}
	log.Println(reflect.TypeOf(count))
	time.Sleep(1 * time.Second)
	testing := fmt.Sprintf("%i homies are online", count)
	
	s.ChannelMessageSend(m.ChannelID, testing)
	
	log.Println(count, " homies are online")	
	}
	if strings.HasPrefix(m.Content, "!test") == true {
		s.AddHandler(killmepls)
		f := s.RequestGuildMembers(m.GuildID, "", 0, true)
		
		if f != nil {
			log.Println(err)
			
		}
	}


	if strings.HasPrefix(m.Content, "!time") == true {
		

		s.ChannelMessageSend(m.ChannelID, "Hol' up")
		log.Println(m.Content, "HERE")
		args := strings.SplitAfter(m.Content, " ")
		t, _ := discordgo.SnowflakeTimestamp(args[1])
		
		t.String()
		stamp := t.Format("2006-01-02 15:04:05")


		mValue, _ := s.ChannelMessage(m.ChannelID, args[1])

		log.Println(mValue.Content)
		time.Sleep(5 & time.Second)

		embedded.Title = "Time Stamp"

		message := fmt.Sprintf("%s | was sent at %s EST", mValue.Content, stamp)
		embedded.Description = message
		
		s.ChannelMessageSendEmbed(m.ChannelID, &embedded)

		
	}

	//!search "booster" "boostee"
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
			}
			

			s.ChannelMessageSend(m.ChannelID, boosted.UsrSearch(args[1], args[2], index, key))

		}else if len(args) == 3 {
			

			s.ChannelMessageSend(m.ChannelID, boosted.UsrSearch(args[1], args[2], 50, key))
		}else {

			s.ChannelMessageSend(m.ChannelID, "Whoops, double check your request, something is amiss")
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

		if len(args) == 2 {
			log.Println(args)
			x, y := boosted.SpectGame(args[1], key)
			if y[0] != "" {

				fmt.Fprintln(w,"```" + x[0]+ "\t\t\t\t" + y[0] + "\n" + x[1] + "\t" + y[1] + "\n" + x[2] + "\t" + y[2] + "\n" + x[3] + "\t" + y[3] + "\n" + x[4]  +"\t" + y[4] + "```")
				w.Flush()
				
				log.Println(b.String)
				s.ChannelMessageSend(m.ChannelID, b.String())
			}

		} else {
			s.ChannelMessageSend(m.ChannelID, "Whoops, double check your request, something is amiss")

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
	
	log.Println(reflect.TypeOf(dg))
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

