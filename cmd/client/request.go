package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log/slog"
	"net"

	"wordofwisdom/internal/service/pow"
)

// sendRequest sends a request to the server
func sendRequest(serverAddr string, handlerID byte) error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	buffer := make([]byte, 9) // 8 bytes for the challenge and 1 byte for the difficulty
	if _, err = conn.Read(buffer); err != nil {
		return fmt.Errorf("failed to read challenge and difficulty: %w", err)
	}

	challenge := buffer[:8] // first 8 bytes are the challenge
	difficulty := buffer[8] // last byte is the difficulty

	powAlgo := pow.NewProofOfWork(difficulty)
	solution := powAlgo.Solve(challenge)

	var requestBuffer bytes.Buffer
	requestBuffer.Write(solution)

	// write the handlerID (1 byte) to the buffer
	if err = requestBuffer.WriteByte(handlerID); err != nil {
		return fmt.Errorf("failed to write handler ID: %w", err)
	}

	// send the binary request to the server
	if _, err = conn.Write(requestBuffer.Bytes()); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	// read the response from the server
	var response []byte
	if response, err = bufio.NewReader(conn).ReadBytes('\n'); err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	slog.Info("server response", "response", string(response))
	return nil
}
