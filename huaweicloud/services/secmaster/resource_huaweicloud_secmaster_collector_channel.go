package secmaster

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var collectorChannelNonUpdatableParams = []string{"workspace_id"}

var collectorChannelActions = []string{"START", "STOP", "REMOVE", "RESTART", "REFRESH", "INSTALL"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/collector/channels
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/collector/channels
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/collector/channels/{channel_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/collector/channels/{channel_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/collector/channels/{channel_id}
func ResourceCollectorChannel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCollectorChannelCreate,
		ReadContext:   resourceCollectorChannelRead,
		UpdateContext: resourceCollectorChannelUpdate,
		DeleteContext: resourceCollectorChannelDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCollectorChannelImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(collectorChannelNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parser_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"input": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     collectorChannelModuleSchema(),
			},
			"output": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     collectorChannelModuleSchema(),
			},
			"nodes": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     collectorChannelNodeSchema(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(collectorChannelActions, false),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"create_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operation_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parser_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func collectorChannelNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"args": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func collectorChannelModuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_module_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"children": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     collectorChannelChildModuleSchema(),
			},
			"fields": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     collectorChannelFieldSchema(),
			},
		},
	}
}

func collectorChannelChildModuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_module_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fields": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     collectorChannelFieldSchema(),
			},
		},
	}
}

func collectorChannelFieldSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"other": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_field_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_module_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func buildCollectorChannelBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"title":       d.Get("title"),
		"group_id":    d.Get("group_id"),
		"parser_id":   d.Get("parser_id"),
		"input":       buildCollectorChannelModulesBodyParams(d.Get("input").([]interface{})),
		"output":      buildCollectorChannelModulesBodyParams(d.Get("output").([]interface{})),
		"nodes":       buildCollectorChannelNodesBodyParams(d.Get("nodes").([]interface{})),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"action":      utils.ValueIgnoreEmpty(d.Get("action")),
	}

	return bodyParams
}

func buildCollectorChannelNodesBodyParams(nodes []interface{}) []map[string]interface{} {
	if len(nodes) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(nodes))
	for _, v := range nodes {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		node := map[string]interface{}{
			"node_id":     utils.ValueIgnoreEmpty(rawMap["node_id"]),
			"node_status": utils.ValueIgnoreEmpty(rawMap["node_status"]),
			"args":        buildCollectorChannelArgsBodyParams(rawMap["args"].([]interface{})),
		}
		rst = append(rst, node)
	}

	return rst
}

func buildCollectorChannelArgsBodyParams(args []interface{}) []map[string]interface{} {
	if len(args) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(args))
	for _, v := range args {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		arg := map[string]interface{}{
			"key":   utils.ValueIgnoreEmpty(rawMap["key"]),
			"value": utils.ValueIgnoreEmpty(rawMap["value"]),
		}
		rst = append(rst, arg)
	}

	return rst
}

func buildCollectorChannelModulesBodyParams(modules []interface{}) []map[string]interface{} {
	if len(modules) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(modules))
	for _, v := range modules {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		module := map[string]interface{}{
			"name":                 utils.ValueIgnoreEmpty(rawMap["name"]),
			"template_id":          utils.ValueIgnoreEmpty(rawMap["template_id"]),
			"connection_module_id": utils.ValueIgnoreEmpty(rawMap["connection_module_id"]),
			"children":             buildCollectorChannelChildModulesBodyParams(rawMap["children"].([]interface{})),
			"fields":               buildCollectorChannelFieldsBodyParams(rawMap["fields"].([]interface{})),
		}
		rst = append(rst, module)
	}

	return rst
}

func buildCollectorChannelChildModulesBodyParams(modules []interface{}) []map[string]interface{} {
	if len(modules) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(modules))
	for _, v := range modules {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		module := map[string]interface{}{
			"name":                 utils.ValueIgnoreEmpty(rawMap["name"]),
			"template_id":          utils.ValueIgnoreEmpty(rawMap["template_id"]),
			"connection_module_id": utils.ValueIgnoreEmpty(rawMap["connection_module_id"]),
			"fields":               buildCollectorChannelFieldsBodyParams(rawMap["fields"].([]interface{})),
		}
		rst = append(rst, module)
	}

	return rst
}

