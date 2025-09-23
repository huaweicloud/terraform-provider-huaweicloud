package ims

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
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
			"__support_agent_list": {
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
			"__support_agent_list": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildImageTypeParam(imageType string) string {
	if imageType == "public" {
		imageType = "gold"
	}

	return imageType
}

func buildImagesQueryParamsWithMarker(requestPath, marker string) string {
	if marker == "" {
		return requestPath
	}

	return fmt.Sprintf("%s&marker=%s", requestPath, marker)
}

func buildImagesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?status=active&enterprise_project_id=%v", epsId)

	if v, ok := d.GetOk("image_id"); ok {
		queryParams = fmt.Sprintf("%s&id=%v", queryParams, v)
	}

	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}

	if v, ok := d.GetOk("image_type"); ok {
		queryParams = fmt.Sprintf("%s&virtual_env_type=%v", queryParams, v)
	}

	if d.Get("is_whole_image").(bool) {
		queryParams = fmt.Sprintf("%s&__whole_image=%v", queryParams, true)
	}

	if v, ok := d.GetOk("visibility"); ok {
		queryParams = fmt.Sprintf("%s&__imagetype=%v", queryParams, buildImageTypeParam(v.(string)))
	}

	if v, ok := d.GetOk("owner"); ok {
		queryParams = fmt.Sprintf("%s&owner=%v", queryParams, v)
	}

	if v, ok := d.GetOk("flavor_id"); ok {
		queryParams = fmt.Sprintf("%s&flavor_id=%v", queryParams, v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}

	if v, ok := d.GetOk("sort_direction"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}

	if v, ok := d.GetOk("os"); ok {
		queryParams = fmt.Sprintf("%s&__platform=%v", queryParams, v)
	}

	if v, ok := d.GetOk("architecture"); ok {
		queryParams = fmt.Sprintf("%s&architecture=%v", queryParams, v)
	}

	if v, ok := d.GetOk("tag"); ok {
		queryParams = fmt.Sprintf("%s&tag=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceImagesImagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	for _, image := range allImages {
		if nameRegexRes != nil &&
			!nameRegexRes.MatchString(utils.PathSearch("name", image, "").(string)) {
			continue
		}

		if v, ok := d.GetOk("__support_agent_list"); ok &&
			v.(string) != utils.PathSearch("__support_agent_list", image, "").(string) {
			continue
		}

		ids = append(ids, utils.PathSearch("id", image, "").(string))
		resultImages = append(resultImages, flattenImage(image))
	}

	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("images", resultImages),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenImage(image interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"id":                    utils.PathSearch("id", image, nil),
		"name":                  utils.PathSearch("name", image, nil),
		"image_type":            utils.PathSearch("virtual_env_type", image, nil),
		"visibility":            flattenVisibility(utils.PathSearch("__imagetype", image, "").(string)),
		"owner":                 utils.PathSearch("owner", image, nil),
		"os":                    utils.PathSearch("__platform", image, nil),
		"architecture":          flattenArchitecture(utils.PathSearch("__support_arm", image, "").(string)),
		"enterprise_project_id": utils.PathSearch("enterprise_project_id", image, nil),
		"os_version":            utils.PathSearch("__os_version", image, nil),
		"file":                  utils.PathSearch("file", image, nil),
		"schema":                utils.PathSearch("schema", image, nil),
		"status":                utils.PathSearch("status", image, nil),
		"description":           utils.PathSearch("__description", image, nil),
		"protected":             utils.PathSearch("protected", image, nil),
		"container_format":      utils.PathSearch("container_format", image, nil),
		"min_ram_mb":            utils.PathSearch("min_ram", image, nil),
		"max_ram_mb":            flattenMaxRAM(utils.PathSearch("max_ram", image, "").(string)),
		"min_disk_gb":           utils.PathSearch("min_disk", image, nil),
		"disk_format":           utils.PathSearch("disk_format", image, nil),
		"data_origin":           utils.PathSearch("__data_origin", image, nil),
		"backup_id":             utils.PathSearch("__backup_id", image, nil),
		"active_at":             utils.PathSearch("active_at", image, nil),
		"created_at":            utils.PathSearch("created_at", image, nil),
		"updated_at":            utils.PathSearch("updated_at", image, nil),
		"__support_agent_list":  utils.PathSearch("__support_agent_list", image, nil),
	}

	if size, err := strconv.Atoi(utils.PathSearch("__image_size", image, "").(string)); err == nil {
		result["size_bytes"] = size
	}

	return result
}

func flattenVisibility(imageType string) string {
	if imageType == "gold" {
		imageType = "public"
	}

	return imageType
}
