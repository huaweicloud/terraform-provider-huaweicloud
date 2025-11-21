package iam

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

// @API IAM GET /v5/authorization-schemas/registered-services
func DataSourceIdentityV5RegisteredServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5RegisteredServicesRead,

		Schema: map[string]*schema.Schema{
			"service_codes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceIdentityV5RegisteredServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var allServices []interface{}
	var marker string
	var path string

	for {
		path = client.Endpoint + "v5/authorization-schemas/registered-services" + buildListRegisteredServicesV5Params(marker)
		reqOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		r, err := client.Request("GET", path, &reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving registered services: %s", err)
		}
		resp, err := utils.FlattenResponse(r)
		if err != nil {
			return diag.FromErr(err)
		}

		services := flattenListRegisteredServicesV5Response(resp)
		allServices = append(allServices, services...)

		marker = utils.PathSearch("page_info.next_marker", resp, "").(string)
		if marker == "" {
			break
		}
	}

	id, _ := uuid.GenerateUUID()
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("service_codes", allServices),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListRegisteredServicesV5Params(marker string) string {
	res := "?limit=100"
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func flattenListRegisteredServicesV5Response(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	serviceCodes := utils.PathSearch("service_codes", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(serviceCodes))
	copy(result, serviceCodes)
	return result
}
