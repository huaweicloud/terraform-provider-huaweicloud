package rgc

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

// @API RGC GET /v1/governance/control-violations
func DataSourceControlViolations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceControlViolationsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organizational_unit_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"control_violations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     ControlViolationsSchema(),
			},
		},
	}
}

func ControlViolationsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_organizational_unit_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_organizational_unit_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceControlViolationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listControlViolationsHttpUrl := "v1/governance/control-violations"
	listControlViolationsProduct := "rgc"
	listControlViolationsClient, err := cfg.NewServiceClient(listControlViolationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	listControlViolationsPath := listControlViolationsClient.Endpoint + listControlViolationsHttpUrl
	listControlViolationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var controlViolations []interface{}
	var marker string
	var queryPath string

	for {
		queryPath = listControlViolationsPath + buildListControlViolationsQueryParams(marker, d)
		listControlViolationsResp, err := listControlViolationsClient.Request("GET", queryPath, &listControlViolationsOpt)
		if err != nil {
			return diag.Errorf("error retrieving RGC control violations: %s", err)
		}

		listControlViolationsRespBody, err := utils.FlattenResponse(listControlViolationsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageOrganizationalUnitControlViolations := FlattenControlViolationsResp(listControlViolationsRespBody)
		controlViolations = append(controlViolations, onePageOrganizationalUnitControlViolations...)
		marker = utils.PathSearch("page_info.next_marker", listControlViolationsRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("control_violations", controlViolations),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListControlViolationsQueryParams(marker string, d *schema.ResourceData) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	if v, ok := d.GetOk("account_id"); ok {
		res = fmt.Sprintf("%s&account_id=%v", res, v)
	}

	if v, ok := d.GetOk("organizational_unit_id"); ok {
		res = fmt.Sprintf("%s&organizational_unit_id=%v", res, v)
	}

	return res
}

func FlattenControlViolationsResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("control_violations", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"account_id":                      utils.PathSearch("account_id", v, nil),
			"account_name":                    utils.PathSearch("account_name", v, nil),
			"display_name":                    utils.PathSearch("display_name", v, nil),
			"name":                            utils.PathSearch("name", v, nil),
			"control_id":                      utils.PathSearch("control_id", v, nil),
			"parent_organizational_unit_id":   utils.PathSearch("parent_organizational_unit_id", v, nil),
			"parent_organizational_unit_name": utils.PathSearch("parent_organizational_unit_name", v, nil),
			"region":                          utils.PathSearch("region", v, nil),
			"resource":                        utils.PathSearch("resource", v, nil),
			"resource_name":                   utils.PathSearch("resource_name", v, nil),
			"resource_type":                   utils.PathSearch("resource_type", v, nil),
			"service":                         utils.PathSearch("service", v, nil),
		})
	}
	return rst
}
