package ims

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var imageMetadataNonUpdatableParams = []string{"__os_version", "visibility", "name1", "protected", "container_format",
	"disk_format", "tags", "min_ram", "min_disk"}

// @API IMS POST /v2/images
// ResourceImageMetadata is a definition of the one-time action resource that used to manage image metadata.
func ResourceImageMetadata() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImageMetadataCreate,
		ReadContext:   resourceImageMetadataRead,
		UpdateContext: resourceImageMetadataUpdate,
		DeleteContext: resourceImageMetadataDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"__os_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name1": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protected": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"container_format": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_format": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"min_ram": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_disk": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildCreateImageMetadataBodyParam(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"__os_version": utils.ValueIgnoreEmpty(d.Get("__os_version")),
		"visibility":   utils.ValueIgnoreEmpty(d.Get("visibility")),
		"name":         utils.ValueIgnoreEmpty(d.Get("name1")),
		// The default value of the `protected` field in the API documentation is **false**.
		"protected":        d.Get("protected"),
		"container_format": utils.ValueIgnoreEmpty(d.Get("container_format")),
		"disk_format":      utils.ValueIgnoreEmpty(d.Get("disk_format")),
		"tags":             utils.ExpandToStringList(d.Get("tags").([]interface{})),
		"min_ram":          utils.ValueIgnoreEmpty(d.Get("min_ram")),
		"min_disk":         utils.ValueIgnoreEmpty(d.Get("min_disk")),
	}

	return bodyParams
}

func resourceImageMetadataCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ims"
		httpUrl = "v2/images"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IMS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateImageMetadataBodyParam(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IMS image metadata: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	imageId := utils.PathSearch("id", createRespBody, "").(string)
	if imageId == "" {
		return diag.Errorf("error creating IMS image metadata: ID is not found in API response")
	}

	d.SetId(imageId)

	return resourceImageMetadataRead(ctx, d, meta)
}

func resourceImageMetadataRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceImageMetadataUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceImageMetadataDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a one-time action resource.
	return nil
}
