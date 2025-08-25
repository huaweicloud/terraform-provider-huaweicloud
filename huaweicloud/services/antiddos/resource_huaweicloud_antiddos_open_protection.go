package antiddos

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableAntiDdosOpenProtectionParams = []string{
	"floating_ip_id",
	"app_type_id",
	"cleaning_access_pos_id",
	"enable_l7",
	"http_request_pos_id",
	"traffic_pos_id",
	"antiddos_config_id",
}

// @API ANTI-DDOS POST /v1/{project_id}/antiddos/{floating_ip_id}
func ResourceAntiDdosOpenProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAntiDdosOpenProtectionCreate,
		ReadContext:   resourceAntiDdosOpenProtectionRead,
		UpdateContext: resourceAntiDdosOpenProtectionUpdate,
		DeleteContext: resourceAntiDdosOpenProtectionDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableAntiDdosOpenProtectionParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"floating_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the EIP ID.`,
			},
			"app_type_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the application type ID.`,
			},
			"cleaning_access_pos_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the cleaning access position ID.`,
			},
			"enable_l7": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether to enable L7 protection.`,
			},
			"http_request_pos_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the HTTP request position ID.`,
			},
			"traffic_pos_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the traffic position ID.`,
			},
			"antiddos_config_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the anti-DDoS configuration ID.`,
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

func buildAntiDdosOpenProtectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"floating_ip_id":         d.Get("floating_ip_id"),
		"app_type_id":            d.Get("app_type_id"),
		"cleaning_access_pos_id": d.Get("cleaning_access_pos_id"),
		"enable_L7":              d.Get("enable_l7"),
		"http_request_pos_id":    d.Get("http_request_pos_id"),
		"traffic_pos_id":         d.Get("traffic_pos_id"),
		"antiddos_config_id":     utils.ValueIgnoreEmpty(d.Get("antiddos_config_id")),
	}
}

func waitingAntiDdosOpenProtectionTaskSuccess(ctx context.Context, client *golangsdk.ServiceClient, taskID string,
	timeout time.Duration) error {
	var (
		errorStatuses = []string{"failed"}
		httpUrl       = "v2/{project_id}/query-task-status"
	)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (result interface{}, state string, err error) {
			requestPath := client.Endpoint + httpUrl
			requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
			requestPath = fmt.Sprintf("%s?task_id=%s", requestPath, taskID)
			requestOpt := golangsdk.RequestOpts{
				MoreHeaders: map[string]string{
					"Content-Type": "application/json;charset=utf8",
				},
				KeepResponseBody: true,
			}

			resp, err := client.Request("GET", requestPath, &requestOpt)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error retrieving Anti-DDoS task: %s", err)
			}

			respBody, err := utils.FlattenResponse(resp)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error retrieving Anti-DDoS task: Failed to flatten response (%s)", err)
			}

			taskStatus := utils.PathSearch("task_status", respBody, "").(string)
			if taskStatus == "" {
				return nil, "ERROR", errors.New("error retrieving Anti-DDoS task: Task status is not found in API response")
			}

			if utils.StrSliceContains(errorStatuses, taskStatus) {
				return respBody, "ERROR", fmt.Errorf("unexpected status: '%s'", taskStatus)
			}

			if taskStatus == "success" {
				return respBody, "COMPLETED", nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceAntiDdosOpenProtectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/antiddos/{floating_ip_id}"
		product = "anti-ddos"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Anti-DDoS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{floating_ip_id}", d.Get("floating_ip_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAntiDdosOpenProtectionBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error opening Anti-DDoS protection: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := utils.PathSearch("task_id", respBody, "").(string)
	if taskID == "" {
		return diag.Errorf("error opening Anti-DDoS protection: task_id is not found in API response")
	}

	if err := waitingAntiDdosOpenProtectionTaskSuccess(ctx, client, taskID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for Anti-DDoS task to become success: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	return resourceAntiDdosOpenProtectionRead(ctx, d, meta)
}

func resourceAntiDdosOpenProtectionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceAntiDdosOpenProtectionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceAntiDdosOpenProtectionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to open Anti-DDoS protection. Deleting this 
resource will not change the current Anti-DDoS opening status, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
