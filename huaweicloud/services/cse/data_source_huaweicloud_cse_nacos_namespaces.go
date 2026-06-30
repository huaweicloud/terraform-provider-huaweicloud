package cse

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSE GET /v1/{project_id}/nacos/v1/console/namespaces
func DataSourceNacosNamespaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNacosNamespacesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the Nacos namespaces are located.`,
			},

			// Required parameters.
			"engine_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Nacos microservice engine to which the namespaces belong.`,
			},

			// Optional parameters.
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the enterprise project to which the Nacos namespaces belong.`,
			},

			// Attributes.
			"namespaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the Nacos namespace.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the Nacos namespace.`,
						},
					},
				},
				Description: `All queried Nacos namespaces that match the filter parameters.`,
			},
		},
	}
}

func flattenNacosNamespaces(namespaces []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(namespaces))

	for _, namespace := range namespaces {
		result = append(result, map[string]interface{}{
			"id":   utils.PathSearch("namespace", namespace, nil),
			"name": utils.PathSearch("namespaceShowName", namespace, nil),
		})
	}

	return result
}

func dataSourceNacosNamespacesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf                = meta.(*config.Config)
		region              = conf.GetRegion(d)
		engineId            = d.Get("engine_id").(string)
		enterpriseProjectId = conf.GetEnterpriseProjectID(d)
	)
	client, err := conf.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	namespaces, err := listNacosNamespaces(client, engineId, enterpriseProjectId)
	if err != nil {
		return diag.Errorf("error querying namespaces under Nacos engine (%s): %s", engineId, err)
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Required parameters.
		d.Set("engine_id", engineId),
		// Optional parameters.
		d.Set("enterprise_project_id", enterpriseProjectId),
		// Attributes.
		d.Set("namespaces", flattenNacosNamespaces(namespaces)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
