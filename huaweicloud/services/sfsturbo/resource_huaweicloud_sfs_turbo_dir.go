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

// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir
// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir
// @API SFSTurbo DELETE /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir
func ResourceSfsTurboDir() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSfsTurboDirCreate,
		ReadContext:   resourceSfsTurboDirRead,
		DeleteContext: resourceSfsTurboDirDelete,

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
				Description: `Specifies the valid full path of SFS Turbo directory.`,
			},
			"mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the SFS Turbo directory permissions.`,
			},
			"uid": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the user ID of the SFS Turbo directory.`,
			},
			"gid": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the group ID of the SFS Turbo directory.`,
			},
		},
	}
}

func resourceSfsTurboDirCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir"
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
		JSONBody:         utils.RemoveNil(buildCreateSfsTurboDirBodyParams(d)),
		OkCodes: []int{
			204,
		},
	}

	_, err = sfsClient.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo directory: %s", err)
	}

	path := d.Get("path").(string)
	d.SetId(path)

	return resourceSfsTurboDirRead(ctx, d, meta)
}

func buildCreateSfsTurboDirBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"path": d.Get("path"),
		"mode": d.Get("mode"),
		"uid":  d.Get("uid"),
		"gid":  d.Get("gid"),
	}

	return params
}

func resourceSfsTurboDirRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir"
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
		return common.CheckDeletedDiag(d, err, "error retrieving SFS Turbo directory")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("path", utils.PathSearch("path", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDeleteSfsTurboDirBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"path": d.Id(),
	}
}

func resourceSfsTurboDirDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir"
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
		JSONBody:         buildDeleteSfsTurboDirBodyParams(d),
	}

	_, err = sfsClient.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		if utils.IsResourceNotFound(err) {
			return nil
		}
		return diag.Errorf("error deleting SFS Turbo directory: %s", err)
	}

	return nil
}
