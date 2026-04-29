package drs

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var downloadBatchCreateTemplateNonUpdatableParams = []string{"engine_type", "template_file_name"}

// @API DRS GET /v5/{project_id}/jobs/template
func ResourceDownloadBatchCreateTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDownloadBatchCreateTemplateCreate,
		ReadContext:   resourceDownloadBatchCreateTemplateRead,
		UpdateContext: resourceDownloadBatchCreateTemplateUpdate,
		DeleteContext: resourceDownloadBatchCreateTemplateDelete,

		CustomizeDiff: config.FlexibleForceNew(downloadBatchCreateTemplateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_file_name": {
				Type:     schema.TypeString,
				Optional: true,
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

func buildDownloadBatchCreateTemplateQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("engine_type"); ok && v.(string) != "" {
		return fmt.Sprintf("?engine_type=%s", v)
	}

	return ""
}

func resourceDownloadBatchCreateTemplateCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/template"
	)

	client, err := cfg.DrsV5Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS v5 client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDownloadBatchCreateTemplateQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error downloading DRS batch create template: %s", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.Errorf("error reading template download response body: %s", err)
	}

	var outputFile string
	if v, ok := d.GetOk("template_file_name"); ok {
		outputFile = v.(string)
		if !strings.HasSuffix(outputFile, ".zip") {
			outputFile += ".zip"
		}
	} else {
		outputFile = "drs-batch-create-template.zip"
	}

	if err := os.WriteFile(outputFile, bodyBytes, 0600); err != nil {
		return diag.Errorf("failed to write batch create template file to (%s): %s", outputFile, err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return nil
}

func resourceDownloadBatchCreateTemplateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceDownloadBatchCreateTemplateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceDownloadBatchCreateTemplateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to download the batch import task template. Deleting
    this resource will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
