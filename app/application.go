package application

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	calc "github.com/Danilka776/web_go_calc/calculate"

	//calc "sample-app/Desktop/service_GO/web_go_calc/calculate"
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

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (a *Application) Run() error {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Failed to read expression from console")
		}
		text = strings.TrimSpace(text)
		if text == "exit" {
			log.Println("application was successfully closed")
			return nil
		}
		result, err := calc.Calc(text)
		if err != nil {
			log.Println(text, " calculation failed")
		} else {
			log.Println(text, " = ", result)
		}
	}
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
	res, err := calc.Calc(infix)
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

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // 500
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("Internal server error")
	}
	return nil
}
