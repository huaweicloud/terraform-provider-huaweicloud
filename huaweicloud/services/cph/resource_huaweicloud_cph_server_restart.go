package cph

import (
	"context"
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

var serverRestartNonUpdatableParams = []string{"server_id"}

// @API CPH POST /v1/{project_id}/cloud-phone/servers/batch-restart
// @API CPH GET /v1/{project_id}/cloud-phone/servers/{server_id}
func ResourceServerRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerRestartCreate,
		UpdateContext: resourceServerRestartUpdate,
		ReadContext:   resourceServerRestartRead,
		DeleteContext: resourceServerRestartDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(serverRestartNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of CPH server.`,
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

func resourceServerRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("cph", region)
	if err != nil {
		return diag.Errorf("error creating CPH client: %s", err)
	}

	// createServerRestart: create CPH server restart
	createServerRestartHttpUrl := "v1/{project_id}/cloud-phone/servers/batch-restart"
	createServerRestartPath := client.Endpoint + createServerRestartHttpUrl
	createServerRestartPath = strings.ReplaceAll(createServerRestartPath, "{project_id}", client.ProjectID)

	createServerRestartOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"server_ids": []string{d.Get("server_id").(string)},
		},
	}

	createServerRestartResp, err := client.Request("POST", createServerRestartPath, &createServerRestartOpt)
	if err != nil {
		return diag.Errorf("error restarting CPH server: %s", err)
	}

	createServerRestartRespBody, err := utils.FlattenResponse(createServerRestartResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("jobs|[0].server_id", createServerRestartRespBody, "").(string)
	if id == "" {
		return diag.Errorf("Unable to find the server ID from the API response")
	}
	d.SetId(id)

	errorCode := utils.PathSearch("jobs|[0].error_code", createServerRestartRespBody, "").(string)
	if errorCode != "" {
		errorMsg := utils.PathSearch("jobs|[0].error_msg", createServerRestartRespBody, "").(string)
		return diag.Errorf("failed to stop CPH server (server_id: %s), error_code: %s, error_msg: %s", id, errorCode, errorMsg)
	}

	err = checkServerStatus(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceServerRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceServerRestartUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceServerRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting CPH server restart resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func checkServerStatus(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      serverStateRefreshFunc(client, id),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CPH server restart to be completed: %s", err)
	}
	return nil
}

func serverStateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getServerRespBody, err := getServerDetail(client, id)
		if err != nil {
			return nil, "ERROR", err
		}

		// Status is 5, indicates the server is running normally.
		serverStatus := utils.PathSearch("status", getServerRespBody, float64(0)).(float64)
		if int(serverStatus) == 5 {
			return getServerRespBody, "COMPLETED", nil
		}
		return getServerRespBody, "PENDING", nil
	}
}

func getServerDetail(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	getServerHttpUrl := "v1/{project_id}/cloud-phone/servers/{server_id}"
	getServerPath := client.Endpoint + getServerHttpUrl
	getServerPath = strings.ReplaceAll(getServerPath, "{project_id}", client.ProjectID)
	getServerPath = strings.ReplaceAll(getServerPath, "{server_id}", id)

	getServerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getServerResp, err := client.Request("GET", getServerPath, &getServerOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getServerResp)
}
