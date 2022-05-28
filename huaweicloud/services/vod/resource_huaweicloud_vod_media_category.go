package vod

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	vod "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vod/v1/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceMediaCategory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMediaCategoryCreate,
		ReadContext:   resourceMediaCategoryRead,
		UpdateContext: resourceMediaCategoryUpdate,
		DeleteContext: resourceMediaCategoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		//request and response parameters
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "0",
			},
			"children": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMediaCategoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	parentId, err := strconv.ParseInt(d.Get("parent_id").(string), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := vod.CreateCategoryReq{
		Name:     d.Get("name").(string),
		ParentId: utils.Int32(int32(parentId)),
	}

	createReq := vod.CreateAssetCategoryRequest{
		Body: &createOpts,
	}

	resp, err := client.CreateAssetCategory(&createReq)
	if err != nil {
		return diag.Errorf("error creating VOD media category: %s", err)
	}

	d.SetId(strconv.FormatInt(int64(*resp.Id), 10))

	return resourceMediaCategoryRead(ctx, d, meta)
}

func resourceMediaCategoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListAssetCategory(&vod.ListAssetCategoryRequest{Id: int32(id)})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VOD media category")
	}

	categoryList := *resp.Body
	if len(categoryList) == 0 {
		log.Printf("unable to retrieve VOD media category: %d", id)
		d.SetId("")
		return nil
	}
	category := categoryList[0]

	children, err := utils.JsonMarshal(category.Children)
	if err != nil {
		log.Printf("error marshaling children: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", category.Name),
		d.Set("children", string(children)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting VOD media category fields: %s", err)
	}

	return nil
}

func resourceMediaCategoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	updateOpts := vod.UpdateCategoryReq{
		Name: d.Get("name").(string),
		Id:   int32(id),
	}

	updateReq := vod.UpdateAssetCategoryRequest{
		Body: &updateOpts,
	}

	_, err = client.UpdateAssetCategory(&updateReq)
	if err != nil {
		return diag.Errorf("error updating VOD media category: %s", err)
	}

	return resourceMediaCategoryRead(ctx, d, meta)
}

func resourceMediaCategoryDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.DeleteAssetCategory(&vod.DeleteAssetCategoryRequest{Id: int32(id)})
	if err != nil {
		return diag.Errorf("error deleting VOD media category: %s", err)
	}

	return nil
}
