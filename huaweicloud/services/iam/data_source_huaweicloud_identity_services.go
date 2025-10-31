package iam

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v3/services/{service_id}
// @API IAM GET /v3/services
func DataSourceIdentityServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityServicesRead,

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"link": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	getIdentityServiceBasePath := iamClient.Endpoint + "v3/services"
	getIdentityServiceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	serviceId := d.Get("service_id").(string)
	if serviceId != "" {
		getIdentityServiceBasePath = getIdentityServiceBasePath + "/" + serviceId
	}
	response, err := iamClient.Request("GET", getIdentityServiceBasePath, &getIdentityServiceOpt)
	if err != nil {
		return diag.Errorf("error getting identity services: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	if serviceId != "" {
		resp := utils.PathSearch("service", respBody, nil)
		if resp == nil {
			return nil
		}
		services := flattenServicesList([]interface{}{resp})
		if e := d.Set("services", services); e != nil {
			return diag.FromErr(e)
		}
	} else {
		services := flattenServicesList(utils.PathSearch("services", respBody, make([]interface{}, 0)).([]interface{}))
		if e := d.Set("services", services); e != nil {
			return diag.FromErr(e)
		}
	}
	return nil
}

func flattenServicesList(servicesBody []interface{}) []map[string]interface{} {
	services := make([]map[string]interface{}, len(servicesBody))
	for i, service := range servicesBody {
		services[i] = map[string]interface{}{
			"name":    utils.PathSearch("name", service, ""),
			"link":    utils.PathSearch("links.self", service, nil),
			"id":      utils.PathSearch("id", service, ""),
			"type":    utils.PathSearch("type", service, ""),
			"enabled": utils.PathSearch("enabled", service, nil),
		}
	}
	return services
}
