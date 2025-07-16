package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"sync"
	"time"

	"github.com/beevik/ntp"
)

// LogEntry represents a timestamp entry with HMAC chain.
type LogEntry struct {
	Server     string        `json:"server"`
	RemoteTime time.Time     `json:"remote_time"`
	LocalTime  time.Time     `json:"local_time"`
	Offset     time.Duration `json:"offset"`
	KernelDiff bool          `json:"kernel_changed"`
	MAC        string        `json:"mac"`
}

// Logger manages chained log entries.
type Logger struct {
	mu       sync.Mutex
	file     *os.File
	macKey   []byte
	lastMAC  []byte
	lastKxor uint64
}

func NewLogger(path string, key []byte) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{file: f, macKey: key}, nil
}

func (l *Logger) Close() error {
	return l.file.Close()
}

func (l *Logger) kernelXor() uint64 {
	var uts unix.Utsname
	if err := unix.Uname(&uts); err != nil {
		return 0
	}
	fields := [][]byte{uts.Sysname[:], uts.Nodename[:], uts.Release[:], uts.Version[:], uts.Machine[:]}
	xorVal := uint64(0)
	for _, f := range fields {
		for _, b := range f {
			xorVal ^= uint64(b)
		}
	}
	return xorVal
}

func (l *Logger) append(entry LogEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	if _, err := l.file.Write(data); err != nil {
		return err
	}
	if _, err := l.file.Write([]byte("\n")); err != nil {
		return err
	}
	return nil
}

func (l *Logger) log(server string, resp *ntp.Response) error {
	kxor := l.kernelXor()
	diff := kxor != l.lastKxor
	l.lastKxor = kxor

	mac := hmac.New(sha256.New, l.macKey)
	mac.Write(l.lastMAC)
	fmt.Fprintf(mac, "%s|%d|%d", server, resp.Time.UnixNano(), resp.ClockOffset)
	macSum := mac.Sum(nil)

	entry := LogEntry{
		Server:     server,
		RemoteTime: resp.Time.UTC(),
		LocalTime:  time.Now().UTC(),
		Offset:     resp.ClockOffset,
		KernelDiff: diff,
		MAC:        hex.EncodeToString(macSum),
	}

	if err := l.append(entry); err != nil {
		return err
	}

	l.lastMAC = macSum
	return nil
}

func runOnce(l *Logger, server string, timeout time.Duration) error {
	resp, err := ntp.QueryWithOptions(server, ntp.QueryOptions{Timeout: timeout})
	if err != nil {
		return err
	}
	return l.log(server, resp)
}

func main() {
	logPath := flag.String("log", "ntp.log", "path to log file")
	interval := flag.Duration("interval", 2*time.Second, "query interval")
	keyFile := flag.String("key", "hmac.key", "HMAC key file")
	flag.Parse()

	key, err := os.ReadFile(*keyFile)
	if err != nil {
		log.Fatalf("read key: %v", err)
	}

	logger, err := NewLogger(*logPath, key)
	if err != nil {
		log.Fatalf("open log: %v", err)
	}
	defer logger.Close()

	servers := []string{
		"time.google.com",
		"time.apple.com",
		"time.windows.com",
		"0.ru.pool.ntp.org",
		"1.ru.pool.ntp.org",
		"2.ru.pool.ntp.org",
		"3.ru.pool.ntp.org",
	}

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	for {
		for _, s := range servers {
			if err := runOnce(logger, s, *interval); err != nil {
				fmt.Fprintln(os.Stderr, "query error", s, err)
			}
		}
		<-ticker.C
	}
}
