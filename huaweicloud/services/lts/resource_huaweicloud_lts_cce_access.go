package lts

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS DELETE /v3/{project_id}/lts/access-config
// @API LTS POST /v3/{project_id}/lts/access-config
// @API LTS PUT /v3/{project_id}/lts/access-config
// @API LTS POST /v3/{project_id}/lts/access-config-list
func ResourceCceAccessConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCceAccessConfigCreate,
		UpdateContext: resourceCceAccessConfigUpdate,
		ReadContext:   resourceCceAccessConfigRead,
		DeleteContext: resourceHostAccessConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: cceAccessConfigResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the CCE access.`,
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The log group ID.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The log stream ID.`,
			},
			"access_config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem:        cceAccessConfigDeatilSchema(),
				Description: `The configurations of CCE access.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The CCE cluster ID.`,
			},
			"host_group_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `The log access host group ID list.`,
			},
			"tags": common.TagsSchema(),
			"binary_collect": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether collect in binary format. Default is false.`,
			},
			"log_split": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to split log. Default is false.`,
			},
			"access_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The log access type.`,
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The log group name.`,
			},
			"log_stream_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The log stream name.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the CCE access, in RFC3339 format.`,
			},
		},
	}
}

func cceAccessConfigDeatilSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"path_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the CCE access.`,
			},
			"paths": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `The collection paths.`,
			},
			"black_paths": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `The collection path blacklist.`,
			},
			"windows_log_info": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        hostAccessConfigWindowsLogInfoSchema(),
				Optional:    true,
				Description: `The configuration of Windows event logs.`,
			},
			"single_log_format": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				Elem:         hostAccessConfigSingleLogFormatSchema(),
				ExactlyOneOf: []string{"access_config.0.multi_log_format"},
				Description:  `The configuration single-line logs.`,
			},
			"multi_log_format": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem:        hostAccessConfigMultiLogFormatSchema(),
				Description: `The configuration multi-line logs.`,
			},
			"stdout": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether output is standard. Default is false.`,
			},
			"stderr": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether error output is standard. Default is false.`,
			},
			"name_space_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The regular expression matching of kubernetes namespaces.`,
			},
			"pod_name_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The regular expression matching of kubernetes pods.`,
			},
			"container_name_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The regular expression matching of kubernetes container names.`,
			},
			"log_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The container label log tag.`,
			},
			"include_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The container label whitelist.`,
			},
			"exclude_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The container label blacklist.`,
			},
			"log_envs": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The environment variable tag.`,
			},
			"include_envs": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The environment variable whitelist.`,
			},
			"exclude_envs": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The environment variable blacklist.`,
			},
			"log_k8s": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The kubernetes label log tag.`,
			},
			"include_k8s_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The kubernetes label whitelist.`,
			},
			"exclude_k8s_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The kubernetes label blacklist.`,
			},
		},
	}
	return &sc
}

func resourceCceAccessConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                          = meta.(*config.Config)
		region                       = cfg.GetRegion(d)
		createCceAccessConfigHttpUrl = "v3/{project_id}/lts/access-config"
		createCceAccessConfigProduct = "lts"
	)
	ltsClient, err := cfg.NewServiceClient(createCceAccessConfigProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	createCceAccessConfigPath := ltsClient.Endpoint + createCceAccessConfigHttpUrl
	createCceAccessConfigPath = strings.ReplaceAll(createCceAccessConfigPath, "{project_id}", ltsClient.ProjectID)

	createCceAccessConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createCceAccessConfigOpt.JSONBody = utils.RemoveNil(buildCreateCceAccessConfigBodyParams(d))
	createCceAccessConfigResp, err := ltsClient.Request("POST", createCceAccessConfigPath, &createCceAccessConfigOpt)
	if err != nil {
		return diag.Errorf("error creating CCE access config: %s", err)
	}

	createCceAccessConfigRespBody, err := utils.FlattenResponse(createCceAccessConfigResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("access_config_id", createCceAccessConfigRespBody)
	if err != nil {
		return diag.Errorf("error creating CCE access config: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceCceAccessConfigRead(ctx, d, meta)
}

func buildCreateCceAccessConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	logInfoOpts := map[string]interface{}{
		"log_group_id":  d.Get("log_group_id"),
		"log_stream_id": d.Get("log_stream_id"),
	}

	bodyParams := map[string]interface{}{
		"access_config_type":   "K8S_CCE",
		"access_config_name":   d.Get("name"),
		"access_config_detail": buildCceAccessConfigDeatilRequestBody(d.Get("access_config")),
		"log_info":             logInfoOpts,
		"host_group_info":      buildHostGroupInfoRequestBody(d),
		"access_config_tag":    utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		"cluster_id":           d.Get("cluster_id"),
		"binary_collect":       utils.ValueIgnoreEmpty(d.Get("binary_collect")),
		"log_split":            utils.ValueIgnoreEmpty(d.Get("log_split")),
	}
	return bodyParams
}

func buildCceAccessConfigDeatilRequestBody(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok || len(rawArray) < 1 {
		log.Printf("[WARN] access configuration is empty or the type is not a standard '[]interface{}'")
		return nil
	}

	raw, ok := rawArray[0].(map[string]interface{})
	if !ok {
		log.Printf("[WARN] access configuration sub node type is not a standard 'map[string]interface{}'")
		return nil
	}

	params := map[string]interface{}{
		"pathType":           raw["path_type"],
		"paths":              utils.ValueIgnoreEmpty(raw["paths"]),
		"black_paths":        utils.ValueIgnoreEmpty(raw["black_paths"]),
		"format":             buildHostAccessConfigFormatRequestBody(raw),
		"windows_log_info":   buildHostAccessConfigWindowsLogInfoRequestBody(raw["windows_log_info"]),
		"stdout":             utils.ValueIgnoreEmpty(raw["stdout"]),
		"stderr":             utils.ValueIgnoreEmpty(raw["stderr"]),
		"namespaceRegex":     utils.ValueIgnoreEmpty(raw["name_space_regex"]),
		"podNameRegex":       utils.ValueIgnoreEmpty(raw["pod_name_regex"]),
		"containerNameRegex": utils.ValueIgnoreEmpty(raw["container_name_regex"]),
		"logLabels":          raw["log_labels"],
		"includeLabels":      raw["include_labels"],
		"excludeLabels":      raw["exclude_labels"],
		"logEnvs":            raw["log_envs"],
		"includeEnvs":        raw["include_envs"],
		"excludeEnvs":        raw["exclude_envs"],
		"logK8s":             raw["log_k8s"],
		"includeK8sLabels":   raw["include_k8s_labels"],
		"excludeK8sLabels":   raw["exclude_k8s_labels"],
	}
	return params
}

func resourceCceAccessConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                        = meta.(*config.Config)
		region                     = cfg.GetRegion(d)
		listCceAccessConfigHttpUrl = "v3/{project_id}/lts/access-config-list"
		listCceAccessConfigProduct = "lts"
	)
	ltsClient, err := cfg.NewServiceClient(listCceAccessConfigProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	listCceAccessConfigPath := ltsClient.Endpoint + listCceAccessConfigHttpUrl
	listCceAccessConfigPath = strings.ReplaceAll(listCceAccessConfigPath, "{project_id}", ltsClient.ProjectID)

	listCceAccessConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	name := d.Get("name").(string)
	listCceAccessConfigOpt.JSONBody = map[string]interface{}{
		"access_config_name_list": []string{name},
	}
	listCceAccessConfigResp, err := ltsClient.Request("POST", listCceAccessConfigPath, &listCceAccessConfigOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE access config")
	}

	listCceAccessConfigRespBody, err := utils.FlattenResponse(listCceAccessConfigResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("result[?access_config_name=='%s']|[0]", name)
	listCceAccessConfigRespBody = utils.PathSearch(jsonPath, listCceAccessConfigRespBody, nil)
	if listCceAccessConfigRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "Cce access config")
	}
	created := utils.PathSearch("create_time", listCceAccessConfigRespBody, float64(0)).(float64)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("access_config_name", listCceAccessConfigRespBody, nil)),
		d.Set("access_type", utils.PathSearch("access_config_type", listCceAccessConfigRespBody, nil)),
		d.Set("log_group_id", utils.PathSearch("log_info.log_group_id", listCceAccessConfigRespBody, nil)),
		d.Set("log_stream_id", utils.PathSearch("log_info.log_stream_id", listCceAccessConfigRespBody, nil)),
		d.Set("log_group_name", utils.PathSearch("log_info.log_group_name", listCceAccessConfigRespBody, nil)),
		d.Set("log_stream_name", utils.PathSearch("log_info.log_stream_name", listCceAccessConfigRespBody, nil)),
		d.Set("host_group_ids", utils.PathSearch("host_group_info.host_group_id_list", listCceAccessConfigRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("access_config_tag", listCceAccessConfigRespBody, nil))),
		d.Set("access_config", flattenCceAccessConfigDetail(listCceAccessConfigRespBody)),
		d.Set("cluster_id", utils.PathSearch("cluster_id", listCceAccessConfigRespBody, nil)),
		d.Set("binary_collect", utils.PathSearch("binary_collect", listCceAccessConfigRespBody, nil)),
		d.Set("log_split", utils.PathSearch("log_split", listCceAccessConfigRespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(created)/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCceAccessConfigDetail(resp interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"path_type":            utils.PathSearch("access_config_detail.pathType", resp, nil),
			"paths":                utils.PathSearch("access_config_detail.paths", resp, nil),
			"black_paths":          utils.PathSearch("access_config_detail.black_paths", resp, nil),
			"single_log_format":    flattenHostAccessConfigLogFormat(utils.PathSearch("access_config_detail.format.single", resp, nil)),
			"multi_log_format":     flattenHostAccessConfigLogFormat(utils.PathSearch("access_config_detail.format.multi", resp, nil)),
			"windows_log_info":     flattenHostAccessConfigWindowsLogInfo(utils.PathSearch("access_config_detail.windows_log_info", resp, nil)),
			"stdout":               utils.PathSearch("access_config_detail.stdout", resp, nil),
			"stderr":               utils.PathSearch("access_config_detail.stderr", resp, nil),
			"name_space_regex":     utils.PathSearch("access_config_detail.namespaceRegex", resp, nil),
			"pod_name_regex":       utils.PathSearch("access_config_detail.podNameRegex", resp, nil),
			"container_name_regex": utils.PathSearch("access_config_detail.containerNameRegex", resp, nil),
			"log_labels":           utils.PathSearch("access_config_detail.logLabels", resp, nil),
			"include_labels":       utils.PathSearch("access_config_detail.includeLabels", resp, nil),
			"exclude_labels":       utils.PathSearch("access_config_detail.excludeLabels", resp, nil),
			"log_envs":             utils.PathSearch("access_config_detail.logEnvs", resp, nil),
			"include_envs":         utils.PathSearch("access_config_detail.includeEnvs", resp, nil),
			"exclude_envs":         utils.PathSearch("access_config_detail.excludeEnvs", resp, nil),
			"log_k8s":              utils.PathSearch("access_config_detail.logK8s", resp, nil),
			"include_k8s_labels":   utils.PathSearch("access_config_detail.includeK8sLabels", resp, nil),
			"exclude_k8s_labels":   utils.PathSearch("access_config_detail.excludeK8sLabels", resp, nil),
		},
	}
}

func resourceCceAccessConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	updateCceAccessConfigChanges := []string{
		"access_config",
		"host_group_ids",
		"tags",
		"binary_collect",
		"log_split",
	}

	if d.HasChanges(updateCceAccessConfigChanges...) {
		var (
			cfg                          = meta.(*config.Config)
			region                       = cfg.GetRegion(d)
			updateCceAccessConfigHttpUrl = "v3/{project_id}/lts/access-config"
			updateCceAccessConfigProduct = "lts"
		)
		ltsClient, err := cfg.NewServiceClient(updateCceAccessConfigProduct, region)
		if err != nil {
			return diag.Errorf("error creating LTS client: %s", err)
		}

		updateCceAccessConfigPath := ltsClient.Endpoint + updateCceAccessConfigHttpUrl
		updateCceAccessConfigPath = strings.ReplaceAll(updateCceAccessConfigPath, "{project_id}", ltsClient.ProjectID)

		updateCceAccessConfigOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		updateCceAccessConfigOpt.JSONBody = utils.RemoveNil(buildUpdateCceAccessConfigBodyParams(d))
		_, err = ltsClient.Request("PUT", updateCceAccessConfigPath, &updateCceAccessConfigOpt)
		if err != nil {
			return diag.Errorf("error updating CCE access config: %s", err)
		}
	}
	return resourceCceAccessConfigRead(ctx, d, meta)
}

func buildUpdateCceAccessConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"access_config_id":     d.Id(),
		"access_config_detail": buildCceAccessConfigDeatilRequestBody(d.Get("access_config")),
		"host_group_info":      buildHostGroupInfoRequestBody(d),
		"access_config_tag":    utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		"cluster_id":           d.Get("cluster_id"),
		"binary_collect":       utils.ValueIgnoreEmpty(d.Get("binary_collect")),
		"log_split":            utils.ValueIgnoreEmpty(d.Get("log_split")),
	}
	return bodyParams
}

func cceAccessConfigResourceImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var (
		cfg                        = meta.(*config.Config)
		region                     = cfg.GetRegion(d)
		listCceAccessConfigProduct = "lts"
		name                       = d.Id()
	)
	d.Set("name", name)
	ltsClient, err := cfg.NewServiceClient(listCceAccessConfigProduct, region)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating LTS client: %s", err)
	}

	return []*schema.ResourceData{d}, refreshCceAccessID(ltsClient, d)
}

func refreshCceAccessID(ltsClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		listCceAccessConfigHttpUrl = "v3/{project_id}/lts/access-config-list"
	)
	listCceAccessConfigPath := ltsClient.Endpoint + listCceAccessConfigHttpUrl
	listCceAccessConfigPath = strings.ReplaceAll(listCceAccessConfigPath, "{project_id}", ltsClient.ProjectID)

	listCceAccessConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	name := d.Get("name").(string)
	listCceAccessConfigOpt.JSONBody = map[string]interface{}{
		"access_config_name_list": []string{name},
	}
	listCceAccessConfigResp, err := ltsClient.Request("POST", listCceAccessConfigPath, &listCceAccessConfigOpt)
	if err != nil {
		return fmt.Errorf("error retrieving CCE access config: %s", err)
	}

	listCceAccessConfigRespBody, err := utils.FlattenResponse(listCceAccessConfigResp)
	if err != nil {
		return err
	}

	jsonPath := fmt.Sprintf("result[?access_config_name=='%s']|[0]", name)
	listCceAccessConfigRespBody = utils.PathSearch(jsonPath, listCceAccessConfigRespBody, nil)
	if listCceAccessConfigRespBody == nil {
		return fmt.Errorf("the CCE access config (%s) does not exist", name)
	}

	configID := utils.PathSearch("access_config_id", listCceAccessConfigRespBody, "")
	if configID == "" {
		return fmt.Errorf("error retrieving CCE access config: ID is not found in API response")
	}

	d.SetId(configID.(string))

	return nil
}
