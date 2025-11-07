package geminidb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	retryErrCodes = map[string]struct{}{
		"DBS.200019": {},
		"DBS.201015": {},
		"DBS.280005": {},
	}
)

// The RDS instance is limited to only one operation at a time.
// In addition to locking and waiting between multiple operations, a retry method is required to ensure that the
// request can be executed correctly.
func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code||errCode", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			return true, err
		}
	}
	if errCode, ok := err.(golangsdk.ErrDefault409); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("errCode||error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	if errCode, ok := err.(golangsdk.ErrDefault403); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("errCode||error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	if errCode, ok := err.(golangsdk.ErrDefault500); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("errCode||error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	// Operation execution failed due to some resource or server issues, no need to try again.
	return false, err
}

type updateInstanceFieldParams struct {
	httpUrl              string
	httpMethod           string
	pathParams           map[string]string
	updateBodyParams     interface{}
	isRetry              bool
	timeout              string
	checkJobExpression   string
	checkOrderExpression string
	bssClient            *golangsdk.ServiceClient
	isWaitInstanceReady  bool
}

func updateGeminiDbInstanceField(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	params updateInstanceFieldParams) (interface{}, error) {
	updatePath := client.Endpoint + params.httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	for pathParam, value := range params.pathParams {
		updatePath = strings.ReplaceAll(updatePath, fmt.Sprintf("{%s}", pathParam), value)
	}

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 201, 202, 204},
		JSONBody:         params.updateBodyParams,
	}

	var res interface{}
	var err error
	if params.isRetry {
		retryFunc := func() (interface{}, bool, error) {
			r, err := client.Request(params.httpMethod, updatePath, &updateOpt)
			retry, err := handleMultiOperationsError(err)
			return r, retry, err
		}
		res, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     geminiDbInstanceStatusRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(params.timeout),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
	} else {
		res, err = client.Request(params.httpMethod, updatePath, &updateOpt)
	}
	if err != nil {
		return nil, err
	}

	updateRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}
	if params.checkJobExpression == "" && params.checkOrderExpression == "" && !params.isWaitInstanceReady {
		return updateRespBody, nil
	}

	jobId := utils.PathSearch(params.checkJobExpression, updateRespBody, "").(string)
	orderId := utils.PathSearch(params.checkOrderExpression, updateRespBody, "").(string)
	err = checkJobAndOrderExpression(params, jobId, orderId)
	if err != nil {
		return nil, err
	}

	if jobId != "" {
		err = checkGeminiDbInstanceJobFinish(ctx, client, jobId, d.Timeout(params.timeout))
		if err != nil {
			return nil, err
		}
	}
	if orderId != "" {
		err = common.WaitOrderComplete(ctx, params.bssClient, orderId, d.Timeout(params.timeout))
		if err != nil {
			return nil, err
		}
	}
	if params.isWaitInstanceReady {
		err = waitForGeminiDBInstanceReady(ctx, d, client)
		if err != nil {
			return nil, err
		}
	}

	return updateRespBody, err
}

func checkJobAndOrderExpression(params updateInstanceFieldParams, jobId, orderId string) error {
	switch {
	case params.checkJobExpression != "" && params.checkOrderExpression != "":
		if jobId == "" && orderId == "" {
			return fmt.Errorf(" %s and %s is not found in the API response", params.checkJobExpression,
				params.checkOrderExpression)
		}
	case params.checkJobExpression != "" && params.checkOrderExpression == "":
		if jobId == "" {
			return fmt.Errorf(" %s is not found in the API response", params.checkJobExpression)
		}
	case params.checkJobExpression == "" && params.checkOrderExpression != "":
		if orderId == "" {
			return fmt.Errorf(" %s is not found in the API response", params.checkOrderExpression)
		}
	}
	return nil
}

func waitForGeminiDBInstanceReady(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Refresh:      geminiDbInstanceStatusRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        1 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for GeminiDB instance (%s) to ready: %s", d.Id(), err)
	}
	return nil
}

// The RDS instance can not be deleted or unsubscribe if another operation is being performed.
func handleDeletionError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	// unsubscribe fail
	var errCode400 golangsdk.ErrDefault400
	if errors.As(err, &errCode400) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode400.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		// CBC.99003651: Another operation is being performed.
		if errorCode == "CBC.99003651" {
			return true, err
		}
	}
	// delete fail
	var errCode409 golangsdk.ErrDefault409
	if errors.As(err, &errCode409) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode409.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok := retryErrCodes[errorCode.(string)]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	return false, err
}
