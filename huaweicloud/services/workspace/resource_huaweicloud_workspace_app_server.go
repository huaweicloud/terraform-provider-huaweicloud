package workspace

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v1/{project_id}/app-servers/actions/create
// @API Workspace GET /v2/{project_id}/job/{job_id}
// @API Workspace GET /v1/{project_id}/app-servers/{server_id}
// @API Workspace PATCH /v1/{project_id}/app-servers/{server_id}
// @API Workspace DELETE /v1/{project_id}/app-servers/{server_id}
func ResourceAppServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppServerCreate,
		ReadContext:   resourceAppServerRead,
		UpdateContext: resourceAppServerUpdate,
		DeleteContext: resourceAppServerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The server group ID to which the server belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the server.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The type of the server.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The flavor ID of the server.`,
			},
			"root_volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The disk type of the server.",
						},
						"size": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "The disk size of the server, in GB.",
						},
					},
				},
				Description: `The system disk configuration of the server.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The VPC ID to which the server belongs.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The subnet ID to which the server belongs.`,
			},
			"os_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The operating system type of the server.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The availability zone of the server.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the server.`,
			},
			"maintain_status": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable maintenance mode.`,
			},
			"ou_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The OU name corresponding to the AD server.`,
			},
			"update_access_agent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether to automatically upgrade protocol component.`,
			},
			"scheduler_hints": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `The configuration of the dedicate host.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dedicated_host_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The ID of the dedicate host.`,
						},
						"tenancy": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The type of the dedicate host.`,
						},
					},
				},
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
		},
	}
}

func buildCreateServerBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"server_group_id":     d.Get("server_group_id"),
		"type":                d.Get("type"),
		"product_id":          d.Get("flavor_id"),
		"subscription_num":    1,
		"root_volume":         buildAppServerRootVolume(d.Get("root_volume").([]interface{})),
		"vpc_id":              d.Get("vpc_id"),
		"subnet_id":           d.Get("subnet_id"),
		"os_type":             utils.ValueIgnoreEmpty(d.Get("os_type")),
		"availability_zone":   utils.ValueIgnoreEmpty(d.Get("availability_zone")),
		"update_access_agent": d.Get("update_access_agent"),
		"ou_name":             utils.ValueIgnoreEmpty(d.Get("ou_name")),
		"scheduler_hints":     buildAppServerSchedulerHints(d.Get("scheduler_hints").([]interface{})),
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		params["create_server_extend_param"] = buildPrepaidOptionsBodyParams(d)
	}

	return params
}

func buildAppServerSchedulerHints(schedulerHint []interface{}) map[string]interface{} {
	if len(schedulerHint) == 0 || schedulerHint[0] == nil {
		return nil
	}

	return map[string]interface{}{
		"dedicated_host_id": utils.PathSearch("dedicated_host_id", schedulerHint[0], nil),
		"tenancy":           utils.PathSearch("tenancy", schedulerHint[0], nil),
	}
}

func buildAppServerRootVolume(rootVolume []interface{}) map[string]interface{} {
	if len(rootVolume) == 0 {
		return nil
	}

	return map[string]interface{}{
		"type": utils.PathSearch("type", rootVolume[0], nil),
		"size": utils.PathSearch("size", rootVolume[0], nil),
	}
}

func buildPrepaidOptionsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"charging_mode": "prePaid",
		"period_type":   buildCreatePeriodTypeParam(d.Get("period_unit").(string)),
		"period_num":    d.Get("period"),
		"is_auto_renew": parseAutoRenewValue(d.Get("auto_renew").(string)),
		"is_auto_pay":   true,
	}
}

func buildCreatePeriodTypeParam(periodUnit string) interface{} {
	switch periodUnit {
	case "month":
		return 2
	case "year":
		return 3
	}

	return nil
}

func parseAutoRenewValue(autoRenew string) bool {
	result, err := strconv.ParseBool(autoRenew)
	if err != nil {
		log.Printf("[ERROR] unable to convert auto_renew to bool value")
		return false
	}
	return result
}

func resourceAppServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	httpUrl := "v1/{project_id}/app-servers/actions/create"
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	serverGroupId := d.Get("server_group_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateServerBodyParams(d)),
	}
	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Workspace APP server under specified server group (%s): %s", serverGroupId, err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		orderId := utils.PathSearch("order_id", respBody, "").(string)
		resourceId, err := waitForPrePaidServerComplete(ctx, cfg, region, d.Timeout(schema.TimeoutCreate), orderId)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resourceId)
	} else {
		jobId := utils.PathSearch("job_id", respBody, "").(string)
		if jobId == "" {
			return diag.Errorf("unable to find job ID from API response")
		}

		serverResp, err := waitForAppServerJobCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), jobId)
		if err != nil {
			return diag.Errorf("error waiting for the job (%s) completed: %s", jobId, err)
		}

		serverId := utils.PathSearch("sub_jobs|[0].job_resource_info.resource_id", serverResp, "").(string)
		if serverId == "" {
			return diag.Errorf("unable to find server ID from API response")
		}

		d.SetId(serverId)
	}

	if err := updateAppServer(client, d, d.Id()); err != nil {
		return diag.Errorf("error updating the server (%s): %s", d.Id(), err)
	}

	return resourceAppServerRead(ctx, d, meta)
}

func waitForPrePaidServerComplete(ctx context.Context, cfg *config.Config, region string, timeOut time.Duration, orderId string) (string, error) {
	if orderId == "" {
		return "", fmt.Errorf("unable to find order ID from API response")
	}

	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return "", fmt.Errorf("error creating BSS v2 client: %s", err)
	}

	if err := common.WaitOrderComplete(ctx, bssClient, orderId, timeOut); err != nil {
		return "", err
	}

	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, timeOut)
	if err != nil {
		return "", fmt.Errorf("error waiting for Workspace APP server order (%s) complete: %s", orderId, err)
	}

	return resourceId, nil
}

func waitForAppServerJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration, jobId string) (interface{},
	error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"SUCCESS"},
		Refresh:      refreshAppServerJobStatusFunc(client, jobId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 30 * time.Second,
	}

	serverResp, err := stateConf.WaitForStateContext(ctx)
	return serverResp, err
}

func refreshAppServerJobStatusFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		httpUrl := "v2/{project_id}/job/{job_id}"
		getJobPath := client.Endpoint + httpUrl
		getJobPath = strings.ReplaceAll(getJobPath, "{project_id}", client.ProjectID)
		getJobPath = strings.ReplaceAll(getJobPath, "{job_id}", jobId)
		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		resp, err := client.Request("GET", getJobPath, &getOpt)
		if err != nil {
			return resp, "ERROR", err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return resp, "ERROR", err
		}

		return respBody, utils.PathSearch("status", respBody, nil).(string), nil
	}
}

func resourceAppServerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	server, err := GetServerById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Workspace APP server")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("server_group_id", utils.PathSearch("server_group_id", server, nil)),
		d.Set("name", utils.PathSearch("name", server, nil)),
		d.Set("flavor_id", utils.PathSearch("product_info.product_id", server, nil)),
		d.Set("root_volume", flattenAppServerRootVolume(utils.PathSearch("product_info", server, nil))),
		d.Set("os_type", utils.PathSearch("os_type", server, nil)),
		d.Set("charging_mode", flattenServerChargingMode(server)),
		d.Set("availability_zone", utils.PathSearch("availability_zone", server, nil)),
		d.Set("description", utils.PathSearch("description", server, nil)),
		d.Set("ou_name", utils.PathSearch("ou_name", server, nil)),
		d.Set("maintain_status", utils.PathSearch("maintain_status", server, nil)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenServerChargingMode(getRespBody interface{}) string {
	chargingMode := utils.PathSearch("metadata.charging_mode", getRespBody, "").(string)
	if chargingMode == "1" {
		return "prePaid"
	}
	if chargingMode == "0" {
		return "postPaid"
	}

	log.Printf("[WARN] error parsing charging_mode from API response")
	return ""
}

func flattenAppServerRootVolume(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	diskSize := utils.PathSearch("system_disk_size", resp, "").(string)
	return []map[string]interface{}{
		{
			"type": utils.PathSearch("system_disk_type", resp, nil),
			"size": utils.StringToInt(&diskSize),
		},
	}
}

// GetServerById is amethod used to query server detail by server ID.
func GetServerById(client *golangsdk.ServiceClient, serverId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/app-servers/{server_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{server_id}", serverId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceAppServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	serverId := d.Id()
	if err := updateAppServer(client, d, serverId); err != nil {
		return diag.Errorf("error updating server (%s): %s", serverId, err)
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), serverId); err != nil {
			return diag.Errorf("error updating the auto-renew of the server (%s): %s", serverId, err)
		}
	}

	return resourceAppServerRead(ctx, d, meta)
}

func updateAppServer(client *golangsdk.ServiceClient, d *schema.ResourceData, serverId string) error {
	updateServerOpt := map[string]interface{}{}
	if d.HasChange("name") {
		updateServerOpt["name"] = d.Get("name")
	}

	if d.HasChange("description") {
		updateServerOpt["description"] = d.Get("description")
	}

	if d.HasChange("maintain_status") {
		updateServerOpt["maintain_status"] = d.Get("maintain_status")
	}

	if len(updateServerOpt) == 0 {
		return nil
	}

	httpUrl := "v1/{project_id}/app-servers/{server_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{server_id}", serverId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         updateServerOpt,
	}

	_, err := client.Request("PATCH", updatePath, &updateOpt)
	if err != nil {
		return err
	}

	return nil
}

func resourceAppServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	serverId := d.Id()
	if d.Get("charging_mode").(string) == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
			return diag.Errorf("error unsubscribing Workspace APP server (%s) under speficied server group (%s): %s",
				serverId, d.Get("server_group_id").(string), err)
		}
	} else {
		if err := deleteAppServerById(ctx, client, serverId, d.Timeout(schema.TimeoutDelete)); err != nil {
			return common.CheckDeletedDiag(d, err, "delete Workspace APP server")
		}
	}

	if err := waitingForServerDeleteCompleted(ctx, client, serverId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for Workspace APP server (%s) deleted: %s", d.Id(), err)
	}

	return nil
}

func deleteAppServerById(ctx context.Context, client *golangsdk.ServiceClient, serverId string, timeout time.Duration) error {
	httpUrl := "v1/{project_id}/app-servers/{server_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{server_id}", serverId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// When deleting a non-existent server, the response status code is 200.
	requestResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		log.Printf("[ERROR] Unable to find job ID from API response")
		return nil
	}

	_, err = waitForAppServerJobCompleted(ctx, client, timeout, jobId)
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) completed: %s", jobId, err)
	}
	return nil
}

func waitingForServerDeleteCompleted(ctx context.Context, client *golangsdk.ServiceClient, serverId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetServerById(client, serverId)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return "deleted", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}
			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
