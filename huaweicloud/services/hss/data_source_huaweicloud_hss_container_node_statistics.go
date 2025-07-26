package hss

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

// @API HSS GET /v5/{project_id}/container/node-statistics
func DataSourceContainerNodeStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerNodeStatisticsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Attribute.
			"unprotected_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protected_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protected_num_on_demand": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protected_num_packet_cycle": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_node_not_installed_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"not_cluster_node_not_installed_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceContainerNodeStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/container/node-statistics"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		requestPath = fmt.Sprintf("%s?enterprise_project_id=%v", requestPath, epsId)
	}
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS container node statistics: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("unprotected_num", utils.PathSearch("unprotected_num", respBody, nil)),
		d.Set("protected_num", utils.PathSearch("protected_num", respBody, nil)),
		d.Set("protected_num_on_demand", utils.PathSearch("protected_num_on_demand", respBody, nil)),
		d.Set("protected_num_packet_cycle", utils.PathSearch("protected_num_packet_cycle", respBody, nil)),
		d.Set("cluster_node_not_installed_num", utils.PathSearch("cluster_node_not_installed_num", respBody, nil)),
		d.Set("not_cluster_node_not_installed_num", utils.PathSearch("not_cluster_node_not_installed_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
