package ebitensim

import (
	"bytes"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

func parseWavFile(path string) []byte {
	rawData, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	s, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(rawData))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	b, err := io.ReadAll(s)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return b
}

func parseOggFile(path string) []byte {
	rawData, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	s, err := vorbis.DecodeWithSampleRate(sampleRate, bytes.NewReader(rawData))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	b, err := io.ReadAll(s)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return b
}

func parseMp3File(path string) []byte {
	rawData, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	s, err := mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(rawData))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	b, err := io.ReadAll(s)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return b
}

func (g *EbitenGame) addSound(pathStr, soundName string) {
	fileName := strings.ToLower(path.Base(pathStr))

	var soundData []byte
	switch {
	case strings.Contains(fileName, ".wav"):
		soundData = parseWavFile(pathStr)
	case strings.Contains(fileName, ".ogg"):
		soundData = parseOggFile(pathStr)
	case strings.Contains(fileName, ".mp3"):
		soundData = parseMp3File(pathStr)
	default:
		log.Printf("Unable to decode sound file: %s\n", pathStr)
		return
	}

	g.sounds[soundName] = soundData
}

func (g *EbitenGame) playSound(soundName string, volume float64) {
	if volume < 0 {
		volume = 0
	}
	if volume > 1 {
		volume = 1
	}

	soundData, ok := g.sounds[soundName]
	if !ok {
		log.Printf("Sound %s not found.\n", soundName)
		return
	}

	p := g.audioContext.NewPlayerFromBytes(soundData)
	p.Play()
	p.SetVolume(volume)
}
