// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CDM
// ---------------------------------------------------------------

package cdm

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

// @API CDM GET /v1.1/{project_id}/clusters
func DataSourceCdmClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceCdmClustersRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Cluster name.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The AZ name.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Cluster status.`,
			},
			"clusters": {
				Type:        schema.TypeList,
				Elem:        cdmClustersClusterSchema(),
				Computed:    true,
				Description: `The list of clusters.`,
			},
		},
	}
}

func cdmClustersClusterSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `cluster ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `cluster name.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The AZ name.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster version.`,
			},
			"is_auto_off": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable auto shutdown.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster status.`,
			},
			"recent_event": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of events.`,
			},
			"public_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `EIP bound to the cluster.`,
			},
			"is_frozen": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether the cluster is frozen. The value can be 0 (not frozen) or 1 (frozen).`,
			},
			"is_failure_remind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether to notifications when a table/file migration job fails or an EIP exception occurs.`,
			},
			"instances": {
				Type:        schema.TypeList,
				Elem:        cdmClustersClusterInstanceSchema(),
				Computed:    true,
				Description: `The list of instance nodes.`,
			},
		},
	}
	return &sc
}

func cdmClustersClusterInstanceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Instance ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Instance name.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Private IP.`,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Public IP.`,
			},
			"manage_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Management IP address.`,
			},
			"traffic_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Traffic IP.`,
			},
			"role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Instance role.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Instance type.`,
			},
			"is_frozen": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether the node is frozen. The value can be 0 (not frozen) or 1 (frozen).`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Instance node status.`,
			},
		},
	}
	return &sc
}

func resourceCdmClustersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listClusters: Query the list of CDM clusters
	var (
		listClustersHttpUrl = "v1.1/{project_id}/clusters"
		listClustersProduct = "cdm"
	)
	listClustersClient, err := cfg.NewServiceClient(listClustersProduct, region)
	if err != nil {
		return diag.Errorf("error creating CDM client: %s", err)
	}

	listClustersPath := listClustersClient.Endpoint + listClustersHttpUrl
	listClustersPath = strings.ReplaceAll(listClustersPath, "{project_id}", listClustersClient.ProjectID)

	listClustersOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	listClustersResp, err := listClustersClient.Request("GET", listClustersPath, &listClustersOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CdmClusters")
	}

	listClustersRespBody, err := utils.FlattenResponse(listClustersResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("clusters", filterListClustersCluster(
			flattenListClustersCluster(listClustersRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListClustersCluster(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("clusters", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"availability_zone": utils.PathSearch("azName", v, nil),
			"version":           utils.PathSearch("datastore.version", v, nil),
			"is_auto_off":       utils.PathSearch("isAutoOff", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"recent_event":      utils.PathSearch("recentEvent", v, nil),
			"public_endpoint":   utils.PathSearch("publicEndpoint", v, nil),
			"is_frozen":         utils.PathSearch("isFrozen", v, nil),
			"is_failure_remind": utils.PathSearch("customerConfig.failureRemind", v, nil),
			"instances":         flattenClusterInstance(v),
		})
	}
	return rst
}

func flattenClusterInstance(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"name":       utils.PathSearch("name", v, nil),
			"private_ip": utils.PathSearch("private_ip", v, nil),
			"public_ip":  utils.PathSearch("publicIp", v, nil),
			"manage_ip":  utils.PathSearch("manageIp", v, nil),
			"traffic_ip": utils.PathSearch("trafficIp", v, nil),
			"role":       utils.PathSearch("role", v, nil),
			"type":       utils.PathSearch("type", v, nil),
			"is_frozen":  utils.PathSearch("isFrozen", v, nil),
			"status":     utils.PathSearch("status", v, nil),
		})
	}
	return rst
}

func filterListClustersCluster(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("name"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("name", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("availability_zone"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("availability_zone", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("status"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("status", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}
