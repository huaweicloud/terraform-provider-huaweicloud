package cse

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
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
			"engine_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Nacos microservice engine to which the namespaces belong.`,
			},
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
				Description: `All queried Nacos namespaces.`,
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
		conf     = meta.(*config.Config)
		region   = conf.GetRegion(d)
		engineId = d.Get("engine_id").(string)
	)
	client, err := conf.NewServiceClient("cse", region)
	if err != nil {
		return diag.Errorf("error creating CSE client: %s", err)
	}

	namespaces, err := listNacosNamespaces(client, engineId)
	if err != nil {
		return diag.Errorf("error querying namespaces under Nacos engine (%s): %s", engineId, err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("namespaces", flattenNacosNamespaces(namespaces)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
