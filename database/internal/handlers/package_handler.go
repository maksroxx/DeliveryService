package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/maksroxx/DeliveryService/database/internal/models"
	"github.com/maksroxx/DeliveryService/database/internal/repository"
)

type PackageHandler struct {
	rep repository.RouteRepository
}

func NewPackageHandler(rep repository.RouteRepository) *PackageHandler {
	return &PackageHandler{
		rep: rep,
	}
}

func (h *PackageHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/packages/{packageID}", h.GetPackage).Methods("GET")
	router.HandleFunc("/packages", h.GetAllPackages).Methods("GET")
	router.HandleFunc("/packages", h.CreatePackage).Methods("POST")
	router.HandleFunc("/packages/{packageID}", h.UpdatePackage).Methods("PUT")
	router.HandleFunc("/packages/{packageID}", h.DeletePackage).Methods("DELETE")
}

func (h *PackageHandler) GetPackage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	packageID := vars["packageID"]

	pkg, err := h.rep.GetByID(r.Context(), packageID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Package not found")
		return
	}

	respondWithJSON(w, http.StatusOK, pkg)
}

func (h *PackageHandler) GetAllPackages(w http.ResponseWriter, r *http.Request) {
	filter := models.RouteFilter{
		Status: r.URL.Query().Get("status"),
	}

	if createdAfter := r.URL.Query().Get("created_after"); createdAfter != "" {
		if t, err := time.Parse(time.RFC3339, createdAfter); err == nil {
			filter.CreatedAfter = t
		}
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.ParseInt(limit, 10, 64); err == nil {
			filter.Limit = l
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.ParseInt(offset, 10, 64); err == nil {
			filter.Offset = o
		}
	}

	packages, err := h.rep.GetAllRoutes(r.Context(), filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	respondWithJSON(w, http.StatusOK, packages)
}

func (h *PackageHandler) CreatePackage(w http.ResponseWriter, r *http.Request) {
	var req models.Route
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	pkg, err := h.rep.Create(r.Context(), &req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create package")
		return
	}

	respondWithJSON(w, http.StatusCreated, pkg)
}

func (h *PackageHandler) UpdatePackage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	packageID := vars["packageID"]

	var update models.RouteUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedPkg, err := h.rep.UpdateRoute(r.Context(), packageID, update)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update package")
		return
	}

	respondWithJSON(w, http.StatusOK, updatedPkg)
}

func (h *PackageHandler) DeletePackage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	packageID := vars["packageID"]

	if err := h.rep.DeleteRoute(r.Context(), packageID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete package")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
