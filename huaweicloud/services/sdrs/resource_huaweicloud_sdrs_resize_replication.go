package sdrs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SDRS POST /v1/{project_id}/replications/{replication_id}/action
// @API SDRS GET /v1/{project_id}/jobs/{job_id}
func ResourceResizeReplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResizeReplicationCreate,
		ReadContext:   resourceResizeReplicationRead,
		UpdateContext: resourceResizeReplicationUpdate,
		DeleteContext: resourceResizeReplicationDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"replication_id",
			"new_size",
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"replication_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the replication pair to resize.`,
			},
			"new_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the new size of the replication pair's disk in GB. Must be greater than the current size.`,
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

func resourceResizeReplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/replications/{replication_id}/action"
		product = "sdrs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{replication_id}", d.Get("replication_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"extend-replication": map[string]interface{}{
				"new_size": d.Get("new_size"),
			},
		},
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error resizing replication pair: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobID := utils.PathSearch("job_id", respBody, "").(string)
	if jobID == "" {
		return diag.Errorf("error resizing replication pair: job ID not found in API response")
	}

	if err := waitingForResizeReplicationSuccess(ctx, client, d.Timeout(schema.TimeoutCreate), jobID); err != nil {
		return diag.Errorf("error waiting for replication resize to complete: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	return resourceResizeReplicationRead(ctx, d, meta)
}

func waitingForResizeReplicationSuccess(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration, jobID string) error {
	unexpectedStatus := []string{"FAIL"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := queryJobStatus(client, jobID)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", respBody, "").(string)
			if status == "" {
				return nil, "ERROR", errors.New("status is not found in API response")
			}

			if status == "SUCCESS" {
				return respBody, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return respBody, status, fmt.Errorf("job failed with status: %s", status)
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceResizeReplicationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceResizeReplicationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceResizeReplicationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to resize a replication pair. Deleting this 
resource will not change the current configuration, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
