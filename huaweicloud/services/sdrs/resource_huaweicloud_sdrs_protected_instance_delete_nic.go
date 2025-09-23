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

// @API SDRS POST /v1/{project_id}/protected-instances/{protected_instance_id}/nic/delete
// @API SDRS GET /v1/{project_id}/jobs/{job_id}
func ResourceProtectedInstanceDeleteNIC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProtectedInstanceDeleteNICCreate,
		ReadContext:   resourceProtectedInstanceDeleteNICRead,
		UpdateContext: resourceProtectedInstanceDeleteNICUpdate,
		DeleteContext: resourceProtectedInstanceDeleteNICDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"protected_instance_id",
			"nic_id",
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level
region will be used.`,
			},
			"protected_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the protected instance to delete the NIC from.`,
			},
			"nic_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the NIC to delete.`,
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

func resourceProtectedInstanceDeleteNICCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/protected-instances/{protected_instance_id}/nic/delete"
		product = "sdrs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{protected_instance_id}", d.Get("protected_instance_id").(string))

	requestBody := map[string]interface{}{
		"nic_id": d.Get("nic_id"),
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         requestBody,
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting NIC from SDRS protected instance: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobID := utils.PathSearch("job_id", respBody, "").(string)
	if jobID == "" {
		return diag.Errorf("error deleting NIC from SDRS protected instance: job ID not found in API response")
	}

	if err := waitingForDeleteNICSuccess(ctx, client, d.Timeout(schema.TimeoutCreate), jobID); err != nil {
		return diag.Errorf("error waiting for SDRS NIC deletion to complete: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	return resourceProtectedInstanceDeleteNICRead(ctx, d, meta)
}

func waitingForDeleteNICSuccess(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration, jobID string) error {
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

func resourceProtectedInstanceDeleteNICRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceProtectedInstanceDeleteNICUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceProtectedInstanceDeleteNICDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to delete a NIC from a protected instance. Deleting this 
resource will not change the current NIC configuration, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
