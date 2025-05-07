package vod

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VOD POST /v1.0/{project_id}/asset/category
// @API VOD GET /v1.0/{project_id}/asset/category
// @API VOD PUT /v1.0/{project_id}/asset/category
// @API VOD DELETE /v1.0/{project_id}/asset/category
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
				Type:     schema.TypeString,
				Required: true,
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

func buildCreateMediaCategoryBodyParams(d *schema.ResourceData, parentId int64) map[string]interface{} {
	return map[string]interface{}{
		"name":      d.Get("name"),
		"parent_id": parentId,
	}
}

func resourceMediaCategoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/asset/category"
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	parentId, err := strconv.ParseInt(d.Get("parent_id").(string), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateMediaCategoryBodyParams(d, parentId),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating VOD media category: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", respBody, nil)
	if id == nil {
		return diag.Errorf("error creating VOD media category: ID is not found in API response")
	}

	d.SetId(strconv.FormatInt(int64(id.(float64)), 10))

	return resourceMediaCategoryRead(ctx, d, meta)
}

func flattenChildrenAttribute(categoryResp interface{}) string {
	children, err := utils.JsonMarshal(utils.PathSearch("children", categoryResp, nil))
	if err != nil {
		log.Printf("error marshaling children: %s", err)
	}

	return string(children)
}

func resourceMediaCategoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/asset/category"
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += fmt.Sprintf("?id=%s", d.Id())
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VOD media category")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	categoryResp := utils.PathSearch("[0]", respBody, nil)
	if categoryResp == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", categoryResp, nil)),
		d.Set("children", flattenChildrenAttribute(categoryResp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateMediaCategoryBodyParams(d *schema.ResourceData, id int64) map[string]interface{} {
	return map[string]interface{}{
		"name": d.Get("name"),
		"id":   id,
	}
}

func resourceMediaCategoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/asset/category"
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateMediaCategoryBodyParams(d, id),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating VOD media category: %s", err)
	}

	return resourceMediaCategoryRead(ctx, d, meta)
}

func resourceMediaCategoryDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/asset/category"
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += fmt.Sprintf("?id=%s", d.Id())
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting VOD media category: %s", err)
	}

	return nil
}
