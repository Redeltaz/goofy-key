package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func main() {
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	sampleRate := beep.SampleRate(44100)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))

	fmt.Println("Press ESC to quit")
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}
		fmt.Println("Key pressed")
		if event.Key == keyboard.KeyEsc {
			break
		}

		go playRandomSound()
	}
}

func playRandomSound() {
	fileNumber := rand.Intn(23) + 1
	fileName := fmt.Sprintf("./sounds/sound%d.wav", fileNumber)

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("failed to open sound file: %v", err)
	}
	defer f.Close()

	streamer, _, err := wav.Decode(f)
	if err != nil {
		log.Fatalf("failed to decode wav file: %v", err)
	}
	defer streamer.Close()

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}
