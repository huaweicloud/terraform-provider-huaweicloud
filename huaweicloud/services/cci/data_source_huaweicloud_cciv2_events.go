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

// @API CCI GET /apis/cci/v2/namespaces/{namespace}/events
func DataSourceV2Events() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2EventsRead,

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
			"events": {
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
					},
				},
			},
		},
	}
}

func dataSourceV2EventsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}
	listEventsHttpUrl := "apis/cci/v2/namespaces/{namespace}/events"
	listEventsPath := client.Endpoint + listEventsHttpUrl
	listEventsPath = strings.ReplaceAll(listEventsPath, "{namespace}", d.Get("namespace").(string))
	listEventsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listEventsResp, err := client.Request("GET", listEventsPath, &listEventsOpt)
	if err != nil {
		return diag.Errorf("error getting CCI events: %s", err)
	}

	listEventsRespBody, err := utils.FlattenResponse(listEventsResp)
	if err != nil {
		return diag.Errorf("error retrieving CCI events: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	events := utils.PathSearch("items", listEventsRespBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("events", flattenEvents(events)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEvents(events []interface{}) []interface{} {
	if len(events) == 0 {
		return nil
	}

	rst := make([]interface{}, len(events))
	for i, v := range events {
		rst[i] = map[string]interface{}{
			"name":               utils.PathSearch("metadata.name", v, nil),
			"namespace":          utils.PathSearch("metadata.namespace", v, nil),
			"annotations":        utils.PathSearch("metadata.annotations", v, nil),
			"labels":             utils.PathSearch("metadata.labels", v, nil),
			"creation_timestamp": utils.PathSearch("metadata.creationTimestamp", v, nil),
			"resource_version":   utils.PathSearch("metadata.resourceVersion", v, nil),
			"uid":                utils.PathSearch("metadata.uid", v, nil),
		}
	}
	return rst
}
