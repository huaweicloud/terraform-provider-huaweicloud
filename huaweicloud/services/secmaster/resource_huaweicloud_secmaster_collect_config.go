package secmaster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var collectConfigNonUpdatableParams = []string{
	"workspace_id",
	"dataspace_id",
	"dataspace_name",
	"region_id",
	"domain_id",
	"config",
	"config.*.csvc_display",
	"config.*.csvc",
	"config.*.shards",
	"config.*.source_display",
	"config.*.source_id",
	"config.*.ttl",
	"config.*.enable",
	"config.*.accounts",
	"config.*.action",
	"config.*.alert",
	"config.*.all_accounts",
	"config.*.new_account_auto_access",
	"config.*.source_name",
	"lts_config",
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
func ResourceCollectConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCollectConfigCreate,
		ReadContext:   resourceCollectConfigRead,
		UpdateContext: resourceCollectConfigUpdate,
		DeleteContext: resourceCollectConfigDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
			"enable": {
				Type:     schema.TypeInt,
				Optional: true,
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

func buildCollectConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"workspace_id":   d.Get("workspace_id"),
		"dataspace_id":   d.Get("dataspace_id"),
		"dataspace_name": d.Get("dataspace_name"),
		"domain_id":      utils.ValueIgnoreEmpty(d.Get("domain_id")),
		"config":         buildConfigListBodyParams(d.Get("config").([]interface{})),
	}

	if v, ok := d.GetOk("lts_config"); ok {
		bodyParams["lts_config"] = buildLtsConfigBodyParams(v.([]interface{}))
	}

	return bodyParams
}

func buildConfigListBodyParams(configList []interface{}) []map[string]interface{} {
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
			"enable":                  raw["enable"],
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

func resourceCollectConfigCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/collector/cloudlogs/config"
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
		JSONBody:         utils.RemoveNil(buildCollectConfigBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster collect config: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	return nil
}

func resourceCollectConfigRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCollectConfigUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCollectConfigDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to set SecMaster collect config. 
Deleting this resource will not undo the disable action or restore the collect config, but will only 
remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
