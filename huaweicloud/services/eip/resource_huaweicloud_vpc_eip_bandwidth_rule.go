package eip

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eipBandwidthRuleNonUpdatableParams = []string{
	"bandwidth_id",
	"name",
	"egress_size",
	"egress_guarented_size",
	"description",
	"publicip_info",
	"publicip_info.*.publicip_id",
}

// Currently, only the API creation and deployment are complete.
// We're doing a one-time resource allocation; further improvements will be needed
// after the three subsequent APIs are deployed.
// @API EIP POST /v3/{project_id}/eip/bandwidths/{bandwidth_id}/bandwidth-rules
func ResourceEipBandwidthRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEipBandwidthRuleCreate,
		ReadContext:   resourceEipBandwidthRuleRead,
		UpdateContext: resourceEipBandwidthRuleUpdate,
		DeleteContext: resourceEipBandwidthRuleDelete,

		CustomizeDiff: config.FlexibleForceNew(eipBandwidthRuleNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bandwidth_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"egress_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"egress_guarented_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"publicip_info": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"publicip_id": {
							Type:     schema.TypeString,
							Required: true,
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

func buildEipBandwidthRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	publicIpInfoRaw := d.Get("publicip_info").([]interface{})

	var publicIpInfoList []map[string]interface{}
	for _, v := range publicIpInfoRaw {
		if item, ok := v.(map[string]interface{}); ok {
			publicIpInfoList = append(publicIpInfoList, map[string]interface{}{
				"publicip_id": item["publicip_id"],
			})
		}
	}

	bodyParams := map[string]interface{}{
		"bandwidth_rule": map[string]interface{}{
			"name":                  d.Get("name"),
			"egress_size":           d.Get("egress_size"),
			"egress_guarented_size": d.Get("egress_guarented_size"),
			"description":           utils.ValueIgnoreEmpty(d.Get("description")),
			"publicip_info":         publicIpInfoList,
		},
	}
	return bodyParams
}

func resourceEipBandwidthRuleCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "vpc"
		bandwidthId = d.Get("bandwidth_id").(string)
		httpUrl     = "v3/{project_id}/eip/bandwidths/{bandwidth_id}/bandwidth-rules"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	requestId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate request ID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{bandwidth_id}", bandwidthId)

	opt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type":      "application/json",
			"Client-Request-Id": requestId,
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildEipBandwidthRuleBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &opt)
	if err != nil {
		return diag.Errorf("error creating EIP bandwidth rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("bandwidth_rule.id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("unable to find the bandwidth rule ID from the API response")
	}

	d.SetId(ruleId)

	return nil
}

func resourceEipBandwidthRuleRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Read()' method because resource is a one-time action resource.
	return nil
}

func resourceEipBandwidthRuleUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Update()' method because resource is a one-time action resource.
	return nil
}

func resourceEipBandwidthRuleDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to create a bandwidth rule for shared bandwidth. 
Please enable the enterprise-level QoS feature before creating the bandwidth rule.
Deleting this resource will not delete the actual bandwidth rule on the cloud, but will only remove the resource 
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
