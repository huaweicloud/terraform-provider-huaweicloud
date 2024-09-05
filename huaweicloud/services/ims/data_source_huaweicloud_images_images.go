package ims

import (
	"context"
	"log"
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

			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"os_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
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
			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_disk_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"min_ram_mb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protected": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_id": {
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
			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

// dataSourceImagesImagesRead performs the image lookup.
func dataSourceImagesImagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	imageClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	imageType := d.Get("visibility").(string)
	if imageType == "public" {
		imageType = "gold"
	}

	listOpts := cloudimages.ListOpts{
		Name:           d.Get("name").(string),
		Owner:          d.Get("owner").(string),
		SortKey:        d.Get("sort_key").(string),
		SortDir:        d.Get("sort_direction").(string),
		Tag:            d.Get("tag").(string),
		Platform:       d.Get("os").(string),
		OsVersion:      d.Get("os_version").(string),
		Architecture:   d.Get("architecture").(string),
		VirtualEnvType: d.Get("image_type").(string),
		FlavorId:       d.Get("flavor_id").(string),
		Imagetype:      imageType,
		Status:         "active",
	}

	if epsId := cfg.GetEnterpriseProjectID(d); epsId != "" {
		listOpts.EnterpriseProjectID = epsId
	} else {
		listOpts.EnterpriseProjectID = "all_granted_eps"
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	allPages, err := cloudimages.List(imageClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("unable to query images: %s", err)
	}

	allImages, err := cloudimages.ExtractImages(allPages)
	if err != nil {
		return diag.Errorf("unable to retrieve images: %s", err)
	}

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
	mErr := multierror.Append(d.Set("images", resultImages))

	d.SetId(hashcode.Strings(ids))

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenImage(image *cloudimages.Image) map[string]interface{} {
	res := map[string]interface{}{
		"id":                    image.ID,
		"name":                  image.Name,
		"visibility":            image.Imagetype,
		"container_format":      image.ContainerFormat,
		"disk_format":           image.DiskFormat,
		"min_disk_gb":           image.MinDisk,
		"min_ram_mb":            image.MinRam,
		"owner":                 image.Owner,
		"protected":             image.Protected,
		"image_type":            image.VirtualEnvType,
		"os":                    image.Platform,
		"os_version":            image.OsVersion,
		"checksum":              image.Checksum,
		"enterprise_project_id": image.EnterpriseProjectID,
		"status":                image.Status,
		"backup_id":             image.BackupID,
		"created_at":            image.CreatedAt.Format(time.RFC3339),
		"updated_at":            image.UpdatedAt.Format(time.RFC3339),
	}

	if image.Imagetype == "gold" {
		res["visibility"] = "public"
	}
	if size, err := strconv.Atoi(image.ImageSize); err == nil {
		res["size_bytes"] = size
	}

	return res
}
