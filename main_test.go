package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestCafeWhenOk(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cafe?city=moscow", nil)
	w := httptest.NewRecorder()

	cafeHandler(w, req) 

	require.Equal(t, http.StatusOK, w.Code)
	body := strings.TrimSpace(w.Body.String())
	assert.NotEmpty(t, body)

	cafes := strings.Split(body, ",")
	assert.GreaterOrEqual(t, len(cafes), 1)
}


func TestCafeNegative(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantBody   string
	}{
		{"нет параметра city", "/cafe", http.StatusBadRequest, "city not found"},
		{"неверный город", "/cafe?city=spb", http.StatusBadRequest, "wrong city value"},
		
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			cafeHandler(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
			body := strings.TrimSpace(w.Body.String())
			assert.Equal(t, tt.wantBody, body)
		})
	}
}


func TestCafeCount(t *testing.T) {
	city := "moscow"
	total := len(cafeList[city]) // полное количество кафе в городе

	tests := []struct {
		name  string
		count string
		want  int
	}{
		{"count=0", "0", 0},
		{"count=1", "1", 1},
		{"count=2", "2", 2},
		{"count=100", "100", min(total, 100)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cafe?city="+city+"&count="+tt.count, nil)
			w := httptest.NewRecorder()

			cafeHandler(w, req)

			require.Equal(t, http.StatusOK, w.Code)

			body := strings.TrimSpace(w.Body.String())
			var cafes []string
			if body == "" {
				cafes = []string{}
			} else {
				cafes = strings.Split(body, ",")
			}

			assert.Equal(t, tt.want, len(cafes), "неправильное количество кафе")
		})
	}
}


func TestCafeSearch(t *testing.T) {
	city := "moscow"

	tests := []struct {
		name      string
		search    string
		wantCount int
	}{
		{"поиск 'фасоль' (нет результатов)", "фасоль", 0},
		{"поиск 'кофе' (2 результата)", "кофе", 2},
		{"поиск 'вилка' (1 результат)", "вилка", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/cafe?city="+city+"&search="+tt.search, nil)
			w := httptest.NewRecorder()

			cafeHandler(w, req)

			require.Equal(t, http.StatusOK, w.Code)

			body := strings.TrimSpace(w.Body.String())
			var cafes []string
			if body == "" {
				cafes = []string{}
			} else {
				cafes = strings.Split(body, ",")
			}

		
			assert.Equal(t, tt.wantCount, len(cafes), "неправильное количество кафе")

			
			if tt.wantCount > 0 {
				lowerSearch := strings.ToLower(tt.search)
				for _, cafe := range cafes {
					assert.Contains(t, strings.ToLower(cafe), lowerSearch,
						"кафе %q не содержит подстроку %q", cafe, tt.search)
				}
			}
		})
	}
}


func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
