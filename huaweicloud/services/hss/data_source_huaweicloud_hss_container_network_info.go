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

// @API HSS GET /v5/{project_id}/container-network/{cluster_id}/network-info
func DataSourceContainerNetworkInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerNetworkInfoRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidrs": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kube_proxy_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_support_egress": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildContainerNetworkInfoQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func dataSourceContainerNetworkInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/container-network/{cluster_id}/network-info"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", d.Get("cluster_id").(string))
	requestPath += buildContainerNetworkInfoQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS container network info: %s", err)
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

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("mode", utils.PathSearch("mode", respBody, nil)),
		d.Set("vpc", utils.PathSearch("vpc", respBody, nil)),
		d.Set("subnet", utils.PathSearch("subnet", respBody, nil)),
		d.Set("security_group", utils.PathSearch("security_group", respBody, nil)),
		d.Set("ipv4_cidr", utils.PathSearch("ipv4_cidr", respBody, nil)),
		d.Set("cidrs", utils.PathSearch("cidrs", respBody, nil)),
		d.Set("kube_proxy_mode", utils.PathSearch("kube_proxy_mode", respBody, nil)),
		d.Set("is_support_egress", utils.PathSearch("is_support_egress", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
