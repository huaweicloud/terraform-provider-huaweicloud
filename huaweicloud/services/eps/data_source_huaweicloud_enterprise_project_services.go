package eps

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

// @API EPS GET /v1.0/enterprise-projects/providers
func DataSourceEpsServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEpsServicesRead,

		Schema: map[string]*schema.Schema{
			"locale": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildEpsServicesSchema(),
			},
		},
	}
}

func buildEpsServicesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_i18n_display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildEpsResourceTypeSchema(),
			},
		},
	}
}

func buildEpsResourceTypeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type_i18n_display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildEpsServicesQueryParams(d *schema.ResourceData, offset int) string {
	rst := "?limit=200"

	if v, ok := d.GetOk("locale"); ok {
		rst += fmt.Sprintf("&locale=%v", v)
	}

	if v, ok := d.GetOk("service"); ok {
		rst += fmt.Sprintf("&provider=%v", v)
	}

	if offset > 0 {
		rst += fmt.Sprintf("&offset=%d", offset)
	}

	return rst
}

func dataSourceEpsServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1.0/enterprise-projects/providers"
		offset     = 0
		allResults = make([]interface{}, 0)
	)
	client, err := cfg.NewServiceClient("eps", region)
	if err != nil {
		return diag.Errorf("error creating EPS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		requestPathWithOffset := requestPath + buildEpsServicesQueryParams(d, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving EPS services: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		providers := utils.PathSearch("providers", respBody, make([]interface{}, 0)).([]interface{})
		if len(providers) == 0 {
			break
		}

		allResults = append(allResults, providers...)
		offset += len(providers)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("services", flattenListEpsServicesResponseBody(allResults)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListEpsServicesResponseBody(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		resourceTypes := utils.PathSearch("resource_types", v, make([]interface{}, 0)).([]interface{})
		rst = append(rst, map[string]interface{}{
			"service":                   utils.PathSearch("provider", v, nil),
			"service_i18n_display_name": utils.PathSearch("provider_i18n_display_name", v, nil),
			"resource_types":            flattenResourceTypeBody(resourceTypes),
		})
	}
	return rst
}

func flattenResourceTypeBody(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"resource_type":                   utils.PathSearch("resource_type", v, nil),
			"resource_type_i18n_display_name": utils.PathSearch("resource_type_i18n_display_name", v, nil),
			"global":                          utils.PathSearch("global", v, nil),
			"regions":                         utils.PathSearch("regions", v, nil),
		})
	}
	return rst
}
