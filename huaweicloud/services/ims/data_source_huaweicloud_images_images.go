package ims

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// @API IMS GET /v2/cloudimages
func DataSourceImagesImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImagesImagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validation.StringIsValidRegExp,
			},
			"image_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_whole_image": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "name",
			},
			"sort_direction": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "asc",
			},
			"os": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "schema: Computed",
			},
			"images": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     ImagesImageRefSchema(),
			},
		},
	}
}

func ImagesImageRefSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protected": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"container_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_ram_mb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_ram_mb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"min_disk_gb": {
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
			"backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size_bytes": {
				Type:     schema.TypeInt,
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

func buildImageTypeParam(d *schema.ResourceData) string {
	imageType := d.Get("visibility").(string)
	if imageType == "public" {
		imageType = "gold"
	}

	return imageType
}

// dataSourceImagesImagesRead performs the image lookup.
func dataSourceImagesImagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	imageClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS v2 client: %s", err)
	}

	listOpts := cloudimages.ListOpts{
		ID:                  d.Get("image_id").(string),
		Name:                d.Get("name").(string),
		VirtualEnvType:      d.Get("image_type").(string),
		WholeImage:          d.Get("is_whole_image").(bool),
		Imagetype:           buildImageTypeParam(d),
		Owner:               d.Get("owner").(string),
		FlavorId:            d.Get("flavor_id").(string),
		SortKey:             d.Get("sort_key").(string),
		SortDir:             d.Get("sort_direction").(string),
		Platform:            d.Get("os").(string),
		Architecture:        d.Get("architecture").(string),
		Tag:                 d.Get("tag").(string),
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d, "all_granted_eps"),
		// Default query status for **active** images.
		Status: "active",
	}

	allPages, err := cloudimages.List(imageClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("unable to query images: %s", err)
	}

	allImages, err := cloudimages.ExtractImages(allPages)
	if err != nil {
		return diag.Errorf("unable to extract images: %s", err)
	}

	// Filter images by `name_regex`.
	var nameRegexRes *regexp.Regexp
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		nameRegexRes, err = regexp.Compile(nameRegex.(string))
		if err != nil {
			return diag.Errorf("name_regex format error: %s", err)
		}
	}
	var resultImages []interface{}
	var ids []string
	for _, item := range allImages {
		image := item
		if nameRegexRes != nil && !nameRegexRes.MatchString(image.Name) {
			continue
		}
		ids = append(ids, image.ID)
		resultImages = append(resultImages, flattenImage(&image))
	}

	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("images", resultImages),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenImage(image *cloudimages.Image) map[string]interface{} {
	res := map[string]interface{}{
		"id":                    image.ID,
		"name":                  image.Name,
		"image_type":            image.VirtualEnvType,
		"visibility":            flattenVisibility(image.Imagetype),
		"owner":                 image.Owner,
		"os":                    image.Platform,
		"architecture":          flattenArchitecture(image.SupportArm),
		"enterprise_project_id": image.EnterpriseProjectID,
		"os_version":            image.OsVersion,
		"file":                  image.File,
		"schema":                image.Schema,
		"status":                image.Status,
		"description":           image.Description,
		"protected":             image.Protected,
		"container_format":      image.ContainerFormat,
		"min_ram_mb":            image.MinRam,
		"max_ram_mb":            flattenMaxRAM(image.MaxRam),
		"min_disk_gb":           image.MinDisk,
		"disk_format":           image.DiskFormat,
		"data_origin":           image.DataOrigin,
		"backup_id":             image.BackupID,
		"active_at":             image.ActiveAt,
		"created_at":            image.CreatedAt.Format(time.RFC3339),
		"updated_at":            image.UpdatedAt.Format(time.RFC3339),
	}

	if size, err := strconv.Atoi(image.ImageSize); err == nil {
		res["size_bytes"] = size
	}

	return res
}

func flattenVisibility(imageType string) string {
	if imageType == "gold" {
		imageType = "public"
	}

	return imageType
}
