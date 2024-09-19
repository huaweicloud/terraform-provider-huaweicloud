package ims

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IMS POST /v1/cloudimages/{image_id}/file
// @API IMS GET /v1/{project_id}/jobs/{job_id}
// ResourceImageExport is a definition of the one-time action resource that used to manage image export.
func ResourceImageExport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImageExportCreate,
		ReadContext:   resourceImageExportRead,
		DeleteContext: resourceImageExportDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bucket_url": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"file_format": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_quick_export": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func buildCreateImageExportBodyParam(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"bucket_url":      d.Get("bucket_url"),
		"file_format":     utils.ValueIgnoreEmpty(d.Get("file_format")),
		"is_quick_export": utils.ValueIgnoreEmpty(d.Get("is_quick_export")),
	}
}

func resourceImageExportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		imageId = d.Get("image_id").(string)
	)

	client, err := cfg.ImageV1Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS v1 client: %s", err)
	}

	createPath := client.Endpoint + "v1/cloudimages/{image_id}/file"
	createPath = strings.ReplaceAll(createPath, "{image_id}", imageId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildCreateImageExportBodyParam(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error exporting IMS image, %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error exporting IMS image: job_id is not found in API response")
	}

	err = cloudimages.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutCreate)/time.Second), jobId)
	if err != nil {
		return diag.Errorf("error waiting for IMS image export to complete: %s", err)
	}

	d.SetId(imageId)

	return resourceImageExportRead(ctx, d, meta)
}

func resourceImageExportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceImageExportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a one-time action resource.
	return nil
}
