package mrs

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

// @API MRS GET /v2/{project_id}/clusters/{cluster_id}/iam-sync-user
func DataSourceClusterSynchronizedIamObjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterSynchronizedIamObjectsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the synchronized IAM users and user groups are located.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster.`,
			},
			"user_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The synchronized IAM user name list.`,
			},
			"group_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The synchronized IAM user group name list.`,
			},
		},
	}
}

func listClusterSynchronizedIamObjects(client *golangsdk.ServiceClient, clusterID string) (interface{}, error) {
	httpURL := "v2/{project_id}/clusters/{cluster_id}/iam-sync-user"
	listPath := client.Endpoint + httpURL
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{cluster_id}", clusterID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func dataSourceClusterSynchronizedIamObjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	respBody, err := listClusterSynchronizedIamObjects(client, clusterId)
	if err != nil {
		return diag.Errorf("error querying synchronized IAM objects of cluster (%s): %s", clusterId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate data source ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("user_names", utils.PathSearch("user_names", respBody, make([]interface{}, 0))),
		d.Set("group_names", utils.PathSearch("group_names", respBody, make([]interface{}, 0))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
