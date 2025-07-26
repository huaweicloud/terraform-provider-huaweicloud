package vpcep

import (
	"context"
	"encoding/json"
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
	"github.com/chnsz/golangsdk/openstack/vpcep/v1/endpoints"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPCEP POST /v1/{project_id}/vpc-endpoints
// @API VPCEP GET /v1/{project_id}/vpc-endpoints/{vpc_endpoint_id}
// @API VPCEP PUT /v1/{project_id}/vpc-endpoints/{vpc_endpoint_id}
// @API VPCEP DELETE /v1/{project_id}/vpc-endpoints/{vpc_endpoint_id}
// @API VPCEP PUT /v1/{project_id}/vpc-endpoints/{vpc_endpoint_id}/policy
// @API VPCEP PUT /v1/{project_id}/vpc-endpoints/{vpc_endpoint_id}/routetables
// @API VPCEP POST /v1/{project_id}/{resource_type}/{resource_id}/tags/action
func ResourceVPCEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCEndpointCreate,
		ReadContext:   resourceVPCEndpointRead,
		UpdateContext: resourceVPCEndpointUpdate,
		DeleteContext: resourceVPCEndpointDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_dns": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"routetables": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_whitelist": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"policy_statement": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			// Deprecated
			// The field type provided in the API document is different from the actual returned type.
			// As a result, an error is reported when the resource is imported.
			"policy_document": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`Specifies the endpoint policy information`, utils.SchemaDescInput{Deprecated: true},
				),
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ipv6_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": common.TagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"packet_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"private_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildPolicyStatement(d *schema.ResourceData) ([]endpoints.PolicyStatement, error) {
	if d.Get("policy_statement").(string) == "" {
		return nil, nil
	}

	var statements []endpoints.PolicyStatement
	err := json.Unmarshal([]byte(d.Get("policy_statement").(string)), &statements)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling policy, please check the format of the policy statement: %s", err)
	}
	return statements, nil
}

func resourceVPCEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcepClient, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	policyStatementOpts, err := buildPolicyStatement(d)
	if err != nil {
		return diag.FromErr(err)
	}

	enableACL := d.Get("enable_whitelist").(bool)
	createOpts := endpoints.CreateOpts{
		ServiceID:       d.Get("service_id").(string),
		VpcID:           d.Get("vpc_id").(string),
		SubnetID:        d.Get("network_id").(string),
		PortIP:          d.Get("ip_address").(string),
		Description:     d.Get("description").(string),
		IPVersion:       d.Get("ip_version").(string),
		IPv6Address:     d.Get("ipv6_address").(string),
		EnableDNS:       utils.Bool(d.Get("enable_dns").(bool)),
		EnableWhitelist: utils.Bool(enableACL),
		Tags:            utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		PolicyStatement: policyStatementOpts,
	}

	routeTables := d.Get("routetables").(*schema.Set)
	if routeTables.Len() > 0 {
		createOpts.RouteTables = utils.ExpandToStringList(routeTables.List())
	}

	raw := d.Get("whitelist").(*schema.Set).List()
	if enableACL && len(raw) > 0 {
		createOpts.Whitelist = utils.ExpandToStringList(raw)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	ep, err := endpoints.Create(vpcepClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating VPC endpoint: %s", err)
	}

	d.SetId(ep.ID)
	log.Printf("[INFO] Waiting for VPC endpoint(%s) to become accepted", ep.ID)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"creating"},
		Target:       []string{"accepted", "pendingAcceptance"},
		Refresh:      waitForVPCEndpointStatus(vpcepClient, ep.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf("error waiting for VPC endpoint(%s) to become accepted: %s", ep.ID, stateErr)
	}

	return resourceVPCEndpointRead(ctx, d, meta)
}

func resourceVPCEndpointRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcepClient, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	ep, err := endpoints.Get(vpcepClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC endpoint")
	}

	log.Printf("[DEBUG] retrieving VPC endpoint: %#v", ep)

	policyStatements, err := json.Marshal(ep.PolicyStatement)
	if err != nil {
		return diag.Errorf("error marshaling policy statement: %s", err)
	}

	serviceType := ep.ServiceType
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("status", ep.Status),
		d.Set("service_id", ep.ServiceID),
		d.Set("service_name", ep.ServiceName),
		d.Set("service_type", serviceType),
		d.Set("vpc_id", ep.VpcID),
		d.Set("network_id", ep.SubnetID),
		d.Set("ip_address", ep.IPAddr),
		d.Set("description", ep.Description),
		d.Set("enable_whitelist", ep.EnableWhitelist),
		d.Set("packet_id", ep.MarkerID),
		d.Set("tags", utils.TagsToMap(ep.Tags)),
		d.Set("policy_statement", string(policyStatements)),
		d.Set("ip_version", ep.IpVersion),
		d.Set("ipv6_address", ep.Ipv6Address),
	)

	// if the VPC endpoint type is interface, the field is used and need to be set
	if serviceType == "interface" {
		mErr = multierror.Append(mErr, d.Set("enable_dns", ep.EnableDNS))
	}

	// if the VPC endpoint type is interface, the parameter can be ignored
	// the api will return an empty array
	if len(ep.RouteTables) == 0 {
		mErr = multierror.Append(mErr, d.Set("routetables", nil))
	} else {
		mErr = multierror.Append(mErr, d.Set("routetables", ep.RouteTables))
	}

	if len(ep.Whitelist) == 0 {
		// if the "whitelist" is not specified, the api will return an empty array
		mErr = multierror.Append(mErr, d.Set("whitelist", nil))
	} else {
		mErr = multierror.Append(mErr, d.Set("whitelist", ep.Whitelist))
	}

	if len(ep.DNSNames) > 0 {
		mErr = multierror.Append(mErr, d.Set("private_domain_name", ep.DNSNames[0]))
	} else {
		mErr = multierror.Append(mErr, d.Set("private_domain_name", nil))
	}
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVPCEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcepClient, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}
	if d.HasChanges("enable_whitelist", "whitelist") {
		updateOpts := endpoints.UpdateOpts{
			EnableWhitelist: utils.Bool(d.Get("enable_whitelist").(bool)),
			Whitelist:       utils.ExpandToStringList(d.Get("whitelist").(*schema.Set).List()),
		}
		_, err := endpoints.Update(vpcepClient, updateOpts, d.Id()).Extract()
		if err != nil {
			return diag.Errorf("error updating VPC endpoint whitelist: %s", err)
		}
	}
	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(vpcepClient, d, tagVPCEP, d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of VPC endpoint %s: %s", d.Id(), tagErr)
		}
	}
	if d.HasChanges("policy_statement") {
		policyStatementOpts, err := buildPolicyStatement(d)
		if err != nil {
			return diag.FromErr(err)
		}

		updatePolicyOpts := endpoints.UpdatePolicyOpts{
			PolicyStatement: policyStatementOpts,
		}
		_, err = endpoints.UpdatePolicy(vpcepClient, updatePolicyOpts, d.Id()).Extract()
		if err != nil {
			return diag.Errorf("error updating VPC endpoint policy: %s", err)
		}
	}

	// binding or unbinding routetables
	if d.HasChanges("routetables") {
		routeTables := utils.ExpandToStringListBySet(d.Get("routetables").(*schema.Set))
		err = updateRouteTables(d, vpcepClient, routeTables)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVPCEndpointRead(ctx, d, meta)
}

func updateRouteTables(d *schema.ResourceData, client *golangsdk.ServiceClient, routeTables []string) error {
	routeTablesHttpUrl := "v1/{project_id}/vpc-endpoints/{vpc_endpoint_id}/routetables"
	routeTablesPath := client.Endpoint + routeTablesHttpUrl
	routeTablesPath = strings.ReplaceAll(routeTablesPath, "{project_id}", client.ProjectID)
	routeTablesPath = strings.ReplaceAll(routeTablesPath, "{vpc_endpoint_id}", d.Id())

	routeTablesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	routeTablesOpt.JSONBody = utils.RemoveNil(map[string]interface{}{
		"routetables": routeTables,
	})

	_, err := client.Request("PUT", routeTablesPath, &routeTablesOpt)
	if err != nil {
		return fmt.Errorf("error updating routetables: %s", err)
	}

	return nil
}

func resourceVPCEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	vpcepClient, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	err = endpoints.Delete(vpcepClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting VPC endpoint %s: %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"deleting"},
		Target:       []string{"deleted"},
		Refresh:      waitForVPCEndpointStatus(vpcepClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting VPC endpoint %s: %s", d.Id(), err)
	}

	return nil
}

func waitForVPCEndpointStatus(vpcepClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ep, err := endpoints.Get(vpcepClient, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted VPC endpoint %s", id)
				return ep, "deleted", nil
			}
			return ep, "error", err
		}

		return ep, ep.Status, nil
	}
}
