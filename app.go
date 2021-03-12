package main

import (
	"net/http"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	fmt "fmt"

	utils "yvd/src/utils"
)

func main() {
	var videoProps utils.VideoProp
	var sAudio bool
	var aDownload bool
	ffmpegInstalled := utils.Installed()
	debugMessage := widget.NewLabel("Ez Youtube Video Downloader!")
	yvd := app.New()
	window := yvd.NewWindow("EZYVD(Ez Youtube video Downloader)")
	entryLink := widget.NewEntry()
	downloadProgress := widget.NewProgressBar()
	downloadAudio := widget.NewCheck("Baixar Somente Audio", func(b bool) {
		sAudio = b
	})
	AQS := []string{"nil"}
	selectQuality := widget.NewSelectEntry(AQS)

	aDownload = false
	SaveFileDialog := dialog.NewFileSave(func(uc fyne.URIWriteCloser, e error) {
		downloadProgress.Show()
		if aDownload {
			debugMessage.SetText("Já esta baixando!")
			return
		}
		aDownload = true
		URL := uc.URI()
		path := fyne.URI(URL).Path()
		if path != "" {
			audioStream := utils.RankAudio(&videoProps.Client, videoProps.Video, videoProps.FormatList)
			var video http.Response
			ava, _ := utils.GetAllFormats(videoProps.Client, videoProps.Video, videoProps.FormatList)
			fmt.Println(selectQuality.Text)
			for _, item := range ava {
				if item.Quality == selectQuality.Text {
					video = *item.Stream
				}
			}
			downloadProgress.Max = 100

			downloadProgress.Value = 0

			debugMessage.SetText("Baixando audio...")
			audioPath := strings.Replace(path, ".mp4", ".mp3", -1)
			utils.Download(audioStream, audioPath)
			downloadProgress.SetValue(10)
			if !ffmpegInstalled || !sAudio {
				debugMessage.SetText("Baixando video...")
				downloadProgress.SetValue(50)
				fmt.Println("Baixando video...")
				utils.Download(video, "video.cache.mp4")
				downloadProgress.SetValue(75)
				debugMessage.SetText("Processando video...")
				os.Remove(path)

				utils.MergeFiles("video.cache.mp4", audioPath, path)
				os.Remove("video.cache.mp4")
				debugMessage.SetText("Pronto!")
				downloadProgress.Hide()
			} else {

				fmt.Println("Audio Baixado!")
				debugMessage.SetText("Audio Baixado!")
			}
		}
		aDownload = true
		downloadProgress.SetValue(100)
	}, window)
	propsLabel := widget.NewLabel("")
	propsLabel.Hide()
	entryLink.SetPlaceHolder("Digite o Link")
	downloadButton := widget.NewButton("Fazer Download", func() {
		SaveFileDialog.Show()
	})
	downloadProgress.Hide()
	container_messages := container.New(layout.NewBorderLayout(nil, debugMessage, nil, nil), debugMessage)

	container_download := container.New(layout.NewBorderLayout(selectQuality, downloadButton, downloadAudio, nil), selectQuality, downloadButton, downloadAudio, downloadProgress)
	if !ffmpegInstalled {
		selectQuality.Hide()
		debugMessage.SetText("AVISO: Você não está com o ffmpeg instalado por isso não será possivel baixar videos.")
		fmt.Println("AVISO: Você não está com o ffmpeg instalado por isso não será possivel baixar videos.")
	}
	container_download.Hide()
	searchButton := widget.NewButton("Pesquisar", func() {
		container_download.Hide()
		propsLabel.Hide()
		debugMessage.SetText("Pesquisando...")
		props, err := utils.GetVideoProps(entryLink.Text)
		if err != nil {

			debugMessage.SetText(fmt.Sprintf("Erro: %v", err))
		}

		videoProps = props
		if props.Title == "" {
			entryLink.SetText("Link Invalido!")

			debugMessage.SetText("404-Não Encontrado.")
			return
		}
		propsLabel.SetText(fmt.Sprintf("Titulo: %s\nAutor: %s\nDuração: %s", videoProps.Title, videoProps.Author, videoProps.Duration))
		propsLabel.Show()
		debugMessage.SetText("Encontrado!")
		container_download.Show()
		selectQuality.SetPlaceHolder("Selecione uma resolução")
		selectQuality.SetOptions(utils.GetAllAvaliableQualitys(&videoProps.FormatList))

	})

	container_main := container.New(layout.NewVBoxLayout(),
		entryLink,
		searchButton,
		propsLabel,
	)

	content := container.New(layout.NewBorderLayout(container_main, container_messages, nil, nil), container_main, container_download, container_messages)
	window.SetContent(content)
	window.Resize(fyne.NewSize(800, 300))
	window.SetFixedSize(true)
	window.ShowAndRun()
}
