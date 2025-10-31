package rgc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var controlNonUpdatableParams = []string{"identifier", "target_identifier", "parameters"}

// @API RGC POST /v1/governance/controls/enable
// @API RGC POST /v1/governance/controls/disable
// @API RGC GET /v1/governance/operation-control-status/{operation_control_status_id}
func ResourceControl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceControlCreate,
		UpdateContext: resourceControlUpdate,
		ReadContext:   resourceControlRead,
		DeleteContext: resourceControlDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(controlNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identifier": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_identifier": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceControlCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		enableControlHttpUrl = "v1/governance/controls/enable"
		enableControlProduct = "rgc"
	)
	enableControlClient, err := cfg.NewServiceClient(enableControlProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	enableControlPath := enableControlClient.Endpoint + enableControlHttpUrl

	body, err := buildControlBodyParams(d)
	if err != nil {
		diag.FromErr(err)
	}
	enableControlOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         body,
	}
	enableControlResp, err := enableControlClient.Request("POST", enableControlPath, &enableControlOpt)
	if err != nil {
		return diag.Errorf("error creating control: %s", err)
	}

	enableControlRespBody, err := utils.FlattenResponse(enableControlResp)
	if err != nil {
		return diag.FromErr(err)
	}

	operationId := utils.PathSearch("control_operate_request_id", enableControlRespBody, "").(string)

	if operationId == "" {
		return diag.Errorf("get empty operationId")
	}

	stateConf := resource.StateChangeConf{
		Pending:      []string{"IN_PROGRESS"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      controlStateRefreshFunc(enableControlClient, operationId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	controlId := d.Get("identifier").(string)
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for RGC control (%s) to create: %s", controlId, err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return resourceControlRead(ctx, d, meta)
}

func resourceControlRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getControlHttpUrl = "v1/governance/managed-organizational-units/{managed_organizational_unit_id}/controls/{control_id}"
		getControlProduct = "rgc"
	)
	getControlClient, err := cfg.NewServiceClient(getControlProduct, region)
	if err != nil {
		return diag.Errorf("error creating rgc client: %s", err)
	}

	controlId := d.Get("identifier").(string)
	ouId := d.Get("target_identifier").(string)

	getControlPath := getControlClient.Endpoint + getControlHttpUrl
	getControlPath = strings.ReplaceAll(getControlPath, "{managed_organizational_unit_id}", ouId)
	getControlPath = strings.ReplaceAll(getControlPath, "{control_id}", controlId)

	getControlOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getControlResp, err := getControlClient.Request("GET", getControlPath, &getControlOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving control")
	}

	getControlRespBody, err := utils.FlattenResponse(getControlResp)
	if err != nil {
		return diag.FromErr(err)
	}

	state := utils.PathSearch("state", getControlRespBody, "").(string)
	if state == "DISABLED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("state", state),
		d.Set("version", utils.PathSearch("version", getControlRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceControlUpdate(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceControlDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		disableControlHttpUrl = "v1/governance/controls/disable"
		disableControlProduct = "rgc"
	)
	disableControlClient, err := cfg.NewServiceClient(disableControlProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	disableControlPath := disableControlClient.Endpoint + disableControlHttpUrl

	body, err := buildControlBodyParams(d)
	if err != nil {
		diag.FromErr(err)
	}
	disableControlOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         body,
	}

	disableControlResp, err := disableControlClient.Request("POST", disableControlPath, &disableControlOpt)
	if err != nil {
		return diag.Errorf("error disable control: %s", err)
	}

	disableControlRespBody, err := utils.FlattenResponse(disableControlResp)
	if err != nil {
		return diag.FromErr(err)
	}

	operationId := utils.PathSearch("control_operate_request_id", disableControlRespBody, "").(string)

	stateConf := resource.StateChangeConf{
		Pending:      []string{"IN_PROGRESS"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      controlStateRefreshFunc(disableControlClient, operationId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	controlId := d.Get("identifier").(string)
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for RGC control (%s) to create: %s", controlId, err)
	}

	return nil
}

func buildControlBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	bodyParams := map[string]interface{}{
		"identifier":        d.Get("identifier").(string),
		"target_identifier": d.Get("target_identifier").(string),
	}

	if v, ok := d.GetOk("parameters"); ok {
		var params []map[string]interface{}
		err := json.Unmarshal([]byte(v.(string)), &params)
		if err != nil {
			return nil, fmt.Errorf("error unmarshal parameters: %s", err)
		}

		bodyParams["parameters"] = params
	}

	return bodyParams, nil
}

func controlStateRefreshFunc(client *golangsdk.ServiceClient, operationId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getControlStatusHttpUrl := "v1/governance/operation-control-status/{operation_control_status_id}"
		getControlStatusPath := client.Endpoint + getControlStatusHttpUrl
		getControlStatusPath = strings.ReplaceAll(getControlStatusPath, "{operation_control_status_id}", operationId)

		getControlStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getControlStatusResp, err := client.Request("GET", getControlStatusPath, &getControlStatusOpt)
		if err != nil {
			return nil, "", err
		}

		getControlStatusRespBody, err := utils.FlattenResponse(getControlStatusResp)
		if err != nil {
			return nil, "", err
		}
		controlOperation := utils.PathSearch("control_operation", getControlStatusRespBody, nil).(map[string]interface{})
		status := controlOperation["status"].(string)
		if status == "FAILED" {
			message := utils.PathSearch("message", getControlStatusRespBody, nil)
			return nil, "FAILED", fmt.Errorf("status: %s; message: %s", status, message)
		}

		return getControlStatusRespBody, status, nil
	}
}
