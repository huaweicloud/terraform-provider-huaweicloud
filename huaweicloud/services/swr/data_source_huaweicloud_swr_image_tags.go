package swr

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR GET /v2/manage/namespaces/{namespace}/repos/{repository}/tags
func DataSourceImageTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"digest": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"digest": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_trusted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"manifest": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scanned": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"docker_schema": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
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
						"deleted_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceImageTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listImageTagsHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/tags"
		listImageTagsProduct = "swr"
	)

	listImageTagsClient, err := cfg.NewServiceClient(listImageTagsProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	listImageTagsPath := listImageTagsClient.Endpoint + listImageTagsHttpUrl
	listImageTagsPath = strings.ReplaceAll(listImageTagsPath, "{namespace}", organization)
	listImageTagsPath = strings.ReplaceAll(listImageTagsPath, "{repository}", repository)

	offset := 0
	listImageTagsPath += fmt.Sprintf("?limit=10&offset=%v", offset)

	listImageTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	results := make([]map[string]interface{}, 0)
	for {
		listImageTagsResp, err := listImageTagsClient.Request("GET", listImageTagsPath, &listImageTagsOpt)
		if err != nil {
			return diag.Errorf("error querying SWR image tags: %s", err)
		}

		listImageTagsRespBody, err := utils.FlattenResponse(listImageTagsResp)
		if err != nil {
			return diag.Errorf("error retrieving SWR image tags: %s", err)
		}
		imageTags := listImageTagsRespBody.([]interface{})
		if len(imageTags) == 0 {
			break
		}
		for _, imageTag := range imageTags {
			name := utils.PathSearch("Tag", imageTag, "").(string)
			digest := utils.PathSearch("digest", imageTag, "").(string)
			if val, ok := d.GetOk("name"); ok && name != val {
				continue
			}
			if val, ok := d.GetOk("digest"); ok && digest != val {
				continue
			}

			tagType := int(utils.PathSearch("tag_type", imageTag, float64(0)).(float64))
			results = append(results, map[string]interface{}{
				"name":          name,
				"size":          int(utils.PathSearch("size", imageTag, float64(0)).(float64)),
				"path":          utils.PathSearch("path", imageTag, nil),
				"internal_path": utils.PathSearch("internal_path", imageTag, nil),
				"digest":        utils.PathSearch("digest", imageTag, nil),
				"image_id":      utils.PathSearch("image_id", imageTag, nil),
				"is_trusted":    utils.PathSearch("is_trusted", imageTag, false),
				"manifest":      utils.PathSearch("manifest", imageTag, nil),
				"scanned":       utils.PathSearch("scanned", imageTag, false),
				"docker_schema": int(utils.PathSearch("schema", imageTag, float64(0)).(float64)),
				"type":          convertTagType(tagType),
				"created_at":    utils.PathSearch("created", imageTag, nil),
				"updated_at":    utils.PathSearch("updated", imageTag, nil),
				"deleted_at":    utils.PathSearch("deleted", imageTag, nil),
			})
		}
		offset += len(imageTags)
		index := strings.Index(listImageTagsPath, "offset")
		listImageTagsPath = fmt.Sprintf("%soffset=%v", listImageTagsPath[:index], offset)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("image_tags", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func convertTagType(tagType int) string {
	if tagType == 1 {
		return "manifest list"
	}
	return "manifest"
}
