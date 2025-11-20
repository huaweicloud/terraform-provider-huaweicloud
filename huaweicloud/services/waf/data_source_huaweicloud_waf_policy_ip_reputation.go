package waf

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

// @API WAF GET /v1/{project_id}/waf/tag/ip-reputation/map
func DataSourcePolicyIpReputation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyIpReputationRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_reputation_map": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildIpReputationMapSchema(),
			},
			"locale": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildIpReputationMapSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"idc": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourcePolicyIpReputationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1/{project_id}/waf/tag/ip-reputation/map"
		langValue = d.Get("lang").(string)
		typeValue = d.Get("type").(string)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?lang=%s&type=%s", langValue, typeValue)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving WAF policy IP reputation: %s", err)
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

	ipReputationMap := utils.PathSearch("ip_reputation_map", respBody, nil)
	localeMap := utils.PathSearch("locale", respBody, make(map[string]interface{})).(map[string]interface{})

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("ip_reputation_map", flattenPolicyIpReputation(ipReputationMap)),
		d.Set("locale", utils.ExpandToStringMap(localeMap)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPolicyIpReputation(ipReputationMap interface{}) []interface{} {
	if ipReputationMap == nil {
		return nil
	}

	idcArray := utils.PathSearch("idc", ipReputationMap, make([]interface{}, 0)).([]interface{})
	return []interface{}{
		map[string]interface{}{
			"idc": utils.ExpandToStringList(idcArray),
		},
	}
}
