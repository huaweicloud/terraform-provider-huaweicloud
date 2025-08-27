package aad

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

// Due to limited testing conditions, this data source cannot be tested and the API was not successfully called.

// @API AAD GET /v2/aad/instances/{instance_id}/domains
func DataSourceInstanceDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceDomainsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the instance ID.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance name.",
			},
			"domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The domain information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain ID.",
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain name.",
						},
						"cname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain CNAME.",
						},
						"domain_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain status. `0` represents normal, `1` represents freeze.",
						},
						"cc_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The CC protection status.",
						},
						"https_cert_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The certificate status. `1` represents uploaded, `2` represents not uploaded",
						},
						"cert_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate name.",
						},
						"protocol_type": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The domain protocol list.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"real_server_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The real server type.",
						},
						"real_servers": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The real servers.",
						},
						"waf_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The WAF protection status.",
						},
					},
				},
			},
		},
	}
}

func dataSourceInstanceDomainsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/instances/{instance_id}/domains"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", d.Get("instance_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving AAD instance domains: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("instance_name", utils.PathSearch("instance_name", respBody, nil)),
		d.Set("domains", flattenInstanceDomains(utils.PathSearch("domains", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstanceDomains(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"domain_id":         utils.PathSearch("domain_id", v, nil),
			"domain_name":       utils.PathSearch("domain_name", v, nil),
			"cname":             utils.PathSearch("cname", v, nil),
			"domain_status":     utils.PathSearch("domain_status", v, nil),
			"cc_status":         utils.PathSearch("cc_status", v, nil),
			"https_cert_status": utils.PathSearch("https_cert_status", v, nil),
			"cert_name":         utils.PathSearch("cert_name", v, nil),
			"protocol_type":     utils.ExpandToStringList(utils.PathSearch("protocol_type", v, make([]interface{}, 0)).([]interface{})),
			"real_server_type":  utils.PathSearch("real_server_type", v, nil),
			"real_servers":      utils.PathSearch("real_servers", v, nil),
			"waf_status":        utils.PathSearch("waf_status", v, nil),
		})
	}

	return rst
}
