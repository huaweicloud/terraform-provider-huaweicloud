package dataarts

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

// @API DataArtsStudio GET /v1/{project_id}/service/servicecatalogs/{catalog_id}/apis
func DataSourceDataServiceCatalogApis() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataServiceCatalogApisRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the catalog APIs are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the catalog belongs.`,
			},
			"catalog_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the catalog.`,
			},

			// Attributes.
			"apis": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the catalog APIs that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the API.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the API.`,
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the API group.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the API.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the API.`,
						},
						"debug_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The debug status of the API.`,
						},
						"publish_messages": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The publish information list of the API.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the published message.`,
									},
									"api_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the API.`,
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the instance.`,
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the instance.`,
									},
									"api_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The status of the API.`,
									},
									"api_debug": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The debug status of the API.`,
									},
								},
							},
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the API.`,
						},
						"manager": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The manager of the API.`,
						},
						"create_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the API.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the API.`,
						},
						"authorization_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The authorization status of the API.`,
						},
					},
				},
			},
		},
	}
}

func flattenDataServiceCatalogApiPublishMessages(publishMessages []interface{}) []map[string]interface{} {
	if len(publishMessages) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(publishMessages))
	for _, msg := range publishMessages {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", msg, nil),
			"api_id":        utils.PathSearch("api_id", msg, nil),
			"instance_id":   utils.PathSearch("instance_id", msg, nil),
			"instance_name": utils.PathSearch("instance_name", msg, nil),
			"api_status":    utils.PathSearch("api_status", msg, nil),
			"api_debug":     utils.PathSearch("api_debug", msg, nil),
		})
	}

	return result
}

func flattenDataServiceCatalogApis(apis []interface{}) []map[string]interface{} {
	if len(apis) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(apis))
	for _, api := range apis {
		result = append(result, map[string]interface{}{
			"id":                   utils.PathSearch("id", api, nil),
			"name":                 utils.PathSearch("name", api, nil),
			"group_id":             utils.PathSearch("group_id", api, nil),
			"description":          utils.PathSearch("description", api, nil),
			"status":               utils.PathSearch("status", api, nil),
			"debug_status":         utils.PathSearch("debug_status", api, nil),
			"authorization_status": utils.PathSearch("authorization_status", api, nil),
			"type":                 utils.PathSearch("type", api, nil),
			"manager":              utils.PathSearch("manager", api, nil),
			"create_user":          utils.PathSearch("create_user", api, nil),
			"publish_messages": flattenDataServiceCatalogApiPublishMessages(utils.PathSearch("publish_messages", api,
				make([]interface{}, 0)).([]interface{})),
			"create_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", api,
				float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func listDataServiceCatalogApis(client *golangsdk.ServiceClient, catalogId, workspaceId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/service/servicecatalogs/{catalog_id}/apis?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{catalog_id}", catalogId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    workspaceId,
			"Dlm-Type":     "EXCLUSIVE",
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

		apis := utils.PathSearch("apis", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, apis...)
		if len(apis) < limit {
			break
		}

		offset += len(apis)
	}

	return result, nil
}

func dataSourceDataServiceCatalogApisRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	catalogId := d.Get("catalog_id").(string)
	workspaceId := d.Get("workspace_id").(string)

	apis, err := listDataServiceCatalogApis(client, catalogId, workspaceId)
	if err != nil {
		return diag.Errorf("error querying catalog APIs: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("apis", flattenDataServiceCatalogApis(apis)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
