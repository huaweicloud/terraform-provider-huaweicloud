package dcs

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var dcsNodeStatusChangeNonUpdatableParams = []string{"instance_id", "node_ids", "action"}

// @API DCS PUT /v2/{project_id}/instances/{instance_id}/nodes/status
// @API DCS GET /v2/{project_id}/jobs/{job_id}
func ResourceDcsNodeStatusChange() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsNodeStatusChangeCreate,
		ReadContext:   resourceDcsNodeStatusChangeRead,
		UpdateContext: resourceDcsNodeStatusChangeUpdate,
		DeleteContext: resourceDcsNodeStatusChangeDelete,

		CustomizeDiff: config.FlexibleForceNew(dcsNodeStatusChangeNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"action": {
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

func resourceDcsNodeStatusChangeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/nodes/status"
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateDcsNodeStatusChangeBodyParams(d))

	createResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DCS node status change: %s", err)
	}

	defer createResp.Body.Close()

	body, err := io.ReadAll(createResp.Body)
	if err != nil {
		return diag.Errorf("error reading response body: %s", err)
	}

	jobId := strings.TrimPrefix(strings.TrimSpace(string(body)), "job_id:")
	jobId = strings.TrimSpace(jobId)

	if jobId == "" {
		return diag.Errorf("error creating DCS node status change: job_id is not found in API response")
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID.String())

	err = checkDcsInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildCreateDcsNodeStatusChangeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action": d.Get("action"),
	}

	if nodeIds, ok := d.GetOk("node_ids"); ok {
		bodyParams["node_ids"] = nodeIds
	}

	return bodyParams
}

func resourceDcsNodeStatusChangeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsNodeStatusChangeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsNodeStatusChangeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS node status change resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
