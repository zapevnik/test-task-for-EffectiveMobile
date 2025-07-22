package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"subscription-service/internal/domain"
	"subscription-service/internal/usecase/subscription"

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling Create request")

	var req CreateSubscriptionDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid request body", slog.String("error", err.Error()))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sub := &domain.Subscription{
		ID:          uuid.New(),
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	h.logger.Debug("creating subscription", slog.Any("subscription", sub))

	if err := h.service.Create(r.Context(), sub); err != nil {
		h.logger.Error("failed to create subscription", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("subscription created successfully", slog.String("id", sub.ID.String()))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling GetAll request")

	subs, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Error("failed to get all subscriptions", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("subscriptions retrieved", slog.Int("count", len(subs)))
	json.NewEncoder(w).Encode(subs)
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
	json.NewEncoder(w).Encode(sub)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	h.logger.Info("handling Update request", slog.String("id", idStr))

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("invalid UUID", slog.String("id", idStr))
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req CreateSubscriptionDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("invalid request body", slog.String("error", err.Error()))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	sub := &domain.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	h.logger.Debug("updating subscription", slog.Any("subscription", sub))

	if err := h.service.Update(r.Context(), sub); err != nil {
		h.logger.Error("failed to update subscription", slog.String("id", id.String()), slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("subscription updated successfully", slog.String("id", id.String()))
	json.NewEncoder(w).Encode(sub)
}

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
		domain.YearMonth{Time: from},
		domain.YearMonth{Time: to},
	)
	if err != nil {
		h.logger.Error("failed to calculate total cost", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("total cost calculated successfully", slog.Int64("total", total))
	json.NewEncoder(w).Encode(map[string]int64{"total": total})
}
