package cci

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCI GET /apis/cci/v2/storageclasses
func DataSourceV2StorageClasses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2StorageClassesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_classes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_volume_expansion": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"allowed_topologies": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_label_expressions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"values": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
						"annotations": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_options": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"parameters": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"provisioner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reclaim_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_binding_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceV2StorageClassesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}
	listStorageClassesHttpUrl := "apis/cci/v2/storageclasses"
	listStorageClassesPath := client.Endpoint + listStorageClassesHttpUrl
	listStorageClassesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listStorageClassesResp, err := client.Request("GET", listStorageClassesPath, &listStorageClassesOpt)
	if err != nil {
		return diag.Errorf("error getting CCI storage classes: %s", err)
	}

	listStorageClassesRespBody, err := utils.FlattenResponse(listStorageClassesResp)
	if err != nil {
		return diag.Errorf("error retrieving CCI storage classes: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	storageClasses := utils.PathSearch("items", listStorageClassesRespBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("storage_classes", flattenStorageClasses(storageClasses)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenStorageClasses(storageClasses []interface{}) []interface{} {
	if len(storageClasses) == 0 {
		return nil
	}

	rst := make([]interface{}, len(storageClasses))
	for i, v := range storageClasses {
		allowedTopologies := utils.PathSearch("allowedTopologies", v, make([]interface{}, 0)).([]interface{})
		rst[i] = map[string]interface{}{
			"name":                   utils.PathSearch("metadata.name", v, nil),
			"annotations":            utils.PathSearch("metadata.annotations", v, nil),
			"labels":                 utils.PathSearch("metadata.labels", v, nil),
			"creation_timestamp":     utils.PathSearch("metadata.creationTimestamp", v, nil),
			"resource_version":       utils.PathSearch("metadata.resourceVersion", v, nil),
			"uid":                    utils.PathSearch("metadata.uid", v, nil),
			"allow_volume_expansion": utils.PathSearch("allowVolumeExpansion", v, nil),
			"mount_options":          utils.PathSearch("mountOptions", v, nil),
			"parameters":             utils.PathSearch("parameters", v, nil),
			"provisioner":            utils.PathSearch("provisioner", v, nil),
			"reclaim_policy":         utils.PathSearch("reclaimPolicy", v, nil),
			"volume_binding_mode":    utils.PathSearch("volumeBindingMode", v, nil),
			"allowed_topologies":     flattenAllowedTopologies(allowedTopologies),
		}
	}
	return rst
}

func flattenAllowedTopologies(allowedTopologies []interface{}) []interface{} {
	if len(allowedTopologies) == 0 {
		return nil
	}

	rst := make([]interface{}, len(allowedTopologies))
	for i, v := range allowedTopologies {
		matchLabelExpressions := utils.PathSearch("matchLabelExpressions", v, make([]interface{}, 0)).([]interface{})
		rst[i] = map[string]interface{}{
			"match_label_expressions": flattenMatchLabelExpressions(matchLabelExpressions),
		}
	}
	return rst
}

func flattenMatchLabelExpressions(matchLabelExpressions []interface{}) []interface{} {
	if len(matchLabelExpressions) == 0 {
		return nil
	}

	rst := make([]interface{}, len(matchLabelExpressions))
	for i, v := range matchLabelExpressions {
		rst[i] = map[string]interface{}{
			"key":    utils.PathSearch("key", v, nil),
			"values": utils.PathSearch("values", v, nil),
		}
	}
	return rst
}
