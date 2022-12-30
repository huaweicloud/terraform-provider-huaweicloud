package dms

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDmsRocketMQTopicUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDmsRocketMQTopicUserRead,
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
				Required:    true,
				Description: `Specifies the name of the RocketMQ topic.`,
			},
			"policies": {
				Type:        schema.TypeList,
				Elem:        DmsRocketMQTopicUserPolicySchema(),
				Computed:    true,
				Description: `Indicates the list of user associated with the topic.`,
			},
		},
	}
}

func DmsRocketMQTopicUserPolicySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the access key of the user.`,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the secret key of the user.`,
			},
			"white_remote_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the IP address whitelist.`,
			},
			"admin": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the user is an administrator.`,
			},
			"perm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the permissions of the user.`,
			},
		},
	}
	return &sc
}

func resourceDmsRocketMQTopicUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// getRocketmqTopicUser: Query DMS RocketMQ users that have been granted permissions for a topic.
	var (
		getRocketmqTopicUserHttpUrl = "v2/{project_id}/instances/{instance_id}/topics/{topic}/accesspolicy"
		getRocketmqTopicUserProduct = "dms"
	)
	getRocketmqTopicUserClient, err := config.NewServiceClient(getRocketmqTopicUserProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQTopicUser Client: %s", err)
	}

	getRocketmqTopicUserPath := getRocketmqTopicUserClient.Endpoint + getRocketmqTopicUserHttpUrl
	getRocketmqTopicUserPath = strings.ReplaceAll(getRocketmqTopicUserPath, "{project_id}",
		getRocketmqTopicUserClient.ProjectID)
	getRocketmqTopicUserPath = strings.ReplaceAll(getRocketmqTopicUserPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))
	getRocketmqTopicUserPath = strings.ReplaceAll(getRocketmqTopicUserPath, "{topic}",
		fmt.Sprintf("%v", d.Get("topic")))

	getRocketmqTopicUserResp, err := pagination.ListAllItems(
		getRocketmqTopicUserClient,
		"offset",
		getRocketmqTopicUserPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DmsRocketMQTopicUser")
	}

	getRocketmqTopicUserRespJson, err := json.Marshal(getRocketmqTopicUserResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getRocketmqTopicUserRespBody interface{}
	err = json.Unmarshal(getRocketmqTopicUserRespJson, &getRocketmqTopicUserRespBody)
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
		d.Set("policies", flattenGetRocketmqTopicUserResponseBodyPolicy(getRocketmqTopicUserRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetRocketmqTopicUserResponseBodyPolicy(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("policies", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"access_key":           utils.PathSearch("access_key", v, nil),
			"secret_key":           utils.PathSearch("secret_key", v, nil),
			"white_remote_address": utils.PathSearch("white_remote_address", v, nil),
			"admin":                utils.PathSearch("admin", v, nil),
			"perm":                 utils.PathSearch("perm", v, nil),
		})
	}
	return rst
}
