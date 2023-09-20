package handler

import (
	"github.com/weeaa/nft/pkg/safemap"
)

type Handler struct {
	M     *safemap.SafeMap[string, interface{}]
	MCopy *safemap.SafeMap[string, interface{}]
}

// New returns an Handler. It is used to store data 🧸.
func New() *Handler {
	return &Handler{
		M:     safemap.New[string, interface{}](),
		MCopy: safemap.New[string, interface{}](),
	}
}

func (h *Handler) Copy() {
	h.M.ForEach(func(k string, v interface{}) {
		h.MCopy.Set(k, v)
	})
}
