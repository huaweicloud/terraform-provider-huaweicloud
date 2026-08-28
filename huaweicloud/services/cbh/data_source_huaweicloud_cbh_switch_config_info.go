package cbh

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBH GET /v2/{project_id}/cbs/feature/config
func DataSourceSwitchConfigInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwitchConfigInfoRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"switch_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_support_unibuy": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_float_ipv6": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_admin_login": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_update_ha": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_tms": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_eps": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_iam_login": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_ipv6": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_ha": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_reset": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_upgrade_instance": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_change_security_group": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_manually_ip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_capacity_expantion": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_ha_expantion": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_agency_authorize": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_change_vpc": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_cluster": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_ondemand": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_period": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"version_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"require_eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"iam_login": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"admin_login": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"float_ipv6": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSwitchConfigInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		mErr       *multierror.Error
		getHttpUrl = "v2/{project_id}/cbs/feature/config"
		product    = "cbh"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving CBH switch config info: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("switch_info", flattenSwitchInfo(getRespBody)),
		d.Set("version_info", flattenVersionInfo(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSwitchInfo(resp interface{}) []interface{} {
	switchInfo := utils.PathSearch("switch_info", resp, nil)
	if switchInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"is_support_unibuy":                utils.PathSearch("is_support_unibuy", switchInfo, nil),
			"is_support_float_ipv6":            utils.PathSearch("is_support_float_ipv6", switchInfo, nil),
			"is_support_admin_login":           utils.PathSearch("is_support_admin_login", switchInfo, nil),
			"is_support_update_ha":             utils.PathSearch("is_support_update_ha", switchInfo, nil),
			"is_support_tms":                   utils.PathSearch("is_support_tms", switchInfo, nil),
			"is_support_eps":                   utils.PathSearch("is_support_eps", switchInfo, nil),
			"is_support_iam_login":             utils.PathSearch("is_support_iam_login", switchInfo, nil),
			"is_support_ipv6":                  utils.PathSearch("is_support_ipv6", switchInfo, nil),
			"is_support_ha":                    utils.PathSearch("is_support_ha", switchInfo, nil),
			"is_support_reset":                 utils.PathSearch("is_support_reset", switchInfo, nil),
			"is_support_upgrade_instance":      utils.PathSearch("is_support_upgrade_instance", switchInfo, nil),
			"is_support_change_security_group": utils.PathSearch("is_support_change_security_group", switchInfo, nil),
			"is_support_manually_ip":           utils.PathSearch("is_support_manually_ip", switchInfo, nil),
			"is_support_capacity_expantion":    utils.PathSearch("is_support_capacity_expantion", switchInfo, nil),
			"is_support_ha_expantion":          utils.PathSearch("is_support_ha_expantion", switchInfo, nil),
			"is_support_agency_authorize":      utils.PathSearch("is_support_agency_authorize", switchInfo, nil),
			"is_support_change_vpc":            utils.PathSearch("is_support_change_vpc", switchInfo, nil),
			"is_support_cluster":               utils.PathSearch("is_support_cluster", switchInfo, nil),
			"is_support_ondemand":              utils.PathSearch("is_support_ondemand", switchInfo, nil),
			"is_support_period":                utils.PathSearch("is_support_period", switchInfo, nil),
		},
	}
}

func flattenVersionInfo(resp interface{}) []interface{} {
	versionInfo := utils.PathSearch("version_info", resp, nil)
	if versionInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"require_eip": utils.PathSearch("require_eip", versionInfo, nil),
			"iam_login":   utils.PathSearch("iam_login", versionInfo, nil),
			"admin_login": utils.PathSearch("admin_login", versionInfo, nil),
			"float_ipv6":  utils.PathSearch("float_ipv6", versionInfo, nil),
		},
	}
}
