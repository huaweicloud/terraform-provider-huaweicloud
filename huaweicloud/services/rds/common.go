package rds

import (
	"context"
	"encoding/json"
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
	// Some error codes that need to be retried coming from https://console-intl.huaweicloud.com/apiexplorer/#/errorcenter/RDS.
	retryErrCodes = map[string]struct{}{
		"DBS.201202":   {},
		"DBS.200011":   {},
		"DBS.200018":   {},
		"DBS.200019":   {},
		"DBS.200047":   {},
		"DBS.200076":   {},
		"DBS.200611":   {},
		"DBS.200080":   {},
		"DBS.200463":   {}, // create replica instance
		"DBS.201015":   {},
		"DBS.201206":   {},
		"DBS.212033":   {}, // http response code is 403
		"DBS.280011":   {},
		"DBS.280343":   {},
		"DBS.280816":   {},
		"DBS.01010337": {},
		"DBS.01280030": {}, // instance status is illegal
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
	if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok && errCode.Actual == 409 {
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

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		// if the error code is RDS.0005, it indicates that the SSL is changed, and the db is rebooted
		if errorCode.(string) == "RDS.0005" {
			return true, err
		}
	}
	// Operation execution failed due to some resource or server issues, no need to try again.
	return false, err
}

// The RDS instance can not be deleted or unsubscribe if another operation is being performed.
func handleDeletionError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	// unsubscribe fail
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
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
	if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok && errCode.Actual == 409 {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	return false, err
}

// The RDS cross region backup strategy can not be updated if another operation is being performed.
func handleCrossRegionBackupStrategyError(err error) (bool, error) {
	if errCode, ok := err.(golangsdk.ErrDefault500); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		// DBS.280228: Another operation is being performed.
		if errorCode == "DBS.280228" {
			return true, err
		}
	}
	return false, err
}

func handleApiNotExistsError(err error) bool {
	if err == nil {
		return false
	}
	if errCode, ok := err.(golangsdk.ErrDefault404); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false
		}
		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false
		}
		if errorCode.(string) == "APIGW.0101" {
			return true
		}
	}
	return false
}

func handleTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	// if the http status code is 500 and the error code is DBS.111205, it indicates timeout for the service
	// error should be ignored, just wait for success
	if errCode, ok := err.(golangsdk.ErrDefault500); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false
		}
		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false
		}
		if errorCode.(string) == "DBS.111205" {
			return true
		}
	}
	return false
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
}

func updateRdsInstanceField(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	params updateInstanceFieldParams) error {
	updatePath := client.Endpoint + params.httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	for pathParam, value := range params.pathParams {
		updatePath = strings.ReplaceAll(updatePath, fmt.Sprintf("{%s}", pathParam), value)
	}

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
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
			WaitFunc:     rdsInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(params.timeout),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
	} else {
		res, err = client.Request(params.httpMethod, updatePath, &updateOpt)
	}
	if err != nil {
		return err
	}

	if params.checkJobExpression == "" && params.checkOrderExpression == "" {
		return nil
	}

	updateRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return err
	}

	jobId := utils.PathSearch(params.checkJobExpression, updateRespBody, "").(string)
	orderId := utils.PathSearch(params.checkOrderExpression, updateRespBody, "").(string)
	if jobId == "" && orderId == "" {
		if params.checkJobExpression != "" && params.checkOrderExpression != "" {
			return fmt.Errorf(" %s and %s is not found in the API response", params.checkJobExpression, params.checkOrderExpression)
		}
		if params.checkJobExpression != "" {
			return fmt.Errorf(" %s is not found in the API response", params.checkJobExpression)
		}
		return fmt.Errorf(" %s is not found in the API response", params.checkOrderExpression)
	}

	if jobId != "" {
		err = checkRDSInstanceJobFinish(client, jobId, d.Timeout(params.timeout))
		if err != nil {
			return err
		}
	}
	if orderId != "" {
		// wait for order success
		err = common.WaitOrderComplete(ctx, params.bssClient, orderId, d.Timeout(params.timeout))
		if err != nil {
			return err
		}

		stateConf := &resource.StateChangeConf{
			Target:       []string{"ACTIVE"},
			Refresh:      rdsInstanceStateRefreshFunc(client, d.Id()),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        1 * time.Second,
			PollInterval: 10 * time.Second,
		}
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("error waiting for instance (%s) to ready: %s", d.Id(), err)
		}
	}

	return nil
}

type getInstanceFieldParams struct {
	httpUrl    string
	httpMethod string
	pathParams map[string]string
}

func getInstanceField(client *golangsdk.ServiceClient, params getInstanceFieldParams) (interface{}, error) {
	getPath := client.Endpoint + params.httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	for pathParam, value := range params.pathParams {
		getPath = strings.ReplaceAll(getPath, fmt.Sprintf("{%s}", pathParam), value)
	}

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request(params.httpMethod, getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(getResp)
}
