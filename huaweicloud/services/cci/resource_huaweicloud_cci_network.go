package cci

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cci/v1/networks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpc"
)

// @API CCI DELETE /apis/networking.cci.io/v1beta1/namespaces/{ns}/networks/{name}
// @API CCI GET /apis/networking.cci.io/v1beta1/namespaces/{ns}/networks/{name}
// @API CCI POST /apis/networking.cci.io/v1beta1/namespaces/{ns}/networks
// @API VPC GET /v1/{project_id}/subnets/{subnet_id}
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
				Optional: true,
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

func resourceNetworkAnnotations(d *schema.ResourceData, conf *config.Config) map[string]string {
	result := map[string]string{
		"network.alpha.kubernetes.io/domain_id":  conf.DomainID,
		"network.alpha.kubernetes.io/project_id": conf.HwClient.ProjectID,
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		result["network.alpha.kubernetes.io/default-security-group"] = v.(string)
	}
	return result
}

func resourceCciNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cciClient, err := conf.CciV1BetaClient(region)
	if err != nil {
		return diag.Errorf("Error creating CCI Beta v1 client: %s", err)
	}

	networkId := d.Get("network_id").(string)
	subnet, err := vpc.GetVpcSubnetById(conf, region, networkId)
	if err != nil {
		return diag.Errorf("The subnet does not exist: %s", err)
	}

	opt := networks.CreateOpts{
		Kind:       "Network",
		ApiVersion: "networking.cci.io/v1beta1",
		Metadata: networks.CreateMetaData{
			Name:        d.Get("name").(string),
			Annotations: resourceNetworkAnnotations(d, conf),
		},
		Spec: networks.Spec{
			NetworkType: "underlay_neutron",
			AttachedVPC: subnet.VPC_ID,
			NetworkID:   networkId,
		},
	}

	if az, ok := d.GetOk("availability_zone"); ok {
		opt.Spec.AvailableZone = az.(string)
	}

	ns := d.Get("namespace").(string)
	create, err := networks.Create(cciClient, ns, opt).Extract()

	if err != nil {
		return diag.Errorf("Error creating CCI Network: %s", err)
	}

	d.SetId(create.Metadata.Name)

	log.Printf("[DEBUG] Waiting for CCI network (%s) to become available", d.Id())
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
		return diag.Errorf("Error obtain CCI network status: %s", err)
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
		return diag.Errorf("Error setting CCI network parameters: %s", mErr)
	}
	return nil
}

func resourceCciNetworkRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	cciClient, err := conf.CciV1BetaClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CCI Beta v1 client: %s", err)
	}

	ns := d.Get("namespace").(string)
	network, err := networks.Get(cciClient, ns, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CCI network")
	}

	return setCciNetworkParms(d, network)
}

func resourceCciNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	cciClient, err := conf.CciV1BetaClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CCI Beta v1 client: %s", err)
	}

	ns := d.Get("namespace").(string)
	err = networks.Delete(cciClient, ns, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting CCI Network: %s", err)
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
		return diag.Errorf("Error obtain CCI network status: %s", err)
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
		log.Printf("[DEBUG] Attempting to delete CCI network %s.", name)

		r, err := networks.Get(cciClient, ns, name).Extract()
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[DEBUG] Successfully deleted CCI network %s", name)
			return r, "Deleted", nil
		}
		if r.Status.State == "Terminating" {
			return r, "Terminating", nil
		}
		log.Printf("[DEBUG] CCI network %s still available.", name)
		return r, "Active", nil
	}
}

func resourceCciNetworkImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <namespace>/<id>")
	}

	d.SetId(parts[1])
	d.Set("namespace", parts[0])

	return []*schema.ResourceData{d}, nil
}
