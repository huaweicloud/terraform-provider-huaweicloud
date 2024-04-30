package dataarts

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/dayu/v1/connections"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/data-connections
func DataSourceDataConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDataConnectionsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the workspace to which the data connection belongs.`,
			},
			"connection_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the data connection.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the data connection.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the data connection.`,
			},
			"connections": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the data connections.`,
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
							Description: `The creation time of the data connection.`,
						},
					},
				},
			},
		},
	}
}

func resourceDataConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DataArtsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio v1 client: %s", err)
	}

	// When limit and offset are used together, limit cannot be equal to the total number of connections.
	// If limit is not specified, all connections are queried by default. A maximum of 200 data connections can be created.
	opts := connections.ListOpts{
		WorkspaceId: d.Get("workspace_id").(string),
		Name:        d.Get("name").(string),
		Type:        d.Get("type").(string),
	}

	resp, err := connections.List(client, opts)
	if err != nil {
		return diag.Errorf("error retrieving data connections: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("connections", filterByConnectionId(flattenConnections(resp), d.Get("connection_id").(string))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConnections(resp []connections.Connection) []map[string]interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(resp))
	for i, connection := range resp {
		result[i] = map[string]interface{}{
			"id":             connection.DwId,
			"name":           connection.DwName,
			"type":           connection.DwType,
			"agent_id":       connection.AgentId,
			"qualified_name": connection.QualifiedName,
			"created_by":     connection.CreateUser,
			"created_at":     utils.FormatTimeStampRFC3339(int64(connection.CreateTime), false),
		}
	}
	return result
}

func filterByConnectionId(all []map[string]interface{}, connectionId string) []map[string]interface{} {
	if connectionId == "" {
		return all
	}

	rst := make([]map[string]interface{}, 0, len(all))
	for _, v := range all {
		if connectionId == fmt.Sprint(utils.PathSearch("id", v, nil)) {
			rst = append(rst, v)
			return rst
		}
	}
	return rst
}
