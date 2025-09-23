package dataarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var modelErrCodes = []string{
	"DLG.0818", // Workspace does not exist.
	"DLG.6026", // Resource does not exist.
	"DLG.3902", // Resource ID value is incorrect.
}

// @API DataArtsStudio POST /v2/{project_id}/design/workspaces
// @API DataArtsStudio GET /v2/{project_id}/design/workspaces
// @API DataArtsStudio PUT /v2/{project_id}/design/workspaces
// @API DataArtsStudio GET /v2/{project_id}/design/workspaces/{model_id}
// @API DataArtsStudio DELETE /v2/{project_id}/design/workspaces
func ResourceArchitectureModel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchitectureModelCreate,
		ReadContext:   resourceArchitectureModelRead,
		UpdateContext: resourceArchitectureModelUpdate,
		DeleteContext: resourceArchitectureModelDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceArchitectureModelImportState,
		},

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
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"physical": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dw_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"level": {
				Type:     schema.TypeString,
				Optional: true,
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
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateModelParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"type":        d.Get("type"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"dw_type":     d.Get("dw_type"),
		"level":       utils.ValueIgnoreEmpty(d.Get("level")),
		"is_physical": d.Get("physical"),
	}
	return bodyParams
}

func resourceArchitectureModelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	createModelHttpUrl := "v2/{project_id}/design/workspaces"
	createModelProduct := "dataarts"

	modelClient, err := cfg.NewServiceClient(createModelProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	createModelPath := modelClient.Endpoint + createModelHttpUrl
	createModelPath = strings.ReplaceAll(createModelPath, "{project_id}", modelClient.ProjectID)
	createModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}
	createModelOpt.JSONBody = utils.RemoveNil(buildCreateModelParams(d))
	createModelResp, err := modelClient.Request("POST", createModelPath, &createModelOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	createModelBody, err := utils.FlattenResponse(createModelResp)
	if err != nil {
		return diag.FromErr(err)
	}
	modelId := utils.PathSearch("data.value.id", createModelBody, "").(string)
	if modelId == "" {
		return diag.Errorf("unable to find the DataArts Studio architecture model ID from the API response")
	}
	d.SetId(modelId)

	return resourceArchitectureModelRead(ctx, d, meta)
}

func resourceArchitectureModelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	readModelProduct := "dataarts"
	modelClient, err := cfg.NewServiceClient(readModelProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	var model interface{}
	if strings.Contains(d.Id(), "/") {
		// if id contains name, it is for import
		model, err = getModelByName(modelClient, d)
	} else {
		model, err = getModelById(modelClient, d)
	}
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DataArts Architecture model")
	}
	d.SetId(utils.PathSearch("id", model, "").(string))
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workspace_id", d.Get("workspace_id").(string)),
		d.Set("name", utils.PathSearch("name", model, "")),
		d.Set("type", utils.PathSearch("type", model, "")),
		d.Set("dw_type", utils.PathSearch("dw_type", model, "")),
		d.Set("level", utils.PathSearch("level", model, "")),
		d.Set("description", utils.PathSearch("description", model, "")),
		d.Set("physical", utils.PathSearch("is_physical", model, nil)),
		d.Set("created_at", utils.PathSearch("create_time", model, "")),
		d.Set("updated_at", utils.PathSearch("update_time", model, "")),
		d.Set("created_by", utils.PathSearch("create_by", model, "")),
		d.Set("updated_by", utils.PathSearch("update_by", model, "")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getModelByName(modelClient *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	readModelHttpUrl := "v2/{project_id}/design/workspaces"
	readModelPath := modelClient.Endpoint + readModelHttpUrl
	readModelPath = strings.ReplaceAll(readModelPath, "{project_id}", modelClient.ProjectID)
	workspaceID := d.Get("workspace_id").(string)
	readModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}
	currentTotal := 0
	readModelPath += fmt.Sprintf("?limit=50&offset=%v", currentTotal)
	for {
		readModelResp, err := modelClient.Request("GET", readModelPath, &readModelOpt)
		if err != nil {
			// Only one scenario where the workspace ID does not exist, the error code is "DLG.0818".
			return nil, common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code")
		}
		readModelBody, err := utils.FlattenResponse(readModelResp)
		if err != nil {
			return nil, err
		}
		models := utils.PathSearch("data.value.records", readModelBody, make([]interface{}, 0)).([]interface{})
		total := utils.PathSearch("data.value.total", readModelBody, 0)
		if len(models) == 0 {
			return nil, golangsdk.ErrDefault404{}
		}
		for _, model := range models {
			// using name to filter result for import, because ID can not be get from console
			name := utils.PathSearch("name", model, "")
			if name != d.Get("name") {
				continue
			}
			return model, nil
		}
		currentTotal += len(models)
		// type of `total` is float64
		if float64(currentTotal) == total {
			break
		}
		index := strings.Index(readModelPath, "offset")
		readModelPath = fmt.Sprintf("%soffset=%v", readModelPath[:index], currentTotal)
	}
	return nil, golangsdk.ErrDefault404{}
}

func getModelById(modelClient *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	readModelHttpUrl := "v2/{project_id}/design/workspaces/{model_id}"
	readModelPath := modelClient.Endpoint + readModelHttpUrl
	readModelPath = strings.ReplaceAll(readModelPath, "{project_id}", modelClient.ProjectID)
	readModelPath = strings.ReplaceAll(readModelPath, "{model_id}", d.Id())
	workspaceID := d.Get("workspace_id").(string)
	readModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}
	readModelResp, err := modelClient.Request("GET", readModelPath, &readModelOpt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", modelErrCodes...)
	}
	readModelBody, err := utils.FlattenResponse(readModelResp)
	if err != nil {
		return nil, err
	}
	model := utils.PathSearch("data.value", readModelBody, nil)
	return model, nil
}

func resourceArchitectureModelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	modelProduct := "dataarts"
	modelClient, err := cfg.NewServiceClient(modelProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	updateModelHttpUrl := "v2/{project_id}/design/workspaces"
	updateModelPath := modelClient.Endpoint + updateModelHttpUrl
	updateModelPath = strings.ReplaceAll(updateModelPath, "{project_id}", modelClient.ProjectID)
	workspaceID := d.Get("workspace_id").(string)
	readModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}
	readModelOpt.JSONBody = utils.RemoveNil(buildUpdateModelBodyParams(d))
	_, err = modelClient.Request("PUT", updateModelPath, &readModelOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceArchitectureModelRead(ctx, d, meta)
}

func buildUpdateModelBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id":          d.Id(),
		"name":        d.Get("name"),
		"type":        d.Get("type"),
		"dw_type":     d.Get("dw_type"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"level":       utils.ValueIgnoreEmpty(d.Get("level")),
		"is_physical": utils.ValueIgnoreEmpty(d.Get("physical")),
	}
	return bodyParams
}

func resourceArchitectureModelDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	modelProduct := "dataarts"
	modelClient, err := cfg.NewServiceClient(modelProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}
	delModelHttpUrl := "v2/{project_id}/design/workspaces?ids="
	delModelPath := modelClient.Endpoint + delModelHttpUrl
	delModelPath = strings.ReplaceAll(delModelPath, "{project_id}", modelClient.ProjectID)
	delModelPath += d.Id()
	workspaceID := d.Get("workspace_id").(string)
	delModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}
	_, err = modelClient.Request("DELETE", delModelPath, &delModelOpt)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceArchitectureModelImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <workspace_id>/<name>")
	}
	d.Set("workspace_id", parts[0])
	d.Set("name", parts[1])
	return []*schema.ResourceData{d}, nil
}
