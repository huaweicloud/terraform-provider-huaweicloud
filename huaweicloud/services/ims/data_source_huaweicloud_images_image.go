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

	"github.com/chnsz/golangsdk"

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

func dataSourceImagesImageV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "ims"
		httpUrl   = "v2/cloudimages"
		epsId     = cfg.GetEnterpriseProjectID(d, "all_granted_eps")
		allImages = make([]interface{}, 0)
		marker    = ""
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath += buildImagesQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		getPathWithMarker := buildImagesQueryParamsWithMarker(getPath, marker)
		getResp, err := client.Request("GET", getPathWithMarker, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving IMS images: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		imagesResp := utils.PathSearch("images", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(imagesResp) == 0 {
			break
		}

		allImages = append(allImages, imagesResp...)
		// Because the API does not return the value of the `marker` field,
		// we can only manually retrieve the last data from the previous page.
		marker = utils.PathSearch("id", imagesResp[len(imagesResp)-1], "").(string)
	}

	if len(allImages) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	// Filter images by `name_regex`.
	var filteredImages []interface{}
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(nameRegex.(string))
		if err != nil {
			return diag.Errorf("name_regex format error: %s", err)
		}

		for _, image := range allImages {
			if r.MatchString(utils.PathSearch("name", image, "").(string)) {
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
	var image interface{}
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

	return dataSourceImagesImageV2Attributes(region, d, image)
}

// dataSourceImagesImageV2Attributes populates the fields of an Image resource.
func dataSourceImagesImageV2Attributes(region string, d *schema.ResourceData, image interface{}) diag.Diagnostics {
	imageId := utils.PathSearch("id", image, "").(string)
	if imageId == "" {
		return diag.Errorf("error retrieving IMS image: ID is not found in API response")
	}

	d.SetId(imageId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("image_id", imageId),
		d.Set("name", utils.PathSearch("name", image, nil)),
		d.Set("image_type", utils.PathSearch("virtual_env_type", image, nil)),
		d.Set("visibility", flattenVisibility(utils.PathSearch("__imagetype", image, "").(string))),
		d.Set("owner", utils.PathSearch("owner", image, nil)),
		d.Set("os", utils.PathSearch("__platform", image, nil)),
		d.Set("architecture", flattenArchitecture(utils.PathSearch("__support_arm", image, "").(string))),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", image, nil)),
		d.Set("os_version", utils.PathSearch("__os_version", image, nil)),
		d.Set("file", utils.PathSearch("file", image, nil)),
		d.Set("schema", utils.PathSearch("schema", image, nil)),
		d.Set("status", utils.PathSearch("status", image, nil)),
		d.Set("description", utils.PathSearch("__description", image, nil)),
		d.Set("protected", utils.PathSearch("protected", image, nil)),
		d.Set("container_format", utils.PathSearch("container_format", image, nil)),
		d.Set("min_ram_mb", utils.PathSearch("min_ram", image, nil)),
		d.Set("max_ram_mb", flattenMaxRAM(utils.PathSearch("max_ram", image, "").(string))),
		d.Set("min_disk_gb", utils.PathSearch("min_disk", image, nil)),
		d.Set("disk_format", utils.PathSearch("disk_format", image, nil)),
		d.Set("data_origin", utils.PathSearch("__data_origin", image, nil)),
		d.Set("backup_id", utils.PathSearch("__backup_id", image, nil)),
		d.Set("active_at", utils.PathSearch("active_at", image, nil)),
		d.Set("created_at", utils.PathSearch("created_at", image, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", image, nil)),
	)

	if size, err := strconv.Atoi(utils.PathSearch("__image_size", image, "").(string)); err == nil {
		mErr = multierror.Append(mErr, d.Set("size_bytes", size))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

type imageSort []interface{}

func (a imageSort) Len() int { return len(a) }

func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a imageSort) Less(i, j int) bool {
	iTimeStr := utils.PathSearch("updated_at", a[i], "").(string)
	jTimeStr := utils.PathSearch("updated_at", a[j], "").(string)

	iTime, err1 := time.Parse(time.RFC3339, iTimeStr)
	if err1 != nil {
		iTime = time.Time{}
	}

	jTime, err2 := time.Parse(time.RFC3339, jTimeStr)
	if err2 != nil {
		jTime = time.Time{}
	}

	return iTime.Unix() < jTime.Unix()
}

// Returns the most recent Image out of a slice of images.
func mostRecentImage(images []interface{}) interface{} {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}
