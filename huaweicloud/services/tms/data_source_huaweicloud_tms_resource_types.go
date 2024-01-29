package tms

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/tms/v1/providers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TMS GET /v1.0/tms/providers
func DataSourceResourceTypes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceTypesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The region name used to filter resource types information.`,
			},
			"service_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The service name used to filter resource types information.`,
			},
			// Attributes
			"types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource type name.`,
						},
						"is_global": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the resource corresponding to this type is a global resource.`,
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The service display name of the resource type.`,
						},
						"service_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the service to which the resource type belong.`,
						},
					},
				},
				Description: `All resource types that match the filter parameters.`,
			},
		},
	}
}

func filterResourceTypes(serviceName, region string, providerList []providers.Provider) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(providerList))
	for _, val := range providerList {
		if serviceName != "" && val.Provider != serviceName {
			continue
		}
		for _, resource := range val.Resources {
			if region != "" && !utils.StrSliceContains(resource.Regions, region) {
				continue
			}
			result = append(result, map[string]interface{}{
				"name":         resource.ResourceType,
				"is_global":    resource.Global,
				"display_name": resource.ResourceTypeI18nDisplayName,
				"service_name": val.Provider,
			})
		}
	}
	return result
}

func dataSourceResourceTypesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg, ok := meta.(*config.Config)
	if !ok {
		return diag.Errorf("invalid type of the meta argument, want '*config.Config', but got '%T'", meta)
	}
	// The region parameter is only used to filter query results and is not used to build the client.
	client, err := cfg.TmsV1Client(cfg.Region)
	if err != nil {
		return diag.Errorf("error creating TMS v1 client: %s", err)
	}

	var (
		region      = d.Get("region").(string)
		serviceName = d.Get("service_name").(string)
		opts        = providers.ListOpts{
			Provider: serviceName,
			Limit:    100,
		}
	)
	resp, err := providers.List(client, opts)
	if err != nil {
		return diag.Errorf("error querying resource types: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	err = d.Set("types", filterResourceTypes(serviceName, region, resp))
	if err != nil {
		return diag.Errorf("error setting resource types field: %s", err)
	}
	return nil
}
