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

var nonUpdatableParamsSiemShipper = []string{
	"workspace_id",
	"shipper_name",
	"consumption_type",
	"dataspace_id",
	"dataspace_name",
	"domain_id",
	"pipe_id",
	"pipe_name",
	"project_id",
	"shipper_destination.*.data_param",
	"shipper_destination.*.destination_dataspace",
	"shipper_destination.*.destination_dataspace_name",
	"shipper_destination.*.destination_identity_role",
	"shipper_destination.*.destination_pipe",
	"shipper_destination.*.destination_pipe_name",
	"shipper_destination.*.destination_region",
	"shipper_destination.*.destination_shipper_type",
	"shipper_destination.*.destination_workspace",
	"shipper_destination.*.destination_workspace_name",
	"shipper_source.*.region",
	"shipper_source.*.source_dataspace",
	"shipper_source.*.source_dataspace_name",
	"shipper_source.*.source_identity_role",
	"shipper_source.*.source_pipe",
	"shipper_source.*.source_pipe_name",
	"shipper_source.*.source_type",
	"shipper_source.*.source_workspace",
	"shipper_source.*.source_workspace_name",
	"version",
	"workspace_name",
}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/siem/shippers
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/shippers
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/siem/shippers
// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/siem/shippers/{shipper_id}/pause
// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/siem/shippers/{shipper_id}/resume
// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/siem/shippers/{shipper_id}/retry
func ResourceSiemShipper() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSiemShipperCreate,
		ReadContext:   resourceSiemShipperRead,
		UpdateContext: resourceSiemShipperUpdate,
		DeleteContext: resourceSiemShipperDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSiemShipperImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsSiemShipper),

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
			"shipper_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"consumption_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"dataspace_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dataspace_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pipe_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pipe_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"shipper_destination": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     buildShipperDestinationSchema(),
			},
			"shipper_source": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     buildShipperSourceSchema(),
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// `workspace_name` is not returned in the response.
			"workspace_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Custom field, no return value.
			"action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"pause", "resume", "retry"}, false),
			},
			"shipper_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
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

func buildShipperDestinationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_param": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination_dataspace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination_dataspace_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination_identity_role": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination_pipe": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination_pipe_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination_region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination_shipper_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"destination_workspace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination_workspace_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func buildShipperSourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_dataspace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_dataspace_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_identity_role": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_pipe": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_pipe_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"source_workspace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_workspace_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func buildShipperDestinationBodyParams(rawList []interface{}) interface{} {
	if len(rawList) == 0 {
		return nil
	}

	item, ok := rawList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"data_param":                 utils.ValueIgnoreEmpty(item["data_param"]),
		"destination_dataspace":      utils.ValueIgnoreEmpty(item["destination_dataspace"]),
		"destination_dataspace_name": utils.ValueIgnoreEmpty(item["destination_dataspace_name"]),
		"destination_identity_role":  utils.ValueIgnoreEmpty(item["destination_identity_role"]),
		"destination_pipe":           utils.ValueIgnoreEmpty(item["destination_pipe"]),
		"destination_pipe_name":      utils.ValueIgnoreEmpty(item["destination_pipe_name"]),
		"destination_region":         utils.ValueIgnoreEmpty(item["destination_region"]),
		"destination_shipper_type":   item["destination_shipper_type"],
		"destination_workspace":      utils.ValueIgnoreEmpty(item["destination_workspace"]),
		"destination_workspace_name": utils.ValueIgnoreEmpty(item["destination_workspace_name"]),
	}
}

func buildShipperSourceBodyParams(rawList []interface{}) interface{} {
	if len(rawList) == 0 {
		return nil
	}

	item, ok := rawList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"region":                utils.ValueIgnoreEmpty(item["region"]),
		"source_dataspace":      utils.ValueIgnoreEmpty(item["source_dataspace"]),
		"source_dataspace_name": utils.ValueIgnoreEmpty(item["source_dataspace_name"]),
		"source_identity_role":  utils.ValueIgnoreEmpty(item["source_identity_role"]),
		"source_pipe":           utils.ValueIgnoreEmpty(item["source_pipe"]),
		"source_pipe_name":      utils.ValueIgnoreEmpty(item["source_pipe_name"]),
		"source_type":           item["source_type"],
		"source_workspace":      utils.ValueIgnoreEmpty(item["source_workspace"]),
		"source_workspace_name": utils.ValueIgnoreEmpty(item["source_workspace_name"]),
	}
}

func buildCreateSiemShipperBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"shipper_name":        d.Get("shipper_name"),
		"consumption_type":    d.Get("consumption_type"),
		"dataspace_id":        utils.ValueIgnoreEmpty(d.Get("dataspace_id")),
		"dataspace_name":      utils.ValueIgnoreEmpty(d.Get("dataspace_name")),
		"domain_id":           utils.ValueIgnoreEmpty(d.Get("domain_id")),
		"pipe_id":             utils.ValueIgnoreEmpty(d.Get("pipe_id")),
		"pipe_name":           utils.ValueIgnoreEmpty(d.Get("pipe_name")),
		"project_id":          utils.ValueIgnoreEmpty(d.Get("project_id")),
		"shipper_destination": buildShipperDestinationBodyParams(d.Get("shipper_destination").([]interface{})),
		"shipper_source":      buildShipperSourceBodyParams(d.Get("shipper_source").([]interface{})),
		"version":             utils.ValueIgnoreEmpty(d.Get("version")),
		"workspace_id":        d.Get("workspace_id"),
		"workspace_name":      utils.ValueIgnoreEmpty(d.Get("workspace_name")),
	}

	return bodyParams
}

func resourceSiemShipperCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		workspaceId   = d.Get("workspace_id").(string)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/siem/shippers"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)

	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateSiemShipperBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster SIEM shipper: %s", err)
	}

	d.SetId(d.Get("shipper_name").(string))

	if action := d.Get("action").(string); action != "" {
		if err := operationShipper(client, d, action); err != nil {
			return diag.Errorf("error operating SecMaster shipper (%s): %s", action, err)
		}
	}

	return resourceSiemShipperRead(ctx, d, meta)
}

func GetSiemShipperInfo(client *golangsdk.ServiceClient, workspaceId, shipperName string) (interface{}, error) {
	getPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/siem/shippers"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath += fmt.Sprintf("?shipper_name=%s&limit=10&offset=0", shipperName)

	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	shipperInfo := utils.PathSearch("data.data[0]", respBody, nil)
	if shipperInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return shipperInfo, nil
}

func flattenShipperDestinationFromList(respBody interface{}) []map[string]interface{} {
	destinationRaw := utils.PathSearch("shipper_destination", respBody, nil)
	if destinationRaw == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"data_param":                 utils.PathSearch("data_param", destinationRaw, nil),
			"destination_dataspace":      utils.PathSearch("dataspace", destinationRaw, nil),
			"destination_dataspace_name": utils.PathSearch("dataspace_name", destinationRaw, nil),
			"destination_identity_role":  utils.PathSearch("identity", destinationRaw, nil),
			"destination_pipe":           utils.PathSearch("pipe", destinationRaw, nil),
			"destination_pipe_name":      utils.PathSearch("pipe_name", destinationRaw, nil),
			"destination_region":         utils.PathSearch("region", destinationRaw, nil),
			"destination_shipper_type":   utils.PathSearch("data_type", destinationRaw, nil),
			"destination_workspace":      utils.PathSearch("workspace", destinationRaw, nil),
			"destination_workspace_name": utils.PathSearch("workspace_name", destinationRaw, nil),
		},
	}
}

func flattenShipperSourceFromList(respBody interface{}) []map[string]interface{} {
	sourceRaw := utils.PathSearch("shipper_source", respBody, nil)
	if sourceRaw == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"region":                utils.PathSearch("region", sourceRaw, nil),
			"source_dataspace":      utils.PathSearch("dataspace", sourceRaw, nil),
			"source_dataspace_name": utils.PathSearch("dataspace_name", sourceRaw, nil),
			"source_identity_role":  utils.PathSearch("identity", sourceRaw, nil),
			"source_pipe":           utils.PathSearch("pipe", sourceRaw, nil),
			"source_pipe_name":      utils.PathSearch("pipe_name", sourceRaw, nil),
			"source_type":           utils.PathSearch("data_type", sourceRaw, nil),
			"source_workspace":      utils.PathSearch("workspace", sourceRaw, nil),
			"source_workspace_name": utils.PathSearch("workspace_name", sourceRaw, nil),
		},
	}
}

func resourceSiemShipperRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	respBody, err := GetSiemShipperInfo(client, workspaceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			err,
			"error retrieving SecMaster SIEM shipper",
		)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("shipper_name", utils.PathSearch("shipper_name", respBody, nil)),
		d.Set("shipper_id", utils.PathSearch("shipper_id", respBody, nil)),
		d.Set("consumption_type", utils.PathSearch("consumption_type", respBody, nil)),
		d.Set("dataspace_id", utils.PathSearch("dataspace_id", respBody, nil)),
		d.Set("dataspace_name", utils.PathSearch("dataspace_name", respBody, nil)),
		d.Set("domain_id", utils.PathSearch("domain_id", respBody, nil)),
		d.Set("pipe_id", utils.PathSearch("pipe_id", respBody, nil)),
		d.Set("pipe_name", utils.PathSearch("pipe_name", respBody, nil)),
		d.Set("project_id", utils.PathSearch("project_id", respBody, nil)),
		d.Set("shipper_destination", flattenShipperDestinationFromList(respBody)),
		d.Set("shipper_source", flattenShipperSourceFromList(respBody)),
		d.Set("version", utils.PathSearch("version", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func operationShipper(client *golangsdk.ServiceClient, d *schema.ResourceData, operateType string) error {
	var (
		workspaceId = d.Get("workspace_id").(string)
		shipperId   = d.Get("shipper_id").(string)
		httpUrl     = fmt.Sprintf("v1/{project_id}/workspaces/{workspace_id}/siem/shippers/{shipper_id}/%s", operateType)
	)

	operatePath := client.Endpoint + httpUrl
	operatePath = strings.ReplaceAll(operatePath, "{project_id}", client.ProjectID)
	operatePath = strings.ReplaceAll(operatePath, "{workspace_id}", workspaceId)
	operatePath = strings.ReplaceAll(operatePath, "{shipper_id}", shipperId)

	operateOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	_, err := client.Request("POST", operatePath, &operateOpt)
	return err
}

func resourceSiemShipperUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		action = d.Get("action").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	if d.HasChange("action") && action != "" {
		if err := operationShipper(client, d, action); err != nil {
			return diag.Errorf("error operating SecMaster SIEM shipper (%s): %s", action, err)
		}
	}

	return resourceSiemShipperRead(ctx, d, meta)
}

func resourceSiemShipperDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/shippers"
		shipperId   = d.Get("shipper_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)

	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         map[string]interface{}{"ids": []string{shipperId}},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster SIEM shipper: %s", err)
	}

	return nil
}

func resourceSiemShipperImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<shipper_name>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("workspace_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
