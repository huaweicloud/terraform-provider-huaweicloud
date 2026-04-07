package organizations

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var policyDryRunConfigurationNonUpdatableParams = []string{"root_id", "policy_type"}

// @API Organizations POST /v1/organizations/dry-run-config
// @API Organizations GET /v1/organizations/dry-run-config
func ResourcePolicyDryRunConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyDryRunConfigurationCreate,
		ReadContext:   resourcePolicyDryRunConfigurationRead,
		UpdateContext: resourcePolicyDryRunConfigurationUpdate,
		DeleteContext: resourcePolicyDryRunConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyDryRunConfigurationImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(policyDryRunConfigurationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"root_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the organization's root.`,
			},
			"policy_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the policy.`,
			},
			// Optional parameters.
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the policy dry-run.`,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the OBS bucket.`,
			},
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the OBS bucket is located.`,
			},
			"bucket_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The prefix of the OBS bucket.`,
			},
			"agency_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the IAM agency.`,
			},
			// Attributes.
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the policy dry-run configuration.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the policy dry-run configuration.`,
			},
			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func resourcePolicyDryRunConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	err = updatePolicyDryRunConfiguration(ctx, client, buildPolicyDryRunConfigurationBodyParams(d, d.Get("status").(string)),
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error creating policy dry-run configuration: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("root_id").(string), d.Get("policy_type").(string)))

	return resourcePolicyDryRunConfigurationRead(ctx, d, meta)
}

func buildPolicyDryRunConfigurationBodyParams(d *schema.ResourceData, status string) map[string]interface{} {
	return map[string]interface{}{
		"root_id":       d.Get("root_id"),
		"policy_type":   d.Get("policy_type"),
		"status":        utils.ValueIgnoreEmpty(status),
		"bucket_name":   utils.ValueIgnoreEmpty(d.Get("bucket_name")),
		"region_id":     utils.ValueIgnoreEmpty(d.Get("region_id")),
		"bucket_prefix": d.Get("bucket_prefix"),
		"agency_name":   utils.ValueIgnoreEmpty(d.Get("agency_name")),
	}
}

func updatePolicyDryRunConfiguration(ctx context.Context, client *golangsdk.ServiceClient, bodyParams map[string]interface{},
	timeout time.Duration) error {
	httpUrl := "v1/organizations/dry-run-config"
	updatePath := client.Endpoint + httpUrl
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(bodyParams),
	}
	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: policyDryRunConfigurationRefreshFunc(client, utils.PathSearch("root_id", bodyParams, "").(string),
			utils.PathSearch("policy_type", bodyParams, "").(string)),
		Timeout:      timeout,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	return err
}

func policyDryRunConfigurationRefreshFunc(client *golangsdk.ServiceClient, rootId, policyType string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetPolicyDryRunConfiguration(client, rootId, policyType)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("root.status", respBody, "").(string)
		if status == "enabled" || status == "disabled" {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func GetPolicyDryRunConfiguration(client *golangsdk.ServiceClient, rootId, policyType string) (interface{}, error) {
	httpUrl := "v1/organizations/dry-run-config"
	getPath := client.Endpoint + httpUrl
	getPath = fmt.Sprintf("%s?root_id=%s&policy_type=%s", getPath, rootId, policyType)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourcePolicyDryRunConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		rootId = d.Get("root_id").(string)
	)

	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	respBody, err := GetPolicyDryRunConfiguration(client, rootId, d.Get("policy_type").(string))
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			fmt.Sprintf("error retrieving policy dry-run configuration under root (%s)", rootId),
		)
	}

	mErr := multierror.Append(
		d.Set("root_id", utils.PathSearch("root.root_id", respBody, nil)),
		d.Set("policy_type", utils.PathSearch("root.policy_type", respBody, nil)),
		d.Set("status", utils.PathSearch("root.status", respBody, nil)),
		d.Set("bucket_name", utils.PathSearch("root.bucket_name", respBody, nil)),
		d.Set("region_id", utils.PathSearch("root.region_id", respBody, nil)),
		d.Set("bucket_prefix", utils.PathSearch("root.bucket_prefix", respBody, nil)),
		d.Set("agency_name", utils.PathSearch("root.agency_name", respBody, nil)),
		d.Set("created_at", utils.PathSearch("root.created_at", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("root.updated_at", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePolicyDryRunConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	err = updatePolicyDryRunConfiguration(ctx, client, buildPolicyDryRunConfigurationBodyParams(d, d.Get("status").(string)),
		d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.Errorf("error updating policy dry-run configuration: %s", err)
	}

	return resourcePolicyDryRunConfigurationRead(ctx, d, meta)
}

func resourcePolicyDryRunConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	err = updatePolicyDryRunConfiguration(ctx, client, buildPolicyDryRunConfigurationBodyParams(d, "disabled"),
		d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			fmt.Sprintf("error disabling policy dry-run configuration under root (%s)", d.Get("root_id").(string)),
		)
	}

	return nil
}

func resourcePolicyDryRunConfigurationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid resource ID format, want '<root_id>/<policy_type>', but got '%s'", importedId)
	}

	mErr := multierror.Append(
		d.Set("root_id", parts[0]),
		d.Set("policy_type", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
