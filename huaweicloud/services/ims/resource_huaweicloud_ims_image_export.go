package ims

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

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
		product = "ims"
		imageId = d.Get("image_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	createPath := client.Endpoint + "v1/cloudimages/{image_id}/file"
	createPath = strings.ReplaceAll(createPath, "{image_id}", imageId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateImageExportBodyParam(d)),
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
		return diag.Errorf("error exporting IMS image: job ID is not found in API response")
	}

	err = waitForCreateImageExportJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for IMS image export to complete: %s", err)
	}

	d.SetId(imageId)

	return resourceImageExportRead(ctx, d, meta)
}

func waitForCreateImageExportJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      imageExportJobStatusRefreshFunc(jobId, client),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for IMS image export job (%s) to succeed: %s", jobId, err)
	}

	return nil
}

func imageExportJobStatusRefreshFunc(jobId string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getPath := client.Endpoint + "v1/{project_id}/jobs/{job_id}"
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{job_id}", fmt.Sprintf("%v", jobId))
		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return getResp, "ERROR", fmt.Errorf("error retrieving IMS image export job: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return getRespBody, "ERROR", err
		}

		status := utils.PathSearch("status", getRespBody, "").(string)
		if status == "SUCCESS" {
			return getRespBody, "COMPLETED", nil
		}

		if status == "FAIL" {
			return getRespBody, "COMPLETED", errors.New("the image export job failed")
		}

		if status == "" {
			return getRespBody, "ERROR", errors.New("status field is not found in API response")
		}

		return getRespBody, "PENDING", nil
	}
}

func resourceImageExportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceImageExportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a one-time action resource.
	return nil
}
