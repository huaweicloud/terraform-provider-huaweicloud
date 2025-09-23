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

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/reboot
// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}
func ResourceLogstashClusterRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogstashClusterRestartCreate,
		ReadContext:   resourceLogstashClusterRestartRead,
		DeleteContext: resourceLogstashClusterRestartDelete,

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
		},
	}
}

func resourceLogstashClusterRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)
	client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	// Check whether the logstash cluster status is available.
	err = checkClusterOperationResult(ctx, client, clusterID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	restartLogstashHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/reboot"
	restartLogstashPath := client.Endpoint + restartLogstashHttpUrl
	restartLogstashPath = strings.ReplaceAll(restartLogstashPath, "{project_id}", client.ProjectID)
	restartLogstashPath = strings.ReplaceAll(restartLogstashPath, "{cluster_id}", clusterID)

	restartLogstashOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("POST", restartLogstashPath, &restartLogstashOpt)
	if err != nil {
		return diag.Errorf("error restart CSS logstash cluster: %s", err)
	}

	// Check whether the logstash cluster restart is complete
	err = checkClusterOperationResult(ctx, client, clusterID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterID)

	return nil
}

func resourceLogstashClusterRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceLogstashClusterRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restart resource is not supported. The restart resource is only removed from the state," +
		" the cluster instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
