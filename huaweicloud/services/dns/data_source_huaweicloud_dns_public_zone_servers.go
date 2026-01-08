package dns

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

// @API DNS GET /v2/public-zones/dns-servers/{domain_name}
func DataSourcePublicZoneServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePublicZoneServersRead,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The domain name of the public zone.`,
			},
			"all_hw_dns": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether all servers are HuaweiCloud DNS servers.`,
			},
			"include_hw_dns": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether HuaweiCloud DNS servers are included.`,
			},
			"dns_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of DNS server addresses.`,
			},
			"expected_dns_servers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of expected DNS server addresses.`,
			},
		},
	}
}

func getPublicZoneServers(client *golangsdk.ServiceClient, domainName string) (interface{}, error) {
	httpUrl := "v2/public-zones/dns-servers/{domain_name}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_name}", domainName)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}
	requestResp, err := client.Request("GET", getPath, &reqOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func dataSourcePublicZoneServersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dns", "")
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	respBody, err := getPublicZoneServers(client, d.Get("domain_name").(string))
	if err != nil {
		return diag.Errorf("error querying public zone DNS servers: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("all_hw_dns", utils.PathSearch("all_hw_dns", respBody, nil)),
		d.Set("include_hw_dns", utils.PathSearch("include_hw_dns", respBody, nil)),
		d.Set("dns_servers", utils.PathSearch("dns_servers", respBody, nil)),
		d.Set("expected_dns_servers", utils.PathSearch("expected_dns_servers", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
