package dcs

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
	operateErrorCode = map[string]struct{}{
		// current state not support
		"DCS.4026": {},
		// instance status is not running
		"DCS.4049": {},
		// backup
		"DCS.4096": {},
		// restore
		"DCS.4097": {},
		// restart
		"DCS.4111": {},
		// resize
		"DCS.4113": {},
		// change config
		"DCS.4114": {},
		// change password
		"DCS.4115": {},
		// upgrade
		"DCS.4116": {},
		// rollback
		"DCS.4117": {},
		// create
		"DCS.4118": {},
		// freeze
		"DCS.4120": {},
		// reset password
		"DCS.4121": {},
		// creating/restarting
		"DCS.4975": {},
	}
)

func handleOperationError(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, jsonErr
		}
		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}
		// CBC.99003651: Another operation is being performed.
		if _, ok = operateErrorCode[errorCode.(string)]; ok || errorCode == "CBC.99003651" {
			return true, err
		}
	}
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

func updateDcsInstanceField(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	params updateInstanceFieldParams) (interface{}, error) {
	updatePath := client.Endpoint + params.httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	for pathParam, value := range params.pathParams {
		updatePath = strings.ReplaceAll(updatePath, fmt.Sprintf("{%s}", pathParam), value)
	}

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         params.updateBodyParams,
		OkCodes:          []int{200, 204},
	}

	var res interface{}
	var err error
	if params.isRetry {
		retryFunc := func() (interface{}, bool, error) {
			r, err := client.Request(params.httpMethod, updatePath, &updateOpt)
			retry, err := handleOperationError(err)
			return r, retry, err
		}
		res, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     refreshDcsInstanceState(client, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			WaitPending:  []string{"PENDING"},
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
	err = checkJobAndOrder(params, jobId, orderId)
	if err != nil {
		return nil, err
	}

	if jobId != "" {
		err = checkDcsInstanceJobFinish(ctx, client, jobId, d.Timeout(params.timeout))
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
		err = waitForDcsInstanceRunning(ctx, client, d.Id(), d.Timeout(params.timeout))
		if err != nil {
			return nil, err
		}
	}

	return updateRespBody, err
}

func checkJobAndOrder(params updateInstanceFieldParams, jobId, orderId string) error {
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

func checkDcsInstanceJobFinish(ctx context.Context, client *golangsdk.ServiceClient, jobId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      dcsInstanceJobRefreshFunc(client, jobId),
		Timeout:      timeout,
		Delay:        2 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for DCS instance job (%s) to be completed: %s ", jobId, err)
	}
	return nil
}

func dcsInstanceJobRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		httpUrl := "v2/{project_id}/jobs/{job_id}"
		getPath := client.Endpoint + httpUrl
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{job_id}", jobId)

		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, "ERROR", err
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", getRespBody, "").(string)
		if status == "SUCCESS" || status == "FAIL" {
			return getRespBody, status, nil
		}

		return getRespBody, "PENDING", nil
	}
}

func waitForDcsInstanceRunning(ctx context.Context, c *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"PENDING"},
		Target:                    []string{"RUNNING"},
		Refresh:                   refreshDcsInstanceState(c, id),
		Timeout:                   timeout,
		Delay:                     10 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting instance(%s) to ready: %s", id, err)
	}
	return nil
}

func refreshDcsInstanceState(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := getDcsInstanceByID(client, id)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "", "DELETED", nil
			}
			return nil, "ERROR", err
		}
		status := utils.PathSearch("status", instance, "").(string)

		failStatus := []string{"CREATEFAILED", "ERROR", "FROZEN"}
		if utils.StrSliceContains(failStatus, status) {
			return instance, status, fmt.Errorf("unexpect status: %s", status)
		}
		if status == "RUNNING" {
			return instance, status, nil
		}
		return instance, "PENDING", nil
	}
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
