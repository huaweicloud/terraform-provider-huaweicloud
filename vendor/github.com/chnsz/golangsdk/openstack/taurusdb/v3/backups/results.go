package backups

import (
	"github.com/chnsz/golangsdk"
)

// UpdateResult represents the result of a update operation.
type UpdateResult struct {
	golangsdk.ErrResult
}

type UpdateEncryptionResponse struct {
	EncryptionStatus string `json:"encryption_status"`
}

type UpdateEncryptionResult struct {
	golangsdk.Result
}

func (r UpdateEncryptionResult) Extract() (*UpdateEncryptionResponse, error) {
	var updateEncryptionResponse UpdateEncryptionResponse
	err := r.ExtractInto(&updateEncryptionResponse)
	return &updateEncryptionResponse, err
}

type GetEncryptionResponse struct {
	EncryptionStatus string `json:"encryption_status"`
}

type GetEncryptionResult struct {
	golangsdk.Result
}

func (r GetEncryptionResult) Extract() (*GetEncryptionResponse, error) {
	var res GetEncryptionResponse
	err := r.ExtractInto(&res)
	return &res, err
}
