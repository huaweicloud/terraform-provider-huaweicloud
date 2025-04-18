package servicestage

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ServiceStage GET /v3/{project_id}/cas/components
func DataSourceV3Components() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV3ComponentsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the components are located.`,
			},
			"application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the application to which the components belong.`,
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
						"environment_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The environment ID where the component is deployed.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the component.`,
						},
						"runtime_stack": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The stack name.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The stack type.`,
									},
									"deploy_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The deploy mode of the stack.`,
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The stack version.`,
									},
								},
							},
							Description: "The configuration of the runtime stack.",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source configuration of the component, in JSON format.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version of the component.`,
						},
						"refer_resources": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource ID.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource type.`,
									},
									"parameters": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource parameters, in JSON format.`,
									},
								},
							},
							Description: `The configuration of the reference resources.`,
						},
						"external_accesses": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The protocol of the external access.`,
									},
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The address of the external access.`,
									},
									"forward_port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The forward port of the external access.`,
									},
								},
							},
							Description: "The configuration of the external accesses.",
						},
						"tags": common.TagsComputedSchema(
							`The key/value pairs to associate with the component.`,
						),
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the component.`,
						},
						// Deprecated attributes.
						"build": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The build configuration of the component, in JSON format.`,
								utils.SchemaDescInput{
									Deprecated: true,
								},
							),
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The creation time of the component, in RFC3339 format.`,
								utils.SchemaDescInput{
									Deprecated: true,
								},
							),
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The latest update time of the component, in RFC3339 format.`,
								utils.SchemaDescInput{
									Deprecated: true,
								},
							),
						},
					},
				},
				Description: "All application details.",
			},
		},
	}
}

func listV3Components(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/cas/components?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
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
		components := utils.PathSearch("components", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, components...)
		if len(components) < limit {
			break
		}
		offset += len(components)
	}

	return result, nil
}

func flattenV3Components(components []interface{}) []map[string]interface{} {
	if len(components) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(components))
	for _, component := range components {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", component, nil),
			"environment_id": utils.PathSearch("environment_id", component, nil),
			"name":           utils.PathSearch("name", component, nil),
			"runtime_stack": flattenV3ComponentRuntimeStackConfig(utils.PathSearch("runtime_stack", component,
				make(map[string]interface{})).(map[string]interface{})),
			"source":  utils.JsonToString(utils.PathSearch("source", component, nil)),
			"version": utils.PathSearch("version", component, nil),
			"refer_resources": flattenV3ComponentReferResources(utils.PathSearch("refer_resources", component,
				make([]interface{}, 0)).([]interface{})),
			"build": utils.JsonToString(utils.PathSearch("build", component, nil)),
			"external_accesses": flattenV3ExternalAccesses(utils.PathSearch("external_accesses", component,
				make([]interface{}, 0)).([]interface{})),
			"tags":   utils.FlattenTagsToMap(utils.PathSearch("labels", component, make([]interface{}, 0))),
			"status": utils.PathSearch("status.component_status", component, nil),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("status.create_time", component,
				float64(0)).(float64))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("status.update_time", component,
				float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func filterV3Components(d *schema.ResourceData, components []interface{}) []interface{} {
	result := components

	if applicationId, ok := d.GetOk("application_id"); ok {
		result = utils.PathSearch(fmt.Sprintf("[?application_id=='%s']", applicationId),
			result, make([]interface{}, 0)).([]interface{})
	}

	return result
}

func dataSourceV3ComponentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	components, err := listV3Components(client)
	if err != nil {
		return diag.Errorf("error getting components: %s", err)
	}

	components = filterV3Components(d, components)

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("components", flattenV3Components(components)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
