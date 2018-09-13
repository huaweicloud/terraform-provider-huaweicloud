package imagedata

import (
	"fmt"
	"io"

	"github.com/huaweicloud/golangsdk"
)

// UploadResult is the result of an upload image operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type UploadResult struct {
	golangsdk.ErrResult
}

// StageResult is the result of a stage image operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type StageResult struct {
	golangsdk.ErrResult
}

// DownloadResult is the result of a download image operation. Call its Extract
// method to gain access to the image data.
type DownloadResult struct {
	golangsdk.Result
}

// Extract builds images model from io.Reader
func (r DownloadResult) Extract() (io.Reader, error) {
	if r, ok := r.Body.(io.Reader); ok {
		return r, nil
	}
	return nil, fmt.Errorf("Expected io.Reader but got: %T(%#v)", r.Body, r.Body)
}
