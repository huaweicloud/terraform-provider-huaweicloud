/* Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights resvered. */
/*
The common package defines some common functions, which are mainly used for the functions of the following services.

The difference between common package and utils:
1. Common functions under common are related to the project, and common functions are placed here.
2. Utils are some stored tool functions, which are not related to the project.
   Such as: date conversion, type conversion.
*/
package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bss/v2/orders"
	"github.com/chnsz/golangsdk/openstack/bss/v2/resources"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// ErrorResp is the response when API failed
type ErrorResp struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func ParseErrorMsg(body []byte) (ErrorResp, error) {
	resp := ErrorResp{}
	err := json.Unmarshal(body, &resp)
	return resp, err
}

// GetRegion returns the region that was specified ina the resource. If a
// region was not set, the provider-level region is checked. The provider-level
// region can either be set by the region argument or by HW_REGION_NAME.
func GetRegion(d *schema.ResourceData, config *config.Config) string {
	if v, ok := d.GetOk("region"); ok {
		return v.(string)
	}

	return config.Region
}

// GetEipIDbyAddress returns the EIP ID of address when success.
func GetEipIDbyAddress(client *golangsdk.ServiceClient, address, epsID string) (string, error) {
	listOpts := &eips.ListOpts{
		PublicIp:            []string{address},
		EnterpriseProjectId: epsID,
	}
	pages, err := eips.List(client, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve eips: %s ", err)
	}

	total := len(allEips)
	if total == 0 {
		return "", fmt.Errorf("queried none results with %s", address)
	} else if total > 1 {
		return "", fmt.Errorf("queried more results with %s", address)
	}

	return allEips[0].ID, nil
}

// CheckDeleted checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeleted(d *schema.ResourceData, err error, msg string) error {
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("%s: %s", msg, err)
}

// CheckDeletedDiag checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeletedDiag(d *schema.ResourceData, err error, msg string) diag.Diagnostics {
	var statusCode int

	// check if the error is raised by **golangsdk**
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		statusCode = http.StatusNotFound
	} else if responseErr, ok := err.(*sdkerr.ServiceResponseError); ok {
		// check if the error is raised by **huaweicloud-sdk-go-v3**
		statusCode = responseErr.StatusCode
	}

	if statusCode == http.StatusNotFound {
		resourceID := d.Id()
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("the resource %s is gone and will be removed in Terraform state.", resourceID),
			},
		}
	}

	return diag.Errorf("%s: %s", msg, err)
}

// UnsubscribePrePaidResource impl the action of unsubscribe resource
func UnsubscribePrePaidResource(d *schema.ResourceData, config *config.Config, resourceIDs []string) error {
	bssV2Client, err := config.BssV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating bss V2 client: %s", err)
	}

	unsubscribeOpts := orders.UnsubscribeOpts{
		ResourceIds:     resourceIDs,
		UnsubscribeType: 1,
	}
	_, err = orders.Unsubscribe(bssV2Client, unsubscribeOpts).Extract()
	return err
}

func CheckForRetryableError(err error) *retry.RetryError {
	switch errCode := err.(type) {
	case golangsdk.ErrDefault500:
		return retry.RetryableError(err)
	case golangsdk.ErrUnexpectedResponseCode:
		switch errCode.Actual {
		case 409, 503:
			return retry.RetryableError(err)
		default:
			return retry.NonRetryableError(err)
		}
	default:
		return retry.NonRetryableError(err)
	}
}

func WaitOrderComplete(ctx context.Context, client *golangsdk.ServiceClient, orderId string,
	timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"2", "3", "6"}, // 2: Pending refund 3: Processing; 6: Pending payment.
		Target:       []string{"5"},           // 5: Completed.
		Refresh:      refreshOrderStatusFunc(client, orderId),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the order (%s) to complete payment: %s", orderId, err)
	}
	return nil
}

func refreshOrderStatusFunc(client *golangsdk.ServiceClient, orderId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := orders.Get(client, orderId).Extract()
		if err != nil {
			return nil, "Error", err
		}
		return r, strconv.Itoa(r.OrderInfo.Status), nil
	}
}

// WaitOrderResourceComplete is the method to wait for the resource to be generated.
// Notes: Note that this method needs to be used in conjunction with method "WaitOrderComplete", because the ID of some
// resources may not be generated when the order is not completed.
func WaitOrderResourceComplete(ctx context.Context, client *golangsdk.ServiceClient, orderId string,
	timeout time.Duration) (string, error) {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DONE"},
		Refresh:      refreshOrderResourceStatusFunc(client, orderId),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	res, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", fmt.Errorf("error while waiting for the order (%s) to complete: %s", orderId, err)
	}

	r := res.(resources.Resource)
	return r.ResourceId, nil
}

