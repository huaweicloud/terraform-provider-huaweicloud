package rocketmq

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

// @API RocketMQ GET /v2/{project_id}/instances/{instance_id}/users
func DataSourceDmsRocketMQUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDmsRocketMQUsersRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"white_remote_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_topic_perm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_group_perm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     usersSchema(),
			},
		},
	}
}

func usersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"white_remote_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"admin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default_topic_perm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_group_perm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic_perms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topicOrGroupPermsSchema(),
			},
			"group_perms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topicOrGroupPermsSchema(),
			},
		},
	}
	return &sc
}

func topicOrGroupPermsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"perm": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDmsRocketMQUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getRocketmqUsersHttpUrl = "v2/{project_id}/instances/{instance_id}/users"
		getRocketmqUsersProduct = "dmsv2"
	)
	getRocketmqUsersClient, err := cfg.NewServiceClient(getRocketmqUsersProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	getRocketmqUsersPath := getRocketmqUsersClient.Endpoint + getRocketmqUsersHttpUrl
	getRocketmqUsersPath = strings.ReplaceAll(getRocketmqUsersPath, "{project_id}", getRocketmqUsersClient.ProjectID)
	getRocketmqUsersPath = strings.ReplaceAll(getRocketmqUsersPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	listUsersResp, err := pagination.ListAllItems(
		getRocketmqUsersClient,
		"offset",
		getRocketmqUsersPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DMS rocketMQ users")
	}

	listUsersRespJson, err := json.Marshal(listUsersResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listUsersRespBody interface{}
	err = json.Unmarshal(listUsersRespJson, &listUsersRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("users", flattenListUsersBody(filterUsers(d, listUsersRespBody))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListUsersBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"access_key":           utils.PathSearch("access_key", v, nil),
			"white_remote_address": utils.PathSearch("white_remote_address", v, nil),
			"admin":                utils.PathSearch("admin", v, nil),
			"default_topic_perm":   utils.PathSearch("default_topic_perm", v, nil),
			"default_group_perm":   utils.PathSearch("default_group_perm", v, nil),
			"topic_perms":          utils.PathSearch("topic_perms", v, nil),
			"group_perms":          utils.PathSearch("group_perms", v, nil),
		})
	}
	return rst
}

func filterUsers(d *schema.ResourceData, resp interface{}) []interface{} {
	userJson := utils.PathSearch("users", resp, make([]interface{}, 0))
	userArray := userJson.([]interface{})
	if len(userArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(userArray))

	rawAccessKey, rawAccessKeyOK := d.GetOk("access_key")
	rawWhiteRemoteAddress, rawWhiteRemoteAddressOK := d.GetOk("white_remote_address")
	rawAdmin := d.Get("admin").(bool)
	rawDefaultTopicPerm, rawDefaultTopicPermOK := d.GetOk("default_topic_perm")
	rawDefaultGroupPerm, rawDefaultGroupPermOK := d.GetOk("default_group_perm")

	for _, user := range userArray {
		accessKey := utils.PathSearch("access_key", user, nil)
		whiteRemoteAddress := utils.PathSearch("white_remote_address", user, nil)
		admin := utils.PathSearch("admin", user, false).(bool)
		defaultTopicPerm := utils.PathSearch("default_topic_perm", user, nil)
		defaultGroupPerm := utils.PathSearch("default_group_perm", user, nil)
		if rawAccessKeyOK && rawAccessKey != accessKey {
			continue
		}
		if rawWhiteRemoteAddressOK && rawWhiteRemoteAddress != whiteRemoteAddress {
			continue
		}
		if (rawAdmin && !admin) || (!rawAdmin && admin) {
			continue
		}
		if rawDefaultTopicPermOK && rawDefaultTopicPerm != defaultTopicPerm {
			continue
		}
		if rawDefaultGroupPermOK && rawDefaultGroupPerm != defaultGroupPerm {
			continue
		}
		result = append(result, user)
	}

	return result
}
