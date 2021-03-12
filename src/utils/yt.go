package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"

	youtube "github.com/kkdai/youtube/v2"
)

type VideoFormat struct {
	Quality       string
	IsAudio       bool
	AudioQuality  string
	ContainsAudio bool
	Stream        *http.Response
}
type VideoProp struct {
	Title      string
	Author     string
	Length     string
	FormatList youtube.FormatList
	Duration   string
	Client     youtube.Client
	Video      *youtube.Video
}

//GetVideoProps Get Props of a video.
func GetVideoProps(url string) (props VideoProp, err error) {

	URL := strings.Replace(url, "https://www.youtube.com/watch?v=", "", -1)
	client := youtube.Client{}

	video, err := client.GetVideo(URL)
	vdProps := VideoProp{}
	if err != nil {
		return vdProps, err
	}

	duration := fmt.Sprintf("%v", video.Duration)
	vdProps.Duration = duration
	vdProps.Title = video.Title
	vdProps.Author = video.Author
	vdProps.FormatList = video.Formats
	vdProps.Client = client
	vdProps.Video = video

	return vdProps, nil

}
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)

	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\r[EZYVD/Download] Downloaded: %v ... ", humanize.Bytes(wc.Total))

	return n, nil
}

type WriteCounter struct {
	Total uint64
}

func Download(video http.Response, Path string) {

	file, err := os.Create(Path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	counter := &WriteCounter{}
	_, err = io.Copy(file, io.TeeReader(video.Body, counter))
	if err != nil {
		panic(err)
	}
	fmt.Println(" Done ")
}

func GetAllFormats(client youtube.Client, video *youtube.Video, formats youtube.FormatList) ([]VideoFormat, []error) {
	var fs []VideoFormat

	var errs []error
	for _, format := range formats {
		if format.AudioQuality != "" && format.QualityLabel != "" {
			s, err := client.GetStream(video, &format)
			if err != nil {
				errs = append(errs, err)
			}
			f := VideoFormat{
				IsAudio:       false,
				ContainsAudio: true,
				Quality:       format.QualityLabel,
				AudioQuality:  format.AudioQuality,
				Stream:        s,
			}
			fs = append(fs, f)
		}
		if format.QualityLabel != "" && format.AudioQuality == "" {
			s, err := client.GetStream(video, &format)
			if err != nil {
				errs = append(errs, err)
			}
			f := VideoFormat{
				IsAudio:       false,
				ContainsAudio: false,
				Quality:       format.QualityLabel,
				Stream:        s,
			}
			fs = append(fs, f)
		}
		if format.QualityLabel == "" && format.AudioQuality != "" {
			s, err := client.GetStream(video, &format)
			if err != nil {
				errs = append(errs, err)
			}
			f := VideoFormat{
				IsAudio:      true,
				AudioQuality: format.AudioQuality,
				Stream:       s,
			}
			fs = append(fs, f)
		}
	}
	return fs, errs
}
func GetAllAvaliableQualitys(formats *youtube.FormatList) []string {
	sup := []string{"2160p", "1440p", "2160p60", "1440p60", "1080p60", "1080p", "720p60", "720p", "480p", "360p", "240p", "144p"}
	var AQ []string
	for _, quality := range sup {
		if formats.FindByQuality(quality) != nil {
			AQ = append(AQ, quality)
		}
	}
	return AQ
}
func RankAudio(client *youtube.Client, video *youtube.Video, formats youtube.FormatList) http.Response {
	fmt.Print("Selecionando AUDIO ... ")
	var selectedAudioStream *http.Response
	totalFindHigh := false
	totalFindMedium := false

	var lowStream *http.Response
	var mediumStream *http.Response
	var highStream *http.Response
	for _, format := range formats {
		if format.AudioQuality == "AUDIO_QUALITY_HIGH" {
			totalFindHigh = true
			s, _ := client.GetStream(video, &format)
			highStream = s

		}
		if format.AudioQuality == "AUDIO_QUALITY_MEDIUM" {
			totalFindMedium = true
			s, _ := client.GetStream(video, &format)
			mediumStream = s

		}
		if format.AudioQuality == "AUDIO_QUALITY_LOW" {
			s, _ := client.GetStream(video, &format)
			lowStream = s
		}
	}
	if totalFindHigh {
		selectedAudioStream = lowStream
		fmt.Println("Alto.")
	} else if totalFindMedium {
		selectedAudioStream = mediumStream
		fmt.Println("Medio.")
	} else {
		selectedAudioStream = highStream
		fmt.Println("Baixo.")
	}

	return *selectedAudioStream
}
