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

var clusterPublicDomainAssociateNonUpdatableParams = []string{
	"cluster_id",
}

// @API DWS POST /v1.0/{project_id}/clusters/{cluster_id}/dns
// @API DWS GET /v1/{project_id}/clusters/{cluster_id}/endpoints
// @API DWS PUT /v1.0/{project_id}/clusters/{cluster_id}/dns
// @API DWS DELETE /v1.0/{project_id}/clusters/{cluster_id}/dns?type=public
func ResourceClusterPublicDomainAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterPublicDomainAssociateCreate,
		ReadContext:   resourceClusterPublicDomainAssociateRead,
		UpdateContext: resourceClusterPublicDomainAssociateUpdate,
		DeleteContext: resourceClusterPublicDomainAssociateDelete,

		CustomizeDiff: config.FlexibleForceNew(clusterPublicDomainAssociateNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: clusterPublicDomainAssociateImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the cluster (to which the public domain belongs) is located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the DWS cluster.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The public domain name.`,
			},

			// Optional parameters.
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The cache period of the SOA record set, in seconds.`,
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

func buildClusterPublicDomainAssociateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name": d.Get("domain_name"),
		// `type` is fixed to `public` for this resource.
		"type": "public",
		"ttl":  utils.ValueIgnoreEmpty(d.Get("ttl").(int)),
	}
}

func createClusterPublicDomainAssociate(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1.0/{project_id}/clusters/{cluster_id}/dns"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", d.Get("cluster_id").(string))

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildClusterPublicDomainAssociateBodyParams(d)),
	}

	_, err := client.Request("POST", createPath, &createOpts)
	return err
}

func resourceClusterPublicDomainAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	err = createClusterPublicDomainAssociate(client, d)
	if err != nil {
		return diag.Errorf("error associating public domain to the cluster (%s): %s", clusterId, err)
	}

	d.SetId(clusterId)

	return resourceClusterPublicDomainAssociateRead(ctx, d, meta)
}

func getClusterEndpointsById(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/clusters/{cluster_id}/endpoints"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func GetClusterPublicEndpointById(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	endpoints, err := getClusterEndpointsById(client, clusterId)
	if err != nil {
		return nil, err
	}

	publicEndpoints := utils.PathSearch("public_endpoints", endpoints, nil)
	if utils.PathSearch("domain_name", publicEndpoints, "").(string) != "" {
		return publicEndpoints, nil
	}
	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v1/{project_id}/clusters/{cluster_id}/endpoints",
			RequestId: "NONE",
			Body:      fmt.Appendf(nil, "the public domain is not associated with the cluster (%s)", clusterId),
		},
	}
}

func resourceClusterPublicDomainAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Id()
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	respBody, err := GetClusterPublicEndpointById(client, clusterId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error retrieving public domain association for the cluster (%s)", clusterId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("cluster_id", clusterId),
		d.Set("domain_name", utils.PathSearch("domain_name", respBody, nil)),
		d.Set("ttl", utils.PathSearch("domain_name_ttl", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func updateClusterPublicDomainAssociate(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1.0/{project_id}/clusters/{cluster_id}/dns"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", d.Get("cluster_id").(string))

	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildClusterPublicDomainAssociateBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateOpts)
	return err
}

func resourceClusterPublicDomainAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Id()
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if d.HasChanges("domain_name", "ttl") {
		err = updateClusterPublicDomainAssociate(client, d)
		if err != nil {
			return diag.Errorf("error updating public domain association for the cluster (%s): %s", clusterId, err)
		}
	}

	return resourceClusterPublicDomainAssociateRead(ctx, d, meta)
}

func deleteClusterPublicDomainAssociate(client *golangsdk.ServiceClient, clusterId string) error {
	httpUrl := "v1.0/{project_id}/clusters/{cluster_id}/dns?type=public"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", clusterId)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpts)
	return err
}

func resourceClusterPublicDomainAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Id()
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	err = deleteClusterPublicDomainAssociate(client, clusterId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error disassociating public domain from the cluster (%s)", clusterId))
	}

	return nil
}

func clusterPublicDomainAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	switch len(parts) {
	case 1:
		return []*schema.ResourceData{d}, nil
	case 2:
		d.SetId(parts[0])
		return []*schema.ResourceData{d}, d.Set("domain_name", parts[1])
	default:
	}
	return nil, fmt.Errorf("invalid format specified for import ID, want '<cluster_id>' or '<cluster_id>/<domain_name>', but got '%s'", d.Id())
}
