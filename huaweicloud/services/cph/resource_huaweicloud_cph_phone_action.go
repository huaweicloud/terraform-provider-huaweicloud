package cph

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var PhoneActionNonUpdatableParams = []string{"action", "phones", "image_id"}
var (
	RESET   = "reset"
	RESTART = "restart"
	START   = "start"
	STOP    = "stop"
)

// @API CPH POST /v1/{project_id}/cloud-phone/phones/batch-restart
// @API CPH POST /v1/{project_id}/cloud-phone/phones/batch-reset
// @API CPH POST /v1/{project_id}/cloud-phone/phones/batch-stop
// @API CPH GET /v1/{project_id}/cloud-phone/jobs
func ResourcePhoneAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePhoneActionCreate,
		UpdateContext: resourcePhoneActionUpdate,
		ReadContext:   resourcePhoneActionRead,
		DeleteContext: resourcePhoneActionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(PhoneActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"reset", "restart", "start", "stop"}, false),
				Description:  `Specifies the CPH phone action.`,
			},
			"phones": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Specifies the CPH phones.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"phone_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the phone ID.`,
						},
						"property": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the phone property.`,
						},
					},
				},
			},
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the image ID of the CPH phone.`,
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

func resourcePhoneActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	action := d.Get("action").(string)

	client, err := cfg.NewServiceClient("cph", region)
	if err != nil {
		return diag.Errorf("error creating CPH client: %s", err)
	}

	var resp interface{}
	var respErr interface{}

	// Call CPH phone operation API
	switch action {
	case RESET:
		resp, respErr = phoneReset(client, d)
	case RESTART:
		resp, respErr = phoneRestart(client, d)
	case STOP:
		resp, respErr = phoneStop(client, d)
	}
	if respErr != nil {
		return diag.FromErr(err)
	}
	id := utils.PathSearch("request_id", resp, "").(string)
	if id == "" {
		return diag.Errorf("Unable to find the request ID from the API response")
	}
	d.SetId(id)

	jobs := utils.PathSearch("jobs", resp, make([]interface{}, 0)).([]interface{})
	for _, v := range jobs {
		errorCode := utils.PathSearch("error_code", v, "").(string)
		if errorCode != "" {
			serverId := utils.PathSearch("phone_id", v, "").(string)
			errorMsg := utils.PathSearch("error_msg", v, "").(string)
			log.Printf("[WARN] Failed to restart CPH phone (phone_id: %s), error_code: %s, error_msg: %s", serverId, errorCode, errorMsg)
		}
	}

	err = checkJobStatus(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("error waiting for CPH phone restart to completed: %s", err),
			},
		}
	}

	return nil
}

func resourcePhoneActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePhoneActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePhoneActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting CPH action resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func phoneRestart(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	// createPhoneRestart: create CPH phone restart
	createPhoneRestartHttpUrl := "v1/{project_id}/cloud-phone/phones/batch-restart"
	createPhoneRestartPath := client.Endpoint + createPhoneRestartHttpUrl
	createPhoneRestartPath = strings.ReplaceAll(createPhoneRestartPath, "{project_id}", client.ProjectID)

	createPhoneRestartOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createPhoneRestartOpt.JSONBody = utils.RemoveNil(map[string]interface{}{
		"image_id": utils.ValueIgnoreEmpty(d.Get("image_id")),
		"phones":   d.Get("phones"),
	})
	createPhoneRestartResp, err := client.Request("POST", createPhoneRestartPath, &createPhoneRestartOpt)
	if err != nil {
		return nil, fmt.Errorf("error creating CPH phone restart: %s", err)
	}

	return utils.FlattenResponse(createPhoneRestartResp)
}

func phoneReset(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	// createPhoneRestart: create CPH phone reset
	createPhoneResetHttpUrl := "v1/{project_id}/cloud-phone/phones/batch-reset"
	createPhoneResetPath := client.Endpoint + createPhoneResetHttpUrl
	createPhoneResetPath = strings.ReplaceAll(createPhoneResetPath, "{project_id}", client.ProjectID)

	createPhoneResetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createPhoneResetOpt.JSONBody = utils.RemoveNil(map[string]interface{}{
		"image_id": utils.ValueIgnoreEmpty(d.Get("image_id")),
		"phones":   d.Get("phones"),
	})
	createPhoneResetResp, err := client.Request("POST", createPhoneResetPath, &createPhoneResetOpt)
	if err != nil {
		return nil, fmt.Errorf("error creating CPH phone reset: %s", err)
	}

	return utils.FlattenResponse(createPhoneResetResp)
}

func phoneStop(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	// createPhoneStop: create CPH phone stop
	createPhoneStopHttpUrl := "v1/{project_id}/cloud-phone/phones/batch-stop"
	createPhoneStopPath := client.Endpoint + createPhoneStopHttpUrl
	createPhoneStopPath = strings.ReplaceAll(createPhoneStopPath, "{project_id}", client.ProjectID)

	createPhoneStopOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createPhoneStopOpt.JSONBody = map[string]interface{}{
		"phone_ids": utils.PathSearch("[*].phone_id", d.Get("phones"), make([]interface{}, 0)),
	}
	createPhoneStopResp, err := client.Request("POST", createPhoneStopPath, &createPhoneStopOpt)
	if err != nil {
		return nil, fmt.Errorf("error creating CPH phone stop: %s", err)
	}

	return utils.FlattenResponse(createPhoneStopResp)
}
