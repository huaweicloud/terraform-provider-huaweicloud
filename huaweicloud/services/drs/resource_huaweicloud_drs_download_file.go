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

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS POST /v5/{project_id}/download
func ResourceDownloadFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDownloadFileCreate,
		ReadContext:   resourceDownloadFileRead,
		UpdateContext: resourceDownloadFileUpdate,
		DeleteContext: resourceDownloadFileDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"files": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildDownloadFileBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"files": utils.ValueIgnoreEmpty(d.Get("files")),
	}
}

func resourceDownloadFileCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/download"
	)

	client, err := cfg.DrsV5Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildDownloadFileBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error downloading DRS files: %s", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.Errorf("error reading download file response body: %s", err)
	}

	outputFile, err := getDownloadFileOutputPath(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := os.WriteFile(outputFile, bodyBytes, 0600); err != nil {
		return diag.Errorf("failed to write downloaded file to (%s): %s", outputFile, err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return nil
}

func getDownloadFileOutputPath(d *schema.ResourceData) (string, error) {
	files := d.Get("files").([]interface{})
	if len(files) == 0 {
		return "", fmt.Errorf("no file name is specified and the 'files' parameter is empty")
	}
	return files[0].(string), nil
}

func resourceDownloadFileRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceDownloadFileUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceDownloadFileDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to download the exported DRS files. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
