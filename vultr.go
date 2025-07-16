package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
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

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfgPath := flag.String("config", "vultr_config.json", "config file path")
	payloadPath := flag.String("payload", "vultr_api_payload.json", "API payload JSON")
	dryRun := flag.Bool("dry", false, "print payload without sending")
	flag.Parse()

	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	payload, err := os.ReadFile(*payloadPath)
	if err != nil {
		log.Fatalf("read payload: %v", err)
	}

	if *dryRun {
		fmt.Println(string(payload))
		return
	}

	req, err := http.NewRequest("POST", "https://api.vultr.com/v2/instances", bytes.NewReader(payload))
	if err != nil {
		log.Fatalf("make request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("status: %s\n%s\n", resp.Status, body)
}
