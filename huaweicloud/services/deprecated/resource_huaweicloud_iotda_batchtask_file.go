package deprecated

import (
	"context"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iotda"
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := iotda.WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	file, err := os.Open(d.Get("content").(string))
	if err != nil {
		return diag.Errorf("error opening batch task file: %s", err)
	}
	defer file.Close()

	body := model.UploadBatchTaskFileRequestBody{
		File: &def.FilePart{
			Content: file,
		},
	}
	createOpt := model.UploadBatchTaskFileRequest{
		Body: &body,
	}
	resp, err := client.UploadBatchTaskFile(&createOpt)
	if err != nil {
		return diag.Errorf("error uploading IoTDA batch task file: %s", err)
	}

	if resp == nil || resp.FileId == nil {
		return diag.Errorf("error uploading IoTDA batch task file: ID is not found in API response")
	}

	d.SetId(*resp.FileId)

	return resourceBatchTaskFileRead(ctx, d, meta)
}

func resourceBatchTaskFileRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := iotda.WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	resp, err := client.ListBatchTaskFiles(&model.ListBatchTaskFilesRequest{})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying IoTDA batch task files")
	}

	allFiles := *resp.Files
	targetFile := new(model.BatchTaskFile)
	for i := range allFiles {
		if *allFiles[i].FileId == d.Id() {
			targetFile = &allFiles[i]
			break
		}
	}

	if targetFile.FileId == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "IoTDA batch task file")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", targetFile.FileName),
		d.Set("created_at", targetFile.UploadTime),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBatchTaskFileDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := iotda.WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	deleteOpt := model.DeleteBatchTaskFileRequest{
		FileId: d.Id(),
	}
	_, err = client.DeleteBatchTaskFile(&deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting IoTDA batch task file: %s", err)
	}

	return nil
}
