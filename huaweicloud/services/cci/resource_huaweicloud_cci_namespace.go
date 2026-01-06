package cci

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cci/v1/namespaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

type stateRefresh struct {
	Pending      []string
	Target       []string
	Delay        time.Duration
	Timeout      time.Duration
	PollInterval time.Duration
}

// @API CCI DELETE /api/v1/namespaces/{name}
// @API CCI GET /api/v1/namespaces/{name}
// @API CCI GET /api/v1/namespaces
// @API CCI POST /api/v1/namespaces
func ResourceCciNamespace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCciNamespaceCreate,
		ReadContext:   resourceCciNamespaceRead,
		DeleteContext: resourceCciNamespaceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCciNamespaceImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"general-computing", "gpu-accelerated",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auto_expend_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"warmup_pool_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"recycling_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"warmup_pool_size"},
			},
			"container_network_enabled": {
				Type:         schema.TypeBool,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"warmup_pool_size"},
			},
			"rbac_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"created_at": {
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

func buildCciNamespaceCreateParams(d *schema.ResourceData, conf *config.Config) (namespaces.CreateOpts,
	error) {
	createOpts := namespaces.CreateOpts{
		Kind:       "Namespace",
		ApiVersion: "v1",
		Metadata: namespaces.Metadata{
			Name: d.Get("name").(string),
			Annotations: namespaces.Annotations{
				Flavor:     d.Get("type").(string),
				AutoExpend: d.Get("auto_expend_enabled").(bool),
			},
			Labels: &namespaces.Labels{
				EnterpriseProjectID: conf.GetEnterpriseProjectID(d),
				RbacEnable:          d.Get("rbac_enabled").(bool),
			},
		},
	}

	if size, isAdvance := d.GetOk("warmup_pool_size"); isAdvance {
		createOpts.Metadata.Annotations.PoolSize = size.(int)
		createOpts.Metadata.Annotations.RecyclingInterval = d.Get("recycling_interval").(int)
		if enabled, ok := d.GetOk("container_network_enabled"); ok && enabled.(bool) {
			createOpts.Metadata.Annotations.NetworkEnable = "vpc-network-ready"
		}
	}

	return createOpts, nil
}

func resourceCciNamespaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.CciV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CCI v1 client: %s", err)
	}

	createOpts, err := buildCciNamespaceCreateParams(d, conf)
	if err != nil {
		return diag.Errorf("Unable to build createOpts of the CCI namespace: %s", err)
	}
	ns := d.Get("name").(string)
	namespace, err := namespaces.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating CCI namespace: %s", err)
	}
	d.SetId(namespace.Metadata.UID)
	stateRef := stateRefresh{
		Pending:      []string{"Pending"},
		Target:       []string{"Active"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        6 * time.Second,
		PollInterval: 5 * time.Second,
	}
	if err := waitForCciNamespacestateRefresh(ctx, client, ns, stateRef); err != nil {
		return err
	}

	return resourceCciNamespaceRead(ctx, d, meta)
}

func isContainNetworkEnabled(network string) bool {
	return network == "vpc-network-ready"
}

func setCciNamespaceParams(d *schema.ResourceData, resp *namespaces.Namespace) error {
	metadata := &resp.Metadata

	mErr := multierror.Append(nil,
		d.Set("name", metadata.Name),
		d.Set("type", metadata.Annotations.Flavor),
		d.Set("enterprise_project_id", metadata.Labels.EnterpriseProjectID),
		d.Set("rbac_enabled", metadata.Labels.RbacEnable),
		d.Set("auto_expend_enabled", metadata.Annotations.AutoExpend),
		d.Set("created_at", metadata.CreationTimestamp),
		d.Set("status", &resp.Status.Phase),
		d.Set("warmup_pool_size", metadata.Annotations.PoolSize),
		d.Set("recycling_interval", metadata.Annotations.RecyclingInterval),
		d.Set("container_network_enabled", isContainNetworkEnabled(metadata.Annotations.NetworkEnable)),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func resourceCciNamespaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CciV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CCI v1 client: %s", err)
	}

	var response *namespaces.Namespace
	response, err = GetCciNamespaceInfoById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error getting the specifies namespace form server")
	}
	if response != nil {
		mErr := multierror.Append(nil,
			d.Set("region", region),
			setCciNamespaceParams(d, response),
		)
		if mErr.ErrorOrNil() != nil {
			return diag.Errorf("Error saving the specifies namespace (%s) to state: %s", d.Id(), mErr)
		}
	}

	return nil
}

func resourceCciNamespaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.CciV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CCI v1 client: %s", err)
	}

	ns := d.Get("name").(string)
	_, err = namespaces.Delete(client, ns).Extract()
	if err != nil {
		return diag.Errorf("Error deleting the specifies namespace (%s): %s", d.Id(), err)
	}

	stateRef := stateRefresh{
		Pending:      []string{"Active", "Terminating"},
		Target:       []string{"DELETED"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        6 * time.Second,
		PollInterval: 5 * time.Second,
	}
	if err := waitForCciNamespacestateRefresh(ctx, client, ns, stateRef); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func waitForCciNamespacestateRefresh(ctx context.Context, c *golangsdk.ServiceClient, ns string,
	s stateRefresh) diag.Diagnostics {
	stateConf := &resource.StateChangeConf{
		Pending:      s.Pending,
		Target:       s.Target,
		Refresh:      namespacestateRefreshFunc(c, ns),
		Timeout:      s.Timeout,
		Delay:        s.Delay,
		PollInterval: s.PollInterval,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Waiting for the status of the namespace (%s) to complete (%s) timeout: %s",
			ns, s.Target, err)
	}
	return nil
}

func namespacestateRefreshFunc(c *golangsdk.ServiceClient, ns string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		response, err := getCciNamespaceInfoByName(c, ns)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return response, "DELETED", nil
			}
			return response, "ERROR", nil
		}
		if response != nil {
			return response, response.Status.Phase, nil
		}
		return response, "ERROR", nil
	}
}

func getCciNamespaceInfoByName(c *golangsdk.ServiceClient, ns string) (*namespaces.Namespace, error) {
	namespace, err := namespaces.Get(c, ns).Extract()
	return namespace, err
}

// GetCciNamespaceInfoById is a method to get namespace informations by client and namespace ID.
func GetCciNamespaceInfoById(c *golangsdk.ServiceClient, id string) (*namespaces.Namespace, error) {
	var response *namespaces.Namespace
	pages, err := namespaces.List(c, namespaces.ListOpts{}).AllPages()
	if err != nil {
		return response, err
	}
	responses, err := namespaces.ExtractNamespaces(pages)
	if err != nil {
		return response, err
	}
	for _, v := range responses {
		if v.Metadata.UID == id {
			return &v, nil
		}
	}
	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/api/v1/namespaces",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the namespace (%s) does not exist", id)),
		},
	}
}

func resourceCciNamespaceImportState(_ context.Context, d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	conf := meta.(*config.Config)
	client, err := conf.CciV1Client(conf.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating CCI v1 client: %s", err)
	}

	response, err := getCciNamespaceInfoByName(client, d.Id()) // The namespace is imported by name.
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("unable to find the CCI namespace by name (%s)", d.Id())
	}
	d.SetId(response.Metadata.UID)

	return []*schema.ResourceData{d}, nil
}
