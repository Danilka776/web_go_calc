package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type RequestBody struct {
	Expression string `json:"expression"`
}

type ResponseBody struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	// Используется не тот метод (не POST)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusInternalServerError) // 500
		return
	}

	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Expression is not valid"})
		return
	}

	infix := strings.TrimSpace(reqBody.Expression)
	res, err := Calc(infix)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Expression is not valid"})
		return
	}
	response := ResponseBody{Result: res}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // 500
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
	})
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic("1Internal server error")
	}
}
