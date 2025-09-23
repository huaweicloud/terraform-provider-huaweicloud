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

var namespaceNonUpdatableParams = []string{"name"}

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
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The annotations of the namespace.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the namespace.`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation timestamp of the namespace.`,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource version of the namespace.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the namespace.`,
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
			"finalizers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The finalizers of the namespace.`,
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

func resourceNamespaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createNamespaceHttpUrl := "apis/cci/v2/namespaces"
	createNamespacePath := client.Endpoint + createNamespaceHttpUrl
	createNamespaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createNamespaceOpt.JSONBody = buildCreateNamespaceParams(d)

	resp, err := client.Request("POST", createNamespacePath, &createNamespaceOpt)
	if err != nil {
		return diag.Errorf("error creating CCI namespace: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ns := utils.PathSearch("metadata.name", respBody, "").(string)
	if ns == "" {
		return diag.Errorf("unable to find namespace name from API response")
	}
	d.SetId(ns)

	err = waitForCreateNamespaceStatus(ctx, client, ns, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceNamespaceRead(ctx, d, meta)
}

func buildCreateNamespaceParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": d.Get("name"),
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	resp, err := GetNamespaceDetail(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting the specifies namespace form server")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("api_version", utils.PathSearch("apiVersion", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("annotations", utils.PathSearch("metadata.annotations", resp, nil)),
		d.Set("labels", utils.PathSearch("metadata.labels", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
		d.Set("finalizers", utils.PathSearch("spec.finalizers", resp, nil)),
		d.Set("status", utils.PathSearch("status.phase", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNamespaceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceNamespaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}
	namespace := d.Id()

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
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return "", "DELETED", nil
				}
				return nil, "ERROR", err
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
		return getNamespaceDetailResp, err
	}

	return utils.FlattenResponse(getNamespaceDetailResp)
}
