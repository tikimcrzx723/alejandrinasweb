package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/tikimcrzx723/alejandrinasweb/internal/dtos"
)

func CreateProduct(
	ctx context.Context,
	baseURL string,
	product dtos.CreateProductRequest,
	token string,
) (dtos.SingleProductResponse, error) {
	if strings.TrimSpace(baseURL) == "" {
		return dtos.SingleProductResponse{}, fmt.Errorf("baseURL is required")
	}

	url := strings.TrimRight(baseURL, "/") + "/products"

	payloadBytes, err := json.Marshal(product)
	if err != nil {
		return dtos.SingleProductResponse{}, fmt.Errorf("marshal create product payload: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return dtos.SingleProductResponse{}, fmt.Errorf("create create product request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return dtos.SingleProductResponse{}, fmt.Errorf("send create product request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return dtos.SingleProductResponse{}, fmt.Errorf("create product failed with status code: %d	", resp.StatusCode)
	}

	var productResp dtos.SingleProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&productResp); err != nil {
		return dtos.SingleProductResponse{}, fmt.Errorf("decode create product response: %w", err)
	}

	return productResp, nil
}

func GetProducts(
	ctx context.Context,
	baseURL string,
) (dtos.ProductResponse, error) {
	if strings.TrimSpace(baseURL) == "" {
		return dtos.ProductResponse{}, fmt.Errorf("baseURL is required")
	}

	url := strings.TrimRight(baseURL, "/") + "/products"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return dtos.ProductResponse{}, fmt.Errorf("create get products request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return dtos.ProductResponse{}, fmt.Errorf("send get products request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return dtos.ProductResponse{}, fmt.Errorf("get products failed with status code: %d	", resp.StatusCode)
	}

	var productsResp dtos.ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&productsResp); err != nil {
		return dtos.ProductResponse{}, fmt.Errorf("decode get products response: %w", err)
	}

	return productsResp, nil
}

func GetProductBySKU(
	ctx context.Context,
	baseURL string,
	sku string,
) (dtos.SingleProductResponse, error) {
	if strings.TrimSpace(baseURL) == "" {
		return dtos.SingleProductResponse{}, fmt.Errorf("baseURL is required")
	}

	url := strings.TrimRight(baseURL, "/") + fmt.Sprintf("/products/sku/%s", sku)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return dtos.SingleProductResponse{}, fmt.Errorf("create get product request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return dtos.SingleProductResponse{}, fmt.Errorf("send get product request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return dtos.SingleProductResponse{}, fmt.Errorf("get product failed with status code: %d	", resp.StatusCode)
	}

	var productResp dtos.SingleProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&productResp); err != nil {
		return dtos.SingleProductResponse{}, fmt.Errorf("decode get product response: %w", err)
	}

	return productResp, nil
}

func AddProductImages(
	ctx context.Context,
	baseURL string,
	productID int,
	images []*multipart.FileHeader,
	token string,
) (dtos.ProductImagesResponse, error) {
	if strings.TrimSpace(baseURL) == "" {
		return dtos.ProductImagesResponse{}, fmt.Errorf("baseURL is required")
	}

	url := strings.TrimRight(baseURL, "/") + fmt.Sprintf("/products/%d/images", productID)
	var lastResp dtos.ProductImagesResponse
	var errMsgs []string

	for idx, image := range images {
		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)

		part, err := writer.CreateFormFile("image", image.Filename)
		if err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("image %d: create form file: %v", idx, err))
			continue
		}

		file, err := image.Open()
		if err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("image %d: open file: %v", idx, err))
			continue
		}
		if _, err := io.Copy(part, file); err != nil {
			file.Close()
			errMsgs = append(errMsgs, fmt.Sprintf("image %d: copy to form file: %v", idx, err))
			continue
		}
		file.Close()

		if err := writer.Close(); err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("image %d: close multipart writer: %v", idx, err))
			continue
		}

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &requestBody)
		if err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("image %d: create request: %v", idx, err))
			continue
		}
		httpReq.Header.Set("Content-Type", writer.FormDataContentType())
		httpReq.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("image %d: send request: %v", idx, err))
			continue
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			msg := strings.TrimSpace(string(bodyBytes))
			if msg == "" {
				msg = "empty response body"
			}
			errMsgs = append(errMsgs, fmt.Sprintf("image %d: status code %d: %s", idx, resp.StatusCode, msg))
			continue
		}

		var imagesResp dtos.ProductImagesResponse
		if err := json.Unmarshal(bodyBytes, &imagesResp); err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("image %d: decode response: %v", idx, err))
			continue
		}
		lastResp = imagesResp
	}

	if len(errMsgs) > 0 {
		return lastResp, fmt.Errorf("add product images completed with errors: %s", strings.Join(errMsgs, "; "))
	}
	return lastResp, nil
}

func UpdateProduct(
	ctx context.Context,
	baseURL string,
	token string,
	id int,
	product dtos.UpdateProductRequest,
) (dtos.SingleProductResponse, error) {
	if strings.TrimSpace(baseURL) == "" {
		return dtos.SingleProductResponse{}, fmt.Errorf("baseURL is required")
	}

	url := strings.TrimRight(baseURL, "/") + fmt.Sprintf("/products/%d", id)
	payloadBytes, err := json.Marshal(product)
	if err != nil {
		return dtos.SingleProductResponse{}, fmt.Errorf("marshal update product payload: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(payloadBytes))
	if err != nil {
		return dtos.SingleProductResponse{}, fmt.Errorf("update product request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return dtos.SingleProductResponse{}, fmt.Errorf("send update product request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return dtos.SingleProductResponse{}, fmt.Errorf("update product failed with status code: %d	", resp.StatusCode)
	}

	var productResp dtos.SingleProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&productResp); err != nil {
		return dtos.SingleProductResponse{}, fmt.Errorf("decode update product response: %w", err)
	}

	return productResp, nil
}
