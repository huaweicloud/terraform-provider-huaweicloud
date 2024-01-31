package waf

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF GET /v1/{project_id}/waf/ip-groups
func DataSourceWafAddressGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceAddressGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the data source.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the address group.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the IP address or IP address ranges.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Elem:        groupsSchema(),
				Computed:    true,
				Description: `Specifies the list of address group.`,
			},
		},
	}
}

func groupsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the ID of the address group.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the name of the address group.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"ip_addresses": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the IP addresses or IP address ranges.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Elem:        rulesInfoSchema(),
				Computed:    true,
				Description: `Specifies the list of rules that use the IP address group.`,
			},
			"share_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the number of the users share the address group.`,
			},
			"accept_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the number of users accept the address group.`,
			},
			"process_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the status of the processing.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the description of the address group.`,
			},
		},
	}
	return &sc
}

func rulesInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of rule.`,
			},
			"rule_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of rule.`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of policy.`,
			},
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of policy.`,
			},
		},
	}
	return &sc
}

func datasourceAddressGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getWAFAddressGroup: Query WAF address group
	var (
		getWAFAddressGroupHttpUrl = "v1/{project_id}/waf/ip-groups"
		getWAFAddressGroupProduct = "waf"
	)
	getWAFAddressGroupClient, err := cfg.NewServiceClient(getWAFAddressGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	getWAFAddressGroupsPath := getWAFAddressGroupClient.Endpoint + getWAFAddressGroupHttpUrl
	getWAFAddressGroupsPath = strings.ReplaceAll(getWAFAddressGroupsPath, "{project_id}",
		getWAFAddressGroupClient.ProjectID)
	getWAFAddressGroupsPath += buildWAFAddressGroupsQueryParams(d, cfg)

	getWAFAddressGroupsResp, err := pagination.ListAllItems(
		getWAFAddressGroupClient,
		"page",
		getWAFAddressGroupsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving address groups, %s", err)
	}

	listWAFAddressGroupsRespJson, err := json.Marshal(getWAFAddressGroupsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listWAFAddressGroupsRespBody interface{}
	err = json.Unmarshal(listWAFAddressGroupsRespJson, &listWAFAddressGroupsRespBody)
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
		d.Set("groups", flattenListAddressGroupsBody(listWAFAddressGroupsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListAddressGroupsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("items", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"ip_addresses":          utils.PathSearch("ips", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"share_count":           utils.PathSearch("share_info.share_count", v, nil),
			"accept_count":          utils.PathSearch("share_info.accept_count", v, nil),
			"process_status":        utils.PathSearch("share_info.process_status", v, nil),
			"rules":                 flattenAddressGroupResponseBodyRules(resp),
		})
	}
	return rst
}

func buildWAFAddressGroupsQueryParams(d *schema.ResourceData, conf *config.Config) string {
	res := ""
	epsId := conf.GetEnterpriseProjectID(d)
	if epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("ip_address"); ok {
		res = fmt.Sprintf("%s&ip=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
