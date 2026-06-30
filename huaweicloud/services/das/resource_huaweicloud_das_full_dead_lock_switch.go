package das

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var fullDeadLockSwitchNonUpdatableParams = []string{
	"instance_id",
	"switch_on",
	"retention_hours",
}

// @API DAS POST /v3/{project_id}/instances/{instance_id}/set-full-dead-lock-switch
// @API DAS GET /v3/{project_id}/instances/{instance_id}/get-full-dead-lock-switch
func ResourceFullDeadLockSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFullDeadLockSwitchCreate,
		ReadContext:   resourceFullDeadLockSwitchRead,
		UpdateContext: resourceFullDeadLockSwitchUpdate,
		DeleteContext: resourceFullDeadLockSwitchDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(fullDeadLockSwitchNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the full dead lock switch is located.",
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the database instance.",
			},
			"switch_on": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable the full dead lock switch.",
			},

			// Optional parameters.
			"retention_hours": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The retention hours of the full dead lock data.",
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					"Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.",
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildFullDeadLockSwitchBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"switch_on": d.Get("switch_on").(bool),

		// The API only supports 'mysql' and is case-sensitive.
		"engine_type": "mysql",
	}

	if v, ok := d.GetOk("retention_hours"); ok {
		bodyParams["retention_hours"] = v.(int)
	}

	return bodyParams
}

func fullDeadLockSwitchRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		httpUrl := "v3/{project_id}/instances/{instance_id}/get-full-dead-lock-switch"
		getPath := client.Endpoint + httpUrl
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type":     "application/json",
				"X-Source-Service": "das",
			},
		}

		requestResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, "ERROR", err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, "ERROR", err
		}

		switchOn := utils.PathSearch("switch_on", respBody, nil)
		if switchOn == nil {
			return nil, "ERROR", errors.New("'switch_on' field is not found in the response")
		}

		currentStatus := strconv.FormatBool(switchOn.(bool))
		targetStatus := strconv.FormatBool(d.Get("switch_on").(bool))

		if currentStatus == targetStatus {
			return respBody, "COMPLETED", nil
		}
		return respBody, "PENDING", nil
	}
}

func waitForFullDeadLockSwitchComplete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      fullDeadLockSwitchRefreshFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the DAS full dead lock switch to become %v: %s", d.Get("switch_on").(bool), err)
	}
	return nil
}

func resourceFullDeadLockSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/set-full-dead-lock-switch"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
		JSONBody: buildFullDeadLockSwitchBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error switching DAS full dead lock: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	if err = waitForFullDeadLockSwitchComplete(ctx, client, d); err != nil {
		return diag.FromErr(err)
	}

	return resourceFullDeadLockSwitchRead(ctx, d, meta)
}

func resourceFullDeadLockSwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFullDeadLockSwitchUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFullDeadLockSwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for switching the full dead lock. Deleting
this resource will not clear the corresponding request record, but will only remove the resource information
from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
