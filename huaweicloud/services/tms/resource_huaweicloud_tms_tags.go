package tms

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TMS POST /v1.0/predefine_tags/action
// @API TMS GET /v1.0/predefine_tags
func ResourceTmsTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTmsTagCreate,
		ReadContext:   resourceTmsTagRead,
		UpdateContext: resourceTmsTagUpdate,
		DeleteContext: resourceTmsTagDelete,

		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceTmsTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "tms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	err = updateTmsTags(client, "create", d.Get("tags").(*schema.Set).List())
	if err != nil {
		return diag.Errorf("error creating TMS tags: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)

	return resourceTmsTagRead(ctx, d, meta)
}

func resourceTmsTagRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1.0/predefine_tags"
		product = "tms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	getBasePath := client.Endpoint + httpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	rawTagsMap := buildRawTagsMap(d)

	res := make([]interface{}, 0)
	marker := ""
	for {
		getPath := getBasePath + buildGetPath(marker)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving TMS tags: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		res = append(res, flattenGetTmsTagsResponseBody(getRespBody, rawTagsMap)...)
		marker = utils.PathSearch("marker", getRespBody, "").(string)
		if marker == "" {
			break
		}
	}
	if len(res) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(nil,
		d.Set("tags", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetPath(marker string) string {
	res := "?limit=100"
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%s", res, marker)
	}
	return res
}

// for result: the key of the map is the tag key, the value of the first map is a map too, the key of this map is the
// tag value, because there can be more than 1 tag values for a tag key
func buildRawTagsMap(d *schema.ResourceData) map[string]map[string]bool {
	rawTags := d.Get("tags").(*schema.Set).List()
	res := make(map[string]map[string]bool)
	for _, rawTag := range rawTags {
		tag := rawTag.(map[string]interface{})
		key := tag["key"].(string)
		value := tag["value"].(string)
		if v, ok := res[key]; ok {
			v[value] = true
		} else {
			res[key] = map[string]bool{
				value: true,
			}
		}
	}
	return res
}

func flattenGetTmsTagsResponseBody(resp interface{}, rawTagsMap map[string]map[string]bool) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("tags", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0)
	for _, v := range curArray {
		key := utils.PathSearch("key", v, "").(string)
		value := utils.PathSearch("value", v, "").(string)
		if valuesMap, ok := rawTagsMap[key]; ok && valuesMap[value] {
			res = append(res, map[string]string{
				"key":   key,
				"value": value,
			})
		}
	}
	return res
}

func resourceTmsTagUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "tms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	oldRaw, newRaw := d.GetChange("tags")
	addTags := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
	deleteTags := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
	if deleteTags.Len() > 0 {
		err = updateTmsTags(client, "delete", deleteTags.List())
		if err != nil {
			return diag.Errorf("error updating TMS tags: %s", err)
		}
	}
	if addTags.Len() > 0 {
		err = updateTmsTags(client, "create", addTags.List())
		if err != nil {
			return diag.Errorf("error updating TMS tags: %s", err)
		}
	}

	return resourceTmsTagRead(ctx, d, meta)
}

func resourceTmsTagDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "tms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating TMS client: %s", err)
	}

	err = updateTmsTags(client, "delete", d.Get("tags").(*schema.Set).List())
	if err != nil {
		return diag.Errorf("error deleting TMS tags: %s", err)
	}

	return nil
}

func updateTmsTags(client *golangsdk.ServiceClient, action string, rawTags []interface{}) error {
	var (
		httpUrl = "v1.0/predefine_tags/action"
	)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}

	updatePath := client.Endpoint + httpUrl
	updateOpt.JSONBody = utils.RemoveNil(buildTmsTagsBodyParams(action, rawTags))
	_, err := client.Request("POST", updatePath, &updateOpt)
	return err
}

func buildTmsTagsBodyParams(action string, rawTags []interface{}) map[string]interface{} {
	if len(rawTags) == 0 {
		return nil
	}
	tags := make([]map[string]interface{}, 0)
	for _, rawTag := range rawTags {
		if tag, ok := rawTag.(map[string]interface{}); ok {
			tags = append(tags, map[string]interface{}{
				"key":   tag["key"],
				"value": tag["value"],
			})
		}
	}
	bodyPrams := map[string]interface{}{
		"action": action,
		"tags":   tags,
	}
	return bodyPrams
}
