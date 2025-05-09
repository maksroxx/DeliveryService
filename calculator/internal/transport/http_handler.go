package transport

import (
	"encoding/json"
	"net/http"

	"github.com/maksroxx/DeliveryService/calculator/internal/metrics"
	"github.com/maksroxx/DeliveryService/calculator/internal/service"
	"github.com/maksroxx/DeliveryService/calculator/models"
)

type HTTPHandler struct {
	service service.Calculator
}

func NewHTTPHandler(s service.Calculator) *HTTPHandler {
	return &HTTPHandler{service: s}
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var pkg models.Package
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		metrics.CalculationFailureTotal.WithLabelValues("POST", "decode").Inc()
		return
	}

	if pkg.Weight <= 0 {
		RespondError(w, http.StatusBadRequest, "Invalid weight")
		metrics.CalculationFailureTotal.WithLabelValues("POST", "validation_weight").Inc()
		return
	}
	if ValidateAddress(pkg) != nil {
		RespondError(w, http.StatusBadRequest, "Invalid location data")
		metrics.CalculationFailureTotal.WithLabelValues("POST", "validation_location").Inc()
		return
	}

	result, err := h.service.Calculate(pkg)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Calculation error: "+err.Error())
		metrics.CalculationFailureTotal.WithLabelValues("POST", "calculation").Inc()
		return
	}

	metrics.CalculationSuccessTotal.WithLabelValues("POST").Inc()
	metrics.CalculatedCostValue.Observe(result.Cost)

	RespondJSON(w, http.StatusOK, result)
}

func RespondJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}
