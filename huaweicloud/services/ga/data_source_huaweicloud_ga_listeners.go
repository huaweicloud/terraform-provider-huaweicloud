package ga

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA GET /v1/listeners
func DataSourceListeners() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceListenersRead,
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the listener.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the listener.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The current status of the listener.",
			},
			"accelerator_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the accelerator to which the listener belongs.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The network transmission protocol type of the listener.",
			},
			"listeners": {
				Type:        schema.TypeList,
				Elem:        listenersSchema(),
				Computed:    true,
				Description: "The list of the listeners.",
			},
		},
	}
}

func listenersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the listener.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the listener.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the listener.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the listener.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network transmission protocol type of the listener.",
			},
			"port_ranges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The listening port range list of the listener.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The listening to start port of the listener.",
						},
						"to_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The listening to end port of the listener.",
						},
					},
				},
			},
			"client_affinity": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The client affinity of the listener.",
			},
			"accelerator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the accelerator to which the listener belongs.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The key/value pairs to associate with the listener.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the listener.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the listener.",
			},
		},
	}
	return &sc
}

func dataSourceListenersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// listListeners: Query the list of listeners
	var (
		listListenersHttpUrl = "v1/listeners"
		listListenersProduct = "ga"
	)
	listListenersClient, err := cfg.NewServiceClient(listListenersProduct, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	listListenersPath := listListenersClient.Endpoint + listListenersHttpUrl

	listListenersqueryParams := buildListListenersQueryParams(d)
	listListenersPath += listListenersqueryParams

	listListenersResp, err := pagination.ListAllItems(
		listListenersClient,
		"marker",
		listListenersPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving listeners")
	}

	listListenersRespJson, err := json.Marshal(listListenersResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listListenersRespBody interface{}
	err = json.Unmarshal(listListenersRespJson, &listListenersRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("listeners", filterListListenersResponseBody(flattenListListenersResponseBody(listListenersRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListListenersResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("listeners", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":              utils.PathSearch("id", v, nil),
			"name":            utils.PathSearch("name", v, nil),
			"status":          utils.PathSearch("status", v, nil),
			"description":     utils.PathSearch("description", v, nil),
			"protocol":        utils.PathSearch("protocol", v, nil),
			"port_ranges":     flattenPortRanges(utils.PathSearch("port_ranges", v, make([]interface{}, 0))),
			"client_affinity": utils.PathSearch("client_affinity", v, nil),
			"accelerator_id":  utils.PathSearch("accelerator_id", v, nil),
			"tags":            utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
			"created_at":      utils.PathSearch("created_at", v, nil),
			"updated_at":      utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func flattenPortRanges(raw interface{}) []map[string]interface{} {
	curArray := raw.([]interface{})
	result := make([]map[string]interface{}, len(curArray))
	for i, portRanges := range curArray {
		result[i] = map[string]interface{}{
			"from_port": utils.PathSearch("from_port", portRanges, nil),
			"to_port":   utils.PathSearch("to_port", portRanges, nil),
		}
	}
	return result
}

func filterListListenersResponseBody(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("protocol"); ok &&
			fmt.Sprint(param) != utils.PathSearch("protocol", v, nil) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListListenersQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("listener_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("accelerator_id"); ok {
		res = fmt.Sprintf("%s&accelerator_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
