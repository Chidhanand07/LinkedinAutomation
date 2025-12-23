package messaging

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"linkedin-automation/config"
	"linkedin-automation/stealth"

	"github.com/go-rod/rod"
)

// AutoReplyToMessages checks for unread messages and replies automatically
func AutoReplyToMessages(page *rod.Page, cfg *config.MessagingConfig) (int, error) {
	if !cfg.AutoReplyEnabled {
		log.Println("⚠️  Auto-reply is disabled in config")
		return 0, nil
	}

	log.Println("💬 Starting auto-reply to messages...")

	// Navigate to messaging page
	page.MustNavigate("https://www.linkedin.com/messaging/")
	time.Sleep(stealth.RandomizedTimingDelay("navigate"))

	repliesSent := 0

	// Find unread message threads
	unreadThreads, err := page.Elements("li.msg-conversation-listitem--unread")
	if err != nil || len(unreadThreads) == 0 {
		log.Println("✓ No unread messages to reply to")
		return 0, nil
	}

	log.Printf("Found %d unread message thread(s)", len(unreadThreads))

	maxReplies := cfg.MaxAutoRepliesPerSession
	if len(unreadThreads) < maxReplies {
		maxReplies = len(unreadThreads)
	}

	for i := 0; i < maxReplies; i++ {
		// Re-query unread threads as DOM changes
		unreadThreads, err = page.Elements("li.msg-conversation-listitem--unread")
		if err != nil || len(unreadThreads) == 0 {
			break
		}

		thread := unreadThreads[0]

		// Get sender name
		senderName := getSenderName(thread)
		log.Printf("  Replying to: %s", senderName)

		// Click on thread to open conversation
		thread.MustClick()
		time.Sleep(stealth.RandomizedTimingDelay("click"))

		// Find message input
		messageBox, err := page.Element("div.msg-form__contenteditable[contenteditable='true']")
		if err != nil {
			log.Printf("  ⚠️  Could not find message input")
			continue
		}

		// Select random auto-reply template
		template := cfg.AutoReplyTemplates[rand.Intn(len(cfg.AutoReplyTemplates))]

		// Type the message
		log.Printf("  Typing reply...")
		stealth.HumanType(messageBox, template)
		time.Sleep(stealth.RandomizedTimingDelay("think"))

		// Find and click send button
		sendBtn, err := page.Element("button.msg-form__send-button[type='submit']")
		if err != nil {
			log.Printf("  ⚠️  Could not find send button")
			continue
		}

		// Move mouse to send button
		box := sendBtn.MustShape().Box()
		targetX := box.X + box.Width/2
		targetY := box.Y + box.Height/2
		stealth.HumanMove(page, targetX, targetY)
		stealth.HumanHover(page, targetX, targetY, time.Duration(rand.Intn(300)+200)*time.Millisecond)

		// Click send
		sendBtn.MustClick()
		repliesSent++
		log.Printf("  ✓ Reply sent to %s", senderName)

		// Delay between replies
		delay := time.Duration(rand.Intn(cfg.DelayBetweenMessagesMax-cfg.DelayBetweenMessagesMin)+cfg.DelayBetweenMessagesMin) * time.Second
		log.Printf("  ⏳ Waiting %v before next reply...", delay.Round(time.Second))
		time.Sleep(delay)
	}

	log.Printf("✓ Sent %d auto-replies", repliesSent)
	return repliesSent, nil
}

// getSenderName extracts sender name from message thread element
func getSenderName(thread *rod.Element) string {
	nameElement, err := thread.Element("h3.msg-conversation-listitem__participant-names")
	if err != nil {
		return "Unknown"
	}

	name, err := nameElement.Text()
	if err != nil {
		return "Unknown"
	}

	return strings.TrimSpace(name)
}

// SendFollowUps sends follow-up messages to new connections (config-driven)
func SendFollowUps(page *rod.Page, db *sql.DB, cfg *config.MessagingConfig) error {
	if !cfg.Enabled || !cfg.SendFollowupToNewConnections {
		log.Println("⚠️  Follow-up messaging is disabled in config")
		return nil
	}

	log.Println("📨 Starting follow-up messages to new connections...")

	page.MustNavigate("https://www.linkedin.com/mynetwork/invite-connect/connections/")
	time.Sleep(stealth.RandomizedTimingDelay("navigate"))

	profiles, err := page.Elements("a.app-aware-link")
	if err != nil {
		return fmt.Errorf("failed to find profiles: %w", err)
	}

	messagesSent := 0
	maxMessages := cfg.MaxMessagesPerSession

	for _, profile := range profiles {
		if messagesSent >= maxMessages {
			log.Printf("✓ Reached max messages per session: %d", maxMessages)
			break
		}

		href, err := profile.Attribute("href")
		if err != nil || href == nil {
			continue
		}

		profileURL := *href
		if !strings.Contains(profileURL, "/in/") {
			continue
		}

		// Check if already messaged
		if cfg.SkipAlreadyMessaged {
			// Check database
			var count int
			err := db.QueryRow("SELECT COUNT(*) FROM profiles WHERE url = ? AND messaged = 1", profileURL).Scan(&count)
			if err == nil && count > 0 {
				continue
			}
		}

		log.Printf("  Sending follow-up to: %s", profileURL)

		page.MustNavigate(profileURL)
		time.Sleep(stealth.RandomizedTimingDelay("navigate"))

		// Find Message button
		msgBtn, err := page.ElementR("button", "Message")
		if err != nil {
			log.Printf("  ⚠️  Message button not found")
			continue
		}

		// Click Message
		box := msgBtn.MustShape().Box()
		stealth.HumanMove(page, box.X+box.Width/2, box.Y+box.Height/2)
		msgBtn.MustClick()
		time.Sleep(stealth.RandomizedTimingDelay("click"))

		// Find message input
		messageBox, err := page.Element("div[contenteditable='true']")
		if err != nil {
			log.Printf("  ⚠️  Message input not found")
			continue
		}

		// Select random template
		template := cfg.FollowupTemplates[rand.Intn(len(cfg.FollowupTemplates))]

		// Type message
		stealth.HumanType(messageBox, template)
		time.Sleep(stealth.RandomizedTimingDelay("think"))

		// Find and click send
		sendBtn, err := page.ElementR("button", "Send")
		if err == nil {
			sendBtn.MustClick()
			messagesSent++
			log.Printf("  ✓ Follow-up sent")

			// Mark as messaged in database
			if cfg.MarkAsMessaged {
				db.Exec("UPDATE profiles SET messaged = 1, message_date = ? WHERE url = ?", time.Now(), profileURL)
			}
		}

		// Delay between messages
		delay := time.Duration(rand.Intn(cfg.DelayBetweenMessagesMax-cfg.DelayBetweenMessagesMin)+cfg.DelayBetweenMessagesMin) * time.Second
		time.Sleep(delay)
	}

	log.Printf("✓ Sent %d follow-up messages", messagesSent)
	return nil
}
