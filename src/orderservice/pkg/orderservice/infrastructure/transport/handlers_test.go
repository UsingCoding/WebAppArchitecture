package transport

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrderById(t *testing.T) {
	w := httptest.NewRecorder()
	orderID := "ADC3412123"
	r := httptest.NewRequest(http.MethodGet, "/api/v1/order/{orderID}", nil)
	r = mux.SetURLVars(r, map[string]string{"orderID": orderID})
	getOrder(w, r)

	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Incorrect response code. Got: %d, expected: %d", response.StatusCode, http.StatusOK)
		return
	}

	bytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Fatal(err)
		return
	}

	orderResponse := getOrderResponse{}
	if err = json.Unmarshal(bytes, &orderResponse); err != nil {
		t.Errorf("Failed to pasrse json response %s", err)
		return
	}

	if orderResponse.ID != orderID {
		t.Errorf("Invalid response, orderID from response must be %s, got: %s", orderID, orderResponse.ID)
	}
}

func TestGetOrders(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	getOrders(w, r)

	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Incorrect response code. Got: %d, expected: %d", response.StatusCode, http.StatusOK)
		return
	}

	bytes, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Fatal(err)
		return
	}

	var ordersResponse getOrdersResponse
	if err = json.Unmarshal(bytes, &ordersResponse); err != nil {
		t.Errorf("Failed to pasrse json response %s", err)
	}
}
