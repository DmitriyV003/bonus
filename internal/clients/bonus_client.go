package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
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
	Order  string  `json:"order,omitempty"`
	Status string  `json:"status,omitempty"`
	Amount float64 `json:"accrual,omitempty"`
}

func NewBonusClient() *BonusClient {
	return &BonusClient{
		client: &http.Client{},
		url:    "http://localhost:8080",
	}
}

func (bc *BonusClient) CreateOrder(orderNumber string) (*Response, error) {
	data := createOrderRequest{Order: orderNumber}
	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, bc.getUrl("api/orders"), bytes.NewBuffer(byteData))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")

	res, err := bc.client.Do(request)
	if err != nil {
		return nil, err
	}
	log.Info().Fields(map[string]interface{}{
		"response": res,
	}).Msgf("order created in black box: ", res.Status)

	return &Response{Code: res.StatusCode}, nil
}

func (bc *BonusClient) GetOrderDetails(orderNumber string) (*OrderDetailsResponse, error) {
	request, _ := http.NewRequest(http.MethodGet, bc.getUrl(fmt.Sprintf("api/orders/%s", orderNumber)), nil)

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
	log.Info().Fields(map[string]interface{}{
		"response": body,
	}).Msg("Order details from service")

	return &response, nil
}

func (bc *BonusClient) getUrl(url string) string {
	return fmt.Sprintf("%s/%s", bc.url, url)
}
