package evs

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EVS GET /v3/{project_id}/os-quota-sets/{target_project_id}
func DataSourceEvsV3Quotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvsV3QuotasRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"usage": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"quota_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     V3QuotaSetSchema(),
			},
		},
	}
}

func V3QuotaSetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"backup_gigabytes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"backups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"gigabytes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"gigabytes_sata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"snapshots_sata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"volumes_sata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"gigabytes_sas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"snapshots_sas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"volumes_sas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"gigabytes_ssd": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"snapshots_ssd": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"volumes_ssd": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"gigabytes_gpssd": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"snapshots_gpssd": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"volumes_gpssd": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
			"per_volume_gigabytes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotaSetSubSchema(),
			},
		},
	}
}

func quotaSetSubSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"in_use": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceEvsV3QuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/os-quota-sets/{target_project_id}"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{target_project_id}", client.ProjectID)
	requestPath += "?usage=" + strconv.FormatBool(d.Get("usage").(bool))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error querying EVS v3 quota: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	quota := utils.PathSearch("quota_set", respBody, nil)

	quotaSet := map[string]interface{}{
		"backup_gigabytes":     flattenQuotaParams(utils.PathSearch("backup_gigabytes", quota, nil)),
		"backups":              flattenQuotaParams(utils.PathSearch("backups", quota, nil)),
		"gigabytes":            flattenQuotaParams(utils.PathSearch("gigabytes", quota, nil)),
		"snapshots":            flattenQuotaParams(utils.PathSearch("snapshots", quota, nil)),
		"volumes":              flattenQuotaParams(utils.PathSearch("volumes", quota, nil)),
		"gigabytes_sata":       flattenQuotaParams(utils.PathSearch("gigabytes_SATA", quota, nil)),
		"snapshots_sata":       flattenQuotaParams(utils.PathSearch("snapshots_SATA", quota, nil)),
		"volumes_sata":         flattenQuotaParams(utils.PathSearch("volumes_SATA", quota, nil)),
		"gigabytes_sas":        flattenQuotaParams(utils.PathSearch("gigabytes_SAS", quota, nil)),
		"snapshots_sas":        flattenQuotaParams(utils.PathSearch("snapshots_SAS", quota, nil)),
		"volumes_sas":          flattenQuotaParams(utils.PathSearch("volumes_SAS", quota, nil)),
		"gigabytes_ssd":        flattenQuotaParams(utils.PathSearch("gigabytes_SSD", quota, nil)),
		"snapshots_ssd":        flattenQuotaParams(utils.PathSearch("snapshots_SSD", quota, nil)),
		"volumes_ssd":          flattenQuotaParams(utils.PathSearch("volumes_SSD", quota, nil)),
		"gigabytes_gpssd":      flattenQuotaParams(utils.PathSearch("gigabytes_GPSSD", quota, nil)),
		"snapshots_gpssd":      flattenQuotaParams(utils.PathSearch("snapshots_GPSSD", quota, nil)),
		"volumes_gpssd":        flattenQuotaParams(utils.PathSearch("volumes_GPSSD", quota, nil)),
		"per_volume_gigabytes": flattenQuotaParams(utils.PathSearch("per_volume_gigabytes", quota, nil)),
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("quota_set", []interface{}{quotaSet}),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving data source fields of the EVS v3 quota: %s", mErr)
	}
	return nil
}

func flattenQuotaParams(params interface{}) []interface{} {
	if params == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"limit":  utils.PathSearch("limit", params, nil),
			"in_use": utils.PathSearch("in_use", params, nil),
		},
	}
}
