package main

import (
	"crypto"
	"encoding/base64"
	"math/rand"
	"strconv"
	"strings"
)

type ShortURL struct {
	URL string `json:"short_url"`
}

type OriginURL struct {
	URL string `json:"url"`
}

func (app *application) createShortURL(URL string) string {
	url := URL + strconv.Itoa(rand.Int())
	h := crypto.SHA3_256.New()
	h.Write([]byte(url))
	hash := h.Sum(nil)
	encoded := base64.RawURLEncoding.EncodeToString(hash)[:6]
	return encoded
}

func getShortURL(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 3 || parts[2] == "" {
		return ""
	}
	return parts[2]
}
