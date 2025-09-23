package sfsturbo

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

// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir-quota
// @API SFSTurbo PUT /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir-quota
// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir-quota
// @API SFSTurbo DELETE /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir-quota
func ResourceSfsTurboDirQuota() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSfsTurboDirQuotaCreate,
		ReadContext:   resourceSfsTurboDirQuotaRead,
		UpdateContext: resourceSfsTurboDirQuotaUpdate,
		DeleteContext: resourceSfsTurboDirQuotaDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"share_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the SFS Turbo ID.`,
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the valid full path of an existing directory.`,
			},
			"capacity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the size of the directory.`,
			},
			"inode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the maximum number of inodes allowed in the directory.`,
			},
			"used_capacity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the size of the used directory.`,
			},
			"used_inode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Specifies the number of used inodes in the directory.`,
			},
		},
	}
}

func resourceSfsTurboDirQuotaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir-quota"
	)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	createPath := sfsClient.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", sfsClient.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{share_id}", d.Get("share_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildSfsTurboDirQuotaBodyParams(d)),
	}

	_, err = sfsClient.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo directory quota: %s", err)
	}

	path := d.Get("path").(string)
	d.SetId(path)

	return resourceSfsTurboDirQuotaRead(ctx, d, meta)
}

func resourceSfsTurboDirQuotaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir-quota"
	)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	updatePath := sfsClient.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", sfsClient.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{share_id}", d.Get("share_id").(string))
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildSfsTurboDirQuotaBodyParams(d)),
	}

	_, err = sfsClient.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SFS Turbo directory quota: %s", err)
	}

	return resourceSfsTurboDirQuotaRead(ctx, d, meta)
}

func buildSfsTurboDirQuotaBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"path":     d.Get("path"),
		"capacity": d.Get("capacity"),
		"inode":    d.Get("inode"),
	}

	return params
}

func resourceSfsTurboDirQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir-quota"
	)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	getPath := sfsClient.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", sfsClient.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{share_id}", d.Get("share_id").(string))
	getPath += fmt.Sprintf("?path=%s", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := sfsClient.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SFS Turbo directory quota")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("path", utils.PathSearch("path", getRespBody, nil)),
		d.Set("capacity", utils.PathSearch("capacity", getRespBody, nil)),
		d.Set("inode", utils.PathSearch("inode", getRespBody, nil)),
		d.Set("used_capacity", utils.PathSearch("used_capacity", getRespBody, nil)),
		d.Set("used_inode", utils.PathSearch("used_inode", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDeleteSfsTurboQuotaDirBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"path": d.Id(),
	}
}

func resourceSfsTurboDirQuotaDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir-quota"
	)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	deletePath := sfsClient.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", sfsClient.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{share_id}", d.Get("share_id").(string))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeleteSfsTurboQuotaDirBodyParams(d),
	}

	_, err = sfsClient.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		if utils.IsResourceNotFound(err) {
			return nil
		}
		return diag.Errorf("error deleting SFS Turbo directory quota: %s", err)
	}

	return nil
}
