package gaussdb

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

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/schema
// @API GaussDB GET /v3/{project_id}/instances
// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/schemas
// @API GaussDB DELETE /v3/{project_id}/instances/{instance_id}/schema
func ResourceOpenGaussSchema() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussSchemaCreate,
		ReadContext:   resourceOpenGaussSchemaRead,
		DeleteContext: resourceOpenGaussSchemaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceOpenGaussSchemaImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOpenGaussSchemaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/schema"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateOpenGaussSchemaBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating GaussDB openGauss schema: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceID, d.Get("db_name").(string), d.Get("name").(string)))

	return resourceOpenGaussSchemaRead(ctx, d, meta)
}

func buildCreateOpenGaussSchemaBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name": d.Get("db_name"),
		"schemas": []map[string]interface{}{
			{
				"name":  d.Get("name"),
				"owner": d.Get("owner"),
			},
		},
	}
	return bodyParams
}

func resourceOpenGaussSchemaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/schemas"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// actually, offset is pageNo in API, not offset
	var offset int
	var dbSchema interface{}
	dbName := d.Get("db_name").(string)
	schemaName := d.Get("name").(string)

	for {
		getPath := getBasePath + buildOpenGaussSchemaQueryParams(dbName, offset)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.200823"),
				"error retrieving GaussDB openGauss schema")
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		dbSchema = utils.PathSearch(fmt.Sprintf("database_schemas[?schema_name=='%s']|[0]", schemaName), getRespBody, nil)
		if dbSchema != nil {
			break
		}
		totalCount := utils.PathSearch("total_count", getRespBody, float64(0)).(float64)
		if int(totalCount) <= (offset+1)*100 {
			break
		}
		offset++
	}
	if dbSchema == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB openGauss schema")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("db_name", dbName),
		d.Set("name", utils.PathSearch("schema_name", dbSchema, nil)),
		d.Set("owner", utils.PathSearch("owner", dbSchema, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildOpenGaussSchemaQueryParams(dbName string, offset int) string {
	return fmt.Sprintf("?db_name=%s&limit=100&offset=%v", dbName, offset)
}

func resourceOpenGaussSchemaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/schema"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteOpenGaussSchemaBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code",
			[]string{"DBS.200823", "DBS.06280007"}...), "error deleting GaussDB openGauss schema")
	}

	return nil
}

func buildDeleteOpenGaussSchemaBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"db_name": d.Get("db_name"),
		"schema":  d.Get("name"),
	}
	return bodyParams
}

func resourceOpenGaussSchemaImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<db_name/><name>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("db_name", parts[1]),
		d.Set("name", parts[2]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
