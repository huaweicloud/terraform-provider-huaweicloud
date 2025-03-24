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

var pvcNonUpdatableParams = []string{"name"}

// @API CCI POST /apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims
// @API CCI GET /apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims/{name}
// @API CCI PUT /apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims/{name}
// @API CCI DELETE /apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims/{name}
func ResourceV2PersistentVolumeClaim() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2PersistentVolumeClaimCreate,
		UpdateContext: resourceV2PersistentVolumeClaimUpdate,
		ReadContext:   resourceV2PersistentVolumeClaimRead,
		DeleteContext: resourceV2PersistentVolumeClaimDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(pvcNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the namespace.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the persistent volume claim in the namespace.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The annotations of the persistent volume claim.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the persistent volume claim.`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation timestamp of the persistent volume claim.`,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource version of the persistent volume claim.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the persistent volume claim.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the persistent volume claim.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kind of the persistent volume claim.`,
			},
			"finalizers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The finalizers of the persistent volume claim.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the persistent volume claim.`,
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

func resourceV2PersistentVolumeClaimCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createV2PersistentVolumeClaimHttpUrl := "apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims"
	createV2PersistentVolumeClaimPath := client.Endpoint + createV2PersistentVolumeClaimHttpUrl
	createV2PersistentVolumeClaimOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createV2PersistentVolumeClaimOpt.JSONBody = buildCreateV2PersistentVolumeClaimParams(d)

	resp, err := client.Request("POST", createV2PersistentVolumeClaimPath, &createV2PersistentVolumeClaimOpt)
	if err != nil {
		return diag.Errorf("error creating CCI namespace: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	name := utils.PathSearch("metadata.name", respBody, "").(string)
	ns := utils.PathSearch("metadata.namespace", respBody, "").(string)
	if ns == "" {
		return diag.Errorf("unable to find V2PersistentVolumeClaim name from API response")
	}
	d.SetId(ns)

	err = waitForCreateV2PersistentVolumeClaimStatus(ctx, client, ns, name, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceV2PersistentVolumeClaimRead(ctx, d, meta)
}

func buildCreateV2PersistentVolumeClaimParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": d.Get("name"),
		},
	}

	return bodyParams
}

func waitForCreateV2PersistentVolumeClaimStatus(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Active"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetV2PersistentVolumeClaimDetail(client, ns, name)
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
		return fmt.Errorf("waiting for the status of the V2PersistentVolumeClaim to complete active timeout: %s", err)
	}
	return nil
}

func resourceV2PersistentVolumeClaimRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	namespace := d.Get("namespace").(string)
	name := d.Get("name").(string)
	resp, err := GetV2PersistentVolumeClaimDetail(client, namespace, name)
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

func resourceV2PersistentVolumeClaimUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2PersistentVolumeClaimDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}
	namespace := d.Get("namespace").(string)
	name := d.Get("name").(string)

	deleteV2PersistentVolumeClaimHttpUrl := "apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims/{name}"
	deleteV2PersistentVolumeClaimPath := client.Endpoint + deleteV2PersistentVolumeClaimHttpUrl
	deleteV2PersistentVolumeClaimPath = strings.ReplaceAll(deleteV2PersistentVolumeClaimPath, "{namespace", namespace)
	deleteV2PersistentVolumeClaimPath = strings.ReplaceAll(deleteV2PersistentVolumeClaimPath, "{name}", name)
	deleteV2PersistentVolumeClaimOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteV2PersistentVolumeClaimPath, &deleteV2PersistentVolumeClaimOpt)
	if err != nil {
		return diag.Errorf("error deleting the specifies namespace (%s): %s", namespace, err)
	}

	err = waitForDeleteV2PersistentVolumeClaimStatus(ctx, client, namespace, name, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForDeleteV2PersistentVolumeClaimStatus(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Active", "Terminating"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetV2PersistentVolumeClaimDetail(client, ns, name)
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

func GetV2PersistentVolumeClaimDetail(client *golangsdk.ServiceClient, namespace, name string) (interface{}, error) {
	getV2PersistentVolumeClaimDetailHttpUrl := "apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims/{name}"
	getV2PersistentVolumeClaimDetailPath := client.Endpoint + getV2PersistentVolumeClaimDetailHttpUrl
	getV2PersistentVolumeClaimDetailPath = strings.ReplaceAll(getV2PersistentVolumeClaimDetailPath, "{namespace}", namespace)
	getV2PersistentVolumeClaimDetailPath = strings.ReplaceAll(getV2PersistentVolumeClaimDetailPath, "{name}", name)
	getV2PersistentVolumeClaimDetailOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getV2PersistentVolumeClaimDetailResp, err := client.Request("GET", getV2PersistentVolumeClaimDetailPath, &getV2PersistentVolumeClaimDetailOpt)
	if err != nil {
		return getV2PersistentVolumeClaimDetailResp, err
	}

	return utils.FlattenResponse(getV2PersistentVolumeClaimDetailResp)
}
