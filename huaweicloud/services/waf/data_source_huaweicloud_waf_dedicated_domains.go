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

// @API WAF GET /v1/{project_id}/premium-waf/host
func DataSourceWafDedicatedDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceDedicatedDomainsRead,

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
			"protect_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Elem:     domainSchema(),
				Computed: true,
			},
		},
	}
}

func domainSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pci_3ds": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"pci_dds": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_dual_az": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"protect_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"access_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"website_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func datasourceDedicatedDomainsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/premium-waf/host"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildWAFDedicatedDomainsQueryParams(d, cfg)
	requestResp, err := pagination.ListAllItems(
		client,
		"page",
		requestPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving WAF dedicated domains, %s", err)
	}

	respJson, err := json.Marshal(requestResp)
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
		d.Set("domains", flattenListDedicatedDomainsBody(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListDedicatedDomainsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("items", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"domain":                utils.PathSearch("hostname", v, nil),
			"policy_id":             utils.PathSearch("policyid", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"protect_status":        utils.PathSearch("protect_status", v, nil),
			"website_name":          utils.PathSearch("web_tag", v, nil),
			"access_status":         utils.PathSearch("access_status", v, nil),
			"pci_3ds":               utils.StringToBool(utils.PathSearch("flag.pci_3ds", v, "")),
			"pci_dds":               utils.StringToBool(utils.PathSearch("flag.pci_dds", v, "")),
			"is_dual_az":            utils.StringToBool(utils.PathSearch("flag.is_dual_az", v, "")),
			"ipv6_enable":           utils.StringToBool(utils.PathSearch("flag.ipv6", v, "")),
		})
	}
	return rst
}

func buildWAFDedicatedDomainsQueryParams(d *schema.ResourceData, conf *config.Config) string {
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
	if v, ok := d.GetOk("protect_status"); ok {
		res = fmt.Sprintf("%s&protect_status=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
