package dsc

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var dscBatchAddDataMaskNonUpdatableParams = []string{
	"mask_strategies",
	"mask_strategies.*.name",
	"mask_strategies.*.algorithm",
	"mask_strategies.*.parameters",
	"data",
}

// @API DSC POST /v1/{project_id}/data/mask
func ResourceDscBatchAddDataMask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDscBatchAddDataMaskCreate,
		ReadContext:   resourceDscBatchAddDataMaskRead,
		UpdateContext: resourceDscBatchAddDataMaskUpdate,
		DeleteContext: resourceDscBatchAddDataMaskDelete,

		CustomizeDiff: config.FlexibleForceNew(dscBatchAddDataMaskNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"mask_strategies": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    100,
				Description: "Specifies the list of mask strategies, each corresponding to a field.",
				Elem:        dscBatchAddDataMaskStrategiesSchema(),
			},
			"data": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specifies the data list to be masked.",
				Elem:        &schema.Schema{Type: schema.TypeMap, Elem: &schema.Schema{Type: schema.TypeString}},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"masked_data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The list of masked data in JSON format.",
			},
		},
	}
}

func dscBatchAddDataMaskStrategiesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the field name to be masked.",
			},
			"algorithm": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the masking algorithm name.",
			},
			"parameters": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the masking algorithm parameters.",
			},
		},
	}
}

func resourceDscBatchAddDataMaskCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/data/mask"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildDscBatchAddDataMaskBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch adding DSC data mask: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	maskedDataRaw := utils.PathSearch("masked_data", respBody, make([]interface{}, 0))
	maskedDataJSON, _ := json.Marshal(maskedDataRaw)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("masked_data", string(maskedDataJSON)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDscBatchAddDataMaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	maskStrategies := make([]map[string]interface{}, 0)
	for _, v := range d.Get("mask_strategies").([]interface{}) {
		strategy := v.(map[string]interface{})
		item := map[string]interface{}{
			"name":      strategy["name"],
			"algorithm": strategy["algorithm"],
		}
		if params, ok := strategy["parameters"].(map[string]interface{}); ok && len(params) > 0 {
			item["parameters"] = params
		}
		maskStrategies = append(maskStrategies, item)
	}

	return map[string]interface{}{
		"mask_strategies": maskStrategies,
		"data":            d.Get("data"),
	}
}

func resourceDscBatchAddDataMaskRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceDscBatchAddDataMaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceDscBatchAddDataMaskDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch add DSC data mask.
Deleting this resource will not restore the masked data or undo the mask action, but will only
remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
