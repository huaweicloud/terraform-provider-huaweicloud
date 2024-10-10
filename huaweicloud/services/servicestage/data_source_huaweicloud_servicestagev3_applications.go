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

// @API ServiceStage GET /v3/{project_id}/cas/applications
func DataSourceV3Applications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV3ApplicationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the applications are located.`,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the application.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the enterprise project to which the application belongs.`,
						},
						"tags": common.TagsComputedSchema(
							`The key/value pairs to associate with the application.`,
						),
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator name of the application.`,
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
					},
				},
				Description: "All application details.",
			},
		},
	}
}

func queryV3Applications(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/cas/applications?limit=100"
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
		apps := utils.PathSearch("applications", respBody, make([]interface{}, 0)).([]interface{})
		if len(apps) < 1 {
			break
		}
		result = append(result, apps...)
		offset += len(apps)
	}

	return result, nil
}

func flattenV3Applications(applications []interface{}) []map[string]interface{} {
	if len(applications) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(applications))
	for _, application := range applications {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", application, nil),
			"name":                  utils.PathSearch("name", application, nil),
			"description":           utils.PathSearch("description", application, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", application, nil),
			"creator":               utils.PathSearch("creator", application, nil),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
				application, float64(0)).(float64))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time",
				application, float64(0)).(float64))/1000, false),
			"tags": utils.FlattenTagsToMap(utils.PathSearch("labels", application, nil)),
		})
	}

	return result
}

func dataSourceV3ApplicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	appList, err := queryV3Applications(client)
	if err != nil {
		return diag.Errorf("error getting applications: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("applications", flattenV3Applications(appList)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
