package tms

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/tms/v1/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceTmsTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTmsTagCreate,
		DeleteContext: resourceTmsTagDelete,
		ReadContext:   resourceTmsTagRead,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.All(
								validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fffA-Za-z0-9-_]+$"),
									"The key can only consist of letters, digits, underscores (_) and hyphens (-)."),
								validation.StringLenBetween(1, 36),
							),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.All(
								validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fffA-Za-z0-9-_.]+$"),
									"The key can only consist of letters, digits, periods (.)underscores (_) and hyphens (-)."),
								validation.StringLenBetween(1, 43),
							),
						},
					},
				},
			},
		},
	}
}

func resourceTmsTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcTmsV1Client(c.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud TMS client: %s", err)
	}

	var tagIds []string
	var predefineTags []model.PredefineTagRequest
	tagsRaw := d.Get("tags").([]interface{})
	for _, v := range tagsRaw {
		tag := v.(map[string]interface{})
		predefineTag := model.PredefineTagRequest{
			Key:   tag["key"].(string),
			Value: tag["value"].(string),
		}
		predefineTags = append(predefineTags, predefineTag)
		tagId := fmt.Sprintf("%s:%s", tag["key"], tag["value"])
		tagIds = append(tagIds, tagId)
	}

	createOpts := &model.CreatePredefineTagsRequest{
		Body: &model.ReqCreatePredefineTag{
			Tags:   predefineTags,
			Action: model.GetReqCreatePredefineTagActionEnum().CREATE,
		},
	}

	logp.Printf("[DEBUG] Create TMS tag options: %#v", createOpts)
	_, err = client.CreatePredefineTags(createOpts)
	if err != nil {
		return fmtp.DiagErrorf("Error creating TMS tag: %s", err)
	}

	d.SetId(hashcode.Strings(tagIds))
	return resourceTmsTagRead(ctx, d, meta)
}

func resourceTmsTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcTmsV1Client(c.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud TMS client: %s", err)
	}

	var marker *string
	var tags []model.PredefineTag
	// List all predefine tags
	for {
		request := &model.ListPredefineTagsRequest{
			Marker: marker,
		}

		response, err := client.ListPredefineTags(request)
		if err != nil {
			return fmtp.DiagErrorf("Error listing TMS tags: %s", err)
		}
		tagsResp := *response.Tags
		if len(tagsResp) == 0 {
			break
		} else {
			marker = response.Marker
			tags = append(tags, tagsResp...)
		}
	}

	// Check if the requested tag is missing on cloud side
	var tagList []map[string]interface{}
	tagsRaw := d.Get("tags").([]interface{})
	for _, v := range tagsRaw {
		tag := v.(map[string]interface{})
		key := tag["key"].(string)
		value := tag["value"].(string)

		for _, t := range tags {
			if key == t.Key && value == t.Value {
				tagFound := map[string]interface{}{
					"key":   key,
					"value": value,
				}
				tagList = append(tagList, tagFound)
			}
		}
	}
	d.Set("tags", tagList)

	return nil
}

func resourceTmsTagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcTmsV1Client(c.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud TMS client: %s", err)
	}

	var predefineTags []model.PredefineTagRequest
	tagsRaw := d.Get("tags").([]interface{})
	if len(tagsRaw) == 0 {
		logp.Printf("[DEBUG] TMS tags are empty, no need to issue delete request")
		return nil
	}
	for _, v := range tagsRaw {
		tag := v.(map[string]interface{})
		predefineTag := model.PredefineTagRequest{
			Key:   tag["key"].(string),
			Value: tag["value"].(string),
		}
		predefineTags = append(predefineTags, predefineTag)
	}

	deleteOpts := &model.DeletePredefineTagsRequest{
		Body: &model.ReqDeletePredefineTag{
			Tags:   predefineTags,
			Action: model.GetReqDeletePredefineTagActionEnum().DELETE,
		},
	}

	logp.Printf("[DEBUG] Delete TMS tag options: %#v", deleteOpts)
	_, err = client.DeletePredefineTags(deleteOpts)
	if err != nil {
		return fmtp.DiagErrorf("Error deleting TMS tag: %s", err)
	}

	return nil
}
