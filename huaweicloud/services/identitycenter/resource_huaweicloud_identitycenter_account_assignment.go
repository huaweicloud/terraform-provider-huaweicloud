// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IdentityCenter
// ---------------------------------------------------------------

package identitycenter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter GET /v1/instances/{instance_id}/account-assignments
// @API IdentityCenter POST /v1/instances/{instance_id}/account-assignments/create
// @API IdentityCenter GET /v1/instances/{instance_id}/account-assignments/creation-status/{request_id}
// @API IdentityCenter POST /v1/instances/{instance_id}/account-assignments/delete
// @API IdentityCenter GET /v1/instances/{instance_id}/account-assignments/deletion-status/{request_id}
func ResourceIdentityCenterAccountAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountAssignmentCreate,
		ReadContext:   resourceAccountAssignmentRead,
		DeleteContext: resourceAccountAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterAccountAssignmentImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the Identity Center instance.`,
			},
			"permission_set_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the permission set.`,
			},
			"principal_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the user or user group that belongs to IAM Identity Center.`,
			},
			"principal_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of the user or user group.`,
			},
			"target_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the target to be bound.`,
			},
			"target_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of the target to be bound.`,
			},
		},
	}
}

func resourceAccountAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAccountAssignment: create Identity Center account assignment
	var (
		createAccountAssignmentHttpUrl = "v1/instances/{instance_id}/account-assignments/create"
		getRequestStatusHttpUrl        = "v1/instances/{instance_id}/account-assignments/creation-status/{request_id}"
		createAccountAssignmentProduct = "identitycenter"
	)
	createAccountAssignmentClient, err := cfg.NewServiceClient(createAccountAssignmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createAccountAssignmentPath := createAccountAssignmentClient.Endpoint + createAccountAssignmentHttpUrl
	createAccountAssignmentPath = strings.ReplaceAll(createAccountAssignmentPath, "{instance_id}", instanceID)

	createAccountAssignmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createAccountAssignmentOpt.JSONBody = utils.RemoveNil(buildAccountAssignmentBodyParams(d))
	createAccountAssignmentResp, err := createAccountAssignmentClient.Request("POST",
		createAccountAssignmentPath, &createAccountAssignmentOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center account assignment: %s", err)
	}

	permissionSetID := d.Get("permission_set_id")
	accountID := d.Get("target_id")
	principalID := d.Get("principal_id")
	d.SetId(fmt.Sprintf("%v/%v/%v", permissionSetID, accountID, principalID))

	createAccountAssignmentRespBody, err := utils.FlattenResponse(createAccountAssignmentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	requestID := utils.PathSearch("account_assignment_creation_status.request_id", createAccountAssignmentRespBody, "").(string)
	if requestID == "" {
		return diag.Errorf("unable to find the request ID of the Identity Center account assignment from the API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"IN_PROGRESS"},
		Target:  []string{"SUCCEEDED"},
		Refresh: identityCenterStatusRefreshFunc(requestID, instanceID, getRequestStatusHttpUrl,
			"account_assignment_creation_status.status", createAccountAssignmentClient),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for request (%s) to complete: %s", requestID, err)
	}

	return resourceAccountAssignmentRead(ctx, d, meta)
}

func buildAccountAssignmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"permission_set_id": utils.ValueIgnoreEmpty(d.Get("permission_set_id")),
		"principal_id":      utils.ValueIgnoreEmpty(d.Get("principal_id")),
		"principal_type":    utils.ValueIgnoreEmpty(d.Get("principal_type")),
		"target_id":         utils.ValueIgnoreEmpty(d.Get("target_id")),
		"target_type":       utils.ValueIgnoreEmpty(d.Get("target_type")),
	}
	return bodyParams
}

func resourceAccountAssignmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	//  getAccountAssignment: Query Identity Center account assignment
	var (
		getAccountAssignmentHttpUrl = "v1/instances/{instance_id}/account-assignments"
		getAccountAssignmentProduct = "identitycenter"
	)
	getAccountAssignmentClient, err := cfg.NewServiceClient(getAccountAssignmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	getAccountAssignmentBasePath := getAccountAssignmentClient.Endpoint + getAccountAssignmentHttpUrl
	getAccountAssignmentBasePath = strings.ReplaceAll(getAccountAssignmentBasePath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))

	getAccountAssignmentPath := getAccountAssignmentBasePath + buildGetAccountAssignmentQueryParams(d, "")

	getAccountAssignmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var accountPermission interface{}
	principalId := d.Get("principal_id")
