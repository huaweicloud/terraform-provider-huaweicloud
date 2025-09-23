package dbss

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DBSS POST /v1/{project_id}/{instance_id}/audit/rule/risk/switch
func ResourceRiskRuleAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRiskRuleActionCreate,
		ReadContext:   resourceRiskRuleActionRead,
		DeleteContext: resourceRiskRuleActionDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"risk_ids": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildRiskRuleActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ids":    d.Get("risk_ids"),
		"status": d.Get("action"),
	}
	return bodyParams
}

func resourceRiskRuleActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		httpUrl    = "v1/{project_id}/{instance_id}/audit/rule/risk/switch"
		mErr       *multierror.Error
	)

	client, err := cfg.NewServiceClient("dbss", region)
	if err != nil {
		return diag.Errorf("error creating DBSS client: %s", err)
	}

	riskRuleActionPath := client.Endpoint + httpUrl
	riskRuleActionPath = strings.ReplaceAll(riskRuleActionPath, "{project_id}", client.ProjectID)
	riskRuleActionPath = strings.ReplaceAll(riskRuleActionPath, "{instance_id}", instanceId)

	riskRuleActionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	riskRuleActionOpt.JSONBody = utils.RemoveNil(buildRiskRuleActionBodyParams(d))
	riskRuleActionResp, err := client.Request("POST", riskRuleActionPath, &riskRuleActionOpt)
	if err != nil {
		return diag.Errorf("error enabling or disabling risk rule: %s", err)
	}

	riskRuleActionRespBody, err := utils.FlattenResponse(riskRuleActionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceId)

	result := utils.PathSearch("status", riskRuleActionRespBody, "").(string)
	if result == "" {
		return diag.Errorf("err searching 'result' in API response: %s", err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("result", result),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRiskRuleActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a action resource.
	return nil
}

func resourceRiskRuleActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a action resource.
	return nil
}
