package ims

import (
	"context"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var iamgeValidSortKeys = []string{
	"name", "container_format", "disk_format", "status", "id", "size",
}
var imageValidVisibilities = []string{
	"public", "private", "community", "shared",
}

func DataSourceImagesImageV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImagesImageV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validation.StringIsValidRegExp,
			},

			"visibility": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(imageValidVisibilities, false),
			},

			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"sort_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "name",
				ValidateFunc: validation.StringInSlice(iamgeValidSortKeys, false),
			},

			"sort_direction": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "asc",
				ValidateFunc: validation.StringInSlice([]string{
					"asc", "desc",
				}, false),
			},

			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"x86", "arm",
				}, false),
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

			// Deprecated values
			"size_min": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "size_min is deprecated",
			},
			"size_max": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "size_max is deprecated",
			},

			// Computed values
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
			"protected": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

// dataSourceImagesImageV2Read performs the image lookup.
func dataSourceImagesImageV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	imageClient, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud image client: %s", err)
	}

	visibility := d.Get("visibility").(string)
	if visibility == "public" {
		visibility = "gold"
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
		Imagetype:      visibility,
		Status:         "active",
	}

	if epsId := common.GetEnterpriseProjectID(d, config); epsId != "" {
		listOpts.EnterpriseProjectID = epsId
	} else {
		listOpts.EnterpriseProjectID = "all_granted_eps"
	}

	logp.Printf("[DEBUG] List Options: %#v", listOpts)

	var image cloudimages.Image
	allPages, err := cloudimages.List(imageClient, listOpts).AllPages()
	if err != nil {
		return fmtp.DiagErrorf("Unable to query images: %s", err)
	}

	allImages, err := cloudimages.ExtractImages(allPages)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve images: %s", err)
	}

	if len(allImages) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	// filter images by name_regex
	var filteredImages []cloudimages.Image
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(nameRegex.(string))
		if err != nil {
			return fmtp.DiagErrorf("name_regex format error: %s", err)
		}

		for _, image := range allImages {
			if r.MatchString(image.Name) {
				filteredImages = append(filteredImages, image)
			}
		}
	} else {
		filteredImages = allImages[:]
	}

	if len(filteredImages) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(filteredImages) > 1 {
		recent := d.Get("most_recent").(bool)
		logp.Printf("[DEBUG] Multiple images %d found and `most_recent` is set to: %t", len(filteredImages), recent)
		if recent {
			image = mostRecentImage(filteredImages)
		} else {
			return fmtp.DiagErrorf("Your query returned more than one result. Please try a more " +
				"specific search criteria, or set `most_recent` attribute to true.")
		}
	} else {
		image = filteredImages[0]
	}

	logp.Printf("[DEBUG] Single Image found: %s", image.ID)
	return dataSourceImagesImageV2Attributes(ctx, d, &image)
}

// dataSourceImagesImageV2Attributes populates the fields of an Image resource.
func dataSourceImagesImageV2Attributes(_ context.Context, d *schema.ResourceData, image *cloudimages.Image) diag.Diagnostics {
	logp.Printf("[DEBUG] huaweicloud_images_image details: %#v", image)
	d.SetId(image.ID)

	mErr := multierror.Append(
		d.Set("name", image.Name),
		d.Set("container_format", image.ContainerFormat),
		d.Set("disk_format", image.DiskFormat),
		d.Set("min_disk_gb", image.MinDisk),
		d.Set("min_ram_mb", image.MinRam),
		d.Set("owner", image.Owner),
		d.Set("protected", image.Protected),
		d.Set("visibility", image.Visibility),
		d.Set("image_type", image.VirtualEnvType),
		d.Set("os", image.Platform),
		d.Set("os_version", image.OsVersion),
		d.Set("checksum", image.Checksum),
		d.Set("file", image.File),
		d.Set("schema", image.Schema),
		d.Set("enterprise_project_id", image.EnterpriseProjectID),
		d.Set("status", image.Status),
		d.Set("created_at", image.CreatedAt.Format(time.RFC3339)),
		d.Set("updated_at", image.UpdatedAt.Format(time.RFC3339)),
	)

	if size, err := strconv.Atoi(image.ImageSize); err == nil {
		mErr = multierror.Append(mErr, d.Set("size_bytes", size))
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting IMS image fields: %s", err)
	}

	return nil
}

type imageSort []cloudimages.Image

func (a imageSort) Len() int      { return len(a) }
func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageSort) Less(i, j int) bool {
	itime := a[i].UpdatedAt
	jtime := a[j].UpdatedAt
	return itime.Unix() < jtime.Unix()
}

// Returns the most recent Image out of a slice of images.
func mostRecentImage(images []cloudimages.Image) cloudimages.Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}
