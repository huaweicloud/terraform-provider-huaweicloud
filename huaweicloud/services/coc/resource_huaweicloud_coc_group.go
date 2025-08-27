package coc

import (
	"context"
	"errors"
	"fmt"
	"log"
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

var groupNonUpdatableParams = []string{"component_id", "region_id", "vendor", "application_id"}

// @API COC POST /v1/groups
// @API COC PUT /v1/groups/{id}
// @API COC GET /v1/groups
// @API COC DELETE /v1/groups/{id}
// @API COC GET /v1/application-model/next
func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGroupImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(groupNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sync_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sync_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rule_tags": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"relation_configurations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameters": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	createHttpUrl := "v1/groups"
	createPath := client.Endpoint + createHttpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateGroupBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC group: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening group: %s", err)
	}

	id := utils.PathSearch("data.id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the COC group ID from the API response")
	}

	d.SetId(id)

	return resourceGroupRead(ctx, d, meta)
}

func buildCreateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":           d.Get("name"),
		"component_id":   d.Get("component_id"),
		"region_id":      d.Get("region_id"),
		"sync_mode":      d.Get("sync_mode"),
		"vendor":         utils.ValueIgnoreEmpty(d.Get("vendor")),
		"application_id": utils.ValueIgnoreEmpty(d.Get("application_id")),
		"sync_rules":     buildCreateGroupSyncRulesBodyParams(d.Get("sync_rules")),
		"relation_configurations": buildCreateGroupRelationConfigurationsBodyParams(
			d.Get("relation_configurations")),
	}

	return bodyParams
}

func buildCreateGroupSyncRulesBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"ep_id":     utils.ValueIgnoreEmpty(raw["enterprise_project_id"]),
				"rule_tags": utils.ValueIgnoreEmpty(raw["rule_tags"]),
			}
		}
		return params
	}

	return nil
}

func buildCreateGroupRelationConfigurationsBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"type":       utils.ValueIgnoreEmpty(raw["type"]),
				"parameters": utils.ValueIgnoreEmpty(raw["parameters"]),
			}
		}
		return params
	}

	return nil
}

func resourceGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	componentID := d.Get("component_id").(string)
	group, err := GetGroup(client, componentID, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving group")
	}
	groupSyncRules, err := GetGroupSyncRules(client, componentID, d.Id())
	if err != nil {
		log.Printf("[WARN] error fetching group sync rules (%s): %s", d.Id(), err)
	}

	mErr = multierror.Append(mErr,
		d.Set("name", utils.PathSearch("name", group, nil)),
		d.Set("vendor", utils.PathSearch("vendor", group, nil)),
		d.Set("code", utils.PathSearch("code", group, nil)),
		d.Set("region_id", utils.PathSearch("region_id", group, nil)),
		d.Set("component_id", utils.PathSearch("component_id", group, nil)),
		d.Set("application_id", utils.PathSearch("application_id", group, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("ep_id", group, nil)),
		d.Set("sync_mode", utils.PathSearch("sync_mode", group, nil)),
		d.Set("relation_configurations", flattenGroupRelationConfigurations(
			utils.PathSearch("relation_configurations", group, nil))),
		d.Set("sync_rules", flattenGroupSyncRules(groupSyncRules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGroupRelationConfigurations(rawParams interface{}) []map[string]interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) < 1 {
			return nil
		}
		configurations := make([]map[string]interface{}, len(paramsList))
		for i, params := range paramsList {
			raw := params.(map[string]interface{})
			configurations[i] = map[string]interface{}{
				"type":       utils.PathSearch("type", raw, nil),
				"parameters": utils.PathSearch("parameters", raw, nil),
			}
		}
		return configurations
	}

	return nil
}

func flattenGroupSyncRules(rawParams interface{}) []map[string]interface{} {
	if rawParams == nil {
		return nil
	}
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) < 1 {
			return nil
		}
		configurations := make([]map[string]interface{}, len(paramsList))
		for i, params := range paramsList {
			raw := params.(map[string]interface{})
			configurations[i] = map[string]interface{}{
				"enterprise_project_id": utils.PathSearch("ep_id", raw, nil),
				"rule_tags":             utils.PathSearch("rule_tags", raw, nil),
			}
		}
		return configurations
	}

	return nil
}

func GetGroup(client *golangsdk.ServiceClient, componentID string, groupID string) (interface{}, error) {
	getHttpUrl := "v1/groups?id_list={group_id}&component_id={component_id}&limit=1"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{group_id}", groupID)
	getPath = strings.ReplaceAll(getPath, "{component_id}", componentID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening group: %s", err)
	}

	group := utils.PathSearch("data|[0]", getRespBody, nil)
	if group == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return group, nil
}

func GetGroupSyncRules(client *golangsdk.ServiceClient, componentID string, groupID string) (interface{}, error) {
	getHttpUrl := "v1/application-model/next?component_id={component_id}&limit=100"
	basePath := client.Endpoint + getHttpUrl
	basePath = strings.ReplaceAll(basePath, "{component_id}", componentID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	var marker string
	for {
		getPath := basePath
		if marker != "" {
			getPath = fmt.Sprintf("%s&marker=%v", getPath, marker)
		}
		getResp, err := client.Request("GET", getPath, &getOpt)

		if err != nil {
			return nil, err
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}

		groups := utils.PathSearch("data.groups", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(groups) < 1 {
			return nil, golangsdk.ErrDefault404{}
		}

		searchGroupPath := fmt.Sprintf("data.groups[?id=='%s']|[0]", groupID)
		group := utils.PathSearch(searchGroupPath, getRespBody, nil)
		if group != nil {
			syncRules := utils.PathSearch("sync_rules", group, nil)
			return syncRules, nil
		}

		marker = utils.PathSearch("data.groups[-1].id", getRespBody, "").(string)
	}
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	changeList := []string{
		"name", "sync_mode", "sync_rules", "relation_configurations",
	}
	if d.HasChanges(changeList...) {
		updateHttpUrl := "v1/groups/{id}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateGroupBodyParams(d)),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating group: %s", err)
		}
	}

	return resourceGroupRead(ctx, d, meta)
}

func buildUpdateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":       d.Get("name"),
		"sync_mode":  d.Get("sync_mode"),
		"sync_rules": buildCreateGroupSyncRulesBodyParams(d.Get("sync_rules")),
		"relation_configurations": buildCreateGroupRelationConfigurationsBodyParams(
			d.Get("relation_configurations")),
	}

	return bodyParams
}

func resourceGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	deleteHttpUrl := "v1/groups/{id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"force_delete": utils.ValueIgnoreEmpty(d.Get("force_delete")),
		}),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code",
			"common.00000400"), "error deleting COC group")
	}

	return nil
}

func resourceGroupImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format of import ID, must be <component_id>/<group_id>")
	}
	d.SetId(parts[1])
	mErr := multierror.Append(
		nil,
		d.Set("component_id", parts[0]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set values in import state, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
