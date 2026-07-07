package gaussdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	// Some error codes that need to be retried coming from https://support.huaweicloud.com/api-gaussdbformysql/ErrorCode.html
	retryErrCodes = map[string]struct{}{
		"DBS.200019":   {},
		"DBS.201014":   {},
		"DBS.201015":   {},
		"DBS.200047":   {},
		"DBS.201202":   {},
		"DBS.216011":   {},
		"DBS.05000084": {},
	}
)

func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault409); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code||errCode", apiError)
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
	delay                int
	checkJobExpression   string
	checkOrderExpression string
	bssClient            *golangsdk.ServiceClient
	isWaitInstanceReady  bool
}

func updateGaussDbInstanceField(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
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
			shouldRetry, err := handleMultiOperationsError(err)
			return r, shouldRetry, err
		}
		res, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     instanceStateRefreshFunc(client, d.Id()),
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
		err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId, params.delay, d.Timeout(params.timeout))
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
		err = waitForInstanceReady(ctx, d, client)
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

func waitForInstanceReady(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	stateConf := &retry.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Refresh:      instanceStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        1 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for GaussDB instance (%s) to ready: %s", d.Id(), err)
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
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request(params.httpMethod, getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(getResp)
}
