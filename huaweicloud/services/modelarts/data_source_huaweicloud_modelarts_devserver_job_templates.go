package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v1/{project_id}/dev-servers/jobs/templates
func DataSourceDevServerJobTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDevServerJobTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the DevServer job templates are located.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the DevServer job template to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the DevServer job template to be queried.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the DevServer job template to be queried.`,
			},

			// Attributes
			"templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of DevServer job templates that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the DevServer job template.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the DevServer job template.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the DevServer job template.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the DevServer job template.`,
						},
						"flavor_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The flavor type of the DevServer job template.`,
						},
						"params": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The parameters of the DevServer job template.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the parameter.`,
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The description of the parameter.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The value of the parameter.`,
									},
									"visible": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Whether the parameter is visible in the console.`,
									},
									"regex": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The regular expression for validating the parameter value.`,
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

func buildListDevServerJobTemplatesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("template_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}

	if len(res) > 0 {
		res = "?" + res[1:]
	}

	return res
}

func listDevServerJobTemplates(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/dev-servers/jobs/templates"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildListDevServerJobTemplatesQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenDevServerJobTemplateParams(params []interface{}) []map[string]interface{} {
	if len(params) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(params))
	for _, param := range params {
		result = append(result, map[string]interface{}{
			"name":        utils.PathSearch("name", param, nil),
			"description": utils.PathSearch("description", param, nil),
			"value":       utils.PathSearch("value", param, nil),
			"visible":     utils.PathSearch("visible", param, nil),
			"regex":       utils.PathSearch("regex", param, nil),
		})
	}

	return result
}

func flattenDevServerJobTemplates(templates []interface{}) []map[string]interface{} {
	if len(templates) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(templates))
	for _, template := range templates {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", template, nil),
			"name":        utils.PathSearch("name", template, nil),
			"description": utils.PathSearch("description", template, nil),
			"type":        utils.PathSearch("type", template, nil),
			"flavor_type": utils.PathSearch("flavor_type", template, nil),
			"params": flattenDevServerJobTemplateParams(
				utils.PathSearch("params", template, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func dataSourceDevServerJobTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	templates, err := listDevServerJobTemplates(client, d)
	if err != nil {
		return diag.Errorf("error querying DevServer job templates: %s", err)
	}

	randomUUID, err := uuid.NewUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("templates", flattenDevServerJobTemplates(templates)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
