package backend

import (
	"fmt"
	"log"
	"net/http"
)

type handler struct {
	esService *esService
}

func NewHandler(esService *esService) *handler {
	return &handler{
		esService: esService,
	}
}

const queryParamData = "data"

func (h *handler) SearchMessages(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get(queryParamData)

	if data == "" {
		log.Println("returning all messages")

		fmt.Fprint(w, string(messagesJSON))

		w.WriteHeader(http.StatusOK)

		return
	}

	err := h.esService.Search(r.Context(), data)
	if err != nil {
		log.Printf("search messages: %v", err)

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

