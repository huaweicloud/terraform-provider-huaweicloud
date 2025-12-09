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

var notifyPolicyNonUpdatableParams = []string{"topic_urn"}

// @API SMN POST /v2/{project_id}/notifications/topics/{topic_urn}/notify-policy
// @API SMN GET /v2/{project_id}/notifications/topics/{topic_urn}/notify-policy
// @API SMN PUT /v2/{project_id}/notifications/topics/{topic_urn}/notify-policy/{notify_policy_id}
// @API SMN DELETE /v2/{project_id}/notifications/topics/{topic_urn}/notify-policy/{notify_policy_id}
func ResourceNotifyPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotifyPolicyCreate,
		ReadContext:   resourceNotifyPolicyRead,
		UpdateContext: resourceNotifyPolicyUpdate,
		DeleteContext: resourceNotifyPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNotifyPolicyImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(notifyPolicyNonUpdatableParams),

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
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"polling": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     notifyPolicyPollingSchema(),
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

func notifyPolicyPollingSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"order": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"subscription_urns": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     notifyPolicyPollingSubscriptionsSchema(),
			},
		},
	}
	return &sc
}

func notifyPolicyPollingSubscriptionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"subscription_urn": {
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
		},
	}
	return &sc
}

func resourceNotifyPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/notifications/topics/{topic_urn}/notify-policy"
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
	createOpt.JSONBody = utils.RemoveNil(buildUpdateNotifyPolicyBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SMN notify policy: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the SMN notify policy ID from the API response")
	}
	d.SetId(id)

	return resourceNotifyPolicyRead(ctx, d, meta)
}

func buildUpdateNotifyPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"protocol": d.Get("protocol"),
		"polling":  buildCreateNotifyPolicyPollingBodyParams(d),
	}
	return bodyParams
}

func buildCreateNotifyPolicyPollingBodyParams(d *schema.ResourceData) []interface{} {
	rawParams := d.Get("polling").(*schema.Set)
	if rawParams.Len() == 0 {
		return nil
	}

	rst := make([]interface{}, 0, rawParams.Len())
	for _, rawParam := range rawParams.List() {
		if v, ok := rawParam.(map[string]interface{}); ok {
			rst = append(rst, map[string]interface{}{
				"order":             v["order"],
				"subscription_urns": v["subscription_urns"].(*schema.Set).List(),
			})
		}
	}

	return rst
}

func resourceNotifyPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "smn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	getRespBody, err := getNotifyPolicy(client, d.Get("topic_urn").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", getRespBody, "").(string)
	if id == "" || id != d.Id() {
		return common.CheckDeletedDiag(d, err, "error retrieving SMN notify policy")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("protocol", utils.PathSearch("protocol", getRespBody, nil)),
		d.Set("polling", flattenNotifyPolicyPolling(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getNotifyPolicy(client *golangsdk.ServiceClient, topicUrn string) (interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/notifications/topics/{topic_urn}/notify-policy"
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{topic_urn}", topicUrn)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func flattenNotifyPolicyPolling(resp interface{}) []interface{} {
	curJson := utils.PathSearch("polling", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		subscriptions, subscriptionUrns := flattenNotifyPolicyPollingSubscriptions(v)
		rst = append(rst, map[string]interface{}{
			"order":             utils.PathSearch("order", v, nil),
			"subscription_urns": subscriptionUrns,
			"subscriptions":     subscriptions,
		})
	}
	return rst
}

func flattenNotifyPolicyPollingSubscriptions(resp interface{}) ([]interface{}, []string) {
	curJson := utils.PathSearch("subscriptions", resp, nil)
	if curJson == nil {
		return nil, nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	subscriptionUrns := make([]string, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"subscription_urn": utils.PathSearch("subscription_urn", v, nil),
			"endpoint":         utils.PathSearch("endpoint", v, nil),
			"remark":           utils.PathSearch("subscription_urn", v, nil),
			"status":           utils.PathSearch("status", v, nil),
		})
		subscriptionUrns = append(subscriptionUrns, utils.PathSearch("subscription_urn", v, "").(string))
	}
	return rst, subscriptionUrns
}

func resourceNotifyPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/notifications/topics/{topic_urn}/notify-policy/{notify_policy_id}"
		product = "smn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SMN Client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{topic_urn}", d.Get("topic_urn").(string))
	updatePath = strings.ReplaceAll(updatePath, "{notify_policy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateNotifyPolicyBodyParams(d)

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SMN notify policy: %s", err)
	}

	return resourceNotifyPolicyRead(ctx, d, meta)
}

func resourceNotifyPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/notifications/topics/{topic_urn}/notify-policy/{notify_policy_id}"
		product = "smn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{topic_urn}", d.Get("topic_urn").(string))
	deletePath = strings.ReplaceAll(deletePath, "{notify_policy_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SMN notify policy")
	}

	return nil
}

func resourceNotifyPolicyImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "smn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SMN client: %s", err)
	}

	getRespBody, err := getNotifyPolicy(client, d.Id())
	if err != nil {
		return nil, err
	}

	id := utils.PathSearch("id", getRespBody, "").(string)
	if id == "" {
		return nil, golangsdk.ErrDefault404{}
	}

	mErr := multierror.Append(nil,
		d.Set("topic_urn", d.Id()),
	)
	d.SetId(id)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
