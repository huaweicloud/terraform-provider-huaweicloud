package dc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DC GET /v3/{project_id}/dcaas/vif-peer-detections
func DataSourceDcVifPeerDetections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcVifPeerDetectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vif_peer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vif_peer_detections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vif_peer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"loss_ratio": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDcVifPeerDetectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/dcaas/vif-peer-detections"
		product = "dc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listQueryParams := buildListVifPeerDetectionsQueryParams(d)
	listPath += listQueryParams

	listResp, err := pagination.ListAllItems(
		client,
		"marker",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving DC vif peer detections: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("vif_peer_detections", flattenListVifPeerDetectionsBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListVifPeerDetectionsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("vif_peer_detections", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"vif_peer_id": utils.PathSearch("vif_peer_id", v, nil),
			"project_id":  utils.PathSearch("project_id", v, nil),
			"domain_id":   utils.PathSearch("domain_id", v, nil),
			"start_time":  utils.PathSearch("start_time", v, nil),
			"end_time":    utils.PathSearch("end_time", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"loss_ratio":  utils.PathSearch("loss_ratio", v, nil),
		})
	}
	return rst
}

func buildListVifPeerDetectionsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("sort_key"); ok {
		res = fmt.Sprintf("%s&sort_key=%v", res, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		for _, raw := range v.([]interface{}) {
			res = fmt.Sprintf("%s&sort_dir=%v", res, raw.(string))
		}
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
