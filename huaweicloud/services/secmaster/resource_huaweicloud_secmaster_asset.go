package secmaster

import (
	"context"
	"fmt"
	"log"
	"strconv"
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

// The API documentation has quality issues, which leads to inaccurate descriptions of the usage of some fields.
// For example, `created` and `created_at` are both considered required fields.
// Keep all fields in this resource consistent with the API documentation. Although some fields are unreasonable, it does
// not affect usage.

// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/sa/resources/{id}
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/sa/resources/{id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/sa/resources
func ResourceAsset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetCreate,
		ReadContext:   resourceAssetRead,
		UpdateContext: resourceAssetUpdate,
		DeleteContext: resourceAssetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAssetImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{"workspace_id", "asset_id"}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the workspace to which the asset belongs.`,
			},
			"asset_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the asset.`,
			},
			// The API used by this resource is quite tolerant, and many fields that are restricted as mandatory in
			// documents can actually be left blank. To better accommodate usage scenarios, all required fields in
			// data_object are modified to required at the Description level.
			"data_object": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Specifies the details of the asset.`,
				Elem:        buildAssetDataObjectSchema(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildAssetDataObjectSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the ID of the asset.", utils.SchemaDescInput{Required: true}),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the name of the asset.", utils.SchemaDescInput{Required: true}),
			},
			"provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the provider of the asset.",
					utils.SchemaDescInput{Required: true}),
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the type of the asset.", utils.SchemaDescInput{Required: true}),
			},
			"environment": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     buildAssetEnvironmentSchema(),
				Description: utils.SchemaDesc("Specifies the environment of the asset.",
					utils.SchemaDescInput{Required: true}),
			},
			"properties": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     buildAssetPropertiesSchema(),
				Description: utils.SchemaDesc("Specifies the properties of the asset.",
					utils.SchemaDescInput{Required: true}),
			},
			"checksum": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the checksum of the asset.`,
			},
			"created": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the creation time of the asset.`,
			},
			"provisioning_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the provisioning state of the asset.`,
			},
			"department": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetDepartmentSchema(),
				Description: `Specifies the department of the asset.`,
			},
			"governance_user": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetGovernanceUserSchema(),
				Description: `Specifies the governance user of the asset.`,
			},
			"level": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the level of the asset.`,
			},
		},
	}
	return sc
}

func buildAssetPropertiesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hwc_ecs": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcEcsSchema(),
				Description: `Specifies the details of the ECS.`,
			},
			"hwc_eip": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcEipSchema(),
				Description: `Specifies the details of the EIP.`,
			},
			"hwc_vpc": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcVpcSchema(),
				Description: `Specifies the details of the VPC.`,
			},
			"hwc_subnet": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcSubnetSchema(),
				Description: `Specifies the details of the subnet.`,
			},
			"hwc_rds": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcRdsSchema(),
				Description: `Specifies the details of the RDS.`,
			},
			"hwc_domain": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcDomainSchema(),
				Description: `Specifies the details of the domain.`,
			},
			"website": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesWebsiteSchema(),
				Description: `Specifies the details of the website.`,
			},
			"oca_ip": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesOcaIpSchema(),
				Description: `Specifies the details of the cloud asset IP.`,
			},
		},
	}
}

func buildAssetPropertiesOcaIpSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the asset value.", utils.SchemaDescInput{Required: true}),
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the asset version.", utils.SchemaDescInput{Required: true}),
			},
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     buildAssetPropertiesOcaIpNetworkSchema(),
				Description: utils.SchemaDesc("Specifies the network information.",
					utils.SchemaDescInput{Required: true}),
			},
			"server_room": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the server room.", utils.SchemaDescInput{Required: true}),
			},
			"server_rack": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the server rack.", utils.SchemaDescInput{Required: true}),
			},
			"data_center": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesOcaIpDataCenterSchema(),
				Description: utils.SchemaDesc("Specifies the data center.", utils.SchemaDescInput{Required: true}),
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the asset remark.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the asset name, default value is asset value.`,
			},
			"relative_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the relative value, such as ipv6 if the asset is ipv4.`,
			},
			"mac_addr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the mac address.`,
			},
			"important": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the importance level, 0: not important, 1: important`,
			},
			"extend_propertites": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesOcaIpExtendPropertitesSchema(),
				Description: `Specifies the other third-party attributes.`,
			},
		},
	}
}

func buildAssetPropertiesOcaIpExtendPropertitesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"device": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesOcaIpExtendPropertitesDeviceSchema(),
				Description: `Specifies the device information.`,
			},
			"system": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesOcaIpExtendPropertitesSystemSchema(),
				Description: `Specifies the system information.`,
			},
			"services": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        buildAssetPropertiesOcaIpExtendPropertitesServicesSchema(),
				Description: `Specifies the application information.`,
			},
		},
	}
}

func buildAssetPropertiesOcaIpExtendPropertitesServicesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the port of the application.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the protocol of the application.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the application.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the version of the application.`,
			},
			"vendor": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesOcaIpExtendPropertitesServicesVendorSchema(),
				Description: `Specifies the application supplier.`,
			},
		},
	}
}

func buildAssetPropertiesOcaIpExtendPropertitesServicesVendorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the supplier.`,
			},
			"is_xc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether the supplier is domestic or not.`,
			},
		},
	}
}

func buildAssetPropertiesOcaIpExtendPropertitesSystemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"family": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the type of the system.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the system.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the version of the system.`,
			},
			"vendor": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesOcaIpExtendPropertitesSystemVendorSchema(),
				Description: `Specifies the system supplier.`,
			},
		},
	}
}

func buildAssetPropertiesOcaIpExtendPropertitesSystemVendorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the supplier.`,
			},
			"is_xc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether the supplier is domestic or not.`,
			},
		},
	}
}

func buildAssetPropertiesOcaIpExtendPropertitesDeviceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the type of the device.`,
			},
			"model": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the model of the device.`,
			},
			"vendor": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesOcaIpExtendPropertitesDeviceVendorSchema(),
				Description: `Specifies the device supplier.`,
			},
		},
	}
}

func buildAssetPropertiesOcaIpExtendPropertitesDeviceVendorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the supplier.`,
			},
			"is_xc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether the supplier is domestic or not.`,
			},
		},
	}
}

func buildAssetPropertiesOcaIpDataCenterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"city_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the city code.", utils.SchemaDescInput{Required: true}),
			},
			"country_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the country code.", utils.SchemaDescInput{Required: true}),
			},
			"latitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the latitude.`,
			},
			"longitude": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the longitude.`,
			},
		},
	}
}

func buildAssetPropertiesOcaIpNetworkSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies whether the IP is public or private.",
					utils.SchemaDescInput{Required: true}),
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the network partition.`,
			},
			"plane": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the network plane (offline has its own code).`,
			},
			"vxlan_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the virtual network ID.`,
			},
		},
	}
}

func buildAssetPropertiesWebsiteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the website URL.", utils.SchemaDescInput{Required: true}),
			},
			"main_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the main domain name.", utils.SchemaDescInput{Required: true}),
			},
			"protected_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the WAF status.", utils.SchemaDescInput{Required: true}),
			},
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies whether the website is public or private.",
					utils.SchemaDescInput{Required: true}),
			},
			"name_server": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: utils.SchemaDesc("Specifies the website server list.", utils.SchemaDescInput{Required: true}),
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the website remark.`,
			},
			"extend_propertites": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesWebsiteExtendPropertitesSchema(),
				Description: `Specifies the other properties.`,
			},
		},
	}
}

func buildAssetPropertiesWebsiteExtendPropertitesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"mac_addr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the MAC address.`,
			},
		},
	}
}

func buildAssetPropertiesHwcDomainSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the domain name.", utils.SchemaDescInput{Required: true}),
			},
			"expire_date": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the domain expiration date.",
					utils.SchemaDescInput{Required: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the domain service status.",
					utils.SchemaDescInput{Required: true}),
			},
			"audit_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the domain real-name authentication status.",
					utils.SchemaDescInput{Required: true}),
			},
			"audit_unpass_reason": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the reason for domain real-name authentication failure.",
					utils.SchemaDescInput{Required: true}),
			},
			"reg_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the registration type.", utils.SchemaDescInput{Required: true}),
			},
			"privacy_protection": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies whether privacy protection is enabled.",
					utils.SchemaDescInput{Required: true}),
			},
			"name_server": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: utils.SchemaDesc("Specifies the domain name server list.",
					utils.SchemaDescInput{Required: true}),
			},
			"credential_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the type of credential.", utils.SchemaDescInput{Required: true}),
			},
			"credential_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the credential ID.", utils.SchemaDescInput{Required: true}),
			},
			"registrar": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the registrar of the domain.",
					utils.SchemaDescInput{Required: true}),
			},
			"contact": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     buildAssetPropertiesHwcDomainContactSchema(),
				Description: utils.SchemaDesc("Specifies the contact information.",
					utils.SchemaDescInput{Required: true}),
			},
			"transfer_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the domain transfer status.`,
			},
		},
	}
}

func buildAssetPropertiesHwcDomainContactSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the email address.", utils.SchemaDescInput{Required: true}),
			},
			"register": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the domain owner.", utils.SchemaDescInput{Required: true}),
			},
			"contact_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the contact name.", utils.SchemaDescInput{Required: true}),
			},
			"phone_num": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the phone number.", utils.SchemaDescInput{Required: true}),
			},
			"province": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the province.", utils.SchemaDescInput{Required: true}),
			},
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the city.", utils.SchemaDescInput{Required: true}),
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the contact address.", utils.SchemaDescInput{Required: true}),
			},
			"zip_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the zip code.", utils.SchemaDescInput{Required: true}),
			},
		},
	}
}

func buildAssetPropertiesHwcRdsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the ID of the RDS.", utils.SchemaDescInput{Required: true}),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the name of the RDS.", utils.SchemaDescInput{Required: true}),
			},
			"protected_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the DBSS opening status of the RDS.",
					utils.SchemaDescInput{Required: true}),
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the status of the RDS.", utils.SchemaDescInput{Required: true}),
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the database port of the RDS.",
					utils.SchemaDescInput{Required: true}),
			},
			"enable_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the SSL flag of the instance.",
					utils.SchemaDescInput{Required: true}),
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the type of the RDS.", utils.SchemaDescInput{Required: true}),
			},
			"ha": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     buildAssetPropertiesHwcRdsHaSchema(),
				Description: utils.SchemaDesc("HA information, returned when the HA instance is obtained.",
					utils.SchemaDescInput{Required: true}),
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the region where the RDS is located.",
					utils.SchemaDescInput{Required: true}),
			},
			"datastore": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcRdsDatastoreSchema(),
				Description: utils.SchemaDesc("Specifies the database information.", utils.SchemaDescInput{Required: true}),
			},
			"created": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the creation time.", utils.SchemaDescInput{Required: true}),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the VPC ID.", utils.SchemaDescInput{Required: true}),
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the subnet ID of the RDS.", utils.SchemaDescInput{Required: true}),
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the security group ID.", utils.SchemaDescInput{Required: true}),
			},
			"flavor_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the flavor of the RDS.", utils.SchemaDescInput{Required: true}),
			},
			"cpu": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the CPU size of the RDS.", utils.SchemaDescInput{Required: true}),
			},
			"mem": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the memory size of the RDS.",
					utils.SchemaDescInput{Required: true}),
			},
			"volume": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcRdsVolumeSchema(),
				Description: utils.SchemaDesc("Specifies the volume information.", utils.SchemaDescInput{Required: true}),
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the project ID.", utils.SchemaDescInput{Required: true}),
			},
			"switch_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the database switch strategy.respectively.",
					utils.SchemaDescInput{Required: true}),
			},
			"read_only_by_user": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					"Specifies the user set read-only status of the RDS. Only supports RDS for MySQL engine.",
					utils.SchemaDescInput{Required: true}),
			},
			"backup_strategy": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcRdsBackupStrategySchema(),
				Description: utils.SchemaDesc("Specifies the backup policy.", utils.SchemaDescInput{Required: true}),
			},
			"maintenance_window": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the maintenance window of the RDS.",
					utils.SchemaDescInput{Required: true}),
			},
			"nodes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     buildAssetPropertiesHwcRdsNodesSchema(),
				Description: utils.SchemaDesc("Main and standby instance information.",
					utils.SchemaDescInput{Required: true}),
			},
			"related_instance": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     buildAssetPropertiesHwcRdsRelatedInstanceSchema(),
				Description: utils.SchemaDesc("Specifies the list of associated database instances.",
					utils.SchemaDescInput{Required: true}),
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the time zone of the RDS.",
					utils.SchemaDescInput{Required: true}),
			},
			"storage_used_space": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the disk space usage, unit is GB.",
					utils.SchemaDescInput{Required: true}),
			},
			"associated_with_ddm": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies whether the instance is associated with DDM.",
					utils.SchemaDescInput{Required: true}),
			},
			"max_iops": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the maximum IOPS of the disk.",
					utils.SchemaDescInput{Required: true}),
			},
			"expiration_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					"Specifies the expiration time of the RDS. Only supports RDS for MySQL engine.",
					utils.SchemaDescInput{Required: true}),
			},
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the alias of the RDS.`,
			},
			"private_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the private IP addresses of the RDS.`,
			},
			"private_dns_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the private DNS names of the RDS.`,
			},
			"public_ips": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of public IP addresses of the instance.`,
			},
			"updated": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the update time of the RDS.`,
			},
			"db_user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the default username of the RDS.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        buildAssetPropertiesHwcRdsTagsSchema(),
				Description: `Specifies the tag list.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project ID of the RDS.`,
			},
			"disk_encryption_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the disk encryption key ID of the RDS.`,
			},
			"backup_used_space": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the backup space usage of the RDS. Only supports RDS for SQL Server engine.`,
			},
		},
	}
}

func buildAssetPropertiesHwcRdsRelatedInstanceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the ID of the associated instance.",
					utils.SchemaDescInput{Required: true}),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the type of the associated instance.",
					utils.SchemaDescInput{Required: true}),
			},
		},
	}
}

func buildAssetPropertiesHwcRdsNodesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the ID of the node.", utils.SchemaDescInput{Required: true}),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the name of the node.", utils.SchemaDescInput{Required: true}),
			},
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the role of the node.", utils.SchemaDescInput{Required: true}),
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the status of the node.", utils.SchemaDescInput{Required: true}),
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the availability zone of the node.",
					utils.SchemaDescInput{Required: true}),
			},
		},
	}
}

func buildAssetPropertiesHwcRdsBackupStrategySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the backup time period. Automatic backup will be triggered in this period.`,
			},
			"keep_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the number of days that the generated backup files can be saved.`,
			},
		},
	}
}

func buildAssetPropertiesHwcRdsTagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the key.", utils.SchemaDescInput{Required: true}),
			},
			"values": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: utils.SchemaDesc("Specifies the value list.", utils.SchemaDescInput{Required: true}),
			},
		},
	}
}

func buildAssetPropertiesHwcRdsVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the type of the volume.",
					utils.SchemaDescInput{Required: true}),
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the size of the volume.",
					utils.SchemaDescInput{Required: true}),
			},
		},
	}
}

func buildAssetPropertiesHwcRdsDatastoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the database engine.", utils.SchemaDescInput{Required: true}),
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the database version.", utils.SchemaDescInput{Required: true}),
			},
			"complete_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the complete version of the database.`,
			},
		},
	}
}

func buildAssetPropertiesHwcRdsHaSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"replication_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the replication mode of the RDS.",
					utils.SchemaDescInput{Required: true}),
			},
		},
	}
}

func buildAssetPropertiesHwcSubnetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the security group ID.", utils.SchemaDescInput{Required: true}),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the security group name.", utils.SchemaDescInput{Required: true}),
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the project ID to which the security group belongs.",
					utils.SchemaDescInput{Required: true}),
			},
			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the security group creation time.",
					utils.SchemaDescInput{Required: true}),
			},
			"updated_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the security group update time.",
					utils.SchemaDescInput{Required: true}),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the security group description.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project ID to which the security group belongs.`,
			},
			"security_group_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        buildAssetPropertiesHwcSubnetSecurityGroupRulesSchema(),
				Description: `Specifies the security group rules.`,
			},
		},
	}
}

func buildAssetPropertiesHwcSubnetSecurityGroupRulesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the unique identifier of the security group rule.",
					utils.SchemaDescInput{Required: true}),
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					"Specifies the security group ID to which the security group rule belongs.",
					utils.SchemaDescInput{Required: true}),
			},
			"direction": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the direction of the security group rule.",
					utils.SchemaDescInput{Required: true}),
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the protocol type.", utils.SchemaDescInput{Required: true}),
			},
			"ethertype": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the IP address protocol type.",
					utils.SchemaDescInput{Required: true}),
			},
			"multiport": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the port range.", utils.SchemaDescInput{Required: true}),
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the security group rule action.",
					utils.SchemaDescInput{Required: true}),
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the priority of the security group rule.",
					utils.SchemaDescInput{Required: true}),
			},
			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the security group rule creation time.",
					utils.SchemaDescInput{Required: true}),
			},
			"updated_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the security group rule update time.",
					utils.SchemaDescInput{Required: true}),
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the project ID to which the security group rule belongs.",
					utils.SchemaDescInput{Required: true}),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the security group rule description.`,
			},
			"remote_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the remote security group ID.`,
			},
			"remote_ip_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the remote IP address.`,
			},
			"remote_address_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the remote address group ID.`,
			},
		},
	}
}

func buildAssetPropertiesHwcVpcSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the ID of the VPC.", utils.SchemaDescInput{Required: true}),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the name of the VPC.", utils.SchemaDescInput{Required: true}),
			},
			"protected_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the protected status of the VPC.",
					utils.SchemaDescInput{Required: true}),
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the status of the VPC.", utils.SchemaDescInput{Required: true}),
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the project ID of the VPC.",
					utils.SchemaDescInput{Required: true}),
			},
			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the created time of the VPC.",
					utils.SchemaDescInput{Required: true}),
			},
			"updated_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the updated time of the VPC.",
					utils.SchemaDescInput{Required: true}),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description of the VPC.`,
			},
			"cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the cidr of the VPC.`,
			},
			"extend_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the extend cidrs of the VPC.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project ID of the VPC.`,
			},
			"cloud_resources": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        buildAssetPropertiesHwcVpcCloudResourcesSchema(),
				Description: `Specifies the cloud resources of the VPC.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        buildAssetPropertiesHwcVpcTagsSchema(),
				Description: `Specifies the tags of the VPC.`,
			},
		},
	}
}

func buildAssetPropertiesHwcVpcTagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the key of the tag.", utils.SchemaDescInput{Required: true}),
			},
			"values": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: utils.SchemaDesc("Specifies the values of the tag.", utils.SchemaDescInput{Required: true}),
			},
		},
	}
}

func buildAssetPropertiesHwcVpcCloudResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the type of the cloud resources.`,
			},
			"resource_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the asset count of the cloud resources.`,
			},
		},
	}
}

func buildAssetPropertiesHwcEipSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the ID of the EIP.", utils.SchemaDescInput{Required: true}),
			},
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the name of the EIP.", utils.SchemaDescInput{Required: true}),
			},
			"protected_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the protection status of the EIP.",
					utils.SchemaDescInput{Required: true}),
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the project ID of the EIP.",
					utils.SchemaDescInput{Required: true}),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the enterprise project ID of the EIP.",
					utils.SchemaDescInput{Required: true}),
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the IP version information.",
					utils.SchemaDescInput{Required: true}),
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the status of the EIP.", utils.SchemaDescInput{Required: true}),
			},
			"public_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the public IP address of the EIP.`,
			},
			"public_ipv6_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the public IPv6 address of the EIP.`,
			},
			"publicip_pool_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the public IP pool name of the EIP.`,
			},
			"publicip_pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the public IP pool ID of the EIP.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description of the EIP.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the tags of the EIP.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the type of the EIP.`,
			},
			"vnic": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcEipVnicSchema(),
				Description: `Specifies the VNIC information of the EIP.`,
			},
			"bandwidth": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcEipBandWidthSchema(),
				Description: `Specifies the bandwidth information of the EIP.`,
			},
			"lock_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the freeze status of the public IP.`,
			},
			"associate_instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the instance type of the public IP.`,
			},
			"associate_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the instance ID of the public IP.`,
			},
			"allow_share_bandwidth_types": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of shared bandwidth types that the public IP can join.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the creation UTC time of the public IP.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the update UTC time of the public IP.`,
			},
			"public_border_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the center site asset or edge site asset. Value range: center, edge site name.`,
			},
		},
	}
}

func buildAssetPropertiesHwcEipBandWidthSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the bandwidth ID of the public IP.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the bandwidth size of the public IP.`,
			},
			"share_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the bandwidth type of the public IP.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the bandwidth name of the public IP.`,
			},
		},
	}
}

func buildAssetPropertiesHwcEipVnicSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"private_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the private IP address of the public IP.`,
			},
			"device_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the device ID of the public IP.`,
			},
			"device_owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the device owner of the public IP.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the virtual private cloud ID of the public IP.`,
			},
			"port_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the port ID of the public IP.`,
			},
			"port_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the port profile information of the public IP.`,
			},
			"mac": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the MAC address of the public IP.`,
			},
			"vtep": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the VTEP IP of the public IP.`,
			},
			"vni": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the VXLAN ID of the public IP.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the instance ID of the public IP.`,
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the instance type of the public IP.`,
			},
		},
	}
}

func buildAssetPropertiesHwcEcsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the ID of the ECS.", utils.SchemaDescInput{Required: true}),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the name of the ECS.", utils.SchemaDescInput{Required: true}),
			},
			"protected_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the protection status of the ECS.",
					utils.SchemaDescInput{Required: true}),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the description of the ECS.",
					utils.SchemaDescInput{Required: true}),
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the status of the ECS.", utils.SchemaDescInput{Required: true}),
			},
			"locked": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies whether the ECS is locked.",
					utils.SchemaDescInput{Required: true}),
			},
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the user ID of the ECS.", utils.SchemaDescInput{Required: true}),
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the project ID of the ECS.",
					utils.SchemaDescInput{Required: true}),
			},
			"host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the host ID of the ECS.", utils.SchemaDescInput{Required: true}),
			},
			"host_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the host name of the ECS.", utils.SchemaDescInput{Required: true}),
			},
			"host_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the host status of the ECS.",
					utils.SchemaDescInput{Required: true}),
			},
			"addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     buildAssetPropertiesHwcEcsAddressesSchema(),
				Description: utils.SchemaDesc("Specifies the addresses of the ECS.",
					utils.SchemaDescInput{Required: true}),
			},
			"security_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     buildAssetPropertiesHwcEcsSecurityGroupsSchema(),
				Description: utils.SchemaDesc("Specifies the security groups of the ECS.",
					utils.SchemaDescInput{Required: true}),
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the availability zone of the ECS.",
					utils.SchemaDescInput{Required: true}),
			},
			"volumes_attached": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     buildAssetPropertiesHwcEcsVolumesAttachedSchema(),
				Description: utils.SchemaDesc("Specifies the volumes attached to the ECS.",
					utils.SchemaDescInput{Required: true}),
			},
			"metadata": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcEcsMetadataSchema(),
				Description: utils.SchemaDesc("Specifies the metadata of the ECS.", utils.SchemaDescInput{Required: true}),
			},
			"updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the last update time of the ECS.",
					utils.SchemaDescInput{Required: true}),
			},
			"created": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the creation time of the ECS server.",
					utils.SchemaDescInput{Required: true}),
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project ID of the ECS.`,
			},
			"flavor": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcEcsFlavorSchema(),
				Description: `Specifies the flavor of the ECS.`,
			},
			"key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the key name of the ECS server.`,
			},
			"scheduler_hints": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcEcsSchedulerHintsSchema(),
				Description: `Specifies the scheduler hints of the ECS.`,
			},
			"hypervisor": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        buildAssetPropertiesHwcEcsHypervisorSchema(),
				Description: `Specifies the virtualization information of the ECS.`,
			},
		},
	}
}

func buildAssetPropertiesHwcEcsHypervisorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"csd_hypervisor": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					"Specifies the Reserved attribute.",
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"hypervisor_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the virtualization type.`,
			},
		},
	}
}

func buildAssetPropertiesHwcEcsSchedulerHintsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the cloud server group ID.`,
			},
			"tenancy": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the tenancy of the ECS.`,
			},
			"dedicated_host_id": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the dedicated host ID.`,
			},
		},
	}
}

func buildAssetPropertiesHwcEcsMetadataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the image ID of the ECS.`,
			},
			"image_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the image type of the ECS.`,
			},
			"image_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the image name of the ECS.`,
			},
			"os_bit": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the OS bit of the ECS.`,
			},
			"os_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the OS type of the ECS.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the VPC ID of the ECS.`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the resource spec code of the ECS.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the resource type of the ECS.`,
			},
			"agency_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the agency name of the ECS.`,
			},
		},
	}
}

func buildAssetPropertiesHwcEcsVolumesAttachedSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the ID of the disk.", utils.SchemaDescInput{Required: true}),
			},
			"delete_on_termination": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to delete the disk when deleting the ECS.`,
			},
			"boot_index": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the boot order of the disk.`,
			},
			"device": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the mount point of the disk.`,
			},
		},
	}
}

func buildAssetPropertiesHwcEcsFlavorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the ID of the flavor.", utils.SchemaDescInput{Required: true}),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the name of the flavor.", utils.SchemaDescInput{Required: true}),
			},
			"disk": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the disk size of the flavor.`,
			},
			"vcpus": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the vcpus of the flavor.`,
			},
			"ram": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the ram of the flavor.`,
			},
		},
	}
}

func buildAssetPropertiesHwcEcsSecurityGroupsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the ID of the security group.",
					utils.SchemaDescInput{Required: true}),
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of the security group.`,
			},
		},
	}
}

func buildAssetPropertiesHwcEcsAddressesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the version of the address.",
					utils.SchemaDescInput{Required: true}),
			},
			"addr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the address.", utils.SchemaDescInput{Required: true}),
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the type of the address.", utils.SchemaDescInput{Required: true}),
			},
			"mac_addr": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the MAC address.", utils.SchemaDescInput{Required: true}),
			},
			"port_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the ID of the port.", utils.SchemaDescInput{Required: true}),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the ID of the VPC.", utils.SchemaDescInput{Required: true}),
			},
		},
	}
}

func buildAssetGovernanceUserSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the governance user type of the asset.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the governance user name of the asset.`,
			},
		},
	}
}

func buildAssetDepartmentSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the department name of the asset.`,
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the department ID of the asset.`,
			},
		},
	}
}

func buildAssetEnvironmentSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vendor_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc("Specifies the environment vendor type.",
					utils.SchemaDescInput{Required: true}),
			},
			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the domain ID.", utils.SchemaDescInput{Required: true}),
			},
			"vendor_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the vendor name.", utils.SchemaDescInput{Required: true}),
			},
			"idc_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("Specifies the IDC name.", utils.SchemaDescInput{Required: true}),
			},
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region ID.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the project ID.`,
			},
			"ep_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"ep_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project name.`,
			},
			"idc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the IDC ID.`,
			},
		},
	}
	return sc
}

func buildDataObjectEnvironmentBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"vendor_type": rawMap["vendor_type"],
		"domain_id":   rawMap["domain_id"],
		"vendor_name": rawMap["vendor_name"],
		"idc_name":    rawMap["idc_name"],
		"region_id":   utils.ValueIgnoreEmpty(rawMap["region_id"]),
		"project_id":  utils.ValueIgnoreEmpty(rawMap["project_id"]),
		"ep_id":       utils.ValueIgnoreEmpty(rawMap["ep_id"]),
		"ep_name":     utils.ValueIgnoreEmpty(rawMap["ep_name"]),
		"idc_id":      utils.ValueIgnoreEmpty(rawMap["idc_id"]),
	}
}

func buildHwcEcsAddressesBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"version":  rawMap["version"],
			"addr":     rawMap["addr"],
			"type":     rawMap["type"],
			"mac_addr": rawMap["mac_addr"],
			"port_id":  rawMap["port_id"],
			"vpc_id":   rawMap["vpc_id"],
		})
	}

	return rst
}

func buildHwcEcsSecurityGroupsBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"id":   rawMap["id"],
			"name": utils.ValueIgnoreEmpty(rawMap["name"]),
		})
	}

	return rst
}

func buildHwcEcsVolumesAttachedBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"id":                    rawMap["id"],
			"delete_on_termination": utils.ValueIgnoreEmpty(rawMap["delete_on_termination"]),
			"boot_index":            utils.ValueIgnoreEmpty(rawMap["boot_index"]),
			"device":                utils.ValueIgnoreEmpty(rawMap["device"]),
		})
	}

	return rst
}

func buildHwcEcsMetadataBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"image_id":           utils.ValueIgnoreEmpty(rawMap["image_id"]),
		"image_type":         utils.ValueIgnoreEmpty(rawMap["image_type"]),
		"image_name":         utils.ValueIgnoreEmpty(rawMap["image_name"]),
		"os_bit":             utils.ValueIgnoreEmpty(rawMap["os_bit"]),
		"os_type":            utils.ValueIgnoreEmpty(rawMap["os_type"]),
		"vpc_id":             utils.ValueIgnoreEmpty(rawMap["vpc_id"]),
		"resource_spec_code": utils.ValueIgnoreEmpty(rawMap["resource_spec_code"]),
		"resource_type":      utils.ValueIgnoreEmpty(rawMap["resource_type"]),
		"agency_name":        utils.ValueIgnoreEmpty(rawMap["agency_name"]),
	}
}

func buildHwcEcsFlavorBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"id":    rawMap["id"],
		"name":  rawMap["name"],
		"disk":  utils.ValueIgnoreEmpty(rawMap["disk"]),
		"vcpus": utils.ValueIgnoreEmpty(rawMap["vcpus"]),
		"ram":   utils.ValueIgnoreEmpty(rawMap["ram"]),
	}
}

func buildHwcEcsSchedulerHintsBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"group":             utils.ValueIgnoreEmpty(rawMap["group"]),
		"tenancy":           utils.ValueIgnoreEmpty(rawMap["tenancy"]),
		"dedicated_host_id": utils.ValueIgnoreEmpty(rawMap["dedicated_host_id"]),
	}
}

func buildHwcEcsHypervisorBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"csd_hypervisor":  utils.ValueIgnoreEmpty(rawMap["csd_hypervisor"]),
		"hypervisor_type": utils.ValueIgnoreEmpty(rawMap["hypervisor_type"]),
	}
}

func buildPropertiesHwcEcsBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"id":                    rawMap["id"],
		"name":                  rawMap["name"],
		"protected_status":      rawMap["protected_status"],
		"description":           rawMap["description"],
		"status":                rawMap["status"],
		"locked":                rawMap["locked"],
		"user_id":               rawMap["user_id"],
		"project_id":            rawMap["project_id"],
		"host_id":               rawMap["host_id"],
		"host_name":             rawMap["host_name"],
		"host_status":           rawMap["host_status"],
		"addresses":             buildHwcEcsAddressesBodyParams(rawMap["addresses"].([]interface{})),
		"security_groups":       buildHwcEcsSecurityGroupsBodyParams(rawMap["security_groups"].([]interface{})),
		"availability_zone":     rawMap["availability_zone"],
		"volumes_attached":      buildHwcEcsVolumesAttachedBodyParams(rawMap["volumes_attached"].([]interface{})),
		"metadata":              buildHwcEcsMetadataBodyParams(rawMap["metadata"].([]interface{})),
		"updated":               rawMap["updated"],
		"created":               rawMap["created"],
		"enterprise_project_id": utils.ValueIgnoreEmpty(rawMap["enterprise_project_id"]),
		"flavor":                buildHwcEcsFlavorBodyParams(rawMap["flavor"].([]interface{})),
		"key_name":              utils.ValueIgnoreEmpty(rawMap["key_name"]),
		"scheduler_hints":       buildHwcEcsSchedulerHintsBodyParams(rawMap["scheduler_hints"].([]interface{})),
		"hypervisor":            buildHwcEcsHypervisorBodyParams(rawMap["hypervisor"].([]interface{})),
	}
}

func buildHwcEipVnicBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"private_ip_address": utils.ValueIgnoreEmpty(rawMap["private_ip_address"]),
		"device_id":          utils.ValueIgnoreEmpty(rawMap["device_id"]),
		"device_owner":       utils.ValueIgnoreEmpty(rawMap["device_owner"]),
		"vpc_id":             utils.ValueIgnoreEmpty(rawMap["vpc_id"]),
		"port_id":            utils.ValueIgnoreEmpty(rawMap["port_id"]),
		"port_profile":       utils.ValueIgnoreEmpty(rawMap["port_profile"]),
		"mac":                utils.ValueIgnoreEmpty(rawMap["mac"]),
		"vtep":               utils.ValueIgnoreEmpty(rawMap["vtep"]),
		"vni":                utils.ValueIgnoreEmpty(rawMap["vni"]),
		"instance_id":        utils.ValueIgnoreEmpty(rawMap["instance_id"]),
		"instance_type":      utils.ValueIgnoreEmpty(rawMap["instance_type"]),
	}
}

func buildHwcEipBandwidthBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"id":         utils.ValueIgnoreEmpty(rawMap["id"]),
		"size":       utils.ValueIgnoreEmpty(rawMap["size"]),
		"share_type": utils.ValueIgnoreEmpty(rawMap["share_type"]),
		"name":       utils.ValueIgnoreEmpty(rawMap["name"]),
	}
}

func buildPropertiesHwcEipBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"id":                          rawMap["id"],
		"alias":                       rawMap["alias"],
		"protected_status":            rawMap["protected_status"],
		"project_id":                  rawMap["project_id"],
		"enterprise_project_id":       rawMap["enterprise_project_id"],
		"ip_version":                  rawMap["ip_version"],
		"status":                      rawMap["status"],
		"public_ip_address":           utils.ValueIgnoreEmpty(rawMap["public_ip_address"]),
		"public_ipv6_address":         utils.ValueIgnoreEmpty(rawMap["public_ipv6_address"]),
		"publicip_pool_name":          utils.ValueIgnoreEmpty(rawMap["publicip_pool_name"]),
		"publicip_pool_id":            utils.ValueIgnoreEmpty(rawMap["publicip_pool_id"]),
		"description":                 utils.ValueIgnoreEmpty(rawMap["description"]),
		"tags":                        utils.ValueIgnoreEmpty(rawMap["tags"]),
		"type":                        utils.ValueIgnoreEmpty(rawMap["type"]),
		"vnic":                        buildHwcEipVnicBodyParams(rawMap["vnic"].([]interface{})),
		"bandwidth":                   buildHwcEipBandwidthBodyParams(rawMap["bandwidth"].([]interface{})),
		"lock_status":                 utils.ValueIgnoreEmpty(rawMap["lock_status"]),
		"associate_instance_type":     utils.ValueIgnoreEmpty(rawMap["associate_instance_type"]),
		"associate_instance_id":       utils.ValueIgnoreEmpty(rawMap["associate_instance_id"]),
		"allow_share_bandwidth_types": utils.ValueIgnoreEmpty(rawMap["allow_share_bandwidth_types"]),
		"created_at":                  utils.ValueIgnoreEmpty(rawMap["created_at"]),
		"updated_at":                  utils.ValueIgnoreEmpty(rawMap["updated_at"]),
		"public_border_group":         utils.ValueIgnoreEmpty(rawMap["public_border_group"]),
	}
}

func buildHwcVpcCloudResourcesBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"resource_type":  utils.ValueIgnoreEmpty(rawMap["resource_type"]),
			"resource_count": utils.ValueIgnoreEmpty(rawMap["resource_count"]),
		})
	}

	return rst
}

func buildHwcVpcTagsBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"key":    rawMap["key"],
			"values": rawMap["values"],
		})
	}

	return rst
}

func buildPropertiesHwcVpcBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"id":                    rawMap["id"],
		"name":                  rawMap["name"],
		"protected_status":      rawMap["protected_status"],
		"status":                rawMap["status"],
		"project_id":            rawMap["project_id"],
		"created_at":            rawMap["created_at"],
		"updated_at":            rawMap["updated_at"],
		"description":           utils.ValueIgnoreEmpty(rawMap["description"]),
		"cidr":                  utils.ValueIgnoreEmpty(rawMap["cidr"]),
		"extend_cidrs":          utils.ValueIgnoreEmpty(rawMap["extend_cidrs"]),
		"enterprise_project_id": utils.ValueIgnoreEmpty(rawMap["enterprise_project_id"]),
		"cloud_resources":       buildHwcVpcCloudResourcesBodyParams(rawMap["cloud_resources"].([]interface{})),
		"tags":                  buildHwcVpcTagsBodyParams(rawMap["tags"].([]interface{})),
	}
}

func buildHwcSubnetSecurityGroupRulesBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"id":                      rawMap["id"],
			"security_group_id":       rawMap["security_group_id"],
			"direction":               rawMap["direction"],
			"protocol":                rawMap["protocol"],
			"ethertype":               rawMap["ethertype"],
			"multiport":               rawMap["multiport"],
			"action":                  rawMap["action"],
			"priority":                rawMap["priority"],
			"created_at":              rawMap["created_at"],
			"updated_at":              rawMap["updated_at"],
			"project_id":              rawMap["project_id"],
			"description":             utils.ValueIgnoreEmpty(rawMap["description"]),
			"remote_group_id":         utils.ValueIgnoreEmpty(rawMap["remote_group_id"]),
			"remote_ip_prefix":        utils.ValueIgnoreEmpty(rawMap["remote_ip_prefix"]),
			"remote_address_group_id": utils.ValueIgnoreEmpty(rawMap["remote_address_group_id"]),
		})
	}

	return rst
}

func buildPropertiesHwcSubnetBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"id":                    rawMap["id"],
		"name":                  rawMap["name"],
		"project_id":            rawMap["project_id"],
		"created_at":            rawMap["created_at"],
		"updated_at":            rawMap["updated_at"],
		"description":           utils.ValueIgnoreEmpty(rawMap["description"]),
		"enterprise_project_id": utils.ValueIgnoreEmpty(rawMap["enterprise_project_id"]),
		"security_group_rules":  buildHwcSubnetSecurityGroupRulesBodyParams(rawMap["security_group_rules"].([]interface{})),
	}
}

func buildHwcRdsHaBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"replication_mode": rawMap["replication_mode"],
	}
}

func buildHwcRdsDatastoreBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"type":             rawMap["type"],
		"version":          rawMap["version"],
		"complete_version": utils.ValueIgnoreEmpty(rawMap["complete_version"]),
	}
}

func buildHwcRdsVolumeBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"type": rawMap["type"],
		"size": rawMap["size"],
	}
}

func buildHwcRdsBackupStrategyBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"start_time": utils.ValueIgnoreEmpty(rawMap["start_time"]),
		"keep_days":  utils.ValueIgnoreEmpty(rawMap["keep_days"]),
	}
}

func buildHwcRdsNodesBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"id":                rawMap["id"],
			"name":              rawMap["name"],
			"role":              rawMap["role"],
			"status":            rawMap["status"],
			"availability_zone": rawMap["availability_zone"],
		})
	}

	return rst
}

func buildHwcRdsRelatedInstanceBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"id":   rawMap["id"],
			"type": rawMap["type"],
		})
	}

	return rst
}

func buildHwcRdsTagsBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"key":    rawMap["key"],
			"values": rawMap["values"],
		})
	}

	return rst
}

func buildPropertiesHwcRdsBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"id":                    rawMap["id"],
		"name":                  rawMap["name"],
		"protected_status":      rawMap["protected_status"],
		"status":                rawMap["status"],
		"port":                  rawMap["port"],
		"enable_ssl":            rawMap["enable_ssl"],
		"type":                  rawMap["type"],
		"ha":                    buildHwcRdsHaBodyParams(rawMap["ha"].([]interface{})),
		"region":                rawMap["region"],
		"datastore":             buildHwcRdsDatastoreBodyParams(rawMap["datastore"].([]interface{})),
		"created":               rawMap["created"],
		"vpc_id":                rawMap["vpc_id"],
		"subnet_id":             rawMap["subnet_id"],
		"security_group_id":     rawMap["security_group_id"],
		"flavor_ref":            rawMap["flavor_ref"],
		"cpu":                   rawMap["cpu"],
		"mem":                   rawMap["mem"],
		"volume":                buildHwcRdsVolumeBodyParams(rawMap["volume"].([]interface{})),
		"project_id":            rawMap["project_id"],
		"switch_strategy":       rawMap["switch_strategy"],
		"read_only_by_user":     rawMap["idread_only_by_user"],
		"backup_strategy":       buildHwcRdsBackupStrategyBodyParams(rawMap["backup_strategy"].([]interface{})),
		"maintenance_window":    rawMap["maintenance_window"],
		"nodes":                 buildHwcRdsNodesBodyParams(rawMap["nodes"].([]interface{})),
		"related_instance":      buildHwcRdsRelatedInstanceBodyParams(rawMap["related_instance"].([]interface{})),
		"time_zone":             rawMap["time_zone"],
		"storage_used_space":    rawMap["storage_used_space"],
		"associated_with_ddm":   rawMap["associated_with_ddm"],
		"max_iops":              rawMap["max_iops"],
		"expiration_time":       rawMap["expiration_time"],
		"alias":                 utils.ValueIgnoreEmpty(rawMap["alias"]),
		"private_ips":           utils.ValueIgnoreEmpty(rawMap["private_ips"]),
		"private_dns_names":     utils.ValueIgnoreEmpty(rawMap["private_dns_names"]),
		"public_ips":            utils.ValueIgnoreEmpty(rawMap["public_ips"]),
		"updated":               utils.ValueIgnoreEmpty(rawMap["updated"]),
		"db_user_name":          utils.ValueIgnoreEmpty(rawMap["db_user_name"]),
		"tags":                  buildHwcRdsTagsBodyParams(rawMap["tags"].([]interface{})),
		"enterprise_project_id": utils.ValueIgnoreEmpty(rawMap["enterprise_project_id"]),
		"disk_encryption_id":    utils.ValueIgnoreEmpty(rawMap["disk_encryption_id"]),
		"backup_used_space":     utils.ValueIgnoreEmpty(rawMap["backup_used_space"]),
	}
}

func buildHwcDomainContactBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"email":        rawMap["email"],
			"register":     rawMap["register"],
			"contact_name": rawMap["contact_name"],
			"phone_num":    rawMap["phone_num"],
			"province":     rawMap["province"],
			"city":         rawMap["city"],
			"address":      rawMap["address"],
			"zip_code":     rawMap["zip_code"],
		})
	}

	return rst
}

func buildPropertiesHwcDomainBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"domain_name":         rawMap["domain_name"],
		"expire_date":         rawMap["expire_date"],
		"status":              rawMap["status"],
		"audit_status":        rawMap["audit_status"],
		"audit_unpass_reason": rawMap["audit_unpass_reason"],
		"reg_type":            rawMap["reg_type"],
		"privacy_protection":  rawMap["privacy_protection"],
		"name_server":         rawMap["name_server"],
		"credential_type":     rawMap["credential_type"],
		"credential_id":       rawMap["credential_id"],
		"registrar":           rawMap["registrar"],
		"contact":             buildHwcDomainContactBodyParams(rawMap["contact"].([]interface{})),
		"transfer_status":     utils.ValueIgnoreEmpty(rawMap["transfer_status"]),
	}
}

func buildHwcWebsiteExtendPropertitesBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"mac_addr": utils.ValueIgnoreEmpty(rawMap["mac_addr"]),
	}
}

func buildPropertiesWebsiteBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"value":              rawMap["value"],
		"main_domain":        rawMap["main_domain"],
		"protected_status":   rawMap["protected_status"],
		"is_public":          rawMap["is_public"],
		"name_server":        rawMap["name_server"],
		"remark":             utils.ValueIgnoreEmpty(rawMap["remark"]),
		"extend_propertites": buildHwcWebsiteExtendPropertitesBodyParams(rawMap["extend_propertites"].([]interface{})),
	}
}

func buildHwcOcaIpNetworkBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"is_public": rawMap["is_public"],
		"partition": utils.ValueIgnoreEmpty(rawMap["partition"]),
		"plane":     utils.ValueIgnoreEmpty(rawMap["plane"]),
		"vxlan_id":  utils.ValueIgnoreEmpty(rawMap["vxlan_id"]),
	}
}

func buildHwcOcaIpDataCenterBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"city_code":    rawMap["city_code"],
		"country_code": rawMap["country_code"],
		"latitude":     utils.ValueIgnoreEmpty(rawMap["latitude"]),
		"longitude":    utils.ValueIgnoreEmpty(rawMap["longitude"]),
	}
}

func buildHwcOcaIpDeviceVendorBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"name":  utils.ValueIgnoreEmpty(rawMap["name"]),
		"is_xc": utils.ValueIgnoreEmpty(rawMap["is_xc"]),
	}
}

func buildHwcOcaIpDeviceBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"type":   utils.ValueIgnoreEmpty(rawMap["type"]),
		"model":  utils.ValueIgnoreEmpty(rawMap["model"]),
		"vendor": buildHwcOcaIpDeviceVendorBodyParams(rawMap["vendor"].([]interface{})),
	}
}

func buildHwcOcaIpSystemVendorBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"name":  utils.ValueIgnoreEmpty(rawMap["name"]),
		"is_xc": utils.ValueIgnoreEmpty(rawMap["is_xc"]),
	}
}

func buildHwcOcaIpSystemBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"family":  utils.ValueIgnoreEmpty(rawMap["family"]),
		"name":    utils.ValueIgnoreEmpty(rawMap["name"]),
		"version": utils.ValueIgnoreEmpty(rawMap["version"]),
		"vendor":  buildHwcOcaIpSystemVendorBodyParams(rawMap["vendor"].([]interface{})),
	}
}

func buildHwcOcaIpServiceVendorBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"name":  utils.ValueIgnoreEmpty(rawMap["name"]),
		"is_xc": utils.ValueIgnoreEmpty(rawMap["is_xc"]),
	}
}

func buildHwcOcaIpServicesBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"port":     utils.ValueIgnoreEmpty(rawMap["port"]),
			"protocol": utils.ValueIgnoreEmpty(rawMap["protocol"]),
			"name":     utils.ValueIgnoreEmpty(rawMap["name"]),
			"version":  utils.ValueIgnoreEmpty(rawMap["version"]),
			"vendor":   buildHwcOcaIpServiceVendorBodyParams(rawMap["vendor"].([]interface{})),
		})
	}

	return rst
}

func buildHwcOcaIpExtendPropertitesBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"device":   buildHwcOcaIpDeviceBodyParams(rawMap["device"].([]interface{})),
		"system":   buildHwcOcaIpSystemBodyParams(rawMap["system"].([]interface{})),
		"services": buildHwcOcaIpServicesBodyParams(rawMap["services"].([]interface{})),
	}
}

func buildPropertiesOcaIpBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"value":              rawMap["value"],
		"version":            rawMap["version"],
		"network":            buildHwcOcaIpNetworkBodyParams(rawMap["network"].([]interface{})),
		"server_room":        rawMap["server_room"],
		"server_rack":        rawMap["server_rack"],
		"data_center":        buildHwcOcaIpDataCenterBodyParams(rawMap["data_center"].([]interface{})),
		"remark":             utils.ValueIgnoreEmpty(rawMap["remark"]),
		"name":               utils.ValueIgnoreEmpty(rawMap["name"]),
		"relative_value":     utils.ValueIgnoreEmpty(rawMap["relative_value"]),
		"mac_addr":           utils.ValueIgnoreEmpty(rawMap["mac_addr"]),
		"important":          utils.ValueIgnoreEmpty(rawMap["important"]),
		"extend_propertites": buildHwcOcaIpExtendPropertitesBodyParams(rawMap["extend_propertites"].([]interface{})),
	}
}

func buildDataObjectPropertiesBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"hwc_ecs":    buildPropertiesHwcEcsBodyParams(rawMap["hwc_ecs"].([]interface{})),
		"hwc_eip":    buildPropertiesHwcEipBodyParams(rawMap["hwc_eip"].([]interface{})),
		"hwc_vpc":    buildPropertiesHwcVpcBodyParams(rawMap["hwc_vpc"].([]interface{})),
		"hwc_subnet": buildPropertiesHwcSubnetBodyParams(rawMap["hwc_subnet"].([]interface{})),
		"hwc_rds":    buildPropertiesHwcRdsBodyParams(rawMap["hwc_rds"].([]interface{})),
		"hwc_domain": buildPropertiesHwcDomainBodyParams(rawMap["hwc_domain"].([]interface{})),
		"website":    buildPropertiesWebsiteBodyParams(rawMap["website"].([]interface{})),
		"oca_ip":     buildPropertiesOcaIpBodyParams(rawMap["oca_ip"].([]interface{})),
	}
}

func buildDataObjectDepartmentBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"name": utils.ValueIgnoreEmpty(rawMap["name"]),
		"id":   utils.ValueIgnoreEmpty(rawMap["id"]),
	}
}

func buildDataObjectGovernanceUserBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"type": utils.ValueIgnoreEmpty(rawMap["type"]),
		"name": utils.ValueIgnoreEmpty(rawMap["name"]),
	}
}

func buildAssetDataObjectBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"id":                 rawMap["id"],
		"name":               rawMap["name"],
		"provider":           rawMap["provider"],
		"type":               rawMap["type"],
		"environment":        buildDataObjectEnvironmentBodyParams(rawMap["environment"].([]interface{})),
		"properties":         buildDataObjectPropertiesBodyParams(rawMap["properties"].([]interface{})),
		"checksum":           utils.ValueIgnoreEmpty(rawMap["checksum"]),
		"created":            utils.ValueIgnoreEmpty(rawMap["created"]),
		"provisioning_state": utils.ValueIgnoreEmpty(rawMap["provisioning_state"]),
		"department":         buildDataObjectDepartmentBodyParams(rawMap["department"].([]interface{})),
		"governance_user":    buildDataObjectGovernanceUserBodyParams(rawMap["governance_user"].([]interface{})),
		"level":              utils.ValueIgnoreEmpty(rawMap["level"]),
	}
}

func buildAssetBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"data_object": buildAssetDataObjectBodyParams(d.Get("data_object").([]interface{})),
	}
}

func updateAsset(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/sa/resources/{id}"
		workspaceID = d.Get("workspace_id").(string)
		assetID     = d.Get("asset_id").(string)
	)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceID)
	requestPath = strings.ReplaceAll(requestPath, "{id}", assetID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildAssetBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourceAssetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	if err := updateAsset(client, d); err != nil {
		return diag.Errorf("error updating SecMaster asset in creation operation: %s", err)
	}

	d.SetId(d.Get("asset_id").(string))

	return resourceAssetRead(ctx, d, meta)
}

func ReadAssetDetail(client *golangsdk.ServiceClient, workspaceID, assetID string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/sa/resources/{id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceID)
	requestPath = strings.ReplaceAll(requestPath, "{id}", assetID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func flattenAssetEnvironmentAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"vendor_type": utils.PathSearch("vendor_type", respBody, nil),
		"domain_id":   utils.PathSearch("domain_id", respBody, nil),
		"region_id":   utils.PathSearch("region_id", respBody, nil),
		"project_id":  utils.PathSearch("project_id", respBody, nil),
		"ep_id":       utils.PathSearch("ep_id", respBody, nil),
		"ep_name":     utils.PathSearch("ep_name", respBody, nil),
		"vendor_name": utils.PathSearch("vendor_name", respBody, nil),
		"idc_name":    utils.PathSearch("idc_name", respBody, nil),
		"idc_id":      utils.PathSearch("idc_id", respBody, nil),
	}

	return []interface{}{rstBody}
}

func flattenAssetDepartmentAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"name": utils.PathSearch("name", respBody, nil),
		"id":   utils.PathSearch("id", respBody, nil),
	}

	return []interface{}{rstBody}
}

func flattenAssetGovernanceUserAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"type": utils.PathSearch("type", respBody, nil),
		"name": utils.PathSearch("name", respBody, nil),
	}

	return []interface{}{rstBody}
}

// Array Objects
func flattenHwcEcsAddressesAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"version":  utils.PathSearch("version", v, nil),
			"addr":     utils.PathSearch("addr", v, nil),
			"type":     utils.PathSearch("type", v, nil),
			"mac_addr": utils.PathSearch("mac_addr", v, nil),
			"port_id":  utils.PathSearch("port_id", v, nil),
			"vpc_id":   utils.PathSearch("vpc_id", v, nil),
		})
	}

	return rst
}

// Array Objects
func flattenHwcEcsSecurityGroupsAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"name": utils.PathSearch("name", v, nil),
		})
	}

	return rst
}

// Object
func flattenHwcEcsFlavorAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"id":    utils.PathSearch("id", respBody, nil),
		"name":  utils.PathSearch("name", respBody, nil),
		"disk":  utils.PathSearch("disk", respBody, nil),
		"vcpus": utils.PathSearch("vcpus", respBody, nil),
		"ram":   utils.PathSearch("ram", respBody, nil),
	}

	return []interface{}{rstBody}
}

// Array Objects
func flattenHwcEcsVolumesAttachedAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"delete_on_termination": utils.PathSearch("delete_on_termination", v, nil),
			"boot_index":            utils.PathSearch("boot_index", v, nil),
			"device":                utils.PathSearch("device", v, nil),
		})
	}

	return rst
}

// Object
func flattenHwcEcsMetadataAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"image_id":           utils.PathSearch("image_id", respBody, nil),
		"image_type":         utils.PathSearch("image_type", respBody, nil),
		"image_name":         utils.PathSearch("image_name", respBody, nil),
		"os_bit":             utils.PathSearch("os_bit", respBody, nil),
		"os_type":            utils.PathSearch("os_type", respBody, nil),
		"vpc_id":             utils.PathSearch("vpc_id", respBody, nil),
		"resource_spec_code": utils.PathSearch("resource_spec_code", respBody, nil),
		"resource_type":      utils.PathSearch("resource_type", respBody, nil),
		"agency_name":        utils.PathSearch("agency_name", respBody, nil),
	}

	return []interface{}{rstBody}
}

// Object
func flattenHwcEcsSchedulerHintsAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"group":             utils.PathSearch("group", respBody, nil),
		"tenancy":           utils.PathSearch("tenancy", respBody, nil),
		"dedicated_host_id": utils.PathSearch("dedicated_host_id", respBody, nil),
	}

	return []interface{}{rstBody}
}

// Object
func flattenHwcEcsHypervisorAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"hypervisor_type": utils.PathSearch("hypervisor_type", respBody, nil),
		"csd_hypervisor":  utils.PathSearch("csd_hypervisor", respBody, nil),
	}

	return []interface{}{rstBody}
}

func flattenPropertiesHwcEcsAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"id":                    utils.PathSearch("id", respBody, nil),
		"name":                  utils.PathSearch("name", respBody, nil),
		"protected_status":      utils.PathSearch("protected_status", respBody, nil),
		"description":           utils.PathSearch("description", respBody, nil),
		"status":                utils.PathSearch("status", respBody, nil),
		"locked":                utils.PathSearch("locked", respBody, nil),
		"enterprise_project_id": utils.PathSearch("enterprise_project_id", respBody, nil),
		"user_id":               utils.PathSearch("user_id", respBody, nil),
		"project_id":            utils.PathSearch("project_id", respBody, nil),
		"host_id":               utils.PathSearch("host_id", respBody, nil),
		"host_name":             utils.PathSearch("host_name", respBody, nil),
		"host_status":           utils.PathSearch("host_status", respBody, nil),
		"addresses":             flattenHwcEcsAddressesAttribute(utils.PathSearch("addresses", respBody, nil)),
		"security_groups":       flattenHwcEcsSecurityGroupsAttribute(utils.PathSearch("security_groups", respBody, nil)),
		"availability_zone":     utils.PathSearch("availability_zone", respBody, nil),
		"flavor":                flattenHwcEcsFlavorAttribute(utils.PathSearch("flavor", respBody, nil)),
		"volumes_attached":      flattenHwcEcsVolumesAttachedAttribute(utils.PathSearch("volumes_attached", respBody, nil)),
		"metadata":              flattenHwcEcsMetadataAttribute(utils.PathSearch("metadata", respBody, nil)),
		"updated":               utils.PathSearch("updated", respBody, nil),
		"created":               utils.PathSearch("created", respBody, nil),
		"key_name":              utils.PathSearch("key_name", respBody, nil),
		"scheduler_hints":       flattenHwcEcsSchedulerHintsAttribute(utils.PathSearch("scheduler_hints", respBody, nil)),
		"hypervisor":            flattenHwcEcsHypervisorAttribute(utils.PathSearch("hypervisor", respBody, nil)),
	}

	return []interface{}{rstBody}
}

// Object
func flattenHwcEipVnicAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"private_ip_address": utils.PathSearch("private_ip_address", respBody, nil),
		"device_id":          utils.PathSearch("device_id", respBody, nil),
		"device_owner":       utils.PathSearch("device_owner", respBody, nil),
		"vpc_id":             utils.PathSearch("vpc_id", respBody, nil),
		"port_id":            utils.PathSearch("port_id", respBody, nil),
		"port_profile":       utils.PathSearch("port_profile", respBody, nil),
		"mac":                utils.PathSearch("mac", respBody, nil),
		"vtep":               utils.PathSearch("vtep", respBody, nil),
		"vni":                utils.PathSearch("vni", respBody, nil),
		"instance_id":        utils.PathSearch("instance_id", respBody, nil),
		"instance_type":      utils.PathSearch("instance_type", respBody, nil),
	}

	return []interface{}{rstBody}
}

// Object
func flattenHwcEipBandwidthAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"id":         utils.PathSearch("id", respBody, nil),
		"size":       utils.PathSearch("size", respBody, nil),
		"share_type": utils.PathSearch("share_type", respBody, nil),
		"name":       utils.PathSearch("name", respBody, nil),
	}

	return []interface{}{rstBody}
}

func flattenPropertiesHwcEipAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"id":                          utils.PathSearch("id", respBody, nil),
		"alias":                       utils.PathSearch("alias", respBody, nil),
		"protected_status":            utils.PathSearch("protected_status", respBody, nil),
		"project_id":                  utils.PathSearch("project_id", respBody, nil),
		"enterprise_project_id":       utils.PathSearch("enterprise_project_id", respBody, nil),
		"ip_version":                  utils.PathSearch("ip_version", respBody, nil),
		"public_ip_address":           utils.PathSearch("public_ip_address", respBody, nil),
		"public_ipv6_address":         utils.PathSearch("public_ipv6_address", respBody, nil),
		"publicip_pool_name":          utils.PathSearch("publicip_pool_name", respBody, nil),
		"publicip_pool_id":            utils.PathSearch("publicip_pool_id", respBody, nil),
		"status":                      utils.PathSearch("status", respBody, nil),
		"description":                 utils.PathSearch("description", respBody, nil),
		"tags":                        utils.PathSearch("tags", respBody, nil),
		"type":                        utils.PathSearch("type", respBody, nil),
		"vnic":                        flattenHwcEipVnicAttribute(utils.PathSearch("vnic", respBody, nil)),
		"bandwidth":                   flattenHwcEipBandwidthAttribute(utils.PathSearch("bandwidth", respBody, nil)),
		"lock_status":                 utils.PathSearch("lock_status", respBody, nil),
		"associate_instance_type":     utils.PathSearch("associate_instance_type", respBody, nil),
		"associate_instance_id":       utils.PathSearch("associate_instance_id", respBody, nil),
		"allow_share_bandwidth_types": utils.PathSearch("allow_share_bandwidth_types", respBody, nil),
		"created_at":                  utils.PathSearch("created_at", respBody, nil),
		"updated_at":                  utils.PathSearch("updated_at", respBody, nil),
		"public_border_group":         utils.PathSearch("public_border_group", respBody, nil),
	}

	return []interface{}{rstBody}
}

