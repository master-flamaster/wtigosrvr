Скомпилируйте приложение для Linux:

GOOS=linux GOARCH=amd64 go build -o ntp-monitor

Используйте cloud-init скрипт для развертывания:

#cloud-config
packages:
  - golang
runcmd:
  - mkdir /opt/ntp-monitor
  - cd /opt/ntp-monitor
  - wget https://example.com/ntp-monitor
  - chmod +x ntp-monitor
  - ./ntp-monitor &


  # /etc/systemd/system/ntp-monitor.service
[Unit]
Description=NTP Monitor Service

[Service]
ExecStart=/opt/ntp-monitor/ntp-monitor
Restart=always

[Install]
WantedBy=multi-user.target


package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("NTP Monitor")
	systray.SetTooltip("NTP Server Monitoring")

	mQuit := systray.AddMenuItem("Quit", "Exit application")

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go startMonitor(ctx)

		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func startMonitor(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	prevXor := false

	for {
		select {
		case <-ticker.C:
			currentXor := xorKernelFlag()
			if currentXor != prevXor {
				showNotification("XOR Flag Changed", 
					fmt.Sprintf("State changed: %v → %v", prevXor, currentXor))
				prevXor = currentXor
			}
		case <-ctx.Done():
			return
		}
	}
}

func showNotification(title, message string) {
	exec.Command("osascript", "-e", 
		fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)).Run()
}

func onExit() {}


