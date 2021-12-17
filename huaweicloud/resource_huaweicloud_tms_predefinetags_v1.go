package huaweicloud

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	tms "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/tms/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/tms/v1/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/tms/v1/region"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceTMSTagsV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTMSTagsV1Create,
		ReadContext:   resourceTMSTagsV1Read,
		UpdateContext: resourceTMSTagsV1Update,
		DeleteContext: resourceTMSTagsV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Read:   schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"tag": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
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

func flattenTMSTagsOptsV1(tagsOpts []model.PredefineTagRequest) []map[string]interface{} {
	tagsOptsList := make([]map[string]interface{}, len(tagsOpts))
	for i, tagOpt := range tagsOpts {
		tagsOptsList[i] = map[string]interface{}{
			"key":   tagOpt.Key,
			"value": tagOpt.Value,
		}
	}
	return tagsOptsList
}

func getTags(d *schema.ResourceData) ([]model.PredefineTagRequest, string) {
	var predefineTagRequest []model.PredefineTagRequest
	tagsId := ""
	if _, ok := d.GetOk("tag"); ok {
		tags := d.Get("tag").([]interface{})
		if len(tags) != 1 {
			logp.Printf("[ERROR] The number of tag labels is abnormal. Check the label attributes")
			return predefineTagRequest, tagsId
		}
		for _, value := range tags {
			tempValue := value.(map[string]interface{})
			predefineTagrequest := model.PredefineTagRequest{
				Key:   tempValue["key"].(string),
				Value: tempValue["value"].(string),
			}
			timeUnixNano := time.Now().UnixNano()
			id := strconv.FormatInt(timeUnixNano, 10)
			tagsId = predefineTagrequest.Key + ":" + predefineTagrequest.Value + ":" + id
			predefineTagRequest = append(predefineTagRequest, predefineTagrequest)
		}
	}

	return predefineTagRequest, tagsId
}

func getTMSClient(meta interface{}) *tms.TmsClient {
	config := meta.(*config.Config)
	ak := config.AccessKey
	sk := config.SecretKey
	securityToken := config.SecurityToken
	auth := global.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		WithSecurityToken(securityToken). // 在临时aksk场景下使用
		Build()
	client := tms.NewTmsClient(
		tms.TmsClientBuilder().
			WithRegion(region.ValueOf("cn-north-4")).
			WithCredential(auth).
			Build())
	return client
}

func resourceTMSTagsV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getTMSClient(meta)
	request := &model.CreatePredefineTagsRequest{}
	var listTagsbody, id = getTags(d)
	if id == "" {
		return fmtp.DiagErrorf("Error creating Huaweicloud predefine tag")
	}
	request.Body = &model.ReqCreatePredefineTag{
		Tags:   listTagsbody,
		Action: model.GetReqCreatePredefineTagActionEnum().CREATE,
	}
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err := client.CreatePredefineTags(request)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud predefine tag")
	}
	createTag := flattenTMSTagsOptsV1(listTagsbody)
	d.SetId(id)
	d.Set("tag", createTag)

	return resourceTMSTagsV1Read(ctx, d, meta)
}

func resourceTMSTagsV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getTMSClient(meta)
	request := &model.ListPredefineTagsRequest{}
	newTagsId := d.Id()
	firstSpilt := strings.Split(newTagsId, ":")
	if len(firstSpilt) != 3 {
		return fmtp.DiagErrorf("The label ID is incorrect. Please check whether it is correct")
	}
	key := firstSpilt[0]
	value := firstSpilt[1]
	request.Key = &key
	request.Value = &value

	err := resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		response, err := client.ListPredefineTags(request)
		logp.Printf("[DEBUG] Read: %#v", response)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmtp.DiagErrorf("Error reading Huaweicloud predefine tag: %s", err)
	}
	d.SetId(newTagsId)
	return nil
}

func resourceTMSTagsV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getTMSClient(meta)
	request := &model.UpdatePredefineTagsRequest{}
	updateId := d.Id()
	firstSpilt := strings.Split(updateId, ":")
	if len(firstSpilt) != 3 {
		return fmtp.DiagErrorf("The label ID is incorrect. Please check whether it is correct")
	}
	oldTagbody := &model.PredefineTagRequest{
		Key:   firstSpilt[0],
		Value: firstSpilt[1],
	}
	var listTagsbody, _ = getTags(d)
	newTag := listTagsbody[0]
	newTagbody := &model.PredefineTagRequest{
		Key:   newTag.Key,
		Value: newTag.Value,
	}
	request.Body = &model.ModifyPrefineTag{
		OldTag: oldTagbody,
		NewTag: newTagbody,
	}

	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err := client.UpdatePredefineTags(request)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	newTagId := newTagbody.Key + ":" + newTagbody.Value + ":" + firstSpilt[2]
	if err != nil {
		return fmtp.DiagErrorf("Error updating Huaweicloud predefine tag: %s", err)
	}
	d.SetId(newTagId)
	return resourceTMSTagsV1Read(ctx, d, meta)
}

func resourceTMSTagsV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getTMSClient(meta)
	request := &model.DeletePredefineTagsRequest{}
	deleteId := d.Id()
	firstSpilt := strings.Split(deleteId, ":")
	if len(firstSpilt) != 3 {
		return fmtp.DiagErrorf("The label ID is incorrect. Please check whether it is correct")
	}
	var deleteTag []model.PredefineTagRequest
	var tag model.PredefineTagRequest
	tag.Key = firstSpilt[0]
	tag.Value = firstSpilt[1]
	deleteTag = append(deleteTag, tag)

	request.Body = &model.ReqDeletePredefineTag{
		Tags:   deleteTag,
		Action: model.GetReqDeletePredefineTagActionEnum().DELETE,
	}

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err := client.DeletePredefineTags(request)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmtp.DiagErrorf("Error deleting Huaweicloud predefine tag: %s", err)
	}
	return nil
}
