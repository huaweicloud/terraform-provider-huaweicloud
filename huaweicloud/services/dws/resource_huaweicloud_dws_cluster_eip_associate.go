package dws

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var clusterEipAssociateNonUpdatableParams = []string{
	"cluster_id",
	"eip_id",
}

// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/eips/{eip_id}
// @API DWS GET /v1/{project_id}/clusters/{cluster_id}/endpoints
// @API DWS DELETE /v2/{project_id}/clusters/{cluster_id}/eips/{eip_id}
func ResourceClusterEipAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterEipAssociateCreate,
		ReadContext:   resourceClusterEipAssociateRead,
		UpdateContext: resourceClusterEipAssociateUpdate,
		DeleteContext: resourceClusterEipAssociateDelete,

		CustomizeDiff: config.FlexibleForceNew(clusterEipAssociateNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the cluster (to which the EIP associated) is located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster to which the EIP is associated.`,
			},
			"eip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the EIP to be associated with the DWS cluster.`,
			},

			// Attributes.
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public IP address of the cluster endpoint.`,
			},
			"public_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The public port of the cluster endpoint.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
}

func associateClusterEip(client *golangsdk.ServiceClient, clusterId, eipId string) error {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/eips/{eip_id}"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", clusterId)
	createPath = strings.ReplaceAll(createPath, "{eip_id}", eipId)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("POST", createPath, &createOpts)
	return err
}

func disassociateClusterEip(client *golangsdk.ServiceClient, clusterId, eipId string) error {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/eips/{eip_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", clusterId)
	deletePath = strings.ReplaceAll(deletePath, "{eip_id}", eipId)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpts)
	return err
}

// GetClusterAssociatedEipById is a method used to query the public endpoint and verify whether the EIP is associated.
func GetClusterAssociatedEipById(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	endpoints, err := getClusterEndpointsById(client, clusterId)
	if err != nil {
		return nil, err
	}

	publicEndpoints := utils.PathSearch("public_endpoints", endpoints, nil)
	currentEipId := utils.PathSearch("ip_id", publicEndpoints, "").(string)
	if currentEipId != "" {
		return publicEndpoints, nil
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v1/{project_id}/clusters/{cluster_id}/endpoints",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the EIP is not associated with the DWS cluster (%s)", clusterId)),
		},
	}
}

func resourceClusterEipAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
		eipId     = d.Get("eip_id").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	err = associateClusterEip(client, clusterId, eipId)
	if err != nil {
		return diag.Errorf("error associating EIP (%s) to the DWS cluster (%s): %s", eipId, clusterId, err)
	}

	d.SetId(clusterId)
	return resourceClusterEipAssociateRead(ctx, d, meta)
}

func resourceClusterEipAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Id()
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	respBody, err := GetClusterAssociatedEipById(client, clusterId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error retrieving EIP association for the DWS cluster (%s)", clusterId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("cluster_id", clusterId),
		d.Set("eip_id", utils.PathSearch("ip_id", respBody, nil)),
		d.Set("public_ip", utils.PathSearch("ip", respBody, nil)),
		d.Set("public_port", utils.PathSearch("port", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceClusterEipAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because all parameters are NonUpdatable.
	return resourceClusterEipAssociateRead(ctx, d, meta)
}

func resourceClusterEipAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Id()
		eipId     = d.Get("eip_id").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	err = disassociateClusterEip(client, clusterId, eipId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error disassociating EIP (%s) from the DWS cluster (%s)", eipId, clusterId))
	}

	return nil
}
