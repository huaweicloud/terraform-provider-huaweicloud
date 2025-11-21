package eps

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EPS GET /v1.0/enterprise-projects/migrate-record/list
func DataSourceMigrateRecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMigrateRecordRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Attributes.
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"associated": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_time": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operate_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
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
						"initiate_ep_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"initiate_ep_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin_ep_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin_ep_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_ep_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_ep_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exist_failed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMigrateRecordRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("eps", region)
	if err != nil {
		return diag.Errorf("error creating EPS client: %s", err)
	}

	httpUrl, err := buildMigrateRecordQueryParameter(d)
	if err != nil {
		return diag.Errorf("error build EPS migrate record url: %s", err)
	}

	listPathBase := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	res := make([]interface{}, 0)
	offset := ""
	for {
		listPath := listPathBase + buildOffsetPath(offset)
		getResp, err := client.Request("GET", listPath, &opt)
		if err != nil {
			return diag.Errorf("error retrieving EPS migrate record: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		res = append(res, flattenListMigrateRecordResponseBody(getRespBody)...)
		offset = utils.PathSearch("offset", getRespBody, "").(string)
		if offset == "" {
			break
		}
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("records", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildMigrateRecordQueryParameter(d *schema.ResourceData) (string, error) {
	httpUrl := "v1.0/enterprise-projects/migrate-record/list?limit=10"

	if d.Get("resource_id") != nil {
		httpUrl += fmt.Sprintf("&resource_id=%s", d.Get("resource_id"))
	}

	if begTime, ok := d.GetOk("start_time"); ok {
		earlierTime, err := utils.FormatUTCTimeStamp(begTime.(string))
		if err != nil {
			return "", fmt.Errorf("unable to parse start_time: %s", err)
		}
		httpUrl += fmt.Sprintf("&start_time=%s", strconv.FormatInt(earlierTime*1000, 10))
	}

	if endTime, ok := d.GetOk("end_time"); ok {
		laterTime, err := utils.FormatUTCTimeStamp(endTime.(string))
		if err != nil {
			return "", fmt.Errorf("unable to parse end_time: %s", err)
		}
		httpUrl += fmt.Sprintf("&end_time=%s", strconv.FormatInt(laterTime*1000, 10))
	}

	return httpUrl, nil
}

func buildOffsetPath(offset string) string {
	res := ""
	if offset != "" {
		res += fmt.Sprintf("&offset=%s", offset)
	}
	return res
}

func flattenListMigrateRecordResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("records", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0)
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"associated":       utils.PathSearch("associated", v, nil),
			"code":             utils.PathSearch("code", v, nil),
			"message":          utils.PathSearch("message", v, nil),
			"project_id":       utils.PathSearch("project_id", v, nil),
			"project_name":     utils.PathSearch("project_name", v, nil),
			"region_id":        utils.PathSearch("region_id", v, nil),
			"event_time":       utils.PathSearch("event_time", v, nil),
			"user_name":        utils.PathSearch("user_name", v, nil),
			"operate_type":     utils.PathSearch("operate_type", v, nil),
			"resource_id":      utils.PathSearch("resource_id", v, nil),
			"resource_name":    utils.PathSearch("resource_name", v, nil),
			"resource_type":    utils.PathSearch("resource_type", v, nil),
			"initiate_ep_id":   utils.PathSearch("initiate_ep_id", v, nil),
			"initiate_ep_name": utils.PathSearch("initiate_ep_name", v, nil),
			"origin_ep_id":     utils.PathSearch("origin_ep_id", v, nil),
			"origin_ep_name":   utils.PathSearch("origin_ep_name", v, nil),
			"target_ep_id":     utils.PathSearch("target_ep_id", v, nil),
			"target_ep_name":   utils.PathSearch("target_ep_name", v, nil),
			"exist_failed":     utils.PathSearch("exist_failed", v, nil),
		})
	}
	return res
}
