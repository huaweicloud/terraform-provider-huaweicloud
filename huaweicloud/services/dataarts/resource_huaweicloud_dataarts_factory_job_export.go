package dataarts

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

var factoryJobExportNonUpdatableParams = []string{
	"job_name",
	"obs_path",
	"workspace_id",
	"export_depend",
	"export_status",
}

// @API DataArtsStudio POST /v1/{project_id}/jobs/{job_name}/export
func ResourceFactoryJobExport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFactoryJobExportCreate,
		ReadContext:   resourceFactoryJobExportRead,
		UpdateContext: resourceFactoryJobExportUpdate,
		DeleteContext: resourceFactoryJobExportDelete,

		CustomizeDiff: config.FlexibleForceNew(factoryJobExportNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the job is located.`,
			},

			// Required parameters.
			"job_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The job name.`,
			},
			"obs_path": {
				Type:     schema.TypeString,
				Optional: true,
				Description: utils.SchemaDesc(
					`The OBS target path for storing the exported job package.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},

			// Optional parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The workspace ID to which the job belongs.`,
			},
			"export_depend": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to export the scripts and resources depended on by the job.`,
			},
			"export_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The job status to export.`,
			},

			// Attributes.
			"folder_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The folder path suffix in the returned OBS path relative to the input OBS path.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true,
						Required: true,
					}),
			},
		},
	}
}

func buildFactoryJobExportBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawConfig := d.GetRawConfig()
	return utils.RemoveNil(map[string]interface{}{
		// When the export_depend parameter is explicitly filled with 'true' or 'false', the exportDepend parameter must be set,
		// and the boolean default value 'false' is not used.
		"exportDepend": utils.ValueIgnoreEmpty(utils.GetNestedObjectFromRawConfig(rawConfig, "export_depend")),
		"obsPath":      utils.ValueIgnoreEmpty(d.Get("obs_path")),
		"exportStatus": utils.ValueIgnoreEmpty(d.Get("export_status")),
	})
}

func buildFactoryJobExportRequestMoreHeaders(workspaceId string) map[string]string {
	moreHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	if workspaceId != "" {
		moreHeaders["workspace"] = workspaceId
	}

	return moreHeaders
}

func buildFactoryJobExportFolderPath(inputObsPath, outputObsPath string) string {
	if outputObsPath == "" {
		return ""
	}
	if inputObsPath == "" {
		return outputObsPath
	}

	normalizedInputObsPath := inputObsPath
	if !strings.HasSuffix(normalizedInputObsPath, "/") {
		normalizedInputObsPath += "/"
	}

	if strings.HasPrefix(outputObsPath, normalizedInputObsPath) {
		return strings.TrimPrefix(outputObsPath, normalizedInputObsPath)
	}
	return outputObsPath
}

func exportFactoryJob(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl     = "v1/{project_id}/jobs/{job_name}/export"
		exportPath  = client.Endpoint + httpUrl
		jobName     = d.Get("job_name").(string)
		workspaceId = d.Get("workspace_id").(string)
	)

	exportPath = strings.ReplaceAll(exportPath, "{project_id}", client.ProjectID)
	exportPath = strings.ReplaceAll(exportPath, "{job_name}", jobName)
	exportOpts := golangsdk.RequestOpts{
		JSONBody:         utils.RemoveNil(buildFactoryJobExportBodyParams(d)),
		KeepResponseBody: true,
		MoreHeaders:      buildFactoryJobExportRequestMoreHeaders(workspaceId),
		OkCodes:          []int{200},
	}

	resp, err := client.Request("POST", exportPath, &exportOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceFactoryJobExportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		jobName = d.Get("job_name").(string)
	)

	client, err := cfg.NewServiceClient("dataarts-dlf", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := exportFactoryJob(client, d)
	if err != nil {
		return diag.Errorf("error exporting DataArts Factory job (%s): %s", jobName, err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	outputObsPath := utils.PathSearch("obsPath", respBody, "").(string)
	if err := d.Set("folder_path", buildFactoryJobExportFolderPath(d.Get("obs_path").(string), outputObsPath)); err != nil {
		return diag.FromErr(err)
	}

	return resourceFactoryJobExportRead(ctx, d, meta)
}

func resourceFactoryJobExportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceFactoryJobExportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceFactoryJobExportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to export a job package to a specified OBS storage
path. Deleting this resource will not clear the export task, but will only remove the resource information from the
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
