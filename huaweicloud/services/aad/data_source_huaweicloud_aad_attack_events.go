package aad

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

// @API AAD GET /v2/aad/domains/waf-info/attack/event
func DataSourceAttackEvents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAttackEventsRead,

		Schema: map[string]*schema.Schema{
			"domains": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the domain name. If not specified, all domains will be included.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the start time.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the end time.",
			},
			"recent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the recent time period.",
			},
			// In the API documentation, `overseas_type` is int type.
			// But it needs to meet the scenario of `0`, so it is defined here as a string type.
			"overseas_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the instance type. `0` represents  mainland China, `1` represents overseas.",
			},
			"sip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the attack source IP.",
			},
			// The field in the API documentation is `list`.
			"events": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The attack events list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The event ID.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The attack target domain.",
						},
						"time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The attack time.",
						},
						"sip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The attack source IP.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The defense action.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The attack URL.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The attack type.",
						},
						"backend": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The current backend information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The current backend protocol.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The current backend port.",
									},
									"host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The current backend host value.",
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

func buildAttackEventsQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("domains"); ok {
		queryParams = fmt.Sprintf("%s&domains=%v", queryParams, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		queryParams = fmt.Sprintf("%s&start_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = fmt.Sprintf("%s&end_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("recent"); ok {
		queryParams = fmt.Sprintf("%s&recent=%v", queryParams, v)
	}
	if v, ok := d.GetOk("overseas_type"); ok {
		queryParams = fmt.Sprintf("%s&overseas_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sip"); ok {
		queryParams = fmt.Sprintf("%s&sip=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceAttackEventsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/domains/waf-info/attack/event"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildAttackEventsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving AAD attack events: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		eventResp := utils.PathSearch("list", respBody, make([]interface{}, 0)).([]interface{})
		if len(eventResp) == 0 {
			break
		}

		result = append(result, eventResp...)
		offset += len(eventResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("events", flattenAttackEvents(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAttackEvents(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":      utils.PathSearch("id", v, nil),
			"domain":  utils.PathSearch("domain", v, nil),
			"time":    utils.PathSearch("time", v, nil),
			"sip":     utils.PathSearch("sip", v, nil),
			"action":  utils.PathSearch("action", v, nil),
			"url":     utils.PathSearch("url", v, nil),
			"type":    utils.PathSearch("type", v, nil),
			"backend": flattenAttackEventBackend(utils.PathSearch("backend", v, nil)),
		})
	}

	return rst
}

func flattenAttackEventBackend(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"protocol": utils.PathSearch("protocol", resp, nil),
			"port":     utils.PathSearch("port", resp, nil),
			"host":     utils.PathSearch("host", resp, nil),
		},
	}
}
