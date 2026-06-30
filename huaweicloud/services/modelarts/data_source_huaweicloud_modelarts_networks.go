package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v1/{project_id}/networks
func DataSourceNetworks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the networks are located.`,
			},

			// Optional parameters.
			"label_selector": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The label selector to filter networks.`,
			},

			// Attributes.
			"networks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of networks that matched filter parameters.`,
				Elem:        networksSchema(),
			},
		},
	}
}

func networksSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The metadata of the network.`,
				Elem:        networksMetadataSchema(),
			},
			"spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The spec of the network.`,
				Elem:        networksSpecSchema(),
			},
			"status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The status of the network.`,
				Elem:        networksStatusSchema(),
			},
		},
	}
}

func networksMetadataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the network.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the network.`,
			},
			"labels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The labels of the network.`,
				Elem:        networksMetadataLabelsSchema(),
			},
			"annotations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The annotations of the network.`,
				Elem:        networksMetadataAnnotationsSchema(),
			},
		},
	}
}

func networksMetadataLabelsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"os_modelarts_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The display name of the resource pool.`,
			},
			"os_modelarts_workspace_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The workspace ID of the resource pool.`,
			},
		},
	}
}

func networksMetadataAnnotationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"os_modelarts_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the network.`,
			},
		},
	}
}

func networksSpecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The CIDR block of the network.`,
			},
			"connection": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The connection information of the network.`,
				Elem:        networksSpecConnectionSchema(),
			},
		},
	}
}

func networksSpecConnectionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"peer_connection_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The peer connection list of the network.`,
				Elem:        networksSpecConnectionPeerConnectionItemSchema(),
			},
		},
	}
}

func networksSpecConnectionPeerConnectionItemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"peer_vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the peer VPC.`,
			},
			"peer_subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the peer subnet.`,
			},
			"default_gateway": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to create the default gateway.`,
			},
		},
	}
}

func networksStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The phase of the network.`,
			},
			"connection_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The connection status of the network.`,
				Elem:        networksStatusConnectionStatusSchema(),
			},
		},
	}
}

func networksStatusConnectionStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"peer_connection_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The peer connection status list of the network.`,
				Elem:        networksStatusPeerConnectionStatusSchema(),
			},
			"sfs_turbo_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The SFS Turbo connection status list of the network.`,
				Elem:        networksStatusSfsTurboStatusSchema(),
			},
		},
	}
}

func networksStatusPeerConnectionStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"peer_vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the peer VPC.`,
			},
			"peer_subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the peer subnet.`,
			},
			"default_gateway": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the default gateway is created.`,
			},
			"phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The connection phase of the peer connection.`,
			},
		},
	}
}

func networksStatusSfsTurboStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the SFS Turbo instance.`,
			},
			"sfs_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the SFS Turbo instance.`,
			},
			"connection_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The connection type of the SFS Turbo.`,
			},
			"ip_addr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP address of the SFS Turbo.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The connection status of the SFS Turbo.`,
			},
		},
	}
}

func buildNetworksQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("label_selector"); ok {
		res = fmt.Sprintf("%s&labelSelector=%v", res, v)
	}

	return res
}

func listNetworks(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl          = "v1/{project_id}/networks?limit=500"
		result           = make([]interface{}, 0)
		metadataContinue string
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildNetworksQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// The value is taken from the value in `metadata.continue` in the user's previous pagination query response,
	// in UUID format, defaults to null.
	for {
		reqPath := listPath
		if metadataContinue != "" {
			reqPath = fmt.Sprintf("%s&continue=%s", listPath, metadataContinue)
		}

		requestResp, err := client.Request("GET", reqPath, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		networks := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, networks...)
		if len(networks) < 1 {
			break
		}

		metadataContinue = utils.PathSearch("metadata.continue", respBody, "").(string)
		if metadataContinue == "" {
			break
		}
	}

	return result, nil
}

func flattenNetworkAnnotations(annotations map[string]interface{}) []map[string]interface{} {
	if len(annotations) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"os_modelarts_description": utils.PathSearch(`"os.modelarts/description"`, annotations, nil),
		},
	}
}

func flattenNetworkLabels(labels map[string]interface{}) []map[string]interface{} {
	if len(labels) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"os_modelarts_name":         utils.PathSearch(`"os.modelarts/name"`, labels, nil),
			"os_modelarts_workspace_id": utils.PathSearch(`"os.modelarts/workspace.id"`, labels, nil),
		},
	}
}

