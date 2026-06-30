package das

import (
	"context"
	"fmt"
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

// @API DAS GET /v3/{project_id}/connections
func DataSourceDatabaseInstanceConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseInstanceConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the database instance connections are located.`,
			},

			// Optional parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the instance to which the database instance connection belongs.`,
			},
			"network_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The network type of the database instance connection.",
			},
			"datastore_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The datastore type of the database instance connection.",
			},
			"connection_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The connection type of the database instance connection.",
			},
			"condition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The keyword used to search for database instance connection address, name, database username, or remarks.",
			},

			// Attributes.
			"connections": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        databaseInstanceConnectionsElem(),
				Description: `The list of connections that matched filter parameters.`,
			},
		},
	}
}

func databaseInstanceConnectionsElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the database instance connection.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance ID of the database instance connection.",
			},
			"engine_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The engine type of the database instance connection.",
			},
			"network_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network type of the database instance connection.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username of the database instance connection.",
			},
			"is_save_password": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to save the password for the database instance connection.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the database instance connection.",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port of the database instance connection.",
			},
			"database_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database name of the database instance connection.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance name of the database instance connection.",
			},
			"datastore_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The datastore version of the database instance connection.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ip address of the database instance connection.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the database instance connection was created, in RFC3339 format.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the database instance connection.",
			},
			"conn_share_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The conn share type of the database instance connection.",
			},
			"shared_user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The shared user name of the database instance connection.",
			},
			"shared_user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The shared user ID of the database instance connection.",
			},
			"expired_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the database instance connection expires, in RFC3339 format.",
			},
		},
	}
	return &sc
}

func buildDatabaseInstanceConnectionsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}
	if v, ok := d.GetOk("network_type"); ok {
		res = fmt.Sprintf("%s&network_type=%v", res, v)
	}
	if v, ok := d.GetOk("datastore_type"); ok {
		res = fmt.Sprintf("%s&datastore_type=%v", res, v)
	}
	if v, ok := d.GetOk("connection_type"); ok {
		res = fmt.Sprintf("%s&connection_type=%v", res, v)
	}
	if v, ok := d.GetOk("condition"); ok {
		res = fmt.Sprintf("%s&condition=%v", res, v)
	}

	return res
}

func listDatabaseInstanceConnections(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/connections?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildDatabaseInstanceConnectionsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
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
		connections := utils.PathSearch("das_conn_info_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, connections...)
		if len(connections) < limit {
			break
		}
		offset += len(connections)
	}

	return result, nil
}

func flattenDatabaseInstanceConnections(connections []interface{}) []map[string]interface{} {
	if len(connections) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(connections))
	for _, connection := range connections {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("connection_id", connection, nil),
			"instance_id":       utils.PathSearch("instance_id", connection, nil),
			"engine_type":       utils.PathSearch("engine_type", connection, nil),
			"network_type":      utils.PathSearch("network_type", connection, nil),
			"username":          utils.PathSearch("user_name", connection, nil),
			"is_save_password":  utils.PathSearch("is_save_password", connection, nil),
			"description":       utils.PathSearch("remarks", connection, nil),
			"port":              utils.PathSearch("port", connection, float64(0)).(float64),
			"database_name":     utils.PathSearch("database_name", connection, nil),
			"instance_name":     utils.PathSearch("instance_name", connection, nil),
			"datastore_version": utils.PathSearch("datastore_version", connection, nil),
			"ip_address":        utils.PathSearch("ip_address", connection, nil),
			"created_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_at",
				connection, float64(0)).(float64))/1000, false),
			"status":           utils.PathSearch("status", connection, nil),
			"conn_share_type":  utils.PathSearch("conn_share_type", connection, nil),
			"shared_user_name": utils.PathSearch("shared_user_name", connection, nil),
			"shared_user_id":   utils.PathSearch("shared_user_id", connection, nil),
			"expired_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("expired_time",
				connection, float64(0)).(float64))/1000, false),
		})
	}
	return result
}

func dataSourceDatabaseInstanceConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	connections, err := listDatabaseInstanceConnections(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS Database instance connections: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("connections", flattenDatabaseInstanceConnections(connections)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
