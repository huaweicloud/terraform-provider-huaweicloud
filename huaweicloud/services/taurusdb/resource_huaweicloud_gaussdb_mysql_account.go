// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GaussDB
// ---------------------------------------------------------------

package taurusdb

import (
	"context"
	"fmt"
	"net/http"
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

// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/db-users
// @API GaussDBforMySQL GET /v3/{project_id}/jobs
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/db-users/comment
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/db-users/password
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/db-users
// @API GaussDBforMySQL DELETE /v3/{project_id}/instances/{instance_id}/db-users
func ResourceGaussDBAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBAccountCreate,
		UpdateContext: resourceGaussDBAccountUpdate,
		ReadContext:   resourceGaussDBAccountRead,
		DeleteContext: resourceGaussDBAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the GaussDB MySQL instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the database username.`,
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `Specifies the password of the database user.`,
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the host IP address.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the database user remarks.`,
			},
		},
	}
}

func resourceGaussDBAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createGaussDBAccount: create a GaussDB MySQL account
	var (
		createGaussDBAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db-users"
		createGaussDBAccountProduct = "gaussdb"
	)
	createGaussDBAccountClient, err := cfg.NewServiceClient(createGaussDBAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createGaussDBAccountPath := createGaussDBAccountClient.Endpoint + createGaussDBAccountHttpUrl
	createGaussDBAccountPath = strings.ReplaceAll(createGaussDBAccountPath, "{project_id}",
		createGaussDBAccountClient.ProjectID)
	createGaussDBAccountPath = strings.ReplaceAll(createGaussDBAccountPath, "{instance_id}", instanceID)

	createGaussDBAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	host := "%"
	if h := d.Get("host").(string); h != "" {
		host = h
	}
	createGaussDBAccountOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBAccountBodyParams(d, host))

	var createGaussDBAccountResp *http.Response
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		createGaussDBAccountResp, err = createGaussDBAccountClient.Request("POST", createGaussDBAccountPath,
			&createGaussDBAccountOpt)
		isRetry, err := handleGaussDBMysqlOperationError(err)
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("error creating GaussDB MySQL account: %s", err)
	}

	createGaussDBAccountRespBody, err := utils.FlattenResponse(createGaussDBAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createGaussDBAccountRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID of the GaussDB MySQL account from the API response")
	}
	err = waitForJobComplete(ctx, createGaussDBAccountClient, d.Timeout(schema.TimeoutCreate), instanceID, jobId)
	if err != nil {
		return diag.FromErr(err)
	}

	accountName := d.Get("name").(string)
	d.SetId(instanceID + "/" + accountName + "/" + host)

	return resourceGaussDBAccountRead(ctx, d, meta)
}

func buildCreateGaussDBAccountBodyParams(d *schema.ResourceData, host string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":     utils.ValueIgnoreEmpty(d.Get("name")),
		"comment":  utils.ValueIgnoreEmpty(d.Get("description")),
		"password": utils.ValueIgnoreEmpty(d.Get("password")),
		"hosts":    []interface{}{host},
	}
	param := map[string]interface{}{
		"users": []interface{}{bodyParams},
	}
	return param
}

func resourceGaussDBAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getGaussDBAccount: Query the GaussDB MySQL account
	var (
		getGaussDBAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db-users"
		getGaussDBAccountProduct = "gaussdb"
	)
	getGaussDBAccountClient, err := cfg.NewServiceClient(getGaussDBAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return diag.Errorf("invalid id format, must be <instance_id>/<name>/<host>")
	}
	instanceID := parts[0]
	accountName := parts[1]
	host := parts[2]

	getGaussDBAccountBasePath := getGaussDBAccountClient.Endpoint + getGaussDBAccountHttpUrl
	getGaussDBAccountBasePath = strings.ReplaceAll(getGaussDBAccountBasePath, "{project_id}",
		getGaussDBAccountClient.ProjectID)
	getGaussDBAccountBasePath = strings.ReplaceAll(getGaussDBAccountBasePath, "{instance_id}", instanceID)

	getGaussDBAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	var currentTotal int
	var account interface{}
	getGaussDBAccountPath := getGaussDBAccountBasePath + buildGaussDBMysqlQueryParams(currentTotal)

	for {
		getGaussDBAccountResp, err := getGaussDBAccountClient.Request("GET", getGaussDBAccountPath,
			&getGaussDBAccountOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB MySQL account")
		}

		getGaussDBAccountRespBody, err := utils.FlattenResponse(getGaussDBAccountResp)
		if err != nil {
			return diag.FromErr(err)
		}
		res, pageNum := flattenGetGaussDBAccountResponseBody(getGaussDBAccountRespBody, accountName, host)
		if res != nil {
			account = res
			break
		}
		total := utils.PathSearch("total_count", getGaussDBAccountRespBody, float64(0)).(float64)
		currentTotal += pageNum
		if currentTotal == int(total) {
			break
		}
		getGaussDBAccountPath = getGaussDBAccountBasePath + buildGaussDBMysqlQueryParams(currentTotal)
	}
	if account == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("name", utils.PathSearch("name", account, nil)),
		d.Set("description", utils.PathSearch("comment", account, nil)),
		d.Set("host", utils.PathSearch("host", account, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetGaussDBAccountResponseBody(resp interface{}, accountName, address string) (interface{}, int) {
	if resp == nil {
		return nil, 0
	}
	curJson := utils.PathSearch("users", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	for _, v := range curArray {
		name := utils.PathSearch("name", v, "").(string)
		host := utils.PathSearch("host", v, "").(string)
		if accountName == name && address == host {
			return v, len(curArray)
		}
	}
	return nil, len(curArray)
}

func resourceGaussDBAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateGaussDBAccount: update the GaussDB MySQL account
	var (
		updateGaussDBAccountProduct = "gaussdb"
	)
	updateGaussDBAccountClient, err := cfg.NewServiceClient(updateGaussDBAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	if d.HasChange("description") {
		if err = updateGaussDBAccountDescription(ctx, d, updateGaussDBAccountClient, instanceID); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("password") {
		if err = updateGaussDBAccountPassword(ctx, d, updateGaussDBAccountClient, instanceID); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGaussDBAccountRead(ctx, d, meta)
}

func updateGaussDBAccountDescription(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) error {
	// updateGaussDBAccountDescription: update the GaussDB MySQL account description
	var (
		updateGaussDBAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db-users/comment"
	)

	updateGaussDBAccountPath := client.Endpoint + updateGaussDBAccountHttpUrl
	updateGaussDBAccountPath = strings.ReplaceAll(updateGaussDBAccountPath, "{project_id}", client.ProjectID)
	updateGaussDBAccountPath = strings.ReplaceAll(updateGaussDBAccountPath, "{instance_id}", instanceID)

	updateGaussDBAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	updateGaussDBAccountOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBAccountDescriptionBodyParams(d))
	updateGaussDBAccountResp, err := client.Request("PUT", updateGaussDBAccountPath, &updateGaussDBAccountOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL account description: %s", err)
	}

	updateGaussDBAccountRespBody, err := utils.FlattenResponse(updateGaussDBAccountResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateGaussDBAccountRespBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("unable to find the job ID of the GaussDB MySQL account from the API response while updating the description")
	}

	return waitForJobComplete(ctx, client, d.Timeout(schema.TimeoutUpdate), instanceID, jobId)
}

func updateGaussDBAccountPassword(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) error {
	// updateGaussDBAccountPassword: update the GaussDB MySQL account password
	var (
		updateGaussDBAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db-users/password"
	)

	updateGaussDBAccountPath := client.Endpoint + updateGaussDBAccountHttpUrl
	updateGaussDBAccountPath = strings.ReplaceAll(updateGaussDBAccountPath, "{project_id}", client.ProjectID)
	updateGaussDBAccountPath = strings.ReplaceAll(updateGaussDBAccountPath, "{instance_id}", instanceID)

	updateGaussDBAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	updateGaussDBAccountOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBAccountPasswordBodyParams(d))

	var updateGaussDBAccountResp *http.Response
	var err error
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		updateGaussDBAccountResp, err = client.Request("PUT", updateGaussDBAccountPath, &updateGaussDBAccountOpt)
		isRetry, err := handleGaussDBMysqlOperationError(err)
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error updating GaussDB account password: %s", err)
	}

	updateGaussDBAccountRespBody, err := utils.FlattenResponse(updateGaussDBAccountResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateGaussDBAccountRespBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("unable to find the job ID of the GaussDB MySQL account from the API response while updating the password")
	}

	return waitForJobComplete(ctx, client, d.Timeout(schema.TimeoutUpdate), instanceID, jobId)
}

func buildUpdateGaussDBAccountDescriptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":    utils.ValueIgnoreEmpty(d.Get("name")),
		"host":    utils.ValueIgnoreEmpty(d.Get("host")),
		"comment": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	param := map[string]interface{}{
		"users": []interface{}{bodyParams},
	}
	return param
}

func buildUpdateGaussDBAccountPasswordBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":     utils.ValueIgnoreEmpty(d.Get("name")),
		"host":     utils.ValueIgnoreEmpty(d.Get("host")),
		"password": utils.ValueIgnoreEmpty(d.Get("password")),
	}
	param := map[string]interface{}{
		"users": []interface{}{bodyParams},
	}
	return param
}

func resourceGaussDBAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteGaussDBAccount: delete the GaussDB MySQL account
	var (
		deleteGaussDBAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db-users"
		deleteGaussDBAccountProduct = "gaussdb"
	)
	deleteGaussDBAccountClient, err := cfg.NewServiceClient(deleteGaussDBAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deleteGaussDBAccountPath := deleteGaussDBAccountClient.Endpoint + deleteGaussDBAccountHttpUrl
	deleteGaussDBAccountPath = strings.ReplaceAll(deleteGaussDBAccountPath, "{project_id}",
		deleteGaussDBAccountClient.ProjectID)
	deleteGaussDBAccountPath = strings.ReplaceAll(deleteGaussDBAccountPath, "{instance_id}", instanceID)

	deleteGaussDBAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	deleteGaussDBAccountOpt.JSONBody = utils.RemoveNil(buildDeleteGaussDBAccountBodyParams(d))

	var deleteGaussDBAccountResp *http.Response
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		deleteGaussDBAccountResp, err = deleteGaussDBAccountClient.Request("DELETE",
			deleteGaussDBAccountPath, &deleteGaussDBAccountOpt)
		isRetry, err := handleGaussDBMysqlOperationError(err)
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("error deleting GaussDB MySQL account: %s", err)
	}

	deleteGaussDBAccountRespBody, err := utils.FlattenResponse(deleteGaussDBAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteGaussDBAccountRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID of the GaussDB MySQL account from the API response")
	}

	err = waitForJobComplete(ctx, deleteGaussDBAccountClient, d.Timeout(schema.TimeoutDelete), instanceID, jobId)

	return nil
}

func buildDeleteGaussDBAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	host := "%"
	if h := d.Get("host").(string); h != "" {
		host = h
	}
	bodyParams := map[string]interface{}{
		"name": utils.ValueIgnoreEmpty(d.Get("name")),
		"host": host,
	}
	param := map[string]interface{}{
		"users": []interface{}{bodyParams},
	}
	return param
}
