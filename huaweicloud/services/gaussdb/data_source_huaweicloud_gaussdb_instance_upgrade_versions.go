package gaussdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB POST /gaussdb/v3.1/{project_id}/instances/db-upgrade/candidate-versions
func DataSourceGaussDBInstanceUpgradeVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceGaussDBKernelVersionUpgradeRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"upgrade_type_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDBKernelVersionUpgradeUpgradeTypeListSchema(),
			},
			"target_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"upgrade_candidate_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hotfix_upgrade_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDBKernelVersionUpgradeHotfixInfosSchema(),
			},
			"hotfix_rollback_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDBKernelVersionUpgradeHotfixInfosSchema(),
			},
		},
	}
}

func gaussDBKernelVersionUpgradeUpgradeTypeListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"upgrade_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"upgrade_action_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDBKernelVersionUpgradeUpgradeActionListSchema(),
			},
			"is_parallel_upgrade": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func gaussDBKernelVersionUpgradeUpgradeActionListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"upgrade_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func gaussDBKernelVersionUpgradeHotfixInfosSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"common_patch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_sensitive": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"descripition": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_upgrade": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceGaussDBKernelVersionUpgradeRead(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3.1/{project_id}/instances/db-upgrade/candidate-versions"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"instance_ids": utils.ExpandToStringList(d.Get("instance_ids").([]interface{})),
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error querying GaussDB kernel version upgrade candidate versions: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("upgrade_type_list", flattenGaussDBKernelVersionUpgradeUpgradeTypeList(getRespBody)),
		d.Set("target_version", utils.PathSearch("target_version", getRespBody, nil)),
		d.Set("upgrade_candidate_versions", utils.PathSearch("upgrade_candidate_versions", getRespBody, make([]interface{}, 0))),
		d.Set("hotfix_upgrade_infos", flattenGaussDBKernelVersionUpgradeHotfixInfos(getRespBody, "hotfix_upgrade_infos")),
		d.Set("hotfix_rollback_infos", flattenGaussDBKernelVersionUpgradeHotfixInfos(getRespBody, "hotfix_rollback_infos")),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGaussDBKernelVersionUpgradeUpgradeTypeList(resp interface{}) []interface{} {
	curJson := utils.PathSearch("upgrade_type_list", resp, make([]interface{}, 0))
	curArray, ok := curJson.([]interface{})
	if !ok {
		return nil
	}

	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"upgrade_type":        utils.PathSearch("upgrade_type", v, nil),
			"enable":              utils.PathSearch("enable", v, nil),
			"upgrade_action_list": flattenGaussDBKernelVersionUpgradeUpgradeActionList(v),
			"is_parallel_upgrade": utils.PathSearch("is_parallel_upgrade", v, nil),
		})
	}
	return res
}

func flattenGaussDBKernelVersionUpgradeUpgradeActionList(resp interface{}) []interface{} {
	curJson := utils.PathSearch("upgrade_action_list", resp, make([]interface{}, 0))
	curArray, ok := curJson.([]interface{})
	if !ok {
		return nil
	}

	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"upgrade_action": utils.PathSearch("upgrade_action", v, nil),
			"enable":         utils.PathSearch("enable", v, nil),
		})
	}
	return res
}

func flattenGaussDBKernelVersionUpgradeHotfixInfos(resp interface{}, path string) []interface{} {
	curJson := utils.PathSearch(path, resp, make([]interface{}, 0))
	curArray, ok := curJson.([]interface{})
	if !ok {
		return nil
	}

	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"version":          utils.PathSearch("version", v, nil),
			"common_patch":     utils.PathSearch("common_patch", v, nil),
			"backup_sensitive": utils.PathSearch("backup_sensitive", v, nil),
			"descripition":     utils.PathSearch("descripition", v, nil),
			"default_upgrade":  utils.PathSearch("default_upgrade", v, nil),
		})
	}
	return res
}
