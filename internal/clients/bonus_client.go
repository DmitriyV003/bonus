package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/applicationerrors"
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
	Order  string  `json:"order,omitempty"`
	Status string  `json:"status,omitempty"`
	Amount float64 `json:"accrual,omitempty"`
}

func NewBonusClient(url string) *BonusClient {
	return &BonusClient{
		client: &http.Client{},
		url:    url,
	}
}

func (bc *BonusClient) CreateOrder(orderNumber string) (*Response, error) {
	data := createOrderRequest{Order: orderNumber}
	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal json: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, bc.getURL("api/orders"), bytes.NewBuffer(byteData))
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("unable to create new request: %w", err)
	}
	request.Header.Add("Content-Type", "application/json")

	res, err := bc.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to send request [POST] to /api/orders/: %w", applicationerrors.ErrServiceUnavailable)
	}
	defer res.Body.Close()
	log.Info().Fields(map[string]interface{}{
		"order_number": orderNumber,
	}).Msg("order created in black box")

	return &Response{Code: res.StatusCode}, nil
}

func (bc *BonusClient) GetOrderDetails(orderNumber string) (*OrderDetailsResponse, error) {
	request, _ := http.NewRequest(http.MethodGet, bc.getURL(fmt.Sprintf("api/orders/%s", orderNumber)), nil)

	res, err := bc.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to send request [GET] to /api/orders/: %w", applicationerrors.ErrServiceUnavailable)
	}
	defer res.Body.Close()

	var response OrderDetailsResponse
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	if res.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusAccepted {
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal json with order details: %w", err)
		}
	}

	return &response, nil
}

func (bc *BonusClient) getURL(url string) string {
	return fmt.Sprintf("%s/%s", bc.url, url)
}
