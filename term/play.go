package term

import (
	"2-some-one-cli/plyaudio"
	"2-some-one-cli/plyvideo"
)

func Play(path, name, filetype string) {
	if filetype == "Audio" {
		plyaudio.PlayAudio(path, name)
	} else if filetype == "Video" {
		plyvideo.PlayVideo(path, name)
	}
}
