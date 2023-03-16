// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IMS
// ---------------------------------------------------------------

package ims

import (
	"context"
	"fmt"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceImsImageCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImsImageCopyCreate,
		UpdateContext: resourceImsImageCopyUpdate,
		ReadContext:   resourceImsImageCopyRead,
		DeleteContext: resourceImsImageCopyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
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
				Description: `Special the ID of the source image.`,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  `Specifies the name of the copy image.`,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"target_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Special the target region name.`,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  `Specifies the description of the copy image.`,
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the master key used for encrypting an image.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the enterprise project id of the image.`,
			},
			"agency_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the agency name.`,
			},
			"vault_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the vault.`,
			},
			"tags": common.TagsSchema(),
		},
	}
}

func resourceImsImageCopyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createImageCopy: create IMS image copy
	var (
		createImageCopyHttpUrl            = "v1/cloudimages/{image_id}/copy"
		createImageCrossRegionCopyHttpUrl = "v1/cloudimages/{image_id}/cross_region_copy"
		createImageCopyProduct            = "ims"
	)

	createImageCopyClient, err := cfg.NewServiceClient(createImageCopyProduct, region)
	if err != nil {
		return diag.Errorf("error creating IMS Client: %s", err)
	}

	createImageCopyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	var createImageCopyPath string
	targetRegion := d.Get("target_region").(string)
	if targetRegion == "" {
		createImageCopyPath = createImageCopyClient.Endpoint + createImageCopyHttpUrl
		createImageCopyOpt.JSONBody = utils.RemoveNil(buildCreateImageCopyBodyParams(d, cfg))
	} else {
		createImageCopyPath = createImageCopyClient.Endpoint + createImageCrossRegionCopyHttpUrl
		createImageCopyOpt.JSONBody = utils.RemoveNil(buildCreateImageCrossRegionCopyBodyParams(d, cfg))
	}

	createImageCopyPath = strings.ReplaceAll(createImageCopyPath, "{image_id}", fmt.Sprintf("%v", d.Get("image_id")))

	createImageCopyResp, err := createImageCopyClient.Request("POST", createImageCopyPath, &createImageCopyOpt)
	if err != nil {
		return diag.Errorf("error creating ImsImageCopy: %s", err)
	}

	createImageCopyRespBody, err := utils.FlattenResponse(createImageCopyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createImageCopyRespBody, nil)
	if err != nil {
		return diag.Errorf("error creating ImsImageCopy: ID is not found in API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INIT", "RUNNING"},
		Target:     []string{"SUCCESS"},
		Refresh:    imageCopyStatusRefreshFunc(createImageCopyClient, cfg.RegionProjectIDMap[region], jobId.(string)),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      60 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	imageId, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for image copy job (%s) to be success: %s", jobId, err)
	}

	d.SetId(imageId.(string))

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		// if the image is copied to new region, the new region client is needed when set tags
		if targetRegion != "" {
			createImageCopyClient, err = cfg.NewServiceClient(createImageCopyProduct, targetRegion)
			if err != nil {
				return diag.Errorf("error creating IMS Client: %s", err)
			}
		}
		tagList := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(createImageCopyClient, "images", d.Id(), tagList).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of images %s: %s", d.Id(), tagErr)
		}
	}

	return resourceImsImageCopyRead(ctx, d, meta)
}

func buildCreateImageCopyBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  utils.ValueIngoreEmpty(d.Get("name")),
		"description":           utils.ValueIngoreEmpty(d.Get("description")),
		"cmk_id":                utils.ValueIngoreEmpty(d.Get("kms_key_id")),
		"enterprise_project_id": utils.ValueIngoreEmpty(common.GetEnterpriseProjectID(d, cfg)),
	}
	return bodyParams
}

func buildCreateImageCrossRegionCopyBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  utils.ValueIngoreEmpty(d.Get("name")),
		"description":           utils.ValueIngoreEmpty(d.Get("description")),
		"region":                utils.ValueIngoreEmpty(d.Get("target_region")),
		"project_name":          utils.ValueIngoreEmpty(d.Get("target_region")),
		"agency_name":           utils.ValueIngoreEmpty(d.Get("agency_name")),
		"enterprise_project_id": utils.ValueIngoreEmpty(common.GetEnterpriseProjectID(d, cfg)),
	}
	return bodyParams
}

func imageCopyStatusRefreshFunc(client *golangsdk.ServiceClient, projectId, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getJobStatusHttpUrl = "v1/{project_id}/jobs/{job_id}"
		)

		getJobStatusPath := client.Endpoint + getJobStatusHttpUrl
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{project_id}", fmt.Sprintf("%v", projectId))
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{job_id}", fmt.Sprintf("%v", jobId))

		getJobStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		getJobStatusResp, err := client.Request("GET", getJobStatusPath, &getJobStatusOpt)
		if err != nil {
			return nil, "", err
		}

		getJobStatusRespBody, err := utils.FlattenResponse(getJobStatusResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("status", getJobStatusRespBody, nil)
		if err != nil {
			return nil, "", err
		}
		if status == "FAIL" {
			return nil, "", fmt.Errorf("creating ImsImageCopy job run fail")
		}
		entities := utils.PathSearch("entities", getJobStatusRespBody, nil)
		if err != nil {
			return nil, "", err
		}
		imageId := utils.PathSearch("image_id", entities, nil)
		return imageId, status.(string), nil
	}
}

func resourceImsImageCopyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateImageCopy: update IMS image copy
	var (
		updateImageCopyHttpUrl = "v2/cloudimages/{image_id}"
		updateImageCopyProduct = "ims"
	)

	var updateImageCopyClient *golangsdk.ServiceClient
	var err error
	targetRegion := d.Get("target_region").(string)
	if targetRegion == "" {
		updateImageCopyClient, err = cfg.NewServiceClient(updateImageCopyProduct, region)
		if err != nil {
			return diag.Errorf("error creating ImsImageCopy Client: %s", err)
		}
	} else {
		updateImageCopyClient, err = cfg.NewServiceClient(updateImageCopyProduct, targetRegion)
		if err != nil {
			return diag.Errorf("error creating ImsImageCopy Client: %s", err)
		}
	}

	if d.HasChanges("name") {

		updateImageCopyPath := updateImageCopyClient.Endpoint + updateImageCopyHttpUrl
		updateImageCopyPath = strings.ReplaceAll(updateImageCopyPath, "{image_id}", fmt.Sprintf("%v", d.Id()))

		updateImageCopyOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateImageCopyOpt.JSONBody = []interface{}{utils.RemoveNil(buildUpdateImageCopyBodyParams(d))}
		_, err = updateImageCopyClient.Request("PATCH", updateImageCopyPath, &updateImageCopyOpt)
		if err != nil {
			return diag.Errorf("error updating ImsImageCopy: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(updateImageCopyClient, d, "images", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of IMS images :%s, err:%s", d.Id(), tagErr)
		}
	}

	return resourceImsImageCopyRead(ctx, d, meta)
}

func buildUpdateImageCopyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"op":    "replace",
		"path":  "/name",
		"value": utils.ValueIngoreEmpty(d.Get("name")),
	}
	return bodyParams
}

func resourceImsImageCopyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getImageCopy: query IMS image copy
	var (
		getImageCopyHttpUrl = "v2/cloudimages"
		getImageCopyProduct = "ims"
	)

	var getImageCopyClient *golangsdk.ServiceClient
	var err error
	targetRegion := d.Get("target_region").(string)
	if targetRegion == "" {
		getImageCopyClient, err = cfg.NewServiceClient(getImageCopyProduct, region)
		if err != nil {
			return diag.Errorf("error creating ImsImageCopy Client: %s", err)
		}
	} else {
		getImageCopyClient, err = cfg.NewServiceClient(getImageCopyProduct, targetRegion)
		if err != nil {
			return diag.Errorf("error creating ImsImageCopy Client: %s", err)
		}
	}

	getImageCopyPath := getImageCopyClient.Endpoint + getImageCopyHttpUrl

	getImageCopyQueryParams := buildGetImageCopyQueryParams(d.Id())
	getImageCopyPath += getImageCopyQueryParams

	getImageCopyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getImageCopyResp, err := getImageCopyClient.Request("GET", getImageCopyPath, &getImageCopyOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ImsImageCopy")
	}

	getImageCopyRespBody, err := utils.FlattenResponse(getImageCopyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	images := utils.PathSearch("images", getImageCopyRespBody, nil).([]interface{})
	if len(images) == 0 {
		return diag.Errorf("copy image is not exists: %s", d.Id())
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", images[0], nil)),
		d.Set("description", utils.PathSearch("__description", images[0], nil)),
		d.Set("kms_key_id", utils.PathSearch("cmk_id", images[0], nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", images[0], nil)),
	)

	// fetch tags
	if resourceTags, err := tags.Get(getImageCopyClient, "image", d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	} else {
		log.Printf("[WARN] Fetching tags of IMS images failed: %s", err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetImageCopyQueryParams(id string) string {
	res := ""
	res = fmt.Sprintf("%s&id=%v", res, id)

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func resourceImsImageCopyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteImageCopy: delete IMS image copy
	var (
		deleteImageCopyHttpUrl = "v2/images/{image_id}"
		deleteImageCopyProduct = "ims"
	)

	var deleteImageCopyClient *golangsdk.ServiceClient
	var err error
	targetRegion := d.Get("target_region").(string)
	if targetRegion == "" {
		deleteImageCopyClient, err = cfg.NewServiceClient(deleteImageCopyProduct, region)
		if err != nil {
			return diag.Errorf("error creating ImsImageCopy Client: %s", err)
		}
	} else {
		deleteImageCopyClient, err = cfg.NewServiceClient(deleteImageCopyProduct, targetRegion)
		if err != nil {
			return diag.Errorf("error creating ImsImageCopy Client: %s", err)
		}
	}

	deleteImageCopyPath := deleteImageCopyClient.Endpoint + deleteImageCopyHttpUrl
	deleteImageCopyPath = strings.ReplaceAll(deleteImageCopyPath, "{image_id}", fmt.Sprintf("%v", d.Id()))

	deleteImageCopyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteImageCopyClient.Request("DELETE", deleteImageCopyPath, &deleteImageCopyOpt)
	if err != nil {
		return diag.Errorf("error deleting ImsImageCopy: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    imageDeleteRefreshFunc(deleteImageCopyClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for delete image (%s) complete: %s", d.Id(), err)
	}

	return nil
}

func imageDeleteRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getImageCopyHttpUrl = "v2/cloudimages"
		)

		getImageCopyPath := client.Endpoint + getImageCopyHttpUrl

		getImageCopyQueryParams := buildGetImageCopyQueryParams(id)
		getImageCopyPath += getImageCopyQueryParams

		getImageCopyOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}

		v, err := client.Request("GET", getImageCopyPath, &getImageCopyOpt)
		if err != nil {
			return nil, "", err
		}
		getImageCopyRespBody, err := utils.FlattenResponse(v)
		if err != nil {
			return nil, "", err
		}
		images := utils.PathSearch("images", getImageCopyRespBody, nil).([]interface{})
		if len(images) == 0 {
			return v, "DELETED", nil
		}
		return v, "ACTIVE", nil
	}
}
