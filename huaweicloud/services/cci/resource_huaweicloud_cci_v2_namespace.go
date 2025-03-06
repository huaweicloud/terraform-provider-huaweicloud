package cci

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var namespaceNonUpdatableParams = []string{
	"name",
	"flavor",
	"enterprise_project_id",
	"warm_pool_size",
	"warm_pool_recycle_interval",
	"container_network_enabled",
	"rbac_enabled",
}

// @API CCI DELETE /apis/cci/v2/namespaces/{name}
// @API CCI GET /apis/cci/v2/namespaces/{name}
// @API CCI POST /apis/cci/v2/namespaces
func ResourceNamespace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNamespaceCreate,
		UpdateContext: resourceNamespaceUpdate,
		ReadContext:   resourceNamespaceRead,
		DeleteContext: resourceNamespaceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(namespaceNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"general-computing", "gpu-accelerated",
				}, false),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"warm_pool_size": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"warm_pool_recycle_interval": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"warm_pool_size"},
			},
			"container_network_enabled": {
				Type:         schema.TypeBool,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"warmup_pool_size"},
			},
			"rbac_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceNamespaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.CciV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CCI v2 client: %s", err)
	}

	createNamespaceHttpUrl := "apis/cci/v2/namespaces"
	createNamespacePath := client.Endpoint + createNamespaceHttpUrl
	createNamespaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createNamespaceOpt.JSONBody = utils.RemoveNil(buildCreateNamespaceParams(d, conf))

	resp, err := client.Request("POST", createNamespacePath, &createNamespaceOpt)
	if err != nil {
		return diag.Errorf("error creating CCI namespace: %s", err)
	}

	ns := utils.PathSearch("metadata.name", resp, "").(string)
	if ns == "" {
		return diag.Errorf("unable to find namespace name from API response")
	}
	d.SetId(ns)

	err = waitForCreateNamespaceStatus(ctx, client, ns, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceCciNamespaceRead(ctx, d, meta)
}

func buildCreateNamespaceParams(d *schema.ResourceData, conf *config.Config) map[string]interface{} {
	annotations := map[string]interface{}{
		"namespace.kubernetes.io/flavor":            d.Get("flavor").(string),
		"network.cci.io/warm-pool-size":             d.Get("warm_pool_size"),
		"network.cci.io/warm-pool-recycle-interval": d.Get("warm_pool_recycle_interval"), // unit: h
	}
	if d.Get("container_network_enabled").(bool) {
		annotations["network.cci.io/ready-before-pod-run"] = "vpc-network-ready"
	}

	bodyParams := map[string]interface{}{
		"kind":       "Namespace",
		"apiVersion": "v2",
		"metadata": map[string]interface{}{
			"name":        d.Get("name"),
			"annotations": annotations,
			"labels": map[string]interface{}{
				"rbac.authorization.cci.io/enable-k8s-rbac": d.Get("rbac_enabled"),
				"sys_enterprise_project_id":                 conf.GetEnterpriseProjectID(d),
			},
			// "clusterName": "",
			// "creationTimestamp": "",
			// "deletionGracePeriodSeconds": 0,
			// "deletionTimestamp": "",
			// "enable": false,
			// "finaliaers": []string{},
			// "generateName": "",
			// "generation": 0,
			// "managedFields": []map[string]interface{}{},
			// "namespace": "",
			// "ownerReferences": []map[string]interface{}{},
			// "resourceVersion": "",
			// "selfLink": "",
			// "uid": "",
		},
	}

	return bodyParams
}

func waitForCreateNamespaceStatus(ctx context.Context, client *golangsdk.ServiceClient, ns string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Active"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetNamespaceDetail(client, ns)
			if err != nil {
				return nil, "failed", err
			}
			return resp, utils.PathSearch("status.phase", resp, "").(string), nil
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the status of the namespace to complete active timeout: %s", err)
	}
	return nil
}

func resourceNamespaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CciV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CCI v1 client: %s", err)
	}

	resp, err := GetNamespaceDetail(client, d.Get("name").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting the specifies namespace form server")
	}

	annotations := utils.PathSearch("metadata.annotations", resp, nil)
	containNetworkEnabled := isContainNetworkEnabled(utils.PathSearch("network.cci.io/ready-before-pod-run", annotations, "").(string))
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("flavor", utils.PathSearch("namespace.kubernetes.io/flavor", annotations, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("metadata.labels.sys_enterprise_project_id", resp, nil)),
		d.Set("warm_pool_size", utils.PathSearch("network.cci.io/warm-pool-size", annotations, nil)),
		d.Set("warm_pool_recycle_interval", utils.PathSearch("network.cci.io/warm-pool-recycle-interval", annotations, nil)),
		d.Set("container_network_enabled", containNetworkEnabled),
		d.Set("rbac_enabled", utils.PathSearch("metadata.labels.rbac.authorization.cci.io/enable-k8s-rbac", resp, nil)),
		// d.Set("kind", utils.PathSearch("kind", resp, nil)),
		// d.Set("api_version", utils.PathSearch("api_version", resp, nil)),
		d.Set("created_at", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("status", utils.PathSearch("status.phase", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNamespaceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceNamespaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.CciV2Client(conf.GetRegion(d))
	namespace := d.Id()
	if err != nil {
		return diag.Errorf("Error creating CCI v2 client: %s", err)
	}

	deleteNamespaceHttpUrl := "apis/cci/v2/namespaces/{name}"
	deleteNamespacePath := client.Endpoint + deleteNamespaceHttpUrl
	deleteNamespacePath = strings.ReplaceAll(deleteNamespacePath, "{name}", namespace)
	deleteNamespaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	deleteNamespaceOpt.JSONBody = utils.RemoveNil(buildCreateNamespaceParams(d, conf))

	_, err = client.Request("DELETE", deleteNamespacePath, &deleteNamespaceOpt)
	if err != nil {
		return diag.Errorf("error deleting the specifies namespace (%s): %s", namespace, err)
	}

	err = waitForDeleteNamespaceStatus(ctx, client, namespace, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForDeleteNamespaceStatus(ctx context.Context, client *golangsdk.ServiceClient, ns string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Active", "Terminating"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetNamespaceDetail(client, ns)
			if err != nil {
				return nil, "failed", err
			}
			return resp, utils.PathSearch("status.phase", resp, "").(string), nil
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the status of the namespace to complete delete timeout: %s", err)
	}
	return nil
}

func GetNamespaceDetail(client *golangsdk.ServiceClient, namespace string) (interface{}, error) {
	getNamespaceDetailHttpUrl := "apis/cci/v2/namespaces/{name}"
	getNamespaceDetailPath := client.Endpoint + getNamespaceDetailHttpUrl
	getNamespaceDetailPath = strings.ReplaceAll(getNamespaceDetailPath, "{name}", namespace)
	getNamespaceDetailOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getNamespaceDetailResp, err := client.Request("GET", getNamespaceDetailPath, &getNamespaceDetailOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying CCI namespace: %s", err)
	}

	return utils.FlattenResponse(getNamespaceDetailResp)
}
