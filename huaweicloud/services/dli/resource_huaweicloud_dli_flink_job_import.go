package dli

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var flinkJobImportNonUpdatableParams = []string{
	"obs_path",
	"is_cover",
}

// @API DLI POST /v1.0/{project_id}/streaming/jobs/import
func ResourceFlinkJobImport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlinkJobImportCreate,
		ReadContext:   resourceFlinkJobImportRead,
		UpdateContext: resourceFlinkJobImportUpdate,
		DeleteContext: resourceFlinkJobImportDelete,

		CustomizeDiff: config.FlexibleForceNew(flinkJobImportNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the imported flink jobs are located.`,
			},

			// Required
			"obs_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The OBS path for the imported job zip file.`,
			},

			// Optional
			"is_cover": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to overwrite existing jobs with the same name.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
}

func buildFlinkJobImportBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"zip_file": d.Get("obs_path"),
	}

	if v, ok := d.GetOk("is_cover"); ok {
		bodyParams["is_cover"] = v
	}

	return bodyParams
}

func resourceFlinkJobImportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	httpUrl := "v1.0/{project_id}/streaming/jobs/import"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildFlinkJobImportBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error importing DLI Flink jobs: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// is_success: the type of the API is string
	if utils.PathSearch("is_success", respBody, "false").(string) == "false" {
		return diag.Errorf("unable to import the jobs: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceFlinkJobImportRead(ctx, d, meta)
}

func resourceFlinkJobImportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFlinkJobImportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFlinkJobImportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for importing Flink jobs. Deleting this resource will not
clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
