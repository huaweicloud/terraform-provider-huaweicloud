package dcs

import (
	"context"
	"encoding/json"
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

var migrationTaskExchangeIpNonUpdatableParams = []string{"task_id", "exchanged_ip", "is_exchange_domain"}

// @API DCS POST /v2/{project_id}/migration-task/{task_id}/exchange-ip
// @API DCS GET /v2/{project_id}/migration-task/{task_id}
func ResourceDcsMigrationTaskExchangeIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsMigrationTaskExchangeIpCreate,
		ReadContext:   resourceDcsMigrationTaskExchangeIpRead,
		UpdateContext: resourceDcsMigrationTaskExchangeIpUpdate,
		DeleteContext: resourceDcsMigrationTaskExchangeIppDelete,

		CustomizeDiff: config.FlexibleForceNew(migrationTaskExchangeIpNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"exchanged_ip": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_exchange_domain": {
				Type:     schema.TypeBool,
				Optional: true,
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

func resourceDcsMigrationTaskExchangeIpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/migration-task/{task_id}/exchange-ip"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	taskId := d.Get("task_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{task_id}", taskId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 204},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateMigrationTaskExchangeIpBodyParams(d))

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error exchanging migration task(%s) IP: %s", taskId, err)
	}

	d.SetId(taskId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"EXCHANGE_SUCCESS"},
		Refresh:      migrationTaskSwitchIpRefreshFunc(client, taskId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for migration task(%s) exchanging IP to be completed: %s ", taskId, err)
	}
	return nil
}

func buildCreateMigrationTaskExchangeIpBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"exchanged_ip": utils.ValueIgnoreEmpty(d.Get("exchanged_ip").(*schema.Set).List()),
	}
	if d.Get("is_exchange_domain").(bool) {
		bodyParams["is_exchange_domain"] = true
	}
	return bodyParams
}

func migrationTaskSwitchIpRefreshFunc(client *golangsdk.ServiceClient, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRespBody, err := getMigrationTask(client, taskId)
		if err != nil {
			if errCode, ok := err.(golangsdk.ErrDefault400); ok {
				var response interface{}
				if jsonErr := json.Unmarshal(errCode.Body, &response); jsonErr == nil {
					errorCode := utils.PathSearch("error_code", response, "").(string)
					if errorCode == "DCS.4133" {
						return "", "DELETED", nil
					}
				}
			}
			return nil, "ERROR", err
		}

		status := utils.PathSearch("task_status", getRespBody, "").(string)
		successStatus := []string{"EXCHANGE_SUCCESS", "ROLLBACK_SUCCESS"}
		if utils.StrSliceContains(successStatus, status) {
			return getRespBody, status, nil
		}

		return getRespBody, "PENDING", nil
	}
}

func resourceDcsMigrationTaskExchangeIpRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsMigrationTaskExchangeIpUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsMigrationTaskExchangeIppDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS exchanging migration task IP resource is not supported. The resource is only removed from" +
		" the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
