package rds

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS POST /v3/{project_id}/instances/{instance_id}/sql-limit
// @API RDS GET /v3/{project_id}/instances/{instance_id}/sql-limit
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/sql-limit/update
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/sql-limit/switch
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/sql-limit
func ResourcePgSqlLimit() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePgSqlLimitCreate,
		UpdateContext: resourcePgSqlLimitUpdate,
		ReadContext:   resourcePgSqlLimitRead,
		DeleteContext: resourcePgSqlLimitDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRdsSqlLimitImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
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
				Description: `Specifies the ID of the RDS PostgreSQL instance.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the database.`,
			},
			"query_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"query_id", "query_string"},
				Description:  `Specifies the query ID`,
			},
			"query_string": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the text form of SQL statement.`,
			},
			"max_concurrency": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the number of SQL statements executed simultaneously`,
			},
			"max_waiting": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the max waiting time in seconds.`,
			},
			"search_path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the query order for names that are not schema qualified.`,
			},
			"switch": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the SQL limit switch.`,
			},
			"sql_limit_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of SQL limit.`,
			},
			"is_effective": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the SQL limit is effective.`,
			},
		},
	}
}

func resourcePgSqlLimitCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/sql-limit"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateSqlLimitBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating RDS PostgreSQL SQL limit: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resp := utils.PathSearch("resp", createRespBody, nil)
	if resp == nil {
		return diag.Errorf("unable to find the resp from the response: %s", err)
	}
	if resp.(string) != "successful" {
		return diag.Errorf("error creating RDS PostgreSQL SQL limit, the response is: %s", resp.(string))
	}

	queryField := "query_string"
	queryValue := d.Get("query_string")
	if v, ok := d.GetOk("query_id"); ok {
		queryField = "query_id"
		queryValue = v
	}
	sqlLimit, err := getSqlLimit(client, d, queryField, queryValue.(string))
	if err != nil {
		return diag.FromErr(err)
	}
	sqlLimitId := utils.PathSearch("id", sqlLimit, "").(string)

	// it is required when get the SQL limit it read method
	err = d.Set("sql_limit_id", sqlLimitId)
	if err != nil {
		return diag.FromErr(err)
	}

	dbName := d.Get("db_name").(string)
	d.SetId(fmt.Sprintf("%s/%s/%s", instanceID, dbName, sqlLimitId))

	if v, ok := d.GetOk("switch"); ok && v == "open" {
		err = updateSqlLimitSwitch(client, d, sqlLimitId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePgSqlLimitRead(ctx, d, meta)
}

func buildCreateSqlLimitBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name":         d.Get("db_name"),
		"max_concurrency": d.Get("max_concurrency"),
		"max_waiting":     d.Get("max_waiting"),
		"query_id":        utils.ValueIgnoreEmpty(d.Get("query_id")),
		"query_string":    utils.ValueIgnoreEmpty(d.Get("query_string")),
		"search_path":     utils.ValueIgnoreEmpty(d.Get("search_path")),
	}
	return bodyParams
}

func resourcePgSqlLimitRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	sqlLimit, err := getSqlLimit(client, d, "id", d.Get("sql_limit_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.01010340"),
			"error retrieving RDS PostgreSQL SQL limit")
	}

	isEffective := utils.PathSearch("is_effective", sqlLimit, false).(bool)
	sqlLimitSwitch := "close"
	if isEffective {
		sqlLimitSwitch = "open"
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("sql_limit_id", utils.PathSearch("id", sqlLimit, nil)),
		d.Set("query_id", utils.PathSearch("query_id", sqlLimit, nil)),
		d.Set("query_string", utils.PathSearch("query_string", sqlLimit, nil)),
		d.Set("max_concurrency", utils.PathSearch("max_concurrency", sqlLimit, nil)),
		d.Set("max_waiting", utils.PathSearch("max_waiting", sqlLimit, nil)),
		d.Set("search_path", utils.PathSearch("search_path", sqlLimit, nil)),
		d.Set("is_effective", isEffective),
		d.Set("switch", sqlLimitSwitch),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getSqlLimit(client *golangsdk.ServiceClient, d *schema.ResourceData, queryField, queryValue string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/sql-limit?db_name={db_name}"
	)
	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", d.Get("instance_id").(string))
	getBasePath = strings.ReplaceAll(getBasePath, "{db_name}", d.Get("db_name").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var currentTotal int
	var getPath string
	for {
		getPath = getBasePath + buildGetSqlLimitQueryParams(currentTotal)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, err
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}

		sqlLimitObjects := utils.PathSearch("sql_limit_objects", getRespBody, make([]interface{}, 0)).([]interface{})
		sqlLimit := utils.PathSearch(fmt.Sprintf("[?%s == '%s']|[0]", queryField, queryValue), sqlLimitObjects, nil)
		if sqlLimit != nil {
			return sqlLimit, nil
		}
		currentTotal += len(sqlLimitObjects)
		total := utils.PathSearch("total", getRespBody, float64(0)).(float64)
		if currentTotal >= int(total) {
			break
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func buildGetSqlLimitQueryParams(offset int) string {
	return fmt.Sprintf("&limit=100&offset=%v", offset)
}

func resourcePgSqlLimitUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	if d.HasChanges("max_concurrency", "max_waiting") {
		err = updateSqlLimit(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChanges("switch") {
		err = updateSqlLimitSwitch(client, d, d.Get("sql_limit_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePgSqlLimitRead(ctx, d, meta)
}

func updateSqlLimit(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/sql-limit/update"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateSqlLimitBodyParams(d))

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating RDS PostgreSQL SQL limit: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	resp := utils.PathSearch("resp", updateRespBody, nil)
	if resp == nil {
		return fmt.Errorf("unable to find the resp from the response: %s", err)
	}
	if resp.(string) != "successful" {
		return fmt.Errorf("error updating RDS PostgreSQL SQL limit, the response is: %s", resp.(string))
	}
	return nil
}

func buildUpdateSqlLimitBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name":         d.Get("db_name"),
		"id":              d.Get("sql_limit_id"),
		"max_concurrency": d.Get("max_concurrency"),
		"max_waiting":     d.Get("max_waiting"),
	}
	return bodyParams
}

func updateSqlLimitSwitch(client *golangsdk.ServiceClient, d *schema.ResourceData, sqlLimitId string) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/sql-limit/switch"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateSqlLimitSwitchBodyParams(d, sqlLimitId))

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating RDS PostgreSQL SQL limit switch: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	resp := utils.PathSearch("resp", updateRespBody, nil)
	if resp == nil {
		return fmt.Errorf("unable to find the resp from the response: %s", err)
	}
	if resp.(string) != "successful" {
		return fmt.Errorf("error updating RDS PostgreSQL SQL limit switch, the response is: %s", resp.(string))
	}
	return nil
}

func buildUpdateSqlLimitSwitchBodyParams(d *schema.ResourceData, sqlLimitId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name": d.Get("db_name"),
		"id":      sqlLimitId,
		"action":  d.Get("switch"),
	}
	return bodyParams
}

func resourcePgSqlLimitDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/sql-limit"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteSqlLimitBodyParams(d))

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting RDS PostgreSQL SQL limit: %s", err)
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resp := utils.PathSearch("resp", deleteRespBody, nil)
	if resp == nil {
		return diag.Errorf("unable to find the resp from the response: %s", err)
	}
	if resp.(string) != "successful" {
		return diag.Errorf("error deleting RDS PostgreSQL SQL limit, the response is: %s", resp.(string))
	}

	return nil
}

func buildDeleteSqlLimitBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name": d.Get("db_name"),
		"id":      d.Get("sql_limit_id"),
	}
	return bodyParams
}

func resourceRdsSqlLimitImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<db_name>/<sql_limit_id>")
	}

	mErr := multierror.Append(
		d.Set("instance_id", parts[0]),
		d.Set("db_name", parts[1]),
		d.Set("sql_limit_id", parts[2]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
