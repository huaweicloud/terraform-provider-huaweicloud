package secmaster

import (
	"context"
	"errors"
	"fmt"
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

var collectConfigNonUpdatableParams = []string{
	"workspace_id",
	"dataspace_id",
	"dataspace_name",
	"region_id",
	"domain_id",
	"config.*.csvc_display",
	"config.*.csvc",
	"config.*.shards",
	"config.*.source_display",
	"config.*.source_id",
	"config.*.ttl",
	"config.*.accounts",
	"config.*.action",
	"config.*.all_accounts",
	"config.*.new_account_auto_access",
	"config.*.source_name",
	"lts_config.*.config_name",
	"lts_config.*.description",
	"lts_config.*.enable",
	"lts_config.*.log_group_id",
	"lts_config.*.log_stream_id",
	"lts_config.*.log_type",
	"lts_config.*.log_type_prefix",
	"lts_config.*.pipe_alias",
}

// @API SecMaster POST /v1/{project_id}/collector/cloudlogs/config
// @API SecMaster GET /v1/{project_id}/collector/cloudlogs/config
func ResourceCollectConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCollectConfigCreate,
		ReadContext:   resourceCollectConfigRead,
		UpdateContext: resourceCollectConfigUpdate,
		DeleteContext: resourceCollectConfigDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceCollectConfigImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(collectConfigNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dataspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dataspace_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     collectConfigSchema(),
			},
			"lts_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     ltsConfigSchema(),
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

func collectConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"csvc_display": {
				Type:     schema.TypeString,
				Required: true,
			},
			"csvc": {
				Type:     schema.TypeString,
				Required: true,
			},
			"shards": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"source_display": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"accounts": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"alert": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"all_accounts": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"new_account_auto_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"source_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func ltsConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"config_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"log_type_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pipe_alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildCollectConfigQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?region_id=%s", d.Get("region_id").(string))
}

func buildCollectConfigBodyParams(d *schema.ResourceData, enable int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"workspace_id":   d.Get("workspace_id"),
		"dataspace_id":   d.Get("dataspace_id"),
		"dataspace_name": d.Get("dataspace_name"),
		"domain_id":      utils.ValueIgnoreEmpty(d.Get("domain_id")),
		"config":         buildConfigListBodyParams(d.Get("config").([]interface{}), enable),
	}

	if v, ok := d.GetOk("lts_config"); ok {
		bodyParams["lts_config"] = buildLtsConfigBodyParams(v.([]interface{}))
	}

	return bodyParams
}

func buildConfigListBodyParams(configList []interface{}, enable int) []map[string]interface{} {
	if len(configList) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(configList))
	for _, v := range configList {
		raw, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		params := map[string]interface{}{
			"csvc_display":            raw["csvc_display"],
			"csvc":                    raw["csvc"],
			"enable":                  enable,
			"shards":                  raw["shards"],
			"source_display":          raw["source_display"],
			"source_id":               raw["source_id"],
			"ttl":                     raw["ttl"],
			"accounts":                utils.ValueIgnoreEmpty(raw["accounts"]),
			"action":                  utils.ValueIgnoreEmpty(raw["action"]),
			"alert":                   utils.ValueIgnoreEmpty(raw["alert"]),
			"all_accounts":            utils.ValueIgnoreEmpty(raw["all_accounts"]),
			"new_account_auto_access": utils.ValueIgnoreEmpty(raw["new_account_auto_access"]),
			"source_name":             utils.ValueIgnoreEmpty(raw["source_name"]),
		}

		result = append(result, params)
	}

	return result
}

func buildLtsConfigBodyParams(ltsConfigList []interface{}) []map[string]interface{} {
	if len(ltsConfigList) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(ltsConfigList))
	for _, v := range ltsConfigList {
		raw, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		params := map[string]interface{}{
			"config_name":     utils.ValueIgnoreEmpty(raw["config_name"]),
			"description":     utils.ValueIgnoreEmpty(raw["description"]),
			"enable":          utils.ValueIgnoreEmpty(raw["enable"]),
			"log_group_id":    utils.ValueIgnoreEmpty(raw["log_group_id"]),
			"log_stream_id":   utils.ValueIgnoreEmpty(raw["log_stream_id"]),
			"log_type":        utils.ValueIgnoreEmpty(raw["log_type"]),
			"log_type_prefix": utils.ValueIgnoreEmpty(raw["log_type_prefix"]),
			"pipe_alias":      utils.ValueIgnoreEmpty(raw["pipe_alias"]),
		}

		result = append(result, params)
	}

	return result
}

