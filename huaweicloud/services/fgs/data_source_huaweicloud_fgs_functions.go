package fgs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API FunctionGraph GET /v2/{project_id}/fgs/functions
func DataSourceFunctions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFunctionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the functions are located.`,
			},
			"package_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The package name used to query the functions.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The function URN used to query the specified function.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The function name used to query the specified function.`,
			},
			"runtime": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The dependency package runtime used to query the functions.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project to which the functions belong.`,
			},
			"functions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The function name.`,
						},
						"urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The function URN.`,
						},
						"package": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The package name that function used.`,
						},
						"runtime": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The dependency package runtime of the function.`,
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The timeout interval of the function.`,
						},
						"handler": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The entry point of the function.`,
						},
						"memory_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The memory size(MB) allocated to the function.`,
						},
						"code_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The function code type.`,
						},
						"code_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The code URL.`,
						},
						"code_filename": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the function file.`,
						},
						"user_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The custom user data (key/value pairs) defined for the function.`,
						},
						"encrypted_user_data": {
							Type:        schema.TypeString,
							Sensitive:   true,
							Computed:    true,
							Description: `The custom user data (key/value pairs) defined to be encrypted for the function.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The function version.`,
						},
						"agency": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IAM agency name for the function configuration.`,
						},
						"app_agency": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IAM agency name for the function execution.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the function.`,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC ID to which the function belongs.`,
						},
						"network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The network ID of subnet to which the function belongs.`,
						},
						"max_instance_num": {
							// The original type of this parameter is int, but its zero value is meaningful.
							// So, the following types of parameter passing are realized through the logic of terraform's implicit
							// conversion of int:
							//   + -1: the number of instances is unlimited.
							//   + 0: this function is disabled.
							//   + (0, +1000]: Specific value (2023.06.26).
							//   + empty: keep the default (latest updated) value.
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The maximum number of instances for a single function.`,
						},
						"initializer_handler": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The initializer of the function.`,
						},
						"initializer_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum duration the function can be initialized.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID to which the function belongs.`,
						},
						"log_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The LTS log group ID.`,
						},
						"log_stream_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The LTS log stream ID.`,
						},
						"functiongraph_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The functionGraph version.`,
						},
					},
				},
				Description: `All functions that match the filter parameters.`,
			},
		},
	}
}

func buildFunctionsQueryParams(d *schema.ResourceData) string {
	if pkgName, ok := d.GetOk("package_name"); ok {
		return fmt.Sprintf("&package_name=%v", pkgName)
	}
	return ""
}

func getFunctions(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/fgs/functions?maxitems=100"
		marker  float64
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildFunctionsQueryParams(d)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithMarker := fmt.Sprintf("%s&marker=%v", listPath, marker)
		requestResp, err := client.Request("GET", listPathWithMarker, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		functions := utils.PathSearch("functions", respBody, make([]interface{}, 0)).([]interface{})
		if len(functions) < 1 {
			break
		}
		result = append(result, functions...)
		// In this API, marker has the same meaning as offset.
		nextMarker := utils.PathSearch("next_marker", respBody, float64(0)).(float64)
		if nextMarker == marker || nextMarker == 0 {
			// Make sure the next marker value is correct, not the previous marker or zero (in the last page).
			break
		}
		marker = nextMarker
	}

	return result, nil
}

// In-place modification of slices will cause data confusion in concurrent call scenarios.
// Although this scenario does not exist, this method is kept for the sake of understanding and avoiding missing
// considerations when copying other methods.
func filterFunctions(cfg *config.Config, d *schema.ResourceData, functions []interface{}) []interface{} {
	result := functions

	if name, ok := d.GetOk("name"); ok && len(result) > 0 {
		result = utils.PathSearch(fmt.Sprintf("[?func_name=='%v']", name), result, make([]interface{}, 0)).([]interface{})
	}

	if urn, ok := d.GetOk("urn"); ok && len(result) > 0 {
		result = utils.PathSearch(fmt.Sprintf("[?func_urn=='%v']", urn), result, make([]interface{}, 0)).([]interface{})
	}

	if runtime, ok := d.GetOk("runtime"); ok && len(result) > 0 {
		result = utils.PathSearch(fmt.Sprintf("[?runtime=='%v']", runtime), result, make([]interface{}, 0)).([]interface{})
	}

	if epsId := cfg.GetEnterpriseProjectID(d); epsId != "" && len(result) > 0 {
		result = utils.PathSearch(fmt.Sprintf("[?enterprise_project_id=='%v']", epsId), result, make([]interface{}, 0)).([]interface{})
	}

	return result
}

func flattenFunctions(functions []interface{}) []map[string]interface{} {
	if len(functions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(functions))
	for _, function := range functions {
		result = append(result, map[string]interface{}{
			"name":                  utils.PathSearch("func_name", function, nil),
			"urn":                   utils.PathSearch("func_urn", function, nil),
			"package":               utils.PathSearch("package", function, nil),
			"runtime":               utils.PathSearch("runtime", function, nil),
			"timeout":               utils.PathSearch("timeout", function, nil),
			"handler":               utils.PathSearch("handler", function, nil),
			"memory_size":           utils.PathSearch("memory_size", function, nil),
			"code_type":             utils.PathSearch("code_type", function, nil),
			"code_url":              utils.PathSearch("code_url", function, nil),
			"code_filename":         utils.PathSearch("code_filename", function, nil),
			"user_data":             utils.PathSearch("user_data", function, nil),
			"encrypted_user_data":   utils.PathSearch("encrypted_user_data", function, nil),
			"version":               utils.PathSearch("version", function, nil),
			"agency":                utils.PathSearch("xrole", function, nil),
			"app_agency":            utils.PathSearch("app_xrole", function, nil),
			"description":           utils.PathSearch("description", function, nil),
			"vpc_id":                utils.PathSearch("func_vpc.vpc_id", function, nil),
			"network_id":            utils.PathSearch("func_vpc.subnet_id", function, nil),
			"max_instance_num":      strconv.Itoa(int(utils.PathSearch("strategy_config.concurrency", function, float64(0)).(float64))),
			"initializer_handler":   utils.PathSearch("initializer_handler", function, nil),
			"initializer_timeout":   utils.PathSearch("initializer_timeout", function, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", function, nil),
			"log_group_id":          utils.PathSearch("log_group_id", function, nil),
			"log_stream_id":         utils.PathSearch("log_stream_id", function, nil),
			"functiongraph_version": utils.PathSearch("type", function, nil),
		})
	}

	return result
}

func dataSourceFunctionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	functions, err := getFunctions(client, d)
	if err != nil {
		return diag.Errorf("error querying functions: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("functions", flattenFunctions(filterFunctions(cfg, d, functions))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
