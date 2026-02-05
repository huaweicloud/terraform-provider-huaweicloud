package hss

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

// @API HSS GET /v5/{project_id}/event/severity
func DataSourceEventSeverity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventSeverityRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"begin_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"last_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"container_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"handle_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severity_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attack_tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asset_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"att_ck": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"low_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"medium_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"high_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"critical_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildEventSeverityQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?category=%s", d.Get("category"))
	severityList := d.Get("severity_list").([]interface{})
	tagList := d.Get("tag_list").([]interface{})

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	for _, key := range []string{
		"begin_time",
		"end_time",
		"last_days",
		"host_name",
		"host_id",
		"private_ip",
		"public_ip",
		"container_name",
		"event_type",
		"handle_status",
		"severity",
		"attack_tag",
		"asset_value",
		"att_ck",
		"event_name",
	} {
		if v, ok := d.GetOk(key); ok {
			queryParams = fmt.Sprintf("%s&%s=%v", queryParams, key, v)
		}
	}

	if len(severityList) > 0 {
		for _, v := range severityList {
			queryParams = fmt.Sprintf("%s&severity_list=%v", queryParams, v)
		}
	}

	if len(tagList) > 0 {
		for _, v := range tagList {
			queryParams = fmt.Sprintf("%s&tag_list=%v", queryParams, v)
		}
	}

	return queryParams
}

func dataSourceEventSeverityRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/event/severity"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildEventSeverityQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS event severity: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", utils.PathSearch("total_num", respBody, nil)),
		d.Set("low_num", utils.PathSearch("low_num", respBody, nil)),
		d.Set("medium_num", utils.PathSearch("medium_num", respBody, nil)),
		d.Set("high_num", utils.PathSearch("high_num", respBody, nil)),
		d.Set("critical_num", utils.PathSearch("critical_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
