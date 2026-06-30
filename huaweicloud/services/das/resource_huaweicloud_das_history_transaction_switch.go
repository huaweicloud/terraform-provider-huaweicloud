package das

import (
	"context"
	"errors"
	"fmt"
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

var historyTransactionSwitchNonUpdatableParams = []string{
	"instance_id",
	"status",
	"datastore_type",
}

// @API DAS POST /v3/{project_id}/instances/{instance_id}/transaction/switch
// @API DAS GET /v3/{project_id}/instances/{instance_id}/transaction/switch
func ResourceHistoryTransactionSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHistoryTransactionSwitchCreate,
		ReadContext:   resourceHistoryTransactionSwitchRead,
		UpdateContext: resourceHistoryTransactionSwitchUpdate,
		DeleteContext: resourceHistoryTransactionSwitchDelete,

		CustomizeDiff: config.FlexibleForceNew(historyTransactionSwitchNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the history transaction switch is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the database instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The switch status of the history transaction.`,
			},
			"datastore_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The database type.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildHistoryTransactionSwitchBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"switch_status":  d.Get("status").(string),
		"datastore_type": d.Get("datastore_type").(string),
	}
}

func waitForHistoryTransactionSwitchComplete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	targetStatus := d.Get("status").(string)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Switching"},
		Target:       []string{targetStatus},
		Refresh:      historyTransactionSwitchRefreshFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the DAS history transaction switch to become %s: %s", targetStatus, err)
	}
	return nil
}

func resourceHistoryTransactionSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/transaction/switch"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildHistoryTransactionSwitchBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error switching DAS history transaction: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	if err = waitForHistoryTransactionSwitchComplete(ctx, client, d); err != nil {
		return diag.Errorf("error waiting for the DAS history transaction switch to complete: %s", err)
	}

	return resourceHistoryTransactionSwitchRead(ctx, d, meta)
}

func historyTransactionSwitchRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		httpUrl := "v3/{project_id}/instances/{instance_id}/transaction/switch"
		getPath := client.Endpoint + httpUrl
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
		getPath = fmt.Sprintf("%s?datastore_type=%s", getPath, d.Get("datastore_type").(string))

		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
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

		status := utils.PathSearch("switch_status", respBody, "").(string)
		if status == "" {
			return nil, "ERROR", errors.New("switch_status is empty in the response")
		}

		return respBody, status, nil
	}
}

func resourceHistoryTransactionSwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceHistoryTransactionSwitchUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceHistoryTransactionSwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for switching the history transaction. Deleting
this resource will not clear the corresponding request record, but will only remove the resource information
from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