// Array Objects
func flattenHwcVpcCloudResourcesAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"resource_type":  utils.PathSearch("resource_type", v, nil),
			"resource_count": utils.PathSearch("resource_count", v, nil),
		})
	}

	return rst
}

// Array Objects
func flattenHwcVpcTagsAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"key":    utils.PathSearch("key", v, nil),
			"values": utils.PathSearch("values", v, nil),
		})
	}

	return rst
}

func flattenPropertiesHwcVpcAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"id":                    utils.PathSearch("id", respBody, nil),
		"name":                  utils.PathSearch("name", respBody, nil),
		"description":           utils.PathSearch("description", respBody, nil),
		"protected_status":      utils.PathSearch("protected_status", respBody, nil),
		"cidr":                  utils.PathSearch("cidr", respBody, nil),
		"extend_cidrs":          utils.PathSearch("extend_cidrs", respBody, nil),
		"status":                utils.PathSearch("status", respBody, nil),
		"project_id":            utils.PathSearch("project_id", respBody, nil),
		"enterprise_project_id": utils.PathSearch("enterprise_project_id", respBody, nil),
		"created_at":            utils.PathSearch("created_at", respBody, nil),
		"updated_at":            utils.PathSearch("updated_at", respBody, nil),
		"cloud_resources":       flattenHwcVpcCloudResourcesAttribute(utils.PathSearch("cloud_resources", respBody, nil)),
		"tags":                  flattenHwcVpcTagsAttribute(utils.PathSearch("tags", respBody, nil)),
	}

	return []interface{}{rstBody}
}

// Array Objects
func flattenHwcSubnetSecurityGroupRulesAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"id":                      utils.PathSearch("id", v, nil),
			"description":             utils.PathSearch("description", v, nil),
			"security_group_id":       utils.PathSearch("security_group_id", v, nil),
			"direction":               utils.PathSearch("direction", v, nil),
			"protocol":                utils.PathSearch("protocol", v, nil),
			"ethertype":               utils.PathSearch("ethertype", v, nil),
			"multiport":               utils.PathSearch("multiport", v, nil),
			"action":                  utils.PathSearch("action", v, nil),
			"priority":                utils.PathSearch("priority", v, nil),
			"remote_group_id":         utils.PathSearch("remote_group_id", v, nil),
			"remote_ip_prefix":        utils.PathSearch("remote_ip_prefix", v, nil),
			"remote_address_group_id": utils.PathSearch("remote_address_group_id", v, nil),
			"created_at":              utils.PathSearch("created_at", v, nil),
			"updated_at":              utils.PathSearch("updated_at", v, nil),
			"project_id":              utils.PathSearch("project_id", v, nil),
		})
	}

	return rst
}

func flattenPropertiesHwcSubnetAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"id":                    utils.PathSearch("id", respBody, nil),
		"name":                  utils.PathSearch("name", respBody, nil),
		"description":           utils.PathSearch("description", respBody, nil),
		"project_id":            utils.PathSearch("project_id", respBody, nil),
		"created_at":            utils.PathSearch("created_at", respBody, nil),
		"updated_at":            utils.PathSearch("updated_at", respBody, nil),
		"enterprise_project_id": utils.PathSearch("enterprise_project_id", respBody, nil),
		"security_group_rules":  flattenHwcSubnetSecurityGroupRulesAttribute(utils.PathSearch("security_group_rules", respBody, nil)),
	}

	return []interface{}{rstBody}
}

// Object
func flattenHwcRdsHaAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"replication_mode": utils.PathSearch("replication_mode", respBody, nil),
	}

	return []interface{}{rstBody}
}

// Object
func flattenHwcRdsDatastoreAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"type":             utils.PathSearch("type", respBody, nil),
		"version":          utils.PathSearch("version", respBody, nil),
		"complete_version": utils.PathSearch("complete_version", respBody, nil),
	}

	return []interface{}{rstBody}
}

// Object
func flattenHwcRdsVolumeAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"type": utils.PathSearch("type", respBody, nil),
		"size": utils.PathSearch("size", respBody, nil),
	}

	return []interface{}{rstBody}
}

// Array Objects
func flattenHwcRdsTagsAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"key":    utils.PathSearch("key", v, nil),
			"values": utils.PathSearch("values", v, nil),
		})
	}

	return rst
}

// Object
func flattenHwcRdsBackupStrategyAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"start_time": utils.PathSearch("start_time", respBody, nil),
		"keep_days":  utils.PathSearch("keep_days", respBody, nil),
	}

	return []interface{}{rstBody}
}

// Array Objects
func flattenHwcRdsNodesAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"role":              utils.PathSearch("role", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"availability_zone": utils.PathSearch("availability_zone", v, nil),
		})
	}

	return rst
}

// Array Objects
func flattenHwcRdsRelatedInstanceAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"id":   utils.PathSearch("id", v, nil),
			"type": utils.PathSearch("type", v, nil),
		})
	}

	return rst
}

func flattenPropertiesHwcRdsAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"id":                    utils.PathSearch("id", respBody, nil),
		"name":                  utils.PathSearch("name", respBody, nil),
		"protected_status":      utils.PathSearch("protected_status", respBody, nil),
		"status":                utils.PathSearch("status", respBody, nil),
		"alias":                 utils.PathSearch("alias", respBody, nil),
		"private_ips":           utils.PathSearch("private_ips", respBody, nil),
		"private_dns_names":     utils.PathSearch("private_dns_names", respBody, nil),
		"public_ips":            utils.PathSearch("public_ips", respBody, nil),
		"port":                  utils.PathSearch("port", respBody, nil),
		"enable_ssl":            utils.PathSearch("enable_ssl", respBody, nil),
		"type":                  utils.PathSearch("type", respBody, nil),
		"ha":                    flattenHwcRdsHaAttribute(utils.PathSearch("ha", respBody, nil)),
		"region":                utils.PathSearch("region", respBody, nil),
		"datastore":             flattenHwcRdsDatastoreAttribute(utils.PathSearch("datastore", respBody, nil)),
		"created":               utils.PathSearch("created", respBody, nil),
		"updated":               utils.PathSearch("updated", respBody, nil),
		"db_user_name":          utils.PathSearch("db_user_name", respBody, nil),
		"vpc_id":                utils.PathSearch("vpc_id", respBody, nil),
		"subnet_id":             utils.PathSearch("subnet_id", respBody, nil),
		"security_group_id":     utils.PathSearch("security_group_id", respBody, nil),
		"flavor_ref":            utils.PathSearch("flavor_ref", respBody, nil),
		"cpu":                   utils.PathSearch("cpu", respBody, nil),
		"mem":                   utils.PathSearch("mem", respBody, nil),
		"volume":                flattenHwcRdsVolumeAttribute(utils.PathSearch("volume", respBody, nil)),
		"tags":                  flattenHwcRdsTagsAttribute(utils.PathSearch("tags", respBody, nil)),
		"enterprise_project_id": utils.PathSearch("enterprise_project_id", respBody, nil),
		"project_id":            utils.PathSearch("project_id", respBody, nil),
		"switch_strategy":       utils.PathSearch("switch_strategy", respBody, nil),
		"read_only_by_user":     utils.PathSearch("read_only_by_user", respBody, nil),
		"backup_strategy":       flattenHwcRdsBackupStrategyAttribute(utils.PathSearch("backup_strategy", respBody, nil)),
		"maintenance_window":    utils.PathSearch("maintenance_window", respBody, nil),
		"nodes":                 flattenHwcRdsNodesAttribute(utils.PathSearch("nodes", respBody, nil)),
		"related_instance":      flattenHwcRdsRelatedInstanceAttribute(utils.PathSearch("related_instance", respBody, nil)),
		"disk_encryption_id":    utils.PathSearch("disk_encryption_id", respBody, nil),
		"time_zone":             utils.PathSearch("time_zone", respBody, nil),
		"backup_used_space":     utils.PathSearch("backup_used_space", respBody, nil),
		"storage_used_space":    utils.PathSearch("storage_used_space", respBody, nil),
		"associated_with_ddm":   utils.PathSearch("associated_with_ddm", respBody, nil),
		"max_iops":              utils.PathSearch("max_iops", respBody, nil),
		"expiration_time":       utils.PathSearch("expiration_time", respBody, nil),
	}

	return []interface{}{rstBody}
}

// Array Objects
func flattenHwcDomainContactAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"email":        utils.PathSearch("email", v, nil),
			"register":     utils.PathSearch("register", v, nil),
			"contact_name": utils.PathSearch("contact_name", v, nil),
			"phone_num":    utils.PathSearch("phone_num", v, nil),
			"province":     utils.PathSearch("province", v, nil),
			"city":         utils.PathSearch("city", v, nil),
			"address":      utils.PathSearch("address", v, nil),
			"zip_code":     utils.PathSearch("zip_code", v, nil),
		})
	}

	return rst
}

func flattenPropertiesHwcDomainAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"domain_name":         utils.PathSearch("domain_name", respBody, nil),
		"expire_date":         utils.PathSearch("expire_date", respBody, nil),
		"status":              utils.PathSearch("status", respBody, nil),
		"audit_status":        utils.PathSearch("audit_status", respBody, nil),
		"audit_unpass_reason": utils.PathSearch("audit_unpass_reason", respBody, nil),
		"transfer_status":     utils.PathSearch("transfer_status", respBody, nil),
		"reg_type":            utils.PathSearch("reg_type", respBody, nil),
		"privacy_protection":  utils.PathSearch("privacy_protection", respBody, nil),
		"name_server":         utils.PathSearch("name_server", respBody, nil),
		"credential_type":     utils.PathSearch("credential_type", respBody, nil),
		"credential_id":       utils.PathSearch("credential_id", respBody, nil),
		"registrar":           utils.PathSearch("registrar", respBody, nil),
		"contact":             flattenHwcDomainContactAttribute(utils.PathSearch("contact", respBody, nil)),
	}

	return []interface{}{rstBody}
}

func flattenWebsiteExtendPropertitesAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"mac_addr": utils.PathSearch("mac_addr", respBody, nil),
	}

	return []interface{}{rstBody}
}

func flattenPropertiesWebsiteAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"value":              utils.PathSearch("value", respBody, nil),
		"main_domain":        utils.PathSearch("main_domain", respBody, nil),
		"protected_status":   utils.PathSearch("protected_status", respBody, nil),
		"is_public":          utils.PathSearch("is_public", respBody, nil),
		"remark":             utils.PathSearch("remark", respBody, nil),
		"name_server":        utils.PathSearch("name_server", respBody, nil),
		"extend_propertites": flattenWebsiteExtendPropertitesAttribute(utils.PathSearch("extend_propertites", respBody, nil)),
	}

	return []interface{}{rstBody}
}

func flattenOcaIpNetworkAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"is_public": utils.PathSearch("is_public", respBody, nil),
		"partition": utils.PathSearch("partition", respBody, nil),
		"plane":     utils.PathSearch("plane", respBody, nil),
		"vxlan_id":  utils.PathSearch("vxlan_id", respBody, nil),
	}

	return []interface{}{rstBody}
}

func flattenOcaIpDataCenterAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"latitude":     utils.PathSearch("latitude", respBody, nil),
		"longitude":    utils.PathSearch("longitude", respBody, nil),
		"city_code":    utils.PathSearch("city_code", respBody, nil),
		"country_code": utils.PathSearch("country_code", respBody, nil),
	}

	return []interface{}{rstBody}
}

func flattenOcaIpDeviceVendorAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"name":  utils.PathSearch("name", respBody, nil),
		"is_xc": utils.PathSearch("is_xc", respBody, nil),
	}

	return []interface{}{rstBody}
}

func flattenOcaIpDeviceAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"type":   utils.PathSearch("type", respBody, nil),
		"model":  utils.PathSearch("model", respBody, nil),
		"vendor": flattenOcaIpDeviceVendorAttribute(utils.PathSearch("vendor", respBody, nil)),
	}

	return []interface{}{rstBody}
}

func flattenOcaIpSystemVendorAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"name":  utils.PathSearch("name", respBody, nil),
		"is_xc": utils.PathSearch("is_xc", respBody, nil),
	}

	return []interface{}{rstBody}
}

func flattenOcaIpSystemAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"family":  utils.PathSearch("family", respBody, nil),
		"name":    utils.PathSearch("name", respBody, nil),
		"version": utils.PathSearch("version", respBody, nil),
		"vendor":  flattenOcaIpSystemVendorAttribute(utils.PathSearch("vendor", respBody, nil)),
	}

	return []interface{}{rstBody}
}

func flattenOcaIpServicesVendorAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"name":  utils.PathSearch("name", respBody, nil),
		"is_xc": utils.PathSearch("is_xc", respBody, nil),
	}

	return []interface{}{rstBody}
}

func flattenOcaIpServicesAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray, ok := respBody.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"port":     utils.PathSearch("port", v, nil),
			"protocol": utils.PathSearch("protocol", v, nil),
			"name":     utils.PathSearch("name", v, nil),
			"version":  utils.PathSearch("version", v, nil),
			"vendor":   flattenOcaIpServicesVendorAttribute(utils.PathSearch("vendor", v, nil)),
		})
	}

	return rst
}

func flattenOcaIpExtendPropertitesAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"device":   flattenOcaIpDeviceAttribute(utils.PathSearch("device", respBody, nil)),
		"system":   flattenOcaIpSystemAttribute(utils.PathSearch("system", respBody, nil)),
		"services": flattenOcaIpServicesAttribute(utils.PathSearch("services", respBody, nil)),
	}

	return []interface{}{rstBody}
}

func flattenPropertiesOcaIpAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"value":              utils.PathSearch("value", respBody, nil),
		"version":            utils.PathSearch("version", respBody, nil),
		"network":            flattenOcaIpNetworkAttribute(utils.PathSearch("network", respBody, nil)),
		"remark":             utils.PathSearch("remark", respBody, nil),
		"name":               utils.PathSearch("name", respBody, nil),
		"relative_value":     utils.PathSearch("relative_value", respBody, nil),
		"server_room":        utils.PathSearch("server_room", respBody, nil),
		"server_rack":        utils.PathSearch("server_rack", respBody, nil),
		"data_center":        flattenOcaIpDataCenterAttribute(utils.PathSearch("data_center", respBody, nil)),
		"mac_addr":           utils.PathSearch("mac_addr", respBody, nil),
		"important":          utils.PathSearch("important", respBody, nil),
		"extend_propertites": flattenOcaIpExtendPropertitesAttribute(utils.PathSearch("extend_propertites", respBody, nil)),
	}

	return []interface{}{rstBody}
}

func flattenAssetPropertiesAttribute(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"hwc_ecs":    flattenPropertiesHwcEcsAttribute(utils.PathSearch("hwc_ecs", respBody, nil)),
		"hwc_eip":    flattenPropertiesHwcEipAttribute(utils.PathSearch("hwc_eip", respBody, nil)),
		"hwc_vpc":    flattenPropertiesHwcVpcAttribute(utils.PathSearch("hwc_vpc", respBody, nil)),
		"hwc_subnet": flattenPropertiesHwcSubnetAttribute(utils.PathSearch("hwc_subnet", respBody, nil)),
		"hwc_rds":    flattenPropertiesHwcRdsAttribute(utils.PathSearch("hwc_rds", respBody, nil)),
		"hwc_domain": flattenPropertiesHwcDomainAttribute(utils.PathSearch("hwc_domain", respBody, nil)),
		"website":    flattenPropertiesWebsiteAttribute(utils.PathSearch("website", respBody, nil)),
		"oca_ip":     flattenPropertiesOcaIpAttribute(utils.PathSearch("oca_ip", respBody, nil)),
	}

	return []interface{}{rstBody}
}

// There is a problem with the API response for this field, with a horizontal bar between int and string types.
func flattenLevelAttribute(respValue interface{}) int {
	if respValue == nil {
		return 0
	}

	if intValue, ok := respValue.(float64); ok {
		return int(intValue)
	}

	stringValue, ok := respValue.(string)
	if !ok {
		return 0
	}

	r, err := strconv.Atoi(stringValue)
	if err != nil {
		log.Printf("[ERROR] convert the string %q to int failed.", stringValue)
	}

	return r
}

func flattenAssetDataObjectAttribute(respBody interface{}) []interface{} {
	rawMap := utils.PathSearch("data.data_object", respBody, nil)
	if rawMap == nil {
		return nil
	}

	rstBody := map[string]interface{}{
		"id":                 utils.PathSearch("id", rawMap, nil),
		"name":               utils.PathSearch("name", rawMap, nil),
		"provider":           utils.PathSearch("provider", rawMap, nil),
		"type":               utils.PathSearch("type", rawMap, nil),
		"checksum":           utils.PathSearch("checksum", rawMap, nil),
		"created":            utils.PathSearch("created", rawMap, nil),
		"provisioning_state": utils.PathSearch("provisioning_state", rawMap, nil),
		"environment":        flattenAssetEnvironmentAttribute(utils.PathSearch("environment", rawMap, nil)),
		"department":         flattenAssetDepartmentAttribute(utils.PathSearch("department", rawMap, nil)),
		"governance_user":    flattenAssetGovernanceUserAttribute(utils.PathSearch("governance_user", rawMap, nil)),
		"level":              flattenLevelAttribute(utils.PathSearch("level", rawMap, nil)),
		"properties":         flattenAssetPropertiesAttribute(utils.PathSearch("properties", rawMap, nil)),
	}

	return []interface{}{rstBody}
}

func resourceAssetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	respBody, err := ReadAssetDetail(client, d.Get("workspace_id").(string), d.Get("asset_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected500ErrInto404Err(err, "code", "SecMaster.00041004"),
			"error reading SecMaster asset")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("data_object", flattenAssetDataObjectAttribute(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAssetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	if err := updateAsset(client, d); err != nil {
		return diag.Errorf("error updating SecMaster asset in update operation: %s", err)
	}

	return resourceAssetRead(ctx, d, meta)
}

func resourceAssetDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/sa/resources"
		product     = "secmaster"
		workspaceID = d.Get("workspace_id").(string)
		assetID     = d.Get("asset_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceID)
	requestPath = strings.ReplaceAll(requestPath, "{id}", assetID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"batch_ids": []string{assetID},
		},
	}

	resp, err := client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected500ErrInto404Err(err, "code", "SecMaster.00041003"),
			"error deleting SecMaster asset")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	successDeleteID := utils.PathSearch("data.success_ids|[0]", respBody, "").(string)
	if successDeleteID != assetID {
		return diag.Errorf("error deleting SecMaster asset: asset ID %s does not match", assetID)
	}

	return nil
}

func resourceAssetImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	importIdParts := strings.Split(importId, "/")
	if len(importIdParts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, want <workspace_id>/<asset_id>, but got %s", importId)
	}

	mErr := multierror.Append(
		d.Set("workspace_id", importIdParts[0]),
		d.Set("asset_id", importIdParts[1]),
	)

	d.SetId(importIdParts[1])
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
