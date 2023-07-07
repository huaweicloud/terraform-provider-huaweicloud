// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DCS
// ---------------------------------------------------------------

package dcs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDcsInstanceShard() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDcsInstanceShardRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the DCS instance.`,
			},
			"shard_names": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the list of the shard names.`,
			},
			"replica_ips": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the list of the replica ips.`,
			},
			"replica_role": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the role of the replica.`,
			},
			"shards": {
				Type:        schema.TypeList,
				Elem:        InstanceShardShardSchema(),
				Computed:    true,
				Description: `Indicates the list of DCS instance replicas.`,
			},
		},
	}
}

func InstanceShardShardSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"shard_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the shard.`,
			},
			"shard_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the shard.`,
			},
			"replicas": {
				Type:        schema.TypeList,
				Elem:        InstanceShardShardReplicaSchema(),
				Computed:    true,
				Description: `Indicates the list of replicas in the shard.`,
			},
		},
	}
	return &sc
}

func InstanceShardShardReplicaSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the replica.`,
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the IP of the replica.`,
			},
			"role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the role of the replica.`,
			},
			"node_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the node.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the replica.`,
			},
		},
	}
	return &sc
}

func resourceDcsInstanceShardRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDCSInstanceShards: Query the List of DCS instance shards.
	var (
		getDCSInstanceShardsHttpUrl = "v2/{project_id}/instance/{instance_id}/groups"
		getDCSInstanceShardsProduct = "dcs"
	)
	getDCSInstanceShardsClient, err := cfg.NewServiceClient(getDCSInstanceShardsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS Client: %s", err)
	}

	getDCSInstanceShardsPath := getDCSInstanceShardsClient.Endpoint + getDCSInstanceShardsHttpUrl
	getDCSInstanceShardsPath = strings.ReplaceAll(getDCSInstanceShardsPath, "{project_id}",
		getDCSInstanceShardsClient.ProjectID)
	getDCSInstanceShardsPath = strings.ReplaceAll(getDCSInstanceShardsPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))

	getDCSInstanceShardsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getDCSInstanceShardsResp, err := getDCSInstanceShardsClient.Request("GET", getDCSInstanceShardsPath,
		&getDCSInstanceShardsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DcsInstanceShard")
	}

	getDCSInstanceShardsRespBody, err := utils.FlattenResponse(getDCSInstanceShardsResp)
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
		d.Set("shards", flattenGetDCSInstanceShardsResponseBodyShard(d, getDCSInstanceShardsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetDCSInstanceShardsResponseBodyShard(d *schema.ResourceData, resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	shardNames := d.Get("shard_names").([]interface{})
	shardNamesMap := make(map[string]bool)
	for _, shardName := range shardNames {
		shardNamesMap[shardName.(string)] = true
	}
	curJson := utils.PathSearch("group_list", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		shardName := utils.PathSearch("group_name", v, "")
		if len(shardNamesMap) == 0 || shardNamesMap[shardName.(string)] {
			rst = append(rst, map[string]interface{}{
				"shard_id":   utils.PathSearch("group_id", v, nil),
				"shard_name": utils.PathSearch("group_name", v, nil),
				"replicas":   flattenShardReplicas(d, v),
			})
		}
	}
	return rst
}

func flattenShardReplicas(d *schema.ResourceData, resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	replicaIps := d.Get("replica_ips").([]interface{})
	replicaIpsMap := make(map[string]bool)
	for _, replicaIp := range replicaIps {
		replicaIpsMap[replicaIp.(string)] = true
	}
	role := d.Get("replica_role").(string)
	curJson := utils.PathSearch("replication_list", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		ip := utils.PathSearch("replication_ip", v, "")
		if len(replicaIpsMap) > 0 && !replicaIpsMap[ip.(string)] {
			continue
		}
		replicationRole := utils.PathSearch("replication_role", v, "")
		if role != "" && role != replicationRole {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"id":      utils.PathSearch("replication_id", v, nil),
			"ip":      ip,
			"role":    replicationRole,
			"node_id": utils.PathSearch("node_id", v, nil),
			"status":  utils.PathSearch("status", v, nil),
		})
	}
	return rst
}
