package connect

import (
	"errors"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

func SendConnectionRequest(page *rod.Page, profileURL string) error {
	page.MustNavigate(profileURL)
	time.Sleep(5 * time.Second)

	page.Mouse.Scroll(0, 400, 1)
	time.Sleep(2 * time.Second)

	btn, err := page.ElementR("button", "Connect")
	if err != nil {
		return errors.New("connect button not found")
	}

	btn.MustHover()
	time.Sleep(1 * time.Second)
	btn.MustClick()

	time.Sleep(2 * time.Second)

	addNoteBtn, err := page.ElementR("button", "Add a note")
	if err == nil {
		addNoteBtn.MustClick()
		time.Sleep(1 * time.Second)

		noteBox, err := page.Element("textarea")
		if err == nil {
			note := "Hi, I came across your profile and would be happy to connect."
			if len(note) > 300 {
				note = note[:300]
			}
			noteBox.MustInput(note)
			time.Sleep(1 * time.Second)
		}
	}

	sendBtn, err := page.ElementR("button", "Send")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.New("send button not found")
		}
		return err
	}

	sendBtn.MustClick()
	time.Sleep(3 * time.Second)

	return nil
}
