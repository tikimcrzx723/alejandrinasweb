package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tikimcrzx723/alejandrinasweb/internal/dtos"
)

func Login(ctx context.Context, baseURL string, req dtos.LoginRequest) (dtos.LoginResponse, error) {
	if strings.TrimSpace(baseURL) == "" {
		return dtos.LoginResponse{}, fmt.Errorf("baseURL is required")
	}

	url := strings.TrimRight(baseURL, "/") + "/auth/login"
	body, err := json.Marshal(req)
	if err != nil {
		return dtos.LoginResponse{}, fmt.Errorf("marshal login request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return dtos.LoginResponse{}, fmt.Errorf("create login request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return dtos.LoginResponse{}, fmt.Errorf("send login request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return dtos.LoginResponse{}, fmt.Errorf("read login response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return dtos.LoginResponse{}, fmt.Errorf("login failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var loginResp dtos.LoginResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return dtos.LoginResponse{}, fmt.Errorf("decode login response: %w", err)
	}

	return loginResp, nil
}

func Register(ctx context.Context, baseURL string, req dtos.RegisterRequest) (dtos.RegisterResponse, error) {
	if strings.TrimSpace(baseURL) == "" {
		return dtos.RegisterResponse{}, fmt.Errorf("baseURL is required")
	}

	url := strings.TrimRight(baseURL, "/") + "/auth/register"
	body, err := json.Marshal(req)
	if err != nil {
		return dtos.RegisterResponse{}, fmt.Errorf("marshal register request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return dtos.RegisterResponse{}, fmt.Errorf("create register request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return dtos.RegisterResponse{}, fmt.Errorf("send register request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return dtos.RegisterResponse{}, fmt.Errorf("read register response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return dtos.RegisterResponse{}, fmt.Errorf("register failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var registerResp dtos.RegisterResponse
	if err := json.Unmarshal(respBody, &registerResp); err != nil {
		return dtos.RegisterResponse{}, fmt.Errorf("decode register response: %w", err)
	}

	return registerResp, nil
}
