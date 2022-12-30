package utils

import (
	"encoding/json"
)

type vlcJson struct {
	Fullscreen int `json:"fullscreen"`
	Stats      struct {
		Inputbitrate        float64 `json:"inputbitrate"`
		Sentbytes           int     `json:"sentbytes"`
		Lostabuffers        int     `json:"lostabuffers"`
		Averagedemuxbitrate int     `json:"averagedemuxbitrate"`
		Readpackets         int     `json:"readpackets"`
		Demuxreadpackets    int     `json:"demuxreadpackets"`
		Lostpictures        int     `json:"lostpictures"`
		Displayedpictures   int     `json:"displayedpictures"`
		Sentpackets         int     `json:"sentpackets"`
		Demuxreadbytes      int     `json:"demuxreadbytes"`
		Demuxbitrate        float64 `json:"demuxbitrate"`
		Playedabuffers      int     `json:"playedabuffers"`
		Demuxdiscontinuity  int     `json:"demuxdiscontinuity"`
		Decodedaudio        int     `json:"decodedaudio"`
		Sendbitrate         int     `json:"sendbitrate"`
		Readbytes           int     `json:"readbytes"`
		Averageinputbitrate int     `json:"averageinputbitrate"`
		Demuxcorrupted      int     `json:"demuxcorrupted"`
		Decodedvideo        int     `json:"decodedvideo"`
	} `json:"stats"`
	Audiodelay   int  `json:"audiodelay"`
	Apiversion   int  `json:"apiversion"`
	Currentplid  int  `json:"currentplid"`
	Time         int  `json:"time"`
	Volume       int  `json:"volume"`
	Length       int  `json:"length"`
	Random       bool `json:"random"`
	Audiofilters struct {
		Filter0 string `json:"filter_0"`
	} `json:"audiofilters"`
	Rate         int `json:"rate"`
	Videoeffects struct {
		Hue        int `json:"hue"`
		Saturation int `json:"saturation"`
		Contrast   int `json:"contrast"`
		Brightness int `json:"brightness"`
		Gamma      int `json:"gamma"`
	} `json:"videoeffects"`
	State       string  `json:"state"`
	Loop        bool    `json:"loop"`
	Version     string  `json:"version"`
	Position    float64 `json:"position"`
	Information struct {
		Chapter  int           `json:"chapter"`
		Chapters []interface{} `json:"chapters"`
		Title    int           `json:"title"`
		Category struct {
			Meta struct {
				DISCID      string `json:"DISCID"`
				Date        string `json:"date"`
				ArtworkUrl  string `json:"artwork_url"`
				Artist      string `json:"artist"`
				Album       string `json:"album"`
				TrackNumber string `json:"track_number"`
				Filename    string `json:"filename"`
				Title       string `json:"title"`
				Genre       string `json:"genre"`
			} `json:"meta"`
			Stream0 struct {
				Bitrate       string `json:"Bitrate"`
				Codec         string `json:"Codec"`
				Channels      string `json:"Channels"`
				BitsPerSample string `json:"Bits_per_sample"`
				Type          string `json:"Type"`
				SampleRate    string `json:"Sample_rate"`
			} `json:"Stream 0"`
		} `json:"category"`
		Titles []interface{} `json:"titles"`
	} `json:"information"`
	Repeat        bool          `json:"repeat"`
	Subtitledelay int           `json:"subtitledelay"`
	Equalizer     []interface{} `json:"equalizer"`
}

func addOrRemoveClass(add bool) string {
	if add {
		return "addClass"
	} else {
		return "removeClass"
	}
}

func processShuffle(result vlcJson) string {
	return addOrRemoveClass(result.Random) + ",shuffle-btn,checked;"
}

func processLoop(result vlcJson) string {
	return addOrRemoveClass(result.Loop) + ",loop-btn,checked;"
}

func processPausedPlaying(result vlcJson) string {
	isPlaying := result.State == "playing"
	oper := addOrRemoveClass(isPlaying) + ",play-pause-btn,pause-btn;"
	return oper + addOrRemoveClass(!isPlaying) + ",play-pause-btn,play-btn;"
}

func processCurrentSong(result vlcJson) string {
	artist := result.Information.Category.Meta.Artist
	title := result.Information.Category.Meta.Title
	album := result.Information.Category.Meta.Album
	filename := result.Information.Category.Meta.Filename
	value := ""
	if artist != "" {
		value += "songInfo," + artist
		if title != "" {
			value += ": " + title
		}
		if album != "" {
			value += " (" + album + ")"
		}
	} else if filename != "" {
		value += "songInfo," + filename
	}
	return value
}

func ProcessOperations(vlcResponse string) string {
	var result vlcJson
	json.Unmarshal([]byte(vlcResponse), &result)
	operations := ""

	operations += processShuffle(result)
	operations += processLoop(result)
	operations += processPausedPlaying(result)
	operations += processCurrentSong(result)

	return operations
}