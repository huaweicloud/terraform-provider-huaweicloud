package dws

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS GET /v1.0/{project_id}/dms/net
func DataSourceHostNets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostNetsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the host nets are located.`,
			},

			// Optional parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The cluster ID to be queried.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The instance name to be queried.`,
			},

			// Attributes.
			"host_nets": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        hostNetSchema(),
				Description: `The list of the host nets that matched filter parameters.`,
			},
		},
	}
}

func hostNetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"virtual_cluster_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The virtual cluster ID.`,
			},
			"ctime": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The query timestamp in Unix milliseconds.`,
			},
			"host_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The host ID.`,
			},
			"host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The host name.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance name.`,
			},
			"interface_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The network interface name.`,
			},
			"up": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the network interface is up.`,
			},
			"speed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The network interface speed in Mbps.`,
			},
			"recv_packets": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The received packets.`,
			},
			"send_packets": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The sent packets.`,
			},
			"recv_drop": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The dropped packets on receiving.`,
			},
			"recv_rate": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The receiving rate in KB/s.`,
			},
			"send_rate": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The sending rate in KB/s.`,
			},
			"io_rate": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The network IO rate in KB/s.`,
			},
		},
	}
}

func buildHostNetsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("cluster_id"); ok {
		res = fmt.Sprintf("%s&cluster_id=%v", res, v)
	}
	if v, ok := d.GetOk("instance_name"); ok {
		res = fmt.Sprintf("%s&instance_name=%v", res, v)
	}

	return res
}

func listHostNets(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1.0/{project_id}/dms/net?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildHostNetsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPathWithLimit, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		hostNets, ok := respBody.([]interface{})
		if !ok {
			return nil, errors.New("error querying host nets: unexpected response type")
		}

		result = append(result, hostNets...)
		if len(hostNets) < limit {
			break
		}
		offset += len(hostNets)
	}

	return result, nil
}

func flattenHostNets(hostNets []interface{}) []map[string]interface{} {
	if len(hostNets) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(hostNets))
	for _, item := range hostNets {
		result = append(result, map[string]interface{}{
			"virtual_cluster_id": utils.PathSearch("virtual_cluster_id", item, nil),
			"ctime":              utils.PathSearch("ctime", item, nil),
			"host_id":            utils.PathSearch("host_id", item, nil),
			"host_name":          utils.PathSearch("host_name", item, nil),
			"instance_name":      utils.PathSearch("instance_name", item, nil),
			"interface_name":     utils.PathSearch("interface_name", item, nil),
			"up":                 utils.PathSearch("up", item, nil),
			"speed":              utils.PathSearch("speed", item, nil),
			"recv_packets":       utils.PathSearch("recv_packets", item, nil),
			"send_packets":       utils.PathSearch("send_packets", item, nil),
			"recv_drop":          utils.PathSearch("recv_drop", item, nil),
			"recv_rate":          utils.PathSearch("recv_rate", item, float64(0)),
			"send_rate":          utils.PathSearch("send_rate", item, float64(0)),
			"io_rate":            utils.PathSearch("io_rate", item, float64(0)),
		})
	}

	return result
}

func dataSourceHostNetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	hostNets, err := listHostNets(client, d)
	if err != nil {
		return diag.Errorf("error querying host nets: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("host_nets", flattenHostNets(hostNets)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
