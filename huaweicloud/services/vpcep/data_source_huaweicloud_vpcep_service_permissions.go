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

// @API VPCEP GET /v1/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/permissions
func DataSourceVPCEPServicePermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceVpcepServicePermissionsRead,

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
			"permission": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Elem:     permissionsSchema(),
				Computed: true,
			},
		},
	}
}

func permissionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"permission_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func datasourceVpcepServicePermissionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcepClient, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	serviceId := d.Get("service_id").(string)
	listPermOpts := services.ListPermOpts{
		Permission: d.Get("permission").(string),
	}

	allPermissions, err := services.ListAllPermissions(vpcepClient, serviceId, listPermOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve VPC endpoint service permissions: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("permissions", flattenListVPCEPServicePermissions(allPermissions)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListVPCEPServicePermissions(allPermissions []services.Permission) []map[string]interface{} {
	if allPermissions == nil {
		return nil
	}
	permissions := make([]map[string]interface{}, len(allPermissions))
	for i, v := range allPermissions {
		permissions[i] = map[string]interface{}{
			"permission_id":   v.ID,
			"permission":      v.Permission,
			"permission_type": v.PermissionType,
			"created_at":      v.Created,
		}
	}
	return permissions
}
