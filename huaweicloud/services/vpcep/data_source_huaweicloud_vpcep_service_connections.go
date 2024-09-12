package vpcep

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/vpcep/v1/services"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API VPCEP GET /v1/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/connections
func DataSourceVPCEPServiceConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceVpcepServiceConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"marker_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connections": {
				Type:     schema.TypeList,
				Elem:     connectionsSchema(),
				Computed: true,
			},
		},
	}
}

func connectionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"marker_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func datasourceVpcepServiceConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcepClient, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	serviceId := d.Get("service_id").(string)
	listConnOpts := services.ListConnOpts{
		EndpointID: d.Get("endpoint_id").(string),
		MarkerID:   d.Get("marker_id").(string),
		Status:     d.Get("status").(string),
	}

	allConnections, err := services.ListAllConnections(vpcepClient, serviceId, listConnOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve VPC endpoint service connections: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("connections", flattenListVPCEPServiceConnections(allConnections)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListVPCEPServiceConnections(allConnections []services.Connection) []map[string]interface{} {
	if allConnections == nil {
		return nil
	}
	connections := make([]map[string]interface{}, len(allConnections))
	for i, v := range allConnections {
		connections[i] = map[string]interface{}{
			"endpoint_id": v.EndpointID,
			"marker_id":   v.MarkerID,
			"domain_id":   v.DomainID,
			"status":      v.Status,
			"created_at":  v.Created,
			"updated_at":  v.Updated,
			"description": v.Description,
		}
	}
	return connections
}
