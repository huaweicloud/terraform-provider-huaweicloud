package das

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	databaseInstanceConnectionNonUpdatableParams = []string{
		"instance_id",
		"engine_type",
		"network_type",
	}
)

// @API DAS POST /v3/{project_id}/connections
// @API DAS GET /v3/{project_id}/connections/{connection_id}
// @API DAS PUT /v3/{project_id}/connections/{connection_id}
// @API DAS DELETE /v3/{project_id}/batch-delete-connections
func ResourceDatabaseInstanceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseInstanceConnectionCreate,
		ReadContext:   resourceDatabaseInstanceConnectionRead,
		UpdateContext: resourceDatabaseInstanceConnectionUpdate,
		DeleteContext: resourceDatabaseInstanceConnectionDelete,

		CustomizeDiff: config.FlexibleForceNew(databaseInstanceConnectionNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the database instance connection is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the instance to which the database instance connection belongs.`,
			},
			"engine_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The engine type of the database instance connection.",
			},
			"network_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network type of the database instance connection.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username of the database instance connection.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The password of the database instance connection.",
			},
			"is_save_password": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to save the password for the database instance connection.",
			},

			// Optional parameters.
			"node_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The unique identifiers of the instance nodes.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the database instance connection.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The port of the database instance connection.",
			},
			"database_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The database name of the database instance connection.",
			},
			"sql_record_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether SQL recording is enabled for the database instance connection.",
			},

			// Attributes.
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource ID, in UUID format.",
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
				Description: "The timestamp when the database instance connection was created, in RFC3339 format.",
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
				Description: "The timestamp when the database instance connection expires, in RFC3339 format.",
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildCreateDatabaseInstanceConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"instance_id":      d.Get("instance_id"),
		"engine_type":      d.Get("engine_type"),
		"network_type":     d.Get("network_type"),
		"username":         d.Get("username"),
		"password":         d.Get("password"),
		"is_save_password": d.Get("is_save_password"),
		"node_ids":         d.Get("node_ids"),
		"remarks":          d.Get("description"),
		"port":             utils.ValueIgnoreEmpty(d.Get("port")),
		"database_name":    d.Get("database_name"),
		"sql_record_flag":  d.Get("sql_record_flag"),
	}
}

func createDatabaseInstanceConnection(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v3/{project_id}/connections"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createDatabaseInstanceConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
		JSONBody: utils.RemoveNil(buildCreateDatabaseInstanceConnectionBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createDatabaseInstanceConnectionOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(resp)
}

func resourceDatabaseInstanceConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	respBody, err := createDatabaseInstanceConnection(client, d)
	if err != nil {
		return diag.Errorf("error creating DAS Database instance connection: %s", err)
	}

	connectionId := utils.PathSearch("connection_id", respBody, "").(string)
	if connectionId == "" {
		return diag.Errorf("unable to find the ID of the DAS Database instance connection from the API response")
	}
	d.SetId(connectionId)

	return resourceDatabaseInstanceConnectionRead(ctx, d, meta)
}

func GetDatabaseInstanceConnectionById(client *golangsdk.ServiceClient, connectionId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/connections/{connection_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{connection_id}", connectionId)
	getDatabaseInstanceConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
	}

	resp, err := client.Request("GET", getPath, &getDatabaseInstanceConnectionOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(resp)
}

func resourceDatabaseInstanceConnectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		connectionId = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	resp, err := GetDatabaseInstanceConnectionById(client, connectionId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving DAS Database instance connection (%s)", connectionId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("das_conn_info.instance_id", resp, nil)),
		d.Set("engine_type", utils.PathSearch("das_conn_info.engine_type", resp, nil)),
		d.Set("network_type", utils.PathSearch("das_conn_info.network_type", resp, nil)),
		d.Set("username", utils.PathSearch("das_conn_info.user_name", resp, nil)),
		d.Set("is_save_password", utils.PathSearch("das_conn_info.is_save_password", resp, nil)),
		d.Set("description", utils.PathSearch("das_conn_info.remarks", resp, nil)),
		d.Set("port", utils.PathSearch("das_conn_info.port", resp, float64(0)).(float64)),
		d.Set("database_name", utils.PathSearch("das_conn_info.database_name", resp, nil)),
		d.Set("instance_name", utils.PathSearch("das_conn_info.instance_name", resp, nil)),
		d.Set("datastore_version", utils.PathSearch("das_conn_info.datastore_version", resp, nil)),
		d.Set("ip_address", utils.PathSearch("das_conn_info.ip_address", resp, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("das_conn_info.create_at",
			resp, float64(0)).(float64))/1000, false)),
		d.Set("status", utils.PathSearch("das_conn_info.status", resp, nil)),
		d.Set("conn_share_type", utils.PathSearch("das_conn_info.conn_share_type", resp, nil)),
		d.Set("shared_user_name", utils.PathSearch("das_conn_info.shared_user_name", resp, nil)),
		d.Set("shared_user_id", utils.PathSearch("das_conn_info.shared_user_id", resp, nil)),
		d.Set("expired_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("das_conn_info.expired_time",
			resp, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateDatabaseInstanceConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"username":         d.Get("username"),
		"is_save_password": d.Get("is_save_password"),
		"password":         d.Get("password"),
		"node_ids":         d.Get("node_ids"),
		"remarks":          d.Get("description"),
		"port":             utils.ValueIgnoreEmpty(d.Get("port")),
		"database_name":    d.Get("database_name"),
		"sql_record_flag":  d.Get("sql_record_flag"),
	}
}

func updateDatabaseInstanceConnection(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v3/{project_id}/connections/{connection_id}"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{connection_id}", d.Id())

	updateDatabaseInstanceConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
		JSONBody: utils.RemoveNil(buildUpdateDatabaseInstanceConnectionBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateDatabaseInstanceConnectionOpt)
	return err
}

func resourceDatabaseInstanceConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS Client: %s", err)
	}

	err = updateDatabaseInstanceConnection(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDatabaseInstanceConnectionRead(ctx, d, meta)
}

func deleteDatabaseInstanceConnection(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v3/{project_id}/batch-delete-connections"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{connection_id}", d.Id())

	deleteDatabaseInstanceConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
		JSONBody: map[string]interface{}{
			"connection_ids": []interface{}{d.Id()},
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteDatabaseInstanceConnectionOpt)
	return err
}

func resourceDatabaseInstanceConnectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		connectionId = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	err = deleteDatabaseInstanceConnection(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting DAS Database instance connection (%s)", connectionId))
	}

	return nil
}