func buildCollectorChannelFieldsBodyParams(fields []interface{}) []map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(fields))
	for _, v := range fields {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		field := map[string]interface{}{
			"name":                 utils.ValueIgnoreEmpty(rawMap["name"]),
			"type":                 utils.ValueIgnoreEmpty(rawMap["type"]),
			"value":                utils.ValueIgnoreEmpty(rawMap["value"]),
			"other":                utils.ValueIgnoreEmpty(rawMap["other"]),
			"template_field_id":    utils.ValueIgnoreEmpty(rawMap["template_field_id"]),
			"connection_module_id": utils.ValueIgnoreEmpty(rawMap["connection_module_id"]),
		}
		rst = append(rst, field)
	}

	return rst
}

func resourceCollectorChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		title       = d.Get("title").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/collector/channels"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCollectorChannelBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating collector channel: %s", err)
	}

	channel, err := GetCollectorChannelByTitle(client, workspaceId, title)
	if err != nil {
		return diag.FromErr(err)
	}

	channelId := utils.PathSearch("channel_id", channel, "").(string)
	if channelId == "" {
		return diag.Errorf("error creating collector channel: unable to find channel ID")
	}

	d.SetId(channelId)

	return resourceCollectorChannelRead(ctx, d, meta)
}

func GetCollectorChannelByTitle(client *golangsdk.ServiceClient, workspaceId, title string) (interface{}, error) {
	httpUrl := "v1/{project_id}/workspaces/{workspace_id}/collector/channels"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{workspace_id}", workspaceId)
	queryParams := url.Values{}
	queryParams.Add("title", title)
	queryParams.Add("limit", "200")
	queryParams.Add("offset", "0")
	listPath = fmt.Sprintf("%s?%s", listPath, queryParams.Encode())

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", listPath, &listOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	record := utils.PathSearch("records|[0]", respBody, nil)
	if record == nil {
		// When the collector channel does not exist, the status code is `200` and the records list is empty.
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/workspaces/{workspace_id}/collector/channels",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the collector channel (%s) does not exist", title)),
			},
		}
	}
	return record, nil
}

func GetCollectorChannelById(client *golangsdk.ServiceClient, workspaceId, channelId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/workspaces/{workspace_id}/collector/channels/{channel_id}"
	readPath := client.Endpoint + httpUrl
	readPath = strings.ReplaceAll(readPath, "{project_id}", client.ProjectID)
	readPath = strings.ReplaceAll(readPath, "{workspace_id}", workspaceId)
	readPath = strings.ReplaceAll(readPath, "{channel_id}", channelId)

	readOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", readPath, &readOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceCollectorChannelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	respBody, err := GetCollectorChannelById(client, workspaceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving collector channel")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workspace_id", workspaceId),
		d.Set("title", utils.PathSearch("title", respBody, nil)),
		d.Set("group_id", utils.PathSearch("group_id", respBody, nil)),
		d.Set("parser_id", utils.PathSearch("parser_id", respBody, nil)),
		d.Set("input", flattenCollectorChannelModules(utils.PathSearch("input", respBody, nil))),
		d.Set("output", flattenCollectorChannelModules(utils.PathSearch("output", respBody, nil))),
		d.Set("nodes", flattenCollectorChannelNodes(utils.PathSearch("nodes", respBody, nil))),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("create_by", utils.PathSearch("create_by", respBody, nil)),
		d.Set("error", utils.PathSearch("error", respBody, nil)),
		d.Set("operation_status", utils.PathSearch("operation_status", respBody, nil)),
		d.Set("parser_name", utils.PathSearch("parser_name", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCollectorChannelNodes(nodesResp interface{}) []interface{} {
	if nodesResp == nil {
		return nil
	}

	nodes, ok := nodesResp.([]interface{})
	if !ok || len(nodes) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(nodes))
	for _, v := range nodes {
		rst = append(rst, map[string]interface{}{
			"node_id":     utils.ValueIgnoreEmpty(utils.PathSearch("node_id", v, "").(string)),
			"node_status": utils.PathSearch("node_status", v, nil),
			"args":        flattenCollectorChannelArgs(utils.PathSearch("args", v, nil)),
		})
	}

	return rst
}

func flattenCollectorChannelArgs(argsResp interface{}) []interface{} {
	if argsResp == nil {
		return nil
	}

	args, ok := argsResp.([]interface{})
	if !ok || len(args) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(args))
	for _, v := range args {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}

	return rst
}

