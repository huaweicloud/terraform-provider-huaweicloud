package smn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var topicSubscriberNonUpdatableParams = []string{"topic_urn", "subscribe_id"}

// @API SMN POST /v2/{project_id}/notifications/topics/{topic_urn}/subscriptions/from-subscription-users
// @API SMN GET /v2/{project_id}/notifications/topics/{topic_urn}/subscriptions
// @API SMN DELETE /v2/{project_id}/notifications/topics/{topic_urn}/subscriptions/{subscription_urn}
func ResourceTopicSubscriber() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTopicSubscriberCreate,
		ReadContext:   resourceTopicSubscriberRead,
		UpdateContext: resourceTopicSubscriberUpdate,
		DeleteContext: resourceTopicSubscriberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTopicSubscriberImport,
		},

		CustomizeDiff: config.FlexibleForceNew(topicSubscriberNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subscribe_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"filter_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topicSubscribeFilterPoliciesSchema(),
			},
		},
	}
}

func topicSubscribeFilterPoliciesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"string_equals": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func resourceTopicSubscriberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/notifications/topics/{topic_urn}/subscriptions/from-subscription-users"
		product = "smn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{topic_urn}", d.Get("topic_urn").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateTopicSubscriberBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SMN topic subscriber: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	subscriptionUrn := utils.PathSearch("subscriptions_result[0].subscription_urn", createRespBody, "").(string)
	if subscriptionUrn == "" {
		return diag.Errorf("error creating SMN topic subscriber: unable to find the subscription_urn from " +
			"the API response")
	}
	d.SetId(subscriptionUrn)

	return resourceTopicSubscriberRead(ctx, d, meta)
}

func buildCreateTopicSubscriberBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ids": []string{d.Get("subscribe_id").(string)},
	}
	return bodyParams
}

func resourceTopicSubscriberRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v2/{project_id}/notifications/topics/{topic_urn}/subscriptions"
		product = "smn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{topic_urn}", d.Get("topic_urn").(string))

	getDwsAlarmSubsResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SMN topic subscriber")
	}

	getRespJson, err := json.Marshal(getDwsAlarmSubsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	subscription := utils.PathSearch(fmt.Sprintf("subscriptions[?subscription_urn=='%s']|[0]", d.Id()), getRespBody, nil)
	if subscription == nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SMN topic subscriber")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("topic_urn", utils.PathSearch("topic_urn", subscription, nil)),
		d.Set("protocol", utils.PathSearch("protocol", subscription, nil)),
		d.Set("owner", utils.PathSearch("owner", subscription, nil)),
		d.Set("endpoint", utils.PathSearch("endpoint", subscription, nil)),
		d.Set("remark", utils.PathSearch("remark", subscription, nil)),
		d.Set("status", utils.PathSearch("status", subscription, nil)),
		d.Set("filter_policies", flattenTopicSubscriberFilterPolicies(subscription)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTopicSubscriberFilterPolicies(resp interface{}) []interface{} {
	curJson := utils.PathSearch("filter_policies", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":          utils.PathSearch("name", v, nil),
			"string_equals": utils.PathSearch("string_equals", v, nil),
		})
	}
	return rst
}

func resourceTopicSubscriberUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTopicSubscriberDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/notifications/topics/{topic_urn}/subscriptions/{subscription_urn}"
		product = "smn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{topic_urn}", d.Get("topic_urn").(string))
	deletePath = strings.ReplaceAll(deletePath, "{subscription_urn}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SMN topic subscriber")
	}

	return nil
}

func resourceTopicSubscriberImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	subscriptionUrn := d.Id()
	index := strings.LastIndex(subscriptionUrn, ":")
	if index == -1 {
		return nil, errors.New("invalid format: the subscription URN is invalid")
	}
	topicUrn := subscriptionUrn[:index]

	mErr := multierror.Append(nil,
		d.Set("topic_urn", topicUrn),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
