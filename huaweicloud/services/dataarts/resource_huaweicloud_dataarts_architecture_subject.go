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

// @API DataArtsStudio POST /v3/{project_id}/design/subjects
// @API DataArtsStudio DELETE /v3/{project_id}/design/subjects
// @API DataArtsStudio GET /v3/{project_id}/design/subjects
// @API DataArtsStudio PUT /v3/{project_id}/design/subjects
func ResourceArchitectureSubject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchitectureSubjectCreate,
		ReadContext:   resourceArchitectureSubjectRead,
		UpdateContext: resourceArchitectureSubjectUpdate,
		DeleteContext: resourceArchitectureSubjectDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceArchitectureSubjectImportState,
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
			"code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Required: true,
			},
			"level": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"department": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
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
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArchitectureSubjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	createSubjectHttpUrl := "v3/{project_id}/design/subjects"
	createSubjectProduct := "dataarts"

	createSubjectClient, err := cfg.NewServiceClient(createSubjectProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}
	createSubjectPath := createSubjectClient.Endpoint + createSubjectHttpUrl
	createSubjectPath = strings.ReplaceAll(createSubjectPath, "{project_id}", createSubjectClient.ProjectID)

	createSubjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}
	createSubjectOpt.JSONBody = utils.RemoveNil(buildCreateArchitectureSubjectBodyParams(d))
	createSubjectResp, err := createSubjectClient.Request("POST", createSubjectPath, &createSubjectOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	createSubjectRespBody, err := utils.FlattenResponse(createSubjectResp)
	if err != nil {
		return diag.FromErr(err)
	}

	subjectId := utils.PathSearch("data.value.id", createSubjectRespBody, "").(string)
	if subjectId == "" {
		return diag.Errorf("unable to find the DataArts Architecture subject ID from the API response")
	}

	d.SetId(subjectId)

	return resourceArchitectureSubjectRead(ctx, d, meta)
}

func buildCreateArchitectureSubjectBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name_ch":         d.Get("name"),
		"name_en":         d.Get("code"),
		"data_owner_list": d.Get("owner"),
		"level":           d.Get("level"),
		"data_owner":      utils.ValueIgnoreEmpty(d.Get("department")),
		"description":     utils.ValueIgnoreEmpty(d.Get("description")),
		"parent_id":       utils.ValueIgnoreEmpty(d.Get("parent_id")),
	}
	return bodyParams
}

func resourceArchitectureSubjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	workspaceID := d.Get("workspace_id").(string)

	getSubjectHttpUrl := "v3/{project_id}/design/subjects"
	getSubjectProduct := "dataarts"

	getSubjectClient, err := cfg.NewServiceClient(getSubjectProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}

	getSubjectPath := getSubjectClient.Endpoint + getSubjectHttpUrl
	getSubjectPath = strings.ReplaceAll(getSubjectPath, "{project_id}", getSubjectClient.ProjectID)

	getSubjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}

	getSubjectPath += "?limit=10"

	// using name to reduce the search results
	if val, ok := d.GetOk("name"); ok {
		getSubjectPath += fmt.Sprintf("&name=%s", val)
	}

	currentTotal := 0
	for {
		path := fmt.Sprintf("%s&offset=%v", getSubjectPath, currentTotal)
		getSubjectResp, err := getSubjectClient.Request("GET", path, &getSubjectOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", workspaceIdNotFound),
				"error retrieving DataArts Architecture subject")
		}
		getSubjectRespBody, err := utils.FlattenResponse(getSubjectResp)
		if err != nil {
			return diag.FromErr(err)
		}
		subjects := utils.PathSearch("data.value.records", getSubjectRespBody, make([]interface{}, 0)).([]interface{})
		total := utils.PathSearch("data.value.total", getSubjectRespBody, 0)
		var mErr *multierror.Error
		for _, subject := range subjects {
			// using path to filter result for import, because ID can not be got from console
			// format of path in results using `/` to split, format of path from import using `.` to split
			id := utils.PathSearch("id", subject, "")
			path := strings.ReplaceAll(utils.PathSearch("path", subject, "").(string), "/", ".")
			if val, ok := d.GetOk("path"); ok {
				if path != val {
					continue
				}
			} else if id != d.Id() {
				continue
			}

			// set ID once more to cover import ID
			d.SetId(utils.PathSearch("id", subject, "").(string))

			mErr = multierror.Append(nil,
				d.Set("region", region),
				d.Set("workspace_id", workspaceID),
				d.Set("name", utils.PathSearch("name_ch", subject, nil)),
				d.Set("code", utils.PathSearch("name_en", subject, nil)),
				d.Set("owner", utils.PathSearch("data_owner_list", subject, nil)),
				d.Set("level", utils.PathSearch("level", subject, nil)),
				d.Set("department", utils.PathSearch("data_owner", subject, nil)),
				d.Set("description", utils.PathSearch("description", subject, nil)),
				d.Set("parent_id", utils.PathSearch("parent_id", subject, nil)),
				d.Set("guid", utils.PathSearch("guid", subject, nil)),
				d.Set("path", path),
				d.Set("status", utils.PathSearch("status", subject, nil)),
				d.Set("created_at", utils.PathSearch("create_time", subject, nil)),
				d.Set("updated_at", utils.PathSearch("update_time", subject, nil)),
				d.Set("created_by", utils.PathSearch("create_by", subject, nil)),
				d.Set("updated_by", utils.PathSearch("update_by", subject, nil)),
			)
			return diag.FromErr(mErr.ErrorOrNil())
		}
		currentTotal += len(subjects)
		// type of `total` is float64
		if float64(currentTotal) == total {
			break
		}
	}

	return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
}

func resourceArchitectureSubjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateSubjectHttpUrl := "v3/{project_id}/design/subjects"
	updateSubjectProduct := "dataarts"

	updateSubjectClient, err := cfg.NewServiceClient(updateSubjectProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}
	updateSubjectPath := updateSubjectClient.Endpoint + updateSubjectHttpUrl
	updateSubjectPath = strings.ReplaceAll(updateSubjectPath, "{project_id}", updateSubjectClient.ProjectID)

	updateSubjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	updateSubjectOpt.JSONBody = utils.RemoveNil(buildUpdateArchitectureSubjectBodyParams(d))
	updateSubjectResp, err := updateSubjectClient.Request("PUT", updateSubjectPath, &updateSubjectOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	updateSubjectRespBody, err := utils.FlattenResponse(updateSubjectResp)
	if err != nil {
		return diag.FromErr(err)
	}

	path := strings.ReplaceAll(utils.PathSearch("data.value.path", updateSubjectRespBody, "").(string), "/", ".")
	d.Set("path", path)

	return resourceArchitectureSubjectRead(ctx, d, meta)
}

func buildUpdateArchitectureSubjectBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id":              d.Id(),
		"name_ch":         d.Get("name"),
		"name_en":         d.Get("code"),
		"data_owner_list": d.Get("owner"),
		"level":           d.Get("level"),
		"data_owner":      utils.ValueIgnoreEmpty(d.Get("department")),
		"description":     utils.ValueIgnoreEmpty(d.Get("description")),
		"parent_id":       utils.ValueIgnoreEmpty(d.Get("parent_id")),
	}
	return bodyParams
}

func resourceArchitectureSubjectDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	deleteSubjectHttpUrl := "v3/{project_id}/design/subjects"
	deleteSubjectProduct := "dataarts"

	deleteSubjectClient, err := cfg.NewServiceClient(deleteSubjectProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}
	deleteSubjectPath := deleteSubjectClient.Endpoint + deleteSubjectHttpUrl
	deleteSubjectPath = strings.ReplaceAll(deleteSubjectPath, "{project_id}", deleteSubjectClient.ProjectID)

	deleteSubjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	deleteSubjectOpt.JSONBody = utils.RemoveNil(map[string]interface{}{
		"ids": []string{d.Id()},
	})

	_, err = deleteSubjectClient.Request("DELETE", deleteSubjectPath, &deleteSubjectOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceArchitectureSubjectImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <workspace_id>/<path>")
	}

	names := strings.Split(parts[1], ".")

	d.Set("name", names[len(names)-1])
	d.Set("workspace_id", parts[0])
	d.Set("path", parts[1])

	return []*schema.ResourceData{d}, nil
}
