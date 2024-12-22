package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/veliashev/web-calculation/internal/application"
)

func TestCalculationHandler(t *testing.T) {
	tests := []struct {
		name           string
		expression     string
		expectedResult float64
		expectedStatus int
		wantErr        bool
	}{
		{
			name:           "simple addition",
			expression:     "2 + 2",
			expectedResult: 4,
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "subtraction",
			expression:     "10 - 5",
			expectedResult: 5,
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "multiplication",
			expression:     "3 * 4",
			expectedResult: 12,
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "division",
			expression:     "15 / 3",
			expectedResult: 5,
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "invalid expression",
			expression:     "invalid",
			expectedStatus: http.StatusUnprocessableEntity,
			wantErr:        true,
		},
		{
			name:           "division by zero",
			expression:     "10 / 0",
			expectedStatus: http.StatusUnprocessableEntity,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := application.Request{Expression: tt.expression}
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(application.CalculationHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if !tt.wantErr {
				var response application.Response
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatal(err)
				}

				if response.Result != tt.expectedResult {
					t.Errorf("handler returned wrong result: got %v want %v", response.Result, tt.expectedResult)
				}
			} else {
				var errorResponse application.Error
				err = json.Unmarshal(rr.Body.Bytes(), &errorResponse)
				if err != nil {
					t.Fatal(err)
				}

				if errorResponse.Error == "" {
					t.Error("expected error message in response, got empty string")
				}
			}
		})
	}
}
