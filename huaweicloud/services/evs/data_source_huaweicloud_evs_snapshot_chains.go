package evs

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

// @API EVS GET /v5/{project_id}/snapshot-chains
func DataSourceEvsSnapshotChains() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvsSnapshotChainsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_chains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     snapshotChainSchema(),
			},
		},
	}
}

func snapshotChainSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snapshot_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEvsSnapshotChainsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/snapshot-chains"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	var allChains []interface{}
	marker := ""
	for {
		pagedPath := fmt.Sprintf("%s?limit=1000%s", requestPath, buildEvsSnapshotChainsQueryParams(d))
		if marker != "" {
			pagedPath += fmt.Sprintf("&marker=%s", marker)
		}
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		resp, err := client.Request("GET", pagedPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error querying EVS snapshot chains: %s", err)
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}
		chains := utils.PathSearch("snapshot_chains", respBody, []interface{}{}).([]interface{})
		if len(chains) == 0 {
			break
		}
		allChains = append(allChains, chains...)

		marker = utils.PathSearch("[-1].id", chains, "").(string)
		if marker == "" {
			break
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("snapshot_chains", flattenEvsSnapshotChains(allChains)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildEvsSnapshotChainsQueryParams(d *schema.ResourceData) string {
	params := ""
	if v, ok := d.GetOk("id"); ok {
		params += fmt.Sprintf("&id=%v", v)
	}
	if v, ok := d.GetOk("volume_id"); ok {
		params += fmt.Sprintf("&volume_id=%v", v)
	}
	if v, ok := d.GetOk("category"); ok {
		params += fmt.Sprintf("&category=%v", v)
	}
	return params
}

func flattenEvsSnapshotChains(chains []interface{}) []map[string]interface{} {
	if len(chains) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(chains))
	for _, c := range chains {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", c, nil),
			"availability_zone": utils.PathSearch("availability_zone", c, nil),
			"snapshot_count":    utils.PathSearch("snapshot_count", c, nil),
			"capacity":          utils.PathSearch("capacity", c, nil),
			"volume_id":         utils.PathSearch("volume_id", c, nil),
			"category":          utils.PathSearch("category", c, nil),
			"created_at":        utils.PathSearch("created_at", c, nil),
			"updated_at":        utils.PathSearch("updated_at", c, nil),
		})
	}
	return rst
}
