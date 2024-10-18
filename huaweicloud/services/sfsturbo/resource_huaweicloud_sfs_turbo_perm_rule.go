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

// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/perm-rules
// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/perm-rules/{rule_id}
// @API SFSTurbo PUT /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/perm-rules/{rule_id}
// @API SFSTurbo DELETE /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/perm-rules/{rule_id}
func ResourceSFSTurboPermRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSFSTurboPermRuleCreate,
		ReadContext:   resourceSFSTurboPermRuleRead,
		UpdateContext: resourceSFSTurboPermRuleUpdate,
		DeleteContext: resourceSFSTurboPermRuleDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"share_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rw_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func buildCreateSFSTurboPermRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"rules": buildPermRuleListBodyParam(d),
	}
	return params
}

func buildPermRuleListBodyParam(d *schema.ResourceData) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"ip_cidr":   d.Get("ip_cidr"),
			"rw_type":   d.Get("rw_type"),
			"user_type": d.Get("user_type"),
		},
	}
}

func resourceSFSTurboPermRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/perm-rules"
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
		JSONBody:         buildCreateSFSTurboPermRuleBodyParams(d),
	}

	createResp, err := sfsClient.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo permission rule: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// When the creation interface is called for the first time, the response body in addition to the permission rules
	// added by the user, an additional default permission rule will be added.
	id := utils.PathSearch(fmt.Sprintf("rules[?ip_cidr=='%s'].id|[0]", d.Get("ip_cidr")),
		createRespBody, "").(string)

	if id == "" {
		return diag.Errorf("error creating SFS Turbo permission rule:" +
			" ID is not found in API response")
	}

	d.SetId(id)

	return resourceSFSTurboPermRuleRead(ctx, d, meta)
}

func resourceSFSTurboPermRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/perm-rules/{rule_id}"
	)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	getPath := sfsClient.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", sfsClient.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{share_id}", d.Get("share_id").(string))
	getPath = strings.ReplaceAll(getPath, "{rule_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := sfsClient.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errCode", "SFS.TURBO.0001"),
			"error retrieving SFS Turbo permission rule")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("ip_cidr", utils.PathSearch("ip_cidr", getRespBody, nil)),
		d.Set("rw_type", utils.PathSearch("rw_type", getRespBody, nil)),
		d.Set("user_type", utils.PathSearch("user_type", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateSFSTurboPermRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"rw_type":   d.Get("rw_type"),
		"user_type": d.Get("user_type"),
	}

	return params
}

func resourceSFSTurboPermRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/perm-rules/{rule_id}"
	)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	updatePath := sfsClient.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", sfsClient.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{share_id}", d.Get("share_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{rule_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateSFSTurboPermRuleBodyParams(d),
	}

	_, err = sfsClient.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SFS Turbo permission rule: %s", err)
	}

	return resourceSFSTurboPermRuleRead(ctx, d, meta)
}

func resourceSFSTurboPermRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/perm-rules/{rule_id}"
	)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	deletePath := sfsClient.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", sfsClient.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{share_id}", d.Get("share_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = sfsClient.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting SFS Turbo permission rule: %s", err)
	}

	return nil
}
