// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDM
// ---------------------------------------------------------------

package ddm

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type ddmError struct {
	ErrorCode string `json:"errCode"`
	ErrorMsg  string `json:"externalMessage"`
}

// @API DDM POST /v1/{project_id}/instances/{instance_id}/databases
// @API DDM GET /v1/{project_id}/instances/{instance_id}/databases/{ddm_dbname}
// @API DDM DELETE /v1/{project_id}/instances/{instance_id}/databases/{ddm_dbname}
func ResourceDdmSchema() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdmSchemaCreate,
		ReadContext:   resourceDdmSchemaRead,
		DeleteContext: resourceDdmSchemaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
				Description: `Specifies the ID of a DDM instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the DDM schema.`,
			},
			"shard_mode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the sharding mode of the schema.`,
				ValidateFunc: validation.StringInSlice([]string{
					"cluster", "single",
				}, false),
			},
			"shard_number": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the number of shards in the same working mode.`,
			},
			"data_nodes": {
				Type:        schema.TypeList,
				Elem:        SchemaDataNodeSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the RDS instances associated with the schema.`,
			},
			"delete_rds_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether data stored on the associated DB instances is deleted`,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the schema status.`,
			},
			"shards": {
				Type:        schema.TypeList,
				Elem:        SchemaShardSchema(),
				Computed:    true,
				Description: `Indicates the sharding information of the schema.`,
			},
			"data_vips": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the IP address and port number for connecting to the schema.`,
			},
		},
	}
}

func SchemaDataNodeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the RDS instance associated with the schema.`,
			},
			"admin_user": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the username for logging in to the associated RDS instance.`,
			},
			"admin_password": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: `Specifies the password for logging in to the associated RDS instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the associated RDS instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the associated RDS instance.`,
			},
		},
	}
	return &sc
}

func SchemaShardSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_slot": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of shards.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the shard name.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the shard status.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the RDS instance where the shard is located.`,
			},
		},
	}
	return &sc
}

func resourceDdmSchemaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createSchema: create DDM schema
	var (
		createSchemaHttpUrl = "v1/{project_id}/instances/{instance_id}/databases"
		createSchemaProduct = "ddm"
	)
	createSchemaClient, err := cfg.NewServiceClient(createSchemaProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createSchemaPath := createSchemaClient.Endpoint + createSchemaHttpUrl
	createSchemaPath = strings.ReplaceAll(createSchemaPath, "{project_id}", createSchemaClient.ProjectID)
	createSchemaPath = strings.ReplaceAll(createSchemaPath, "{instance_id}", fmt.Sprintf("%v", instanceID))

	createSchemaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createSchemaOpt.JSONBody = utils.RemoveNil(buildCreateSchemaBodyParams(d))

	var createSchemaResp *http.Response
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		createSchemaResp, err = createSchemaClient.Request("POST", createSchemaPath, &createSchemaOpt)
		isRetry, err := handleOperationError(err, "creating", "schema")
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	createSchemaRespBody, err := utils.FlattenResponse(createSchemaResp)
	if err != nil {
		return diag.FromErr(err)
	}

	schemaName := utils.PathSearch("databases[0].name", createSchemaRespBody, "").(string)
	if schemaName == "" {
		return diag.Errorf("unable to find the DDM schema name from the API response %s", err)
	}

	err = waitForInstanceRunning(ctx, d, createSchemaClient, instanceID, schema.TimeoutCreate)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceID + "/" + schemaName)

	return resourceDdmSchemaRead(ctx, d, meta)
}

func buildCreateSchemaBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":         utils.ValueIgnoreEmpty(d.Get("name")),
		"shard_mode":   utils.ValueIgnoreEmpty(d.Get("shard_mode")),
		"shard_number": utils.ValueIgnoreEmpty(d.Get("shard_number")),
		"used_rds":     buildCreateSchemaUsedRdsChildBody(d),
	}
	params := map[string]interface{}{
		"databases": []interface{}{bodyParams},
	}
	return params
}

func buildCreateSchemaUsedRdsChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("data_nodes").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}
	params := make([]map[string]interface{}, 0)
	for _, param := range rawParams {
		perm := make(map[string]interface{})
		perm["id"] = utils.PathSearch("id", param, nil)
		perm["adminUser"] = utils.PathSearch("admin_user", param, nil)
		perm["adminPassword"] = utils.PathSearch("admin_password", param, nil)
		params = append(params, perm)
	}

	return params
}

func resourceDdmSchemaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSchema: Query DDM schema
	var (
		getSchemaHttpUrl = "v1/{project_id}/instances/{instance_id}/databases/{ddm_dbname}"
		getSchemaProduct = "ddm"
	)
	getSchemaClient, err := cfg.NewServiceClient(getSchemaProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<schema_name>")
	}
	instanceID := parts[0]
	schemaName := parts[1]
	getSchemaPath := getSchemaClient.Endpoint + getSchemaHttpUrl
	getSchemaPath = strings.ReplaceAll(getSchemaPath, "{project_id}", getSchemaClient.ProjectID)
	getSchemaPath = strings.ReplaceAll(getSchemaPath, "{instance_id}", fmt.Sprintf("%v", instanceID))
	getSchemaPath = strings.ReplaceAll(getSchemaPath, "{ddm_dbname}", fmt.Sprintf("%v", schemaName))

	getSchemaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSchemaResp, err := getSchemaClient.Request("GET", getSchemaPath, &getSchemaOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DdmSchema")
	}

	getSchemaRespBody, err := utils.FlattenResponse(getSchemaResp)
	if err != nil {
		return diag.FromErr(err)
	}

	schemas := utils.PathSearch("database", getSchemaRespBody, nil)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", schemas, nil)),
		d.Set("status", utils.PathSearch("status", schemas, nil)),
		d.Set("shards", flattenGetSchemaResponseBodyShard(schemas)),
		d.Set("shard_mode", utils.PathSearch("shard_mode", schemas, nil)),
		d.Set("shard_number", utils.PathSearch("shard_number", schemas, nil)),
		d.Set("data_vips", utils.PathSearch("dataVips", schemas, nil)),
		d.Set("data_nodes", flattenGetSchemaResponseBodyDataNode(schemas)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetSchemaResponseBodyShard(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"db_slot": utils.PathSearch("dbslot", v, nil),
			"name":    utils.PathSearch("name", v, nil),
			"status":  utils.PathSearch("status", v, nil),
			"id":      utils.PathSearch("id", v, nil),
		})
	}
	return rst
}

func flattenGetSchemaResponseBodyDataNode(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("used_rds", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":     utils.PathSearch("id", v, nil),
			"name":   utils.PathSearch("name", v, nil),
			"status": utils.PathSearch("status", v, nil),
		})
	}
	return rst
}

func resourceDdmSchemaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteSchema: Delete DDM schema
	var (
		deleteSchemaHttpUrl = "v1/{project_id}/instances/{instance_id}/databases/{ddm_dbname}"
		deleteSchemaProduct = "ddm"
	)
	deleteSchemaClient, err := cfg.NewServiceClient(deleteSchemaProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<schema_name>")
	}
	instanceID := parts[0]
	schemaName := parts[1]
	deleteSchemaPath := deleteSchemaClient.Endpoint + deleteSchemaHttpUrl
	deleteSchemaPath = strings.ReplaceAll(deleteSchemaPath, "{project_id}", deleteSchemaClient.ProjectID)
	deleteSchemaPath = strings.ReplaceAll(deleteSchemaPath, "{instance_id}", fmt.Sprintf("%v", instanceID))
	deleteSchemaPath = strings.ReplaceAll(deleteSchemaPath, "{ddm_dbname}", schemaName)

	deleteSchemaQueryParams := buildDeleteSchemaQueryParams(d)
	deleteSchemaPath += deleteSchemaQueryParams

	deleteSchemaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err = deleteSchemaClient.Request("DELETE", deleteSchemaPath, &deleteSchemaOpt)
		isRetry, err := handleOperationError(err, "deleting", "schema")
		if isRetry {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitForInstanceRunning(ctx, d, deleteSchemaClient, instanceID, schema.TimeoutDelete)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildDeleteSchemaQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("delete_rds_data"); ok {
		res = fmt.Sprintf("%s&delete_rds_data=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func handleOperationError(err error, operateType string, operateObj string) (bool, error) {
	if err == nil {
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault409); ok {
		var apiError ddmError
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("error %s DDM %s: error format error: %s", operateType,
				operateObj, jsonErr)
		}
		if apiError.ErrorCode == "DBS.200019" {
			return true, err
		}
	}
	return false, err
}

func waitForInstanceRunning(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, instanceID string,
	timeout string) error {

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"RUNNING"},
		Refresh:      ddmInstanceStatusRefreshFunc(instanceID, client),
		Timeout:      d.Timeout(timeout),
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) to running: %s", instanceID, err)
	}
	return nil
}
