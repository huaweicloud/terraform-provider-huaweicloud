package vpcep

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/vpcep/v1/services"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPCEP GET /v1/{project_id}/vpc-endpoint-services
func DataSourceVPCEPServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceVpcepServicesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint_services": {
				Type:     schema.TypeList,
				Elem:     servicesSchema(),
				Computed: true,
			},
		},
	}
}

func servicesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"approval_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
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
			"tags": common.TagsComputedSchema(),
			"connection_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tcp_proxy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_policy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
	return &sc
}

func datasourceVpcepServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcepClient, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	listOpts := services.ListOpts{
		ServiceName:       d.Get("service_name").(string),
		ID:                d.Get("service_id").(string),
		Status:            d.Get("status").(string),
		ServerType:        d.Get("server_type").(string),
		PublicBorderGroup: d.Get("public_border_group").(string),
	}

	allServices, err := services.ListAllServices(vpcepClient, listOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve VPC endpoint services: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("endpoint_services", flattenListVpcepServices(allServices)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListVpcepServices(allServices []services.Service) []map[string]interface{} {
	if allServices == nil {
		return nil
	}
	endpointServices := make([]map[string]interface{}, len(allServices))
	for i, v := range allServices {
		endpointServices[i] = map[string]interface{}{
			"id":                  v.ID,
			"service_name":        v.ServiceName,
			"service_type":        v.ServiceType,
			"server_type":         v.ServerType,
			"vpc_id":              v.VpcID,
			"approval_enabled":    v.Approval,
			"status":              v.Status,
			"created_at":          v.Created,
			"updated_at":          v.Updated,
			"tags":                utils.TagsToMap(v.Tags),
			"tcp_proxy":           v.TCPProxy,
			"description":         v.Description,
			"enable_policy":       v.EnablePolicy,
			"public_border_group": v.PublicBorderGroup,
			"connection_count":    v.ConnectionCount,
		}
	}
	return endpointServices
}
