package dds

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS POST /v3/{project_id}/instances/recovery
// @API DDS GET /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/jobs
func ResourceDDSInstanceRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDDSInstanceRestoreCreate,
		ReadContext:   resourceDDSInstanceRestoreRead,
		DeleteContext: resourceDDSInstanceRestoreDelete,

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
			"target_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"backup_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"restore_time", "backup_id"},
			},
			"restore_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func buildInstanceRestoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	source := map[string]interface{}{
		"instance_id": d.Get("source_id"),
	}
	if backupID, ok := d.GetOk("backup_id"); ok {
		source["backup_id"] = backupID
	} else {
		source["type"] = "timestamp"
		source["restore_time"] = d.Get("restore_time")
	}
	opts := map[string]interface{}{
		"source": source,
		"target": map[string]interface{}{
			"instance_id": d.Get("target_id"),
		},
	}
	return opts
}

func resourceDDSInstanceRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS client: %s ", err)
	}

	// restore instance
	restoreHttpUrl := "v3/{project_id}/instances/recovery"
	restorePath := client.Endpoint + restoreHttpUrl
	restorePath = strings.ReplaceAll(restorePath, "{project_id}", client.ProjectID)
	restoreOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildInstanceRestoreBodyParams(d),
	}

	// retry
	retryFunc := func() (interface{}, bool, error) {
		resp, err := client.Request("POST", restorePath, &restoreOpt)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	instId := d.Get("target_id").(string)
	restoreResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error restoring to instance(%s): %s", instId, err)
	}

	// get job ID
	restoreRespBody, err := utils.FlattenResponse(restoreResp.(*http.Response))
	if err != nil {
		return diag.Errorf("error flatten response: %s", err)
	}
	jobID := utils.PathSearch("job_id", restoreRespBody, "")
	if jobID.(string) == "" {
		return diag.Errorf("error restoring to instance(%s): job_id not found", instId)
	}

	// wait for job complete
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      JobStateRefreshFunc(client, jobID.(string)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        60 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the job (%s) completed: %s ", jobID.(string), err)
	}

	d.SetId(instId)

	return nil
}

func resourceDDSInstanceRestoreRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDDSInstanceRestoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restore resource is not supported. The restore resource is only removed from the state," +
		" the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
