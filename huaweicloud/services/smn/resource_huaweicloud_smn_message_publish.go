package smn

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var messagePublishNonUpdatableParams = []string{
	"topic_urn",
	"subject",
	"message",
	"message_structure",
	"message_template_name",
	"tag",
	"time_to_live",
	"message_attributes",
}

// @API SMN POST /v2/{project_id}/notifications/topics/{topic_urn}/publish
func ResourceMessagePublish() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMessagePublishCreate,
		UpdateContext: resourceMessagePublishUpdate,
		ReadContext:   resourceMessagePublishRead,
		DeleteContext: resourceMessagePublishDelete,

		CustomizeDiff: config.FlexibleForceNew(messagePublishNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource identifier of a topic.`,
			},
			"subject": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the message title.`,
			},
			"message": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"message", "message_structure", "message_template_name"},
				Description:  `Specifies the message content.`,
			},
			"message_structure": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"message", "message_structure", "message_template_name"},
				Description:  `Specifies the message structure.`,
			},
			"message_template_name": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"tags"},
				ExactlyOneOf: []string{"message", "message_structure", "message_template_name"},
				Description:  `Specifies the message template name.`,
			},
			"tags": common.TagsSchema(),
			"time_to_live": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the maximum retention time of the message within the SMN system.`,
			},
			"message_attributes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the message filter policies of a subscriber.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the property name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the property type.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the property value.`,
						},
						"values": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the property values.`,
						},
					},
				},
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

func resourceMessagePublishCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	topicUrn := d.Get("topic_urn").(string)

	client, err := cfg.NewServiceClient("smn", region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	// createMessagePublish: create SMN message publishment
	createMessagePublishHttpUrl := "v2/{project_id}/notifications/topics/{topic_urn}/publish"
	createMessagePublishPath := client.Endpoint + createMessagePublishHttpUrl
	createMessagePublishPath = strings.ReplaceAll(createMessagePublishPath, "{project_id}", client.ProjectID)
	createMessagePublishPath = strings.ReplaceAll(createMessagePublishPath, "{topic_urn}", topicUrn)

	createMessagePublishOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createMessagePublishOpt.JSONBody = utils.RemoveNil(buildCreateMessagePublishBodyParams(d))
	createMessagePublishResp, err := client.Request("POST", createMessagePublishPath, &createMessagePublishOpt)
	if err != nil {
		return diag.Errorf("error creating SMN message publishment: %s", err)
	}

	createMessagePublishRespBody, err := utils.FlattenResponse(createMessagePublishResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("message_id", createMessagePublishRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SMN message publishment: message ID is not found in API response")
	}
	d.SetId(id)

	return nil
}

func buildCreateMessagePublishBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"subject":               utils.ValueIgnoreEmpty(d.Get("subject")),
		"message":               utils.ValueIgnoreEmpty(d.Get("message")),
		"message_structure":     utils.ValueIgnoreEmpty(d.Get("message_structure")),
		"message_template_name": utils.ValueIgnoreEmpty(d.Get("message_template_name")),
		"tags":                  utils.ValueIgnoreEmpty(d.Get("tags")),
		"time_to_live":          utils.ValueIgnoreEmpty(d.Get("time_to_live")),
	}

	if arr, ok := d.GetOk("message_attributes"); ok {
		curArr := arr.([]interface{})
		rst := make([]interface{}, len(curArr))
		for i, v := range curArr {
			attrType := utils.PathSearch("type", v, "").(string)
			var attrValue interface{}
			if attrType == "STRING" {
				attrValue = utils.PathSearch("value", v, nil)
			} else {
				attrValue = utils.PathSearch("values", v, nil)
			}
			rst[i] = map[string]interface{}{
				"name":  utils.PathSearch("name", v, nil),
				"type":  attrType,
				"value": attrValue,
			}
		}
		bodyParams["message_attributes"] = rst
	}

	return bodyParams
}

func resourceMessagePublishRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMessagePublishUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMessagePublishDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting message pbulish resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
