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

// @API CFW GET /v2/{project_id}/rule/domain/resolve-ip-list/{domain_address_id}
func DataSourceCfwDomainResolveIpList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCfwDomainResolveIpListRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"fw_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the firewall instance ID.`,
			},
			"domain_address_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the domain address ID.`,
			},
			"address_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the address type.`,
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The data of domain resolve IP list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"excess_ip": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of excess IPs.`,
						},
						"parsed_success_ip": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of parsed success IPs.`,
						},
					},
				},
			},
		},
	}
}

func buildDomainResolveIpListQueryParams(d *schema.ResourceData) string {
	fwInstanceId := d.Get("fw_instance_id").(string)
	queryParams := fmt.Sprintf("?fw_instance_id=%s", fwInstanceId)

	addressType, ok := d.GetOk("address_type")
	if ok {
		queryParams += fmt.Sprintf("&address_type=%v", addressType)
	}

	return queryParams
}

func dataSourceCfwDomainResolveIpListRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		domainAddressId = d.Get("domain_address_id").(string)
		httpUrl         = "v2/{project_id}/rule/domain/resolve-ip-list/{domain_address_id}"
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{domain_address_id}", domainAddressId)
	requestPath += buildDomainResolveIpListQueryParams(d)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW domain resolve IP list: %s", err)
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
		d.Set("data", flattenDomainResolveIpListData(utils.PathSearch("data", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDomainResolveIpListData(data interface{}) []interface{} {
	if data == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"excess_ip":         utils.PathSearch("excess_ip", data, make([]interface{}, 0)),
			"parsed_success_ip": utils.PathSearch("parsed_success_ip", data, make([]interface{}, 0)),
		},
	}
}
