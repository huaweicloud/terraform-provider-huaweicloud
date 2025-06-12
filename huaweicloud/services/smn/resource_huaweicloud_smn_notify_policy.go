package smn

import (
	"context"
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
				Elem: &schema.Resource{
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

func resourceNotifyPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("smn", region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	createNotifyPolicyHttpUrl := "v2/{project_id}/notifications/topics/{topic_urn}/notify-policy"
	createNotifyPolicyPath := client.Endpoint + createNotifyPolicyHttpUrl
	createNotifyPolicyPath = strings.ReplaceAll(createNotifyPolicyPath, "{project_id}", client.ProjectID)
	createNotifyPolicyPath = strings.ReplaceAll(createNotifyPolicyPath, "{topic_urn}", d.Get("topic_urn").(string))

	createNotifyPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createNotifyPolicyOpt.JSONBody = map[string]interface{}{
		"protocol": d.Get("protocol").(string),
		"polling":  buildPolling(d.Get("polling").(*schema.Set).List()),
	}

	createNotifyPolicyResp, err := client.Request("POST", createNotifyPolicyPath, &createNotifyPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating SMN notify policy: %s", err)
	}

	createNotifyPolicyRespBody, err := utils.FlattenResponse(createNotifyPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createNotifyPolicyRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SMN notify policy: ID is not found in API response")
	}
	d.SetId(id)
	return resourceNotifyPolicyRead(ctx, d, meta)
}

func buildPolling(polling []interface{}) []map[string]interface{} {
	opts := make([]map[string]interface{}, len(polling))
	for i, v := range polling {
		opts[i] = map[string]interface{}{
			"order":             utils.PathSearch("order", v, nil),
			"subscription_urns": utils.PathSearch("subscription_urns", v, &schema.Set{}).(*schema.Set).List(),
		}
	}
	return opts
}

func resourceNotifyPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("smn", region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	readNotifyPolicyHttpUrl := "v2/{project_id}/notifications/topics/{topic_urn}/notify-policy"
	readNotifyPolicyPath := client.Endpoint + readNotifyPolicyHttpUrl
	readNotifyPolicyPath = strings.ReplaceAll(readNotifyPolicyPath, "{project_id}", client.ProjectID)
	readNotifyPolicyPath = strings.ReplaceAll(readNotifyPolicyPath, "{topic_urn}", d.Get("topic_urn").(string))

	readNotifyPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	readNotifyPolicyResp, err := client.Request("POST", readNotifyPolicyPath, &readNotifyPolicyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying the subscription filter policy")
	}

	readNotifyPolicyRespBody, err := utils.FlattenResponse(readNotifyPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("protocol", utils.PathSearch("protocol", readNotifyPolicyRespBody, nil)),
		d.Set("polling", flattenPolling(utils.PathSearch("polling", readNotifyPolicyRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPolling(polling []interface{}) []map[string]interface{} {
	opts := make([]map[string]interface{}, len(polling))
	for i, v := range polling {
		opts[i] = map[string]interface{}{
			"order":             utils.PathSearch("order", v, nil),
			"subscription_urns": utils.PathSearch("subscription_urns", v, nil),
		}
	}
	return opts
}

func resourceNotifyPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	updateNotifyPolicyHttpUrl := "v2/{project_id}/notifications/topics/{topic_urn}/notify-policy/{notify_policy_id}"
	updateNotifyPolicyPath := client.Endpoint + updateNotifyPolicyHttpUrl
	updateNotifyPolicyPath = strings.ReplaceAll(updateNotifyPolicyPath, "{project_id}", client.ProjectID)
	updateNotifyPolicyPath = strings.ReplaceAll(updateNotifyPolicyPath, "{topic_urn}", d.Get("topic_urn").(string))
	updateNotifyPolicyPath = strings.ReplaceAll(updateNotifyPolicyPath, "{notify_policy_id}", d.Id())

	updateNotifyPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateNotifyPolicyOpt.JSONBody = map[string]interface{}{
		"polices": []map[string]interface{}{
			{
				"subscription_urn": d.Id(),
				"filter_polices":   buildFilterPolicies(d.Get("filter_policies").(*schema.Set).List()),
			},
		},
	}

	_, err = client.Request("PUT", updateNotifyPolicyPath, &updateNotifyPolicyOpt)
	if err != nil {
		return diag.Errorf("error updating SMN subscription filter policy: %s", err)
	}

	return resourceNotifyPolicyRead(ctx, d, meta)
}

func resourceNotifyPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	deleteNotifyPolicyHttpUrl := "v2/{project_id}/notifications/topics/{topic_urn}/notify-policy/{notify_policy_id}"
	deleteNotifyPolicyPath := client.Endpoint + deleteNotifyPolicyHttpUrl
	deleteNotifyPolicyPath = strings.ReplaceAll(deleteNotifyPolicyPath, "{project_id}", client.ProjectID)
	deleteNotifyPolicyPath = strings.ReplaceAll(deleteNotifyPolicyPath, "{topic_urn}", d.Get("topic_urn").(string))
	deleteNotifyPolicyPath = strings.ReplaceAll(deleteNotifyPolicyPath, "{notify_policy_id}", d.Id())

	deleteNotifyPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deleteNotifyPolicyPath, &deleteNotifyPolicyOpt)
	if err != nil {
		return diag.Errorf("error deleting SMN notify policy: %s", err)
	}

	// Successful deletion API call does not guarantee that the resource is successfully deleted.
	// Call the details API to confirm that the resource has been successfully deleted.
	_, err = GetSubscriptionFilterPolicies(client, d.Id())
	if err == nil {
		return diag.Errorf("error deleting the subscription filter policy")
	}

	return nil
}
