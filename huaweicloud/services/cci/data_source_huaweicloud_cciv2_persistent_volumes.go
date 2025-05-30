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

// @API CCI GET /apis/cci/v2/persistentvolumes
func DataSourceV2PersistentVolumes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2PersistentVolumesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"persistent_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
						"access_modes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"capacity": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"claim_ref": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"field_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"namespace": {
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
								},
							},
						},
						"csi": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     pvsCSISchema(),
						},
						"mount_options": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"node_affinity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     pvsVolumeNodeAffinitySchema(),
						},
						"reclaim_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_class_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_mode": {
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
						"finalizers": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"phase": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"reason": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func pvsCSISchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"driver": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_handle": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fs_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"volume_attributes": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"controller_expand_secret_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     pvsSecretReferenceSchema(),
			},
			"controller_publish_secret_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     pvsSecretReferenceSchema(),
			},
			"node_expand_secret_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     pvsSecretReferenceSchema(),
			},
			"node_publish_secret_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     pvsSecretReferenceSchema(),
			},
			"node_stage_secret_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     pvsSecretReferenceSchema(),
			},
		},
	}

	return &sc
}

func pvsSecretReferenceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func pvsVolumeNodeAffinitySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"required": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     pvsNodeSelectorSchema(),
			},
		},
	}

	return &sc
}

func pvsNodeSelectorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_selector_terms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     pvsNodeSelectorTermSchema(),
			},
		},
	}

	return &sc
}

func pvsNodeSelectorTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_expressions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     pvsNodeSelectorRequirementSchema(),
			},
		},
	}

	return &sc
}

func pvsNodeSelectorRequirementSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

	return &sc
}

func dataSourceV2PersistentVolumesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}
	listPersistentVolumesHttpUrl := "apis/cci/v2/persistentvolumes"
	listPersistentVolumesPath := client.Endpoint + listPersistentVolumesHttpUrl
	listPersistentVolumesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listPersistentVolumesResp, err := client.Request("GET", listPersistentVolumesPath, &listPersistentVolumesOpt)
	if err != nil {
		return diag.Errorf("error getting CCI persistent volumes list: %s", err)
	}

	listPersistentVolumesRespBody, err := utils.FlattenResponse(listPersistentVolumesResp)
	if err != nil {
		return diag.Errorf("error retrieving CCI persistent volumes: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	pvs := utils.PathSearch("items", listPersistentVolumesRespBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("persistent_volumes", flattenPersistentVolumes(pvs)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPersistentVolumes(pvs []interface{}) []interface{} {
	if len(pvs) == 0 {
		return nil
	}

	rst := make([]interface{}, len(pvs))
	for i, v := range pvs {
		rst[i] = map[string]interface{}{
			"name":               utils.PathSearch("metadata.name", v, nil),
			"annotations":        utils.PathSearch("metadata.annotations", v, nil),
			"labels":             utils.PathSearch("metadata.labels", v, nil),
			"creation_timestamp": utils.PathSearch("metadata.creationTimestamp", v, nil),
			"resource_version":   utils.PathSearch("metadata.resourceVersion", v, nil),
			"uid":                utils.PathSearch("metadata.uid", v, nil),
			"finalizers":         utils.PathSearch("metadata.finalizers", v, nil),
			"access_modes":       utils.PathSearch("spec.accessModes", v, nil),
			"capacity":           utils.PathSearch("spec.capacity", v, nil),
			"claim_ref":          flattenClaimRef(utils.PathSearch("spec.claimRef", v, nil)),
			"csi":                flattenCSI(utils.PathSearch("spec.csi", v, nil)),
			"mount_options":      utils.PathSearch("spec.mountOptions", v, nil),
			"node_affinity":      flattenPvNodeAffinity(utils.PathSearch("spec.nodeAffinity", v, nil)),
			"reclaim_policy":     utils.PathSearch("spec.reclaimPolicy", v, nil),
			"storage_class_name": utils.PathSearch("spec.storageClassName", v, nil),
			"volume_mode":        utils.PathSearch("spec.volumeMode", v, nil),
			"status":             flattenPersistentVolumeStatus(utils.PathSearch("status", v, nil)),
		}
	}
	return rst
}
