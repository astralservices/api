package utils

type Response[T any] struct {
	Result T      `json:"result"`
	Error  string `json:"error"`
	Code   int    `json:"code"`
}
