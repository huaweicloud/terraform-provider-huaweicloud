package workspace

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v1/{project_id}/image-servers/{server_id}/actions/recreate-image
// @API Workspace GET /v1/{project_id}/image-server-jobs/{job_id}
// @API IMS GET /v2/cloudimages
// @API IMS DELETE /v2/images/{image_id}
func ResourceAppImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppImageCreate,
		ReadContext:   resourceAppImageRead,
		DeleteContext: resourceAppImageDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The image server ID for generating a private image.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the image.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The description of the image.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The enterprise project ID to which the image belongs.",
			},
		},
	}
}

func buildAppImageCreateOpts(name, description interface{}, epsId string) map[string]interface{} {
	return map[string]interface{}{
		"name":                  name,
		"description":           utils.ValueIgnoreEmpty(description),
		"enterprise_project_id": utils.ValueIgnoreEmpty(epsId),
	}
}

func resourceAppImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v1/{project_id}/image-servers/{server_id}/actions/recreate-image"
		imageServerId = d.Get("server_id").(string)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{server_id}", imageServerId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAppImageCreateOpts(d.Get("name"), d.Get("description"), cfg.GetEnterpriseProjectID(d))),
	}
	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating private image from image server (%s): %s", imageServerId, err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The job ID consists of 18 digits, e.g. `773698274006671360`.
	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find job ID from API response")
	}
	// Set the resource ID in advance to avoid dirty data when the image task fails.
	d.SetId(jobId)

	serverResp, err := waitForImageServerJobCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), jobId)
	if err != nil {
		return diag.Errorf("error waiting for creating image job (%s) completed: %s", jobId, err)
	}

	imageId := utils.PathSearch("sub_jobs|[0].job_resource_info.resource_id", serverResp, "").(string)
	if imageId == "" {
		return diag.Errorf("unable to find generated image ID from API response")
	}

	d.SetId(imageId)

	return resourceAppImageRead(ctx, d, meta)
}

func resourceAppImageRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/images/{image_id}"
		imageId = d.Id()
	)

	imsClient, err := cfg.NewServiceClient("ims", region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}
	deletePath := imsClient.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", imsClient.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{image_id}", imageId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// If the image creation fails, the resource ID consists of 18 digits, e.g. `773698274006671360`.
	// At this time, the deletion interface response status code is still 204.
	_, err = imsClient.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting generated image (%s)", imageId))
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      refreshImageStatusFunc(imsClient, imageId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 30 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for deleting generated image (%s) completed: %s", imageId, err)
	}

	return nil
}

func refreshImageStatusFunc(client *golangsdk.ServiceClient, imageId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		httpUrl := "v2/cloudimages?id={image_id}"
		listPath := client.Endpoint + httpUrl
		listPath = strings.ReplaceAll(listPath, "{image_id}", imageId)
		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		resp, err := client.Request("GET", listPath, &getOpt)
		if err != nil {
			return resp, "ERROR", err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return resp, "ERROR", err
		}

		images := utils.PathSearch("images", respBody, make([]interface{}, 0)).([]interface{})
		if len(images) < 1 {
			return "", "DELETED", nil
		}
		return images, "PENDING", nil
	}
}
