package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"orderservice/pkg/orderservice/app/query"
	"orderservice/pkg/orderservice/domain/service"
	"strconv"
	"time"
)

type Server struct {
	orderService      service.OrderService
	orderQueryService query.OrderQueryService
}

func NewServer(service service.OrderService, queryService query.OrderQueryService) *Server {
	return &Server{
		orderService:      service,
		orderQueryService: queryService,
	}
}

func (s *Server) getOrder(w http.ResponseWriter, request *http.Request) {
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

func (s *Server) getOrders(w http.ResponseWriter, _ *http.Request) {
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

func (s *Server) createOrder(w http.ResponseWriter, req *http.Request) {
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

	order, err := s.orderService.CreateOrder(nil, 50)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeJsonResponse(w, createOrderResponse{ID: order.ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
