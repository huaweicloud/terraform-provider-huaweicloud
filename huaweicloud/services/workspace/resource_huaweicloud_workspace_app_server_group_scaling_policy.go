package workspace

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

var appServerGroupScalingPolicyNonUpdatableParams = []string{"server_group_id"}

// @API WorkspaceApp PUT /v1/{project_id}/scaling-policy
// @API WorkspaceApp GET /v1/{project_id}/scaling-policy
// @API WorkspaceApp DELETE /v1/{project_id}/scaling-policy
func ResourceAppServerGroupScalingPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppServerGroupScalingPolicyCreate,
		ReadContext:   resourceAppServerGroupScalingPolicyRead,
		UpdateContext: resourceAppServerGroupScalingPolicyUpdate,
		DeleteContext: resourceAppServerGroupScalingPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAppServerGroupScalingPolicyImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(appServerGroupScalingPolicyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the scaling policy is located.`,
			},
			// Required parameters
			"server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the server group to which the scaling policy belongs.`,
			},
			"max_scaling_amount": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The maximum number of instances that can be scaled out.`,
			},
			"single_expansion_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The number of instances to scale out in a single scaling operation.`,
			},
			"scaling_policy_by_session": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"session_usage_threshold": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The total session usage threshold of the server group.`,
						},
						"shrink_after_session_idle_minutes": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The number of minutes to wait before releasing instances with no session connections.`,
						},
					},
				},
				Description: `The session-based scaling policy configuration.`,
			},
			"enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the scaling policy is enabled.`,
			},
			// Internal parameters
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildScalingPolicyBySession(scalingPolicyBySession []interface{}) map[string]interface{} {
	if len(scalingPolicyBySession) < 1 {
		return nil
	}

	policy := scalingPolicyBySession[0]
	return map[string]interface{}{
		"session_usage_threshold":           utils.PathSearch("session_usage_threshold", policy, nil),
		"shrink_after_session_idle_minutes": utils.PathSearch("shrink_after_session_idle_minutes", policy, nil),
	}
}

func buildScalingPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"server_group_id":           d.Get("server_group_id"),
		"max_scaling_amount":        d.Get("max_scaling_amount"),
		"single_expansion_count":    d.Get("single_expansion_count"),
		"enable":                    true,
		"scaling_policy_by_session": buildScalingPolicyBySession(d.Get("scaling_policy_by_session").([]interface{})),
	}
	return bodyParams
}

func resourceAppServerGroupScalingPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	httpUrl := "v1/{project_id}/scaling-policy"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildScalingPolicyBodyParams(d),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("PUT", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error enabling Workspace APP scaling policy under a specified server group (%s): %s",
			d.Get("server_group_id"), err)
	}

	d.SetId(d.Get("server_group_id").(string))
	return resourceAppServerGroupScalingPolicyRead(ctx, d, meta)
}

// GetAppServerGroupScalingPolicy is a method used to get the scaling policy of a server group.
func GetAppServerGroupScalingPolicy(client *golangsdk.ServiceClient, serverGroupId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/scaling-policy"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?server_group_id=%s", getPath, serverGroupId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	response, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return nil, err
	}

	isScalingPolicyEnabled := utils.PathSearch("enable", respBody, false).(bool)
	if !isScalingPolicyEnabled {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/scaling-policy",
				RequestId: "NONE",
				Body:      []byte("the scaling policy has been disabled"),
			},
		}
	}
	return respBody, nil
}

func resourceAppServerGroupScalingPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	respBody, err := GetAppServerGroupScalingPolicy(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error retrieving Workspace APP scaling policy under a specified server group (%s)", d.Id()))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("enable", utils.PathSearch("enable", respBody, true)),
		d.Set("max_scaling_amount", utils.PathSearch("max_scaling_amount", respBody, nil)),
		d.Set("single_expansion_count", utils.PathSearch("single_expansion_count", respBody, nil)),
		d.Set("scaling_policy_by_session", flattenScalingPolicyBySession(utils.PathSearch("scaling_policy_by_session",
			respBody, make(map[string]interface{})).(map[string]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAppServerGroupScalingPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	httpUrl := "v1/{project_id}/scaling-policy"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildScalingPolicyBodyParams(d),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error updating Workspace APP scaling policy under a specified server group (%s): %s",
			d.Get("server_group_id"), err)
	}

	return resourceAppServerGroupScalingPolicyRead(ctx, d, meta)
}

func resourceAppServerGroupScalingPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	httpUrl := "v1/{project_id}/scaling-policy"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = fmt.Sprintf("%s?server_group_id=%s", deletePath, d.Id())

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return diag.Errorf("error disabling Workspace APP scaling policy from the specified server group (%s): %s",
			d.Get("server_group_id"), err)
	}

	return nil
}

func resourceAppServerGroupScalingPolicyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("server_group_id", d.Id())
}