func refreshOrderResourceStatusFunc(client *golangsdk.ServiceClient, orderId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		if strings.TrimSpace(orderId) == "" {
			return nil, "ERROR", fmt.Errorf("order id is empty")
		}
		listOpts := resources.ListOpts{
			OrderId:          orderId,
			OnlyMainResource: 1,
		}
		resp, err := resources.List(client, listOpts)
		if err != nil || resp == nil {
			return nil, "ERROR", fmt.Errorf("error waiting for the order (%s) to complete: %s", orderId, err)
		}
		if resp.TotalCount < 1 {
			return nil, "PENDING", nil
		}
		return resp.Resources[0], "DONE", nil
	}
}

// WaitOrderAllResourceComplete is the method to wait for the non-main resource to be generated.
// Notes: Note that this method needs to be used in conjunction with method "WaitOrderComplete", because the ID of some
// resources may not be generated when the order is not completed.
func WaitOrderAllResourceComplete(ctx context.Context, client *golangsdk.ServiceClient, orderId string,
	timeout time.Duration) (string, error) {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DONE"},
		Refresh:      refreshOrderAllResourceStatusFunc(client, orderId),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	res, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", fmt.Errorf("error while waiting for the order (%s) to complete: %s", orderId, err)
	}

	r := res.(resources.Resource)
	return r.ResourceId, nil
}

func refreshOrderAllResourceStatusFunc(client *golangsdk.ServiceClient, orderId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		if strings.TrimSpace(orderId) == "" {
			return nil, "ERROR", fmt.Errorf("order id is empty")
		}
		listOpts := resources.ListOpts{
			OrderId: orderId,
		}
		resp, err := resources.List(client, listOpts)
		if err != nil || resp == nil {
			return nil, "ERROR", fmt.Errorf("error waiting for the order (%s) to complete: %s", orderId, err)
		}
		if resp.TotalCount < 1 {
			return nil, "PENDING", nil
		}
		return resp.Resources[0], "DONE", nil
	}
}

func CaseInsensitiveFunc() schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return strings.EqualFold(old, new)
	}
}

func HasFilledOpt(d *schema.ResourceData, param string) bool {
	_, b := d.GetOk(param)
	return b
}

// RetryFunc is the function retried until it succeeds.
// The first return parameter is the result of the retry func.
// The second return parameter indicates whether a retry is required.
// The last return parameter is the error of the func.
type RetryFunc func() (res interface{}, retry bool, err error)

type RetryContextWithWaitForStateParam struct {
	Ctx context.Context
	// The func that need to be retried
	RetryFunc RetryFunc
	// The wait func when the retry which returned by the retry func is true
	WaitFunc retry.StateRefreshFunc
	// The target of the wait func
	WaitTarget []string
	// The pending of the wait func
	WaitPending []string
	// The timeout of the retry func and wait func
	Timeout time.Duration
	// The delay timeout of the retry func and wait func
	DelayTimeout time.Duration
	// The poll interval of the retry func and wait func
	PollInterval time.Duration
}

// RetryContextWithWaitForState The RetryFunc will be called first
// if the error of the return is nil, the retry will be ended and the res of the return will be returned
// if the retry of the return is true, the RetryFunc will be retried, and the WaitFunc will be called if it is not nil
// if the retry of the return is false, the retry will be ended and the error of the retry func will be returned
func RetryContextWithWaitForState(param *RetryContextWithWaitForStateParam) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"retryable"},
		Target:       []string{"success"},
		Timeout:      param.Timeout,
		PollInterval: param.PollInterval,
		Refresh: func() (interface{}, string, error) {
			res, shouldRetry, err := param.RetryFunc()
			if err == nil {
				if res != nil {
					return res, "success", nil
				}
				// If we didn't find the resource, convert it to "", otherwise,
				// it will report an error in WaitForStateContext.
				return "", "success", nil
			}

			if !shouldRetry {
				return nil, "quit", err
			}

			if param.WaitFunc != nil {
				stateConf := &retry.StateChangeConf{
					Target:       param.WaitTarget,
					Pending:      param.WaitPending,
					Refresh:      param.WaitFunc,
					Timeout:      param.Timeout,
					Delay:        param.DelayTimeout,
					PollInterval: param.PollInterval,
				}
				if _, err := stateConf.WaitForStateContext(param.Ctx); err != nil {
					return nil, "quit", err
				}
			}
			return "", "retryable", nil
		},
	}

	return stateConf.WaitForStateContext(param.Ctx)
}

// GetEipsbyAddresses returns the EIPs of addresses when success.
func GetEipsbyAddresses(client *golangsdk.ServiceClient, addresses []string, epsID string) ([]eips.PublicIp, error) {
	listOpts := &eips.ListOpts{
		PublicIp:            addresses,
		EnterpriseProjectId: epsID,
	}
	pages, err := eips.List(client, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve eips: %s ", err)
	}
	return allEips, nil
}