func waitingForCollectConfigApplied(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, sourceId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetCollectConfigInfo(client, d.Get("region_id").(string), sourceId)
			if err != nil {
				return nil, "ERROR", err
			}

			configData := utils.PathSearch("config", respBody, nil)
			if configData == nil {
				return nil, "ERROR", errors.New("unable to find config in API response")
			}

			status := utils.PathSearch("process_status", configData, "").(string)
			if status == "" {
				return nil, "ERROR", errors.New("unable to find process_status in config")
			}

			if status == "SUCCESS" {
				return respBody, "COMPLETED", nil
			}

			// The documentation does not provide the status values for application failures;
			// To be on the safe side, all other statuses will be treated as PENDING.
			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCollectConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/collector/cloudlogs/config"
		sourceId      = fmt.Sprintf("%d", d.Get("config.0.source_id"))
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createPath += buildCollectConfigQueryParams(d)

	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCollectConfigBodyParams(d, 1)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster collect config: %s", err)
	}

	_, err = utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sourceId)

	if err := waitingForCollectConfigApplied(ctx, client, d, d.Timeout(schema.TimeoutCreate), d.Id()); err != nil {
		return diag.Errorf("error waiting for SecMaster collect config applied: %s", err)
	}

	return resourceCollectConfigRead(ctx, d, meta)
}

func resourceCollectConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		sourceId = d.Id()
		regionId = d.Get("region_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	respBody, err := GetCollectConfigInfo(client, regionId, sourceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SecMaster collect config")
	}

	configData := utils.PathSearch("config", respBody, nil)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workspace_id", utils.PathSearch("workspace_id", configData, nil)),
		d.Set("dataspace_id", utils.PathSearch("dataspace_id", configData, nil)),
		d.Set("dataspace_name", utils.PathSearch("dataspace_name", configData, nil)),
		d.Set("region_id", utils.PathSearch("region_id", configData, nil)),
		d.Set("domain_id", utils.PathSearch("domain_id", respBody, nil)),
		d.Set("config", flattenCollectConfig(configData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCollectConfig(configData interface{}) []map[string]interface{} {
	if configData == nil {
		return nil
	}

	cfg := map[string]interface{}{
		"csvc_display":            utils.PathSearch("reference.csvc_display", configData, nil),
		"csvc":                    utils.PathSearch("csvc", configData, nil),
		"shards":                  utils.PathSearch("target.shards", configData, nil),
		"source_display":          utils.PathSearch("reference.source_display", configData, nil),
		"source_id":               utils.PathSearch("source_id", configData, nil),
		"ttl":                     utils.PathSearch("target.ttl", configData, nil),
		"accounts":                utils.PathSearch("accounts", configData, nil),
		"alert":                   utils.PathSearch("alert", configData, nil),
		"all_accounts":            utils.PathSearch("all_accounts", configData, nil),
		"new_account_auto_access": utils.PathSearch("new_account_auto_access", configData, nil),
		"source_name":             utils.PathSearch("source_name", configData, nil),
	}

	return []map[string]interface{}{cfg}
}

func buildGetCollectConfigInfoQueryParams(regionId string, limit, offset int) string {
	queryParams := fmt.Sprintf("?query_statistics=true&region_id=%s&limit=%d&offset=%d",
		regionId, limit, offset)

	return queryParams
}

func GetCollectConfigInfo(client *golangsdk.ServiceClient, regionId, sourceId string) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/collector/cloudlogs/config"
		limit   = 500
		offset  = 0
	)
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	getOpts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
		KeepResponseBody: true,
	}

	for {
		currentRequestPath := requestPath + buildGetCollectConfigInfoQueryParams(regionId, limit, offset)
		resp, err := client.Request("GET", currentRequestPath, &getOpts)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		data := utils.PathSearch(fmt.Sprintf("datasets[?source_id==`%s`]|[0]", sourceId), respBody, nil)
		if data != nil {
			result := make(map[string]interface{})
			for k, v := range respBody.(map[string]interface{}) {
				if k == "datasets" || k == "total" {
					continue
				}
				result[k] = v
			}
			result["config"] = data
			return result, nil
		}

		datasets := utils.PathSearch("datasets", respBody, make([]interface{}, 0)).([]interface{})
		if len(datasets) < limit {
			break
		}

		offset += len(datasets)
	}

	return nil, golangsdk.ErrDefault404{}
}

func resourceCollectConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/collector/cloudlogs/config"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updatePath += buildCollectConfigQueryParams(d)

	updateOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCollectConfigBodyParams(d, 1)),
	}

	_, err = client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster collect config: %s", err)
	}

	return resourceCollectConfigRead(ctx, d, meta)
}

func waitingForCollectConfigDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration, sourceId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetCollectConfigInfo(client, d.Get("region_id").(string), sourceId)
			if err != nil {
				return nil, "ERROR", err
			}

			configData := utils.PathSearch("config", respBody, nil)
			if configData == nil {
				return nil, "ERROR", errors.New("unable to find config in API response")
			}

			status := utils.PathSearch("process_status", configData, "").(string)
			if status == "" {
				return nil, "ERROR", errors.New("unable to find process_status in config")
			}

			if status == "DEFAULT" {
				return respBody, "COMPLETED", nil
			}

			// The documentation does not provide the status values for application failures;
			// To be on the safe side, all other statuses will be treated as PENDING.
			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCollectConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/collector/cloudlogs/config"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath += buildCollectConfigQueryParams(d)

	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCollectConfigBodyParams(d, 0)),
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster collect config: %s", err)
	}

	if err := waitingForCollectConfigDeleted(ctx, client, d, d.Timeout(schema.TimeoutCreate), d.Id()); err != nil {
		return diag.Errorf("error waiting for SecMaster collect config deleted: %s", err)
	}

	return nil
}

func resourceCollectConfigImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<region_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("region_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
