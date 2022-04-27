package plyvideo

import (
	"2-some-one-cli/util"
	"2-some-one-cli/wsclient"
	vlc "github.com/adrg/libvlc-go"
	"log"
	"strconv"
	"time"
)

type VlcPlayer struct {
	VP *vlc.Player
}

func NewPlayer() (*VlcPlayer, error) {
	err := vlc.Init("--quiet", "--no-xlib")
	util.AssertErr(err)

	// Create a new player.
	player, err := vlc.NewPlayer()
	util.AssertErr(err)
	player.SetVolume(80)

	return &VlcPlayer{VP: player}, nil
}

func (p VlcPlayer) Play() (err error) {
	err = p.VP.Play()
	util.AssertErr(err)
	return nil
}

func (p VlcPlayer) Close() (err error) {
	err = p.VP.Stop()
	util.AssertErr(err)
	return nil
}

func (p VlcPlayer) UploadTime(wsc *wsclient.WsClient) {
	nowTime := time.Now().UnixNano() / 1e6
	sec, err := p.VP.MediaTime()
	util.AssertErr(err)
	var msg wsclient.Message
	msg.MediaTime = sec
	msg.Ts = nowTime
	wsc.Writer(msg)
}

func (p VlcPlayer) SyncTime(statu chan wsclient.Message) {
	Statu := <-statu
	nowTime := time.Now().UnixNano() / 1e6
	settime := int(nowTime-Statu.Ts) + (Statu.MediaTime)
	p.VP.SetMediaTime(settime)
}

func (p VlcPlayer) Pause() {
	if p.VP.IsPlaying() {
		err := p.VP.SetPause(true)
		util.AssertErr(err)
	} else {
		err := p.VP.SetPause(false)
		util.AssertErr(err)
	}
	return
}

func (p VlcPlayer) Forward() {
	sec, err := p.VP.MediaTime()
	util.AssertErr(err)
	p.VP.SetMediaTime(sec + 10000)
}

func (p VlcPlayer) Rewind() {
	sec, err := p.VP.MediaTime()
	util.AssertErr(err)
	p.VP.SetMediaTime(sec - 10000)
	p.VP.Play()
}

func (p VlcPlayer) DownVolume() (text string) {
	v, err := p.VP.Volume()
	util.AssertErr(err)
	if v <= 100 && v > 0 {
		p.VP.SetVolume(v - 10)
		t := strconv.Itoa(v - 10)
		text = "Volume Now: " + t
		return text
	}
	return "Volume Now: 0"
}

func (p VlcPlayer) UpVolume() (text string) {
	v, err := p.VP.Volume()
	util.AssertErr(err)
	if v < 100 && v >= 0 {
		p.VP.SetVolume(v + 10)
		t := strconv.Itoa(v + 10)
		text = "Volume Now: " + t
		return text
	}
	return "Volume Now: 100"
}

func (p VlcPlayer) MuteVolume(volume int) (text string, volumeNew int) {
	v, err := p.VP.Volume()
	util.AssertErr(err)
	if v != 0 {
		p.VP.SetVolume(0)
		volumeNew = v
		text := "Volume Now: 0"
		return text, volumeNew
	} else if v == 0 {
		p.VP.SetVolume(volume)
		t := strconv.Itoa(volume)
		text := "Volume Now: " + t
		volumeNew = 0
		return text, volumeNew
	}
	return
}

func (p VlcPlayer) Load(path string) (err error) {
	p.playerReleaseMedia()
	if util.IsUrl(path) {
		if _, err := p.VP.LoadMediaFromURL(path); err != nil {
			log.Printf("Cannot load selected media: %s\n", err)
			return err
		}
	} else {
		if _, err := p.VP.LoadMediaFromPath(path); err != nil {
			log.Printf("Cannot load selected media: %s\n", err)
			return err
		}
	}
	return nil
}

func (p VlcPlayer) playerReleaseMedia() {
	p.VP.Stop()
	if media, _ := p.VP.Media(); media != nil {
		media.Release()
	}
}

func (p VlcPlayer) FullScreen() {
	v, _ := p.VP.IsFullScreen()
	if v {
		p.VP.SetFullScreen(false)
	} else {
		p.VP.SetFullScreen(true)
	}
}

func (p VlcPlayer) Release() {
	err := p.VP.Release()
	util.AssertErr(err)
}

func (p VlcPlayer) TimeNow() (t string, ps float64) {
	f, _ := p.VP.MediaTime()
	l, _ := p.VP.MediaLength()
	min, sec := util.ResolveTime(f / 1000)
	te := strconv.Itoa(min) + ":" + strconv.Itoa(sec)
	min, sec = util.ResolveTime(l / 1000)
	tx := strconv.Itoa(min) + ":" + strconv.Itoa(sec)
	t = "Time: " + te + " // " + tx
	ps32, _ := p.VP.MediaPosition()
	return t, float64(ps32)
}
