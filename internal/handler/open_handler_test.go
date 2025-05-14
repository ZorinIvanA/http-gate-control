package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZorinIvanA/http-gate-control/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestOpenHandler(t *testing.T) {
	mockSvc := new(mock.MockGateService)
	mockSvc.On("ProcessOpen").Return(true, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/open", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	OpenHandler(mockSvc)(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"opened"`)
}
