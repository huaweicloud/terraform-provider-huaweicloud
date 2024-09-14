package apig

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/environments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/envs
func DataSourceEnvironments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnvironmentsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flattenEnvironments(envList []environments.Environment) ([]map[string]interface{}, []string) {
	if len(envList) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(envList))
	ids := make([]string, len(envList))
	for i, env := range envList {
		result[i] = map[string]interface{}{
			"id":          env.Id,
			"name":        env.Name,
			"description": env.Description,
			"create_time": env.CreateTime,
		}
		ids[i] = env.Id
	}
	return result, ids
}

func dataSourceEnvironmentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opt := environments.ListOpts{
		Name: d.Get("name").(string),
	}
	pages, err := environments.List(client, d.Get("instance_id").(string), opt).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := environments.ExtractEnvironments(pages)
	if err != nil {
		return diag.Errorf("unable to get the environment list form server: %v", err)
	}

	envResult, ids := flattenEnvironments(resp)
	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("environments", envResult),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
