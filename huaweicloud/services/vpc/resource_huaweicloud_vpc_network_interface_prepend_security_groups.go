package vpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	networkInterfaceTypeNormal       = "normal"
	networkInterfaceTypeSupplemental = "supplemental"
)

// @API VPC GET /v1/{project_id}/ports/{port_id}
// @API VPC PUT /v1/{project_id}/ports/{port_id}
// @API VPC GET /v3/{project_id}/vpc/sub-network-interfaces/{sub_network_interface_id}
// @API VPC PUT /v3/{project_id}/vpc/sub-network-interfaces/{sub_network_interface_id}
func ResourceNetworkInterfacePrependSecurityGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkInterfacePrependSecurityGroupsCreate,
		ReadContext:   resourceNetworkInterfacePrependSecurityGroupsRead,
		UpdateContext: resourceNetworkInterfacePrependSecurityGroupsUpdate,
		DeleteContext: resourceNetworkInterfacePrependSecurityGroupsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNetworkInterfacePrependSecurityGroupsImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_interface_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      networkInterfaceTypeNormal,
				ValidateFunc: validation.StringInSlice([]string{networkInterfaceTypeNormal, networkInterfaceTypeSupplemental}, false),
			},
			"prepend_security_group_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"original_security_group_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"effective_security_group_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNetworkInterfacePrependSecurityGroupsImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.SplitN(d.Id(), "/", 2)
	if len(idParts) != 2 {
		return nil, fmt.Errorf("invalid import ID format (%s), want <network_interface_type>/<network_interface_id>", d.Id())
	}

	interfaceType := strings.TrimSpace(idParts[0])
	interfaceID := strings.TrimSpace(idParts[1])
	if interfaceType == "" || interfaceID == "" {
		return nil, fmt.Errorf("invalid import ID format (%s), want <network_interface_type>/<network_interface_id>", d.Id())
	}

	if err := d.Set("network_interface_type", interfaceType); err != nil {
		return nil, err
	}
	if err := d.Set("network_interface_id", interfaceID); err != nil {
		return nil, err
	}
	d.SetId(fmt.Sprintf("%s/%s", interfaceType, interfaceID))
	return []*schema.ResourceData{d}, nil
}

func resourceNetworkInterfacePrependSecurityGroupsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	originSecurityGroups, err := getCurrentInterfaceSecurityGroups(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("original_security_group_ids", originSecurityGroups); err != nil {
		return diag.FromErr(err)
	}

	finalSecurityGroups := reorderPrependSecurityGroups(
		utils.ExpandToStringList(d.Get("prepend_security_group_ids").([]interface{})),
		originSecurityGroups,
	)
	if err := updateInterfaceSecurityGroups(d, meta, finalSecurityGroups); err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%s/%s", d.Get("network_interface_type").(string), d.Get("network_interface_id").(string))
	d.SetId(id)

	return resourceNetworkInterfacePrependSecurityGroupsRead(ctx, d, meta)
}

func resourceNetworkInterfacePrependSecurityGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	currentSecurityGroups, err := getCurrentInterfaceSecurityGroups(d, meta)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "network interface security group prepend")
	}

	mErr := multierror.Append(nil,
		d.Set("region", meta.(*config.Config).GetRegion(d)),
		d.Set("effective_security_group_ids", currentSecurityGroups),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNetworkInterfacePrependSecurityGroupsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	currentSecurityGroups, err := getCurrentInterfaceSecurityGroups(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	finalSecurityGroups := reorderPrependSecurityGroups(
		utils.ExpandToStringList(d.Get("prepend_security_group_ids").([]interface{})),
		currentSecurityGroups,
	)
	if err := updateInterfaceSecurityGroups(d, meta, finalSecurityGroups); err != nil {
		return diag.FromErr(err)
	}

	return resourceNetworkInterfacePrependSecurityGroupsRead(ctx, d, meta)
}

func resourceNetworkInterfacePrependSecurityGroupsDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	originSecurityGroups := utils.ExpandToStringList(d.Get("original_security_group_ids").([]interface{}))
	if len(originSecurityGroups) == 0 {
		return nil
	}

	if err := updateInterfaceSecurityGroups(d, meta, originSecurityGroups); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getCurrentInterfaceSecurityGroups(d *schema.ResourceData, meta interface{}) ([]string, error) {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	networkInterfaceID := d.Get("network_interface_id").(string)

	if d.Get("network_interface_type").(string) == networkInterfaceTypeSupplemental {
		client, err := cfg.NewServiceClient("vpcv3", region)
		if err != nil {
			return nil, fmt.Errorf("error creating VPC v3 client: %s", err)
		}

		resp, err := getSubNetworkInterfaceInfo(client, networkInterfaceID)
		if err != nil {
			return nil, err
		}

		body, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		return utils.ExpandToStringList(utils.PathSearch("sub_network_interface.security_groups", body, []interface{}{}).([]interface{})), nil
	}

	client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC network client: %s", err)
	}

	port, err := ports.Get(client, networkInterfaceID)
	if err != nil {
		return nil, err
	}
	return port.SecurityGroups, nil
}

func updateInterfaceSecurityGroups(d *schema.ResourceData, meta interface{}, securityGroups []string) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	networkInterfaceID := d.Get("network_interface_id").(string)

	if d.Get("network_interface_type").(string) == networkInterfaceTypeSupplemental {
		client, err := cfg.NewServiceClient("vpcv3", region)
		if err != nil {
			return fmt.Errorf("error creating VPC v3 client: %s", err)
		}

		updateURL := client.ResourceBaseURL() + "vpc/sub-network-interfaces/" + networkInterfaceID
		opts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes:          []int{200},
			JSONBody: map[string]interface{}{
				"sub_network_interface": map[string]interface{}{
					"security_groups": securityGroups,
				},
			},
		}

		_, err = client.Request("PUT", updateURL, &opts)
		return err
	}

	client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return fmt.Errorf("error creating VPC network client: %s", err)
	}

	_, err = ports.Update(client, networkInterfaceID, ports.UpdateOpts{
		SecurityGroups: securityGroups,
	})
	return err
}

func reorderPrependSecurityGroups(prependSecurityGroups, currentSecurityGroups []string) []string {
	if len(prependSecurityGroups) == 0 && len(currentSecurityGroups) == 0 {
		return nil
	}

	result := make([]string, 0, len(prependSecurityGroups)+len(currentSecurityGroups))
	seen := make(map[string]struct{}, len(prependSecurityGroups)+len(currentSecurityGroups))

	for _, securityGroupID := range append(prependSecurityGroups, currentSecurityGroups...) {
		securityGroupID = strings.TrimSpace(securityGroupID)
		if securityGroupID == "" {
			continue
		}
		if _, ok := seen[securityGroupID]; ok {
			continue
		}

		seen[securityGroupID] = struct{}{}
		result = append(result, securityGroupID)
	}

	return result
}
