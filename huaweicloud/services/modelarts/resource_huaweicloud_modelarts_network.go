// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v1/{project_id}/networks/{id}
// @API ModelArts PATCH /v1/{project_id}/networks/{id}
// @API ModelArts DELETE /v1/{project_id}/networks/{id}
// @API ModelArts POST /v1/{project_id}/networks
func ResourceModelartsNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModelartsNetworkCreate,
		UpdateContext: resourceModelartsNetworkUpdate,
		ReadContext:   resourceModelartsNetworkRead,
		DeleteContext: resourceModelartsNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of network.`,
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Network CIDR.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Workspace ID, which defaults to 0.`,
			},
			"peer_connections": {
				Type:        schema.TypeList,
				Elem:        modelartsNetworkPeerConnectionSchema(),
				Optional:    true,
				Description: `List of networks that can be connected in peer mode.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of network.`,
			},
		},
	}
}

func modelartsNetworkPeerConnectionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `ID of the peer VPC.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `ID of the peer subnet.`,
			},
		},
	}
	return &sc
}

func resourceModelartsNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createNetwork: create a Modelarts network.
	var (
		createNetworkHttpUrl = "v1/{project_id}/networks"
		createNetworkProduct = "modelarts"
	)
	createNetworkClient, err := cfg.NewServiceClient(createNetworkProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	createNetworkPath := createNetworkClient.Endpoint + createNetworkHttpUrl
	createNetworkPath = strings.ReplaceAll(createNetworkPath, "{project_id}", createNetworkClient.ProjectID)

	createNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createNetworkOpt.JSONBody = utils.RemoveNil(buildCreateNetworkBodyParams(d))
	createNetworkResp, err := createNetworkClient.Request("POST", createNetworkPath, &createNetworkOpt)
	if err != nil {
		return diag.Errorf("error creating Modelarts network: %s", err)
	}

	createNetworkRespBody, err := utils.FlattenResponse(createNetworkResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("metadata.name", createNetworkRespBody, nil)
	if id == nil {
		return diag.Errorf("error creating Modelarts network: ID is not found in API response")
	}
	d.SetId(id.(string))

	err = createNetworkWaitingForStateCompleted(ctx, createNetworkClient, id.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts network (%s) creation to complete: %s", d.Id(), err)
	}
	return resourceModelartsNetworkRead(ctx, d, meta)
}

func buildCreateNetworkBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Network",
		"metadata": map[string]interface{}{
			"labels": map[string]interface{}{
				"os.modelarts/name":         d.Get("name"),
				"os.modelarts/workspace.id": utils.ValueIgnoreEmpty(d.Get("workspace_id")),
			},
		},
		"spec": map[string]interface{}{
			"cidr": d.Get("cidr"),
			"connection": map[string]interface{}{
				"peerConnectionList": buildNetworkRequestBodyPeerConnection(d.Get("peer_connections")),
			},
		},
	}
	return bodyParams
}

func buildNetworkRequestBodyPeerConnection(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return []map[string]interface{}{}
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"peerVpcId":    utils.ValueIgnoreEmpty(raw["vpc_id"]),
				"peerSubnetId": utils.ValueIgnoreEmpty(raw["subnet_id"]),
			}
		}
		return rst
	}
	return []map[string]interface{}{}
}

func createNetworkWaitingForStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, networkId string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			getModelartsNetworkRespBody, err := getModelartsNetwork(client, networkId)
			if err != nil {
				return nil, "ERROR", err
			}

			statusRaw := utils.PathSearch(`status.phase`, getModelartsNetworkRespBody, "")

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"Active",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return getModelartsNetworkRespBody, "COMPLETED", nil
			}

			pendingStatus := []string{
				"Creating", "",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return getModelartsNetworkRespBody, "PENDING", nil
			}

			return getModelartsNetworkRespBody, status, nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceModelartsNetworkRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}
	// getModelartsNetwork: Query the Modelarts network.
	getModelartsNetworkRespBody, err := getModelartsNetwork(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Modelarts network")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch(`metadata.labels."os.modelarts/name"`, getModelartsNetworkRespBody, nil)),
		d.Set("workspace_id", utils.PathSearch(`metadata.labels."os.modelarts/workspace.id"`, getModelartsNetworkRespBody, nil)),
		d.Set("cidr", utils.PathSearch("spec.cidr", getModelartsNetworkRespBody, nil)),
		d.Set("peer_connections", flattenGetNetworkResponseBodyPeerConnection(getModelartsNetworkRespBody)),
		d.Set("status", utils.PathSearch("status.phase", getModelartsNetworkRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getModelartsNetwork(client *golangsdk.ServiceClient, networkId string) (interface{}, error) {
	var (
		getModelartsNetworkHttpUrl = "v1/{project_id}/networks/{id}"
	)

	getModelartsNetworkPath := client.Endpoint + getModelartsNetworkHttpUrl
	getModelartsNetworkPath = strings.ReplaceAll(getModelartsNetworkPath, "{project_id}", client.ProjectID)
	getModelartsNetworkPath = strings.ReplaceAll(getModelartsNetworkPath, "{id}", networkId)

	getModelartsNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getModelartsNetworkResp, err := client.Request("GET", getModelartsNetworkPath, &getModelartsNetworkOpt)

	if err != nil {
		return nil, err
	}

	getModelartsNetworkRespBody, err := utils.FlattenResponse(getModelartsNetworkResp)
	if err != nil {
		return nil, err
	}
	return getModelartsNetworkRespBody, nil
}

func flattenGetNetworkResponseBodyPeerConnection(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("spec.connection.peerConnectionList", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"vpc_id":    utils.PathSearch("peerVpcId", v, nil),
			"subnet_id": utils.PathSearch("peerSubnetId", v, nil),
		})
	}
	return rst
}

func resourceModelartsNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		networkId = d.Id()
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	updateNetworkChanges := []string{
		"peer_connections",
	}

	if d.HasChanges(updateNetworkChanges...) {
		// updateNetwork: update the ModelArts network.
		updateNetworkHttpUrl := "v1/{project_id}/networks/{id}"

		updateNetworkPath := client.Endpoint + updateNetworkHttpUrl
		updateNetworkPath = strings.ReplaceAll(updateNetworkPath, "{project_id}", client.ProjectID)
		updateNetworkPath = strings.ReplaceAll(updateNetworkPath, "{id}", networkId)

		updateNetworkOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/merge-patch+json"},
		}

		updateNetworkOpt.JSONBody = buildUpdateNetworkBodyParams(d)
		_, err = client.Request("PATCH", updateNetworkPath, &updateNetworkOpt)
		if err != nil {
			return diag.Errorf("error updating Modelarts network: %s", err)
		}
		err = updateNetworkWaitingForStateCompleted(ctx, client, networkId, d.Get("peer_connections").([]interface{}),
			d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Modelarts network (%s) update to complete: %s", d.Id(), err)
		}
	}
	return resourceModelartsNetworkRead(ctx, d, meta)
}

func buildUpdateNetworkBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"spec": map[string]interface{}{
			"connection": map[string]interface{}{
				"peerConnectionList": buildNetworkRequestBodyPeerConnection(d.Get("peer_connections")),
			},
		},
	}
	return bodyParams
}

func updateNetworkWaitingForStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, networkId string, peerConnections []interface{},
	t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			getModelartsNetworkRespBody, err := getModelartsNetwork(client, networkId)
			if err != nil {
				return nil, "ERROR", err
			}

			// if peer_connections is empty, then check status.connectionStatus.peerConnectionStatus is empty
			if len(peerConnections) == 0 {
				if utils.PathSearch(`length(status.connectionStatus.peerConnectionStatus)`, getModelartsNetworkRespBody, float64(0)).(float64) > 0 {
					return getModelartsNetworkRespBody, "PENDING", nil
				}

				return getModelartsNetworkRespBody, "COMPLETED", nil
			}

			// if peer_connections is not empty, then check those connections are all
			// in `status.connectionStatus.peerConnectionStatus` and the status is `Active`
			for _, v := range peerConnections {
				raw := v.(map[string]interface{})
				peerVpcId := utils.ValueIgnoreEmpty(raw["vpc_id"])
				peerSubnetId := utils.ValueIgnoreEmpty(raw["subnet_id"])

				searchAbnormalJsonPath := fmt.Sprintf(`length(status.connectionStatus.peerConnectionStatus[?peerVpcId=='%s'
                     && peerSubnetId=='%s' && phase=='Abnormal'])`, peerVpcId, peerSubnetId)
				if utils.PathSearch(searchAbnormalJsonPath, getModelartsNetworkRespBody, float64(0)).(float64) > 0 {
					return nil, "ERROR", fmt.Errorf("error updating peer_connections, vpc_id: %s, subnet_id: %s", peerVpcId, peerSubnetId)
				}

				searchActiveJsonPath := fmt.Sprintf(`length(status.connectionStatus.peerConnectionStatus[?peerVpcId=='%s'
                     && peerSubnetId=='%s' && phase=='Active'])`, peerVpcId, peerSubnetId)
				if utils.PathSearch(searchActiveJsonPath, getModelartsNetworkRespBody, float64(0)).(float64) == 0 {
					return getModelartsNetworkRespBody, "PENDING", nil
				}
			}

			return getModelartsNetworkRespBody, "COMPLETED", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceModelartsNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// deleteNetwork: delete Modelarts network
	var (
		deleteNetworkHttpUrl = "v1/{project_id}/networks/{id}"
		deleteNetworkProduct = "modelarts"
		cfg                  = meta.(*config.Config)
		region               = cfg.GetRegion(d)
		networkId            = d.Id()
	)

	deleteNetworkClient, err := cfg.NewServiceClient(deleteNetworkProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	deleteNetworkPath := deleteNetworkClient.Endpoint + deleteNetworkHttpUrl
	deleteNetworkPath = strings.ReplaceAll(deleteNetworkPath, "{project_id}", deleteNetworkClient.ProjectID)
	deleteNetworkPath = strings.ReplaceAll(deleteNetworkPath, "{id}", networkId)

	deleteNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteNetworkClient.Request("DELETE", deleteNetworkPath, &deleteNetworkOpt)
	if err != nil {
		return diag.Errorf("error deleting Modelarts network: %s", err)
	}

	err = deleteNetworkWaitingForStateCompleted(ctx, deleteNetworkClient, networkId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts network (%s) deletion to complete: %s", networkId, err)
	}
	return nil
}

func deleteNetworkWaitingForStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, networkId string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			_, err := getModelartsNetwork(client, networkId)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					var obj = map[string]string{"code": "COMPLETED"}
					return obj, "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			return nil, "PENDING", nil
		},
		Timeout:      t,
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
