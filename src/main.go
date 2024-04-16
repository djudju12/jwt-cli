package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
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
	fileInfo, _ := os.Stdin.Stat()
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

	ba, err := json.MarshalIndent(token, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(ba))
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
