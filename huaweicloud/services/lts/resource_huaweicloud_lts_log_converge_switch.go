package lts

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS PUT /v1/{project_id}/lts/log-converge-config/switch
// @API LTS GET /v1/{project_id}/lts/log-converge-config/switch
func ResourceLogConvergeSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogConvergeSwitchCreate,
		ReadContext:   resourceLogConvergeSwitchRead,
		DeleteContext: resourceLogConvergeSwitchDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func modifyLogConvergeConfigsMessageSwitch(client *golangsdk.ServiceClient, switchTarget bool) error {
	httpUrl := "v1/{project_id}/lts/log-converge-config/switch?log_converge_switch={log_converge_switch}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{log_converge_switch}", strconv.FormatBool(switchTarget))

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	requestResp, err := client.Request("PUT", getPath, &opts)
	if err != nil {
		return fmt.Errorf("failed to enable log receiving status (target: %v): %s", switchTarget, err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return fmt.Errorf("error retrieving API response body: %v", err)
	}
	requestResult := utils.PathSearch("result", respBody, nil)
	if requestResult != "success" {
		return fmt.Errorf("failed to enable log receiving status (target: %v), but the result of API request is: %v",
			switchTarget, requestResult)
	}
	return nil
}

func resourceLogConvergeSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	if err = modifyLogConvergeConfigsMessageSwitch(client, true); err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	return resourceLogConvergeSwitchRead(ctx, d, meta)
}

func GetLogConvergeSwitchEnabled(client *golangsdk.ServiceClient) (bool, error) {
	httpUrl := "v1/{project_id}/lts/log-converge-config/switch"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return false, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return false, err
	}
	switchEnabled := utils.PathSearch("log_converge_switch", respBody, false).(bool)
	if !switchEnabled {
		return switchEnabled, golangsdk.ErrDefault404{}
	}
	return switchEnabled, nil
}

func resourceLogConvergeSwitchRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	_, err = GetLogConvergeSwitchEnabled(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "LTS log converge switch")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLogConvergeSwitchDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	if err = modifyLogConvergeConfigsMessageSwitch(client, false); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
