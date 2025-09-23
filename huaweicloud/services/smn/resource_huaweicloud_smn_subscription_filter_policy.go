package smn

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var filterPolicyNonUpdatableParams = []string{"subscription_urn"}

// @API SMN POST /v2/{project_id}/notifications/subscriptions/filter_polices
// @API SMN GET /v2/{project_id}/notifications/topics/{topicUrn}/subscriptions
// @API SMN PUT /v2/{project_id}/notifications/subscriptions/filter_polices
// @API SMN DELETE /v2/{project_id}/notifications/subscriptions/filter_polices
func ResourceSubscriptionFilterPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubscriptionFilterPolicyCreate,
		ReadContext:   resourceSubscriptionFilterPolicyRead,
		UpdateContext: resourceSubscriptionFilterPolicyUpdate,
		DeleteContext: resourceSubscriptionFilterPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(filterPolicyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"subscription_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource identifier of the subscriber.`,
			},
			"filter_policies": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: `Specifies the message filter policies of a subscriber.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the filter policy name. The policy name must be unique.`,
						},
						"string_equals": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the string array for exact match.`,
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

func resourceSubscriptionFilterPolicyCreate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	subscriptionUrn := d.Get("subscription_urn").(string)

	// createFilterPolicy: create SMN subscription filter policy
	createFilterPolicyHttpUrl := "v2/{project_id}/notifications/subscriptions/filter_polices"
	createFilterPolicyPath := client.Endpoint + createFilterPolicyHttpUrl
	createFilterPolicyPath = strings.ReplaceAll(createFilterPolicyPath, "{project_id}", client.ProjectID)

	createFilterPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createFilterPolicyOpt.JSONBody = map[string]interface{}{
		"polices": []map[string]interface{}{
			{
				"subscription_urn": subscriptionUrn,
				"filter_polices":   buildFilterPolicies(d.Get("filter_policies").(*schema.Set).List()),
			},
		},
	}

	_, err = client.Request("POST", createFilterPolicyPath, &createFilterPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating SMN subscription filter policy: %s", err)
	}

	d.SetId(subscriptionUrn)
	return resourceSubscriptionFilterPolicyRead(ctx, d, meta)
}

func resourceSubscriptionFilterPolicyRead(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	filterPolicies, err := GetSubscriptionFilterPolicies(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying the subscription filter policy")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("subscription_urn", d.Id()),
		d.Set("filter_policies", flattenFilterPolicies(filterPolicies)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSubscriptionFilterPolicyUpdate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	// updateFilterPolicy: create SMN subscription filter policy
	updateFilterPolicyHttpUrl := "v2/{project_id}/notifications/subscriptions/filter_polices"
	updateFilterPolicyPath := client.Endpoint + updateFilterPolicyHttpUrl
	updateFilterPolicyPath = strings.ReplaceAll(updateFilterPolicyPath, "{project_id}", client.ProjectID)

	updateFilterPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateFilterPolicyOpt.JSONBody = map[string]interface{}{
		"polices": []map[string]interface{}{
			{
				"subscription_urn": d.Id(),
				"filter_polices":   buildFilterPolicies(d.Get("filter_policies").(*schema.Set).List()),
			},
		},
	}

	_, err = client.Request("PUT", updateFilterPolicyPath, &updateFilterPolicyOpt)
	if err != nil {
		return diag.Errorf("error updating SMN subscription filter policy: %s", err)
	}

	return resourceSubscriptionFilterPolicyRead(ctx, d, meta)
}

func resourceSubscriptionFilterPolicyDelete(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	// deleteFilterPolicy: create SMN subscription filter policy
	deleteFilterPolicyHttpUrl := "v2/{project_id}/notifications/subscriptions/filter_polices"
	deleteFilterPolicyPath := client.Endpoint + deleteFilterPolicyHttpUrl
	deleteFilterPolicyPath = strings.ReplaceAll(deleteFilterPolicyPath, "{project_id}", client.ProjectID)

	deleteFilterPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteFilterPolicyOpt.JSONBody = map[string]interface{}{
		"subscription_urns": []string{d.Id()},
	}

	_, err = client.Request("DELETE", deleteFilterPolicyPath, &deleteFilterPolicyOpt)
	if err != nil {
		return diag.Errorf("error deleting SMN subscription (%s) filter policy: %s", d.Id(), err)
	}

	// Successful deletion API call does not guarantee that the resource is successfully deleted.
	// Call the details API to confirm that the resource has been successfully deleted.
	_, err = GetSubscriptionFilterPolicies(client, d.Id())
	if err == nil {
		return diag.Errorf("error deleting the subscription filter policy")
	}

	return nil
}

func GetSubscriptionFilterPolicies(client *golangsdk.ServiceClient, subscriptionUrn string) ([]interface{}, error) {
	// we can get topic_urn from the subscription_urn.
	// topic_urn        "urn:smn:cn-south-1:09f960944c80f4802f85c003e0ed1d18:test_topic"
	// subscription_urn "urn:smn:cn-south-1:09f960944c80f4802f85c003e0ed1d18:test_topic:699830a8eed442fa93ab41c6bd1fee11"
	index := strings.LastIndex(subscriptionUrn, ":")
	if index == -1 {
		return nil, fmt.Errorf("this is not a subscription URN")
	}
	topicUrn := subscriptionUrn[:index]

	// getFilterPolicy: query SMN subscription filter policy
	getFilterPolicyHttpUrl := "v2/{project_id}/notifications/topics/{topicUrn}/subscriptions"
	getFilterPolicyPath := client.Endpoint + getFilterPolicyHttpUrl
	getFilterPolicyPath = strings.ReplaceAll(getFilterPolicyPath, "{project_id}", client.ProjectID)
	getFilterPolicyPath = strings.ReplaceAll(getFilterPolicyPath, "{topicUrn}", topicUrn)

	getFilterPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	var subscription interface{}
	for {
		curPath := fmt.Sprintf("%s?offset=%d&limit=100", getFilterPolicyPath, offset)
		getFilterPolicyResp, err := client.Request("GET", curPath, &getFilterPolicyOpt)
		if err != nil {
			return nil, err
		}
		getFilterPolicyRespBody, err := utils.FlattenResponse(getFilterPolicyResp)
		if err != nil {
			return nil, fmt.Errorf("error flattening the subscription filter policy: %s", err)
		}
		curSubscriptions := utils.PathSearch("subscriptions", getFilterPolicyRespBody,
			make([]interface{}, 0)).([]interface{})
		expression := fmt.Sprintf("subscriptions|[?subscription_urn=='%s']|[0]", subscriptionUrn)
		subscription = utils.PathSearch(expression, getFilterPolicyRespBody, nil)
		if subscription != nil {
			break
		}

		if len(curSubscriptions) < 100 {
			break
		}
		offset += len(curSubscriptions)
	}

	filterPolicies := utils.PathSearch("filter_polices", subscription, make([]interface{}, 0)).([]interface{})
	if len(filterPolicies) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return filterPolicies, nil
}

func buildFilterPolicies(filterPolicies []interface{}) []map[string]interface{} {
	opts := make([]map[string]interface{}, len(filterPolicies))
	for i, v := range filterPolicies {
		opts[i] = map[string]interface{}{
			"name":          utils.PathSearch("name", v, nil),
			"string_equals": utils.PathSearch("string_equals", v, nil).(*schema.Set).List(),
		}
	}
	return opts
}

func flattenFilterPolicies(filterPolicies []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, len(filterPolicies))
	for i, v := range filterPolicies {
		rst[i] = map[string]interface{}{
			"name":          utils.PathSearch("name", v, nil),
			"string_equals": utils.PathSearch("string_equals", v, nil),
		}
	}
	return rst
}
