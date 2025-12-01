package esw

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

// @API ESW GET /v3/{project_id}/l2cg/instances/{instance_id}/connections
func DataSourceEswConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEswConnectionsRead,

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
			"connection_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fixed_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"remote_infos": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     remoteInfosSchema(),
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virsubnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
				},
			},
		},
	}
}

func remoteInfosSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"segmentation_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tunnel_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tunnel_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceEswConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/l2cg/instances/{instance_id}/connections"
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ESW client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	listQueryParams := buildListConnectionsQueryParams(d)
	listPath += listQueryParams

	listResp, err := pagination.ListAllItems(
		client,
		"marker",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving ESW connections: %s", err)
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

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("connections", flattenEswConnectionsBody(listRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListConnectionsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("connection_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenEswConnectionsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("connections", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":           utils.PathSearch("id", v, nil),
			"instance_id":  utils.PathSearch("instance_id", v, nil),
			"name":         utils.PathSearch("name", v, nil),
			"project_id":   utils.PathSearch("project_id", v, nil),
			"fixed_ips":    utils.PathSearch("fixed_ips", v, nil),
			"remote_infos": flattenConnectionsRemoteInfos(v),
			"vpc_id":       utils.PathSearch("vpc_id", v, nil),
			"virsubnet_id": utils.PathSearch("virsubnet_id", v, nil),
			"status":       utils.PathSearch("status", v, nil),
			"created_at":   utils.PathSearch("created_at", v, nil),
			"updated_at":   utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func flattenConnectionsRemoteInfos(resp interface{}) []interface{} {
	curJson := utils.PathSearch("remote_infos", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"segmentation_id": utils.PathSearch("segmentation_id", v, nil),
			"tunnel_ip":       utils.PathSearch("tunnel_ip", v, nil),
			"tunnel_port":     utils.PathSearch("tunnel_port", v, nil),
			"tunnel_type":     utils.PathSearch("tunnel_type", v, nil),
		})
	}
	return rst
}
