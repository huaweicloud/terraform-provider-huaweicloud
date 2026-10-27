package drs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var drsJobUpdateNonUpdatableParams = []string{
	"job_id",
	"type",
}

// @API DRS PUT /v5/{project_id}/jobs/{job_id}
func ResourceDrsJobUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrsJobUpdateCreate,
		ReadContext:   resourceDrsJobUpdateRead,
		UpdateContext: resourceDrsJobUpdateUpdate,
		DeleteContext: resourceDrsJobUpdateDelete,

		CustomizeDiff: config.FlexibleForceNew(drsJobUpdateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"endpoint", "description", "name", "notify", "expired_days",
					"policy", "speed_limit", "tuning_params", "db_object",
					"user_migration", "precheck", "db_param", "db_password",
					"re_create", "all", "resource_group", "application_id",
					"dr_info", "replay_config",
				}, false),
			},
			"policy_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sync_type_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"conflict_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"filter_ddl_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ddl_trans": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"index_trans": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"struct_trans": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"data_sync_topology_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"support_ddl_info": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"full_mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"export_snapshot": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"replace_definer": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"increment_read_mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"slot_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"transformation_name_case_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_create_table_with_index": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"kafka_data_format": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kafka_compression_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"partition_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"topic_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"topic_name_format": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"topic": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"partitions_num": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"replication_factor": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"filter_user_names": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"base_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"speed_limit": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"speed": {
							Type:     schema.TypeString,
							Required: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"is_effect": {
							Type:     schema.TypeBool,
							Optional: true,
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

func buildDrsJobUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	updateType := d.Get("type").(string)
	params := map[string]interface{}{}

	switch updateType {
	case "policy":
		if v, ok := d.GetOk("policy_config"); ok {
			params["policy_config"] = buildDrsJobUpdatePolicyConfig(v.([]interface{}))
		}
	case "name", "description":
		if v, ok := d.GetOk("base_info"); ok {
			params["base_info"] = buildDrsJobUpdateBaseInfo(v.([]interface{}))
		}
	case "speed_limit":
		if v, ok := d.GetOk("speed_limit"); ok {
			params["speed_limit"] = buildDrsJobUpdateSpeedLimit(v.([]interface{}))
		}
	}

	return map[string]interface{}{
		"job": map[string]interface{}{
			"type":   updateType,
			"params": params,
		},
	}
}

func buildDrsJobUpdatePolicyConfig(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}
	raw, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}
	return map[string]interface{}{
		"sync_type_policy":                utils.ValueIgnoreEmpty(raw["sync_type_policy"]),
		"conflict_policy":                 utils.ValueIgnoreEmpty(raw["conflict_policy"]),
		"filter_ddl_policy":               utils.ValueIgnoreEmpty(raw["filter_ddl_policy"]),
		"ddl_trans":                       utils.ValueIgnoreEmpty(raw["ddl_trans"]),
		"index_trans":                     utils.ValueIgnoreEmpty(raw["index_trans"]),
		"struct_trans":                    utils.ValueIgnoreEmpty(raw["struct_trans"]),
		"data_sync_topology_type":         utils.ValueIgnoreEmpty(raw["data_sync_topology_type"]),
		"support_ddl_info":                utils.ValueIgnoreEmpty(raw["support_ddl_info"]),
		"full_mode":                       utils.ValueIgnoreEmpty(raw["full_mode"]),
		"export_snapshot":                 utils.ValueIgnoreEmpty(raw["export_snapshot"]),
		"replace_definer":                 utils.ValueIgnoreEmpty(raw["replace_definer"]),
		"increment_read_mode":             utils.ValueIgnoreEmpty(raw["increment_read_mode"]),
		"slot_name":                       utils.ValueIgnoreEmpty(raw["slot_name"]),
		"transformation_name_case_policy": utils.ValueIgnoreEmpty(raw["transformation_name_case_policy"]),
		"is_create_table_with_index":      utils.ValueIgnoreEmpty(raw["is_create_table_with_index"]),
		"kafka_data_format":               utils.ValueIgnoreEmpty(raw["kafka_data_format"]),
		"kafka_compression_type":          utils.ValueIgnoreEmpty(raw["kafka_compression_type"]),
		"partition_policy":                utils.ValueIgnoreEmpty(raw["partition_policy"]),
		"topic_policy":                    utils.ValueIgnoreEmpty(raw["topic_policy"]),
		"topic_name_format":               utils.ValueIgnoreEmpty(raw["topic_name_format"]),
		"topic":                           utils.ValueIgnoreEmpty(raw["topic"]),
		"partitions_num":                  utils.ValueIgnoreEmpty(raw["partitions_num"]),
		"replication_factor":              utils.ValueIgnoreEmpty(raw["replication_factor"]),
		"filter_user_names":               utils.ValueIgnoreEmpty(raw["filter_user_names"]),
	}
}

func buildDrsJobUpdateBaseInfo(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}
	raw, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}
	return map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(raw["name"]),
		"description": utils.ValueIgnoreEmpty(raw["description"]),
	}
}

func buildDrsJobUpdateSpeedLimit(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		raw, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"speed":      raw["speed"],
			"start_time": raw["start_time"],
			"end_time":   raw["end_time"],
			"is_effect":  utils.ValueIgnoreEmpty(raw["is_effect"]),
		})
	}
	return rst
}

func resourceDrsJobUpdateCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		jobId   = d.Get("job_id").(string)
		httpUrl = "v5/{project_id}/jobs/{job_id}"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		OkCodes:          []int{200, 202},
		JSONBody:         utils.RemoveNil(buildDrsJobUpdateBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating DRS job: %s", err)
	}

	d.SetId(jobId)

	return nil
}

func resourceDrsJobUpdateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceDrsJobUpdateUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		jobId   = d.Get("job_id").(string)
		httpUrl = "v5/{project_id}/jobs/{job_id}"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		OkCodes:          []int{200, 202},
		JSONBody:         utils.RemoveNil(buildDrsJobUpdateBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating DRS job: %s", err)
	}

	return nil
}

func resourceDrsJobUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to update DRS job details. Deleting this resource
will not rollback the update operation, but will only remove the resource information from the tf
state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
