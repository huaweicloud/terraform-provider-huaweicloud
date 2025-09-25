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

var osReinstallNonUpdatableParams = []string{"cloud_init_installed", "server_id", "os_reinstall", "os_reinstall.*.adminpass",
	"os_reinstall.*.keyname", "os_reinstall.*.userid", "os_reinstall.*.mode", "os_reinstall.*.metadata",
	"os_reinstall.*.metadata.*.user_data", "os_reinstall.*.metadata.*.__system__encrypted",
	"os_reinstall.*.metadata.*.__system__cmkid"}

// @API ECS POST /v1/{project_id}/cloudservers/{server_id}/reinstallos
// @API ECS POST /v2/{project_id}/cloudservers/{server_id}/reinstallos
// @API ECS GET /v1/{project_id}/jobs/{job_id}
func ResourceComputeOsReinstall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeOsReinstallCreate,
		ReadContext:   resourceComputeOsReinstallRead,
		UpdateContext: resourceComputeOsReinstallUpdate,
		DeleteContext: resourceComputeOsReinstallDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(osReinstallNonUpdatableParams),

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
			"os_reinstall": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     osReinstallSchema(),
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

func osReinstallSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
				Elem:     osReinstallMetadataSchema(),
			},
		},
	}
}

func osReinstallMetadataSchema() *schema.Resource {
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

func resourceComputeOsReinstallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrlV1 = "v1/{project_id}/cloudservers/{server_id}/reinstallos"
		httpUrlV2 = "v2/{project_id}/cloudservers/{server_id}/reinstallos"
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
	createOpt.JSONBody = buildOsReinstallBodyParams(d)

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ECS OS reinstall: %s", err)
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
		return diag.Errorf("error waiting for ECS OS reinstall: %s", err)
	}

	return nil
}

func buildOsReinstallBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"os-reinstall": utils.RemoveNil(buildOsReinstallOsReinstallBodyParams(d.Get("os_reinstall"))),
	}
	return bodyParams
}

func buildOsReinstallOsReinstallBodyParams(osReinstallRaw interface{}) map[string]interface{} {
	if osReinstallRaw == nil || len(osReinstallRaw.([]interface{})) == 0 {
		return nil
	}

	osReinstall := osReinstallRaw.([]interface{})[0]
	if v, ok := osReinstall.(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"adminpass": utils.ValueIgnoreEmpty(v["adminpass"]),
			"keyname":   utils.ValueIgnoreEmpty(v["keyname"]),
			"userid":    utils.ValueIgnoreEmpty(v["userid"]),
			"mode":      utils.ValueIgnoreEmpty(v["mode"]),
			"metadata":  buildOsReinstallOsReinstallMetadataBodyParams(v["metadata"]),
		}
		return bodyParams
	}
	return nil
}

func buildOsReinstallOsReinstallMetadataBodyParams(metadataRaw interface{}) map[string]interface{} {
	if metadataRaw == nil || len(metadataRaw.([]interface{})) == 0 {
		return nil
	}

	osReinstall := metadataRaw.([]interface{})[0]
	if v, ok := osReinstall.(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"user_data":           utils.ValueIgnoreEmpty(v["user_data"]),
			"__system__encrypted": utils.ValueIgnoreEmpty(v["__system__encrypted"]),
			"__system__cmkid":     utils.ValueIgnoreEmpty(v["__system__cmkid"]),
		}
		return bodyParams
	}
	return nil
}

func resourceComputeOsReinstallRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceComputeOsReinstallUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceComputeOsReinstallDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting ECS OS reinstall resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
