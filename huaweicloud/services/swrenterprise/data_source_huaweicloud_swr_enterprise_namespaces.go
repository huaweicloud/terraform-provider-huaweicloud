package swrenterprise

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR GET /v2/{project_id}/instances/{instance_id}/namespaces
func DataSourceSwrEnterpriseNamespaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrEnterpriseNamespacesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise instance ID.`,
			},
			"order_column": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the order column.`,
			},
			"order_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the order type.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the namespace name.`,
			},
			"public": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  `Specifies whether the namespace is public.`,
			},
			"namespaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the namespaces.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the namespace ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the namespace name.`,
						},
						"metadata": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the metadata.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"public": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates whether the namespace is public.`,
									},
								},
							},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creation time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the last update time.`,
						},
						"repo_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the repo count of the namespace.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceSwrEnterpriseNamespacesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}
	listNamespacesHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces"
	listNamespacesPath := client.Endpoint + listNamespacesHttpUrl
	listNamespacesPath = strings.ReplaceAll(listNamespacesPath, "{project_id}", client.ProjectID)
	listNamespacesPath = strings.ReplaceAll(listNamespacesPath, "{instance_id}", d.Get("instance_id").(string))
	listNamespacesPath += buildSwrEnterpriseNamespacesQueryParams(d)
	listNamespacesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	offset := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := listNamespacesPath + fmt.Sprintf("&offset=%v", offset)
		listNamespacesResp, err := client.Request("GET", currentPath, &listNamespacesOpt)
		if err != nil {
			return diag.Errorf("error querying SWR namespaces: %s", err)
		}
		listNamespacesRespBody, err := utils.FlattenResponse(listNamespacesResp)
		if err != nil {
			return diag.Errorf("error flattening SWR namespaces response: %s", err)
		}

		namespaces := utils.PathSearch("namespaces", listNamespacesRespBody, make([]interface{}, 0)).([]interface{})
		if len(namespaces) == 0 {
			break
		}
		for _, namespace := range namespaces {
			results = append(results, map[string]interface{}{
				"id":         utils.PathSearch("namespace_id", namespace, nil),
				"name":       utils.PathSearch("name", namespace, nil),
				"metadata":   flattenSwrEnterpriseNamespaceMetadata(namespace),
				"repo_count": utils.PathSearch("repo_count", namespace, nil),
				"created_at": utils.PathSearch("created_at", namespace, nil),
				"updated_at": utils.PathSearch("updated_at", namespace, nil),
			})
		}

		// offset must be the multiple of limit
		offset += 100
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("namespaces", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildSwrEnterpriseNamespacesQueryParams(d *schema.ResourceData) string {
	res := "?limit=100"

	if v, ok := d.GetOk("order_column"); ok {
		res = fmt.Sprintf("%s&order_column=%v", res, v)
	}
	if v, ok := d.GetOk("order_type"); ok {
		res = fmt.Sprintf("%s&order_type=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("public"); ok {
		res = fmt.Sprintf("%s&public=%v", res, v)
	}

	return res
}