func flattenNetworkMetadata(metadata map[string]interface{}) []map[string]interface{} {
	if len(metadata) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":       utils.PathSearch("name", metadata, nil),
			"created_at": utils.PathSearch("creationTimestamp", metadata, nil),
			"labels": flattenNetworkLabels(utils.PathSearch("labels", metadata,
				make(map[string]interface{})).(map[string]interface{})),
			"annotations": flattenNetworkAnnotations(utils.PathSearch("annotations", metadata,
				make(map[string]interface{})).(map[string]interface{})),
		},
	}
}

func flattenNetworkPeerConnectionList(peerConnectionList []interface{}) []map[string]interface{} {
	if len(peerConnectionList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(peerConnectionList))
	for _, conn := range peerConnectionList {
		result = append(result, map[string]interface{}{
			"peer_vpc_id":     utils.PathSearch("peerVpcId", conn, nil),
			"peer_subnet_id":  utils.PathSearch("peerSubnetId", conn, nil),
			"default_gateway": utils.PathSearch("defaultGateWay", conn, nil),
		})
	}

	return result
}

func flattenNetworkConnection(connection map[string]interface{}) []map[string]interface{} {
	if len(connection) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"peer_connection_list": flattenNetworkPeerConnectionList(
				utils.PathSearch("peerConnectionList", connection, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenNetworkSpec(spec map[string]interface{}) []map[string]interface{} {
	if len(spec) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"cidr": utils.PathSearch("cidr", spec, nil),
			"connection": flattenNetworkConnection(utils.PathSearch("connection", spec,
				make(map[string]interface{})).(map[string]interface{})),
		},
	}
}

func flattenNetworkSfsTurboStatusList(sfsTurboStatuses []interface{}) []map[string]interface{} {
	if len(sfsTurboStatuses) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(sfsTurboStatuses))
	for _, status := range sfsTurboStatuses {
		result = append(result, map[string]interface{}{
			"name":            utils.PathSearch("name", status, nil),
			"sfs_id":          utils.PathSearch("sfsId", status, nil),
			"connection_type": utils.PathSearch("connectionType", status, nil),
			"ip_addr":         utils.PathSearch("ipAddr", status, nil),
			"status":          utils.PathSearch("status", status, nil),
		})
	}

	return result
}

func flattenNetworkPeerConnectionStatusList(peerConnectionStatuses []interface{}) []map[string]interface{} {
	if len(peerConnectionStatuses) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(peerConnectionStatuses))
	for _, status := range peerConnectionStatuses {
		result = append(result, map[string]interface{}{
			"peer_vpc_id":     utils.PathSearch("peerVpcId", status, nil),
			"peer_subnet_id":  utils.PathSearch("peerSubnetId", status, nil),
			"default_gateway": utils.PathSearch("defaultGateWay", status, nil),
			"phase":           utils.PathSearch("phase", status, nil),
		})
	}

	return result
}

func flattenNetworkConnectionStatus(connectionStatus map[string]interface{}) []map[string]interface{} {
	if len(connectionStatus) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"peer_connection_status": flattenNetworkPeerConnectionStatusList(
				utils.PathSearch("peerConnectionStatus", connectionStatus, make([]interface{}, 0)).([]interface{})),
			"sfs_turbo_status": flattenNetworkSfsTurboStatusList(
				utils.PathSearch("sfsTurboStatus", connectionStatus, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenNetworkStatus(networkStatus map[string]interface{}) []map[string]interface{} {
	if len(networkStatus) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"phase": utils.PathSearch("phase", networkStatus, nil),
			"connection_status": flattenNetworkConnectionStatus(utils.PathSearch("connectionStatus", networkStatus,
				make(map[string]interface{})).(map[string]interface{})),
		},
	}
}

func flattenNetworks(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		metadata := utils.PathSearch("metadata", item, make(map[string]interface{})).(map[string]interface{})
		spec := utils.PathSearch("spec", item, make(map[string]interface{})).(map[string]interface{})
		networkStatus := utils.PathSearch("status", item, make(map[string]interface{})).(map[string]interface{})

		result = append(result, map[string]interface{}{
			"metadata": flattenNetworkMetadata(metadata),
			"spec":     flattenNetworkSpec(spec),
			"status":   flattenNetworkStatus(networkStatus),
		})
	}

	return result
}

func dataSourceNetworksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	networks, err := listNetworks(client, d)
	if err != nil {
		return diag.Errorf("error querying networks: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("networks", flattenNetworks(networks)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
