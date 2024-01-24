package fgs

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/function"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API FunctionGraph GET /v2/{project_id}/fgs/functions
func DataSourceFunctionGraphFunctions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFunctionGraphFunctionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"package_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"runtime": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"functions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"package": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"handler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"code_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code_filename": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypted_user_data": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_agency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_instance_num": {
							// The original type of this parameter is int, but its zero value is meaningful.
							// So, the following types of parameter passing are realized through the logic of terraform's implicit
							// conversion of int:
							//   + -1: the number of instances is unlimited.
							//   + 0: this function is disabled.
							//   + (0, +1000]: Specific value (2023.06.26).
							//   + empty: keep the default (latest updated) value.
							Type:     schema.TypeString,
							Computed: true,
						},
						"initializer_handler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"initializer_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_stream_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"functiongraph_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceFunctionGraphFunctionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	fgsClient, err := conf.FgsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph V2 client: %s", err)
	}

	// MaxItems and Marker use default values.
	opts := function.ListOpts{
		PackageName: d.Get("package_name").(string),
	}
	allPages, err := function.List(fgsClient, opts).AllPages()
	if err != nil {
		return diag.Errorf("error querying functions: %s", err)
	}
	resp, err := function.ExtractList(allPages)
	if err != nil {
		return diag.Errorf("error querying functions: %s", err)
	}

	filterResult, err := filterFunctionList(d, resp.Functions, conf)
	if err != nil {
		return diag.FromErr(err)
	}
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("functions", flattenFunctions(filterResult)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving datas of FunctionGraph functions: %s", err)
	}

	return nil
}

func filterFunctionList(d *schema.ResourceData, functions []function.Function, conf *config.Config) ([]interface{}, error) {
	filter := map[string]interface{}{
		"FuncName":            d.Get("name"),
		"FuncUrn":             d.Get("urn"),
		"Runtime":             d.Get("runtime"),
		"EnterpriseProjectID": conf.GetEnterpriseProjectID(d),
	}

	filterResult, err := utils.FilterSliceWithField(functions, filter)
	if err != nil {
		return nil, fmt.Errorf("error filtering list of functions: %s", err)
	}

	return filterResult, nil
}

func flattenFunctions(functions []interface{}) []map[string]interface{} {
	if len(functions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(functions))
	for i, val := range functions {
		f := val.(function.Function)
		result[i] = map[string]interface{}{
			"name":                  f.FuncName,
			"urn":                   f.FuncUrn,
			"package":               f.Package,
			"runtime":               f.Runtime,
			"timeout":               f.Timeout,
			"handler":               f.Handler,
			"memory_size":           f.MemorySize,
			"code_type":             f.CodeType,
			"code_url":              f.CodeUrl,
			"code_filename":         f.CodeFileName,
			"user_data":             f.UserData,
			"encrypted_user_data":   f.EncryptedUserData,
			"version":               f.Version,
			"agency":                f.Xrole,
			"app_agency":            f.AppXrole,
			"description":           f.Description,
			"vpc_id":                f.FuncVpc.VpcId,
			"network_id":            f.FuncVpc.SubnetId,
			"max_instance_num":      strconv.Itoa(*f.StrategyConfig.Concurrency),
			"initializer_handler":   f.InitializerHandler,
			"initializer_timeout":   f.InitializerTimeout,
			"enterprise_project_id": f.EnterpriseProjectID,
			"log_group_id":          f.LogGroupId,
			"log_stream_id":         f.LogStreamId,
			"functiongraph_version": f.Type,
		}
	}
	return result
}
