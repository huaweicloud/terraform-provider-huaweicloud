// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CC
// ---------------------------------------------------------------

package cc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CC GET /v3/{domain_id}/gcn/central-network/{central_network_id}/policies
// @API CC DELETE /v3/{domain_id}/gcn/central-network/{central_network_id}/policies/{policy_id}
// @API CC POST /v3/{domain_id}/gcn/central-network/{central_network_id}/policies/{policy_id}/apply
func ResourceCentralNetworkPolicyApply() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCentralNetworkPolicyApplyCreate,
		UpdateContext: resourceCentralNetworkPolicyApplyCreate,
		ReadContext:   resourceCentralNetworkPolicyApplyRead,
		DeleteContext: resourceCentralNetworkPolicyApplyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCentralNetworkPolicyApplyImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"central_network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Central network ID.`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Policy ID.`,
			},
		},
	}
}

func resourceCentralNetworkPolicyApplyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	centralNetworkId := d.Get("central_network_id").(string)

	err := applyCentralNetworkPolicy(ctx, d, meta, d.Get("policy_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(centralNetworkId)

	return resourceCentralNetworkPolicyApplyRead(ctx, d, meta)
}

func applyCentralNetworkPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}, policyId string) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		applyPolicyHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/policies/{policy_id}/apply"
		applyPolicyProduct = "cc"
	)
	applyPolicyClient, err := cfg.NewServiceClient(applyPolicyProduct, region)
	if err != nil {
		return fmt.Errorf("error creating CC client: %s", err)
	}

	applyPolicyPath := applyPolicyClient.Endpoint + applyPolicyHttpUrl
	applyPolicyPath = strings.ReplaceAll(applyPolicyPath, "{domain_id}", cfg.DomainID)
	applyPolicyPath = strings.ReplaceAll(applyPolicyPath, "{central_network_id}", d.Get("central_network_id").(string))
	applyPolicyPath = strings.ReplaceAll(applyPolicyPath, "{policy_id}", policyId)

	applyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = applyPolicyClient.Request("POST", applyPolicyPath, &applyOpt)
	if err != nil {
		return fmt.Errorf("error applying central network policy: %s", err)
	}

	err = centralNetworkPolicyApplyWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate), policyId)
	if err != nil {
		return fmt.Errorf("error waiting for the central network policy(%s) to be applied: %s", d.Id(), err)
	}

	return nil
}

func centralNetworkPolicyApplyWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{},
	t time.Duration, policyId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				createWaitingHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/policies"
				createWaitingProduct = "cc"
			)
			applyPolicyClient, err := cfg.NewServiceClient(createWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CC client: %s", err)
			}

			applyPolicyPath := applyPolicyClient.Endpoint + createWaitingHttpUrl
			applyPolicyPath = strings.ReplaceAll(applyPolicyPath, "{domain_id}", cfg.DomainID)
			applyPolicyPath = strings.ReplaceAll(applyPolicyPath, "{central_network_id}",
				d.Get("central_network_id").(string))

			applyPolicyOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json"},
			}

			applyPolicyResp, err := applyPolicyClient.Request("GET",
				applyPolicyPath, &applyPolicyOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			applyPolicyRespBody, err := utils.FlattenResponse(applyPolicyResp)
			if err != nil {
				return nil, "ERROR", err
			}
			jsonPath := fmt.Sprintf("central_network_policies[?id =='%s']|[0].is_applied", policyId)
			statusRaw := utils.PathSearch(jsonPath, applyPolicyRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `central_network_policies.is_applied`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			if status == "false" {
				return applyPolicyRespBody, "PENDING", nil
			}

			if status == "true" {
				return applyPolicyRespBody, "COMPLETED", nil
			}

			return applyPolicyRespBody, status, nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCentralNetworkPolicyApplyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getCentralNetworkPolicyApplyHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/policies"
		getCentralNetworkPolicyApplyProduct = "cc"
	)
	getCentralNetworkPolicyApplyClient, err := cfg.NewServiceClient(getCentralNetworkPolicyApplyProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	getCentralNetworkPolicyApplyPath := getCentralNetworkPolicyApplyClient.Endpoint + getCentralNetworkPolicyApplyHttpUrl
	getCentralNetworkPolicyApplyPath = strings.ReplaceAll(getCentralNetworkPolicyApplyPath, "{domain_id}", cfg.DomainID)
	getCentralNetworkPolicyApplyPath = strings.ReplaceAll(getCentralNetworkPolicyApplyPath, "{central_network_id}", d.Id())

	getCentralNetworkPolicyApplyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getCentralNetworkPolicyApplyResp, err := getCentralNetworkPolicyApplyClient.Request("GET", getCentralNetworkPolicyApplyPath,
		&getCentralNetworkPolicyApplyOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving central network policy")
	}

	getCentralNetworkPolicyApplyRespBody, err := utils.FlattenResponse(getCentralNetworkPolicyApplyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("central_network_policies[?id =='%s' && is_applied == `true`]|[0]", d.Get("policy_id").(string))
	getCentralNetworkPolicyApplyRespBody = utils.PathSearch(jsonPath, getCentralNetworkPolicyApplyRespBody, nil)
	if getCentralNetworkPolicyApplyRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no policy found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("central_network_id", utils.PathSearch("central_network_id", getCentralNetworkPolicyApplyRespBody, nil)),
		d.Set("policy_id", utils.PathSearch("id", getCentralNetworkPolicyApplyRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCentralNetworkPolicyApplyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// get the default policy id which is not associated with the er instance and the version is 1.
	var (
		getCentralNetworkPolicyApplyHttpUrl = "v3/{domain_id}/gcn/central-network/{central_network_id}/policies"
		getCentralNetworkPolicyApplyProduct = "cc"
	)
	getCentralNetworkPolicyApplyClient, err := cfg.NewServiceClient(getCentralNetworkPolicyApplyProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	getCentralNetworkPolicyApplyPath := getCentralNetworkPolicyApplyClient.Endpoint + getCentralNetworkPolicyApplyHttpUrl
	getCentralNetworkPolicyApplyPath = strings.ReplaceAll(getCentralNetworkPolicyApplyPath, "{domain_id}", cfg.DomainID)
	getCentralNetworkPolicyApplyPath = strings.ReplaceAll(getCentralNetworkPolicyApplyPath, "{central_network_id}",
		d.Get("central_network_id").(string))

	getCentralNetworkPolicyApplyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getCentralNetworkPolicyApplyResp, err := getCentralNetworkPolicyApplyClient.Request("GET", getCentralNetworkPolicyApplyPath,
		&getCentralNetworkPolicyApplyOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving central network policy")
	}

	getCentralNetworkPolicyApplyRespBody, err := utils.FlattenResponse(getCentralNetworkPolicyApplyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("central_network_policies[?id =='%s' && is_applied == `true`]|[0]", d.Get("policy_id").(string))
	appliedID := utils.PathSearch(jsonPath, getCentralNetworkPolicyApplyRespBody, nil)
	if appliedID == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no policy found")
	}

	defaultId := utils.PathSearch("central_network_policies[?version == `1`]|[0].id", getCentralNetworkPolicyApplyRespBody, nil)
	if defaultId == nil {
		return diag.Errorf("error applying central network policy to none: no default policy found")
	}

	// apply default policy
	err = applyCentralNetworkPolicy(ctx, d, meta, defaultId.(string))
	if err != nil {
		return diag.Errorf("error applying central network policy to none: %s", err)
	}
	return nil
}

func resourceCentralNetworkPolicyApplyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <central_network_id>/<policy_id>")
	}

	d.SetId(parts[0])
	d.Set("policy_id", parts[1])

	return []*schema.ResourceData{d}, nil
}
