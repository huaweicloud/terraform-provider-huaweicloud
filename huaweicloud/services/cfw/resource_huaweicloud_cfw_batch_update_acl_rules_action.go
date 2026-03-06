package cfw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var batchUpdateAclRulesActionNonUpdatableParams = []string{
	"object_id", "rule_ids", "action_type", "fw_instance_id", "enterprise_project_id"}

// @API CFW PUT /v1/{project_id}/acl-rule/action
func ResourceBatchUpdateAclRulesAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchUpdateAclRulesActionCreate,
		ReadContext:   resourceBatchUpdateAclRulesActionRead,
		UpdateContext: resourceBatchUpdateAclRulesActionUpdate,
		DeleteContext: resourceBatchUpdateAclRulesActionDelete,

		CustomizeDiff: config.FlexibleForceNew(batchUpdateAclRulesActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildBatchUpdateAclRulesActionQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := ""

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("fw_instance_id"); ok {
		queryParams = fmt.Sprintf("%s&fw_instance_id=%v", queryParams, v)
	}
	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func buildBatchUpdateAclRulesActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"object_id": d.Get("object_id"),
		"action":    d.Get("action"),
		"rule_ids":  utils.ExpandToStringList(d.Get("rule_ids").([]interface{})),
	}
}

func resourceBatchUpdateAclRulesActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		objectId = d.Get("object_id").(string)
		epsId    = cfg.GetEnterpriseProjectID(d)
		httpUrl  = "v1/{project_id}/acl-rule/action"
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildBatchUpdateAclRulesActionQueryParams(d, epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildBatchUpdateAclRulesActionBodyParams(d),
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch updating CFW ACL rules action: %s", err)
	}

	d.SetId(objectId)

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(d.Set("data", flattenBatchUpdateAclRulesActionDataResp(respBody)))
}

func flattenBatchUpdateAclRulesActionDataResp(respBody interface{}) []string {
	dataResp := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
	if len(dataResp) == 0 {
		return nil
	}

	return utils.ExpandToStringList(dataResp)
}

func resourceBatchUpdateAclRulesActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchUpdateAclRulesActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchUpdateAclRulesActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch update ACL rules action. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
