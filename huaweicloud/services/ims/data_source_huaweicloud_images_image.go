package ims

import (
	"context"
	"log"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IMS GET /v2/cloudimages
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

			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
		},
	}
}

// dataSourceImagesImageV2Read performs the image lookup.
func dataSourceImagesImageV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if len(allImages) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	// filter images by name_regex
	var filteredImages []cloudimages.Image
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(nameRegex.(string))
		if err != nil {
			return diag.Errorf("name_regex format error: %s", err)
		}

		for _, image := range allImages {
			if r.MatchString(image.Name) {
				filteredImages = append(filteredImages, image)
			}
		}
	} else {
		filteredImages = allImages
	}

	if len(filteredImages) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	var image cloudimages.Image
	if len(filteredImages) > 1 {
		recent := d.Get("most_recent").(bool)
		log.Printf("[DEBUG] Multiple images %d found and `most_recent` is set to: %t", len(filteredImages), recent)
		if !recent {
			return diag.Errorf("Your query returned more than one result. Please try a more " +
				"specific search criteria, or set `most_recent` attribute to true.")
		}
		image = mostRecentImage(filteredImages)
	} else {
		image = filteredImages[0]
	}

	log.Printf("[DEBUG] Single Image found: %s", image.ID)
	return dataSourceImagesImageV2Attributes(ctx, d, &image)
}

// dataSourceImagesImageV2Attributes populates the fields of an Image resource.
func dataSourceImagesImageV2Attributes(_ context.Context, d *schema.ResourceData, image *cloudimages.Image) diag.Diagnostics {
	log.Printf("[DEBUG] Get IMS image details: %#v", image)
	d.SetId(image.ID)

	imageType := image.Imagetype
	if imageType == "gold" {
		imageType = "public"
	}
	mErr := multierror.Append(
		d.Set("name", image.Name),
		d.Set("visibility", imageType),
		d.Set("container_format", image.ContainerFormat),
		d.Set("disk_format", image.DiskFormat),
		d.Set("min_disk_gb", image.MinDisk),
		d.Set("min_ram_mb", image.MinRam),
		d.Set("owner", image.Owner),
		d.Set("protected", image.Protected),
		d.Set("image_type", image.VirtualEnvType),
		d.Set("os", image.Platform),
		d.Set("os_version", image.OsVersion),
		d.Set("checksum", image.Checksum),
		d.Set("file", image.File),
		d.Set("schema", image.Schema),
		d.Set("enterprise_project_id", image.EnterpriseProjectID),
		d.Set("status", image.Status),
		d.Set("backup_id", image.BackupID),
		d.Set("created_at", image.CreatedAt.Format(time.RFC3339)),
		d.Set("updated_at", image.UpdatedAt.Format(time.RFC3339)),
	)

	if size, err := strconv.Atoi(image.ImageSize); err == nil {
		mErr = multierror.Append(mErr, d.Set("size_bytes", size))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

type imageSort []cloudimages.Image

func (a imageSort) Len() int { return len(a) }

func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a imageSort) Less(i, j int) bool {
	iTime := a[i].UpdatedAt
	jTime := a[j].UpdatedAt
	return iTime.Unix() < jTime.Unix()
}

// Returns the most recent Image out of a slice of images.
func mostRecentImage(images []cloudimages.Image) cloudimages.Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}
