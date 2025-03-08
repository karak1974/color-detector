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
	On  bool `json:"on"`
	Seg []struct {
		Col [][]int `json:"col"`
	} `json:"seg"`
}

func sendToWLED(ip string, r, g, b int) error {
	url := fmt.Sprintf("http://%s/json/state", ip)

	state := WLEDState{
		On: true,
		Seg: []struct {
			Col [][]int `json:"col"`
		}{
			{Col: [][]int{{r, g, b}}},
		},
	}

	jsonData, err := json.Marshal(state)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to set color, status code: %d", resp.StatusCode)
	}

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

			if err := sendToWLED(wledIP, int(clr.R), int(clr.B), int(clr.G)); err != nil {
        		fmt.Println("Error sending request:", err)
    		}
		}
		time.Sleep(interval)
	}
}

