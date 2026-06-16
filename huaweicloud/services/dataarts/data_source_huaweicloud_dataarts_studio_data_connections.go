package dataarts

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/data-connections
func DataSourceStudioDataConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceStudioDataConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the data connections are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the data connections belong.`,
			},

			// Optional parameters.
			"connection_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the data connection to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the data connection to be queried.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the data connection to be queried.`,
			},

			// Attributes.
			"connections": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the data connections that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the data connection.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the data connection.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the data connection.`,
						},
						"agent_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The agent ID corresponding to the data connection.`,
						},
						"qualified_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The qualified name of the data connection.`,
						},
						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the data connection.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the data connection, in RFC3339 format.`,
						},
						"create_timestamp": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The creation timestamp of the data connection.`,
						},
					},
				},
			},
		},
	}
}

func buildStudioDataConnectionsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}

	if res != "" {
		return res[1:]
	}
	return res
}

func buildStudioMoreHeaders(workspaceId string) map[string]string {
	result := map[string]string{
		"Content-Type": "application/json",
	}

	if workspaceId != "" {
		result["workspace"] = workspaceId
	}

	return result
}

func listStudioDataConnections(client *golangsdk.ServiceClient, workspaceId string, queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/data-connections?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if len(queryParams) > 0 {
		listPath = fmt.Sprintf("%s&%s", listPath, queryParams[0])
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildStudioMoreHeaders(workspaceId),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		connections := utils.PathSearch("data_connection_lists", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, connections...)
		if len(connections) < limit {
			break
		}
		offset += len(connections)
	}

	return result, nil
}

func filterByConnectionId(connections []interface{}, connectionId string) []interface{} {
	if connectionId == "" {
		return connections
	}

	return utils.PathSearch(fmt.Sprintf("[?dw_id == '%s']", connectionId), connections, make([]interface{}, 0)).([]interface{})
}

func flattenConnections(connections []interface{}) []map[string]interface{} {
	if len(connections) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(connections))
	for _, connection := range connections {
		result = append(result, map[string]interface{}{
			"id":               utils.PathSearch("dw_id", connection, nil),
			"name":             utils.PathSearch("dw_name", connection, nil),
			"type":             utils.PathSearch("dw_type", connection, nil),
			"agent_id":         utils.PathSearch("agent_id", connection, nil),
			"qualified_name":   utils.PathSearch("qualified_name", connection, nil),
			"created_by":       utils.PathSearch("create_user", connection, nil),
			"created_at":       utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", connection, float64(0)).(float64))/1000, false),
			"create_timestamp": int(utils.PathSearch("create_time", connection, float64(0)).(float64)),
		})
	}

	return result
}

func resourceStudioDataConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		workspaceId  = d.Get("workspace_id").(string)
		connectionId = d.Get("connection_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	connections, err := listStudioDataConnections(client, workspaceId, buildStudioDataConnectionsQueryParams(d))
	if err != nil {
		return diag.Errorf("error retrieving data connections: %s", err)
	}
	log.Printf("[Lance] The list response is: %v", connections)

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("connections", flattenConnections(filterByConnectionId(connections, connectionId))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
