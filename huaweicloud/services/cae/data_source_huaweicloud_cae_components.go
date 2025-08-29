package cae

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CAE GET /v1/{project_id}/cae/applications/{application_id}/components
func DataSourceComponents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComponentsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the components are located.`,
			},

			// Required parameter(s).
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the environment to which the components belong.`,
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the application to which the components belong.`,
			},

			// Optional parameter(s).
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project to which the components belong.`,
			},

			"components": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The component ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the component.`,
						},
						"annotations": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The parameters of key/value pairs related to the component.`,
						},
						"spec": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"runtime": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The component runtime.`,
									},
									"environment_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the environment to which the component belongs.`,
									},
									"replica": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The instance number of the component.`,
									},
									"available_replica": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The available instance number of the component.`,
									},
									"source": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The code source configuration information corresponding to the component.`,
									},
									"build": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The build information of the code source corresponding to the component.`,
									},
									"resource_limit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The instance specification corresponding to the component.`,
									},
									"image_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The image URL that component used.`,
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The status of the component.`,
									},
								},
							},
							Description: "The configuration information of the component.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the component, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the component, in RFC3339 format.`,
						},
					},
				},
				Description: `All queried components.`,
			},
		},
	}
}

func queryComponents(client *golangsdk.ServiceClient, epsId, environmentId, appId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/cae/applications/{application_id}/components?limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{application_id}", appId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(environmentId, epsId),
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		components := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(components) < 1 {
			break
		}
		result = append(result, components...)
		offset += len(components)
	}

	return result, nil
}

func flattenComponentAnnotations(annotations map[string]interface{}) map[string]interface{} {
	if len(annotations) < 1 {
		return nil
	}

	result := make(map[string]interface{})
	// All values convert to the string values.
	for key, value := range annotations {
		result[key] = fmt.Sprint(value)
	}
	return result
}

func flattenComponentSpec(spec interface{}) []map[string]interface{} {
	if spec == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"runtime":           utils.PathSearch("runtime", spec, nil),
			"environment_id":    utils.PathSearch("env_id", spec, nil),
			"replica":           utils.PathSearch("replica", spec, nil),
			"available_replica": utils.PathSearch("available_replica", spec, nil),
			"source":            utils.JsonToString(utils.PathSearch("source", spec, nil)),
			"build":             utils.JsonToString(utils.PathSearch("build", spec, nil)),
			"resource_limit":    utils.JsonToString(utils.PathSearch("resource_limit", spec, nil)),
			"image_url":         utils.PathSearch("image_url", spec, nil),
			"status":            utils.PathSearch("status", spec, nil),
		},
	}
}

func flattenComponents(components []interface{}) []map[string]interface{} {
	if len(components) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(components))
	for _, component := range components {
		result = append(result, map[string]interface{}{
			"id":   utils.PathSearch("id", component, nil),
			"name": utils.PathSearch("name", component, nil),
			"annotations": flattenComponentAnnotations(utils.PathSearch("annotations",
				component, make(map[string]interface{})).(map[string]interface{})),
			"spec": flattenComponentSpec(utils.PathSearch("spec", component, make(map[string]interface{})).(map[string]interface{})),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at",
				component, "").(string))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("updated_at",
				component, "").(string))/1000, false),
		})
	}

	return result
}

func dataSourceComponentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		envId  = d.Get("environment_id").(string)
		appId  = d.Get("application_id").(string)
	)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	components, err := queryComponents(client, cfg.GetEnterpriseProjectID(d), envId, appId)
	if err != nil {
		return diag.Errorf("error getting components: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("components", flattenComponents(components)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
