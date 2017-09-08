package discohook

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Message struct {
	Content   string  `json:"content"`
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	TTS       bool    `json:"tts"`
	File      []byte  `json:"file"`
	Embeds    []Embed `json:"embeds"`
}

type Embed struct {
	Title       string          `json:"title,omitempty"`
	Type        string          `json:"type,omitempty"`
	Description string          `json:"description,omitempty"`
	URL         string          `json:"url,omitempty"`
	Timestamp   string          `json:"timestamp,omitempty"`
	Color       *Color          `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *EmbedVideo     `json:"video,omitempty"`
	Provider    *EmbedProvider  `json:"provider,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedField    `json:"fields,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type EmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type EmbedVideo struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type EmbedProvider struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

type EmbedAuthor struct {
	URL          string `json:"url,omitempty"`
	Name         string `json:"name,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type Color struct {
	R, G, B byte
}

func (c *Color) MarshalJSON() ([]byte, error) {
	number := strconv.Itoa(int(c.R)<<16 | int(c.G)<<8 | int(c.B))
	return []byte(number), nil
}

func (c *Color) UnmarshalJSON(data []byte) error {
	number, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}

	c.R = byte((number >> 16) & 0xFF)
	c.G = byte((number >> 8) & 0xFF)
	c.B = byte(number & 0xFF)

	return nil
}

func Send(url string, msg *Message, wait bool) error {
	body := new(bytes.Buffer)

	err := json.NewEncoder(body).Encode(msg)
	if err != nil {
		return err
	}

	if wait {
		url += "?wait=true"
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "discohook/1.0")
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("discohook: server returned error")
	}

	return nil
}
