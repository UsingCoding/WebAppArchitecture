package transport

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
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

func NewRouter() http.Handler {
	router := mux.NewRouter()
	s := router.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/order/{orderID}", getOrder).Methods(http.MethodGet)
	s.HandleFunc("/orders", getOrders).Methods(http.MethodGet)
	s.HandleFunc("/order", createOrder).Methods(http.MethodPost)

	return logMiddleware(router)
}

func getOrder(w http.ResponseWriter, request *http.Request) {
	orderId, ok := mux.Vars(request)["orderID"]
	if !ok {
		setBadRequestResponse(w, "OrderId not found")
		return
	}

	response := getOrderResponse{
		order: order{
			ID: orderId,
			Items: []orderItem{
				{
					Id:       orderId,
					Quantity: 25,
				},
			},
		},
		OrderedAtTimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		Cost:               999,
	}

	err := writeJsonResponse(w, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getOrders(w http.ResponseWriter, _ *http.Request) {
	response := getOrdersResponse{
		Orders: []order{
			{
				ID: "d290f1ee-6c56-4b01-90e6-d701748f0851",
				Items: []orderItem{{
					Id:       "f290d1ce-6c234-4b31-90e6-d701748f0851",
					Quantity: 1,
				}},
			},
		},
	}

	err := writeJsonResponse(w, response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createOrder(w http.ResponseWriter, req *http.Request) {
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer req.Body.Close()

	var requestData createOrderRequest
	err = json.Unmarshal(bytes, &requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(requestData.Items) == 0 {
		response := struct {
			Message string `json:"message"`
		}{
			Message: "empty order items",
		}
		err = writeJsonResponse(w, response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	orderID := uuid.New()

	err = writeJsonResponse(w, createOrderResponse{ID: orderID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
