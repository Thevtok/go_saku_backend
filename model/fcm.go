package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"

	"github.com/ReygaFitra/inc-final-project.git/utils"
	"github.com/sirupsen/logrus"
)

type FCMNotification struct {
	To           string              `json:"to"`
	Notification FCMNotificationData `json:"notification"`
}

type FCMNotificationData struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func SendFCMNotification(token string, title string, body string) error {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	notification := FCMNotification{
		To: token,
		Notification: FCMNotificationData{
			Title: title,
			Body:  body,
		},
	}

	payload, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %v", err)
	}

	url := "https://fcm.googleapis.com/fcm/send"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create FCM request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", utils.DotEnv("FCM_KEY")))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send FCM request: %v", err)
	}
	defer res.Body.Close()

	// Check the response status
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("FCM request failed with status: %s", res.Status)
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read FCM response body: %v", err)
	}
	logrus.Infof("FCM response: %s", string(responseBody))

	return nil
}
