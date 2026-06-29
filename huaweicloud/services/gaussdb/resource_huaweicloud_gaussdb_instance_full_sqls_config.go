package gaussdb

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var gaussdbFullSqlsNonUpdatableParams = []string{"instance_id", "storage_mode"}

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/full-sqls/start
// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/full-sqls/stop
func ResourceGaussDbFullSqlsConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFullSqlConfigCreate,
		ReadContext:   resourceFullSqlConfigRead,
		UpdateContext: resourceFullSqlConfigUpdate,
		DeleteContext: resourceFullSqlConfigDelete,

		CustomizeDiff: config.FlexibleForceNew(gaussdbFullSqlsNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
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
			},
			"storage_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_exclude_sys_user": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"save_days": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"lts_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_group_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"log_stream_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"sql_type_range": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Required: true,
						},
						"prefixes": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildFullSqlConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"storage_mode":        d.Get("storage_mode"),
		"save_days":           d.Get("save_days"),
		"is_exclude_sys_user": utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "is_exclude_sys_user"),
		"lts_config":          buildLtsConfigBodyParams(d.Get("lts_config").([]interface{})),
	}

	if v, ok := d.GetOk("sql_type_range"); ok {
		params["sql_type_range"] = buildSqlTypeRangeBodyParams(v.([]interface{}))
	}

	return params
}

func buildLtsConfigBodyParams(ltsConfigs []interface{}) map[string]interface{} {
	if len(ltsConfigs) == 0 {
		return nil
	}

	ltsConfig, ok := ltsConfigs[0].(map[string]interface{})
	if !ok {
		return nil
	}
	return map[string]interface{}{
		"log_group_name":  ltsConfig["log_group_name"],
		"log_stream_name": ltsConfig["log_stream_name"],
	}
}

func buildSqlTypeRangeBodyParams(sqlTypeRanges []interface{}) []map[string]interface{} {
	if len(sqlTypeRanges) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(sqlTypeRanges))
	for _, item := range sqlTypeRanges {
		sqlType, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		params := map[string]interface{}{
			"category": sqlType["category"],
			"prefixes": utils.ValueIgnoreEmpty(utils.ExpandToStringList(sqlType["prefixes"].([]interface{}))),
		}

		rst = append(rst, params)
	}

	return rst
}

func resourceFullSqlConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/full-sqls/start"
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
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildFullSqlConfigBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error starting GaussDB full SQL collection: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, nil)
	if jobId == nil {
		return diag.Errorf("error starting GaussDB full SQL collection: unable to find job ID")
	}
	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 5, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceID)

	return nil
}

func resourceFullSqlConfigRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFullSqlConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/full-sqls/start"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildFullSqlConfigBodyParams(d)),
	}

	resp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating GaussDB full SQL configuration: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, nil)
	if jobId == nil {
		return diag.Errorf("error updating GaussDB full SQL collection: unable to find job ID")
	}
	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 5, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceFullSqlConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/full-sqls/stop"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error stopping GaussDB full SQL collection: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, nil)
	if jobId == nil {
		return diag.Errorf("error stopping GaussDB full SQL collection: unable to find job ID")
	}
	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 5, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
