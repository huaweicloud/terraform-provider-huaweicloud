package esw

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var connectionNonUpdatableParams = []string{"instance_id", "vpc_id", "virsubnet_id", "remote_infos", "fixed_ips"}

// @API ESW POST /v3/{project_id}/l2cg/instances/{instance_id}/connections
// @API ESW GET /v3/{project_id}/l2cg/instances/{instance_id}/connections/{connection_id}
// @API ESW PUT /v3/{project_id}/l2cg/instances/{instance_id}/connections/{connection_id}
// @API ESW DELETE /v3/{project_id}/l2cg/instances/{instance_id}/connections/{connection_id}
func ResourceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionCreate,
		ReadContext:   resourceConnectionRead,
		UpdateContext: resourceConnectionUpdate,
		DeleteContext: resourceConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceEswConnectionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(connectionNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"virsubnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remote_infos": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     connectionRemoteInfosSchema(),
			},
			"fixed_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func connectionRemoteInfosSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"segmentation_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"tunnel_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tunnel_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"tunnel_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/l2cg/instances/{instance_id}/connections"
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ESW client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateConnectionBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ESW connection: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	connectionId := utils.PathSearch("connection.id", createRespBody, "").(string)
	if connectionId == "" {
		return diag.Errorf("unable to find the ESW connection ID from the API response")
	}
	d.SetId(connectionId)

	err = waitForConnectionReady(ctx, client, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceConnectionRead(ctx, d, meta)
}

func buildCreateConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":         d.Get("name"),
		"vpc_id":       d.Get("vpc_id"),
		"virsubnet_id": d.Get("virsubnet_id"),
		"remote_infos": buildCreateConnectionRemoteInfosBody(d),
		"fixed_ips":    utils.ValueIgnoreEmpty(d.Get("fixed_ips").(*schema.Set).List()),
	}
	return map[string]interface{}{
		"connection": bodyParams,
	}
}

func buildCreateConnectionRemoteInfosBody(d *schema.ResourceData) []interface{} {
	rawParams := d.Get("remote_infos").(*schema.Set)
	if rawParams.Len() == 0 {
		return nil
	}

	rst := make([]interface{}, 0, rawParams.Len())
	for _, rawParam := range rawParams.List() {
		if remoteInfo, ok := rawParam.(map[string]interface{}); ok {
			rst = append(rst, map[string]interface{}{
				"segmentation_id": remoteInfo["segmentation_id"],
				"tunnel_ip":       remoteInfo["tunnel_ip"],
				"tunnel_port":     utils.ValueIgnoreEmpty(remoteInfo["tunnel_port"]),
			})
		}
	}

	return rst
}

func resourceConnectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ESW client: %s", err)
	}

	connection, err := getConnectionIdById(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ESW connection")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("connection.instance_id", connection, nil)),
		d.Set("name", utils.PathSearch("connection.name", connection, nil)),
		d.Set("vpc_id", utils.PathSearch("connection.vpc_id", connection, nil)),
		d.Set("virsubnet_id", utils.PathSearch("connection.virsubnet_id", connection, nil)),
		d.Set("remote_infos", flattenConnectionRemoteInfos(connection)),
		d.Set("fixed_ips", utils.PathSearch("connection.fixed_ips", connection, nil)),
		d.Set("status", utils.PathSearch("connection.status", connection, nil)),
		d.Set("created_at", utils.PathSearch("connection.created_at", connection, nil)),
		d.Set("updated_at", utils.PathSearch("connection.updated_at", connection, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConnectionRemoteInfos(resp interface{}) []interface{} {
	curJson := utils.PathSearch("connection.remote_infos", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"segmentation_id": utils.PathSearch("segmentation_id", v, nil),
			"tunnel_ip":       utils.PathSearch("tunnel_ip", v, nil),
			"tunnel_port":     utils.PathSearch("tunnel_port", v, nil),
			"tunnel_type":     utils.PathSearch("tunnel_type", v, nil),
		})
	}
	return rst
}

func getConnectionIdById(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/l2cg/instances/{instance_id}/connections/{connection_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{connection_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func resourceConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/l2cg/instances/{instance_id}/connections/{connection_id}"
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ESW Client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{connection_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateConnectionBodyParams(d)

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating ESW connection: %s", err)
	}

	return resourceConnectionRead(ctx, d, meta)
}

func buildUpdateConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": d.Get("name"),
	}
	return map[string]interface{}{
		"connection": bodyParams,
	}
}

func resourceConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/l2cg/instances/{instance_id}/connections/{connection_id}"
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ESW client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{connection_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting ESW connection")
	}

	err = waitForConnectionDeleted(ctx, client, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForConnectionReady(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Ready"},
		Refresh:      connectionStatusRefreshFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for ESW connection(%s) to ready: %s ", d.Id(), err)
	}
	return nil
}

func waitForConnectionDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Deleted"},
		Refresh:      connectionStatusRefreshFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for ESW connection(%s) to be deleted: %s ", d.Id(), err)
	}
	return nil
}

func connectionStatusRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		connection, err := getConnectionIdById(client, d)
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				return "", "Deleted", nil
			}
			return nil, "Failed", err
		}

		status := utils.PathSearch("connection.status", connection, "").(string)
		if status == "" {
			return nil, "Failed", errors.New("status is not found")
		}

		if utils.StrSliceContains([]string{"failed", "abnormal"}, status) {
			return connection, "Failed", fmt.Errorf("the connection status is: %s", status)
		}
		if utils.StrSliceContains([]string{"connected", "disconnect"}, status) {
			return connection, "Ready", nil
		}

		return connection, "Pending", nil
	}
}

func resourceEswConnectionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
