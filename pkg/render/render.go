package render

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

func JSON(w http.ResponseWriter, output any) {
	b, err := json.Marshal(output)
	if err != nil {
		log.Error().Err(err).Msg("render.JSON json.Marshal")
		http.Error(w, "marshal json error", http.StatusInternalServerError)

		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Error().Err(err).Msg("render.JSON  w.Write")
		http.Error(w, "response write error", http.StatusInternalServerError)
	}
}
