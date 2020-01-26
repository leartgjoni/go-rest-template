package utils

import (
	"github.com/go-chi/render"
	"github.com/leartgjoni/go-rest-template/http/payloads"
	"net/http"
)

func Render(w http.ResponseWriter, r *http.Request, v render.Renderer) {
	if err := render.Render(w, r, v); err != nil {
		err := render.Render(w, r, payloads.ErrRender(err))
		panic(err)
	}
}

func RenderList(w http.ResponseWriter, r *http.Request, l []render.Renderer) {
	if err := render.RenderList(w, r, l); err != nil {
		err := render.Render(w, r, payloads.ErrRender(err))
		panic(err)
	}
}
