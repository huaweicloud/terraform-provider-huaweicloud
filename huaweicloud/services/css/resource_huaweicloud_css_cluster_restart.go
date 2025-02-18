package css

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API CSS POST /v2.0/{project_id}/clusters/{cluster_id}/restart
// @API CSS POST /v2.0/{project_id}/clusters/{cluster_id}/rolling_restart
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}
func ResourceCssClusterRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCssClusterRestartCreate,
		ReadContext:   resourceCssClusterRestartRead,
		DeleteContext: resourceCssClusterRestartDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_rolling": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCssClusterRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterId := d.Get("cluster_id").(string)
	restartType := d.Get("type").(string)
	value := d.Get("value").(string)
	v1client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	v2client, err := conf.CssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V2 client: %s", err)
	}

	// Check whether the cluster status is available.
	err = checkClusterOperationResult(ctx, v1client, clusterId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("is_rolling").(bool) {
		err := rollingRestartCluster(v2client, clusterId, restartType, value)
		if err != nil {
			return diag.Errorf("error rolling restart CSS cluster, err: %s", err)
		}
	} else {
		err := restartCluster(v2client, clusterId, restartType, value)
		if err != nil {
			return diag.Errorf("error restart CSS cluster, err: %s", err)
		}
	}

	// Check whether the cluster restart is complete
	err = checkClusterOperationResult(ctx, v1client, clusterId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterId)

	return nil
}

func restartCluster(v2client *golangsdk.ServiceClient, clusterId, restartType, value string) error {
	restartClusterHttpUrl := "v2.0/{project_id}/clusters/{cluster_id}/restart"
	restartClusterPath := v2client.Endpoint + restartClusterHttpUrl
	restartClusterPath = strings.ReplaceAll(restartClusterPath, "{project_id}", v2client.ProjectID)
	restartClusterPath = strings.ReplaceAll(restartClusterPath, "{cluster_id}", clusterId)

	restartClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"type":  restartType,
			"value": value,
		},
	}
	_, err := v2client.Request("POST", restartClusterPath, &restartClusterOpt)
	if err != nil {
		return err
	}

	return nil
}

func rollingRestartCluster(v2client *golangsdk.ServiceClient, clusterId, restartType, value string) error {
	rollingRestartClusterHttpUrl := "v2.0/{project_id}/clusters/{cluster_id}/rolling_restart"
	rollingRestartClusterPath := v2client.Endpoint + rollingRestartClusterHttpUrl
	rollingRestartClusterPath = strings.ReplaceAll(rollingRestartClusterPath, "{project_id}", v2client.ProjectID)
	rollingRestartClusterPath = strings.ReplaceAll(rollingRestartClusterPath, "{cluster_id}", clusterId)

	rollingRestartClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"type":  restartType,
			"value": value,
		},
	}
	_, err := v2client.Request("POST", rollingRestartClusterPath, &rollingRestartClusterOpt)
	if err != nil {
		return err
	}

	return nil
}

func resourceCssClusterRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCssClusterRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restart resource is not supported. The restart resource is only removed from the state," +
		" the cluster instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
