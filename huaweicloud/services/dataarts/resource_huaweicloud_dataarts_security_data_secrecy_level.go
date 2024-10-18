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

var DataSecrecyLevelResourceNotFoundCodes = []string{
	"DLS.6036",
	"DLS.1000",
}

// @API DataArtsStudio POST /v1/{project_id}/security/data-classification/secrecy-level
// @API DataArtsStudio GET /v1/{project_id}/security/data-classification/secrecy-level/{id}
// @API DataArtsStudio PUT /v1/{project_id}/security/data-classification/secrecy-level/{id}
// @API DataArtsStudio DELETE /v1/{project_id}/security/data-classification/secrecy-level/{id}
func ResourceSecurityDataSecrecyLevel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityDataSecrecyLevelCreate,
		ReadContext:   resourceSecurityDataSecrecyLevelRead,
		UpdateContext: resourceSecurityDataSecrecyLevelUpdate,
		DeleteContext: resourceSecurityDataSecrecyLevelDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityDataSecrecyLevelImportState,
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
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"level_number": {
				Type:     schema.TypeInt,
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

func resourceSecurityDataSecrecyLevelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/security/data-classification/secrecy-level"
		product = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	creatPath := client.Endpoint + httpUrl
	creatPath = strings.ReplaceAll(creatPath, "{project_id}", client.ProjectID)
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody:         utils.RemoveNil(buildCreateDataSecrecyLevelBodyParams(d)),
	}
	resp, err := client.Request("POST", creatPath, &opts)
	if err != nil {
		return diag.FromErr(err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error retrieving DataArts Security data secrecy level: %s", err)
	}

	levelId := utils.PathSearch("secrecy_level_id", respBody, "").(string)
	if levelId == "" {
		return diag.Errorf("unable to find the secrecy level ID of the DataArts Security from the API response")
	}

	d.SetId(levelId)
	return resourceSecurityDataSecrecyLevelRead(ctx, d, meta)
}

func buildCreateDataSecrecyLevelBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceSecurityDataSecrecyLevelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/security/data-classification/secrecy-level/{id}"
		product = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	seceryLevelId := d.Id()
	getPath = strings.ReplaceAll(getPath, "{id}", seceryLevelId)
	getDataSecrecyLevelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}
	resp, err := client.Request("GET", getPath, &getDataSecrecyLevelOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", DataSecrecyLevelResourceNotFoundCodes...),
			"error retrieving DataArts Security data secrecy level")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error retrieving DataArts Security data secrecy level (%s): %s", seceryLevelId, err)
	}

	createdAt := utils.PathSearch("created_at", respBody, float64(0)).(float64)
	updatedAt := utils.PathSearch("updated_at", respBody, float64(0)).(float64)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("secrecy_level_name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("level_number", utils.PathSearch("secrecy_level_number", respBody, 0)),
		d.Set("created_by", utils.PathSearch("created_by", respBody, nil)),
		d.Set("updated_by", utils.PathSearch("updated_by", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(createdAt)/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(updatedAt)/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSecurityDataSecrecyLevelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/security/data-classification/secrecy-level/{id}"
		product = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	seceryLevelId := d.Id()
	updatePath = strings.ReplaceAll(updatePath, "{id}", seceryLevelId)
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody: map[string]interface{}{
			"description": d.Get("description"),
		},
	}

	_, err = client.Request("PUT", updatePath, &opts)
	if err != nil {
		return diag.Errorf("error updating DataArts Security data secrecy level (%s): %s", seceryLevelId, err)
	}

	return resourceSecurityDataSecrecyLevelRead(ctx, d, meta)
}

func resourceSecurityDataSecrecyLevelDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/security/data-classification/secrecy-level/{id}"
		product = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	seceryLevelId := d.Id()
	deletePath = strings.ReplaceAll(deletePath, "{id}", seceryLevelId)
	deleteDataSecrecyLevelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}
	_, err = client.Request("DELETE", deletePath, &deleteDataSecrecyLevelOpt)
	if err != nil {
		return diag.Errorf("error deleting DataArts Security secrecy level (%s): %s", seceryLevelId, err)
	}

	return nil
}

func resourceSecurityDataSecrecyLevelImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(
		d.Set("workspace_id", parts[0]),
	)

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
