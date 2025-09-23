package cce

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

// @API CCE POST /api/v3/projects/{project_id}/clusters/{cluster_id}/rotatecredentials
// @API CCE GET /api/v3/projects/{project_id}/jobs/{job_id}
var rotatecredentialsNonUpdatableParams = []string{"cluster_id", "component"}

func ResourceRotatecredentials() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRotatecredentialsCreate,
		ReadContext:   resourceRotatecredentialsRead,
		UpdateContext: resourceRotatecredentialsUpdate,
		DeleteContext: resourceRotatecredentialsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(rotatecredentialsNonUpdatableParams),

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
			},
			"component": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceRotatecredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	createRotatecredentialsHttpUrl := "api/v3/projects/{project_id}/clusters/{cluster_id}/rotatecredentials"
	createRotatecredentialsPath := client.Endpoint + createRotatecredentialsHttpUrl
	createRotatecredentialsPath = strings.ReplaceAll(createRotatecredentialsPath, "{project_id}", client.ProjectID)
	createRotatecredentialsPath = strings.ReplaceAll(createRotatecredentialsPath, "{cluster_id}", d.Get("cluster_id").(string))
	createRotatecredentialsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createRotatecredentialsOpt.JSONBody = map[string]interface{}{
		"component": d.Get("component"),
	}

	resp, err := client.Request("POST", createRotatecredentialsPath, &createRotatecredentialsOpt)
	if err != nil {
		return diag.Errorf("error creating CCI rotatecredentials: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	jobId := utils.PathSearch("jobid", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find job ID from API response")
	}
	d.SetId(jobId)

	err = waitForRotatecredentialsJobStatus(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRotatecredentialsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRotatecredentialsUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRotatecredentialsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `Deleting certificate rotatecredentials resource is not supported. The certificate rotatecredentials resource
		is only removed from the state.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func waitForRotatecredentialsJobStatus(ctx context.Context, client *golangsdk.ServiceClient, jobID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Initializing", "Running"},
		Target:       []string{"Success"},
		Refresh:      waitForJobStatus(client, jobID),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the status of the CCI rotatecredentials job to success: %s", err)
	}
	return nil
}
