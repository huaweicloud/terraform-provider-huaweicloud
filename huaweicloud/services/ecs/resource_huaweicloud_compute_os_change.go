package ecs

import (
	"context"
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

var osChangeNonUpdatableParams = []string{"cloud_init_installed", "server_id", "os_change", "os_change.*.imageid",
	"os_change.*.adminpass", "os_change.*.keyname", "os_change.*.userid", "os_change.*.mode", "os_change.*.metadata",
	"os_change.*.metadata.*.user_data", "os_change.*.metadata.*.__system__encrypted",
	"os_change.*.metadata.*.__system__cmkid"}

// @API ECS POST /v1/{project_id}/cloudservers/{server_id}/changeos
// @API ECS POST /v2/{project_id}/cloudservers/{server_id}/changeos
// @API ECS GET /v1/{project_id}/jobs/{job_id}
func ResourceComputeOsChange() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeOsChangeCreate,
		ReadContext:   resourceComputeOsChangeRead,
		UpdateContext: resourceComputeOsChangeUpdate,
		DeleteContext: resourceComputeOsChangeDelete,

		CustomizeDiff: config.FlexibleForceNew(osChangeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cloud_init_installed": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"os_change": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     osChangeSchema(),
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

func osChangeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"imageid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"adminpass": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keyname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"userid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     osChangeMetadataSchema(),
			},
		},
	}
}

func osChangeMetadataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"__system__encrypted": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"__system__cmkid": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceComputeOsChangeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrlV1 = "v1/{project_id}/cloudservers/{server_id}/changeos"
		httpUrlV2 = "v2/{project_id}/cloudservers/{server_id}/changeos"
		product   = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	httpUrl := httpUrlV1
	if v := d.Get("cloud_init_installed"); v == "true" {
		httpUrl = httpUrlV2
	}
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{server_id}", d.Get("server_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = buildOsChangeBodyParams(d)

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ECS OS change: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("server_id").(string))

	jobID := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobID == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      getJobRefreshFunc(client, jobID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for ECS OS change: %s", err)
	}

	return nil
}

func buildOsChangeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"os-change": utils.RemoveNil(buildOsChangeOsChangeBodyParams(d.Get("os_change"))),
	}
	return bodyParams
}

func buildOsChangeOsChangeBodyParams(osChangeRaw interface{}) map[string]interface{} {
	if osChangeRaw == nil || len(osChangeRaw.([]interface{})) == 0 {
		return nil
	}

	osChange := osChangeRaw.([]interface{})[0]
	if v, ok := osChange.(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"imageid":   v["imageid"],
			"adminpass": utils.ValueIgnoreEmpty(v["adminpass"]),
			"keyname":   utils.ValueIgnoreEmpty(v["keyname"]),
			"userid":    utils.ValueIgnoreEmpty(v["userid"]),
			"mode":      utils.ValueIgnoreEmpty(v["mode"]),
			"isAutoPay": "true",
			"metadata":  buildOsChangeOsChangeMetadataBodyParams(v["metadata"]),
		}
		return bodyParams
	}
	return nil
}

func buildOsChangeOsChangeMetadataBodyParams(metadataRaw interface{}) map[string]interface{} {
	if metadataRaw == nil || len(metadataRaw.([]interface{})) == 0 {
		return nil
	}

	osChange := metadataRaw.([]interface{})[0]
	if v, ok := osChange.(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"user_data":           utils.ValueIgnoreEmpty(v["user_data"]),
			"__system__encrypted": utils.ValueIgnoreEmpty(v["__system__encrypted"]),
			"__system__cmkid":     utils.ValueIgnoreEmpty(v["__system__cmkid"]),
		}
		return bodyParams
	}
	return nil
}

func resourceComputeOsChangeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceComputeOsChangeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceComputeOsChangeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting ECS OS change resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
