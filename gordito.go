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
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"text/tabwriter"
	"bytes"
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var b bytes.Buffer
	w := tabwriter.NewWriter(&b, 1, 1, 1, ' ', 0)
	
	
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
	if strings.HasPrefix(m.Content, "!spect") == true{
		s.ChannelMessageSend(m.ChannelID, "Hol' up, Sir.")

		log.Println(m.Content)

		args := strings.SplitAfter(m.Content, "\"")


		for i := range args{
			args[i] = strings.TrimRight(args[i], "\"")
			args[i] = strings.TrimLeft(args[i], " ")
		}
		args = delete_empty(args)

		if len(args) == 2 {
			log.Println(args)
			x, y := boosted.SpectGame(args[1], key)
			fmt.Fprintln(w,"```" + x[0]+ "\t\t\t\t" + y[0] + "\n" + x[1] + "\t" + y[1] + "\n" + x[2] + "\t" + y[2] + "\n" + x[3] + "\t" + y[3] + "\n" + x[4]  +"\t" + y[4] + "```")
			w.Flush()
			
			log.Println(b.String)
			s.ChannelMessageSend(m.ChannelID, b.String())

		} else if len(args) != 2{
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