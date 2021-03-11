package transport

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type order struct {
	ID    string `json:"id"`
	Items []orderItem
}

type orderItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type getOrderResponse struct {
	order
	OrderedAtTimeStamp string `json:"orderedAtTimeStamp"`
	Cost               int    `json:"cost"`
}

type getOrdersResponse struct {
	Orders []order `json:"orders"`
}

type createOrderRequest struct {
	Items []orderItem `json:"menuItems"`
}

type createOrderResponse struct {
	ID uuid.UUID `json:"id"`
}

func NewRouter(srv *Server) http.Handler {
	router := mux.NewRouter()
	s := router.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/order/{orderID}", srv.getOrder).Methods(http.MethodGet)
	s.HandleFunc("/orders", srv.getOrders).Methods(http.MethodGet)
	s.HandleFunc("/order", srv.createOrder).Methods(http.MethodPost)

	return logMiddleware(router)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		now := time.Now()
		h.ServeHTTP(writer, request)

		log.WithFields(log.Fields{
			"duration":   time.Since(now),
			"method":     request.Method,
			"url":        request.URL,
			"remoteAddr": request.RemoteAddr,
			"userAgent":  request.UserAgent(),
		}).Info("starting the server")
	})
}

func writeJsonResponse(w http.ResponseWriter, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = io.WriteString(w, string(bytes))
	return err
}

func setBadRequestResponse(w http.ResponseWriter, err string) {
	http.Error(w, err, http.StatusBadRequest)
}
