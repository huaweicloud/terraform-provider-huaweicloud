package sfsturbo

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir-usage
func DataSourceDirusage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDirUsageRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"share_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dir_usage": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"used_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDirUsageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		shareId = d.Get("share_id").(string)
		dirPath = d.Get("path").(string)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/dir-usage"
	)

	client, err := cfg.NewServiceClient("sfs-turbo", region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{share_id}", shareId)
	getPath = fmt.Sprintf("%s?path=%s", getPath, dirPath)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving directory used information: %s", err)
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("dir_usage", flattenDirectoryInfo(utils.PathSearch("dir_usage", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDirectoryInfo(dirCapacity interface{}) []map[string]interface{} {
	if dirCapacity == nil {
		return nil
	}

	result := map[string]interface{}{
		"used_capacity": utils.PathSearch("used_capacity", dirCapacity, nil),
	}

	return []map[string]interface{}{result}
}
