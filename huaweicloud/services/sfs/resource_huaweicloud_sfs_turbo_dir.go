package sfs

import (
	"context"
	"fmt"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

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
				Description: `Specifies the created file system ID.`,
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the valid full path of an existing directory.`,
			},
			"mude": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the file directory permissions.`,
			},
			"uid": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the user ID of the file directory.`,
			},
			"gid": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the group ID of the file directory.`,
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
		JSONBody:         buildCreateSfsTurboDirBodyParams(d),
		OkCodes: []int{
			204,
		},
	}

	_, err = sfsClient.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SFS file directory: %s", err)
	}

	path := d.Get("path").(string)
	d.SetId(path)

	return resourceSfsTurboDirRead(ctx, d, meta)
}
func buildCreateSfsTurboDirBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"path": d.Get("path"),
		"mude": d.Get("mude"),
		"uid":  d.Get("uid"),
		"gid":  d.Get("gid"),
	}
}

func resourceSfsTurboDirRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	getPath += fmt.Sprintf("?path=%s", d.Get("path"))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := sfsClient.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SFS file directory")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil && getRespBody != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("path", d.Get("path")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDeleteSfsTurboDirBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"path": d.Get("path"),
	}
}

func resourceSfsTurboDirDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("error deleting SFS file directory: %s", err)
	}

	return nil
}
