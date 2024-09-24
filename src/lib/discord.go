package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

/*

	Enums

*/

// type DiscordStyle int

// const (
// 	DiscordStylePrimary   DiscordStyle = 1 // blue
// 	DiscordStyleSecondary DiscordStyle = 2 // gray
// 	DiscordStyleSuccess   DiscordStyle = 3 // green
// 	DiscordStyleDanger    DiscordStyle = 4 // red
// )

// type DiscordButtonType int

// const (
// 	DiscordButtonSimple DiscordButtonType = 2 // gray
// )

// type DiscordButton struct {
// 	Type     DiscordButtonType `json:"type"`
// 	Label    string            `json:"label"`
// 	Style    DiscordStyle      `json:"style"`
// 	CustomId string            `json:"custom_id,omitempty"`
// 	URL      string            `json:"url,omitempty"`
// }

// type DiscordComponent struct {
// 	Type       int             `json:"type"`
// 	Components []DiscordButton `json:"components"`
// }

/*
Functions
*/
func DiscordSendTextMessage(webhookURL, message string) error {

	payload, err := json.Marshal(map[string]string{
		"content": message,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Failed to send log to Discord, status: %s", resp.Status)
	}

	return nil
}

// func DiscordSendWithComponents(webhookURL, message string, component DiscordComponent) error {

// 	payload, err := json.Marshal(map[string]any{
// 		"content":    message,
// 		"components": component,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
// 		return fmt.Errorf("Failed to send log to Discord, status: %s", resp.Status)
// 	}

// 	return nil
// }
