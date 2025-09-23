package taurusdb

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

// @API GaussDBforMySQL GET /v3/{project_id}/instance/{instance_id}/auditlog/download-link
func DataSourceGaussDBMysqlAuditLogDownloadLinks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBMysqlAuditLogDownloadLinksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the GaussDB MySQL instance.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.`,
			},
			"node_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the GaussDB MySQL instance node.`,
			},
			"links": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of the full SQL file information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the file.`,
						},
						"full_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the full name of the file.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the file size, in KB.`,
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the last modification time of the SQL file.`,
						},
						"download_link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the link for downloading the file.`,
						},
						"link_expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the link expiration time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceGaussDBMysqlAuditLogDownloadLinksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instance/{instance_id}/auditlog/download-link"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", d.Get("instance_id").(string))
	getBasePath += buildGetAuditLogDownloadLinksParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	lastFileName := ""
	limit := 50
	res := make([]map[string]interface{}, 0)
	for {
		getPath := getBasePath + buildGetAuditLogDownloadLinksPageParams(lastFileName, limit)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving GaussDB MySQL audit log download links: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		links, nextLastFileName := flattenGaussDBMysqlAuditLogDownloadLinks(getRespBody)
		res = append(res, links...)
		if len(links) < limit {
			break
		}
		lastFileName = nextLastFileName
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("links", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAuditLogDownloadLinksParams(d *schema.ResourceData) string {
	instanceId := d.Get("instance_id").(string)
	startTime := d.Get("start_time").(string)
	endTime := d.Get("end_time").(string)
	res := fmt.Sprintf("?instance_id=%s&start_time=%s&end_time=%s", instanceId, startTime, endTime)
	if v, ok := d.GetOk("node_id"); ok {
		res = fmt.Sprintf("%s&node_id=%s", res, v.(string))
	}
	return res
}

func buildGetAuditLogDownloadLinksPageParams(lastFileName string, limit int) string {
	return fmt.Sprintf("&limit=%d&last_file_name=%s", limit, lastFileName)
}

func flattenGaussDBMysqlAuditLogDownloadLinks(resp interface{}) ([]map[string]interface{}, string) {
	auditLogDownloadLinksJson := utils.PathSearch("files", resp, make([]interface{}, 0))
	auditLogDownloadLinksArray := auditLogDownloadLinksJson.([]interface{})
	if len(auditLogDownloadLinksArray) == 0 {
		return nil, ""
	}

	result := make([]map[string]interface{}, 0, len(auditLogDownloadLinksArray))
	var lastFileName string
	for _, auditLogDownloadLink := range auditLogDownloadLinksArray {
		result = append(result, map[string]interface{}{
			"name":              utils.PathSearch("name", auditLogDownloadLink, nil),
			"full_name":         utils.PathSearch("full_name", auditLogDownloadLink, nil),
			"size":              utils.PathSearch("size", auditLogDownloadLink, nil),
			"updated_time":      utils.PathSearch("updated_time", auditLogDownloadLink, nil),
			"download_link":     utils.PathSearch("download_link", auditLogDownloadLink, nil),
			"link_expired_time": utils.PathSearch("link_expired_time", auditLogDownloadLink, nil),
		})
		lastFileName = utils.PathSearch("name", auditLogDownloadLink, "").(string)
	}
	return result, lastFileName
}
