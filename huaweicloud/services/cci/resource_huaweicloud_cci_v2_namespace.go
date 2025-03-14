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
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the namespace.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The annotations of the namespace.`,
				DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
					oldVal, newVal := d.GetChange("annotations")
					for key, value := range newVal.(map[string]interface{}) {
						if mapValue, exists := oldVal.(map[string]interface{})[key]; exists && mapValue == value {
							continue
						}
						return false
					}
					return true
				},
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the namespace.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the namespace.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kind of the namespace.`,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cluster name of the namespace.`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation timestamp of the namespace.`,
			},
			"deletion_grace_period_seconds": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The deletion grace period seconds of the namespace.`,
			},
			"deletion_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The deletion timestamp of the namespace.`,
			},
			"finalizers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The finalizers of the namespace.`,
			},
			"generate_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The generate name of the namespace.`,
			},
			"generation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The generation of the namespace.`,
			},
			"managed_fields": {
				Type:        schema.TypeList,
				Elem:        managedFieldsSchema(),
				Computed:    true,
				Description: `The managed fields of the namespace.`,
			},
			"owner_references": {
				Type:        schema.TypeList,
				Elem:        ownerReferencesSchema(),
				Computed:    true,
				Description: `The owner references of the namespace.`,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource version of the namespace.`,
			},
			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The self link of the namespace.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the namespace.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the namespace.`,
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

func managedFieldsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the managed fields.`,
			},
			"fields_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fields type of the managed fields.`,
			},
			"fields_v1": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fields v1 of the managed fields.`,
			},
			"manager": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The manager of the managed fields.`,
			},
			"operation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The operation of the managed fields.`,
			},
			"time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time of the managed fields.`,
			},
		},
	}
	return &sc
}

func ownerReferencesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the owner references.`,
			},
			"block_owner_deletion": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `The block owner deletion of the owner references.`,
			},
			"controller": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `The controller of the owner references.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kind of the owner references.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the owner references.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the owner references.`,
			},
		},
	}
	return &sc
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
	createNamespaceOpt.JSONBody = utils.RemoveNil(buildCreateNamespaceParams(d))

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

func buildCreateNamespaceParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"kind":       "Namespace",
		"apiVersion": "v2",
		"metadata": map[string]interface{}{
			"name":        d.Get("name"),
			"annotations": d.Get("annotations"),
			"labels":      d.Get("labels"),
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

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("api_version", utils.PathSearch("api_version", resp, nil)),
		d.Set("annotations", utils.PathSearch("metadata.annotations", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("cluster_name", utils.PathSearch("metadata.clusterName", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("deletion_grace_period_seconds", utils.PathSearch("metadata.deletionGracePeriodSeconds", resp, nil)),
		d.Set("deletion_timestamp", utils.PathSearch("metadata.deletionTimestamp", resp, nil)),
		d.Set("finaliaers", utils.PathSearch("metadata.finaliaers", resp, nil)),
		d.Set("generate_name", utils.PathSearch("metadata.generateName", resp, nil)),
		d.Set("generation", utils.PathSearch("metadata.generation", resp, nil)),
		d.Set("managed_fields", flattenManagedFields(utils.PathSearch("metadata.managedFields", resp, nil))),
		d.Set("namespace", utils.PathSearch("metadata.namespace", resp, nil)),
		d.Set("owner_references", flattenOwnerReferences(utils.PathSearch("metadata.ownerReferences", resp, nil))),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("self_link", utils.PathSearch("metadata.selfLink", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOwnerReferences(resp interface{}) []interface{} {
	rst := make([]interface{}, 0)
	if resp == nil {
		return nil
	}

	rst = append(rst, map[string]interface{}{
		"api_version":          utils.PathSearch("apiVersion", resp, nil),
		"block_owner_deletion": utils.PathSearch("blockOwnerDeletion", resp, false),
		"controller":           utils.PathSearch("controller", resp, false),
		"kind":                 utils.PathSearch("kind", resp, nil),
		"name":                 utils.PathSearch("name", resp, nil),
		"uid":                  utils.PathSearch("uid", resp, nil),
	})

	return rst
}

func flattenManagedFields(resp interface{}) []interface{} {
	rst := make([]interface{}, 0)
	if resp == nil {
		return nil
	}

	rst = append(rst, map[string]interface{}{
		"api_version": utils.PathSearch("apiVersion", resp, nil),
		"fields_type": utils.PathSearch("fieldsType", resp, nil),
		"fields_v1":   utils.PathSearch("fieldsV1", resp, nil),
		"manager":     utils.PathSearch("manager", resp, nil),
		"operation":   utils.PathSearch("operation", resp, nil),
		"time":        utils.PathSearch("time", resp, nil),
	})

	return rst
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
