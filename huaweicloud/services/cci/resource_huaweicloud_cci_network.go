package cci

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cci/v1/networks"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceCciNetworkV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCciNetworkCreate,
		ReadContext:   resourceCciNetworkRead,
		DeleteContext: resourceCciNetworkDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCciNetworkImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(
						"[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*"),
						"The name can only contains lowercase characters, hyphens (-) and dots (.), and must start "+
							"and end with a character or digit."),
					validation.StringLenBetween(1, 200),
				),
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkAnnotations(d *schema.ResourceData, config *config.Config) map[string]string {
	result := map[string]string{
		"network.alpha.kubernetes.io/domain_id":  config.DomainID,
		"network.alpha.kubernetes.io/project_id": config.HwClient.ProjectID,
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		result["network.alpha.kubernetes.io/default-security-group"] = v.(string)
	}
	return result
}

func resourceCciNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cciClient, err := config.CciV1BetaClient(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCI Beta v1 client: %s", err)
	}

	networkId := d.Get("network_id").(string)
	subnet, err := vpc.GetVpcSubnetById(config, region, networkId)
	if err != nil {
		return fmtp.DiagErrorf("The subnet does not exist: %s", err)
	}

	opt := networks.CreateOpts{
		Kind:       "Network",
		ApiVersion: "networking.cci.io/v1beta1",
		Metadata: networks.CreateMetaData{
			Name:        d.Get("name").(string),
			Annotations: resourceNetworkAnnotations(d, config),
		},
		Spec: networks.Spec{
			AvailableZone: d.Get("availability_zone").(string),
			NetworkType:   "underlay_neutron",
			AttachedVPC:   subnet.VPC_ID,
			NetworkID:     networkId,
		},
	}

	ns := d.Get("namespace").(string)
	create, err := networks.Create(cciClient, ns, opt).Extract()

	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCI Network: %s", err)
	}

	d.SetId(create.Metadata.Name)

	logp.Printf("[DEBUG] Waiting for HuaweiCloud CCI network (%s) to become available", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Initializing", "Pending"},
		Target:       []string{"Active"},
		Refresh:      waitForCciNetworkActive(cciClient, ns, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error obtain HuaweiCloud CCI network status: %s", err)
	}

	return resourceCciNetworkRead(ctx, d, meta)
}

func setCciNetworkParms(d *schema.ResourceData, network *networks.Network) diag.Diagnostics {
	mErr := multierror.Append(nil,
		d.Set("availability_zone", network.Spec.AvailableZone),
		d.Set("name", network.Metadata.Name),
		d.Set("network_id", network.Spec.NetworkID),
		d.Set("security_group_id", network.Metadata.Annotations["network.alpha.kubernetes.io/default-security-group"]),
		d.Set("vpc_id", network.Spec.AttachedVPC),
		d.Set("subnet_id", network.Spec.SubnetID),
		d.Set("cidr", network.Spec.Cidr),
		d.Set("status", network.Status.State),
	)
	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("Error setting CCI network parameters: %s", mErr)
	}
	return nil
}

func resourceCciNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	cciClient, err := config.CciV1BetaClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCI Beta v1 client: %s", err)
	}

	ns := d.Get("namespace").(string)
	network, err := networks.Get(cciClient, ns, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CCI network")
	}

	return setCciNetworkParms(d, network)
}

func resourceCciNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	cciClient, err := config.CciV1BetaClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCI Beta v1 client: %s", err)
	}

	ns := d.Get("namespace").(string)
	err = networks.Delete(cciClient, ns, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloud CCI Network: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Terminating", "Active"},
		Target:     []string{"Deleted"},
		Refresh:    waitForCciNetworkDelete(cciClient, ns, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error obtain HuaweiCloud CCI network status: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCciNetworkActive(cciClient *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := networks.Get(cciClient, ns, name).Extract()
		if err != nil {
			return nil, "", err
		}

		return n, n.Status.State, nil
	}
}

func waitForCciNetworkDelete(cciClient *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud CCI network %s.", name)

		r, err := networks.Get(cciClient, ns, name).Extract()
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			logp.Printf("[DEBUG] Successfully deleted HuaweiCloud CCI network %s", name)
			return r, "Deleted", nil
		}
		if r.Status.State == "Terminating" {
			return r, "Terminating", nil
		}
		logp.Printf("[DEBUG] HuaweiCloud CCI network %s still available.", name)
		return r, "Active", nil
	}
}

func resourceCciNetworkImportState(_ context.Context, d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmtp.Errorf("Invalid format specified for import id, must be <namespace>/<id>")
	}

	d.SetId(parts[1])
	d.Set("namespace", parts[0])

	return []*schema.ResourceData{d}, nil
}
