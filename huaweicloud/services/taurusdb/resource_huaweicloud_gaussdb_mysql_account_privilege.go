// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GaussDB
// ---------------------------------------------------------------

package taurusdb

import (
	"context"
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

// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/db-users/privilege
// @API GaussDBforMySQL GET /v3/{project_id}/jobs
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/db-users
// @API GaussDBforMySQL DELETE /v3/{project_id}/instances/{instance_id}/db-users/privilege
func ResourceGaussDBAccountPrivilege() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBAccountPrivilegeCreate,
		ReadContext:   resourceGaussDBAccountPrivilegeRead,
		DeleteContext: resourceGaussDBAccountPrivilegeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
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
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the database username.`,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `Specifies the host IP address which allow database users to connect to the database
on the current host`,
			},
			"databases": {
				Type:        schema.TypeList,
				Elem:        AccountPrivilegeDatabaseSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the list of the databases.`,
			},
		},
	}
}

func AccountPrivilegeDatabaseSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the database name.`,
			},
			"readonly": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies whether the database permission is read-only.`,
			},
		},
	}
	return &sc
}

func resourceGaussDBAccountPrivilegeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createGaussDBAccount: create a GaussDB MySQL account privilege
	var (
		createGaussDBAccountPrivilegeHttpUrl = "v3/{project_id}/instances/{instance_id}/db-users/privilege"
		createGaussDBAccountPrivilegeProduct = "gaussdb"
	)
	createGaussDBAccountPrivilegeClient, err := cfg.NewServiceClient(createGaussDBAccountPrivilegeProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	accountName := d.Get("account_name").(string)
	host := d.Get("host").(string)
	createGaussDBAccountPrivilegePath := createGaussDBAccountPrivilegeClient.Endpoint + createGaussDBAccountPrivilegeHttpUrl
	createGaussDBAccountPrivilegePath = strings.ReplaceAll(createGaussDBAccountPrivilegePath, "{project_id}",
		createGaussDBAccountPrivilegeClient.ProjectID)
	createGaussDBAccountPrivilegePath = strings.ReplaceAll(createGaussDBAccountPrivilegePath, "{instance_id}", instanceID)

	createGaussDBAccountPrivilegeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createGaussDBAccountPrivilegeOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBAccountPrivilegeBodyParams(d, accountName, host))

	var createGaussDBAccountPrivilegeResp *http.Response
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		createGaussDBAccountPrivilegeResp, err = createGaussDBAccountPrivilegeClient.Request("POST",
			createGaussDBAccountPrivilegePath, &createGaussDBAccountPrivilegeOpt)
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
		return diag.Errorf("error creating GaussDB account privilege: %s", err)
	}

	createGaussDBAccountPrivilegeRespBody, err := utils.FlattenResponse(createGaussDBAccountPrivilegeResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createGaussDBAccountPrivilegeRespBody, "")
	if jobId == "" {
		return diag.Errorf("unable to find the job ID of the GaussDB MySQL account privilege from the API response")
	}
	err = waitForJobComplete(ctx, createGaussDBAccountPrivilegeClient, d.Timeout(schema.TimeoutCreate), instanceID, jobId.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceID + "/" + accountName + "/" + host)

	return resourceGaussDBAccountPrivilegeRead(ctx, d, meta)
}

func buildCreateGaussDBAccountPrivilegeBodyParams(d *schema.ResourceData, accountName, host string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":      accountName,
		"host":      host,
		"databases": buildCreateGaussDBAccountPrivilegeDatabasesChildBody(d),
	}
	params := map[string]interface{}{
		"users": []interface{}{bodyParams},
	}
	return params
}

func buildCreateGaussDBAccountPrivilegeDatabasesChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("databases").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}
	params := make([]map[string]interface{}, 0)
	for _, param := range rawParams {
		perm := make(map[string]interface{})
		perm["name"] = utils.PathSearch("name", param, nil)
		perm["readonly"] = utils.PathSearch("readonly", param, nil)
		params = append(params, perm)
	}
	return params
}

func resourceGaussDBAccountPrivilegeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getGaussDBAccountPrivilege: Query the GaussDB MySQL account privilege
	var (
		getGaussDBAccountPrivilegeHttpUrl = "v3/{project_id}/instances/{instance_id}/db-users"
		getGaussDBAccountPrivilegeProduct = "gaussdb"
	)
	getGaussDBAccountPrivilegeClient, err := cfg.NewServiceClient(getGaussDBAccountPrivilegeProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return diag.Errorf("invalid id format, must be <instance_id>/<account_name>/<host>")
	}
	instanceID := parts[0]
	accountName := parts[1]
	host := parts[2]

	getGaussDBAccountPrivilegeBasePath := getGaussDBAccountPrivilegeClient.Endpoint + getGaussDBAccountPrivilegeHttpUrl
	getGaussDBAccountPrivilegeBasePath = strings.ReplaceAll(getGaussDBAccountPrivilegeBasePath, "{project_id}",
		getGaussDBAccountPrivilegeClient.ProjectID)
	getGaussDBAccountPrivilegeBasePath = strings.ReplaceAll(getGaussDBAccountPrivilegeBasePath, "{instance_id}", instanceID)

	getGaussDBAccountPrivilegeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	var currentTotal int
	var account interface{}
	getGaussDBAccountPrivilegePath := getGaussDBAccountPrivilegeBasePath + buildGaussDBMysqlQueryParams(currentTotal)

	for {
		getGaussDBAccountResp, err := getGaussDBAccountPrivilegeClient.Request("GET", getGaussDBAccountPrivilegePath,
			&getGaussDBAccountPrivilegeOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB MySQL account privilege")
		}

		getGaussDBAccountRespBody, err := utils.FlattenResponse(getGaussDBAccountResp)
		if err != nil {
			return diag.FromErr(err)
		}
		res, pageNum := flattenGetGaussDBAccountPrivilegeResponseBody(getGaussDBAccountRespBody, accountName, host)
		if res != nil {
			account = res
			break
		}
		total := utils.PathSearch("total_count", getGaussDBAccountRespBody, float64(0)).(float64)
		currentTotal += pageNum
		if currentTotal == int(total) {
			break
		}
		getGaussDBAccountPrivilegePath = getGaussDBAccountPrivilegeBasePath + buildGaussDBMysqlQueryParams(currentTotal)
	}
	if account == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	databases := utils.PathSearch("databases", account, make([]interface{}, 0)).([]interface{})
	if len(databases) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("account_name", utils.PathSearch("name", account, nil)),
		d.Set("host", utils.PathSearch("host", account, nil)),
		d.Set("databases", flattenGetGaussDBAccountPrivilegeResponseBodyDatabase(databases)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetGaussDBAccountPrivilegeResponseBody(resp interface{}, accountName, address string) (interface{}, int) {
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

func flattenGetGaussDBAccountPrivilegeResponseBodyDatabase(databases []interface{}) []interface{} {
	if databases == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(databases))
	for _, v := range databases {
		rst = append(rst, map[string]interface{}{
			"name":     utils.PathSearch("name", v, nil),
			"readonly": utils.PathSearch("readonly", v, nil),
		})
	}
	return rst
}

func resourceGaussDBAccountPrivilegeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteGaussDBAccountPrivilege: delete the GaussDB MySQL account privilege
	var (
		deleteGaussDBAccountPrivilegeHttpUrl = "v3/{project_id}/instances/{instance_id}/db-users/privilege"
		deleteGaussDBAccountPrivilegeProduct = "gaussdb"
	)
	deleteGaussDBAccountPrivilegeClient, err := cfg.NewServiceClient(deleteGaussDBAccountPrivilegeProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deleteGaussDBAccountPrivilegePath := deleteGaussDBAccountPrivilegeClient.Endpoint + deleteGaussDBAccountPrivilegeHttpUrl
	deleteGaussDBAccountPrivilegePath = strings.ReplaceAll(deleteGaussDBAccountPrivilegePath, "{project_id}",
		deleteGaussDBAccountPrivilegeClient.ProjectID)
	deleteGaussDBAccountPrivilegePath = strings.ReplaceAll(deleteGaussDBAccountPrivilegePath, "{instance_id}", instanceID)

	deleteGaussDBAccountPrivilegeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	deleteGaussDBAccountPrivilegeOpt.JSONBody = utils.RemoveNil(buildDeleteGaussDBAccountPrivilegeBodyParams(d))

	var deleteGaussDBAccountPrivilegeResp *http.Response
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		deleteGaussDBAccountPrivilegeResp, err = deleteGaussDBAccountPrivilegeClient.Request("DELETE",
			deleteGaussDBAccountPrivilegePath, &deleteGaussDBAccountPrivilegeOpt)
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
		return diag.Errorf("error deleting GaussDB MySQL account privilege: %s", err)
	}

	deleteGaussDBAccountPrivilegeRespBody, err := utils.FlattenResponse(deleteGaussDBAccountPrivilegeResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteGaussDBAccountPrivilegeRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID of the GaussDB MySQL account privilege from the API response")
	}

	err = waitForJobComplete(ctx, deleteGaussDBAccountPrivilegeClient, d.Timeout(schema.TimeoutDelete), instanceID, jobId)

	return nil
}

func buildDeleteGaussDBAccountPrivilegeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":      utils.ValueIgnoreEmpty(d.Get("account_name")),
		"host":      utils.ValueIgnoreEmpty(d.Get("host")),
		"databases": buildCreateGaussDBAccountPrivilegeDatabaseNamesChildBody(d),
	}
	params := map[string]interface{}{
		"users": []interface{}{bodyParams},
	}
	return params
}

func buildCreateGaussDBAccountPrivilegeDatabaseNamesChildBody(d *schema.ResourceData) []interface{} {
	rawParams := d.Get("databases").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}
	params := make([]interface{}, 0)
	for _, param := range rawParams {
		name := utils.PathSearch("name", param, nil)
		params = append(params, name)
	}
	return params
}
