package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// token errors
const (
	ErrMalformedToken  = "malformed token"
	ErrNoTokenProvided = "no token provided"
)

type J map[string]any

type Token struct {
	Header  J `json:"header"`
	Payload J `json:"payload"`
}

func main() {
	var tokenStr string
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	timeFlag := false
	if (len(os.Args) > 1) && (strings.HasPrefix(os.Args[1], "-")) {
		if os.Args[1] == "-t" || os.Args[1] == "--timestamp" {
			timeFlag = true
			os.Args = os.Args[1:]
		} else {
			panic(fmt.Errorf("invalid flag: %s", os.Args[1]))
		}
	}

	if fileInfo.Mode()&os.ModeCharDevice == 0 {
		reader := bufio.NewReader(os.Stdin)
		buffer, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		tokenStr = string(buffer)
	} else if len(os.Args) > 1 {
		tokenStr = os.Args[1]
	} else {
		panic(fmt.Errorf(ErrNoTokenProvided))
	}

	tokenParts, err := breakToken(tokenStr)
	if err != nil {
		panic(err)
	}

	headerStr, err := decodePart(tokenParts[0])
	if err != nil {
		panic(err)
	}

	var header J
	if err = json.Unmarshal(headerStr, &header); err != nil {
		panic(err)
	}

	payloadStr, err := decodePart(tokenParts[1])
	if err != nil {
		panic(err)
	}

	var payload J
	if err = json.Unmarshal(payloadStr, &payload); err != nil {
		panic(err)
	}

	token := Token{
		Header:  header,
		Payload: payload,
	}

	if timeFlag {
		var timestamp int64 = int64(token.Payload["exp"].(float64))
		token.Payload["exp"] = fmt.Sprintf(
			"%d (%s)",
			timestamp,
			time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"),
		)
		timestamp = int64(token.Payload["iat"].(float64))
		token.Payload["iat"] = fmt.Sprintf(
			"%d (%s)",
			timestamp,
			time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"),
		)
	}

	bytea, err := json.MarshalIndent(token, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytea))
}

func breakToken(tokenStr string) ([]string, error) {
	tokenPars := strings.Split(tokenStr, ".")
	if (tokenPars == nil) || (len(tokenPars) != 3) {
		return nil, fmt.Errorf(ErrMalformedToken)
	}

	return tokenPars, nil
}

func decodePart(part string) ([]byte, error) {
	return base64.RawStdEncoding.DecodeString(part)
}
