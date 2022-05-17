package utils

type DocsAPIError struct {
	Result interface{} `json:"result"`
	Code   int         `json:"code"`
	Error  string      `json:"error"`
}