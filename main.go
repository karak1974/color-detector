package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"net/http"
	"time"

	"github.com/kbinani/screenshot"
)

type WLEDState struct {
    On  bool  `json:"on"`
    Bri uint8 `json:"bri"`
    R   uint8 `json:"r"`
    G   uint8 `json:"g"`
    B   uint8 `json:"b"`
}

func sendToWLED(ip string, state WLEDState) error {
    url := fmt.Sprintf("http://%s/json/state", ip)
    jsonData, err := json.Marshal(state)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    fmt.Println("Response Status:", resp.Status)
    return nil
}

func getColorAt(x, y int) (color.RGBA, error) {
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return color.RGBA{}, err
	}

	clr := img.At(x, y).(color.RGBA)
	return clr, nil
}

func printColor(r, g, b uint8, char string) {
    fmt.Printf("\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, char)
}

func main() {
	interval   := 100 * time.Millisecond
	posX, posY := 100, 100
	
	wledIP := "x.x.x.x"

	for {
		clr, err := getColorAt(posX, posY)
		if err != nil {
			fmt.Println("Error capturing screen:", err)
		} else {
			printColor(clr.R, clr.G, clr.B, fmt.Sprintf("█ R=%d G=%d B=%d\n", clr.R, clr.G, clr.B ))

			state := WLEDState{
        		On:  true,
        		Bri: 255,
        		R:   clr.R,
        		G:   clr.G,
        		B:   clr.B,
    		}

    		if err := sendToWLED(wledIP, state); err != nil {
        		fmt.Println("Error sending request:", err)
    		}
		}
		time.Sleep(interval)
	}
}

