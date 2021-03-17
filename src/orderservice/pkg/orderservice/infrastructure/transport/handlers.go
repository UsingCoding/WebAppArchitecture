package transport

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	MenuItemsIDs []uuid.UUID `json:"menu_items_ids"`
}

type createOrderResponse struct {
	ID uuid.UUID `json:"id"`
}

type createMenuItemRequest struct {
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

type createMenuItemResponse struct {
	ID uuid.UUID `json:"id"`
}

type getMenuItemRequest struct {
	ID uuid.UUID `json:"id"`
}

type getMenuItemResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string
}

func NewRouter(srv *Server) http.Handler {
	router := mux.NewRouter()
	s := router.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/order/{orderID}", srv.getOrder).Methods(http.MethodGet)
	s.HandleFunc("/orders", srv.getOrders).Methods(http.MethodGet)
	s.HandleFunc("/order", srv.createOrder).Methods(http.MethodPost)

	s.HandleFunc("/menu-item", srv.createMenuItem).Methods(http.MethodPost)
	s.HandleFunc("/menu-item/{menuItemID}", srv.getMenuItem).Methods(http.MethodGet)

	return logMiddleware(router)
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

	if len(requestData.MenuItemsIDs) == 0 {
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

	order, err := s.orderService.CreateOrder(requestData.MenuItemsIDs)
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

func (s *Server) createMenuItem(w http.ResponseWriter, req *http.Request) {
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer req.Body.Close()

	var requestData createMenuItemRequest
	err = json.Unmarshal(bytes, &requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	menuItem, err := s.menuItemService.CreateMenuItem(requestData.Name, requestData.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeJsonResponse(w, createMenuItemResponse{ID: menuItem.ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getMenuItem(w http.ResponseWriter, req *http.Request) {
	menuItemID, ok := mux.Vars(req)["menuItemID"]
	if !ok {
		setBadRequestResponse(w, "OrderId not found")
		return
	}

	id, err := uuid.Parse(menuItemID)
	if err != nil {
		setBadRequestResponse(w, "OrderId invalid")
		return
	}

	menuItemView, err := s.menuItemQueryService.GetMenuItem(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = writeJsonResponse(w, getMenuItemResponse{
		ID:   menuItemView.ID,
		Name: menuItemView.Name,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
