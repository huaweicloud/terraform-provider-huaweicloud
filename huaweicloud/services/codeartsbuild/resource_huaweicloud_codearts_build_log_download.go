package codeartsbuild

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var logDownloadNonUpdatableParams = []string{
	"record_id", "log_level", "log_file",
}

// @API CodeArtsBuild GET /v4/{record_id}/download-log
func ResourceCodeArtsBuildLogDownload() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBuildLogDownloadCreate,
		ReadContext:   resourceBuildLogDownloadRead,
		UpdateContext: resourceBuildLogDownloadUpdate,
		DeleteContext: resourceBuildLogDownloadDelete,

		CustomizeDiff: config.FlexibleForceNew(logDownloadNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"record_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the record ID.`,
			},
			"log_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the log level.`,
			},
			"log_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the log file path.`,
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

func resourceBuildLogDownloadCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("codearts_build", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CodeArts Build client: %s", err)
	}

	recordId := d.Get("record_id").(string)
	httpUrl := "v4/{record_id}/download-log"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{record_id}", recordId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	if v, ok := d.GetOk("log_level"); ok {
		getPath += fmt.Sprintf("?log_level=%v", v)
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error downloading CodeArts Build task build log: %s", err)
	}

	contentType := getResp.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flattening response: %s", err)
		}
		if err := checkResponseError(getRespBody); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(recordId)
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "no file was saved",
			},
		}
	} else if strings.HasPrefix(contentType, "application/octet-stream") {
		d.SetId(recordId)

		filePath := fmt.Sprintf("%s.txt", recordId)
		if v, ok := d.GetOk("log_file"); ok {
			filePath = v.(string)
		}

		if _, err := os.Stat(filePath); err != nil {
			if !os.IsNotExist(err) {
				return diag.Errorf("error occur while trying to access file %s: %s", filePath, err)
			}

			defer getResp.Body.Close()

			file, err := os.Create(filePath)
			if err != nil {
				return diag.Errorf("error creating file %s: %s", filePath, err)
			}

			if _, err = io.Copy(file, getResp.Body); err != nil {
				return diag.Errorf("error downloading build log to file %s: %s", filePath, err)
			}
			return nil
		}

		errorMsg := fmt.Sprintf("file (%s) exists", filePath)
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  errorMsg,
			},
		}
	}

	return diag.Errorf("unsupported content type: %s", contentType)
}

func resourceBuildLogDownloadRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBuildLogDownloadUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBuildLogDownloadDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting build log download resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
