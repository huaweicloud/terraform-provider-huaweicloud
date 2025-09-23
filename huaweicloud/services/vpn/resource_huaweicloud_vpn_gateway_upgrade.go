package vpn

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var gatewayUpgradeNonUpdatableParams = []string{"vgw_id", "action"}

// @API VPN POST /v5/{project_id}/vpn-gateways/{vgw_id}/upgrade
// @API VPN GET /v5/{project_id}/vpn-gateways/jobs
func ResourceGatewayUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGatewayUpgradeCreate,
		ReadContext:   resourceGatewayUpgradeRead,
		UpdateContext: resourceGatewayUpgradeUpdate,
		DeleteContext: resourceGatewayUpgradeDelete,

		CustomizeDiff: config.FlexibleForceNew(gatewayUpgradeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vgw_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID of a VPN gateway.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies an upgrade operation.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceGatewayUpgradeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	gatewayId := d.Get("vgw_id").(string)

	createHttpUrl := "v5/{project_id}/vpn-gateways/{vgw_id}/upgrade"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{vgw_id}", gatewayId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"action": d.Get("action"),
		},
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating VPN gateway upgrade: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("job_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating VPN gateway upgrade: job ID is not found in API response")
	}
	d.SetId(id)

	if err = waitForGatewayJobComplete(ctx, client, d); err != nil {
		return diag.Errorf("error waiting for gateway job (%s) to be completed: %s", id, err)
	}

	return nil
}

func waitForGatewayJobComplete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var pending []string
	var target []string

	statusType := d.Get("action").(string)
	switch statusType {
	case "start":
		pending = []string{"upgrading"}
		target = []string{"pending_upgrade_confirm"}
	case "finish":
		pending = []string{"pending_upgrade_confirm"}
		target = []string{"success"}
	case "rollback":
		pending = []string{"rolling_back"}
		target = []string{"rollback_success"}
	default:
		return errors.New("unsupport action")
	}
	stateConf := &resource.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			job, err := getGatewayJob(client, d)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", job, "").(string)
			return job, status, nil
		},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Second,
		PollInterval: 3 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func getGatewayJob(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v5/{project_id}/vpn-gateways/jobs?resource_id={resource_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", d.Get("project_id").(string))
	getPath = strings.ReplaceAll(getPath, "{resource_id}", d.Get("vgw_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	searchPath := fmt.Sprintf("jobs[?id=='%s']", d.Id())
	job := utils.PathSearch(searchPath, getRespBody, nil)
	if job == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return job, nil
}

func resourceGatewayUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGatewayUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGatewayUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting VPN gateway upgrade resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
