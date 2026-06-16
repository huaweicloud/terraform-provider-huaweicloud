package mrs

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

// @API MRS GET /v2/{project_id}/clusters/{cluster_id}/tags/quota
func DataSourceTagsQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTagsQuotaRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the tags quota is located.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster.`,
			},
			"total_quota": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total quota of tags.`,
			},
			"available_quota": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The available quota of tags.`,
			},
		},
	}
}

func getTagsQuota(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/tags/quota"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func dataSourceTagsQuotaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	respBody, err := getTagsQuota(client, clusterId)
	if err != nil {
		return diag.Errorf("error querying tags quota of the cluster (%s): %s", clusterId, err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("total_quota", utils.PathSearch("total_quota", respBody, nil)),
		d.Set("available_quota", utils.PathSearch("available_quota", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
