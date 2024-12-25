package application

import (
	"bytes"
	"net/http/httptest"
	"testing"
)

func TestCalcHandler(t *testing.T) {
	type args struct {
		expression string
		code       int
		result     float64
		err        string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "OK",
			args: args{
				expression: "5+10*3",
				code:       200,
				result:     35,
			},
		},
		{
			name: "Invalid expression",
			args: args{
				expression: "2**2",
				code:       422,
				err:        `Expression is not valid`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewReader([]byte("{\"expression\": \""+test.args.expression+"\"}")))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			calculateHandler(recorder, req)
			if recorder.Code != test.args.code {
				t.Errorf("excepted status code %d, got %d", test.args.code, recorder.Code)
			}
		})
	}
}
