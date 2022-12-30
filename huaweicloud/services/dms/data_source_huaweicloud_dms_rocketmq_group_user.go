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

func DataSourceDmsRocketMQGroupUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDmsRocketMQGroupUserRead,
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
			"policies": {
				Type:        schema.TypeList,
				Elem:        DmsRocketMQGroupUserPolicySchema(),
				Computed:    true,
				Description: `Indicates the list of user associated with the consumer group.`,
			},
		},
	}
}

func DmsRocketMQGroupUserPolicySchema() *schema.Resource {
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

func resourceDmsRocketMQGroupUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// getRocketmqGroupUser: Query DMS RocketMQ users that have been granted permissions for a consumer group.
	var (
		getRocketmqGroupUserHttpUrl = "v2/{engine}/{project_id}/instances/{instance_id}/groups/{group_id}/accesspolicy"
		getRocketmqGroupUserProduct = "dms"
	)
	getRocketmqGroupUserClient, err := config.NewServiceClient(getRocketmqGroupUserProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQGroupUser Client: %s", err)
	}

	getRocketmqGroupUserPath := getRocketmqGroupUserClient.Endpoint + getRocketmqGroupUserHttpUrl
	getRocketmqGroupUserPath = strings.ReplaceAll(getRocketmqGroupUserPath, "{engine}", "reliability")
	getRocketmqGroupUserPath = strings.ReplaceAll(getRocketmqGroupUserPath, "{project_id}",
		getRocketmqGroupUserClient.ProjectID)
	getRocketmqGroupUserPath = strings.ReplaceAll(getRocketmqGroupUserPath, "{instance_id}",
		fmt.Sprintf("%v", d.Get("instance_id")))
	getRocketmqGroupUserPath = strings.ReplaceAll(getRocketmqGroupUserPath, "{group_id}",
		fmt.Sprintf("%v", d.Get("group")))

	getRocketmqGroupUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqGroupUserResp, err := getRocketmqGroupUserClient.Request("GET", getRocketmqGroupUserPath,
		&getRocketmqGroupUserOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DmsRocketMQGroupUser")
	}

	getRocketmqGroupUserRespBody, err := utils.FlattenResponse(getRocketmqGroupUserResp)
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
		d.Set("policies", flattenGetRocketmqGroupUserResponseBodyPolicy(getRocketmqGroupUserRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetRocketmqGroupUserResponseBodyPolicy(resp interface{}) []interface{} {
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
