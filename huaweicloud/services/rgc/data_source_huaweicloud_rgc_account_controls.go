package rgc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RGC GET /v1/governance/managed-accounts/{managed_account_id}/controls
func DataSourceAccountControls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountControlRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"managed_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"control_summaries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     AccountControlSchema(),
			},
		},
	}
}

func AccountControlSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"manage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_objective": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"behavior": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"regional_preference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guidance": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"implementation": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceAccountControlRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	managedAccountId := d.Get("managed_account_id").(string)
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listAccountControlHttpUrl := "v1/governance/managed-accounts/{managed_account_id}/controls"
	listAccountControlProduct := "rgc"
	listAccountControlClient, err := cfg.NewServiceClient(listAccountControlProduct, region)
	if err != nil {
		return diag.Errorf("error creating rgc client: %s", err)
	}

	listAccountControlPath := listAccountControlClient.Endpoint + listAccountControlHttpUrl
	listAccountControlPath = strings.ReplaceAll(listAccountControlPath, "{managed_account_id}", managedAccountId)

	listAccountControlOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var accountControls []interface{}
	var marker string
	var queryPath string

	for {
		queryPath = listAccountControlPath + buildListAccountControlQueryParams(marker)
		listAccountControlResp, err := listAccountControlClient.Request("GET", queryPath, &listAccountControlOpt)
		if err != nil {
			return diag.Errorf("error retrieving RGC account controls: %s", err)
		}

		listAccountControlRespBody, err := utils.FlattenResponse(listAccountControlResp)
		if err != nil {
			return diag.FromErr(err)
		}

		onePageAccountAccountControl := FlattenAccountControlResp(listAccountControlRespBody)
		accountControls = append(accountControls, onePageAccountAccountControl...)
		marker = utils.PathSearch("page_info.next_marker", listAccountControlRespBody, "").(string)
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
		d.Set("control_summaries", accountControls),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListAccountControlQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func FlattenAccountControlResp(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("control_summaries", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"manage_account_id":   utils.PathSearch("manage_account_id", v, nil),
			"control_identifier":  utils.PathSearch("control_identifier", v, nil),
			"state":               utils.PathSearch("state", v, nil),
			"version":             utils.PathSearch("version", v, nil),
			"name":                utils.PathSearch("name", v, nil),
			"description":         utils.PathSearch("description", v, nil),
			"control_objective":   utils.PathSearch("control_objective", v, nil),
			"behavior":            utils.PathSearch("behavior", v, nil),
			"owner":               utils.PathSearch("owner", v, nil),
			"regional_preference": utils.PathSearch("regional_preference", v, nil),
			"guidance":            utils.PathSearch("guidance", v, nil),
			"service":             utils.PathSearch("service", v, nil),
			"implementation":      utils.PathSearch("implementation", v, nil),
		})
	}
	return rst
}
