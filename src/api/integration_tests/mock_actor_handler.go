package integration_tests

import (
	"encoding/json"
	"net/http"
)

type RequestPayload struct {
	ActorIds []string `json:"actor_ids"`
}

type ResponsePayload struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func ActorHandler(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var result []ResponsePayload

	for _, actorId := range payload.ActorIds {
		if actorId != ValidActor1Id && actorId != ValidActor2Id {
			w.WriteHeader(http.StatusNotFound)
		}
		result = append(result, ResponsePayload{Id: actorId, Name: "Jhon"})
	}

	resultJson, _ := json.Marshal(&result)
	w.WriteHeader(http.StatusOK)
	w.Write(resultJson)
}
