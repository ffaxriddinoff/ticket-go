package eticket

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (c *Client) SearchTrains(
	ctx context.Context,
	input SearchDto,
	lang string,
) ([]Train, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseUrl+"/handbook/trains/list",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("eticket: search request: %w", err)
	}

	if lang == "" {
		lang = "uz"
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("accept-language", lang)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("eticket: send search request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("eticket: search request failed with status %d", response.StatusCode)
	}

	var result SearchResponseDto

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("eticket: parse search response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("eticket: Api error: %s", string(*result.Error))
	}

	if result.Data == nil {
		return nil, errors.New("eticket: empty result")
	}

	forward := result.Data.Directions.Forward

	if forward == nil || len(forward.Trains) == 0 {
		return []Train{}, nil
	}

	return forward.Trains, nil
}
