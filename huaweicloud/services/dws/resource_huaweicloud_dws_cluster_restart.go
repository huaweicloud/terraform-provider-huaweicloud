package dws

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDM POST /v1.0/{project_id}/clusters/{cluster_id}/restart
// @API DWS GET /v1.0/{project_id}/clusters/{cluster_id}
func ResourceClusterRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterRestartCreate,
		ReadContext:   resourceClusterRestartRead,
		DeleteContext: resourceClusterRestartDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the DWS cluster to be restarted.",
			},
		},
	}
}

func resourceClusterRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterId := d.Get("cluster_id").(string)
	// For the same DWS cluster, it is not supported to run multiple tasks at the same time.
	config.MutexKV.Lock(clusterId)
	defer config.MutexKV.Unlock(clusterId)

	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dws", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	httpUrl := "v1.0/{project_id}/clusters/{cluster_id}/restart"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{cluster_id}", clusterId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
		// The "restart" parameter is a restart flag, which has no practical meaning, but is required.
		JSONBody: map[string]interface{}{
			"restart": map[string]interface{}{},
		},
	}

	_, err = client.Request("POST", path, &opt)
	if err != nil {
		return diag.Errorf("error restarting DWS cluster (%s): %s", clusterId, err)
	}

	// After the interface is sent successfully, the task ID will not be returned, so the query cluster details interface is called
	// to determine whether the restart is successful.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshClusterStateFun(client, clusterId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the status of DWS cluster (%s) to become available: %s", clusterId, err)
	}

	d.SetId(clusterId)

	return nil
}

func refreshClusterStateFun(client *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetClusterInfoByClusterId(client, clusterId)
		if err != nil {
			return nil, "ERROR", err
		}

		taskStatus := utils.PathSearch("cluster.task_status", respBody, "").(string)
		status := utils.PathSearch("cluster.status", respBody, "").(string)
		// After the restart task is completed, "task_status" will become `null` and the cluster status is available.
		if taskStatus == "" && utils.StrSliceContains([]string{"AVAILABLE", "ACTIVE"}, status) {
			return respBody, "COMPLETED", nil
		}

		if utils.StrSliceContains([]string{"REBOOT_FAILURE"}, taskStatus) {
			return respBody, taskStatus, nil
		}

		return respBody, "PENDING", nil
	}
}

func resourceClusterRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for restarting the DWS cluster. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
