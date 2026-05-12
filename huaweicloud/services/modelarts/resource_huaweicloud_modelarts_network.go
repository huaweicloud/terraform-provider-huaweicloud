package modelarts

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var networkNonUpdatableParams = []string{
	"name",
	"cidr",
	"workspace_id",
}

// @API ModelArts POST /v1/{project_id}/networks
// @API ModelArts GET /v1/{project_id}/networks/{id}
// @API ModelArts PATCH /v1/{project_id}/networks/{id}
// @API ModelArts DELETE /v1/{project_id}/networks/{id}
func ResourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(networkNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the network is located.`,
			},

			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the network.`,
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The CIDR of the network.`,
			},

			// Optional parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the workspace to which the network belongs.`,
			},
			"peer_connections": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        networkPeeringConnectionSchema(),
				Description: `The list of networks that can be connected in peering mode.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the network.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true,
						Required: true,
					}),
			},
		},
	}
}

func networkPeeringConnectionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the VPC to which the peering connection belongs.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the subnet to which the peering connection belongs.`,
			},
		},
	}
}

func buildNetworkPeerConnections(peeringConncetions []interface{}) []map[string]interface{} {
	if len(peeringConncetions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(peeringConncetions))
	for _, peeringConnection := range peeringConncetions {
		result = append(result, map[string]interface{}{
			"peerVpcId":    utils.ValueIgnoreEmpty(utils.PathSearch("vpc_id", peeringConnection, "")),
			"peerSubnetId": utils.ValueIgnoreEmpty(utils.PathSearch("subnet_id", peeringConnection, "")),
		})
	}

	return result
}

func buildCreateNetworkBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
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
				"peerConnectionList": buildNetworkPeerConnections(d.Get("peer_connections").([]interface{})),
			},
		},
	}
}

func createNetwork(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpURL := "v1/{project_id}/networks"

	createPath := client.Endpoint + httpURL
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateNetworkBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(createResp)
}

func GetNetworkById(client *golangsdk.ServiceClient, networkId string) (interface{}, error) {
	httpURL := "v1/{project_id}/networks/{id}"

	getPath := client.Endpoint + httpURL
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", networkId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func refreshNetworkStatusFunc(client *golangsdk.ServiceClient, networkId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetNetworkById(client, networkId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "RESOURCE_NOT_FOUND", "COMPLETED", nil
			}
			return respBody, "ERROR", err
		}

		status := utils.PathSearch("status.phase", respBody, "").(string)
		if utils.StrSliceContains(targets, status) {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func waitForNetworkCreateCompleted(ctx context.Context, client *golangsdk.ServiceClient, networkId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshNetworkStatusFunc(client, networkId, []string{"Active"}),
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	respBody, err := createNetwork(client, d)
	if err != nil {
		return diag.Errorf("error creating network: %s", err)
	}

	resourceId := utils.PathSearch("metadata.name", respBody, "").(string)
	if resourceId == "" {
		return diag.Errorf("unable to find the network ID from the API response")
	}
	d.SetId(resourceId)

	if err = waitForNetworkCreateCompleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for the network (%s) creation to complete: %s", d.Id(), err)
	}

	return resourceNetworkRead(ctx, d, meta)
}

func flattenNetworkPeerConnections(peeringConnections []interface{}) []map[string]interface{} {
	if len(peeringConnections) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(peeringConnections))
	for _, peeringConnection := range peeringConnections {
		result = append(result, map[string]interface{}{
			"vpc_id":    utils.PathSearch("peerVpcId", peeringConnection, ""),
			"subnet_id": utils.PathSearch("peerSubnetId", peeringConnection, ""),
		})
	}
	return result
}

func resourceNetworkRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		networkId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	respBody, err := GetNetworkById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving network (%s)", networkId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch(`metadata.labels."os.modelarts/name"`, respBody, nil)),
		d.Set("workspace_id", utils.PathSearch(`metadata.labels."os.modelarts/workspace.id"`, respBody, nil)),
		d.Set("cidr", utils.PathSearch("spec.cidr", respBody, nil)),
		d.Set("peer_connections", flattenNetworkPeerConnections(utils.PathSearch("spec.connection.peerConnectionList",
			respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("status", utils.PathSearch("status.phase", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateNetworkBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"spec": map[string]interface{}{
			"connection": map[string]interface{}{
				"peerConnectionList": buildNetworkPeerConnections(d.Get("peer_connections").([]interface{})),
			},
		},
	}
}

func parseNetworkAssociatedPeerConnections(connections []interface{}) []string {
	if len(connections) < 1 {
		return nil
	}

	result := make([]string, 0, len(connections))
	for _, connection := range connections {
		result = append(result, fmt.Sprintf("%s:%s", utils.PathSearch("peerVpcId|vpc_id", connection, ""), utils.PathSearch("peerSubnetId|subnet_id", connection, "")))
	}

	return result
}

func checkAllPeerConnectionsExist(remoteConnections, localConnections []string) bool {
	if len(localConnections) < 1 || len(remoteConnections) < 1 {
		return true
	}

	return utils.IsSliceContainsAnyAnotherSliceElement(remoteConnections, localConnections, true, true)
}

func refreshNetworkConnectionsFunc(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			networkId        = d.Id()
			localConnections = parseNetworkAssociatedPeerConnections(d.Get("peer_connections").([]interface{}))
		)

		respBody, err := GetNetworkById(client, networkId)
		if err != nil {
			return nil, "ERROR", err
		}

		remoteConnections := parseNetworkAssociatedPeerConnections(utils.PathSearch("spec.connection.peerConnectionList",
			respBody, make([]interface{}, 0)).([]interface{}))
		if !checkAllPeerConnectionsExist(remoteConnections, localConnections) {
			log.Printf("some peer connections are not found, remote: %v, local: %v", remoteConnections, localConnections)
			return "Missing Some Peer Connections", "PENDING", nil
		}

		return respBody, "COMPLETED", nil
	}
}

func waitForNetworkConnectionsUpdateCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshNetworkConnectionsFunc(client, d),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func updateNetworkConnections(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpURL := "v1/{project_id}/networks/{id}"
	updatePath := client.Endpoint + httpURL
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/merge-patch+json",
		},
		JSONBody: buildUpdateNetworkBodyParams(d),
	}

	_, err := client.Request("PATCH", updatePath, &updateOpt)
	if err != nil {
		return err
	}

	return waitForNetworkConnectionsUpdateCompleted(ctx, client, d)
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	if d.HasChange("peer_connections") {
		if err = updateNetworkConnections(ctx, client, d); err != nil {
			return diag.Errorf("error updating network connections: %s", err)
		}
	}
	return resourceNetworkRead(ctx, d, meta)
}

func waitForNetworkDeleteCompleted(ctx context.Context, client *golangsdk.ServiceClient, networkId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshNetworkStatusFunc(client, networkId, nil),
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func deleteNetwork(ctx context.Context, client *golangsdk.ServiceClient, networkID string) error {
	httpURL := "v1/{project_id}/networks/{id}"
	deletePath := client.Endpoint + httpURL
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", networkID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		networkId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = deleteNetwork(ctx, client, networkId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting network (%s)", networkId))
	}

	if err = waitForNetworkDeleteCompleted(ctx, client, networkId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for the network (%s) deletion to complete: %s", networkId, err)
	}
	return nil
}
