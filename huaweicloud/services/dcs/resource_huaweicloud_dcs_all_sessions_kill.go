package dcs

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var dcsAllSessionsKillNonUpdatableParams = []string{"instance_id", "node_id", "kill_all_nodes"}

// @API DCS POST /v2/{project_id}/instances/{instance_id}/clients/kill-all
func ResourceDcsAllSessionsKill() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsAllSessionsKillCreate,
		ReadContext:   resourceDcsAllSessionsKillRead,
		UpdateContext: resourceDcsAllSessionsKillUpdate,
		DeleteContext: resourceDcsAllSessionsKillDelete,

		CustomizeDiff: config.FlexibleForceNew(dcsAllSessionsKillNonUpdatableParams),

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
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kill_all_nodes": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
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

func resourceDcsAllSessionsKillCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/clients/kill-all"
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateDcsAllSessionsKillBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DCS all sessoions kill: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)

	if err != nil {
		return diag.FromErr(err)
	}

	// Kill all sessions API is asynchronous and returns a job_id.
	// Need to wait for the job to complete before finishing Create.
	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error killing all sessions of instance(%s): job_id is not found in API response", instanceID)
	}

	err = checkDcsInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	return nil
}

func resourceDcsAllSessionsKillRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsAllSessionsKillUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsAllSessionsKillDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS all sessoions kill resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func buildCreateDcsAllSessionsKillBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := make(map[string]interface{})
	if v, ok := d.GetOk("node_id"); ok {
		body["node_id"] = v.(string)
	}
	if v, ok := d.GetOk("kill_all_nodes"); ok {
		body["kill_all_nodes"] = v.(string) == "true"
	}
	return body
}
