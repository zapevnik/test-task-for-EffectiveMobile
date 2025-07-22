package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"subscription-service/internal/usecase/subscription"
	"subscription-service/internal/delivery/dto"
	"subscription-service/pkg/dtoConv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service *subscription.Service
	logger  *slog.Logger
}

func NewHandler(service *subscription.Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}
// Create godoc
// @Summary Create subscription
// @Description Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body dto.SubscriptionRequestDTO true "Subscription request"
// @Success 201 {object} dto.SubscriptionResponseDTO
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal error"
// @Router /subscriptions [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling Create request")

	var req dto.SubscriptionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid request body", slog.String("error", err.Error()))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := dtoConv.RequestDtoToDomain(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sub.ID = uuid.New()

	h.logger.Debug("creating subscription", slog.Any("subscription", sub))

	if err := h.service.Create(r.Context(), sub); err != nil {
		h.logger.Error("failed to create subscription", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("subscription created successfully", slog.String("id", sub.ID.String()))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dtoConv.DomainToResponseDTO(sub))
}

// GetAll godoc
// @Summary Get all subscriptions
// @Description Get list of all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} dto.SubscriptionResponseDTO
// @Failure 500 {string} string "internal error"
// @Router /subscriptions [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling GetAll request")

	subs, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Error("failed to get all subscriptions", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result []dto.SubscriptionResponseDTO
	for _, sub := range subs {
		result = append(result, dtoConv.DomainToResponseDTO(sub))
	}

	h.logger.Info("subscriptions retrieved", slog.Int("count", len(result)))
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	h.logger.Info("handling GetByID request", slog.String("id", idStr))

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("invalid UUID", slog.String("id", idStr))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	sub, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("subscription not found", slog.String("id", idStr), slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info("subscription found", slog.String("id", idStr))
	json.NewEncoder(w).Encode(dtoConv.DomainToResponseDTO(sub))
}

// Update godoc
// @Summary Update subscription
// @Description Update a subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body dto.SubscriptionRequestDTO true "Subscription update"
// @Success 200 {object} dto.SubscriptionResponseDTO
// @Failure 400 {string} string "invalid input"
// @Failure 500 {string} string "internal error"
// @Router /subscriptions/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	h.logger.Info("handling Update request", slog.String("id", idStr))

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("invalid UUID", slog.String("id", idStr))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req dto.SubscriptionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid request body", slog.String("error", err.Error()))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sub, err := dtoConv.RequestDtoToDomain(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sub.ID = id

	h.logger.Debug("updating subscription", slog.Any("subscription", sub))

	if err := h.service.Update(r.Context(), sub); err != nil {
		h.logger.Error("failed to update subscription", slog.String("id", id.String()), slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("subscription updated successfully", slog.String("id", id.String()))
	json.NewEncoder(w).Encode(dtoConv.DomainToResponseDTO(sub))
}

// Delete godoc
// @Summary Delete subscription
// @Description Delete a subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "invalid id"
// @Failure 500 {string} string "internal error"
// @Router /subscriptions/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	h.logger.Info("handling Delete request", slog.String("id", idStr))

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("invalid UUID", slog.String("id", idStr))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.logger.Error("failed to delete subscription", slog.String("id", id.String()), slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("subscription deleted successfully", slog.String("id", id.String()))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "subscription deleted successfully"}`))
}

// TotalCost godoc
// @Summary Calculate total subscription cost
// @Description Calculate total subscription cost in date range, optional filters by user and service
// @Tags subscriptions
// @Produce json
// @Param from query string true "Start month in MM-YYYY format"
// @Param to query string true "End month in MM-YYYY format"
// @Param user_id query string false "User UUID"
// @Param service_name query string false "Service name"
// @Success 200 {object} map[string]int64
// @Failure 400 {string} string "invalid input"
// @Failure 500 {string} string "internal error"
// @Router /subscriptions/total [get]
func (h *Handler) TotalCost(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling TotalCost request")

	query := r.URL.Query()
	fromStr := query.Get("from")
	toStr := query.Get("to")

	from, err := time.Parse("01-2006", fromStr)
	if err != nil {
		h.logger.Warn("invalid 'from' date", slog.String("value", fromStr))
		http.Error(w, "invalid from date format, expected MM-YYYY", http.StatusBadRequest)
		return
	}

	to, err := time.Parse("01-2006", toStr)
	if err != nil {
		h.logger.Warn("invalid 'to' date", slog.String("value", toStr))
		http.Error(w, "invalid to date format, expected MM-YYYY", http.StatusBadRequest)
		return
	}

	var userID *uuid.UUID
	if userIDStr := query.Get("user_id"); userIDStr != "" {
		uid, err := uuid.Parse(userIDStr)
		if err != nil {
			h.logger.Warn("invalid user_id", slog.String("value", userIDStr))
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		userID = &uid
	}

	var serviceName *string
	if sn := query.Get("service_name"); sn != "" {
		serviceName = &sn
	}

	h.logger.Debug("calculating total cost", slog.String("from", fromStr), slog.String("to", toStr))

	total, err := h.service.TotalCost(
		r.Context(),
		userID,
		serviceName,
		from,
		to,
	)
	if err != nil {
		h.logger.Error("failed to calculate total cost", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("total cost calculated successfully", slog.Int64("total", total))
	json.NewEncoder(w).Encode(map[string]int64{"total": total})
}
