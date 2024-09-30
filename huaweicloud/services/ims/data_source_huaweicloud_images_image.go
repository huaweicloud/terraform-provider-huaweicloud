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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
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
			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"image_id": {
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
			"image_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_whole_image": {
				Type:     schema.TypeBool,
				Optional: true,
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
				Computed: true,
			},
			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"os_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "schema: Computed",
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
			"metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: utils.SchemaDesc("metadata is deprecated", utils.SchemaDescInput{Internal: true}),
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
			"checksum": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: utils.SchemaDesc("checksum is deprecated", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

// dataSourceImagesImageV2Read performs the image lookup.
func dataSourceImagesImageV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if len(allImages) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	// Filter images by `name_regex`.
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

	// Filter images by `most_recent`.
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

	return dataSourceImagesImageV2Attributes(region, d, &image)
}

// dataSourceImagesImageV2Attributes populates the fields of an Image resource.
func dataSourceImagesImageV2Attributes(region string, d *schema.ResourceData, image *cloudimages.Image) diag.Diagnostics {
	d.SetId(image.ID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("image_id", image.ID),
		d.Set("name", image.Name),
		d.Set("image_type", image.VirtualEnvType),
		d.Set("visibility", flattenVisibility(image.Imagetype)),
		d.Set("owner", image.Owner),
		d.Set("os", image.Platform),
		d.Set("architecture", flattenArchitecture(image.SupportArm)),
		d.Set("enterprise_project_id", image.EnterpriseProjectID),
		d.Set("os_version", image.OsVersion),
		d.Set("file", image.File),
		d.Set("schema", image.Schema),
		d.Set("status", image.Status),
		d.Set("description", image.Description),
		d.Set("protected", image.Protected),
		d.Set("container_format", image.ContainerFormat),
		d.Set("min_ram_mb", image.MinRam),
		d.Set("max_ram_mb", flattenMaxRAM(image.MaxRam)),
		d.Set("min_disk_gb", image.MinDisk),
		d.Set("min_ram_mb", image.MinRam),
		d.Set("disk_format", image.DiskFormat),
		d.Set("data_origin", image.DataOrigin),
		d.Set("backup_id", image.BackupID),
		d.Set("active_at", image.ActiveAt),
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
