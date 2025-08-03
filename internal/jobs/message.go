package jobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"messaging-server/internal/configs"
	"messaging-server/internal/database"
	log "messaging-server/internal/logging"
	"messaging-server/internal/models"
	"net/http"
	"time"
)

// sendViaAPI serializes the payload and posts it to your external URL
func sendViaAPI(client *http.Client, msg models.Message) ([]byte, error) {
	// build JSON body
	body, err := json.Marshal(
		models.SendMessage{
			Content: msg.Content,
			To:      msg.PhoneNumber,
		})
	if err != nil {
		return nil, fmt.Errorf("json.Marshal failed: %w", err)
	}

	// create request body and header
	req, err := http.NewRequest(http.MethodPost, configs.AppConfig.WebhookURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when sending the request: %w", err)
	}
	defer resp.Body.Close()

	// read the response body
	respBody, _ := io.ReadAll(resp.Body)

	// check for not accepted status codes
	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// processMessage logs a single message and marks it as sent
func processMessage(client *http.Client, msg models.Message) {

	log.Logger.Debugf("processing message id=%s to=%s", msg.ID, msg.PhoneNumber)

	// calculate the sending time
	sendingTime := time.Now().Format(time.RFC3339)

	respBody, err := sendViaAPI(client, msg)
	if err != nil {
		log.Logger.Errorf("failed to send message id=%s: %v", msg.ID, err)
		return
	}
	log.Logger.Debugf("message sent successfully at %s", sendingTime)

	// create a RedisRecord
	var redisRecord models.RedisRecord
	if err = json.Unmarshal(respBody, &redisRecord); err != nil {
		log.Logger.Errorf("failed to parse response JSON: %v; body=%s", err, string(respBody))
		return
	}
	redisRecord.SentAt = sendingTime

	//time.Sleep(3 * time.Second) // simulate processing delay

	if err = database.RedisClient.InsertRecord(redisRecord); err != nil {
		log.Logger.Errorf("failed to insert record into Redis: %v", err)
	}

	if err = database.PostgresConnection.MarkSent(msg.ID); err != nil {
		log.Logger.Errorf("failed to mark message %s as sent: %v", msg.ID, err)
	}
}

// SendMessageJob pulls up to FetchLimit messages, logs them, and marks them sent
func SendMessageJob() {

	// create a per-job HTTP client with its own Transport
	transport := &http.Transport{}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	// close idle connections when the job is done
	defer transport.CloseIdleConnections()

	messages, err := database.PostgresConnection.FetchPendingMessages(configs.AppConfig.MessageFetchLimit)
	if err != nil {
		log.Logger.Errorf("failed to fetch pending messages: %v", err)
		return
	}
	if len(messages) == 0 {
		log.Logger.Info("no pending messages to process")
		return
	}

	for _, msg := range messages {
		processMessage(client, msg)
		log.Logger.Debug("Message processed successfully")
	}
}
