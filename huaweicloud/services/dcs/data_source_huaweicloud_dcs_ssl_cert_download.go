package dcs

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS POST /v2/{project_id}/instances/{instance_id}/ssl-certs/download
func DataSourceDcsSslCertDownload() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsSslCertDownloadRead,

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
			"file_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsSslCertDownloadRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	httpUrl := "v2/{project_id}/instances/{instance_id}/ssl-certs/download"
	getSslCertDownloadPath := client.Endpoint + httpUrl
	getSslCertDownloadPath = strings.ReplaceAll(getSslCertDownloadPath, "{project_id}", client.ProjectID)
	getSslCertDownloadPath = strings.ReplaceAll(getSslCertDownloadPath, "{instance_id}", instanceID)

	getSslCertDownloadOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getSslCertDownloadResp, err := client.Request("POST", getSslCertDownloadPath, &getSslCertDownloadOpt)
	if err != nil {
		return diag.Errorf("error creating DCS ssl cert download link: %s", err)
	}

	getSslCertDownloadRespBody, err := utils.FlattenResponse(getSslCertDownloadResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("file_name", utils.PathSearch("file_name", getSslCertDownloadRespBody, nil)),
		d.Set("link", utils.PathSearch("link", getSslCertDownloadRespBody, nil)),
		d.Set("bucket_name", utils.PathSearch("bucket_name", getSslCertDownloadRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
