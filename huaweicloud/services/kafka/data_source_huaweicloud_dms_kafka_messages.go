package kafka

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

// @API Kafka GET /v2/{project_id}/instances/{instance_id}/messages
func DataSourceDmsKafkaMessages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsKafkaMessagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the topic name.`,
			},
			"start_time": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"end_time"},
				Description:  `Specifies the start time, a Unix timestamp in millisecond.`,
			},
			"end_time": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"start_time"},
				Description:  `Specifies the end time, a Unix timestamp in millisecond.`,
			},
			"download": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether download is required.`,
			},
			"message_offset": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the message offset.`,
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the partition.`,
			},
			"keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the keyword.`,
			},
			"messages": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the message list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the message key.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the message content.`,
						},
						"timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the message production time.`,
						},
						"huge_message": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates the big data flag.`,
						},
						"message_offset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the message offset.`,
						},
						"partition": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the partition where the message is located.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the message size.`,
						},
						"message_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the message ID.`,
						},
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the application ID.`,
						},
						"tag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the message label.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDmsKafkaMessagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	listMessagesHttpUrl := "v2/{project_id}/instances/{instance_id}/messages"
	listMessagesPath := client.Endpoint + listMessagesHttpUrl
	listMessagesPath = strings.ReplaceAll(listMessagesPath, "{project_id}", client.ProjectID)
	listMessagesPath = strings.ReplaceAll(listMessagesPath, "{instance_id}", d.Get("instance_id").(string))
	listMessagesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listMessagesPath = buildQueryMessagesListPath(d, listMessagesPath)

	// actually, offset is pageNo in API, not offset
	var offset int
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := listMessagesPath + fmt.Sprintf("&limit=%d&offset=%d", pageLimit, offset)
		listMessagesResp, err := client.Request("GET", currentPath, &listMessagesOpt)
		if err != nil {
			return diag.Errorf("error retrieving messages: %s", err)
		}
		listMessagesRespBody, err := utils.FlattenResponse(listMessagesResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		messages := utils.PathSearch("messages", listMessagesRespBody, make([]interface{}, 0)).([]interface{})
		for _, message := range messages {
			// there are two different kinds of response
			// if `keyword` is specified, the response will be `offset`
			offset := utils.PathSearch("message_offset", message, nil)
			if offset == nil {
				offset = utils.PathSearch("offset", message, nil)
			}

			results = append(results, map[string]interface{}{
				"key":            utils.PathSearch("key", message, nil),
				"value":          utils.PathSearch("value", message, nil),
				"huge_message":   utils.PathSearch("huge_message", message, nil),
				"message_offset": offset,
				"partition":      utils.PathSearch("partition", message, nil),
				"size":           utils.PathSearch("size", message, nil),
				"message_id":     utils.PathSearch("message_id", message, nil),
				"app_id":         utils.PathSearch("app_id", message, nil),
				"tag":            utils.PathSearch("tag", message, nil),
				"timestamp": utils.FormatTimeStampRFC3339(
					int64(utils.PathSearch("timestamp", message, float64(0)).(float64))/1000, true),
			})
		}

		// if keyword is specified, the limit and offset is useless, do not query next page
		if _, ok := d.GetOk("keyword"); ok {
			break
		}

		// `total` means the number of all `messages`, and type is float64.
		total := utils.PathSearch("total", listMessagesRespBody, float64(0))
		if int(total.(float64)) <= (offset+1)*pageLimit {
			break
		}
		offset++
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("messages", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildQueryMessagesListPath(d *schema.ResourceData, listMessagesPath string) string {
	listMessagesPath += fmt.Sprintf("?topic=%v", d.Get("topic"))
	if startTime, ok := d.GetOk("start_time"); ok {
		listMessagesPath += fmt.Sprintf("&start_time=%s", startTime)
		listMessagesPath += fmt.Sprintf("&end_time=%s", d.Get("end_time"))
	}
	if download, ok := d.GetOk("download"); ok {
		listMessagesPath += fmt.Sprintf("&download=%v", download)
	}
	if messageOffset, ok := d.GetOk("message_offset"); ok {
		listMessagesPath += fmt.Sprintf("&message_offset=%v", messageOffset)
	}
	if partition, ok := d.GetOk("partition"); ok {
		listMessagesPath += fmt.Sprintf("&partition=%v", partition)
	}
	if keyword, ok := d.GetOk("keyword"); ok {
		listMessagesPath += fmt.Sprintf("&keyword=%v", keyword)
	}

	return listMessagesPath
}
