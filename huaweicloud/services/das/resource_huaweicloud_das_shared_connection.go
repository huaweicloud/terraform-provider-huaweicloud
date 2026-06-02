package das

import (
	"context"
	"fmt"
	"strconv"
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
	sharedConnectionNonUpdatableParams = []string{
		"connection_id",
		"user_id",
		"user_name",
		"expired_at",
	}

	sharedConnectionNotFoundCodes = []string{
		"DAS.5052", // shared connection does not exist during deletion.
	}
)

// @API DAS POST /v3/{project_id}/connections/share
// @API DAS GET /v3/{project_id}/connections/{connection_id}/get-shared-list
// @API DAS DELETE /v3/{project_id}/connections/share
func ResourceSharedConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedConnectionCreate,
		ReadContext:   resourceSharedConnectionRead,
		UpdateContext: resourceSharedConnectionUpdate,
		DeleteContext: resourceSharedConnectionDelete,

		CustomizeDiff: config.FlexibleForceNew(sharedConnectionNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the shared connection is located.`,
			},

			// Required parameters.
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the connection to which the shared connection belongs.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user ID of the shared connection.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user name of the shared connection.",
			},

			// Optional parameters.
			"expired_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The expiration time of the shared connection, in RFC3339 format.",
			},

			// Attributes.
			"shared_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the shared connection, in RFC3339 format.",
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

func buildCreateSharedConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	users := []map[string]interface{}{
		{
			"user_id":   d.Get("user_id"),
			"user_name": d.Get("user_name"),
		},
	}

	expiredTime := d.Get("expired_at")
	if expiredAtStr := d.Get("expired_at").(string); expiredAtStr != "" {
		timestamp := utils.ConvertTimeStrToNanoTimestamp(expiredAtStr)
		expiredTime = utils.FormatTimeStampRFC3339(timestamp/1000, true, "2006-01-02T15:04:05.000Z")
	}

	return map[string]interface{}{
		"shared_conn_id": d.Get("connection_id"),
		"users":          users,
		"expired_time":   expiredTime,
	}
}

func createSharedConnection(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v3/{project_id}/connections/share"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createSharedConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
		JSONBody: buildCreateSharedConnectionBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createSharedConnectionOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(resp)
}

func resourceSharedConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	respBody, err := createSharedConnection(client, d)
	if err != nil {
		return diag.Errorf("error creating DAS shared connection: %s", err)
	}

	status, ok := utils.PathSearch("status", respBody, false).(bool)
	if !ok || !status {
		return diag.Errorf("failed to create DAS shared connection: expected 'status' to be true in API response")
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("connection_id").(string), d.Get("user_id").(string)))

	return resourceSharedConnectionRead(ctx, d, meta)
}

func listSharedConnections(client *golangsdk.ServiceClient, connectionId, queryParams string) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/connections/{connection_id}/get-shared-list?perPage"
		perPage = 100
		curPage = 1
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{connection_id}", connectionId)
	listPath = strings.ReplaceAll(listPath, "{perPage}", strconv.Itoa(perPage))
	if queryParams != "" {
		listPath += queryParams
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
	}

	for {
		requestURL := fmt.Sprintf("%s&cur_page=%d", listPath, curPage)

		requestResp, err := client.Request("GET", requestURL, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		sharedConnections := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, sharedConnections...)

		if len(sharedConnections) < perPage {
			break
		}

		curPage++
	}

	return result, nil
}

func GetSharedConnectionById(client *golangsdk.ServiceClient, connectionId, userId string) (interface{}, error) {
	sharedConnections, err := listSharedConnections(client, connectionId, fmt.Sprintf("&keywords=%s", userId))
	if err != nil {
		return nil, err
	}

	if len(sharedConnections) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/{project_id}/connections/{connection_id}/get-shared-list",
				RequestId: "NONE",
				Body: []byte(fmt.Sprintf("DAS shared connection (connection_id: %s, user_id: %s) not found in list response",
					connectionId, userId)),
			},
		}
	}

	return sharedConnections[0], nil
}

func resourceSharedConnectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		parts        = strings.Split(d.Id(), "/")
		connectionId = parts[0]
		userId       = parts[1]
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	resp, err := GetSharedConnectionById(client, connectionId, userId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error retrieving DAS shared connection (connection_id: %s, user_id: %s)", connectionId, userId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("user_id", utils.PathSearch("user_id", resp, nil)),
		d.Set("user_name", utils.PathSearch("user_name", resp, nil)),
		d.Set("expired_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("expired_time",
			resp, float64(0)).(float64))/1000, false)),
		d.Set("shared_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("shared_time",
			resp, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSharedConnectionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func buildDeleteSharedConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	users := []map[string]interface{}{
		{
			"user_id":   d.Get("user_id"),
			"user_name": d.Get("user_name"),
		},
	}

	return map[string]interface{}{
		"shared_conn_id": d.Get("connection_id"),
		"users":          users,
	}
}

func deleteSharedConnection(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v3/{project_id}/connections/share"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	deleteSharedConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
		JSONBody: buildDeleteSharedConnectionBodyParams(d),
	}

	_, err := client.Request("DELETE", deletePath, &deleteSharedConnectionOpt)
	return err
}

func resourceSharedConnectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	err = deleteSharedConnection(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errorCodeStr", sharedConnectionNotFoundCodes...),
			fmt.Sprintf("error retrieving DAS shared connection (connection_id: %s, user_id: %s)",
				d.Get("connection_id").(string), d.Get("user_id").(string)))
	}

	return nil
}
