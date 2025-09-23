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

// @API WAF GET /v1/{project_id}/waf/instance
func DataSourceWafDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceDomainsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Elem:     wafDomainSchema(),
				Computed: true,
			},
		},
	}
}

func wafDomainSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protect_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pci_3ds": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"pci_dss": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"access_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"access_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"website_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proxy_layer": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func datasourceDomainsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/waf/instance"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildWAFDomainsQueryParams(d, cfg)

	resp, err := pagination.ListAllItems(
		client,
		"page",
		requestPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving WAF domains, %s", err)
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
		d.Set("domains", flattenListDomainsBody(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListDomainsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("items", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		createTime := utils.PathSearch("create_time", v, float64(0))
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"proxy":                 utils.StringToBool(utils.PathSearch("proxcy", v, "")),
			"domain":                utils.PathSearch("hostname", v, nil),
			"policy_id":             utils.PathSearch("policyid", v, nil),
			"protect_status":        utils.PathSearch("protect_status", v, nil),
			"pci_3ds":               utils.StringToBool(utils.PathSearch("flag.pci_3ds", v, "")),
			"pci_dss":               utils.StringToBool(utils.PathSearch("flag.pci_dss", v, "")),
			"ipv6_enable":           utils.StringToBool(utils.PathSearch("flag.ipv6", v, "")),
			"access_status":         utils.PathSearch("access_status", v, nil),
			"access_code":           utils.PathSearch("access_code", v, nil),
			"charging_mode":         utils.PathSearch("paid_type", v, nil),
			"website_name":          utils.PathSearch("web_tag", v, nil),
			"proxy_layer":           utils.PathSearch("proxy_layer", v, nil),
			"created_at":            utils.FormatTimeStampRFC3339(int64(createTime.(float64))/1000, false),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
		})
	}
	return rst
}

func buildWAFDomainsQueryParams(d *schema.ResourceData, conf *config.Config) string {
	res := ""
	epsId := conf.GetEnterpriseProjectID(d)
	if epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}
	if v, ok := d.GetOk("domain"); ok {
		res = fmt.Sprintf("%s&hostname=%v", res, v)
	}
	if v, ok := d.GetOk("policy_name"); ok {
		res = fmt.Sprintf("%s&policyname=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
