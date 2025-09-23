package dws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDM POST /v2/{project_id}/clusters/{cluster_id}/logical-clusters/{logical_cluster_id}/restart
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/logical-clusters
func ResourceLogicalClusterRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogicalClusterRestartCreate,
		ReadContext:   resourceLogicalClusterRestartRead,
		DeleteContext: resourceLogicalClusterRestartDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
				Description: `Specifies the DWS cluster ID to which the logical cluster belongs.`,
			},
			"logical_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the logical cluster to be restarted.",
			},
		},
	}
}

func resourceLogicalClusterRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterId := d.Get("cluster_id").(string)
	// For the same DWS cluster, it is not supported to run multiple tasks at the same time.
	config.MutexKV.Lock(clusterId)
	defer config.MutexKV.Unlock(clusterId)

	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dws", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	httpUrl := "v2/{project_id}/clusters/{cluster_id}/logical-clusters/{logical_cluster_id}/restart"
	logicaClusterId := d.Get("logical_cluster_id").(string)
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{cluster_id}", clusterId)
	path = strings.ReplaceAll(path, "{logical_cluster_id}", logicaClusterId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}

	resp, err := client.Request("POST", path, &opt)
	if err != nil {
		return diag.Errorf("error restarting logical cluster (%s): %s", logicaClusterId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// When the restart interface is sent again during the restart process, the status code is "DWS.9999" and the error message is "DWS.7023".
	errorMsg := utils.PathSearch("error_msg", respBody, "").(string)
	if errorMsg != "" {
		return diag.Errorf("error restarting logical cluster: %s", errorMsg)
	}

	_, err = waitingForRestartCompleted(ctx, client, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for restarting the logical cluster (%s) to complete: %s", logicaClusterId, err)
	}

	d.SetId(logicaClusterId)

	return nil
}

func waitingForRestartCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) (interface{}, error) {
	logicaClusterId := d.Get("logical_cluster_id").(string)
	expression := fmt.Sprintf("logical_clusters[?logical_cluster_id=='%s']|[0]", logicaClusterId)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshLogicalClusterStateFun(client, d, expression),
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func refreshLogicalClusterStateFun(client *golangsdk.ServiceClient, d *schema.ResourceData, expression string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clusterRespBody, err := readLogicalClusters(client, d)
		if err != nil {
			return nil, "ERROR", err
		}

		cluster := utils.PathSearch(expression, clusterRespBody, nil)
		if cluster == nil {
			return nil, "ERROR", golangsdk.ErrDefault404{}
		}

		completed := utils.PathSearch("action_info.completed", cluster, false).(bool)
		result := utils.PathSearch("action_info.result", cluster, "").(string)
		status := utils.PathSearch("status", cluster, "").(string)
		// In some cases, when the restart task is completed, the logical cluster status is unavailable or abnormal.
		if completed && result == "success" && status == "Normal" {
			return cluster, "COMPLETED", nil
		}

		if completed && result == "failed" {
			errMsg := fmt.Errorf("the logical cluster (%s) is restarted failed", utils.PathSearch("logical_cluster_name", cluster, "").(string))
			return cluster, "ERROR", errMsg
		}
		return cluster, "PENDING", nil
	}
}

func resourceLogicalClusterRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceLogicalClusterRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for restarting the logical cluster. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
