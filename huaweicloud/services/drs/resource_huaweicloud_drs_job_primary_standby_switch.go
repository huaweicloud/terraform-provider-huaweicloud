package drs

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
	"github.com/chnsz/golangsdk/openstack/drs/v3/jobs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DRS POST /v3/{project_id}/jobs/batch-switchover
// @API DRS GET /v3/{project_id}/jobs
func ResourceDRSPrimaryStandbySwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDRSPrimaryStandbySwitchCreate,
		ReadContext:   resourceDRSPrimaryStandbySwitchRead,
		DeleteContext: resourceDRSPrimaryStandbySwitchDelete,

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
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDRSPrimaryStandbySwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DrsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS v3 client, error: %s", err)
	}

	createHttpUrl := "v3/{project_id}/jobs/batch-switchover"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"jobs": []string{d.Get("job_id").(string)},
		},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error performing primary standby switch for job: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	// wait for switch complete
	err = waitingforJobActionComplete(ctx, client, d.Get("job_id").(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitingforJobActionComplete(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		// when the action completed, `current_action` is not in return
		Pending: []string{"SWITCH_OVER"},
		Target:  []string{""},
		Refresh: func() (interface{}, string, error) {
			resp, err := jobs.List(client, jobs.ListJobsReq{
				CurPage:   1,
				PerPage:   1,
				Name:      id,
				DbUseType: "cloudDataGuard",
			})
			if err != nil {
				return nil, "FAILED", err
			}
			if len(resp.Jobs) == 0 {
				return resp, "FAILED", fmt.Errorf("unable to get job info from API")
			}

			return resp, resp.Jobs[0].JobAction.CurrentAction, nil
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DRS job (%s) primary standby switch to be completed: %s", id, err)
	}

	return nil
}

func resourceDRSPrimaryStandbySwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDRSPrimaryStandbySwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting primary standby switch resource is not supported. The resource is only removed from the state," +
		" the job remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
