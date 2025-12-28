package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tikimcrzx723/alejandrinasweb/internal/dtos"
)

func GetAllCategories(ctx context.Context, baseURL string) (dtos.CategoryResponse, error) {
	if strings.TrimSpace(baseURL) == "" {
		return dtos.CategoryResponse{}, fmt.Errorf("baseURL is required")
	}

	url := strings.TrimRight(baseURL, "/") + "/categories"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return dtos.CategoryResponse{}, fmt.Errorf("create get categories request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return dtos.CategoryResponse{}, fmt.Errorf("send get categories request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return dtos.CategoryResponse{}, fmt.Errorf("get categories failed with status code: %d", resp.StatusCode)
	}

	var categoryResp dtos.CategoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&categoryResp); err != nil {
		fmt.Println(err)
		return dtos.CategoryResponse{}, fmt.Errorf("decode get categories response: %w", err)
	}

	return categoryResp, nil
}

func CreateCategory(
	ctx context.Context,
	baseURL string,
	category dtos.CreateCategoryRequest,
	token string,
) (dtos.SingleCategoryResponse, error) {
	if strings.TrimSpace(baseURL) == "" {
		return dtos.SingleCategoryResponse{}, fmt.Errorf("baseURL is required")
	}

	url := strings.TrimRight(baseURL, "/") + "/categories"

	payloadBytes, err := json.Marshal(category)
	if err != nil {
		return dtos.SingleCategoryResponse{}, fmt.Errorf("marshal create category payload: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return dtos.SingleCategoryResponse{}, fmt.Errorf("create create category request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return dtos.SingleCategoryResponse{}, fmt.Errorf("send create category request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return dtos.SingleCategoryResponse{}, fmt.Errorf("create category failed with status code: %d", resp.StatusCode)
	}

	var categoryResp dtos.SingleCategoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&categoryResp); err != nil {
		return dtos.SingleCategoryResponse{}, fmt.Errorf("decode create category response: %w", err)
	}

	return categoryResp, nil
}
