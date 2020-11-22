package main

import (
	"log"
	"strings"
	"strconv"
	"fmt"
	"os/exec"
	"os"
	"google.golang.org/api/option"
	"context"
	"os/signal"
	"syscall"
	//"reflect"
	//"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"time"
	"bufio"
	"encoding/binary"
	"layeh.com/gopus"
	"io"
    "google.golang.org/api/youtube/v3"
)
const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
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
	embedded.URL = "https://github.com/Alonzo-x0/Go-Rito"

	if m.Author.ID == s.State.User.ID {

		return
	}

	serverID := "690961298384486410"

	test := make(chan bool)
	if strings.HasPrefix(m.Content, "!1") {
		voiceConn, err := s.ChannelVoiceJoin("690961298384486410", "690961298892259421", true, false)
		if err != nil {
			log.Println(err)
			return
		}
		defer voiceConn.Close()

		ctx := context.Background()
		
		var query []string
		
		service, err := youtube.NewService(ctx, option.WithAPIKey(GoogleKey))
		

		args := strings.SplitAfter(m.Content, "!1 ")
		for x, y := range args {
			log.Println(x, y)
		}
		if err != nil {
			log.Println(err)
			return
		}
		
		query = append(query, "snippet")
	
		call := service.Search.List(query).Q(args[0]).MaxResults(1)

		response, err := call.Do()
		if err != nil {
			log.Println(err)
			return
		}
		//log.Println(response.Items)
		videos := make(map[string]string)
		

        // Iterate through each item and add it to the correct list.
        for _, item := range response.Items {
                switch item.Id.Kind {
                case "youtube#video":
                        videos[item.Id.VideoId] = item.Snippet.Title
                }
        }
        id, title := printIDs(videos)
        s.ChannelMessageSend(m.ChannelID, "Now loading! >>> " + title)
        log.Println(id, title)

        vi := new(VoiceInstance)
        voiceInstances[serverID] = vi
        
        vi.PlayAudioFile(voiceConn, "https://www.youtube.com/watch?v=YJVmu6yttiw", test)
        time.Sleep(35 * time.Second)
        //test <- true


	}
	
	if strings.HasPrefix(m.Content, "!stop"){
		log.Println(voiceInstances[serverID] != nil)
		voiceInstances[serverID].StopVideo()


	}
	}



func printIDs(matches map[string]string) (string, string){
        for id, title := range matches {
                return id, title
        }
        return "", ""
}
	
type VoiceInstance struct {
	serverID     string
	skip         bool
	stop         bool
	trackPlaying bool
}
	
func (vi *VoiceInstance) StopVideo() {
	vi.stop = true
}

func main() {

	DiscordKey := os.Getenv("DisKey")

	dg, err := discordgo.New("Bot " + DiscordKey)
	
	//log.Println(reflect.TypeOf(dg))
	if err != nil {
		fmt.Println(err)
		 
	}

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	dg.AddHandler(messageCreate)
	
	

	dg.State.MaxMessageCount = 50
	discordgo.NewState()


	err1 := dg.Open()
	if err1 != nil {
		fmt.Println(err1)
		 
	}
	
	fmt.Println("CTRL-C to exit")
	
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	
	dg.Close()
	//messageCreate()
}

func SendPCM(v *discordgo.VoiceConnection, pcm <- chan []int16) {
	if pcm == nil {
		log.Println("PCM chan is nil")
		
		
	}

	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)

	if err != nil {
		log.Println(err)
		
	}

	opusEncoder.SetBitrate(384000)

	for {
		recv, ok := <-pcm
		if !ok {
			fmt.Println("song ended, or error kekw")
			return
			
		}

		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			log.Println("error encoding receiving chan bytes")
			log.Println(err)
			
		}

		if v.OpusSend == nil {
			//try doing != nil later
			log.Println("Discordgo not ready for opus packets. %+v : %+v", v.Ready, v.OpusSend)
			
		}
		v.OpusSend <- opus
	}
}



func (vi *VoiceInstance) PlayAudioFile(v *discordgo.VoiceConnection, link string, closer <- chan bool) {
	youtubeDl := exec.Command("youtube-dl", "--no-color", "--audio-format", "best", "--audio-format", "opus", link, "-o", "-")
	youtubeOut, err := youtubeDl.StdoutPipe()
	if err != nil {
		log.Println(err)
	}

	err = youtubeDl.Start()
	if err != nil {
		log.Println(err)
		
	}

	ffmpegRun := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	ffmpegRun.Stdin = youtubeOut
	ffmpegRun.Stderr = os.Stderr

	ffmpegout, err := ffmpegRun.StdoutPipe()

	if err != nil {
		log.Println(err)
		
	}
	
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	if err := ffmpegRun.Start(); err != nil {
		log.Println(err)
		
	}

	//go func() {
		//log.Println("In goFunc")
		//<- closer
		//err = youtubeDl.Process.Kill()
		//err = ffmpegRun.Process.Kill()
	//}()
	
	if err := v.Speaking(true); err != nil {
		log.Println(err)
		
	}
	
	defer func () {
		if err := v.Speaking(false); err != nil {
			log.Println(err)
			
		}
	}()

	pcmChan := make(chan []int16, 2)
	defer close(pcmChan)

	go func() {
		SendPCM(v, pcmChan)
		
	}()
	
	for {
		audiobuf := make([]int16, frameSize*channels)
		err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
			
		}
		
		if err != nil {
			log.Println(err)
			
		}
		if vi.stop == true {
			ffmpegRun.Process.Kill()
			break
		}

		select{
		case pcmChan <- audiobuf:
		case <- closer:
			err = youtubeDl.Process.Kill()
			err = ffmpegRun.Process.Kill()
			return
			
		}
	}
}

var (
	opusEncoder      *gopus.Encoder
	onIndex          int
	onChannel        int
	DiscordKey string
	LeagueKey string
	WeatherKey string
	GoogleKey string
	client     *discordgo.Session
	vconn      *discordgo.VoiceConnection
	youtubeKey string
	discordKey string
	chans      string
	channel    string
	guildId    string
	limitVids  int64
	bitRate    int64
	voiceInstances = map[string]*VoiceInstance{}
)


func InitApp() (string, string, string, string) {
	err := godotenv.Load("C:/Users/Alonzo/Programming/Go-Rito/isHeBoosted/killerkeys.env")
	if err != nil {
		log.Fatal(err)
	} 
	dkey := os.Getenv("DisKey")
	rkey := os.Getenv("APIkey")
	wkey := os.Getenv("WeatherKey")
	gkey := os.Getenv("youtubeKey")
	return dkey, rkey, wkey, gkey
}

func init() {
	DiscordKey, LeagueKey, WeatherKey, GoogleKey = InitApp()
}
