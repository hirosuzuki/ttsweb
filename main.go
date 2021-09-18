package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var content []byte

func synthesizeSsml(ssml string) error {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Ssml{Ssml: ssml},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "ja-JP",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return err
	}

	content = resp.AudioContent

	return nil
}

func cutString(s string, n int) string {
	runes := []rune(s)
	if n > len(runes) {
		n = len(runes)
	}
	return string(runes[:n])
}

func synthesizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(1024 * 1024)
		ssml := r.Form.Get("ssml")
		if ssml != "" {
			log.Println("Synthesize:", cutString(strings.Replace(ssml, "\r\n", " ", -1), 40))
			err := synthesizeSsml(ssml)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprint(w, err.Error())
			} else {
				w.WriteHeader(200)
			}
		}
	} else {
		w.Header().Add("Content-Type", "audio/mpeg")
		w.WriteHeader(200)
		w.Write(content)
	}
}

//go:embed static/*
var staticFS embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000"
	}

	if os.Getenv("DEV") == "" {
		staticDir, _ := fs.Sub(staticFS, "static")
		http.Handle("/", http.FileServer(http.FS(staticDir)))
	} else {
		staticDir := http.FileServer(http.Dir("static"))
		http.Handle("/", staticDir)
	}

	http.HandleFunc("/api/synthesize", synthesizeHandler)
	http.HandleFunc("/synthesize.mp3", synthesizeHandler)
	log.Printf("Listen: http://127.0.0.1%s", port)
	http.ListenAndServe(port, nil)
}
