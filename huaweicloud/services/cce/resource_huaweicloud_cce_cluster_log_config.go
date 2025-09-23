package cce

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var logConfigNonUpdatableParams = []string{"cluster_id"}

// @API CCE PUT /api/v3/projects/{project_id}/cluster/{cluster_id}/log-configs
// @API CCE GET /api/v3/projects/{project_id}/cluster/{cluster_id}/log-configs
func ResourceClusterLogConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterLogConfigCreateOrUpdate,
		UpdateContext: resourceClusterLogConfigCreateOrUpdate,
		ReadContext:   resourceClusterLogConfigRead,
		DeleteContext: resourceClusterLogConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceClusterLogConfigImport,
		},

		CustomizeDiff: config.FlexibleForceNew(logConfigNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"log_configs": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func buildClusterLogConfigCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ttl_in_days": utils.ValueIgnoreEmpty(d.Get("ttl_in_days")),
		"log_configs": buildLogConfigsCreateBodyParams(d),
	}
	return bodyParams
}

func buildLogConfigsCreateBodyParams(d *schema.ResourceData) []map[string]interface{} {
	logConfigsRaw := d.Get("log_configs").(*schema.Set).List()
	if len(logConfigsRaw) == 0 {
		return nil
	}

	bodyParams := make([]map[string]interface{}, 0, len(logConfigsRaw))
	for _, v := range logConfigsRaw {
		logConfig, ok := v.(map[string]interface{})
		if ok {
			bodyParams = append(bodyParams, map[string]interface{}{
				"name":   utils.ValueIgnoreEmpty(logConfig["name"]),
				"enable": utils.ValueIgnoreEmpty(logConfig["enable"]),
			})
		}
	}
	return bodyParams
}

func resourceClusterLogConfigCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		createClusterLogConfigHttpUrl = "api/v3/projects/{project_id}/cluster/{cluster_id}/log-configs"
		createClusterLogConfigProduct = "cce"
	)

	clusterID := d.Get("cluster_id").(string)
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient(createClusterLogConfigProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	createClusterLogConfigPath := client.Endpoint + createClusterLogConfigHttpUrl
	createClusterLogConfigPath = strings.ReplaceAll(createClusterLogConfigPath, "{project_id}", client.ProjectID)
	createClusterLogConfigPath = strings.ReplaceAll(createClusterLogConfigPath, "{cluster_id}", clusterID)

	createClusterLogConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createClusterLogConfigOpt.JSONBody = utils.RemoveNil(buildClusterLogConfigCreateBodyParams(d))
	_, err = client.Request("PUT", createClusterLogConfigPath, &createClusterLogConfigOpt)
	if err != nil {
		return diag.Errorf("error updating CCE cluster log config: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(clusterID)
	}
	return resourceClusterLogConfigRead(ctx, d, meta)
}

func resourceClusterLogConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		getClusterLogConfigHttpUrl = "api/v3/projects/{project_id}/cluster/{cluster_id}/log-configs"
		getClusterLogConfigProduct = "cce"
	)

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient(getClusterLogConfigProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	getClusterLogConfigPath := client.Endpoint + getClusterLogConfigHttpUrl
	getClusterLogConfigPath = strings.ReplaceAll(getClusterLogConfigPath, "{project_id}", client.ProjectID)
	getClusterLogConfigPath = strings.ReplaceAll(getClusterLogConfigPath, "{cluster_id}", d.Id())

	getClusterLogConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getClusterLogConfigResp, err := client.Request("GET", getClusterLogConfigPath, &getClusterLogConfigOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving log config")
	}

	getClusterLogConfigRespBody, err := utils.FlattenResponse(getClusterLogConfigResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("ttl_in_days", utils.PathSearch("ttl_in_days", getClusterLogConfigRespBody, nil)),
		d.Set("log_configs", flattenLogConfigs(d, getClusterLogConfigRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLogConfigs(d *schema.ResourceData, getClusterLogConfigRespBody interface{}) []map[string]interface{} {
	logConfigsRaw := utils.PathSearch("log_configs", getClusterLogConfigRespBody, nil)
	if logConfigsRaw == nil {
		return nil
	}

	originalLogConfigs := buildLogConfigsCreateBodyParams(d)
	originalNamesRaw := utils.PathSearch("[*].name", originalLogConfigs, []interface{}{}).([]interface{})
	originalNames := utils.ExpandToStringList(originalNamesRaw)
	logConfigs := logConfigsRaw.([]interface{})

	res := make([]map[string]interface{}, 0, len(logConfigs))
	if len(originalNames) != 0 {
		for _, v := range logConfigs {
			logConfig := v.(map[string]interface{})
			if utils.StrSliceContains(originalNames, logConfig["name"].(string)) {
				res = append(res, map[string]interface{}{
					"name":   logConfig["name"],
					"enable": logConfig["enable"],
				})
			}
		}
	} else {
		for _, v := range logConfigs {
			logConfig := v.(map[string]interface{})
			res = append(res, map[string]interface{}{
				"name":   logConfig["name"],
				"enable": logConfig["enable"],
			})
		}
	}

	return res
}

func buildClusterLogConfigDeleteBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ttl_in_days": 0,
		"log_configs": buildLogConfigsDeleteBodyParams(d),
	}
	return bodyParams
}

func buildLogConfigsDeleteBodyParams(d *schema.ResourceData) []map[string]interface{} {
	logConfigsOld, _ := d.GetChange("log_configs")
	logConfigsRaw := logConfigsOld.(*schema.Set).List()
	if len(logConfigsRaw) == 0 {
		return nil
	}

	bodyParams := make([]map[string]interface{}, 0, len(logConfigsRaw))
	for _, v := range logConfigsRaw {
		logConfig, ok := v.(map[string]interface{})
		if ok {
			bodyParams = append(bodyParams, map[string]interface{}{
				"name":   utils.ValueIgnoreEmpty(logConfig["name"]),
				"enable": false,
			})
		}
	}
	return bodyParams
}

func resourceClusterLogConfigDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		deleteClusterLogConfigHttpUrl = "api/v3/projects/{project_id}/cluster/{cluster_id}/log-configs"
		deleteClusterLogConfigProduct = "cce"
	)

	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient(deleteClusterLogConfigProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	deleteClusterLogConfigPath := client.Endpoint + deleteClusterLogConfigHttpUrl
	deleteClusterLogConfigPath = strings.ReplaceAll(deleteClusterLogConfigPath, "{project_id}", client.ProjectID)
	deleteClusterLogConfigPath = strings.ReplaceAll(deleteClusterLogConfigPath, "{cluster_id}", d.Id())

	deleteClusterLogConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteClusterLogConfigOpt.JSONBody = utils.RemoveNil(buildClusterLogConfigDeleteBodyParams(d))
	_, err = client.Request("PUT", deleteClusterLogConfigPath, &deleteClusterLogConfigOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CCE cluster log config")
	}

	return nil
}

func resourceClusterLogConfigImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	err := d.Set("cluster_id", d.Id())
	return []*schema.ResourceData{d}, err
}
