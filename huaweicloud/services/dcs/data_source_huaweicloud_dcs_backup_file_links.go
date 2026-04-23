package dcs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS POST /v2/{project_id}/instances/{instance_id}/backups/{backup_id}/links
func DataSourceBackupFileLinks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsBackupFileLinksRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expiration": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"file_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"link": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func dataSourceDcsBackupFileLinksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listBackupFileLinksHttpUrl = "v2/{project_id}/instances/{instance_id}/backups/{backup_id}/links"
		listBackupFileLinksProduct = "dcs"
	)
	listBackupFileLinksClient, err := cfg.NewServiceClient(listBackupFileLinksProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	listBackupFileLinksPath := listBackupFileLinksClient.Endpoint + listBackupFileLinksHttpUrl
	listBackupFileLinksPath = strings.ReplaceAll(
		listBackupFileLinksPath, "{project_id}", listBackupFileLinksClient.ProjectID)
	listBackupFileLinksPath = strings.ReplaceAll(
		listBackupFileLinksPath, "{instance_id}", d.Get("instance_id").(string))
	listBackupFileLinksPath = strings.ReplaceAll(
		listBackupFileLinksPath, "{backup_id}", d.Get("backup_id").(string))

	getBody := map[string]interface{}{
		"expiration": d.Get("expiration").(int),
	}

	getBackupFileLinksOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getBackupFileLinksOpt.JSONBody = utils.RemoveNil(getBody)
	listBackupFileLinksResp, err := listBackupFileLinksClient.Request(
		"POST", listBackupFileLinksPath, &getBackupFileLinksOpt)

	if err != nil {
		return diag.Errorf("error retrieving DCS backup file links: %s", err)
	}

	listBackupFileLinksRespBody, err := utils.FlattenResponse(listBackupFileLinksResp)
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
		d.Set("bucket_name",
			utils.PathSearch("bucket_name", listBackupFileLinksRespBody, nil)),
		d.Set("file_path", utils.PathSearch("file_path", listBackupFileLinksRespBody, nil)),
		d.Set("links", flattenListBackupFileLinksBody(
			utils.PathSearch("links", listBackupFileLinksRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListBackupFileLinksBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		item := map[string]interface{}{
			"file_name": utils.PathSearch("file_name", v, ""),
			"link":      utils.PathSearch("link", v, ""),
		}
		rst = append(rst, item)
	}
	return rst
}
