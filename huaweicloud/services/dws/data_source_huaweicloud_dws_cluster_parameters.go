package dws

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS GET /v1.0/{project_id}/clusters/{cluster_id}/configurations/{configuration_id}
func DataSourceClusterParameters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterParametersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The DWS cluster ID.`,
			},
			"parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the parameters under specified DWS cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the parameter.`,
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of the parameter values.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the parameter.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The value of the parameter.`,
									},
									"default_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The default value of the parameter.`,
									},
								},
							},
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unit of the parameter.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the parameter value.`,
						},
						"readonly": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the parameter is read-only.`,
						},
						"value_range": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The range of the parameter value.`,
						},
						"restart_required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the DWS cluster needs to be restarted after modifying the parameter value.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the parameter.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceClusterParametersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	resp, err := GetParameterConfigurations(client, d.Get("cluster_id").(string))
	if err != nil {
		return diag.Errorf("error retrieving DWS cluster parameters: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("parameters",
			flattenClusterParameters(utils.PathSearch("configurations", resp, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenClusterParameters(all []interface{}) []map[string]interface{} {
	if len(all) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(all))
	for i, v := range all {
		result[i] = map[string]interface{}{
			"name":             utils.PathSearch("name", v, nil),
			"values":           flattenParamterValues(utils.PathSearch("values", v, make([]interface{}, 0)).([]interface{})),
			"unit":             utils.PathSearch("unit", v, nil),
			"type":             utils.PathSearch("type", v, nil),
			"readonly":         utils.PathSearch("readonly", v, false),
			"value_range":      utils.PathSearch("value_range", v, nil),
			"restart_required": utils.PathSearch("restart_required", v, false),
			"description":      utils.PathSearch("description", v, nil),
		}
	}
	return result
}

func flattenParamterValues(values []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(values))
	for i, v := range values {
		result[i] = map[string]interface{}{
			"type":          utils.PathSearch("type", v, nil),
			"value":         utils.PathSearch("value", v, nil),
			"default_value": utils.PathSearch("default_value", v, false),
		}
	}
	return result
}
