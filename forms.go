package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"leoj.de/virbin/templates"
)

func httpError(w http.ResponseWriter, r *http.Request, err error, message string, public string) {
	slog.Error(message, err)
	templates.Error(public).Render(r.Context(), w)
}

func NewPasteForm(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 10<<20 {
		httpError(w, r, errors.New("to large"), "too large", "file is too large")
		return
	}
	// read form data from request
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		httpError(w, r, err, "parsing form data", "bad request")
		return
	}

	upload, header, err := r.FormFile("file")
	if err != nil {
		httpError(w, r, err, "getting file", "bad request")
		return
	}
	defer upload.Close()
	_ = header

	store, err := os.CreateTemp(".", "tmp-*")
	if err != nil {
		httpError(w, r, err, "creating store file", "internal error")
		return
	}

	reader := io.TeeReader(upload, store)

	h := sha256.New()
	_, err = io.Copy(h, reader)
	if err != nil {
		httpError(w, r, err, "hashing file", "internal error")
		os.Remove(store.Name())
		return
	}

	hash := hex.EncodeToString(h.Sum(nil))

	ok, err := checkFile(hash)
	if err != nil {
		httpError(w, r, err, "storing file ", "internal error")
		return
	}

	_ = os.MkdirAll(getFolderPath(hash, !ok), 0777)
	_ = os.Rename(store.Name(), getPath(hash, !ok))

	if ok {
		templates.Error("File is safe to upload").Render(r.Context(), w)
	} else {
		templates.UploadSuccessFul(hash).Render(r.Context(), w)
	}
}

func getFolderPath(hash string, pub bool) string {
	if pub {

		return fmt.Sprintf("./store/%s", hash[:2])
	}
	return fmt.Sprintf("./private/%s", hash[:2])
}

func getPath(hash string, pub bool) string {
	return getFolderPath(hash, pub) + "/" + hash[2:]
}
