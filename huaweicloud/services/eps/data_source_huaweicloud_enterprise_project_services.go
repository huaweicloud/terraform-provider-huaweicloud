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

			// Attributes.
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
							Elem: &schema.Resource{
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
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceEpsServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("eps", region)
	if err != nil {
		return diag.Errorf("error creating EPS client: %s", err)
	}

	listPathBase := client.Endpoint + buildEpsServicesQueryParameter(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	res := make([]interface{}, 0)
	offset := ""
	for {
		listPath := listPathBase + buildEpsServicesOffsetPath(offset)
		getResp, err := client.Request("GET", listPath, &opt)
		if err != nil {
			return diag.Errorf("error retrieving EPS services: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		res = append(res, flattenListEpsServicesResponseBody(getRespBody)...)
		offset = utils.PathSearch("offset", getRespBody, "").(string)
		if offset == "" {
			break
		}
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("services", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildEpsServicesQueryParameter(d *schema.ResourceData) string {
	httpUrl := "v1.0/enterprise-projects/providers?limit=200"

	if d.Get("locale") != "" {
		httpUrl += fmt.Sprintf("&locale=%s", d.Get("locale"))
	}

	if d.Get("service") != "" {
		httpUrl += fmt.Sprintf("&provider=%s", d.Get("service"))
	}

	return httpUrl
}

func buildEpsServicesOffsetPath(offset string) string {
	res := ""
	if offset != "" {
		res += fmt.Sprintf("&offset=%s", offset)
	}
	return res
}

func flattenListEpsServicesResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("providers", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0)
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"service":                   utils.PathSearch("provider", v, nil),
			"service_i18n_display_name": utils.PathSearch("provider_i18n_display_name", v, nil),
			"resource_types":            flattenResourceTypeBody(utils.PathSearch("resource_types", v, nil)),
		})
	}
	return res
}

func flattenResourceTypeBody(resourceTypeBody interface{}) []interface{} {
	if resourceTypeBody == nil {
		return nil
	}
	resourceTypeArray := resourceTypeBody.([]interface{})
	res := make([]interface{}, 0)
	for _, v := range resourceTypeArray {
		res = append(res, map[string]interface{}{
			"resource_type":                   utils.PathSearch("resource_type", v, nil),
			"resource_type_i18n_display_name": utils.PathSearch("resource_type_i18n_display_name", v, nil),
			"global":                          utils.PathSearch("global", v, nil),
			"regions":                         utils.PathSearch("regions", v, nil),
		})
	}
	return res
}
