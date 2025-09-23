package cci

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCI GET /apis/cci/v2/namespaces/{namespace}/services
func DataSourceV2Services() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2ServicesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"session_affinity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ports": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"app_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"selector": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"conditions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     servicesStatusConditionsSchema(),
									},
									"loadbalancer": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     servicesStatusLoadBalancerSchema(),
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

func servicesStatusConditionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"observe_generation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_transition_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func servicesStatusLoadBalancerSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"ingress": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ports": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"protocol": {
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

	return &sc
}

func dataSourceV2ServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}
	listServicesHttpUrl := "apis/cci/v2/namespaces/{namespace}/services"
	listServicesPath := client.Endpoint + listServicesHttpUrl
	listServicesPath = strings.ReplaceAll(listServicesPath, "{namespace}", d.Get("namespace").(string))
	listServicesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listServicesResp, err := client.Request("GET", listServicesPath, &listServicesOpt)
	if err != nil {
		return diag.Errorf("error getting CCI services list: %s", err)
	}

	listServicesRespBody, err := utils.FlattenResponse(listServicesResp)
	if err != nil {
		return diag.Errorf("error retrieving CCI services: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	services := utils.PathSearch("items", listServicesRespBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("services", flattenServices(services)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenServices(services []interface{}) []interface{} {
	if len(services) == 0 {
		return nil
	}

	rst := make([]interface{}, len(services))
	for i, v := range services {
		rst[i] = map[string]interface{}{
			"name":               utils.PathSearch("metadata.name", v, nil),
			"namespace":          utils.PathSearch("metadata.namespace", v, nil),
			"annotations":        utils.PathSearch("metadata.annotations", v, nil),
			"labels":             utils.PathSearch("metadata.labels", v, nil),
			"creation_timestamp": utils.PathSearch("metadata.creationTimestamp", v, nil),
			"resource_version":   utils.PathSearch("metadata.resourceVersion", v, nil),
			"uid":                utils.PathSearch("metadata.uid", v, nil),
			"finalizers":         utils.PathSearch("metadata.finalizers", v, nil),
			"ports":              flattenServiceSpecPorts(utils.PathSearch("spec.ports", v, make([]interface{}, 0)).([]interface{})),
			"selector":           utils.PathSearch("spec.selector", v, nil),
			"session_affinity":   utils.PathSearch("spec.sessionAffinity", v, nil),
			"type":               utils.PathSearch("spec.type", v, nil),
			"status":             flattenServiceStatus(utils.PathSearch("status", v, nil)),
		}
	}
	return rst
}
