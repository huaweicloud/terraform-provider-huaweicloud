// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IMS
// ---------------------------------------------------------------

package ims

import (
	"context"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IMS PUT /v1/cloudimages/members
// @API IMS GET /v1/{project_id}/jobs/{job_id}
func ResourceImsImageShareAccepter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImsImageShareAccepterCreate,
		ReadContext:   resourceImsImageShareAccepterRead,
		DeleteContext: resourceImsImageShareAccepterDelete,

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
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the image.`,
			},
			"vault_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a vault.`,
			},
		},
	}
}

func resourceImsImageShareAccepterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                             = meta.(*config.Config)
		region                          = cfg.GetRegion(d)
		createImageShareAccepterHttpUrl = "v1/cloudimages/members"
		createImageShareAccepterProduct = "ims"
	)
	createImageShareAccepterClient, err := cfg.NewServiceClient(createImageShareAccepterProduct, region)
	if err != nil {
		return diag.Errorf("error creating IMS Client: %s", err)
	}

	createImageShareAccepterPath := createImageShareAccepterClient.Endpoint + createImageShareAccepterHttpUrl
	createImageShareAccepterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createImageShareAccepterOpt.JSONBody = utils.RemoveNil(buildCreateImageShareAccepterBodyParams(d,
		createImageShareAccepterClient.ProjectID))

	createImageShareAccepterResp, err := createImageShareAccepterClient.Request("PUT",
		createImageShareAccepterPath, &createImageShareAccepterOpt)
	if err != nil {
		return diag.Errorf("error creating IMS image share accepter: %s", err)
	}

	createImageShareAccepterRespBody, err := utils.FlattenResponse(createImageShareAccepterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createImageShareAccepterRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating IMS image share accepter: job_id is not found in API response")
	}

	err = waitForImageShareOrAcceptJobSuccess(ctx, d, createImageShareAccepterClient, jobId, schema.TimeoutCreate)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)

	return resourceImsImageShareAccepterRead(ctx, d, meta)
}

func buildCreateImageShareAccepterBodyParams(d *schema.ResourceData, projectId string) map[string]interface{} {
	imagesParams := []interface{}{
		utils.ValueIgnoreEmpty(d.Get("image_id")),
	}
	bodyParams := map[string]interface{}{
		"images":     imagesParams,
		"project_id": projectId,
		"status":     "accepted",
		"vault_id":   utils.ValueIgnoreEmpty(d.Get("vault_id")),
	}
	return bodyParams
}

func resourceImsImageShareAccepterRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceImsImageShareAccepterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                             = meta.(*config.Config)
		region                          = cfg.GetRegion(d)
		deleteImageShareAccepterHttpUrl = "v1/cloudimages/members"
		deleteImageShareAccepterProduct = "ims"
	)
	deleteImageShareAccepterClient, err := cfg.NewServiceClient(deleteImageShareAccepterProduct, region)
	if err != nil {
		return diag.Errorf("error creating IMS Client: %s", err)
	}

	deleteImageShareAccepterPath := deleteImageShareAccepterClient.Endpoint + deleteImageShareAccepterHttpUrl
	deleteImageShareAccepterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	deleteImageShareAccepterOpt.JSONBody = utils.RemoveNil(buildDeleteImageShareAccepterBodyParams(d,
		deleteImageShareAccepterClient.ProjectID))
	deleteImageShareAccepterResp, err := deleteImageShareAccepterClient.Request("PUT",
		deleteImageShareAccepterPath, &deleteImageShareAccepterOpt)
	if err != nil {
		return diag.Errorf("error deleting IMS image share accepter: %s", err)
	}

	deleteImageShareAccepterRespBody, err := utils.FlattenResponse(deleteImageShareAccepterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteImageShareAccepterRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID of IMS image share accepter from the API response")
	}

	err = waitForImageShareOrAcceptJobSuccess(ctx, d, deleteImageShareAccepterClient, jobId, schema.TimeoutDelete)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildDeleteImageShareAccepterBodyParams(d *schema.ResourceData, projectId string) map[string]interface{} {
	imagesParams := []interface{}{
		utils.ValueIgnoreEmpty(d.Get("image_id")),
	}
	bodyParams := map[string]interface{}{
		"images":     imagesParams,
		"project_id": projectId,
		"status":     "rejected",
	}
	return bodyParams
}
