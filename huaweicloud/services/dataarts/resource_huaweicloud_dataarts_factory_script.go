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

// @API DataArtsStudio POST /v1/{project_id}/scripts
// @API DataArtsStudio DELETE /v1/{project_id}/scripts/{script_name}
// @API DataArtsStudio GET /v1/{project_id}/scripts/{script_name}
// @API DataArtsStudio PUT /v1/{project_id}/scripts/{script_name}
func ResourceDataArtsFactoryScript() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScriptCreate,
		ReadContext:   resourceScriptRead,
		UpdateContext: resourceScriptUpdate,
		DeleteContext: resourceScriptDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceScriptImportState,
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
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"connection_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"directory": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"queue_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// actually not in the GET response
			"target_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// actually not in the GET response
			"approvers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"approver_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"configuration": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_acquire_lock": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceScriptCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/scripts"
		product = "dataarts-dlf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateScriptBodyParams(d)),
		OkCodes:          []int{204},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DataArts script: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("workspace_id"), d.Get("name").(string)))
	return resourceScriptRead(ctx, d, meta)
}

func resourceScriptRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                      = meta.(*config.Config)
		region                   = cfg.GetRegion(d)
		mErr                     *multierror.Error
		httpUrl                  = "v1/{project_id}/scripts/{script_name}"
		product                  = "dataarts-dlf"
		resourceNotFoundErrCodes = []string{"DLF.0819", "DLF.6201"}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{script_name}", d.Get("name").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", resourceNotFoundErrCodes...),
			"error retrieving DataArts script")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getRespBody, nil)),
		d.Set("content", utils.PathSearch("content", getRespBody, nil)),
		d.Set("connection_name", utils.PathSearch("connectionName", getRespBody, nil)),
		d.Set("directory", utils.PathSearch("directory", getRespBody, nil)),
		d.Set("database", utils.PathSearch("database", getRespBody, nil)),
		d.Set("queue_name", utils.PathSearch("queueName", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("target_status", utils.PathSearch("targetStatus", getRespBody, nil)),
		d.Set("approvers", utils.PathSearch("approvers", getRespBody, nil)),
		d.Set("configuration", utils.PathSearch("configuration", getRespBody, nil)),
		d.Set("created_by", utils.PathSearch("createUser", getRespBody, nil)),
		d.Set("auto_acquire_lock", utils.PathSearch("autoAcquireLock", getRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceScriptUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/scripts/{script_name}"
		product = "dataarts-dlf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{script_name}", d.Get("name").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateScriptBodyParams(d)),
		OkCodes:          []int{204},
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating DataArts script: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("workspace_id"), d.Get("name").(string)))
	return resourceScriptRead(ctx, d, meta)
}

func resourceScriptDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/scripts/{script_name}"
		product = "dataarts-dlf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{script_name}", d.Get("name").(string))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceScriptImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	partLength := len(parts)

	if partLength == 2 {
		mErr := multierror.Append(nil,
			d.Set("workspace_id", parts[0]),
			d.Set("name", parts[1]),
		)
		return []*schema.ResourceData{d}, mErr.ErrorOrNil()
	}
	return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<name>")
}

func buildCreateOrUpdateScriptBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":           d.Get("name"),
		"type":           d.Get("type"),
		"content":        d.Get("content"),
		"connectionName": d.Get("connection_name"),
		"directory":      utils.ValueIgnoreEmpty(d.Get("directory")),
		"database":       utils.ValueIgnoreEmpty(d.Get("database")),
		"queueName":      utils.ValueIgnoreEmpty(d.Get("queue_name")),
		"configuration":  utils.ValueIgnoreEmpty(d.Get("configuration")),
		"description":    utils.ValueIgnoreEmpty(d.Get("description")),
		"targetStatus":   utils.ValueIgnoreEmpty(d.Get("target_status")),
		"approvers":      utils.ValueIgnoreEmpty(d.Get("approvers")),
	}
}
