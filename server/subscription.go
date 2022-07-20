package server

import (
	"net/http"

	"github.com/frain-dev/convoy/server/models"

	"github.com/frain-dev/convoy/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"

	m "github.com/frain-dev/convoy/internal/pkg/middleware"
)

// GetSubscriptions
// @Summary Get all subscriptions
// @Description This endpoint fetches all the subscriptions
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param perPage query string false "results per page"
// @Param page query string false "page number"
// @Param sort query string false "sort order"
// @Param q query string false "subscription title"
// @Param groupId query string true "group id"
// @Success 200 {object} serverResponse{data=pagedResponse{content=[]datastore.Subscription}}
// @Failure 400,401,500 {object} serverResponse{data=Stub}
// @Security ApiKeyAuth
// @Router /subscriptions [get]
func (s *Server) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	pageable := m.GetPageableFromContext(r.Context())
	group := m.GetGroupFromContext(r.Context())

	apps, paginationData, err := s.subService.LoadSubscriptionsPaged(r.Context(), group.UID, pageable)
	if err != nil {
		log.WithError(err).Error("failed to load subscriptions")
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	_ = render.Render(w, r, util.NewServerResponse("Subscriptions fetched successfully",
		pagedResponse{Content: &apps, Pagination: &paginationData}, http.StatusOK))
}

// GetSubscription
// @Summary Gets a subscription
// @Description This endpoint fetches an Subscription by it's id
// @Tags Subscription
// @Accept json
// @Produce  json
// @Param groupId query string true "group id"
// @Param subscriptionID path string true "application id"
// @Success 200 {object} serverResponse{data=datastore.Subscription}
// @Failure 400,401,500 {object} serverResponse{data=Stub}
// @Security ApiKeyAuth
// @Router /subscriptions/{subscriptionID} [get]
func (s *Server) GetSubscription(w http.ResponseWriter, r *http.Request) {
	subId := chi.URLParam(r, "subscriptionID")
	group := m.GetGroupFromContext(r.Context())

	subscription, err := s.subService.FindSubscriptionByID(r.Context(), group, subId, false)
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	_ = render.Render(w, r, util.NewServerResponse("Subscription fetched successfully", subscription, http.StatusOK))
}

// CreateSubscription
// @Summary Creates a subscription
// @Description This endpoint creates a subscriptions
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param groupId query string true "group id"
// @Success 200 {object} serverResponse{data=pagedResponse{content=[]datastore.Subscription}}
// @Failure 400,401,500 {object} serverResponse{data=Stub}
// @Security ApiKeyAuth
// @Router /subscriptions [post]
func (s *Server) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	group := m.GetGroupFromContext(r.Context())

	var sub models.Subscription
	err := util.ReadJSON(r, &s)
	if err != nil {
		_ = render.Render(w, r, util.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	sub.Type = string(group.Type)

	subscription, err := s.subService.CreateSubscription(r.Context(), group, &sub)
	if err != nil {
		log.WithError(err).Error("failed to create subscription")
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	_ = render.Render(w, r, util.NewServerResponse("Subscriptions created successfully", subscription, http.StatusCreated))
}

// DeleteSubscription
// @Summary Delete subscription
// @Description This endpoint deletes a subscription
// @Tags Application
// @Accept json
// @Produce json
// @Param groupId query string true "group id"
// @Param subscriptionID path string true "subscription id"
// @Success 200 {object} serverResponse{data=Stub}
// @Failure 400,401,500 {object} serverResponse{data=Stub}
// @Security ApiKeyAuth
// @Router /subscriptions/{subscriptionID} [delete]
func (s *Server) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	group := m.GetGroupFromContext(r.Context())

	sub, err := s.subService.FindSubscriptionByID(r.Context(), group, chi.URLParam(r, "subscriptionID"), true)
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	err = s.subService.DeleteSubscription(r.Context(), group.UID, sub)
	if err != nil {
		log.Errorln("failed to delete subscription - ", err)
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	_ = render.Render(w, r, util.NewServerResponse("Subscription deleted successfully", nil, http.StatusOK))
}

// UpdateSubscription
// @Summary Update a subscription
// @Description This endpoint updates a subscription
// @Tags Subscription
// @Accept json
// @Produce json
// @Param subscriptionID path string true "subscription id"
// @Param subscription body models.Subscription true "Subscription Details"
// @Success 200 {object} serverResponse{data=datastore.Subscription}
// @Failure 400,401,500 {object} serverResponse{data=Stub}
// @Security ApiKeyAuth
// @Router /subscriptions/{subscriptionID} [put]
func (s *Server) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var update models.UpdateSubscription
	err := util.ReadJSON(r, &update)
	if err != nil {
		log.WithError(err).Error(err.Error())
		_ = render.Render(w, r, util.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	g := m.GetGroupFromContext(r.Context())
	subscription := chi.URLParam(r, "subscriptionID")

	sub, err := s.subService.UpdateSubscription(r.Context(), g.UID, subscription, &update)
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	_ = render.Render(w, r, util.NewServerResponse("Subscription updated successfully", sub, http.StatusAccepted))
}

// ToggleSubscriptionStatus
// @Summary Toggles a subscription's status from active <-> inactive
// @Description This endpoint updates a subscription
// @Tags Subscription
// @Accept json
// @Produce json
// @Param subscriptionID path string true "subscription id"
// @Success 200 {object} serverResponse{data=datastore.Subscription}
// @Failure 400,401,500 {object} serverResponse{data=Stub}
// @Security ApiKeyAuth
// @Router /subscriptions/{subscriptionID}/toggle_status [put]
func (s *Server) ToggleSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	g := m.GetGroupFromContext(r.Context())
	subscription := chi.URLParam(r, "subscriptionID")

	sub, err := s.subService.ToggleSubscriptionStatus(r.Context(), g.UID, subscription)
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	_ = render.Render(w, r, util.NewServerResponse("Subscription status updated successfully", sub, http.StatusAccepted))
}
