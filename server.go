// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"path"
	"strings"
	"time"
)

type captchaHandler struct {
	imgWidth    int
	imgHeight   int
	forceReload bool
}

// Server returns a handler that serves HTTP requests with image or
// audio representations of captchas. Image dimensions are accepted as
// arguments. The server decides which captcha to serve based on the last URL
// path component: file name part must contain a captcha id, file extension —
// its format (PNG or WAV).
//
// For example, for file name "LBm5vMjHDtdUfaWYXiQX.png" it serves an image captcha
// with id "LBm5vMjHDtdUfaWYXiQX", and for "LBm5vMjHDtdUfaWYXiQX.wav" it serves the
// same captcha in audio format.
//
// To serve a captcha as a downloadable file, the URL must be constructed in
// such a way as if the file to serve is in the "download" subdirectory:
// "/download/LBm5vMjHDtdUfaWYXiQX.wav".
//
// To reload captcha (get a different solution for the same captcha id), append
// "?reload=x" to URL, where x may be anything (for example, current time or a
// random number to make browsers refetch an image instead of loading it from
// cache).
//
// By default, the Server serves audio in English language. To serve audio
// captcha in one of the other supported languages, append "lang" value, for
// example, "?lang=ru".
func Server(imgWidth, imgHeight int, forceReload bool) http.Handler {
	return &captchaHandler{imgWidth, imgHeight, forceReload}
}

func (h *captchaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir, file := path.Split(r.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext != "" && id != "" {
		if r.FormValue("reload") != "" || h.forceReload {
			Reload(id, h.forceReload)
		}
		if d := globalStore.Get(id, false); d != nil {
			download := path.Base(dir) == "download"
			if handler, err := h.createMediaHandler(r); err == nil {
				writeMedia(handler.createMedia(d), download, w, r, file)
			} else {
				http.NotFound(w, r)
			}
		} else {
			http.NotFound(w, r)
		}
	} else {
		http.NotFound(w, r)
	}

	// Ignore other errors.
}

func writeMedia(writer io.WriterTo, download bool, w http.ResponseWriter, r *http.Request, name string) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	var content bytes.Buffer
	writer.WriteTo(&content)
	http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(content.Bytes()))
}

func (h *captchaHandler) createMediaHandler(r *http.Request) (MediaHandler, error) {
	_, file := path.Split(r.URL.Path)
	ext := path.Ext(file)
	lang := strings.ToLower(r.FormValue("lang"))

	switch ext {
	case ".png":
		return &ImgHandler{h.imgWidth, h.imgHeight}, nil
	case ".wav":
		return &AudioHandler{lang}, nil
	default:
		return nil, errors.New("not supported type " + ext)
	}
}
