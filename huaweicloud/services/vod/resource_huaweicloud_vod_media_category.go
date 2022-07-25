package vod

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

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
	region := config.GetRegion(d)
	client, err := config.Client("vod", region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	parentId, err := strconv.ParseInt(d.Get("parent_id").(string), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	url := client.ServiceURL("asset/category")

	reqBody := map[string]interface{}{
		"name":      d.Get("name").(string),
		"parent_id": utils.Int32(int32(parentId)),
	}

	reqOpt := golangsdk.RequestOpts{
		JSONBody:         reqBody,
		KeepResponseBody: true,
	}

	resp, err := client.Request("POST", url, &reqOpt)
	if err != nil {
		return diag.Errorf("error creating VOD media category: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(int64(respBody.(map[string]interface{})["id"].(float64)), 10))

	return resourceMediaCategoryRead(ctx, d, meta)
}

func resourceMediaCategoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.Client("vod", config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	url := client.ServiceURL(fmt.Sprintf("asset/category?id=%s", d.Id()))

	resp, err := client.Request("GET", url, &golangsdk.RequestOpts{KeepResponseBody: true})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VOD media category")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	categoryList := respBody.([]interface{})
	if len(categoryList) == 0 {
		log.Printf("unable to retrieve VOD media category: %s", d.Id())
		d.SetId("")
		return nil
	}
	category := categoryList[0].(map[string]interface{})

	children, err := utils.JsonMarshal(category["children"])
	if err != nil {
		log.Printf("error marshaling children: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", category["name"]),
		d.Set("children", string(children)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting VOD media category fields: %s", err)
	}

	return nil
}

func resourceMediaCategoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.Client("vod", config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	url := client.ServiceURL("asset/category")

	reqBody := map[string]interface{}{
		"name": d.Get("name").(string),
		"id":   d.Id(),
	}

	reqOpt := golangsdk.RequestOpts{
		JSONBody:         reqBody,
		KeepResponseBody: true,
	}

	_, err = client.Request("PUT", url, &reqOpt)
	if err != nil {
		return diag.Errorf("error updating VOD media category: %s", err)
	}

	return resourceMediaCategoryRead(ctx, d, meta)
}

func resourceMediaCategoryDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.Client("vod", config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	url := client.ServiceURL(fmt.Sprintf("asset/category?id=%s", d.Id()))

	_, err = client.Request("DELETE", url, &golangsdk.RequestOpts{})
	if err != nil {
		return diag.Errorf("error deleting VOD media category: %s", err)
	}

	return nil
}
