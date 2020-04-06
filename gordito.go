package main

import (
	"log"
	"flag"
	"strings"
	"strconv"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	boosted "./isHeBoosted/lib"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var results []string
	
	err := godotenv.Load("killerkeys.env")
	if err != nil {
		log.Fatal(err)
	} 
	key := os.Getenv("APIkey")

	if m.Author.ID == s.State.User.ID {
		return
	}


	if strings.HasPrefix(m.Content, "!search") == true{
		s.ChannelMessageSend(m.ChannelID, "Hol' up")
		log.Println(m.Content)
		
		results = strings.SplitAfter(m.Content, "h")
		args := strings.SplitAfter(results[1], "/")

		for i := range results {
			results[i] = strings.TrimRight(results[i], " ")
		}

		for i := range args {
			args[i] = strings.TrimRight(args[i], "/")
		}
		fmt.Println(args[0], args[1], args[2])
		fmt.Println(len(args))

		if len(args) == 3 {
			index, err := strconv.Atoi(args[2])
			if err != nil {
				log.Println(err)
				s.ChannelMessageSend(m.ChannelID, "Error with index range, using default value of 50")
				s.ChannelMessageSend(m.ChannelID, boosted.UsrSearch(args[0], args[1], 50, key))
			}
			s.ChannelMessageSend(m.ChannelID, boosted.UsrSearch(args[0], args[1], index, key))
		} else {
			s.ChannelMessageSend(m.ChannelID, "Error with index range, using default value of 50")
			s.ChannelMessageSend(m.ChannelID, boosted.UsrSearch(args[0], args[1], 50, key))
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
	
	if err != nil {
		fmt.Println(err)
		return 
	}

	dg.AddHandler(messageCreate)
	
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