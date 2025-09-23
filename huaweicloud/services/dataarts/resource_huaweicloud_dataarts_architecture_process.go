// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DataArtsStudio
// ---------------------------------------------------------------

package dataarts

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio POST /v2/{project_id}/design/biz/catalogs
// @API DataArtsStudio DELETE /v2/{project_id}/design/biz/catalogs
// @API DataArtsStudio GET /v2/{project_id}/design/biz/catalogs/tree
// @API DataArtsStudio GET /v2/{project_id}/design/biz/catalogs/{id}
// @API DataArtsStudio PUT /v2/{project_id}/design/biz/catalogs
func ResourceArchitectureProcess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchitectureProcessCreate,
		ReadContext:   resourceArchitectureProcessRead,
		UpdateContext: resourceArchitectureProcessUpdate,
		DeleteContext: resourceArchitectureProcessDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceArchitectureProcessImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the process is located.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of DataArts Studio workspace.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of catalog.",
			},
			"owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of person responsible for process.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of process.",
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The parent process ID of process.",
			},
			"prev_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the previous node in the process.",
			},
			"next_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the next node in the process.",
			},
			"qualified_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of all superior processes.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the process.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the process.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator of the process.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last editor of the process.",
			},
			"children": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The name list of subordinate process.",
			},
		},
	}
}

func resourceArchitectureProcessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/design/biz/catalogs"
		product = "dataarts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateArchitectureProcessBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DataArts Architecture process: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}
	processID := utils.PathSearch("data.value.id", createRespBody, "").(string)
	if processID == "" {
		return diag.Errorf("unable to find the DataArts Architecture process ID from the API response")
	}

	// need to set qualified ID to filter result in READ.
	qualifiedID := utils.PathSearch("data.value.qualified_id", createRespBody, "").(string)
	if qualifiedID == "" {
		return diag.Errorf("unable to find the qualified ID of the DataArts Architecture process from the API response")
	}

	d.SetId(processID)
	d.Set("qualified_id", qualifiedID)

	return resourceArchitectureProcessRead(ctx, d, meta)
}

func buildCreateOrUpdateArchitectureProcessBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"owner":       d.Get("owner"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"parent_id":   utils.ValueIgnoreEmpty(d.Get("parent_id")),
		"prev_id":     utils.ValueIgnoreEmpty(d.Get("prev_id")),
		"next_id":     utils.ValueIgnoreEmpty(d.Get("next_id")),
	}
}

func resourceArchitectureProcessRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/design/biz/catalogs/tree"
		product = "dataarts"
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DataArts Architecture process")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	paths := strings.Split(d.Get("qualified_id").(string), ".")
	jsonPaths := fmt.Sprintf("data.value[?id=='%s']", paths[0])
	for _, path := range paths[1:] {
		jsonPaths += fmt.Sprintf("[children][][?id=='%s'][]", path)
	}

	processes := utils.PathSearch(jsonPaths, getRespBody, make([]interface{}, 0)).([]interface{})
	if len(processes) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "DataArts Architecture process")
	}

	process := processes[0]
	mErr = multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", process, nil)),
		d.Set("owner", utils.PathSearch("owner", process, nil)),
		d.Set("description", utils.PathSearch("description", process, nil)),
		d.Set("parent_id", utils.PathSearch("parent_id", process, nil)),
		d.Set("prev_id", utils.PathSearch("prev_id", process, nil)),
		d.Set("next_id", utils.PathSearch("next_id", process, nil)),
		d.Set("qualified_id", utils.PathSearch("qualified_id", process, nil)),
		d.Set("created_at", utils.PathSearch("create_time", process, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", process, nil)),
		d.Set("created_by", utils.PathSearch("create_by", process, nil)),
		d.Set("updated_by", utils.PathSearch("update_by", process, nil)),
		d.Set("children", utils.PathSearch(`children[*].name`, process, make([]interface{}, 0))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceArchitectureProcessUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v2/{project_id}/design/biz/catalogs"
		product     = "dataarts"
		workspaceID = d.Get("workspace_id").(string)
		qualifiedID string
		processID   = d.Id()
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateBody := utils.RemoveNil(buildCreateOrUpdateArchitectureProcessBodyParams(d))
	updateBody["id"] = d.Id()
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
		JSONBody:         updateBody,
	}

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating DataArts Architecture process: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// After calling the update API, qualified ID always returns null.
	// Therefore, it is necessary to recombine the values of this field to filter the result in READ.
	parentID := utils.PathSearch("data.value.parent_id", updateRespBody, "").(string)
	if parentID != "" {
		getResp, err := readArchitectureProcess(client, parentID, workspaceID)
		if err != nil {
			return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", "DLG.3903"),
				"error retrieving DataArts Architecture process")
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		parentProcessParentID := utils.PathSearch("data.value.parent_id", getRespBody, "").(string)
		if parentProcessParentID != "" {
			qualifiedID = fmt.Sprintf("%v.%v.%v", parentProcessParentID, parentID, processID)
		} else {
			log.Printf("[DEBUG] unable to find the parent ID of the DataArts Architecture parent process from the API response")
			qualifiedID = fmt.Sprintf("%v.%v", parentID, processID)
		}
	} else {
		qualifiedID = processID
	}

	d.Set("qualified_id", qualifiedID)

	return resourceArchitectureProcessRead(ctx, d, meta)
}

func resourceArchitectureProcessDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v2/{project_id}/design/biz/catalogs"
		product     = "dataarts"
		workspaceID = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
		JSONBody:         buildDeleteArchitectureProcessBodyParams(d),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting DataArts Architecture process: %s", err)
	}

	// Successful deletion API call does not guarantee that the resource is successfully deleted.
	// Call the details API to confirm that the resource has been successfully deleted.
	_, err = readArchitectureProcess(client, d.Id(), workspaceID)
	if err == nil {
		return diag.Errorf("error deleting DataArts Architecture process: the process still exists")
	}

	return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", "DLG.3903"),
		"error deleting DataArts Architecture process")
}

func buildDeleteArchitectureProcessBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"ids": []string{d.Id()},
	}
}

func readArchitectureProcess(client *golangsdk.ServiceClient, catalogID, workspaceID string) (*http.Response, error) {
	getPath := client.Endpoint + "v2/{project_id}/design/biz/catalogs/{id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", catalogID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}

	return client.Request("GET", getPath, &getOpt)
}

func resourceArchitectureProcessImportState(_ context.Context, d *schema.ResourceData, meta interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <workspace_id>/<qualified_id>")
	}

	var (
		workspaceID = parts[0]
		qualifiedID = parts[1]
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v2/{project_id}/design/biz/catalogs/tree"
		product     = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error query all DataArts Architecture process trees: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}

	paths := strings.Split(qualifiedID, ".")
	jsonPaths := fmt.Sprintf("data.value[?id=='%s']", paths[0])
	for _, path := range paths[1:] {
		jsonPaths += fmt.Sprintf("[children][][?id=='%s'][]", path)
	}

	processes := utils.PathSearch(jsonPaths, getRespBody, make([]interface{}, 0)).([]interface{})
	if len(processes) < 1 {
		return []*schema.ResourceData{d}, fmt.Errorf("data process not found: %s", err)
	}

	d.SetId(utils.PathSearch("id", processes[0], nil).(string))
	d.Set("qualified_id", qualifiedID)
	d.Set("workspace_id", workspaceID)
	return []*schema.ResourceData{d}, nil
}
