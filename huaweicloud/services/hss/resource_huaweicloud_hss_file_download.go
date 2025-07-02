package hss

import (
	"archive/zip"
	"bytes"
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

var fileDownloadNonUpdatableParams = []string{"file_id", "enterprise_project_id", "export_file_name"}

// @API HSS POST /v5/{project_id}/vul/export
func ResourceFileDownload() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFileDownloadCreate,
		ReadContext:   resourceFileDownloadRead,
		UpdateContext: resourceFileDownloadUpdate,
		DeleteContext: resourceFileDownloadDelete,

		CustomizeDiff: config.FlexibleForceNew(fileDownloadNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"file_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"export_file_name": {
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

func resourceFileDownloadCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		fileId  = d.Get("file_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/download/{file_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{file_id}", fileId)
	if epsId != "" {
		requestPath = fmt.Sprintf("%s?enterprise_project_id=%v", requestPath, epsId)
	}
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error downloading HSS export file: %s", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.Errorf("error reading response body: %s", err)
	}

	var outputFile string
	if v, ok := d.GetOk("export_file_name"); ok {
		outputFile = v.(string)
		if !strings.HasSuffix(outputFile, ".zip") {
			outputFile += ".zip"
		}
	} else {
		outputFile = fmt.Sprintf("hss-export-%s.zip", fileId)
	}

	if err := writeToZipFile(bodyBytes, outputFile); err != nil {
		return diag.Errorf(err.Error())
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return nil
}

func writeToZipFile(bodyBytes []byte, outputFile string) error {
	reader, err := zip.NewReader(bytes.NewReader(bodyBytes), int64(len(bodyBytes)))
	if err != nil {
		return fmt.Errorf("error parsing zip file: %s", err)
	}

	for _, f := range reader.File {
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("error opening zip file entry: %s", err)
		}

		_, err = io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return fmt.Errorf("error reading zip file content: %s", err)
		}
	}

	err = os.WriteFile(outputFile, bodyBytes, 0600)
	if err != nil {
		return fmt.Errorf("failed to write zip file to (%s): %s", outputFile, err)
	}

	return nil
}

func resourceFileDownloadRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceFileDownloadUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceFileDownloadDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to download HSS export file. Deleting this resource
    will not clear the corresponding downloaded record, but will only remove the resource information from the
    tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
