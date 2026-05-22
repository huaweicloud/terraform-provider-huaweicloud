package modelarts

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

var (
	v2WorkflowSubscriptionNonUpdatableParams = []string{
		"workflow_id",
	}
	v2WorkflowSubscriptionNotFoundErrCodes = []string{
		"ModelArts.7512", // The workflow does not exist.
		"ModelArts.7519", // The workflow subscription does not exist.
	}
)

// @API ModelArts POST /v2/{project_id}/workflows/{workflow_id}/subscriptions
// @API ModelArts GET /v2/{project_id}/workflows/{workflow_id}/subscriptions/{subscription_id}
// @API ModelArts PUT /v2/{project_id}/workflows/{workflow_id}/subscriptions/{subscription_id}
// @API ModelArts DELETE /v2/{project_id}/workflows/{workflow_id}/subscriptions/{subscription_id}
func ResourceV2WorkflowSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2WorkflowSubscriptionCreate,
		ReadContext:   resourceV2WorkflowSubscriptionRead,
		UpdateContext: resourceV2WorkflowSubscriptionUpdate,
		DeleteContext: resourceV2WorkflowSubscriptionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2WorkflowSubscriptionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(v2WorkflowSubscriptionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the workflow subscription is located.`,
			},

			// Required parameters.
			"workflow_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workflow to which the subscription belongs.`,
			},
			"topic_urns": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of SMN topic URNs to subscribe.`,
			},

			// Optional parameters.
			"events": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of workflow events to subscribe.`,
			},

			// Attributes.
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the workflow subscription, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildV2WorkflowSubscriptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"topic_urns": d.Get("topic_urns"),
		"events":     utils.ValueIgnoreEmpty(d.Get("events")),
	}
}

func createV2WorkflowSubscription(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/workflows/{workflow_id}/subscriptions"
		workflowId = d.Get("workflow_id").(string)
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workflow_id}", workflowId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV2WorkflowSubscriptionBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceV2WorkflowSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := createV2WorkflowSubscription(client, d)
	if err != nil {
		return diag.Errorf("error creating ModelArts workflow subscription: %s", err)
	}

	subscriptionId := utils.PathSearch("subscription_id", resp, "").(string)
	if subscriptionId == "" {
		return diag.Errorf("unable to find the ModelArts workflow subscription ID from the API response")
	}
	d.SetId(subscriptionId)

	return resourceV2WorkflowSubscriptionRead(ctx, d, meta)
}

func GetV2WorkflowSubscriptionById(client *golangsdk.ServiceClient, workflowId, subscriptionId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/workflows/{workflow_id}/subscriptions/{subscription_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workflow_id}", workflowId)
	getPath = strings.ReplaceAll(getPath, "{subscription_id}", subscriptionId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceV2WorkflowSubscriptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		workflowId     = d.Get("workflow_id").(string)
		subscriptionId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := GetV2WorkflowSubscriptionById(client, workflowId, subscriptionId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", v2WorkflowSubscriptionNotFoundErrCodes...),
			fmt.Sprintf("error retrieving ModelArts workflow subscription (%s)", subscriptionId))
	}

	createdAt := utils.PathSearch("created_at", resp, "").(string)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("topic_urns", utils.PathSearch("topic_urns", resp, make([]interface{}, 0))),
		d.Set("events", utils.PathSearch("events", resp, make([]interface{}, 0))),
	)
	if createdAt != "" {
		mErr = multierror.Append(mErr,
			d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(createdAt)/1000, false)),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func updateV2WorkflowSubscription(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl        = "v2/{project_id}/workflows/{workflow_id}/subscriptions/{subscription_id}"
		workflowId     = d.Get("workflow_id").(string)
		subscriptionId = d.Id()
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workflow_id}", workflowId)
	updatePath = strings.ReplaceAll(updatePath, "{subscription_id}", subscriptionId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV2WorkflowSubscriptionBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceV2WorkflowSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = updateV2WorkflowSubscription(client, d)
	if err != nil {
		return diag.Errorf("error updating ModelArts workflow subscription: %s", err)
	}

	return resourceV2WorkflowSubscriptionRead(ctx, d, meta)
}

func deleteV2WorkflowSubscription(client *golangsdk.ServiceClient, workflowId, subscriptionId string) error {
	httpUrl := "v2/{project_id}/workflows/{workflow_id}/subscriptions/{subscription_id}"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workflow_id}", workflowId)
	deletePath = strings.ReplaceAll(deletePath, "{subscription_id}", subscriptionId)

	opt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func resourceV2WorkflowSubscriptionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		workflowId     = d.Get("workflow_id").(string)
		subscriptionId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = deleteV2WorkflowSubscription(client, workflowId, subscriptionId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", v2WorkflowSubscriptionNotFoundErrCodes...),
			fmt.Sprintf("error deleting ModelArts workflow subscription (%s)", subscriptionId))
	}

	return nil
}

func resourceV2WorkflowSubscriptionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workflow_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("workflow_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
