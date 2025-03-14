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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCI POST /apis/yangtse/v2/namespaces/{namespace}/networks
// @API CCI GET /apis/yangtse/v2/namespaces/{namespace}/networks/{name}
// @API CCI PATCH /apis/yangtse/v2/namespaces/{namespace}/networks/{name}
// @API CCI PUT /apis/yangtse/v2/namespaces/{namespace}/networks/{name}
// @API CCI DELETE /apis/yangtse/v2/namespaces/{namespace}/networks/{name}
func ResourceV2Network() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2NetworkCreate,
		ReadContext:   resourceV2NetworkRead,
		// UpdateContext: resourceV2NetworkUpdate,
		DeleteContext: resourceV2NetworkDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2NetworkImportState,
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
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the CCI network.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the annotations of the CCI network.`,
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
			"ip_families": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the IP families of the CCI network.`,
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the security group IDs of the CCI network.`,
			},
			"subnets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the subnet ID of the CCI network.`,
						},
					},
				},
				Description: `Specifies the subnets of the CCI network.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the CCI network.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kind of the CCI network.`,
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

func buildCreateV2NetworkParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"kind":       "Namespace",
		"apiVersion": "v2",
		"metadata": map[string]interface{}{
			"name":        d.Get("name"),
			"annotations": d.Get("annotations"),
			"labels":      d.Get("labels"),
		},
		"spec": map[string]interface{}{
			"ipFamilies":     d.Get("ip_families"),
			"networkType":    "underlay_neutron",
			"securityGroups": d.Get("security_group_ids"),
			"subnets":        d.Get("subnets"),
		},
	}

	return bodyParams
}

func resourceV2NetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.CciV1Client(conf.GetRegion(d))
	// client, err := conf.CciV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CCI v2 client: %s", err)
	}

	createNetworkHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/networks"
	createNetworkPath := client.Endpoint + createNetworkHttpUrl
	createNetworkPath = strings.ReplaceAll(createNetworkPath, "{namespace}", d.Get("namespace").(string))
	createNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createNetworkOpt.JSONBody = utils.RemoveNil(buildCreateV2NetworkParams(d))

	resp, err := client.Request("POST", createNetworkPath, &createNetworkOpt)
	if err != nil {
		return diag.Errorf("error creating CCI Network: %s", err)
	}

	ns := utils.PathSearch("metadata.name", resp, "").(string)
	if ns == "" {
		return diag.Errorf("unable to find CCI Network name from API response")
	}
	d.SetId(ns)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Initializing", "Pending"},
		Target:       []string{"Active"},
		Refresh:      waitForV2NetworkActive(client, ns, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waitting CCI network status: %s", err)
	}
	return resourceCciNetworkRead(ctx, d, meta)
}

func resourceV2NetworkRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.CciV1BetaClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CCI Beta v1 client: %s", err)
	}

	ns := d.Get("namespace").(string)
	resp, err := GetNetwork(client, ns, d.Get("name").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying CCI network")
	}

	mErr := multierror.Append(
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("api_version", utils.PathSearch("api_version", resp, nil)),
		d.Set("namespace", utils.PathSearch("metadata.namespace", resp, nil)),
		d.Set("annotations", utils.PathSearch("metadata.annotations", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("self_link", utils.PathSearch("metadata.selfLink", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
		d.Set("ip_families", utils.PathSearch("metadata.spec.ipFamilies", resp, nil)),
		d.Set("security_group_ids", utils.PathSearch("metadata.spec.securityGroups", resp, nil)),
		d.Set("subnets", utils.PathSearch("metadata.spec.subnets", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV2NetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.CciV1BetaClient(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CCI Beta v1 client: %s", err)
	}

	ns := d.Get("namespace").(string)
	deleteNetworkHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/networks/{name}"
	deleteNetworkPath := client.Endpoint + deleteNetworkHttpUrl
	deleteNetworkPath = strings.ReplaceAll(deleteNetworkPath, "{namespace}", ns)
	deleteNetworkPath = strings.ReplaceAll(deleteNetworkPath, "{name}", d.Get("name").(string))
	deleteNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteNetworkPath, &deleteNetworkOpt)
	if err != nil {
		return diag.Errorf("error deleting CCI network: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Terminating", "Active"},
		Target:     []string{"Deleted"},
		Refresh:    waitForV2NetworkDelete(client, ns, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waitting CCI network status: %s", err)
	}

	return nil
}

func waitForV2NetworkActive(client *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetNetwork(client, ns, name)
		status := utils.PathSearch("status.state", resp, "").(string)
		if err != nil {
			return nil, "", err
		}

		return resp, status, nil
	}
}

func waitForV2NetworkDelete(client *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete CCI network %s.", name)

		resp, err := GetNetwork(client, ns, name)
		status := utils.PathSearch("status.state", resp, "").(string)
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[DEBUG] successfully deleted CCI network: %s", name)
			return resp, "Deleted", nil
		}
		if status == "Terminating" {
			return resp, "Terminating", nil
		}
		log.Printf("[DEBUG] CCI network %s still available", name)
		return resp, "Active", nil
	}
}

func resourceV2NetworkImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<namespace>/<id>', but '%s'", importedId)
	}

	d.SetId(parts[1])
	d.Set("namespace", parts[0])

	return []*schema.ResourceData{d}, nil
}

func GetNetwork(client *golangsdk.ServiceClient, namespace, name string) (interface{}, error) {
	getNetworkHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/networks/{name}"
	getNetworkPath := client.Endpoint + getNetworkHttpUrl
	getNetworkPath = strings.ReplaceAll(getNetworkPath, "{namespace}", namespace)
	getNetworkPath = strings.ReplaceAll(getNetworkPath, "{name}", name)
	getNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getNetworkResp, err := client.Request("GET", getNetworkPath, &getNetworkOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying CCI Network: %s", err)
	}

	return utils.FlattenResponse(getNetworkResp)
}
