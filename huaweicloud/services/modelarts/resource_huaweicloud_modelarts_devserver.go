package modelarts

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

var devServerNotExistCode = "ModelArts.6404" // DevServer does not exist.

// @API ModelArts POST /v1/{project_id}/dev-servers
// @API ModelArts GET /v1/{project_id}/dev-servers
// @API ModelArts GET /v1/{project_id}/dev-servers/{id}
// @API ModelArts DELETE /v1/{project_id}/dev-servers/{id}
func ResourceDevServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDevServerCreate,
		ReadContext:   resourceDevServerRead,
		UpdateContext: resourceDevServerUpdate,
		DeleteContext: resourceDevServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the DevServer.`,
			},
			"flavor": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The flavor of the DevServer.`,
			},
			"architecture": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The architecture of the DevServer.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the VPC to which the DevServer belongs.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the subnet to which the DevServer belongs.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of security group to which the DevServer belongs.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The image ID of the DevServer.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The type of the DevServer.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The availability zone where the DevServer is located.`,
			},
			"admin_pass": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ForceNew:      true,
				ConflictsWith: []string{"key_pair_name"},
				Description:   `The login password for logging in to the server.`,
			},
			"key_pair_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The key pair name for logging in to the server.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The enterprise project ID to which the DevServer belongs.`,
			},
			"root_volume": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: `The size of system disk.`,
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The type of system disk.`,
						},
					},
				},
				Description: `The system disk configuration of the DevServer.`,
			},
			"ipv6_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Whether to enable IPv6.`,
			},
			"roce_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The RoCE network ID of the DevServer.`,
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The user data defined for the server.`,
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The charging mode of the DevServer.`,
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"period"},
				Description:  `The charging period unit of the DevServer.`,
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"period_unit"},
				Description:  `The charging period of the DevServer.`,
			},
			"auto_renew": common.SchemaAutoRenewUpdatable(nil),
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the DevServer, in RFC3339 format.`,
			},
		},
	}
}

func buildCreateDevServerBodyParams(d *schema.ResourceData, epsId string) map[string]interface{} {
	params := map[string]interface{}{
		"name":                  d.Get("name"),
		"flavor":                d.Get("flavor"),
		"network":               buildDevServerNetwork(d),
		"image_id":              d.Get("image_id"),
		"server_type":           utils.ValueIgnoreEmpty(d.Get("type")),
		"root_volume":           buildDevServerRootVolume(d.Get("root_volume").([]interface{})),
		"admin_pass":            utils.ValueIgnoreEmpty(d.Get("admin_pass")),
		"arch":                  utils.ValueIgnoreEmpty(d.Get("architecture")),
		"availability_zone":     utils.ValueIgnoreEmpty(d.Get("availability_zone")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(epsId),
		"key_pair_name":         utils.ValueIgnoreEmpty(d.Get("key_pair_name")),
		"userdata":              utils.ValueIgnoreEmpty(d.Get("user_data")),
	}

	if d.Get("charging_mode").(string) == "PRE_PAID" {
		params["charging_info"] = buildPrepaidOptionsBodyParams(d)
	}

	return params
}

func buildDevServerNetwork(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"security_group_id": d.Get("security_group_id"),
		"vpc_id":            d.Get("vpc_id"),
		"subnet_id":         d.Get("subnet_id"),
		"ipv6_enable":       d.Get("ipv6_enable"),
		"roce_id":           utils.ValueIgnoreEmpty(d.Get("roce_id")),
	}
}

func buildDevServerRootVolume(rootVolume []interface{}) map[string]interface{} {
	if len(rootVolume) == 0 || rootVolume[0] == nil {
		return nil
	}

	return map[string]interface{}{
		"type": utils.PathSearch("type", rootVolume[0], nil),
		"size": utils.PathSearch("size", rootVolume[0], nil),
	}
}

func buildPrepaidOptionsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"charging_mode": d.Get("charging_mode"),
		"period_type":   d.Get("period_unit"),
		"period_num":    d.Get("period"),
		"is_auto_renew": parseAutoRenewValue(d.Get("auto_renew").(string)),
		"is_auto_pay":   true,
	}
}

func parseAutoRenewValue(autoRenew string) bool {
	result, err := strconv.ParseBool(autoRenew)
	if err != nil {
		log.Printf("[ERROR] unable to convert auto_renew to bool value")
		return false
	}
	return result
}

func resourceDevServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dev-servers"
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDevServerBodyParams(d, cfg.GetEnterpriseProjectID(d))),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ModelArts DevServer: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("charging_mode").(string) == "PRE_PAID" {
		bssClient, err := cfg.NewServiceClient("bssv2", region)
		if err != nil {
			return diag.Errorf("error creating BSS client: %s", err)
		}

		orderId := utils.PathSearch("order_id", respBody, "").(string)
		devServerId, err := waitForPrePaidDevServerComplete(ctx, client, bssClient, d.Timeout(schema.TimeoutCreate), orderId)
		if devServerId != "" {
			// When a DevServer fails to be created, the failed server will exist in the cloud. so set the resource ID in advance to
			// prevent dirty data.
			d.SetId(devServerId)
		}

		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		devServerId := utils.PathSearch("id", respBody, "").(string)
		if devServerId == "" {
			return diag.Errorf("unable to find DevServer ID from API response")
		}

		d.SetId(devServerId)
		err := waitForDevServerCreated(ctx, client, d.Timeout(schema.TimeoutCreate), devServerId)
		if err != nil {
			return diag.Errorf("error waiting for the DevServer (%s) creation to complete: %s", devServerId, err)
		}
	}
	return resourceDevServerRead(ctx, d, meta)
}

func waitForDevServerCreated(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration, devServerId string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshDevServerStatusFunc(client, devServerId, []string{"RUNNING"}),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshDevServerStatusFunc(client *golangsdk.ServiceClient, serverId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetDevServerById(client, serverId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "deleted", "DELETED", nil
			}
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", respBody, "").(string)
		if utils.StrSliceContains([]string{"CREATE_FAILED", "ERROR", "DELETE_FAILED"}, status) {
			return respBody, "ERROR", fmt.Errorf("unexpected status (%s)", status)
		}

		if utils.StrSliceContains(targets, status) {
			return respBody, "COMPLETED", nil
		}
		return "continue", "PENDING", nil
	}
}

func waitForPrePaidDevServerComplete(ctx context.Context, client, bssClient *golangsdk.ServiceClient, timeOut time.Duration,
	orderId string) (string, error) {
	if orderId == "" {
		return "", fmt.Errorf("unable to find order ID from API response")
	}

	if err := common.WaitOrderComplete(ctx, bssClient, orderId, timeOut); err != nil {
		return "", err
	}

	// When DevServer creation fails, the `common.WaitOrderResourceComplete` method cannot obtain the resource ID.
	// so the Devserver list query interface is called here.
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			devServers, err := getDevServersByOrderId(client, orderId)
			if err != nil {
				return "", "ERROR", err
			}

			devServer := utils.PathSearch(fmt.Sprintf("[?order_id=='%s']|[0]", orderId), devServers, nil)
			devServerId := utils.PathSearch("id", devServer, "").(string)
			status := utils.PathSearch("status", devServer, "").(string)
			if utils.StrSliceContains([]string{"CREATE_FAILED", "ERROR"}, status) {
				return devServerId, "ERROR", fmt.Errorf("unexpected status (%s)", status)
			}

			if status == "RUNNING" {
				return devServerId, "COMPLETED", nil
			}
			return devServerId, "PENDING", nil
		},
		Timeout:      timeOut,
		Delay:        10 * time.Second,
		PollInterval: 30 * time.Second,
	}

	devServerId, err := stateConf.WaitForStateContext(ctx)
	resourceId := devServerId.(string)
	if err != nil {
		return resourceId, fmt.Errorf("error waiting for DevService order (%s) to complete: %s", orderId, err)
	}

	if resourceId == "" {
		return "", fmt.Errorf("unable to find DevServer ID from API response")
	}
	return resourceId, nil
}

func GetDevServerById(client *golangsdk.ServiceClient, serverId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/dev-servers/{id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", serverId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", devServerNotExistCode)
	}

	return utils.FlattenResponse(requestResp)
}

func getDevServersByOrderId(client *golangsdk.ServiceClient, orderId string) (interface{}, error) {
	var (
		httpUrl  = "v1/{project_id}/dev-servers?limit=100"
		offset   = 0
		result   = make([]interface{}, 0)
		listOpts = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &listOpts)
		if err != nil {
			return nil, fmt.Errorf("error getting list of the DevServers by order ID (%s): %s", orderId, err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		devServers := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(devServers) < 1 {
			break
		}
		result = append(result, devServers...)
		if len(result) == int(utils.PathSearch("total", respBody, float64(0)).(float64)) {
			break
		}
		offset += len(devServers)
	}

	return result, nil
}

func resourceDevServerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		devServerId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	respBody, err := GetDevServerById(client, devServerId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ModelArts DevServer")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("flavor", utils.PathSearch("flavor", respBody, nil)),
		d.Set("architecture", utils.PathSearch("image.arch", respBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", respBody, nil)),
		d.Set("image_id", utils.PathSearch("image.image_id", respBody, nil)),
		d.Set("type", utils.PathSearch("cloud_server.type", respBody, nil)),
		d.Set("key_pair_name", utils.PathSearch("key_pair_name", respBody, nil)),
		d.Set("charging_mode", utils.PathSearch("charging_mode", respBody, nil)),
		d.Set("created_at",
			utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_at", respBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDevServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("auto_renew") {
		var (
			cfg         = meta.(*config.Config)
			region      = cfg.GetRegion(d)
			devServerId = d.Id()
		)
		bssClient, err := cfg.NewServiceClient("bssv2", region)
		if err != nil {
			return diag.Errorf("error creating BSS client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), devServerId); err != nil {
			return diag.Errorf("error updating the auto_renew of the DevServer (%s): %s", devServerId, err)
		}
	}
	return resourceDevServerRead(ctx, d, meta)
}

func resourceDevServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		devServerId = d.Id()
	)

	bssV2Client, err := cfg.NewServiceClient("bssv2", region)
	if err != nil {
		return diag.Errorf("error creating BSS client: %s", err)
	}
	resources, err := getPrePaidResourcesById(bssV2Client, devServerId)
	if err != nil {
		return diag.FromErr(err)
	}

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	if len(resources) != 0 {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{devServerId}); err != nil {
			// CBC.30000067: The resource has been unsubscribed
			return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "CBC.30000067"),
				fmt.Sprintf("error unsubscribing DevServer (%s)", devServerId))
		}
	} else {
		if err := deleteDevServerById(client, devServerId); err != nil {
			return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting DevServer (%s)", devServerId))
		}
	}

	err = waitForDevServerDeleted(ctx, client, devServerId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		diag.Errorf("error waiting for the DevServer (%s) deletion to complete: %s", devServerId, err)
	}

	return nil
}

func getPrePaidResourcesById(client *golangsdk.ServiceClient, devServerId string) ([]interface{}, error) {
	httpUrl := "v2/orders/suscriptions/resources/query"
	getPath := client.Endpoint + httpUrl
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"resource_ids": []string{devServerId},
		},
	}

	requestResp, err := client.Request("POST", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func deleteDevServerById(client *golangsdk.ServiceClient, devServerId string) error {
	httpUrl := "v1/{project_id}/dev-servers/{id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", devServerId)
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.ConvertExpected400ErrInto404Err(err, "error_code", devServerNotExistCode)
	}

	return nil
}

func waitForDevServerDeleted(ctx context.Context, client *golangsdk.ServiceClient, devServerId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING", "DELETING"},
		Target:       []string{"DELETED"},
		Refresh:      refreshDevServerStatusFunc(client, devServerId, nil),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
