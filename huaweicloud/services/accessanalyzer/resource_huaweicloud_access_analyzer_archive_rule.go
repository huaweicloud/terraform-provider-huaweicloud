package accessanalyzer

import (
	"context"
	"fmt"
	"strconv"
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

// @API AccessAnalyzer POST /v5/analyzers/{analyzer_id}/archive-rules
// @API AccessAnalyzer GET /v5/analyzers/{analyzer_id}/archive-rules/{archive_rule_id}
// @API AccessAnalyzer PUT /v5/analyzers/{analyzer_id}/archive-rules/{archive_rule_id}
// @API AccessAnalyzer DELETE /v5/analyzers/{analyzer_id}/archive-rules/{archive_rule_id}
var nonUpdatableParamsArchiveRule = []string{"name", "analyzer_id"}

func ResourceArchiveRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchiveRuleCreate,
		ReadContext:   resourceArchiveRuleRead,
		UpdateContext: resourceArchiveRuleUpdate,
		DeleteContext: resourceArchiveRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceArchiveRuleImport,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsArchiveRule),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"analyzer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filters": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the cluster informations in the mesh.`,
						},
						"organization_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the cluster informations in the mesh.`,
						},
						"criterion": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: `Specifies the extend parameters of the mesh.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"contains": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `Specifies the cluster informations in the mesh.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"eq": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `Specifies the cluster informations in the mesh.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"neq": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `Specifies the cluster informations in the mesh.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"exists": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  `Specifies the cluster informations in the mesh.`,
										ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
									},
								},
							},
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
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArchiveRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("accessanalyzer", region)
	if err != nil {
		return diag.Errorf("error creating Access Analyzer client: %s", err)
	}

	createArchiveRuleHttpUrl := "v5/analyzers/{analyzer_id}/archive-rules"
	createArchiveRulePath := client.Endpoint + createArchiveRuleHttpUrl
	createArchiveRulePath = strings.ReplaceAll(createArchiveRulePath, "{analyzer_id}", d.Get("analyzer_id").(string))
	createArchiveRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateArchiveRuleBodyParams(d)),
	}
	createArchiveRuleResp, err := client.Request("POST", createArchiveRulePath, &createArchiveRuleOpt)
	if err != nil {
		return diag.Errorf("error creating Access Analyzer archive rule: %s", err)
	}
	createArchiveRuleRespBody, err := utils.FlattenResponse(createArchiveRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createArchiveRuleRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating archive rule: id is not found in API response")
	}
	d.SetId(id)

	return resourceArchiveRuleRead(ctx, d, meta)
}

func buildCreateArchiveRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":    d.Get("name"),
		"filters": buildCreateArchiveRuleFiltersBodyParams(d),
	}
	return bodyParams
}

func buildCreateArchiveRuleFiltersBodyParams(d *schema.ResourceData) []map[string]interface{} {
	filtersRaw := d.Get("filters").([]interface{})
	if len(filtersRaw) == 0 {
		return nil
	}

	bodyParams := make([]map[string]interface{}, len(filtersRaw))
	for i, v := range filtersRaw {
		filter := v.(map[string]interface{})
		bodyParams[i] = map[string]interface{}{
			"key":            filter["key"],
			"OrganizationId": utils.ValueIgnoreEmpty(filter["organization_id"]),
			"criterion":      buildCreateArchiveRuleCriterionBodyParams(filter["criterion"].([]interface{})),
		}
	}
	return bodyParams
}

func buildCreateArchiveRuleCriterionBodyParams(criterionRaw []interface{}) map[string]interface{} {
	if len(criterionRaw) == 0 || criterionRaw[0] == nil {
		return nil
	}

	criterion := criterionRaw[0].(map[string]interface{})
	bodyParams := map[string]interface{}{
		"contains": utils.ValueIgnoreEmpty(criterion["contains"]),
		"eq":       utils.ValueIgnoreEmpty(criterion["eq"]),
		"neq":      utils.ValueIgnoreEmpty(criterion["neq"]),
	}

	if criterion["exists"] != "" {
		bodyParams["exists"] = utils.StringToBool(criterion["exists"])
	}

	return bodyParams
}

func resourceArchiveRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("accessanalyzer", region)
	if err != nil {
		return diag.Errorf("error creating Access Analyzer client: %s", err)
	}

	getArchiveRuleHttpUrl := "v5/analyzers/{analyzer_id}/archive-rules/{archive_rule_id}"
	getArchiveRulePath := client.Endpoint + getArchiveRuleHttpUrl
	getArchiveRulePath = strings.ReplaceAll(getArchiveRulePath, "{analyzer_id}", d.Get("analyzer_id").(string))
	getArchiveRulePath = strings.ReplaceAll(getArchiveRulePath, "{archive_rule_id}", d.Id())
	getArchiveRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getArchiveRuleResp, err := client.Request("GET", getArchiveRulePath, &getArchiveRuleOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving archive rule")
	}
	getArchiveRuleRespBody, err := utils.FlattenResponse(getArchiveRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	archiveRule := utils.PathSearch("archive_rule", getArchiveRuleRespBody, nil)
	if archiveRule == nil {
		return diag.Errorf("error getting archive rule: archive_rule is not found in API response")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", archiveRule, nil)),
		d.Set("urn", utils.PathSearch("urn", archiveRule, nil)),
		d.Set("created_at", utils.PathSearch("created_at", archiveRule, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", archiveRule, nil)),
		d.Set("filters", flattenFilters(utils.PathSearch("filters", archiveRule, nil))),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting archive rule fields: %s", err)
	}

	return nil
}

func flattenFilters(filtersRaw interface{}) []map[string]interface{} {
	if filtersRaw == nil {
		return nil
	}

	filters := filtersRaw.([]interface{})
	res := make([]map[string]interface{}, len(filters))
	for i, v := range filters {
		filter := v.(map[string]interface{})
		res[i] = map[string]interface{}{
			"key":             filter["key"],
			"organization_id": filter["OrganizationId"],
			"criterion":       flattenCriterion(filter["criterion"]),
		}
	}

	return res
}

func flattenCriterion(criterionRaw interface{}) []map[string]interface{} {
	if criterionRaw == nil {
		return nil
	}

	criterion := criterionRaw.(map[string]interface{})

	var exists string
	if v, ok := criterion["exists"]; ok {
		exists = strconv.FormatBool(v.(bool))
	}

	res := []map[string]interface{}{
		{
			"contains": criterion["contains"],
			"eq":       criterion["eq"],
			"neq":      criterion["neq"],
			"exists":   exists,
		},
	}

	return res
}

func resourceArchiveRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("accessanalyzer", region)
	if err != nil {
		return diag.Errorf("error creating Access Analyzer client: %s", err)
	}

	updateArchiveRuleHttpUrl := "v5/analyzers/{analyzer_id}/archive-rules/{archive_rule_id}"
	updateArchiveRulePath := client.Endpoint + updateArchiveRuleHttpUrl
	updateArchiveRulePath = strings.ReplaceAll(updateArchiveRulePath, "{analyzer_id}", d.Get("analyzer_id").(string))
	updateArchiveRulePath = strings.ReplaceAll(updateArchiveRulePath, "{archive_rule_id}", d.Id())
	updateArchiveRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateArchiveRuleBodyParams(d)),
	}
	_, err = client.Request("PUT", updateArchiveRulePath, &updateArchiveRuleOpt)
	if err != nil {
		return diag.Errorf("error updating Access Analyzer archive rule: %s", err)
	}

	return resourceArchiveRuleRead(ctx, d, meta)
}

func buildUpdateArchiveRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"filters": buildCreateArchiveRuleFiltersBodyParams(d),
	}
	return bodyParams
}

func resourceArchiveRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("accessanalyzer", region)
	if err != nil {
		return diag.Errorf("error creating Access Analyzer client: %s", err)
	}

	deleteArchiveRuleHttpUrl := "v5/analyzers/{analyzer_id}/archive-rules/{archive_rule_id}"
	deleteArchiveRulePath := client.Endpoint + deleteArchiveRuleHttpUrl
	deleteArchiveRulePath = strings.ReplaceAll(deleteArchiveRulePath, "{analyzer_id}", d.Get("analyzer_id").(string))
	deleteArchiveRulePath = strings.ReplaceAll(deleteArchiveRulePath, "{archive_rule_id}", d.Id())
	deleteArchiveRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deleteArchiveRulePath, &deleteArchiveRuleOpt)
	if err != nil {
		return diag.Errorf("error deleting archive rule: %s", err)
	}

	return nil
}

func resourceArchiveRuleImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for archive rule. Format must be <analyzer_id>/<id>")
		return nil, err
	}

	analyzerID := parts[0]
	id := parts[1]

	d.SetId(id)
	d.Set("analyzer_id", analyzerID)

	return []*schema.ResourceData{d}, nil
}
