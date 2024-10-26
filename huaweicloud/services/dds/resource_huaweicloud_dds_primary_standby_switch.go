package dds

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS POST /v3/{project_id}/instances/{instance_id}/switchover
// @API DDS POST /v3/{project_id}/instances/{instance_id}/nodes/{node_id}/primary
// @API DDS GET /v3/{project_id}/jobs
func ResourceDDSPrimaryStandbySwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDDSPrimaryStandbySwitchCreate,
		ReadContext:   resourceDDSPrimaryStandbySwitchRead,
		DeleteContext: resourceDDSPrimaryStandbySwitchDelete,

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
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDDSPrimaryStandbySwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	var jobID string

	if _, ok := d.GetOk("node_id"); ok {
		// switch by node, it's a synchronous task
		// but need to wait the node become primary for a few seconds
		err = promoteStandbyNodeToPrimary(client, d)
		if err != nil {
			return diag.FromErr(err)
		}

		// lintignore:R018
		time.Sleep(30 * time.Second)
	} else {
		// switch instance, it's a asynchronous task
		jobID, err = performSwitchForInstance(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	// wait for job complete
	if jobID != "" {
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"Running"},
			Target:       []string{"Completed"},
			Refresh:      JobStateRefreshFunc(client, jobID),
			Timeout:      d.Timeout(schema.TimeoutCreate),
			Delay:        10 * time.Second,
			PollInterval: 10 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for the job (%s) completed: %s ", jobID, err)
		}
	}

	return nil
}

func promoteStandbyNodeToPrimary(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	createHttpUrl := "v3/{project_id}/instances/{instance_id}/nodes/{node_id}/primary"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createPath = strings.ReplaceAll(createPath, "{node_id}", d.Get("node_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error promoting standby node to primary: %s", err)
	}

	return nil
}

func performSwitchForInstance(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	createHttpUrl := "v3/{project_id}/instances/{instance_id}/switchover"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return "", fmt.Errorf("error performing primary standby switch for instance: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return "", fmt.Errorf("error flattening response: %s", err)
	}

	jobID := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobID == "" {
		return "", fmt.Errorf("unable to find job ID in API response")
	}

	return jobID, nil
}

func resourceDDSPrimaryStandbySwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDDSPrimaryStandbySwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting primary standby switch resource is not supported. The resource is only removed from the state," +
		" the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
