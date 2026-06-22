package secmaster

import (
	"context"
	"fmt"
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

var collectorConnectionNonUpdatableParams = []string{
	"workspace_id",
}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/collector/connections
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/collector/connections/{connection_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/collector/connections/{connection_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/collector/connections/{connection_id}
func ResourceCollectorConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCollectorConnectionCreate,
		ReadContext:   resourceCollectorConnectionRead,
		UpdateContext: resourceCollectorConnectionUpdate,
		DeleteContext: resourceCollectorConnectionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCollectorConnectionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(collectorConnectionNonUpdatableParams),

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"other": {
							Type:     schema.TypeString,
							Required: true,
						},
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
						"template_field_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"field_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"fields_attribute": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"other": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_field_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"field_id": {
							Type:     schema.TypeString,
							Computed: true,
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
			// module_id is not required to be provided by the user, but it is required by the Update API
			// and automatically obtained from the Read response.
			"module_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCollectorConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"title":       d.Get("title"),
		"name":        d.Get("name"),
		"template_id": d.Get("template_id"),
		"module_id":   utils.ValueIgnoreEmpty(d.Get("module_id")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"fields":      buildCollectorConnectionFieldsBodyParams(d.Get("fields").([]interface{})),
	}
}

func buildCollectorConnectionFieldsBodyParams(fields []interface{}) []map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(fields))
	for _, v := range fields {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"title":             utils.ValueIgnoreEmpty(rawMap["title"]),
			"other":             utils.ValueIgnoreEmpty(rawMap["other"]),
			"name":              utils.ValueIgnoreEmpty(rawMap["name"]),
			"type":              utils.ValueIgnoreEmpty(rawMap["type"]),
			"value":             utils.ValueIgnoreEmpty(rawMap["value"]),
			"template_field_id": utils.ValueIgnoreEmpty(rawMap["template_field_id"]),
			"field_id":          utils.ValueIgnoreEmpty(rawMap["field_id"]),
		})
	}

	return rst
}

func resourceCollectorConnectionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/collector/connections"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", d.Get("workspace_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCollectorConnectionBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating collector connection: %s", err)
	}

	respBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	connectionId := utils.PathSearch("connection_id", respBody, "").(string)
	if connectionId == "" {
		return diag.Errorf("unable to find connection ID from the API response")
	}

	d.SetId(connectionId)

	return resourceCollectorConnectionRead(context.Background(), d, meta)
}

func GetCollectorConnectionById(client *golangsdk.ServiceClient, workspaceId, connectionId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/workspaces/{workspace_id}/collector/connections/{connection_id}"
	readPath := client.Endpoint + httpUrl
	readPath = strings.ReplaceAll(readPath, "{project_id}", client.ProjectID)
	readPath = strings.ReplaceAll(readPath, "{workspace_id}", workspaceId)
	readPath = strings.ReplaceAll(readPath, "{connection_id}", connectionId)

	readOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", readPath, &readOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceCollectorConnectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/collector/connections/{connection_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	readPath := client.Endpoint + httpUrl
	readPath = strings.ReplaceAll(readPath, "{project_id}", client.ProjectID)
	readPath = strings.ReplaceAll(readPath, "{workspace_id}", workspaceId)
	readPath = strings.ReplaceAll(readPath, "{connection_id}", d.Id())

	readOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", readPath, &readOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving collector connection")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workspace_id", workspaceId),
		d.Set("title", utils.PathSearch("title", respBody, nil)),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("template_id", utils.PathSearch("template_id", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("module_id", utils.PathSearch("module_id", respBody, nil)),
		d.Set("fields_attribute", flattenFieldsAttribute(utils.PathSearch("fields", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCollectorConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/collector/connections/{connection_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workspace_id}", workspaceId)
	updatePath = strings.ReplaceAll(updatePath, "{connection_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCollectorConnectionBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating collector connection: %s", err)
	}

	return resourceCollectorConnectionRead(ctx, d, meta)
}

func flattenFieldsAttribute(fieldsResp interface{}) []interface{} {
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
			"name":              utils.PathSearch("name", v, nil),
			"title":             utils.PathSearch("title", v, nil),
			"type":              utils.PathSearch("type", v, nil),
			"value":             utils.PathSearch("value", v, nil),
			"other":             utils.PathSearch("other", v, nil),
			"template_field_id": utils.PathSearch("template_field_id", v, nil),
			"field_id":          utils.PathSearch("field_id", v, nil),
		})
	}

	return rst
}

func resourceCollectorConnectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/collector/connections/{connection_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)
	deletePath = strings.ReplaceAll(deletePath, "{connection_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting collector connection (%s)", d.Id()))
	}

	return nil
}

func resourceCollectorConnectionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("workspace_id", parts[0])
}
