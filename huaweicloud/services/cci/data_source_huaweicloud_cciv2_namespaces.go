package cci

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
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
							Elem:        dataSourceNanagedFieldsSchema(),
							Computed:    true,
							Description: `The managed fields of the namespace.`,
						},
						"owner_references": {
							Type:        schema.TypeList,
							Elem:        dataSourceOwnerReferencesSchema(),
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
					},
				},
			},
		},
	}
}

func dataSourceNanagedFieldsSchema() *schema.Resource {
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

func dataSourceOwnerReferencesSchema() *schema.Resource {
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

func dataSourceV2NamespacesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CciV1Client(conf.GetRegion(d))
	// client, err := conf.CciV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCI v2 client: %s", err)
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
			"name":                          utils.PathSearch("metadata.name", v, nil),
			"api_version":                   utils.PathSearch("api_version", v, nil),
			"kind":                          utils.PathSearch("kind", v, nil),
			"annotations":                   utils.PathSearch("metadata.annotations", v, nil),
			"creation_timestamp":            utils.PathSearch("metadata.creationTimestamp", v, nil),
			"cluster_name":                  utils.PathSearch("metadata.clusterName", v, nil),
			"deletion_grace_period_seconds": utils.PathSearch("metadata.deletionGracePeriodSeconds", v, nil),
			"deletion_timestamp":            utils.PathSearch("metadata.deletionTimestamp", v, nil),
			"finaliaers":                    utils.PathSearch("metadata.finaliaers", v, nil),
			"generate_name":                 utils.PathSearch("metadata.generateName", v, nil),
			"generation":                    utils.PathSearch("metadata.generation", v, nil),
			"managed_fields":                flattenNsListManagedFields(utils.PathSearch("metadata.managedFields", v, nil)),
			"owner_references":              flattenNsListOwnerReferences(utils.PathSearch("metadata.ownerReferences", v, nil)),
			"resource_version":              utils.PathSearch("metadata.resourceVersion", v, nil),
			"self_link":                     utils.PathSearch("metadata.selfLink", v, nil),
			"uid":                           utils.PathSearch("metadata.uid", v, nil),
		})
	}
	return rst
}

func flattenNsListOwnerReferences(resp interface{}) []interface{} {
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

func flattenNsListManagedFields(resp interface{}) []interface{} {
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