getAccountPermissionsLoop:
	for {
		getAccountAssignmentResp, err := getAccountAssignmentClient.Request("GET", getAccountAssignmentPath,
			&getAccountAssignmentOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving Identity Center account assignment")
		}

		getAccountAssignmentRespBody, err := utils.FlattenResponse(getAccountAssignmentResp)
		if err != nil {
			return diag.FromErr(err)
		}

		accountPermissions := utils.PathSearch("account_assignments", getAccountAssignmentRespBody, nil)

		if accountPermissions == nil {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
		}

		for _, v := range accountPermissions.([]interface{}) {
			if principalId == utils.PathSearch("principal_id", v, "") {
				accountPermission = v
				break getAccountPermissionsLoop
			}
		}
		marker := utils.PathSearch("page_info.next_marker", getAccountAssignmentRespBody, nil)
		if marker == nil {
			break
		}
		getAccountAssignmentPath = getAccountAssignmentBasePath + buildGetAccountAssignmentQueryParams(d, marker.(string))
	}
	if accountPermission == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("permission_set_id", utils.PathSearch("permission_set_id", accountPermission, nil)),
		d.Set("principal_id", utils.PathSearch("principal_id", accountPermission, nil)),
		d.Set("principal_type", utils.PathSearch("principal_type", accountPermission, nil)),
		d.Set("target_id", utils.PathSearch("account_id", accountPermission, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAccountAssignmentQueryParams(d *schema.ResourceData, marker string) string {
	res := "?limit=100"
	if v, ok := d.GetOk("target_id"); ok {
		res = fmt.Sprintf("%s&account_id=%v", res, v)
	}

	if v, ok := d.GetOk("permission_set_id"); ok {
		res = fmt.Sprintf("%s&permission_set_id=%v", res, v)
	}

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func resourceAccountAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAccountAssignment: Delete Identity Center account assignment
	var (
		deleteAccountAssignmentHttpUrl = "v1/instances/{instance_id}/account-assignments/delete"
		getRequestStatusHttpUrl        = "v1/instances/{instance_id}/account-assignments/deletion-status/{request_id}"
		deleteAccountAssignmentProduct = "identitycenter"
	)
	deleteAccountAssignmentClient, err := cfg.NewServiceClient(deleteAccountAssignmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deleteAccountAssignmentPath := deleteAccountAssignmentClient.Endpoint + deleteAccountAssignmentHttpUrl
	deleteAccountAssignmentPath = strings.ReplaceAll(deleteAccountAssignmentPath, "{instance_id}", instanceID)

	deleteAccountAssignmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteAccountAssignmentOpt.JSONBody = utils.RemoveNil(buildAccountAssignmentBodyParams(d))
	deleteAccountAssignmentResp, err := deleteAccountAssignmentClient.Request("POST",
		deleteAccountAssignmentPath, &deleteAccountAssignmentOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center account assignment: %s", err)
	}

	deleteAccountAssignmentRespBody, err := utils.FlattenResponse(deleteAccountAssignmentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	requestID := utils.PathSearch("account_assignment_deletion_status.request_id",
		deleteAccountAssignmentRespBody, "").(string)
	if requestID == "" {
		return diag.Errorf("unable to find the request ID of the Identity Center account assignment from the API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"IN_PROGRESS"},
		Target:  []string{"SUCCEEDED"},
		Refresh: identityCenterStatusRefreshFunc(requestID, instanceID, getRequestStatusHttpUrl,
			"account_assignment_deletion_status.status", deleteAccountAssignmentClient),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for request (%s) to complete: %s", requestID, err)
	}

	return nil
}

func identityCenterStatusRefreshFunc(requestID, instanceID, getRequestStatusHttpUrl, searchExpression string,
	client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRequestStatusPath := client.Endpoint + getRequestStatusHttpUrl
		getRequestStatusPath = strings.ReplaceAll(getRequestStatusPath, "{instance_id}", instanceID)
		getRequestStatusPath = strings.ReplaceAll(getRequestStatusPath, "{request_id}", requestID)

		getRequestStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getRequestStatusResp, err := client.Request("GET", getRequestStatusPath, &getRequestStatusOpt)
		if err != nil {
			return nil, "", err
		}
		getRequestStatusRespBody, err := utils.FlattenResponse(getRequestStatusResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch(searchExpression, getRequestStatusRespBody, "").(string)
		if status == "SUCCEEDED" || status == "FAILED" {
			return getRequestStatusRespBody, status, nil
		}
		return getRequestStatusRespBody, "IN_PROGRESS", nil
	}
}

func resourceIdentityCenterAccountAssignmentImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid id format, must be " +
			"<instance_id>/<permission_set_id>/<target_id>/<principal_id>")
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", parts[1], parts[2], parts[3]))
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("permission_set_id", parts[1]),
		d.Set("target_id", parts[2]),
		d.Set("principal_id", parts[3]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
