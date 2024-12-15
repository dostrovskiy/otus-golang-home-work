package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeoutStr string
	flag.StringVar(&timeoutStr, "timeout", "5s", "connection timeout in seconds")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "%s\n", "Two arguments required: <host> <port>")
		return
	}
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	in := bytes.Buffer{}

	client := NewTelnetClient(address, timeout, io.NopCloser(&in), os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect:", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			in.Reset()
			if _, err := in.WriteString(scanner.Text() + "\n"); err != nil {
				fmt.Fprintln(os.Stderr, "Error writing to input buffer:", err)
				cancel()
				return
			}
			if err := client.Send(); err != nil {
				fmt.Fprintln(os.Stderr, "Send error:", err)
				cancel()
				return
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Scanner error:", err)
		}
		fmt.Fprintln(os.Stderr, "Detected EOF (Ctrl+D), closing connection...")
		cancel()
	}()

	go func() {
		if err := client.Receive(); err != nil {
			fmt.Fprintln(os.Stderr, "Receive error:", err)
			cancel()
		}
	}()

	select {
	case <-ctx.Done():
	case <-sigChan:
		fmt.Fprintln(os.Stderr, "\nReceived SIGINT, shutting down...")
		cancel()
	}
}
