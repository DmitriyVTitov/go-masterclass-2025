package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"ugc/internal/api/middleware"
	"ugc/internal/config"
	"ugc/internal/errs"
	"ugc/internal/reviews"
	"ugc/internal/types"
	docs "ugc/openapi"
)

type API struct {
	cfg     config.Config
	reviews *reviews.Service

	router *chi.Mux
}

func New(cfg config.Config, reviews *reviews.Service) *API {
	api := API{
		cfg:     cfg,
		reviews: reviews,
		router:  chi.NewRouter(),
	}

	api.router.Use(middleware.RequestID)
	api.router.Use(middleware.Telemetry)
	api.router.Use(middleware.Metrics)

	api.registerEndpoints()

	return &api
}

func (api *API) Serve(ctx context.Context) error {
	log.Info().Msgf("HTTP server listen on port: %v", api.cfg.Server.ListenPort)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", api.cfg.Server.ListenPort),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		Handler:      api.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return srv.ListenAndServe()
}

// @title Go Masterclass
// @version 1.0
// @contact.name Dmitriy Titov
// @BasePath /api/v1
// @Produce      json
func (api *API) registerEndpoints() {
	// Profiler
	api.router.HandleFunc("/debug/pprof/", pprof.Index)
	api.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	api.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	api.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	api.router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	api.router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	api.router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	api.router.Handle("/debug/pprof/block", pprof.Handler("block"))

	// Prometheus metrics
	api.router.Handle("/metrics", promhttp.Handler())

	// K8s probes
	api.router.HandleFunc("/ready", api.readinessProbeHandle)
	api.router.HandleFunc("/alive", api.livenessProbeHandle)

	// Swagger
	docs.SwaggerInfo.BasePath = api.cfg.AppBasePath
	api.router.Get("/swagger/*", func(w http.ResponseWriter, r *http.Request) {
		swaggerHandler := httpSwagger.Handler(
			httpSwagger.URL(
				fmt.Sprintf("%s/swagger/doc.json", api.cfg.AppBaseURL),
			),
		)
		swaggerHandler.ServeHTTP(w, r)
	})

	// Application handlers
	api.router.Get("/api/v1/reviews/object/{id}", api.objectReviewsHandler)
	api.router.Post("/api/v1/reviews", api.addReviewHandler)
}

// objectReviewsHandler godoc
// @Summary      Get revies by object ID
// @Tags         get
// @Param        id   query    int  true  "Object ID"
// @Success      200  {array}  types.Review
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /reviews/object/{id} [get]
func (api *API) objectReviewsHandler(w http.ResponseWriter, r *http.Request) {
	s := chi.URLParam(r, "id")
	id, err := strconv.Atoi(s)
	if err != nil {
		span := trace.SpanFromContext(r.Context())
		e := errs.NewErrBadRequest(err.Error())
		span.SetStatus(codes.Error, e.Error())
		api.WriteError(w, r, e)
		return
	}

	resp, err := api.reviews.ObjectReviews(r.Context(), id)
	if err != nil {
		span := trace.SpanFromContext(r.Context())
		span.SetStatus(codes.Error, err.Error())
		api.WriteError(w, r, err)
		return
	}

	api.WriteJSON(w, r, resp)
}

// addReviewHandler godoc
// @Summary      Add Review
// @Tags         post
// @Param        request   body types.Review  true  "Review"
// @Success      200
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /reviews [post]
func (api *API) addReviewHandler(w http.ResponseWriter, r *http.Request) {
	var req types.Review
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		api.WriteError(w, r, err)
		return
	}

	err = api.reviews.AddReview(r.Context(), req)
	if err != nil {
		api.WriteError(w, r, err)
		return
	}
}
