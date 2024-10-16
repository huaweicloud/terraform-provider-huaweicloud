// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IMS
// ---------------------------------------------------------------

package ims

import (
	"context"
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

// @API IMS DELETE /v1/cloudimages/members
// @API IMS POST /v1/cloudimages/members
// @API IMS GET /v1/{project_id}/jobs/{job_id}
func ResourceImsImageShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImsImageShareCreate,
		UpdateContext: resourceImsImageShareUpdate,
		ReadContext:   resourceImsImageShareRead,
		DeleteContext: resourceImsImageShareDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source_image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the source image.`,
			},
			"target_project_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: `Specifies the IDs of the target projects.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceImsImageShareCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		projectIds    = d.Get("target_project_ids")
		sourceImageId = d.Get("source_image_id").(string)
	)
	err := dealImageMembers(ctx, d, cfg, "POST", sourceImageId, projectIds.(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sourceImageId)

	return resourceImsImageShareRead(ctx, d, meta)
}

func resourceImsImageShareUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	if d.HasChange("target_project_ids") {
		oProjectIdsRaw, nProjectIdsRaw := d.GetChange("target_project_ids")
		shareProjectIds := nProjectIdsRaw.(*schema.Set).Difference(oProjectIdsRaw.(*schema.Set))
		unShareProjectIds := oProjectIdsRaw.(*schema.Set).Difference(nProjectIdsRaw.(*schema.Set))
		if shareProjectIds.Len() > 0 {
			err := dealImageMembers(ctx, d, cfg, "POST", d.Id(), shareProjectIds.List())
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if unShareProjectIds.Len() > 0 {
			err := dealImageMembers(ctx, d, cfg, "DELETE", d.Id(), unShareProjectIds.List())
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceImsImageShareRead(ctx, d, meta)
}

func resourceImsImageShareRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceImsImageShareDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	projectIds := d.Get("target_project_ids")
	err := dealImageMembers(ctx, d, cfg, "DELETE", d.Id(), projectIds.(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func dealImageMembers(ctx context.Context, d *schema.ResourceData, cfg *config.Config, requestMethod,
	imageId string, projectIds []interface{}) error {
	var (
		region             = cfg.GetRegion(d)
		imageMemberHttpUrl = "v1/cloudimages/members"
		imageMemberProduct = "ims"
	)

	imageMemberClient, err := cfg.NewServiceClient(imageMemberProduct, region)
	if err != nil {
		return fmt.Errorf("error creating IMS Client: %s", err)
	}

	imageMemberPath := imageMemberClient.Endpoint + imageMemberHttpUrl
	imageMemberOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	imageMemberOpt.JSONBody = utils.RemoveNil(buildImageMemberBodyParams(imageId, projectIds))
	imageMemberResp, err := imageMemberClient.Request(requestMethod, imageMemberPath, &imageMemberOpt)

	var (
		operateMethod = "creating"
		timeout       = schema.TimeoutCreate
	)
	if requestMethod == "DELETE" {
		operateMethod = "deleting"
		timeout = schema.TimeoutDelete
	}

	if err != nil {
		return fmt.Errorf("error %s IMS image share: %s", operateMethod, err)
	}

	imageMemberRespBody, err := utils.FlattenResponse(imageMemberResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", imageMemberRespBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("unable to find the job ID of the IMS image share from the API response")
	}

	return waitForImageShareOrAcceptJobSuccess(ctx, d, imageMemberClient, jobId, timeout)
}

func buildImageMemberBodyParams(imageId string, projectIds []interface{}) map[string]interface{} {
	imagesParams := []interface{}{
		utils.ValueIgnoreEmpty(imageId),
	}
	bodyParams := map[string]interface{}{
		"images":   imagesParams,
		"projects": projectIds,
	}
	return bodyParams
}

func waitForImageShareOrAcceptJobSuccess(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	jobId, timeout string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INIT", "RUNNING"},
		Target:     []string{"SUCCESS"},
		Refresh:    imageShareOrAcceptJobStatusRefreshFunc(jobId, client),
		Timeout:    d.Timeout(timeout),
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for job (%s) success: %s", jobId, err)
	}

	return nil
}

func imageShareOrAcceptJobStatusRefreshFunc(jobId string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getJobStatusHttpUrl = "v1/{project_id}/jobs/{job_id}"
		)

		getJobStatusPath := client.Endpoint + getJobStatusHttpUrl
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{project_id}", client.ProjectID)
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{job_id}", fmt.Sprintf("%v", jobId))

		getJobStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		getJobStatusResp, err := client.Request("GET", getJobStatusPath, &getJobStatusOpt)
		if err != nil {
			return getJobStatusResp, "ERROR", nil
		}

		getJobStatusRespBody, err := utils.FlattenResponse(getJobStatusResp)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", getJobStatusRespBody, "")
		return getJobStatusRespBody, status.(string), nil
	}
}
