package dcs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS POST /v2/{project_id}/instances/{instance_id}/redislog/{id}/links
func DataSourceDcsRedisLogDownloadLink() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsRedisLogDownloadLinkRead,

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
			"log_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsRedisLogDownloadLinkRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/redislog/{id}/links"
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	logID := d.Get("log_id").(string)

	getRedisLogLinksPath := client.Endpoint + httpUrl
	getRedisLogLinksPath = strings.ReplaceAll(getRedisLogLinksPath, "{project_id}", client.ProjectID)
	getRedisLogLinksPath = strings.ReplaceAll(getRedisLogLinksPath, "{instance_id}", instanceID)
	getRedisLogLinksPath = strings.ReplaceAll(getRedisLogLinksPath, "{id}", logID)

	getRedisLogLinksOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getRedisLogLinksResp, err := client.Request("POST", getRedisLogLinksPath, &getRedisLogLinksOpt)
	if err != nil {
		return diag.Errorf("error retrieving DCS redis log download link: %s", err)
	}

	getRedisLogLinksRespBody, err := utils.FlattenResponse(getRedisLogLinksResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(logID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("link", utils.PathSearch("link", getRedisLogLinksRespBody, nil)),
		d.Set("backup_id", utils.PathSearch("backup_id", getRedisLogLinksRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
