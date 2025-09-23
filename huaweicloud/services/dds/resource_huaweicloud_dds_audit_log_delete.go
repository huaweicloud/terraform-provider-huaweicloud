package dds

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS DELETE /v3/{project_id}/instances/{instance_id}/auditlog
// @API DDS GET /v3/{project_id}/instances
func ResourceDDSAuditLogDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDDSAuditLogDeleteCreate,
		ReadContext:   resourceDDSAuditLogDeleteRead,
		DeleteContext: resourceDDSAuditLogDeleteDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
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
			"file_names": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDDSAuditLogDeleteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	instId := d.Get("instance_id").(string)

	deleteHttpUrl := "v3/{project_id}/instances/{instance_id}/auditlog"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAuditLogDeleteBodyParams(d)),
	}

	// retry
	retryFunc := func() (interface{}, bool, error) {
		resp, err := client.Request("DELETE", deletePath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	deleteResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(client, instId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error deleting audit log: %s", err)
	}

	// job ID is the ID for background task that are not perceived by users
	deleteRespBody, err := utils.FlattenResponse(deleteResp.(*http.Response))
	if err != nil {
		return diag.Errorf("error flattening response: %s", err)
	}
	jobID := utils.PathSearch("job_id", deleteRespBody, "").(string)
	if jobID == "" {
		return diag.Errorf("unable to find job ID from API response")
	}

	d.SetId(jobID)

	return nil
}

func buildAuditLogDeleteBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"file_names": d.Get("file_names"),
	}
}

func resourceDDSAuditLogDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDDSAuditLogDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting audit log delete resource is not supported. The delete resource is only removed from the state," +
		" the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
