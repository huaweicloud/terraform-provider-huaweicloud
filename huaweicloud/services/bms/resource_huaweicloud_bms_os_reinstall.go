package bms

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var osReinstallNonUpdatableParams = []string{"server_id", "os_reinstall", "os_reinstall.*.admin_pass",
	"os_reinstall.*.key_name", "os_reinstall.*.user_id", "os_reinstall.*.metadata", "os_reinstall.*.metadata.*.user_data",
	"os_reinstall.*.metadata.*.__system__encrypted", "os_reinstall.*.metadata.*.__system__cmkid",
	"os_reinstall.*.metadata.*.__system__encrypted"}

// @API BMS POST /v1/{project_id}/baremetalservers/{server_id}/reinstallos
// @API BMS GET /v1/{project_id}/jobs/{job_id}
func ResourceBmsOsReinstall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBmsOsReinstallCreate,
		ReadContext:   resourceBmsOsReinstallRead,
		UpdateContext: resourceBmsOsReinstallUpdate,
		DeleteContext: resourceBmsOsReinstallDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(osReinstallNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
			"admin_pass": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": {
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
			"__system__encryption_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceBmsOsReinstallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/reinstallos"
		product = "bms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating BMS client: %s", err)
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
		return diag.Errorf("error creating BMS OS reinstall: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("server_id").(string))

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	err = waitForJobComplete(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for BMS OS reinstall: %s", err)
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
			"adminpass": utils.ValueIgnoreEmpty(v["admin_pass"]),
			"keyname":   utils.ValueIgnoreEmpty(v["key_name"]),
			"userid":    utils.ValueIgnoreEmpty(v["user_id"]),
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
		userData := v["user_data"].(string)
		if _, err := base64.StdEncoding.DecodeString(userData); err != nil {
			userData = base64.StdEncoding.EncodeToString([]byte(userData))
		}
		bodyParams := map[string]interface{}{
			"user_data":                      userData,
			"__system__encrypted":            utils.ValueIgnoreEmpty(v["__system__encrypted"]),
			"__system__cmkid":                utils.ValueIgnoreEmpty(v["__system__cmkid"]),
			"__system__encryption_algorithm": utils.ValueIgnoreEmpty(v["__system__encryption_algorithm"]),
		}
		return bodyParams
	}
	return nil
}

func resourceBmsOsReinstallRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBmsOsReinstallUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBmsOsReinstallDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting BMS OS reinstall resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
