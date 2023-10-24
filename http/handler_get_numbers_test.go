package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerGetNumbers(t *testing.T) {
	expectedURLs := []string{"http://test1.com", "http://test2.com"}
	expectedNumbers := []int{1, 2, 3}

	mockGetNumbers := func(urls []string) []int {
		assert.ElementsMatch(t, expectedURLs, urls)
		return expectedNumbers
	}

	handler := NewHandlerGetNumbers(mockGetNumbers)

	req, err := http.NewRequest("GET", "/numbers?u=http://test1.com&u=http://test2.com", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string][]int
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, map[string][]int{"numbers": expectedNumbers}, response)
}
