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

// @API SDRS POST /v1/{project_id}/protected-instances/{protected_instance_id}/resize
// @API SDRS GET /v1/{project_id}/jobs/{job_id}
func ResourceProtectedInstanceResize() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProtectedInstanceResizeCreate,
		ReadContext:   resourceProtectedInstanceResizeRead,
		UpdateContext: resourceProtectedInstanceResizeUpdate,
		DeleteContext: resourceProtectedInstanceResizeDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"protected_instance_id",
			"flavor_ref",
			"production_flavor_ref",
			"dr_flavor_ref",
			"production_dedicated_host_id",
			"dr_dedicated_host_id",
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
			"protected_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the protected instance to resize.`,
			},
			"flavor_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the flavor ID of the production and DR site servers after the modification.`,
			},
			"production_flavor_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the flavor ID of the production site server after the modification.`,
			},
			"dr_flavor_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the flavor ID of the DR site server after the modification.`,
			},
			"production_dedicated_host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the new DeH ID for the production site.`,
			},
			"dr_dedicated_host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the new DeH ID for the DR site.`,
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

func buildResizeBodyParams(d *schema.ResourceData) map[string]interface{} {
	resizeMap := map[string]interface{}{
		"flavorRef":                    utils.ValueIgnoreEmpty(d.Get("flavor_ref")),
		"production_flavorRef":         utils.ValueIgnoreEmpty(d.Get("production_flavor_ref")),
		"dr_flavorRef":                 utils.ValueIgnoreEmpty(d.Get("dr_flavor_ref")),
		"production_dedicated_host_id": utils.ValueIgnoreEmpty(d.Get("production_dedicated_host_id")),
		"dr_dedicated_host_id":         utils.ValueIgnoreEmpty(d.Get("dr_dedicated_host_id")),
	}

	return map[string]interface{}{
		"resize": resizeMap,
	}
}

func resourceProtectedInstanceResizeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/protected-instances/{protected_instance_id}/resize"
		product = "sdrs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{protected_instance_id}", d.Get("protected_instance_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildResizeBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error resizing SDRS protected instance: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobID := utils.PathSearch("job_id", respBody, "").(string)
	if jobID == "" {
		return diag.Errorf("error resizing SDRS protected instance: job ID not found in API response")
	}

	if err := waitingForResizeSuccess(ctx, client, d.Timeout(schema.TimeoutCreate), jobID); err != nil {
		return diag.Errorf("error waiting for SDRS resize to complete: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	return resourceProtectedInstanceResizeRead(ctx, d, meta)
}

func waitingForResizeSuccess(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration, jobID string) error {
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

func resourceProtectedInstanceResizeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceProtectedInstanceResizeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceProtectedInstanceResizeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to resize a protected instance. Deleting this 
resource will not change the current configuration, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
