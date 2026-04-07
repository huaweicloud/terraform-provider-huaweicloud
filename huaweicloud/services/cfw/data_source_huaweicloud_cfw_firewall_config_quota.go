package cfw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CFW GET /v2/{project_id}/cfw/{fw_instance_id}/quota
func DataSourceFirewallConfigQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFirewallConfigQuotaRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"set_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: firewallConfigQuotaDataSchema(),
				},
			},
		},
	}
}

func firewallConfigQuotaDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"item_info": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"max_quota": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"used_quota": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"extras_info": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
		"max_quota": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"quota_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"used_quota": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func buildFirewallConfigQuotaQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?config_type=%s", d.Get("config_type"))
	if v, ok := d.GetOk("set_id"); ok {
		res = fmt.Sprintf("%s&set_id=%s", res, v.(string))
	}

	return res
}

func dataSourceFirewallConfigQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v2/{project_id}/cfw/{fw_instance_id}/quota"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{fw_instance_id}", d.Get("fw_instance_id").(string))
	requestPath += buildFirewallConfigQuotaQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW firewall config quota: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenShowConfigQuotaData(utils.PathSearch("data", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenShowConfigQuotaData(dataResp interface{}) []interface{} {
	if dataResp == nil {
		return nil
	}

	itemInfoRaw := utils.PathSearch("item_info", dataResp, nil)
	var itemInfoList []interface{}
	if itemInfoRaw != nil {
		itemInfoList = []interface{}{map[string]interface{}{
			"max_quota":  utils.PathSearch("max_quota", itemInfoRaw, nil),
			"used_quota": utils.PathSearch("used_quota", itemInfoRaw, nil),
			"extras_info": utils.ExpandToStringMap(
				utils.PathSearch("extras_info", itemInfoRaw, make(map[string]interface{})).(map[string]interface{})),
		}}
	}

	return []interface{}{map[string]interface{}{
		"item_info":  itemInfoList,
		"max_quota":  utils.PathSearch("max_quota", dataResp, nil),
		"quota_type": utils.PathSearch("quota_type", dataResp, nil),
		"used_quota": utils.PathSearch("used_quota", dataResp, nil),
	}}
}
