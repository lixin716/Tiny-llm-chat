package handlers

import (
	"encoding/json"
	"net/http"
)

// 响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 成功响应
func SuccessResponse(w http.ResponseWriter, data interface{}) {
	resp := Response{
		Code:    200,
		Message: "成功",
		Data:    data,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// 错误响应
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	resp := Response{
		Code:    statusCode,
		Message: message,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

// JSON请求解析
func ParseJSON(r *http.Request, dest interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(dest)
} 