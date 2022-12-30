package dms

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDmsRocketMQGroupTopics() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDmsRocketMQGroupTopicsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RocketMQ instance.`,
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the RocketMQ consumer group.`,
			},
			"topics": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the list of RocketMQ topics associated with the consumer group.`,
			},
		},
	}
}

func resourceDmsRocketMQGroupTopicsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// getRocketmqGroupTopic: Query DMS RocketMQ topics list associated with consumer group
	var (
		getRocketmqGroupTopicHttpUrl = "v2/{project_id}/instances/{instance_id}/groups/{group}/topics"
		getRocketmqGroupTopicProduct = "dms"
	)
	getRocketmqGroupTopicClient, err := config.NewServiceClient(getRocketmqGroupTopicProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQGroupTopics Client: %s", err)
	}

	getRocketmqGroupTopicPath := getRocketmqGroupTopicClient.Endpoint + getRocketmqGroupTopicHttpUrl
	getRocketmqGroupTopicPath = strings.ReplaceAll(getRocketmqGroupTopicPath, "{project_id}",
		getRocketmqGroupTopicClient.ProjectID)
	getRocketmqGroupTopicPath = strings.ReplaceAll(getRocketmqGroupTopicPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))
	getRocketmqGroupTopicPath = strings.ReplaceAll(getRocketmqGroupTopicPath, "{group}",
		fmt.Sprintf("%v", d.Get("group")))

	getRocketmqGroupTopicOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqGroupTopicResp, err := getRocketmqGroupTopicClient.Request("GET", getRocketmqGroupTopicPath,
		&getRocketmqGroupTopicOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DmsRocketMQGroupTopics")
	}

	getRocketmqGroupTopicRespBody, err := utils.FlattenResponse(getRocketmqGroupTopicResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("topics", utils.PathSearch("topics", getRocketmqGroupTopicRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
