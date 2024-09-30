package ims

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cbr/v3/backups"
	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IMS POST /v1/cloudimages/wholeimages/action
// @API IMS GET /v1/{project_id}/jobs/{job_id}
// @API IMS GET /v2/cloudimages
// @API IMS GET /v2/{project_id}/images/{image_id}/tags
// @API CBR GET /v3/{project_id}/backups/{backup_id}
// @API IMS PATCH /v2/cloudimages/{image_id}
// @API IMS POST /v2/{project_id}/images/{image_id}/tags/action
// @API IMS DELETE /v2/images/{image_id}
func ResourceEcsWholeImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEcsWholeImageCreate,
		ReadContext:   resourceEcsWholeImageRead,
		UpdateContext: resourceEcsWholeImageUpdate,
		DeleteContext: resourceWholeImageDelete,

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
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vault_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// The `description` field can be left blank, so the `Computed` attribute is not used.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_ram": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"min_ram": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"is_delete_backup": {
				Type:     schema.TypeBool,
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
			"backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_disk": {
				Type:     schema.TypeInt,
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
			"os_version": {
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

func resourceEcsWholeImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		createResp *cloudimages.JobResponse
	)

	client, err := cfg.ImageV1Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS v1 client: %s", err)
	}

	imageTags := buildCreateImageTagsParam(d)
	createOpts := &cloudimages.CreateWholeImageOpts{
		Name:                d.Get("name").(string),
		InstanceId:          d.Get("instance_id").(string),
		VaultId:             d.Get("vault_id").(string),
		Description:         d.Get("description").(string),
		MaxRam:              d.Get("max_ram").(int),
		MinRam:              d.Get("min_ram").(int),
		ImageTags:           imageTags,
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
	}

	createResp, err = cloudimages.CreateWholeImageByServer(client, createOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error creating IMS ECS whole image: %s", err)
	}

	imageId, err := waitForCreateImageCompleted(client, d, createResp.JobID)
	if err != nil {
		return diag.Errorf("error waiting for IMS ECS whole image to complete: %s", err)
	}

	d.SetId(imageId)

	return resourceEcsWholeImageRead(ctx, d, meta)
}

func resourceEcsWholeImageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("error retrieving IMS ECS whole images: %s", err)
	}

	// If the list API return empty, then process `CheckDeleted` logic.
	if len(imageList) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "IMS ECS whole image")
	}

	image := imageList[0]
	imageTags := flattenImageTags(d, client)
	result, err := getBackupDetail(cfg, region, image.BackupID)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("name", image.Name),
		d.Set("instance_id", result[0]),
		d.Set("vault_id", result[1]),
		d.Set("description", image.Description),
		d.Set("max_ram", flattenMaxRAM(image.MaxRam)),
		d.Set("min_ram", image.MinRam),
		d.Set("tags", imageTags),
		d.Set("enterprise_project_id", image.EnterpriseProjectID),
		d.Set("status", image.Status),
		d.Set("visibility", image.Visibility),
		d.Set("backup_id", image.BackupID),
		d.Set("min_disk", image.MinDisk),
		d.Set("disk_format", image.DiskFormat),
		d.Set("data_origin", image.DataOrigin),
		d.Set("os_version", image.OsVersion),
		d.Set("active_at", image.ActiveAt),
		d.Set("created_at", image.CreatedAt.Format(time.RFC3339)),
		d.Set("updated_at", image.UpdatedAt.Format(time.RFC3339)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getBackupDetail(cfg *config.Config, region, backupId string) ([]string, error) {
	cbrClient, err := cfg.CbrV3Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating CBR v3 client: %s", err)
	}

	backup, err := backups.Get(cbrClient, backupId)
	if err != nil {
		return nil, fmt.Errorf("error querying backup detail: %s", err)
	}

	return []string{backup.ResourceId, backup.VaultId}, nil
}

func resourceEcsWholeImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("error updating IMS ECS whole image: %s", err)
	}

	return resourceEcsWholeImageRead(ctx, d, meta)
}

func resourceWholeImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		imageId = d.Id()
	)

	client, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS v2 client: %s", err)
	}

	// Before deleting, call the query API first, if the query result is empty, then process `CheckDeleted` logic.
	imageList, err := GetImageList(client, imageId)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(imageList) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "IMS whole image")
	}

	// For the whole image, need to use `delete_backup` to control whether to delete backup when deleting image.
	deletePath := client.Endpoint + "v2/images/{image_id}"
	deletePath = strings.ReplaceAll(deletePath, "{image_id}", imageId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"delete_backup": d.Get("is_delete_backup").(bool),
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting IMS whole image: %s", err)
	}

	// Because the delete API always return `204` status code,
	// so we need to call the list query API to check if the image has been successfully deleted.
	err = waitForDeleteImageCompleted(ctx, client, d)
	if err != nil {
		return diag.Errorf("error waiting for IMS whole image deleted: %s", err)
	}

	return nil
}
