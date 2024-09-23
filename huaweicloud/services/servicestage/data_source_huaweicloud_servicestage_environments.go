package servicestage

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

// @API ServiceStage GET /v2/{project_id}/cas/environments
func DataSourceEnvironments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnvironmentsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The region where the environments are located.`,
			},
			"environments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the environment.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the environment.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the environment.`,
						},
						"deploy_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The deploy mode of the environment.`,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC ID to which the environment belongs.`,
						},
						"basic_resources": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataEnvResourcesSchema(),
							Description: `The basic resources.`,
						},
						"optional_resources": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataEnvResourcesSchema(),
							Description: `The optional resources.`,
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator name.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the environment, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the environment, in RFC3339 format.`,
						},
					},
				},
				Description: `All queried environments.`,
			},
		},
	}
}

func dataEnvResourcesSchema() *schema.Resource {
	sc := schema.Resource{
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
		},
	}
	return &sc
}

func queryEnvironments(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/cas/environments?limit=1000"
		offset  = 0
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

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
		environments := utils.PathSearch("environments", respBody, make([]interface{}, 0)).([]interface{})
		if len(environments) < 1 {
			break
		}
		result = append(result, environments...)
		offset += len(environments)
	}

	return result, nil
}

func flattenDataEnvResources(resources []interface{}) []interface{} {
	result := make([]interface{}, 0, len(resources))

	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"id":   utils.PathSearch("id", resource, nil),
			"type": utils.PathSearch("type", resource, nil),
		})
	}

	return result
}

func flattenDataEnvironments(environments []interface{}) []interface{} {
	result := make([]interface{}, 0, len(environments))

	for _, env := range environments {
		result = append(result, map[string]interface{}{
			"id":                 utils.PathSearch("id", env, nil),
			"name":               utils.PathSearch("name", env, nil),
			"alias":              utils.PathSearch("alias", env, nil),
			"description":        utils.PathSearch("description", env, nil),
			"deploy_mode":        utils.PathSearch("deploy_mode", env, nil),
			"vpc_id":             utils.PathSearch("vpc_id", env, nil),
			"basic_resources":    flattenDataEnvResources(utils.PathSearch("base_resources", env, make([]interface{}, 0)).([]interface{})),
			"optional_resources": flattenDataEnvResources(utils.PathSearch("optional_resources", env, make([]interface{}, 0)).([]interface{})),
			"creator":            utils.PathSearch("creator", env, nil),
			"created_at":         utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", env, float64(0)).(float64)/1000), false),
			"updated_at":         utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", env, float64(0)).(float64)/1000), false),
		})
	}

	return result
}

func dataSourceEnvironmentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	environments, err := queryEnvironments(client)
	if err != nil {
		return diag.Errorf("error querying environments: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("environments", flattenDataEnvironments(environments)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving data source fields of ServiceStage environments: %s", err)
	}
	return nil
}
