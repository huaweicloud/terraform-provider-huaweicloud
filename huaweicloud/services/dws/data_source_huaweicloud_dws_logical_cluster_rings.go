// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DWS
// ---------------------------------------------------------------

package dws

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/logical-clusters
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/logical-clusters/rings
func DataSourceLogicalClusterRings() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceLogicalClusterRingsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the DWS cluster ID.`,
			},
			"cluster_rings": {
				Type:        schema.TypeList,
				Elem:        clusterRingsSchema(),
				Computed:    true,
				Description: `Indicates the cluster ring list information.`,
			},
		},
	}
}

func clusterRingsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"is_available": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the cluster host ring is available.`,
			},
			"ring_hosts": {
				Type:        schema.TypeList,
				Elem:        ringHostsSchema(),
				Computed:    true,
				Description: `Indicates the cluster host ring list information.`,
			},
		},
	}
	return &sc
}

func ringHostsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the host name.`,
			},
			"back_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the backend IP address.`,
			},
			"cpu_cores": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of CPU cores.`,
			},
			"memory": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `Indicates the host memory.`,
			},
			"disk_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `Indicates the host disk size.`,
			},
		},
	}
	return &sc
}

func readLogicalClusters(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getPath := client.Endpoint + "v2/{project_id}/clusters/{cluster_id}/logical-clusters"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", d.Get("cluster_id").(string))
	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, err
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return nil, err
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return nil, err
	}
	return getRespBody, nil
}

func readLogicalClusterRings(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getPath := client.Endpoint + "v2/{project_id}/clusters/{cluster_id}/logical-clusters/rings"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", d.Get("cluster_id").(string))
	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, err
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return nil, err
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return nil, err
	}
	return getRespBody, nil
}

func resourceLogicalClusterRingsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr    *multierror.Error
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dws"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	clusterRespBody, err := readLogicalClusters(client, d)
	if err != nil {
		return diag.Errorf("error retrieving DWS logical clusters: %s", err)
	}
	logicalClusterRings := flattenLogicalClusterRings(clusterRespBody)
	if len(logicalClusterRings) == 0 {
		// When the logical cluster list API cannot query the data, try to obtain it from logical cluster ring list API.
		clusterRingRespBody, err := readLogicalClusterRings(client, d)
		if err != nil {
			return diag.Errorf("error retrieving DWS logical cluster rings: %s", err)
		}
		logicalClusterRings = flattenClusterRings(clusterRingRespBody, true)
	}
	sort.Sort(RingsSlice(logicalClusterRings))

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("cluster_rings", logicalClusterRings),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLogicalClusterRings(resp interface{}) []interface{} {
	var rst []interface{}

	availableExpression := "logical_clusters[?logical_cluster_name == 'elastic_group']"
	availableCurJson := utils.PathSearch(availableExpression, resp, make([]interface{}, 0))
	availableCurArray := availableCurJson.([]interface{})
	for _, v := range availableCurArray {
		rst = append(rst, flattenClusterRings(v, true)...)
	}

	unAvailableExpression := "logical_clusters[?logical_cluster_name != 'elastic_group']"
	unAvailableCurJson := utils.PathSearch(unAvailableExpression, resp, make([]interface{}, 0))
	unAvailableCurArray := unAvailableCurJson.([]interface{})
	for _, v := range unAvailableCurArray {
		rst = append(rst, flattenClusterRings(v, false)...)
	}
	return rst
}

func flattenClusterRings(resp interface{}, isAvailable bool) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("cluster_rings", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"ring_hosts":   flattenRingHosts(v),
			"is_available": isAvailable,
		}
	}
	return rst
}

func flattenRingHosts(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("ring_hosts", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"host_name": utils.PathSearch("host_name", v, nil),
			"back_ip":   utils.PathSearch("back_ip", v, nil),
			"cpu_cores": utils.PathSearch("cpu_cores", v, nil),
			"memory":    utils.PathSearch("memory", v, nil),
			"disk_size": utils.PathSearch("disk_size", v, nil),
		}
	}
	return rst
}

type RingsSlice []interface{}

func (x RingsSlice) Len() int {
	return len(x)
}

// Less Sort by comparing the largest IP value in the array
func (x RingsSlice) Less(i, j int) bool {
	ringHostsI := utils.PathSearch("ring_hosts", x[i], make([]interface{}, 0)).([]interface{})
	ringHostsJ := utils.PathSearch("ring_hosts", x[j], make([]interface{}, 0)).([]interface{})

	return searchMaxBackIp(ringHostsI) > searchMaxBackIp(ringHostsJ)
}

func searchMaxBackIp(ringHosts []interface{}) string {
	var maxBackIp string
	for _, v := range ringHosts {
		backIp := utils.PathSearch("back_ip", v, "").(string)
		if backIp > maxBackIp {
			maxBackIp = backIp
		}
	}
	return maxBackIp
}

func (x RingsSlice) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