func flattenCollectorChannelModules(modulesResp interface{}) []interface{} {
	if modulesResp == nil {
		return nil
	}

	modules, ok := modulesResp.([]interface{})
	if !ok || len(modules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(modules))
	for _, v := range modules {
		rst = append(rst, map[string]interface{}{
			"name":                 utils.PathSearch("name", v, nil),
			"template_id":          utils.PathSearch("template_id", v, nil),
			"connection_module_id": utils.PathSearch("connection_module_id", v, nil),
			"children":             flattenCollectorChannelChildModules(utils.PathSearch("children", v, nil)),
			"fields":               flattenCollectorChannelFields(utils.PathSearch("fields", v, nil)),
		})
	}

	return rst
}

func flattenCollectorChannelChildModules(modulesResp interface{}) []interface{} {
	if modulesResp == nil {
		return nil
	}

	modules, ok := modulesResp.([]interface{})
	if !ok || len(modules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(modules))
	for _, v := range modules {
		rst = append(rst, map[string]interface{}{
			"name":                 utils.PathSearch("name", v, nil),
			"template_id":          utils.PathSearch("template_id", v, nil),
			"connection_module_id": utils.PathSearch("connection_module_id", v, nil),
			"fields":               flattenCollectorChannelFields(utils.PathSearch("fields", v, nil)),
		})
	}

	return rst
}

func flattenCollectorChannelFields(fieldsResp interface{}) []interface{} {
	if fieldsResp == nil {
		return nil
	}

	fields, ok := fieldsResp.([]interface{})
	if !ok || len(fields) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(fields))
	for _, v := range fields {
		rst = append(rst, map[string]interface{}{
			"name":                 utils.PathSearch("name", v, nil),
			"type":                 utils.PathSearch("type", v, nil),
			"value":                utils.PathSearch("value", v, nil),
			"other":                utils.PathSearch("other", v, nil),
			"template_field_id":    utils.PathSearch("template_field_id", v, nil),
			"connection_module_id": utils.PathSearch("connection_module_id", v, nil),
		})
	}

	return rst
}

func resourceCollectorChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/collector/channels/{channel_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workspace_id}", workspaceId)
	updatePath = strings.ReplaceAll(updatePath, "{channel_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCollectorChannelBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating collector channel: %s", err)
	}

	return resourceCollectorChannelRead(ctx, d, meta)
}

func resourceCollectorChannelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/collector/channels/{channel_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)
	deletePath = strings.ReplaceAll(deletePath, "{channel_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting collector channel (%s)", d.Id()))
	}

	if err := waitCollectorChannelDeleted(ctx, client, workspaceId, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for collector channel (%s) deleted: %s", d.Id(), err)
	}

	return nil
}

func waitCollectorChannelDeleted(ctx context.Context, client *golangsdk.ServiceClient, workspaceId, channelId string,
	timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      collectorChannelDeleteStateRefreshFunc(client, workspaceId, channelId),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func collectorChannelDeleteStateRefreshFunc(client *golangsdk.ServiceClient, workspaceId, channelId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		channel, err := GetCollectorChannelById(client, workspaceId, channelId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return map[string]interface{}{}, "DELETED", nil
			}
			return channel, "ERROR", err
		}

		return channel, "PENDING", nil
	}
}

func resourceCollectorChannelImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("workspace_id", parts[0])
}
