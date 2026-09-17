package drs

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

// @API DRS POST /v5/{project_id}/jobs/{job_id}/check-filter-rules
func ResourceCheckFilterRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCheckFilterRulesCreate,
		ReadContext:   resourceCheckFilterRulesRead,
		UpdateContext: resourceCheckFilterRulesUpdate,
		DeleteContext: resourceCheckFilterRulesDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_filter_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "common",
				ValidateFunc: validation.StringInSlice([]string{"common", "config"}, false),
			},
			"source": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "job",
				ValidateFunc: validation.StringInSlice([]string{"job", "compare"}, false),
			},
			"general_filtering_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_object_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"parent_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"object_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"object_alias_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"object_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"source_filter_condition": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"advanced_setting_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"col_names": {
							Type:     schema.TypeString,
							Required: true,
						},
						"prim_key_or_indexes": {
							Type:     schema.TypeString,
							Required: true,
						},
						"indexes": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"query_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCheckFilterRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"general_filtering_list": buildGeneralFilteringListBodyParams(d.Get("general_filtering_list")),
		"advanced_setting_list":  buildAdvancedSettingListBodyParams(d.Get("advanced_setting_list")),
		"source":                 utils.ValueIgnoreEmpty(d.Get("source")),
		"data_filter_type":       utils.ValueIgnoreEmpty(d.Get("data_filter_type")),
	}
}

func resourceCheckFilterRulesCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}/check-filter-rules"
	)

	client, err := cfg.DrsV5Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCheckFilterRulesBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error checking DRS data filter rules: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	queryId := utils.PathSearch("id", respBody, "").(string)
	if queryId == "" {
		return diag.Errorf("error checking DRS data filter rules: ID is not found in API response")
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)
	if err = d.Set("query_id", queryId); err != nil {
		return diag.Errorf("error setting query_id: %s", err)
	}

	return nil
}

func resourceCheckFilterRulesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceCheckFilterRulesUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceCheckFilterRulesDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to check the DRS data filter rules. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
