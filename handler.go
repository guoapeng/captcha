package captcha

import (
	"io"
	"net/http"
)

type MediaHandler interface {
	handle(w http.ResponseWriter, d []byte, r *http.Request, id string)
	createMedia(d []byte) io.WriterTo
}

type AudioHandler struct {
	lang     string
}

func (h *AudioHandler) handle(w http.ResponseWriter, d []byte, r *http.Request, id string) {
	w.Header().Set("Content-Type", "audio/x-wav")
	writeMedia(h.createMedia(d), false, w, r, id+".wav")
}

func (h *AudioHandler) createMedia(d []byte) io.WriterTo {
	return NewAudio(d, h.lang)
}

type ImgHandler struct {
	imgWidth  int
	imgHeight int
}

func (h *ImgHandler) handle(w http.ResponseWriter, d []byte, r *http.Request, id string) {
	w.Header().Set("Content-Type", "image/png")
	writeMedia(h.createMedia(d), false, w, r, id+".png")
}

func (h *ImgHandler) createMedia(d []byte) io.WriterTo {
	return NewImage(d, h.imgWidth, h.imgHeight)
}

