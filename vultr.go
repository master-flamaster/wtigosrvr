package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
)

type Config struct {
	APIKey       string `json:"api_key"`
	InstanceID   string `json:"instance_id"`
	SSHPublicKey string `json:"ssh_public_key"`
}

func readJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func main() {
	cfgPath := flag.String("config", "vultr_config.json", "config file")
	payloadPath := flag.String("payload", "vultr_api_payload.json", "payload file")
	flag.Parse()

	var cfg Config
	if err := readJSON(*cfgPath, &cfg); err != nil {
		log.Fatalf("read config: %v", err)
	}
	var payload map[string]interface{}
	if err := readJSON(*payloadPath, &payload); err != nil {
		log.Fatalf("read payload: %v", err)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.vultr.com/v2/instances", bytes.NewReader(body))
	if err != nil {
		log.Fatalf("create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("send request: %v", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("status: %s", resp.Status)
	log.Printf("response: %s", respBody)
}
