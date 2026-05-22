package pkg

import (
	"encoding/json"
	"net/http"
)

type HTTPRouter interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

func RegisterHTTPAPI(mux HTTPRouter) {
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, Message{Message: "ok"})
	})
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			WriteJSON(w, http.StatusOK, []User{{ID: "usr_1", Name: "Coco", Email: "coco@example.com"}})
		case http.MethodPost:
			var req CreateUserRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				WriteJSON(w, http.StatusBadRequest, Message{Message: "invalid request"})
				return
			}
			WriteJSON(w, http.StatusCreated, User{ID: "usr_2", Name: req.Name, Email: req.Email})
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
