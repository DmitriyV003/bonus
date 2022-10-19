package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DmitriyV003/bonus/cmd/gophermart/models"
	"io"
	"net/http"
)

type BonusClient struct {
	client *http.Client
	url    string
}

type createOrderRequest struct {
	Order string `json:"order"`
}

type Response struct {
	Code int
}

type OrderDetailsResponse struct {
	Response
	Order  string
	Status string
	Amount float64
}

func NewBonusClient() *BonusClient {
	return &BonusClient{
		client: &http.Client{},
		url:    "http://localhost:8080",
	}
}

func (bc *BonusClient) CreateOrder(order *models.Order) (*Response, error) {
	data := createOrderRequest{Order: order.Number}
	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	request, _ := http.NewRequest(http.MethodPost, bc.getUrl("/api/orders"), bytes.NewBuffer(byteData))
	request.Header.Add("Content-Type", "application/json")

	res, err := bc.client.Do(request)
	if err != nil {
		return nil, err
	}

	return &Response{Code: res.StatusCode}, nil
}

func (bc *BonusClient) GetOrderDetails(order *models.Order) (*OrderDetailsResponse, error) {
	request, _ := http.NewRequest(http.MethodGet, bc.getUrl("/api/orders/"+order.Number), nil)

	res, err := bc.client.Do(request)
	if err != nil {
		return nil, err
	}

	response := OrderDetailsResponse{}
	response.Response = Response{Code: res.StatusCode}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (bc *BonusClient) getUrl(url string) string {
	return fmt.Sprintf("%s/%s", bc.url, url)
}
