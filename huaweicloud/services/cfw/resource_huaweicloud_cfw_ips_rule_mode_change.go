package cfw

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var ruleModeChangeNonUpdatableFields = []string{
	"object_id", "ips_ids", "status", "enterprise_project_id",
}

func ResourceCfwIpsRuleModeChange() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCfwIpsRuleModeChangeCreate,
		UpdateContext: resourceCfwIpsRuleModeChangeUpdate,
		ReadContext:   resourceCfwIpsRuleModeChangeRead,
		DeleteContext: resourceCfwIpsRuleModeChangeDelete,

		CustomizeDiff: config.FlexibleForceNew(ruleModeChangeNonUpdatableFields),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"object_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the protected object ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the IPS rule status.`,
			},
			"ips_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the IPS ID list.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID.`,
			},
		},
	}
}

type IpsRuleModeChangeRSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newIpsRuleModeChangeRSWrapper(d *schema.ResourceData, meta interface{}) *IpsRuleModeChangeRSWrapper {
	return &IpsRuleModeChangeRSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func resourceCfwIpsRuleModeChangeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newIpsRuleModeChangeRSWrapper(d, meta)
	_, err := wrapper.ChangeIpsRuleMode()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return resourceCfwIpsRuleModeChangeRead(ctx, d, meta)
}

// @API CFW POST /v1/{project_id}/ips-rule/mode
func (w *IpsRuleModeChangeRSWrapper) ChangeIpsRuleMode() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cfw")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/ips-rule/mode"
	params := map[string]any{
		"ips_ids":   w.ListToArray("ips_ids"),
		"object_id": w.Get("object_id"),
		"status":    w.Get("status"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("POST").
		URI(uri).
		Body(params).
		Request().
		Result()
}

func resourceCfwIpsRuleModeChangeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCfwIpsRuleModeChangeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCfwIpsRuleModeChangeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the API. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
