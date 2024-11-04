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
				Description: `The ID of the address group.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the address group.`,
			},
			"ip_addresses": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP addresses or IP address ranges.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Elem:        rulesInfoSchema(),
				Computed:    true,
				Description: `The list of rules that use the IP address group.`,
			},
			"share_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of the users share the address group.`,
			},
			"accept_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of users accept the address group.`,
			},
			"process_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The status of the processing.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the address group.`,
			},

			// Deprecated
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `schema: Deprecated; The enterprise project ID.`,
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
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/waf/ip-groups"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildWAFAddressGroupsQueryParams(d, cfg)

	resp, err := pagination.ListAllItems(
		client,
		"page",
		requestPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving WAF address groups, %s", err)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
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
		d.Set("groups", flattenListAddressGroupsBody(respBody)),
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
