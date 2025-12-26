package cce

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE GET /v5/imagecaches
func DataSourceCCEImageCaches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEImageCachesRead,

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
			"image_caches": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"images": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_cache_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"retention_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type ImageCachesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newCCEImageCachesDSWrapper(d *schema.ResourceData, meta interface{}) *ImageCachesDSWrapper {
	return &ImageCachesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCCEImageCachesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newCCEImageCachesDSWrapper(d, meta)
	imageCachesRst, err := wrapper.getImageCaches()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.imageCachesToSchema(imageCachesRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCE GET /v5/imagecaches
func (w *ImageCachesDSWrapper) getImageCaches() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cce")
	if err != nil {
		return nil, err
	}

	uri := "/v5/imagecaches"
	params := map[string]any{
		"name": w.Get("name"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *ImageCachesDSWrapper) imageCachesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("image_caches", schemas.SliceToList(body.Get("image_caches"),
			func(imagecache gjson.Result) any {
				return map[string]any{
					"name":             imagecache.Get("kind").Value(),
					"id":               imagecache.Get("apiVersion").Value(),
					"created_at":       imagecache.Get("name").Value(),
					"images":           imagecache.Get("policyId").Value(),
					"image_cache_size": imagecache.Get("policyType").Value(),
					"retention_days":   imagecache.Get("createTime").Value(),
					"status":           imagecache.Get("updateTime").Value(),
					"message":          imagecache.Get("updateTime").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
