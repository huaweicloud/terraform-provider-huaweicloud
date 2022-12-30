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

func DataSourceDmsRocketMQTopicGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDmsRocketMQTopicGroupRead,
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
			"topic": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the RocketMQ topic.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the list of RocketMQ consumer groups associated with the topic.`,
			},
		},
	}
}

func resourceDmsRocketMQTopicGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// getRocketmqTopicGroup: Query DMS RocketMQ consumer group list associated with topic
	var (
		getRocketmqTopicGroupHttpUrl = "v2/{project_id}/instances/{instance_id}/topics/{topic}/groups"
		getRocketmqTopicGroupProduct = "dms"
	)
	getRocketmqTopicGroupClient, err := config.NewServiceClient(getRocketmqTopicGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQTopicGroup Client: %s", err)
	}

	getRocketmqTopicGroupPath := getRocketmqTopicGroupClient.Endpoint + getRocketmqTopicGroupHttpUrl
	getRocketmqTopicGroupPath = strings.ReplaceAll(getRocketmqTopicGroupPath, "{project_id}",
		getRocketmqTopicGroupClient.ProjectID)
	getRocketmqTopicGroupPath = strings.ReplaceAll(getRocketmqTopicGroupPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))
	getRocketmqTopicGroupPath = strings.ReplaceAll(getRocketmqTopicGroupPath, "{topic}",
		fmt.Sprintf("%v", d.Get("topic")))

	getRocketmqTopicGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqTopicGroupResp, err := getRocketmqTopicGroupClient.Request("GET", getRocketmqTopicGroupPath,
		&getRocketmqTopicGroupOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DmsRocketMQTopicGroup")
	}

	getRocketmqTopicGroupRespBody, err := utils.FlattenResponse(getRocketmqTopicGroupResp)
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
		d.Set("groups", utils.PathSearch("groups", getRocketmqTopicGroupRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
