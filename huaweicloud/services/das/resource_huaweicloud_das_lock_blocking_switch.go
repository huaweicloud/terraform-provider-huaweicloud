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

var lockBlockingSwitchNonUpdatableParams = []string{
	"instance_id",
	"switch_on",
	"retention_hours",
}

// @API DAS POST /v3/{project_id}/lock-blocking/switch
// @API DAS GET /v3/{project_id}/lock-blocking/switch
func ResourceLockBlockingSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLockBlockingSwitchCreate,
		ReadContext:   resourceLockBlockingSwitchRead,
		UpdateContext: resourceLockBlockingSwitchUpdate,
		DeleteContext: resourceLockBlockingSwitchDelete,

		CustomizeDiff: config.FlexibleForceNew(lockBlockingSwitchNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the lock blocking switch is located.",
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
				Description: "Whether to enable the lock blocking switch.",
			},

			// Optional parameters.
			"retention_hours": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The retention hours of the lock blocking data.",
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

func buildLockBlockingSwitchBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_id": d.Get("instance_id").(string),
		"switch_on":   d.Get("switch_on").(bool),

		// The API only supports 'SQLServer' and is case-insensitive.
		"engine_type": "SQLServer",
	}

	if v, ok := d.GetOk("retention_hours"); ok {
		bodyParams["retention_hours"] = v.(int)
	}

	return bodyParams
}

func lockBlockingSwitchRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		httpUrl := "v3/{project_id}/lock-blocking/switch"
		getPath := client.Endpoint + httpUrl
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = fmt.Sprintf("%s?instance_id=%s&engine_type=%s", getPath, d.Get("instance_id").(string), "SQLServer")

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

func waitForLockBlockingSwitchComplete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      lockBlockingSwitchRefreshFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the DAS lock blocking switch to become %v: %s", d.Get("switch_on").(bool), err)
	}
	return nil
}

func resourceLockBlockingSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/lock-blocking/switch"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
		JSONBody: buildLockBlockingSwitchBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error switching DAS lock blocking: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	if err = waitForLockBlockingSwitchComplete(ctx, client, d); err != nil {
		return diag.Errorf("error waiting for the DAS lock blocking switch to complete: %s", err)
	}

	return resourceLockBlockingSwitchRead(ctx, d, meta)
}

func resourceLockBlockingSwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceLockBlockingSwitchUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceLockBlockingSwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for switching the lock blocking. Deleting
this resource will not clear the corresponding request record, but will only remove the resource information
from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
