package servicestage

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ServiceStage GET /v2/{project_id}/cas/applications
// @API ServiceStage GET /v2/{project_id}/cas/applications/{application_id}/configuration
func DataSourceApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The region where the applications are located.`,
			},
			"ignore_environments": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to ignore environments query. `,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the application.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the application.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the application.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterpeise project ID to which the application belongs.`,
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator name.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the application, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the application, in RFC3339 format.`,
						},
						"component_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of the components associated with the application.`,
						},
						"environments": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataAppEnvVariablesSchema(),
							Description: `The environment configuration associated with the application.`,
						},
					},
				},
				Description: `All queried applications.`,
			},
		},
	}
}

func dataAppEnvVariablesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The environment ID.`,
			},
			"variables": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The variable name.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The variable value.`,
						},
					},
				},
				Description: `The variables of the environment.`,
			},
		},
	}
	return &sc
}

func queryApplications(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/cas/applications?limit=1000"
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
		applications := utils.PathSearch("applications", respBody, make([]interface{}, 0)).([]interface{})
		if len(applications) < 1 {
			break
		}
		result = append(result, applications...)
		offset += len(applications)
	}

	return result, nil
}

func queryAppConfigurations(client *golangsdk.ServiceClient, appId string) ([]interface{}, error) {
	httpUrl := "v2/{project_id}/cas/applications/{application_id}/configuration"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{application_id}", appId)

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
	configurations := utils.PathSearch("configuration", respBody, make([]interface{}, 0)).([]interface{})
	return configurations, nil
}

func flattenAppConfigurationVariables(variables []interface{}) []interface{} {
	result := make([]interface{}, 0, len(variables))

	for _, variable := range variables {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", variable, nil),
			"value": utils.PathSearch("value", variable, nil),
		})
	}
	return result
}

func flattenAppConfigurations(configurations []interface{}) []interface{} {
	result := make([]interface{}, 0, len(configurations))

	for _, configuration := range configurations {
		result = append(result, map[string]interface{}{
			"id": utils.PathSearch("environment_id", configuration, nil),
			"variables": flattenAppConfigurationVariables(utils.PathSearch("configuration.env",
				configuration, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenDataApplications(client *golangsdk.ServiceClient, applications []interface{}, ignoreVars bool) []interface{} {
	result := make([]interface{}, 0, len(applications))

	for _, app := range applications {
		appId := utils.PathSearch("id", app, "").(string)
		appElem := map[string]interface{}{
			"id":                    appId,
			"name":                  utils.PathSearch("name", app, nil),
			"description":           utils.PathSearch("description", app, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", app, nil),
			"creator":               utils.PathSearch("creator", app, nil),
			"created_at":            utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", app, float64(0)).(float64)/1000), false),
			"updated_at":            utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", app, float64(0)).(float64)/1000), false),
			"component_count":       utils.PathSearch("component_count", app, nil),
		}
		if !ignoreVars {
			configurations, err := queryAppConfigurations(client, appId)
			if err != nil {
				log.Printf("[ERROR] unable to find the application (%s) configuration: %s", appId, err)
			} else {
				appElem["environments"] = flattenAppConfigurations(configurations)
			}
		}
		result = append(result, appElem)
	}

	return result
}

func dataSourceApplicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	applications, err := queryApplications(client)
	if err != nil {
		return diag.Errorf("error querying applications: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("applications", flattenDataApplications(client, applications, d.Get("ignore_variables").(bool))),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving data source fields of ServiceStage environments: %s", err)
	}
	return nil
}
