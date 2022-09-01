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
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bss/v2/orders"
	"github.com/chnsz/golangsdk/openstack/bss/v2/resources"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

// GetRegion returns the region that was specified ina the resource. If a
// region was not set, the provider-level region is checked. The provider-level
// region can either be set by the region argument or by HW_REGION_NAME.
func GetRegion(d *schema.ResourceData, config *config.Config) string {
	if v, ok := d.GetOk("region"); ok {
		return v.(string)
	}

	return config.Region
}

// GetEnterpriseProjectID returns the enterprise_project_id that was specified in the resource.
// If it was not set, the provider-level value is checked. The provider-level value can
// either be set by the `enterprise_project_id` argument or by HW_ENTERPRISE_PROJECT_ID.
func GetEnterpriseProjectID(d *schema.ResourceData, config *config.Config) string {
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		return v.(string)
	}

	return config.EnterpriseProjectID
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
		return "", fmtp.Errorf("Unable to retrieve eips: %s ", err)
	}

	total := len(allEips)
	if total == 0 {
		return "", fmtp.Errorf("queried none results with %s", address)
	} else if total > 1 {
		return "", fmtp.Errorf("queried more results with %s", address)
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

	return fmtp.Errorf("%s: %s", msg, err)
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

	return fmtp.DiagErrorf("%s: %s", msg, err)
}

// UnsubscribePrePaidResource impl the action of unsubscribe resource
func UnsubscribePrePaidResource(d *schema.ResourceData, config *config.Config, resourceIDs []string) error {
	bssV2Client, err := config.BssV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud bss V2 client: %s", err)
	}

	unsubscribeOpts := orders.UnsubscribeOpts{
		ResourceIds:     resourceIDs,
		UnsubscribeType: 1,
	}
	_, err = orders.Unsubscribe(bssV2Client, unsubscribeOpts).Extract()
	return err
}

func CheckForRetryableError(err error) *resource.RetryError {
	switch errCode := err.(type) {
	case golangsdk.ErrDefault500:
		return resource.RetryableError(err)
	case golangsdk.ErrUnexpectedResponseCode:
		switch errCode.Actual {
		case 409, 503:
			return resource.RetryableError(err)
		default:
			return resource.NonRetryableError(err)
		}
	default:
		return resource.NonRetryableError(err)
	}
}

func WaitOrderComplete(ctx context.Context, d *schema.ResourceData, config *config.Config, orderNum string) error {
	bssV2Client, err := config.BssV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud bss V2 client: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"3", "6"}, // 3: Processing; 6: Pending payment.
		Target:       []string{"5"},      // 5: Completed.
		Refresh:      refreshOrderStatus(bssV2Client, orderNum),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("Error while waiting for the order (%s) to complete payment: %#v", orderNum, err)
	}
	return nil
}

func refreshOrderStatus(c *golangsdk.ServiceClient, orderNum string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := orders.Get(c, orderNum).Extract()
		if err != nil {
			return nil, "Error", err
		}
		return r, strconv.Itoa(r.OrderInfo.Status), nil
	}
}

func CaseInsensitiveFunc() schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return strings.EqualFold(old, new)
	}
}

// GetAutoPay is a method to return whether order is auto pay according to the user input.
// auto_pay parameter inputs and returns:
//
//	false: false
//	true, empty: true
//
// Before using this function, make sure the parameter behavior is auto pay (the default value is "true").
func GetAutoPay(d *schema.ResourceData) string {
	if d.Get("auto_pay").(string) == "false" {
		return "false"
	}
	return "true"
}

func UpdateAutoRenew(c *golangsdk.ServiceClient, enabled, resourceId string) error {
	if enabled == "true" {
		return resources.EnableAutoRenew(c, resourceId)
	}
	return resources.DisableAutoRenew(c, resourceId)
}
