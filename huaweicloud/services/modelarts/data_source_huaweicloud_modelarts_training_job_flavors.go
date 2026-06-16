package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v2/{project_id}/training-job-flavors
func DataSourceTrainingJobFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrainingJobFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the training job flavors are located.`,
			},
			"flavor_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the flavor.`,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flavor_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the flavor.`,
						},
						"flavor_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the flavor.`,
						},
						"flavor_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the flavor. The valid values are **CPU**, **GPU**, and **Ascend**.`,
						},
						"billing": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The billing code.`,
									},
									"unit_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The billing unit.`,
									},
								},
							},
							Description: `The billing information of the flavor.`,
						},
						"flavor_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The maximum number of nodes that can be selected.`,
									},
									"cpu": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"arch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The CPU architecture.`,
												},
												"core_num": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: `The number of CPU cores.`,
												},
											},
										},
										Description: `The CPU information of the flavor.`,
									},
									"gpu": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit_num": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: `The number of GPUs.`,
												},
												"product_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The GPU product name.`,
												},
												"memory": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The GPU memory.`,
												},
											},
										},
										Description: `The GPU information of the flavor.`,
									},
									"npu": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit_num": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: `The number of NPUs.`,
												},
												"product_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The NPU product name.`,
												},
												"memory": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The NPU memory.`,
												},
											},
										},
										Description: `The Ascend information of the flavor.`,
									},
									"memory": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: `The memory size.`,
												},
												"unit": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The memory unit.`,
												},
											},
										},
										Description: `The memory information of the flavor.`,
									},
									"disk": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: `The disk size.`,
												},
												"unit": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The disk unit.`,
												},
											},
										},
										Description: `The disk information of the flavor.`,
									},
								},
							},
							Description: `The detailed information of the flavor.`,
						},
						"attributes": {
							Type:        schema.TypeMap,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `The other attributes of the flavor.`,
						},
						"support_engines": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The engines supported by the flavor.`,
						},
					},
				},
				Description: `The list of training job flavors that match the filter parameters.`,
			},
		},
	}
}

func buildTrainingJobFlavorsQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("flavor_type"); ok {
		return fmt.Sprintf("?flavor_type=%v", v)
	}

	return ""
}

func listTrainingJobFlavors(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "v2/{project_id}/training-job-flavors"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildTrainingJobFlavorsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("flavors", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func dataSourceTrainingJobFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	flavors, err := listTrainingJobFlavors(client, d)
	if err != nil {
		return diag.Errorf("error querying training job flavors: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("flavors", flattenTrainingJobFlavors(flavors)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTrainingJobFlavors(flavors []interface{}) []map[string]interface{} {
	if len(flavors) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(flavors))
	for _, flavor := range flavors {
		result = append(result, map[string]interface{}{
			"flavor_id":   utils.PathSearch("flavor_id", flavor, nil),
			"flavor_name": utils.PathSearch("flavor_name", flavor, nil),
			"flavor_type": utils.PathSearch("flavor_type", flavor, nil),
			"billing": flattenTrainingJobFlavorsBilling(utils.PathSearch("billing", flavor,
				make(map[string]interface{})).(map[string]interface{})),
			"flavor_info": flattenTrainingJobFlavorsInfo(utils.PathSearch("flavor_info", flavor,
				make(map[string]interface{})).(map[string]interface{})),
			"attributes":      utils.PathSearch("attributes", flavor, nil),
			"support_engines": utils.PathSearch("support_engines", flavor, nil),
		})
	}

	return result
}

func flattenTrainingJobFlavorsBilling(billing map[string]interface{}) []map[string]interface{} {
	if len(billing) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"code":     utils.PathSearch("code", billing, nil),
			"unit_num": utils.PathSearch("unit_num", billing, nil),
		},
	}
}

func flattenTrainingJobFlavorsInfo(flavorInfo map[string]interface{}) []map[string]interface{} {
	if len(flavorInfo) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"max_num": utils.PathSearch("max_num", flavorInfo, nil),
			"cpu": flattenTrainingJobFlavorsInfoCpu(utils.PathSearch("cpu", flavorInfo,
				make(map[string]interface{})).(map[string]interface{})),
			"gpu": flattenTrainingJobFlavorsInfoGpu(utils.PathSearch("gpu", flavorInfo,
				make(map[string]interface{})).(map[string]interface{})),
			"npu": flattenTrainingJobFlavorsInfoNpu(utils.PathSearch("npu", flavorInfo,
				make(map[string]interface{})).(map[string]interface{})),
			"memory": flattenTrainingJobFlavorsInfoMemory(utils.PathSearch("memory", flavorInfo,
				make(map[string]interface{})).(map[string]interface{})),
			"disk": flattenTrainingJobFlavorsInfoDisk(utils.PathSearch("disk", flavorInfo,
				make(map[string]interface{})).(map[string]interface{})),
		},
	}
}

func flattenTrainingJobFlavorsInfoCpu(cpu map[string]interface{}) []map[string]interface{} {
	if len(cpu) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"arch":     utils.PathSearch("arch", cpu, nil),
			"core_num": utils.PathSearch("core_num", cpu, nil),
		},
	}
}

func flattenTrainingJobFlavorsInfoGpu(gpu map[string]interface{}) []map[string]interface{} {
	if len(gpu) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"unit_num":     utils.PathSearch("unit_num", gpu, nil),
			"product_name": utils.PathSearch("product_name", gpu, nil),
			"memory":       utils.PathSearch("memory", gpu, nil),
		},
	}
}

func flattenTrainingJobFlavorsInfoNpu(npu map[string]interface{}) []map[string]interface{} {
	if len(npu) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"unit_num":     utils.PathSearch("unit_num", npu, nil),
			"product_name": utils.PathSearch("product_name", npu, nil),
			"memory":       utils.PathSearch("memory", npu, nil),
		},
	}
}

func flattenTrainingJobFlavorsInfoMemory(memory map[string]interface{}) []map[string]interface{} {
	if len(memory) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"size": utils.PathSearch("size", memory, nil),
			"unit": utils.PathSearch("unit", memory, nil),
		},
	}
}

func flattenTrainingJobFlavorsInfoDisk(disk map[string]interface{}) []map[string]interface{} {
	if len(disk) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"size": utils.PathSearch("size", disk, nil),
			"unit": utils.PathSearch("unit", disk, nil),
		},
	}
}
