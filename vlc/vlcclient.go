package vlc

import (
	vlc "github.com/adrg/libvlc-go/v2"
	"log"
)

type VLCControl struct {
}

func (v VLCControl) Init() {
	// Initialize libVLC. Additional command line arguments can be passed in
	// to libVLC by specifying them in the Init function.
	// if err := vlc.Init("--no-video", "--quiet"); err != nil {
	if err := vlc.Init(); err != nil {
		log.Fatal(err)
	}
	defer vlc.Release()
}

func (v VLCControl) Play() {
	// Create a new player.
	player, err := vlc.NewPlayer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		player.Stop()
		player.Release()
	}()

	media, err := player.LoadMediaFromPath("d:/Hudba/JINE/Guild Of Lore - Winterstead.mp3")
	if err != nil {
		log.Fatal(err)
	}
	defer media.Release()

	// Retrieve player event manager.
	manager, err := player.EventManager()
	if err != nil {
		log.Fatal(err)
	}

	// Register the media end reached event with the event manager.
	quit := make(chan struct{})
	eventCallback := func(event vlc.Event, userData interface{}) {
		close(quit)
	}

	eventID, err := manager.Attach(vlc.MediaPlayerEndReached, eventCallback, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer manager.Detach(eventID)

	// Start playing the media.
	if err = player.Play(); err != nil {
		log.Fatal(err)
	}

	<-quit
}
