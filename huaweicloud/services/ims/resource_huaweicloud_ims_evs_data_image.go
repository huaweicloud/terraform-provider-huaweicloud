package ims

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IMS POST /v2/cloudimages/action
// @API IMS GET /v1/{project_id}/jobs/{job_id}
// @API IMS GET /v2/cloudimages
// @API IMS GET /v2/{project_id}/images/{image_id}/tags
// @API IMS PATCH /v2/cloudimages/{image_id}
// @API IMS POST /v2/{project_id}/images/{image_id}/tags/action
// @API IMS DELETE /v2/images/{image_id}
func ResourceEvsDataImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvsDataImageCreate,
		ReadContext:   resourceEvsDataImageRead,
		UpdateContext: resourceEvsDataImageUpdate,
		DeleteContext: resourceImageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// The `description` field can be left blank, so the `Computed` attribute is not used.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Attributes
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_disk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEvsDataImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS v2 client: %s", err)
	}

	tags := buildCreateTagsParam(d)
	dataImageOpts := []cloudimages.DataImage{
		{
			Name:        d.Get("name").(string),
			VolumeId:    d.Get("volume_id").(string),
			Description: d.Get("description").(string),
			Tags:        tags,
		},
	}
	createOpts := &cloudimages.CreateDataImageByServerOpts{
		DataImages:          dataImageOpts,
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
	}

	createResp, err := cloudimages.CreateDataImageByServer(client, createOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error creating IMS EVS data image: %s", err)
	}

	imageId, err := waitForCreateImageCompleted(client, d, createResp.JobID)
	if err != nil {
		return diag.Errorf("error waiting for IMS EVS data image to complete: %s", err)
	}

	d.SetId(imageId)

	return resourceEvsDataImageRead(ctx, d, meta)
}

func buildCreateTagsParam(d *schema.ResourceData) []string {
	rawTags := d.Get("tags").(map[string]interface{})
	var tagStrings []string
	for key, val := range rawTags {
		tagStrings = append(tagStrings, fmt.Sprintf("%s.%s", key, val))
	}

	return tagStrings
}

func resourceEvsDataImageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		mErr   *multierror.Error
	)

	client, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS v2 client: %s", err)
	}

	imageList, err := GetImageList(client, d.Id())
	if err != nil {
		return diag.Errorf("error retrieving IMS EVS data images: %s", err)
	}

	// If the list API return empty, then process `CheckDeleted` logic.
	if len(imageList) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "IMS EVS data image")
	}

	image := imageList[0]
	imageTags := getImageTags(d, client)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("name", image.Name),
		d.Set("volume_id", getSpecificValueFormDataOrigin(image.DataOrigin, "volume")),
		d.Set("description", image.Description),
		d.Set("tags", imageTags),
		d.Set("enterprise_project_id", image.EnterpriseProjectID),
		d.Set("status", image.Status),
		d.Set("visibility", image.Visibility),
		d.Set("image_size", image.ImageSize),
		d.Set("min_disk", image.MinDisk),
		d.Set("os_type", image.OsType),
		d.Set("disk_format", image.DiskFormat),
		d.Set("data_origin", image.DataOrigin),
		d.Set("active_at", image.ActiveAt),
		d.Set("created_at", image.CreatedAt.Format(time.RFC3339)),
		d.Set("updated_at", image.UpdatedAt.Format(time.RFC3339)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEvsDataImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS v2 client: %s", err)
	}

	err = updateImage(ctx, cfg, client, d)
	if err != nil {
		return diag.Errorf("error updating IMS EVS data image: %s", err)
	}

	return resourceEvsDataImageRead(ctx, d, meta)
}
