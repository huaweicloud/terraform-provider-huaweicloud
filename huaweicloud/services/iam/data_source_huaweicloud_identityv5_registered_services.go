package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/authorization-schemas/registered-services
func DataSourceV5RegisteredServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5RegisteredServicesRead,

		Schema: map[string]*schema.Schema{
			"service_codes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: `The list of service codes.`,
			},
		},
	}
}

func dataSourceV5RegisteredServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var (
		httpUrl     = "v5/authorization-schemas/registered-services"
		allServices = make([]interface{}, 0)
		// Default limit is 100, maximum is 200.
		limit  = 200
		marker = ""
	)
	listPath := client.Endpoint + fmt.Sprintf("%s?limit=%d", httpUrl, limit)
	for {
		listPathWithMarker := listPath + buildV5RegisteredServicesQueryParams(marker)
		reqOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		r, err := client.Request("GET", listPathWithMarker, &reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving registered services: %s", err)
		}

		resp, err := utils.FlattenResponse(r)
		if err != nil {
			return diag.FromErr(err)
		}

		services := utils.PathSearch("service_codes", resp, make([]interface{}, 0)).([]interface{})
		allServices = append(allServices, services...)
		if len(services) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", resp, "").(string)
		if marker == "" {
			break
		}
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomId)

	return diag.FromErr(d.Set("service_codes", allServices))
}

func buildV5RegisteredServicesQueryParams(marker string) string {
	if marker != "" {
		return fmt.Sprintf("&marker=%v", marker)
	}
	return ""
}
