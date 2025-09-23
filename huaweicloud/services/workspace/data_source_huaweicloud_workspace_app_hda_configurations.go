package workspace

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v1/{project_id}/app-servers/access-agent/list
func DataSourceAppHdaConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppHdaConfigurationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the HDA configurations are located.`,
			},

			// Optional parameters.
			"server_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the server group to be queried.`,
			},
			"server_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the server to be queried.`,
			},

			// Attributes.
			"configurations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of HDA configurations that match the query parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the server.`,
						},
						"machine_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The machine name of the server.`,
						},
						"maintain_status": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the server is in maintenance status.`,
						},
						"server_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the server.`,
						},
						"server_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the server group.`,
						},
						"server_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the server group.`,
						},
						"sid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The SID of the server.`,
						},
						"session_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of sessions.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the server.`,
						},
						"current_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current version of the access agent.`,
						},
					},
				},
			},
		},
	}
}

func buildAppHdaConfigurationsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("server_group_id"); ok {
		queryParams = fmt.Sprintf("%s&server_group_id=%v", queryParams, v.(string))
	}
	if v, ok := d.GetOk("server_name"); ok {
		queryParams = fmt.Sprintf("%s&server_name=%v", queryParams, v.(string))
	}

	return queryParams
}

func listAppHdaConfigurations(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/app-servers/access-agent/list?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildAppHdaConfigurationsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithPagination := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithPagination, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		configurations := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, configurations...)
		if len(configurations) < limit {
			break
		}
		offset += len(configurations)
	}

	return result, nil
}

func flattenAppHdaConfigurations(configurations []interface{}) []map[string]interface{} {
	if len(configurations) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(configurations))
	for _, configuration := range configurations {
		result = append(result, map[string]interface{}{
			"server_id":         utils.PathSearch("server_id", configuration, nil),
			"machine_name":      utils.PathSearch("machine_name", configuration, nil),
			"maintain_status":   utils.PathSearch("maintain_status", configuration, nil),
			"server_name":       utils.PathSearch("server_name", configuration, nil),
			"server_group_id":   utils.PathSearch("server_group_id", configuration, nil),
			"server_group_name": utils.PathSearch("server_group_name", configuration, nil),
			"sid":               utils.PathSearch("sid", configuration, nil),
			"session_count":     utils.PathSearch("session_count", configuration, nil),
			"status":            utils.PathSearch("status", configuration, nil),
			"current_version":   utils.PathSearch("current_version", configuration, nil),
		})
	}

	return result
}

func dataSourceAppHdaConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	configurations, err := listAppHdaConfigurations(client, d)
	if err != nil {
		return diag.Errorf("error retrieving HDA configurations: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("configurations", flattenAppHdaConfigurations(configurations)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
