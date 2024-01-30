// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDM
// ---------------------------------------------------------------

package ddm

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDM GET /v1/{project_id}/instances/{instance_id}/nodes
func DataSourceDdmInstanceNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdmInstanceNodesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of DDM instance.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Elem:        InstanceNodesNodeSchema(),
				Computed:    true,
				Description: `Indicates the list of DDM instance node.`,
			},
		},
	}
}

func InstanceNodesNodeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the DDM instance node.`,
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the IP address of the DDM instance node.`,
			},
			"port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the port of the DDM instance node.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the DDM instance node.`,
			},
		},
	}
	return &sc
}

func resourceDdmInstanceNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDdmInstanceNodes: Query the list of DDM instance nodes
	var (
		getDdmInstanceNodesHttpUrl = "v1/{project_id}/instances/{instance_id}/nodes"
		getDdmInstanceNodesProduct = "ddm"
	)
	getDdmInstanceNodesClient, err := cfg.NewServiceClient(getDdmInstanceNodesProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM Client: %s", err)
	}

	getDdmInstanceNodesPath := getDdmInstanceNodesClient.Endpoint + getDdmInstanceNodesHttpUrl
	getDdmInstanceNodesPath = strings.ReplaceAll(getDdmInstanceNodesPath, "{project_id}", getDdmInstanceNodesClient.ProjectID)
	getDdmInstanceNodesPath = strings.ReplaceAll(getDdmInstanceNodesPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	getDdmInstanceNodesResp, err := pagination.ListAllItems(
		getDdmInstanceNodesClient,
		"offset",
		getDdmInstanceNodesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DdmInstanceNodes")
	}

	getDdmInstanceNodesRespJson, err := json.Marshal(getDdmInstanceNodesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getDdmInstanceNodesRespBody interface{}
	err = json.Unmarshal(getDdmInstanceNodesRespJson, &getDdmInstanceNodesRespBody)
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
		d.Set("nodes", flattenGetInstanceNodesResponseBodyNode(getDdmInstanceNodesRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetInstanceNodesResponseBodyNode(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("nodes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":     utils.PathSearch("node_id", v, nil),
			"ip":     utils.PathSearch("ip", v, nil),
			"port":   utils.PathSearch("port", v, nil),
			"status": utils.PathSearch("status", v, nil),
		})
	}
	return rst
}
