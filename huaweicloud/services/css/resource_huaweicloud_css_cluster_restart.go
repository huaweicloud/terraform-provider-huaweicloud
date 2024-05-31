package css

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cssv2model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/css/v2/model"

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
	clusterID := d.Get("cluster_id").(string)
	cssV1Client, err := conf.HcCssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}
	cssV2Client, err := conf.HcCssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V2 client: %s", err)
	}

	// Check whether the cluster status is available.
	err = checkClusterOperationCompleted(ctx, cssV1Client, clusterID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("is_rolling").(bool) {
		rollingRestartClusterOpts := cssv2model.RollingRestartRequest{
			ClusterId: clusterID,
			Body: &cssv2model.RollingRestartReq{
				Type:  d.Get("type").(string),
				Value: d.Get("value").(string),
			},
		}

		_, err = cssV2Client.RollingRestart(&rollingRestartClusterOpts)
		if err != nil {
			return diag.Errorf("error rolling restart CSS cluster, err: %s", err)
		}
	} else {
		restartClusterOpts := cssv2model.RestartClusterRequest{
			ClusterId: clusterID,
			Body: &cssv2model.RestartClusterReq{
				Type:  d.Get("type").(string),
				Value: d.Get("value").(string),
			},
		}

		_, err = cssV2Client.RestartCluster(&restartClusterOpts)
		if err != nil {
			return diag.Errorf("error restart CSS cluster, err: %s", err)
		}
	}

	// Check whether the cluster restart is complete
	err = checkClusterOperationCompleted(ctx, cssV1Client, clusterID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterID)

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
