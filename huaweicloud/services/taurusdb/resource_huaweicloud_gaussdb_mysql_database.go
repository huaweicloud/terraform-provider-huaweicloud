// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GaussDB
// ---------------------------------------------------------------

package taurusdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/databases
// @API GaussDBforMySQL GET /v3/{project_id}/jobs
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/databases/comment
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/databases
// @API GaussDBforMySQL DELETE /v3/{project_id}/instances/{instance_id}/databases
func ResourceGaussDBDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBDatabaseCreate,
		UpdateContext: resourceGaussDBDatabaseUpdate,
		ReadContext:   resourceGaussDBDatabaseRead,
		DeleteContext: resourceGaussDBDatabaseDelete,
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
				Description: `Specifies the database name.`,
			},
			"character_set": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the database character set.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the database remarks.`,
			},
		},
	}
}

func resourceGaussDBDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createGaussDBDatabase: create a GaussDB MySQL database
	var (
		createGaussDBDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/databases"
		createGaussDBDatabaseProduct = "gaussdb"
	)
	createGaussDBDatabaseClient, err := cfg.NewServiceClient(createGaussDBDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createGaussDBDatabasePath := createGaussDBDatabaseClient.Endpoint + createGaussDBDatabaseHttpUrl
	createGaussDBDatabasePath = strings.ReplaceAll(createGaussDBDatabasePath, "{project_id}",
		createGaussDBDatabaseClient.ProjectID)
	createGaussDBDatabasePath = strings.ReplaceAll(createGaussDBDatabasePath, "{instance_id}", instanceID)

	createGaussDBDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createGaussDBDatabaseOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBDatabaseBodyParams(d))

	var createGaussDBDatabaseResp *http.Response
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		createGaussDBDatabaseResp, err = createGaussDBDatabaseClient.Request("POST",
			createGaussDBDatabasePath, &createGaussDBDatabaseOpt)
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
		return diag.Errorf("error creating GaussDB MySQL database: %s", err)
	}

	createGaussDBDatabaseRespBody, err := utils.FlattenResponse(createGaussDBDatabaseResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createGaussDBDatabaseRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID of the GaussDB MySQL database from the API response")
	}

	err = waitForJobComplete(ctx, createGaussDBDatabaseClient, d.Timeout(schema.TimeoutCreate), instanceID, jobId)
	if err != nil {
		return diag.FromErr(err)
	}

	databaseName := d.Get("name").(string)
	d.SetId(instanceID + "/" + databaseName)

	return resourceGaussDBDatabaseRead(ctx, d, meta)
}

func buildCreateGaussDBDatabaseBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name":          utils.ValueIgnoreEmpty(d.Get("name")),
		"comment":       utils.ValueIgnoreEmpty(d.Get("description")),
		"character_set": utils.ValueIgnoreEmpty(d.Get("character_set")),
	}
	bodyParams := map[string]interface{}{
		"databases": []interface{}{params},
	}
	return bodyParams
}

func resourceGaussDBDatabaseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getGaussDBDatabase: Query the GaussDB MySQL database
	var (
		getGaussDBDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/databases"
		getGaussDBDatabaseProduct = "gaussdb"
	)
	getGaussDBDatabaseClient, err := cfg.NewServiceClient(getGaussDBDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<name>")
	}
	instanceID := parts[0]
	databaseName := parts[1]

	getGaussDBDatabaseBasePath := getGaussDBDatabaseClient.Endpoint + getGaussDBDatabaseHttpUrl
	getGaussDBDatabaseBasePath = strings.ReplaceAll(getGaussDBDatabaseBasePath, "{project_id}",
		getGaussDBDatabaseClient.ProjectID)
	getGaussDBDatabaseBasePath = strings.ReplaceAll(getGaussDBDatabaseBasePath, "{instance_id}", instanceID)

	getGaussDBDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	var currentTotal int
	var database interface{}
	getGaussDBDatabasePath := getGaussDBDatabaseBasePath + buildGaussDBMysqlQueryParams(currentTotal)

	for {
		getGaussDBDatabaseResp, err := getGaussDBDatabaseClient.Request("GET", getGaussDBDatabasePath,
			&getGaussDBDatabaseOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB MySQL database")
		}

		getGaussDBDatabaseRespBody, err := utils.FlattenResponse(getGaussDBDatabaseResp)
		if err != nil {
			return diag.FromErr(err)
		}
		db, pageNum := flattenGetGaussDBDatabaseResponseBody(getGaussDBDatabaseRespBody, databaseName)
		if db != nil {
			database = db
			break
		}
		total := utils.PathSearch("total_count", getGaussDBDatabaseRespBody, float64(0)).(float64)
		currentTotal += pageNum
		if currentTotal == int(total) {
			break
		}
		getGaussDBDatabasePath = getGaussDBDatabaseBasePath + buildGaussDBMysqlQueryParams(currentTotal)
	}
	if database == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("name", utils.PathSearch("name", database, nil)),
		d.Set("description", utils.PathSearch("comment", database, nil)),
		d.Set("character_set", utils.PathSearch("charset", database, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGaussDBMysqlQueryParams(offset int) string {
	return fmt.Sprintf("?limit=100&offset=%v", offset)
}

func flattenGetGaussDBDatabaseResponseBody(resp interface{}, databaseName string) (interface{}, int) {
	if resp == nil {
		return nil, 0
	}
	curJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	for _, v := range curArray {
		name := utils.PathSearch("name", v, "").(string)
		if databaseName == name {
			return v, len(curArray)
		}
	}
	return nil, len(curArray)
}

func resourceGaussDBDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateGaussDBDatabase: update the GaussDB MySQL database
	var (
		updateGaussDBDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/databases/comment"
		updateGaussDBDatabaseProduct = "gaussdb"
	)
	updateGaussDBDatabaseClient, err := cfg.NewServiceClient(updateGaussDBDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	updateGaussDBDatabasePath := updateGaussDBDatabaseClient.Endpoint + updateGaussDBDatabaseHttpUrl
	updateGaussDBDatabasePath = strings.ReplaceAll(updateGaussDBDatabasePath, "{project_id}",
		updateGaussDBDatabaseClient.ProjectID)
	updateGaussDBDatabasePath = strings.ReplaceAll(updateGaussDBDatabasePath, "{instance_id}", instanceID)

	updateGaussDBDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	updateGaussDBDatabaseOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBDatabaseBodyParams(d))

	var updateGaussDBDatabaseResp *http.Response
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		updateGaussDBDatabaseResp, err = updateGaussDBDatabaseClient.Request("PUT",
			updateGaussDBDatabasePath, &updateGaussDBDatabaseOpt)
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
		return diag.Errorf("error updating GaussDB MySQL database: %s", err)
	}

	updateGaussDBDatabaseRespBody, err := utils.FlattenResponse(updateGaussDBDatabaseResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", updateGaussDBDatabaseRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID of the GaussDB MySQL database from the API response")
	}

	err = waitForJobComplete(ctx, updateGaussDBDatabaseClient, d.Timeout(schema.TimeoutUpdate), instanceID, jobId)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGaussDBDatabaseRead(ctx, d, meta)
}

func buildUpdateGaussDBDatabaseBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name":    utils.ValueIgnoreEmpty(d.Get("name")),
		"comment": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	bodyParams := map[string]interface{}{
		"databases": []interface{}{params},
	}
	return bodyParams
}

func resourceGaussDBDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteGaussDBDatabase: delete the GaussDB MySQL database
	var (
		deleteGaussDBDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/databases"
		deleteGaussDBDatabaseProduct = "gaussdb"
	)
	deleteGaussDBDatabaseClient, err := cfg.NewServiceClient(deleteGaussDBDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deleteGaussDBDatabasePath := deleteGaussDBDatabaseClient.Endpoint + deleteGaussDBDatabaseHttpUrl
	deleteGaussDBDatabasePath = strings.ReplaceAll(deleteGaussDBDatabasePath, "{project_id}",
		deleteGaussDBDatabaseClient.ProjectID)
	deleteGaussDBDatabasePath = strings.ReplaceAll(deleteGaussDBDatabasePath, "{instance_id}", instanceID)

	deleteGaussDBDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	deleteGaussDBDatabaseOpt.JSONBody = utils.RemoveNil(buildDeleteGaussDBDatabaseBodyParams(d))

	var deleteGaussDBDatabaseResp *http.Response
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		deleteGaussDBDatabaseResp, err = deleteGaussDBDatabaseClient.Request("DELETE",
			deleteGaussDBDatabasePath, &deleteGaussDBDatabaseOpt)
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
		return diag.Errorf("error deleting GaussDB MySQL database: %s", err)
	}

	deleteGaussDBDatabaseRespBody, err := utils.FlattenResponse(deleteGaussDBDatabaseResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteGaussDBDatabaseRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID of the GaussDB MySQL database from the API response")
	}

	err = waitForJobComplete(ctx, deleteGaussDBDatabaseClient, d.Timeout(schema.TimeoutDelete), instanceID, jobId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildDeleteGaussDBDatabaseBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"databases": []interface{}{utils.ValueIgnoreEmpty(d.Get("name"))},
	}
	return bodyParams
}

func waitForJobComplete(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instanceId, jobId string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending", "Running", "Failed"},
		Target:       []string{"Completed"},
		Refresh:      gaussDBMysqlDatabaseStatusRefreshFunc(client, jobId),
		Timeout:      timeout,
		Delay:        1 * time.Second,
		PollInterval: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for GaussDB MySQL job (%s) to complete: %s", instanceId, err)
	}
	return nil
}

func handleGaussDBMysqlOperationError(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault409); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, jsonErr
		}
		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}
		if errorCode == "DBS.200047" {
			return true, err
		}
	}
	return false, err
}
