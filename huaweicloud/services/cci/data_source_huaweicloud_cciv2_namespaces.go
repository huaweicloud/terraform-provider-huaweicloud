package cci

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCI GET /apis/cci/v2/namespaces
// @API CCI GET /apis/cci/v2/namespaces/{name}
func DataSourceV2Namespaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2NamespacesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the namespace.`,
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
						"finalizers": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `The finalizers of the namespace.`,
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
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the namespace.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceV2NamespacesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	nsList := make([]interface{}, 0)
	if ns, ok := d.GetOk("name"); ok {
		resp, err := GetNamespaceDetail(client, ns.(string))
		if err != nil {
			return diag.Errorf("error getting the namespace (%s) from the server: %s", ns.(string), err)
		}
		nsList = append(nsList, resp)
	} else {
		resp, err := listNamespaces(client)
		if err != nil {
			return diag.Errorf("error finding the namespace list from the server: %s", err)
		}
		nsList = utils.PathSearch("items", resp, make([]interface{}, 0)).([]interface{})
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("namespaces", flattenNamespace(nsList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNamespace(nsList []interface{}) []interface{} {
	if len(nsList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(nsList))
	for _, v := range nsList {
		rst = append(rst, map[string]interface{}{
			"name":               utils.PathSearch("metadata.name", v, nil),
			"api_version":        utils.PathSearch("apiVersion", v, nil),
			"kind":               utils.PathSearch("kind", v, nil),
			"annotations":        utils.PathSearch("metadata.annotations", v, nil),
			"labels":             utils.PathSearch("metadata.labels", v, nil),
			"creation_timestamp": utils.PathSearch("metadata.creationTimestamp", v, nil),
			"resource_version":   utils.PathSearch("metadata.resourceVersion", v, nil),
			"uid":                utils.PathSearch("metadata.uid", v, nil),
			"finalizers":         utils.PathSearch("spec.finalizers", v, nil),
			"status":             utils.PathSearch("status.phase", v, nil),
		})
	}
	return rst
}

func listNamespaces(client *golangsdk.ServiceClient) (interface{}, error) {
	listNamespaceHttpUrl := "apis/cci/v2/namespaces"
	listNamespacePath := client.Endpoint + listNamespaceHttpUrl
	listNamespaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listNamespacesResp, err := client.Request("GET", listNamespacePath, &listNamespaceOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying CCI namespaces: %s", err)
	}

	return utils.FlattenResponse(listNamespacesResp)
}
