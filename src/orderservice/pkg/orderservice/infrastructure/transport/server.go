package transport

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"orderservice/pkg/orderservice/app/query"
	"orderservice/pkg/orderservice/domain/service"
	"time"
)

type Server struct {
	orderService         service.OrderService
	orderQueryService    query.OrderQueryService
	menuItemService      service.MenuItemService
	menuItemQueryService query.MenuItemQueryService
}

func NewServer(
	service service.OrderService,
	queryService query.OrderQueryService,
	menuItemService service.MenuItemService,
	menuItemQueryService query.MenuItemQueryService,
) *Server {
	return &Server{
		orderService:         service,
		orderQueryService:    queryService,
		menuItemService:      menuItemService,
		menuItemQueryService: menuItemQueryService,
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
