package lib

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	//"strings"

	//"golang.org/x/text/cases"
	"google.golang.org/api/option"

	//"reflect"
	//"github.com/bwmarrin/dgvoice"
	"bufio"
	"encoding/binary"
	"io"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"google.golang.org/api/youtube/v3"
	"layeh.com/gopus"
)

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data

)

//deleteEmpty deletes empty elements in a slice
func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

//Zoop is the main function to play music to whatever channel the command was sent from
func Zoop(s *discordgo.Session, m *discordgo.MessageCreate, title string) {
	voiceConn, err := s.ChannelVoiceJoin("690961298384486410", "690961298892259421", true, false)
	vi.stop = false

	if err != nil && vi.trackPlaying != false {
		return
	}

	defer voiceConn.Close()
	ctx := context.Background()

	var query []string
	closer := make(chan bool)

	service, err := youtube.NewService(ctx, option.WithAPIKey(GoogleKey))

	if err != nil {
		log.Println(err)
		return
	}

	if vi.trackPlaying == false {

		query = append(query, "snippet")

		log.Println(title)
		call := service.Search.List(query).Q(title).MaxResults(1)

		response, err := call.Do()
		if err != nil {
			log.Println(err)
			return
		}
		//log.Println(response.Items)
		videos := make(map[string]string)

		// Iterate through each item and add it to the correct list.
		for _, item := range response.Items {
			//log.Println(x, item)
			switch item.Id.Kind {
			case "youtube#video":
				videos[item.Id.VideoId] = item.Snippet.Title
			}
		}
		id, title := printIDs(videos)
		s.ChannelMessageSend(m.ChannelID, "Now loading! >>> "+title)
		//TODO FUCK WITH PLAYAUDIOFILE TO PROPERLY NOTIFY AND STOP STREAM
		retard := vi.PlayAudioFile(voiceConn, "https://www.youtube.com/watch?v="+id, closer)
		log.Println("TRACK PLAYING: ", vi.trackPlaying)
		if retard == "" {
			vi.StopVideo()
			err := voiceConn.Disconnect()
			if err != nil {
				log.Println(err)
			}
		}

	}
}

//if strings.HasPrefix(m.Content, "!stop"){
//log.Println(voiceInstances[serverID] != nil)
//vi.stop = true
//
//
////}
//if strings.HasPrefix(m.Content, "!stop") {
//voiceInstances[gID].StopVideo()
//}
//}

func printIDs(matches map[string]string) (string, string) {
	for id, title := range matches {
		return id, title
	}
	return "", ""
}

//VoiceInstance contains information regarding music and if its being played or not
type VoiceInstance struct {
	serverID     string
	skip         bool
	stop         bool
	trackPlaying bool
	queue        []string
	curPlay      string
}

//StopVideo stops the video
func (vi *VoiceInstance) StopVideo() {
	vi.stop = true
	vi.trackPlaying = false
}

//SendPCM sends pcm to channel for a continuous stream
func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16, ffmpegRun *exec.Cmd) {
	vi.trackPlaying = true
	if pcm == nil {
		log.Println("PCM chan is nil")

	}

	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)

	log.Println("track playing ", vi.trackPlaying)

	if err != nil {
		log.Println(err)

	}

	opusEncoder.SetBitrate(384000)

	for {
		if vi.stop == true {
			log.Println("here2")
			ffmpegRun.Process.Kill()
			//os.Exit(1)

			break
		}

		recv, ok := <-pcm
		if !ok {
			fmt.Println("song ended, or error kekw")
			vi.trackPlaying = false
			return

		}

		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			log.Println("error encoding receiving chan bytes")
			log.Println(err)

		}

		if v.OpusSend == nil {
			//try doing != nil later
			fmt.Printf("Discordgo not ready for opus packets. %v : %v", v.Ready, v.OpusSend)

		}
		v.OpusSend <- opus
	}
	vi.trackPlaying = false
	log.Println("track playing ", vi.trackPlaying)

}

//PlayAudioFile plays audio search by LINK url
func (vi *VoiceInstance) PlayAudioFile(v *discordgo.VoiceConnection, link string, closer <-chan bool) string {
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

	if err := v.Speaking(true); err != nil {
		log.Println(err)

	}

	defer func() {
		if err := v.Speaking(false); err != nil {
			log.Println(err)

		}
	}()

	pcmChan := make(chan []int16, 2)
	defer close(pcmChan)

	go func() {
		SendPCM(v, pcmChan, ffmpegRun)

	}()

	for {
		audiobuf := make([]int16, frameSize*channels)
		err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			log.Println("EOFOMEGGAAA")
			vi.trackPlaying = false

			return ""

		}

		if err != nil {
			log.Println(err)

		}

		select {
		case pcmChan <- audiobuf:
		case <-closer:
			log.Println("here4")
			err = youtubeDl.Process.Kill()
			err = ffmpegRun.Process.Kill()
			return "ERR"

		}

	}

	//if len(vi.queue) == 0 {
	//copy(vi.queue[0:], vi.queue[0+1:])
	//vi.queue[len(vi.queue)-1] = " "
	//vi.queue = vi.queue[:len(vi.queue)-1]
	//log.Println(vi.queue)

}

var (
	opusEncoder *gopus.Encoder
	onIndex     int
	onChannel   int
	//DiscordKey is my discord api key
	DiscordKey string
	//LeagueKey is riotAPI key
	LeagueKey string
	//WeatherKey is accuweather api key
	WeatherKey string
	//GoogleKey is google api key for youtube lookups
	GoogleKey      string
	client         *discordgo.Session
	vconn          *discordgo.VoiceConnection
	youtubeKey     string
	discordKey     string
	chans          string
	channel        string
	guildID        string
	limitVids      int64
	bitRate        int64
	voiceInstances = map[string]*VoiceInstance{}
	vi             *VoiceInstance
)

//InitApp sets initial variable values like keys
func InitApp() (*VoiceInstance, string, string, string, string) {
	err := godotenv.Load("C:/Users/Alonzo/Programming/Go-Rito/isHeBoosted/killerkeys.env")
	if err != nil {
		log.Fatal(err)
	}
	vi := new(VoiceInstance)
	dkey := os.Getenv("DisKey")
	rkey := os.Getenv("APIkey")
	wkey := os.Getenv("WeatherKey")
	gkey := os.Getenv("youtubeKey")
	return vi, dkey, rkey, wkey, gkey
}

func init() {
	vi, DiscordKey, LeagueKey, WeatherKey, GoogleKey = InitApp()
}
