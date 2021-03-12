package utils

import (
	"fmt"
	system "os/exec"
	"strings"
)

//MergeFiles Merge Audio and Video to one file.
func MergeFiles(videoPath string, audioPath string, outputFilePath string) (Exit []byte, err error) {
	command := fmt.Sprintf("ffmpeg -i %s -i %s -c copy %s", videoPath, audioPath, outputFilePath)
	args := strings.Split(command, " ")
	result := system.Command(args[0], args[1:]...)
	exit, err := result.CombinedOutput()
	fmt.Println(string(exit))
	if err != nil {
		panic(err)
	}
	return exit, err
}
func Installed() bool {
	r := system.Command("ffmpeg", "-version")
	_, err := r.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return true
	}
}
