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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ServiceStage GET /v3/{project_id}/cas/environments
func DataSourceV3Environments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV3EnvironmentsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the environments are located.`,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the environment to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the environment to be queried.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project to which the environments belong.`,
			},
			"environments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The environment ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The environment name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the environment.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the enterprise project to which the environment belongs.`,
						},
						"deploy_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The deploy mode of the environment.`,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the VPC to which the environment belongs.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version number of the environment.`,
						},
						"tags": common.TagsComputedSchema(
							`The key/value pairs to associate with the environment.`,
						),
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator name of the environment.`,
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
				Description: "All runtime stack details.",
			},
		},
	}
}

func buildV3EnvironmentsQueryParams(d *schema.ResourceData) string {
	res := ""
	if envName, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, envName)
	}
	if envId, ok := d.GetOk("environment_id"); ok {
		res = fmt.Sprintf("%s&environment_id=%v", res, envId)
	}
	if epsId, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}
	return res
}

func queryV3Environments(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/cas/environments?limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildV3EnvironmentsQueryParams(d)

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
		envs := utils.PathSearch("environments", respBody, make([]interface{}, 0)).([]interface{})
		if len(envs) < 1 {
			break
		}
		result = append(result, envs...)
		offset += len(envs)
	}

	return result, nil
}

func flattenV3Environments(environments []interface{}) []map[string]interface{} {
	if len(environments) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(environments))
	for _, environment := range environments {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", environment, nil),
			"name":                  utils.PathSearch("name", environment, nil),
			"description":           utils.PathSearch("description", environment, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", environment, nil),
			"deploy_mode":           utils.PathSearch("deploy_mode", environment, nil),
			"vpc_id":                utils.PathSearch("vpc_id", environment, nil),
			"creator":               utils.PathSearch("creator", environment, nil),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
				environment, float64(0)).(float64))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time",
				environment, float64(0)).(float64))/1000, false),
			"tags": utils.FlattenTagsToMap(utils.PathSearch("labels", environment, nil)),
		})
	}
	return result
}

func dataSourceV3EnvironmentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	envList, err := queryV3Environments(client, d)
	if err != nil {
		return diag.Errorf("error getting environments: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("environments", flattenV3Environments(envList)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
