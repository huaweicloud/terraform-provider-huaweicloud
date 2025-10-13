package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/ip-info
func DataSourceIpInformation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIpInformationRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the queried IP attribution information are located.`,
			},

			// Required parameters.
			"ips": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The list of IP addresses to be queried.`,
			},

			// Optional parameters.
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project to which the IP attribution information belongs.`,
			},

			// Attributes.
			"information": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        informationSchema(),
				Description: `The list of IP attribution information that matched filter parameters.`,
			},
		},
	}
}

func informationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP address to be queried.`,
			},
			"belongs": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the IP belongs to CDN nodes.`,
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The province where the IP is located.`,
			},
			"isp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ISP name.`,
			},
			"platform": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The platform name.`,
			},
		},
	}
}

func buildIpInformationQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("ips"); ok {
		res = fmt.Sprintf("%sips=%v", res, v)
	}

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}

	if len(res) < 1 {
		return res
	}
	return "?" + res
}

func flattenIpsInformation(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"ip":       utils.PathSearch("ip", item, nil),
			"belongs":  utils.PathSearch("belongs", item, nil),
			"region":   utils.ValueIgnoreEmpty(utils.PathSearch("region", item, nil)),
			"isp":      utils.ValueIgnoreEmpty(utils.PathSearch("isp", item, nil)),
			"platform": utils.ValueIgnoreEmpty(utils.PathSearch("platform", item, nil)),
		})
	}

	return result
}

func dataSourceIpInformationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/cdn/ip-info"
	)

	client, err := cfg.NewServiceClient("cdn", region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath += buildIpInformationQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return diag.Errorf("error querying IP attribution information: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing IP attribution information response: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("information", flattenIpsInformation(utils.PathSearch("cdn_ips", respBody,
			make([]interface{}, 0)).([]interface{}))))
	return diag.FromErr(mErr.ErrorOrNil())
}
