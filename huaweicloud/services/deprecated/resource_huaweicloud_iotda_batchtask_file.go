package deprecated

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iotda"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA POST /v5/iot/{project_id}/batchtask-files
// @API IoTDA GET /v5/iot/{project_id}/batchtask-files
// @API IoTDA DELETE /v5/iot/{project_id}/batchtask-files/{file_id}
func ResourceBatchTaskFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchTaskFileCreate,
		ReadContext:   resourceBatchTaskFileRead,
		DeleteContext: resourceBatchTaskFileDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		DeprecationMessage: "batchtask file has been deprecated.",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBatchTaskFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/batchtask-files"
		product = "iotda"
	)

	isDerived := iotda.WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	file, err := os.Open(d.Get("content").(string))
	if err != nil {
		return diag.Errorf("error opening batch task file: %s", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)

	formFile, err := multiPartWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = io.Copy(formFile, file)
	if err != nil {
		return diag.FromErr(err)
	}

	err = multiPartWriter.Close()
	if err != nil {
		return diag.FromErr(err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":         multiPartWriter.FormDataContentType(),
			"X-Sdk-Content-Sha256": "UNSIGNED-PAYLOAD",
		},
		RawBody: &requestBody,
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error uploading IoTDA batch task file: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	fileID := utils.PathSearch("file_id", respBody, "").(string)
	if fileID == "" {
		return diag.Errorf("error uploading IoTDA batch task file: ID is not found in API response")
	}

	d.SetId(fileID)

	return resourceBatchTaskFileRead(ctx, d, meta)
}

func resourceBatchTaskFileRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/batchtask-files"
		product = "iotda"
	)

	isDerived := iotda.WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error querying IoTDA batch task files: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskFile := utils.PathSearch(fmt.Sprintf("files[?file_id == '%s']|[0]", d.Id()), respBody, nil)
	if taskFile == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("file_name", taskFile, nil)),
		d.Set("created_at", utils.PathSearch("upload_time", taskFile, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBatchTaskFileDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/batchtask-files/{file_id}"
		product = "iotda"
	)

	isDerived := iotda.WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{file_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		// When the resource does not exist, the API response status code is `404`.
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA batch task file")
	}

	return nil
}
