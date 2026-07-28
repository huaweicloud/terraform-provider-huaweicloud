package geminidb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var memoryRuleNonUpdatableParams = []string{
	"dbcache_mapping_id",
	"name",
	"source_db_schema",
	"source_db_table",
	"storage_type",
	"target_database",
	"key_prefix",
	"key_columns",
	"key_separator",
}

// @API GeminiDB POST /v3/{project_id}/dbcache/rule
// @API GeminiDB GET /v3/{project_id}/dbcache/rules
// @API GeminiDB PUT /v3/{project_id}/dbcache/rule
// @API GeminiDB DELETE /v3/{project_id}/dbcache/rule
func ResourceMemoryRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemoryRuleCreate,
		ReadContext:   resourceMemoryRuleRead,
		UpdateContext: resourceMemoryRuleUpdate,
		DeleteContext: resourceMemoryRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMemoryRuleImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(memoryRuleNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dbcache_mapping_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_db_schema": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_db_table": {
				Type:     schema.TypeString,
				Required: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_columns": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"value_columns": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"key_separator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value_separator": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ttl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateMemoryRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"dbcache_mapping_id": d.Get("dbcache_mapping_id"),
		"name":               d.Get("name"),
		"source_db_schema":   d.Get("source_db_schema"),
		"source_db_table":    d.Get("source_db_table"),
		"storage_type":       d.Get("storage_type"),
		"target_database":    d.Get("target_database"),
		"key_prefix":         d.Get("key_prefix"),
		"key_columns":        utils.ExpandToStringList(d.Get("key_columns").([]interface{})),
		"value_columns":      utils.ExpandToStringList(d.Get("value_columns").([]interface{})),
		"key_separator":      d.Get("key_separator"),
		"value_separator":    utils.ValueIgnoreEmpty(d.Get("value_separator")),
		"ttl":                utils.ValueIgnoreEmpty(d.Get("ttl")),
	}

	return bodyParams
}

func resourceMemoryRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v3/{project_id}/dbcache/rule"
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateMemoryRuleBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating memory acceleration rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("error creating memory acceleration mapping: unable to find rule ID from API response")
	}

	d.SetId(ruleId)

	return resourceMemoryRuleRead(ctx, d, meta)
}

func resourceMemoryRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		mappingId = d.Get("dbcache_mapping_id").(string)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	memoryRuleInfo, err := GetMemoryRuleInfo(client, mappingId, d.Id())
	if err != nil {
		// When the memory rule does not exist, the response HTTP status code of the query API is 200
		// and return nil
		return common.CheckDeletedDiag(d, err, "error retrieving memory acceleration rule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("dbcache_mapping_id", mappingId),
		d.Set("name", utils.PathSearch("name", memoryRuleInfo, nil)),
		d.Set("status", utils.PathSearch("status", memoryRuleInfo, nil)),
		d.Set("source_db_schema", utils.PathSearch("source_db_schema", memoryRuleInfo, nil)),
		d.Set("source_db_table", utils.PathSearch("source_db_table", memoryRuleInfo, nil)),
		d.Set("storage_type", utils.PathSearch("storage_type", memoryRuleInfo, nil)),
		d.Set("target_database", utils.PathSearch("target_database", memoryRuleInfo, nil)),
		d.Set("key_prefix", utils.PathSearch("key_prefix", memoryRuleInfo, nil)),
		d.Set("key_columns", utils.PathSearch("key_columns", memoryRuleInfo, nil)),
		d.Set("value_columns", utils.PathSearch("value_columns", memoryRuleInfo, nil)),
		d.Set("key_separator", utils.PathSearch("key_separator", memoryRuleInfo, nil)),
		d.Set("value_separator", utils.PathSearch("value_separator", memoryRuleInfo, nil)),
		d.Set("ttl", utils.PathSearch("ttl", memoryRuleInfo, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetMemoryRuleInfo(client *golangsdk.ServiceClient, mappingId, ruleId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/dbcache/rules?dbcache_mapping_id={dbcache_mapping_id}&rule_id={rule_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{dbcache_mapping_id}", mappingId)
	getPath = strings.ReplaceAll(getPath, "{rule_id}", ruleId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	ruleInfo := utils.PathSearch("rules|[0]", respBody, nil)
	if ruleInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return ruleInfo, nil
}

func buildUpdateMemoryRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"dbcache_rule_id": d.Id(),
		"value_columns":   utils.ExpandToStringList(d.Get("value_columns").([]interface{})),
		"value_separator": d.Get("value_separator"),
		"ttl":             utils.ValueIgnoreEmpty(d.Get("ttl")),
	}

	return bodyParams
}

func resourceMemoryRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	if d.HasChangeExcept("enable_force_new") {
		err = updateMemoryRule(client, d)
		if err != nil {
			return diag.Errorf("error updating memory acceleration rule: %s", err)
		}
	}

	return resourceMemoryRuleRead(ctx, d, meta)
}

func updateMemoryRule(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v3/{project_id}/dbcache/rule"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
		JSONBody:    utils.RemoveNil(buildUpdateMemoryRuleBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &opt)

	return err
}

func resourceMemoryRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/dbcache/rule"
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"id": d.Id(),
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// whether the memory rule exist or not, the response HTTP status code of the deletion API is 200.
		return diag.Errorf("error deleting memory acceleration rule: %s", err)
	}

	return nil
}

func resourceMemoryRuleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<dbcache_mapping_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("dbcache_mapping_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
