package dws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dws/v1/cluster"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceDwsCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDwsClusterCreate,
		ReadContext:   resourceDwsClusterRead,
		DeleteContext: resourceDwsClusterDelete,
		UpdateContext: resourceDwsClusterUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"node_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"number_of_node": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"user_pwd": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"public_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eip_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
						"public_bind_type": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice(
								[]string{
									cluster.PublicBindTypeBindExisting,
									cluster.PublicBindTypeAuto,
									cluster.PublicBindTypeNotUse,
								},
								false),
						},
					},
				},
			},

			"number_of_cn": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"tags": common.TagsSchema(),

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connect_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"jdbc_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"public_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"jdbc_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_connect_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"recent_event": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sub_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"task_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDwsClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DwsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DWS v1 client, err=%s", err)
	}

	opts := &cluster.CreateOpts{
		Name:                d.Get("name").(string),
		AvailabilityZone:    d.Get("availability_zone").(string),
		SubnetID:            d.Get("network_id").(string),
		NodeType:            d.Get("node_type").(string),
		UserName:            d.Get("user_name").(string),
		UserPwd:             d.Get("user_pwd").(string),
		NumberOfNode:        d.Get("number_of_node").(int),
		Port:                d.Get("port").(int),
		SecurityGroupID:     d.Get("security_group_id").(string),
		VpcID:               d.Get("vpc_id").(string),
		EnterpriseProjectId: config.GetEnterpriseProjectID(d),
	}

	if obj, ok := d.GetOk("number_of_cn"); ok {
		opts.NumberOfCn = utils.IntIgnoreEmpty(obj.(int))
	}

	publicIPProp, err := buildDwsClusterPublicIP(d)
	if err != nil {
		return diag.FromErr(err)
	}
	opts.PublicIp = publicIPProp

	rst, createErr := cluster.Create(client, opts)
	if createErr != nil {
		return diag.Errorf("error creating DWS Cluster: %s", createErr)
	}

	clusterId := rst.Cluster.Id

	d.SetId(clusterId)

	checkCreateErr := checkClusterCreateResult(ctx, client, clusterId, d.Timeout(schema.TimeoutCreate))
	if checkCreateErr != nil {
		return diag.FromErr(checkCreateErr)
	}

	// tags
	if v, ok := d.GetOk("tags"); ok {
		tagRaw := v.(map[string]interface{})
		if len(tagRaw) > 0 {
			if err = addDwsClusterTags(client, clusterId, utils.ExpandResourceTags(tagRaw)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceDwsClusterRead(ctx, d, meta)
}

func resourceDwsClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DwsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DWS v1 client, err=%s", err)
	}

	clusterDetail, err := cluster.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, parseDwsClusterNotFoundError(err), "DWS cluster")
	}

	mErr := multierror.Append(
		d.Set("name", clusterDetail.Name),
		d.Set("network_id", clusterDetail.SubnetID),
		d.Set("node_type", clusterDetail.NodeType),
		d.Set("number_of_node", clusterDetail.NumberOfNode),
		d.Set("security_group_id", clusterDetail.SecurityGroupID),
		d.Set("user_name", clusterDetail.UserName),
		d.Set("vpc_id", clusterDetail.VpcID),
		d.Set("availability_zone", clusterDetail.AvailabilityZone),
		d.Set("port", clusterDetail.Port),
		d.Set("public_ip", flattenPublicIpToState(d, clusterDetail.PublicIp)),
		d.Set("enterprise_project_id", clusterDetail.EnterpriseProjectId),
		d.Set("tags", utils.TagsToMap(clusterDetail.Tags)),
		d.Set("created", clusterDetail.Created),
		d.Set("endpoints", flattenEndpointsToState(d, clusterDetail.Endpoints)),
		d.Set("public_endpoints", flattenPublicEndpointsToState(d, clusterDetail.PublicEndpoints)),
		d.Set("recent_event", clusterDetail.RecentEvent),
		d.Set("status", clusterDetail.Status),
		d.Set("sub_status", clusterDetail.SubStatus),
		d.Set("task_status", clusterDetail.TaskStatus),
		d.Set("updated", clusterDetail.Updated),
		d.Set("version", clusterDetail.Version),
		d.Set("private_ip", clusterDetail.PrivateIp),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPublicIpToState(d *schema.ResourceData, publicIp *cluster.PublicIp) []interface{} {
	if publicIp == nil {
		return nil
	}
	return []interface{}{map[string]string{
		"eip_id":           publicIp.EipID,
		"public_bind_type": publicIp.PublicBindType,
	}}
}

func resourceDwsClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DwsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DWS v1 client, err=%s", err)
	}

	clusterId := d.Id()
	errResult := cluster.Delete(client, clusterId)
	if errResult.Err != nil {
		return diag.Errorf("deleting DWS Cluster failed. %s", errResult.Err)
	}

	errCheckRt := checkClusterDeleteResult(ctx, client, clusterId, d.Timeout(schema.TimeoutDelete))
	if errCheckRt != nil {
		return diag.Errorf("failed to check the result of deletion %s", errCheckRt)
	}
	return nil
}

func resourceDwsClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DwsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DWS v1 client, err=%s", err)
	}

	clusterId := d.Id()
	// check cluster state is available before update
	checkErr := checkAndWaitClusterStateAvailable(ctx, client, clusterId, true, d.Timeout(schema.TimeoutUpdate))
	if checkErr != nil {
		return diag.Errorf("cluster state is not available to update. cluster_id:%s,error:%s", clusterId, checkErr)
	}

	// extend cluster
	if d.HasChange("number_of_node") {
		oldValue, newValue := d.GetChange("number_of_node")
		num := newValue.(int) - oldValue.(int)
		_, extendErr := cluster.Resize(client, clusterId, num)
		if extendErr != nil {
			return diag.Errorf("extend DWS cluster failed.cluster_id:%s, error:%s", clusterId, extendErr)
		}
		checkErr = checkAndWaitClusterStateAvailable(ctx, client, clusterId, true, d.Timeout(schema.TimeoutUpdate))
		if checkErr != nil {
			return diag.Errorf("extend DWS cluster failed. cluster_id:%s, error:%s", clusterId, checkErr)
		}
	}

	// change pwd
	if d.HasChange("user_pwd") {
		newValue := d.Get("user_pwd")

		var opts = cluster.ResetPasswordOpts{
			NewPassword: newValue.(string),
		}

		_, rErr := cluster.ResetPassword(client, clusterId, opts)
		if rErr != nil {
			return diag.Errorf("reset password of DWS cluster failed. cluster_id:%s, error:%s", clusterId, rErr)
		}

		checkErr = checkAndWaitClusterStateAvailable(ctx, client, clusterId, false, d.Timeout(schema.TimeoutUpdate))
		if checkErr != nil {
			return diag.Errorf("reset password of dws cluster failed. cluster_id:%s, error:%s", clusterId, checkErr)
		}
	}

	// change tag
	if d.HasChange("tags") {
		err = updateDwsTags(client, d, clusterId)
		if err != nil {
			return diag.Errorf("error updating tags of DWS cluster:%s, err:%s", clusterId, err)
		}
	}

	return resourceDwsClusterRead(ctx, d, meta)
}

func addDwsClusterTags(client *golangsdk.ServiceClient, clusterId string, tags []tags.ResourceTag) error {
	var (
		addTagsHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/tags/batch-create"
	)

	addTagsPath := client.Endpoint + addTagsHttpUrl
	addTagsPath = strings.ReplaceAll(addTagsPath, "{project_id}", client.ProjectID)
	addTagsPath = strings.ReplaceAll(addTagsPath, "{cluster_id}", clusterId)

	addTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	addTagsOpt.JSONBody = map[string]interface{}{
		"tags": tags,
	}
	_, err := client.Request("POST", addTagsPath, &addTagsOpt)
	if err != nil {
		return fmt.Errorf("error setting tags of DWS cluster: %s", err)
	}

	return nil
}

func deleteDwsClusterTags(client *golangsdk.ServiceClient, clusterId string, tags []tags.ResourceTag) error {
	var (
		deleteTagsHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/tags/batch-delete"
	)

	deleteTagsPath := client.Endpoint + deleteTagsHttpUrl
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{project_id}", client.ProjectID)
	deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{cluster_id}", clusterId)

	deleteTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	deleteTagsOpt.JSONBody = map[string]interface{}{
		"tags": tags,
	}
	_, err := client.Request("POST", deleteTagsPath, &deleteTagsOpt)
	if err != nil {
		return fmt.Errorf("error deleting tags of DWS cluster: %s", err)
	}

	return nil
}

func updateDwsTags(client *golangsdk.ServiceClient, d *schema.ResourceData, id string) error {
	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	// remove old tags
	if len(oMap) > 0 {
		taglist := utils.ExpandResourceTags(oMap)
		err := deleteDwsClusterTags(client, id, taglist)
		if err != nil {
			return err
		}
	}

	// set new tags
	if len(nMap) > 0 {
		taglist := utils.ExpandResourceTags(nMap)
		err := addDwsClusterTags(client, id, taglist)
		if err != nil {
			return err
		}
	}

	return nil
}

func flattenEndpointsToState(d *schema.ResourceData, endpoints []cluster.Endpoints) []interface{} {
	if len(endpoints) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(endpoints))
	for _, endpoint := range endpoints {
		transformed := map[string]interface{}{
			"connect_info": endpoint.ConnectInfo,
			"jdbc_url":     endpoint.JdbcUrl,
		}
		result = append(result, transformed)
	}

	return result
}

func flattenPublicEndpointsToState(d *schema.ResourceData, endpoints []cluster.PublicEndpoints) []interface{} {
	if len(endpoints) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(endpoints))
	for _, endpoint := range endpoints {
		transformed := map[string]interface{}{
			"public_connect_info": endpoint.PublicConnectInfo,
			"jdbc_url":            endpoint.JdbcUrl,
		}
		result = append(result, transformed)
	}

	return result
}

func buildDwsClusterPublicIP(d *schema.ResourceData) (*cluster.PublicIpOpts, error) {
	var rst cluster.PublicIpOpts
	if obj, ok := d.GetOk("public_ip.0.public_bind_type"); ok {
		publicBindType := obj.(string)

		switch publicBindType {
		case cluster.PublicBindTypeBindExisting:
			if obj, ok := d.GetOk("public_ip.0.eip_id"); ok {
				rst.EipID = obj.(string)
				rst.PublicBindType = publicBindType
			} else {
				return nil, fmt.Errorf("illegal parameter:When public_bind_type is equal '%s', eip_id is required",
					publicBindType)
			}

		default:
			rst.PublicBindType = publicBindType
		}

	}
	return &rst, nil
}

func checkClusterCreateResult(ctx context.Context, client *golangsdk.ServiceClient, clusterId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"CREATING"},
		Target:  []string{"AVAILABLE"},
		Refresh: func() (interface{}, string, error) {
			resp, err := cluster.Get(client, clusterId)
			if err != nil {
				return nil, "failed", err
			}

			if resp.FailedReasons != nil && resp.FailedReasons.ErrorCode != "" {
				return nil, "failed", fmt.Errorf("create DWS failed. error_code: %s, error_msg: %s",
					resp.FailedReasons.ErrorCode, resp.FailedReasons.ErrorMsg)
			}
			return resp, resp.Status, err
		},
		Timeout:      timeout,
		PollInterval: 20 * time.Second,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DWS (%s) to be created: %s", clusterId, err)
	}
	return nil
}

func checkClusterDeleteResult(ctx context.Context, client *golangsdk.ServiceClient, clusterId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			_, err := cluster.Get(client, clusterId)
			if err != nil {
				if _, ok := parseDwsClusterNotFoundError(err).(golangsdk.ErrDefault404); ok {
					return true, "Done", nil
				}
				return nil, "failed", err
			}
			return true, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 20 * time.Second,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DWS (%s) to be delete: %s", clusterId, err)
	}
	return nil
}

// when extend=true: if TaskStatus = RESIZE_FAILURE ,return error; else just check cluster is no task running
func parseClusterStatus(detail *cluster.ClusterDetail, extend bool) (bool, error) {
	// actions --- the behaviors on a cluster
	if len(detail.ActionProgress) > 0 {
		return false, nil
	}

	if detail.Status != "AVAILABLE" {
		return false, nil
	}

	if extend && detail.TaskStatus == "RESIZE_FAILURE" {
		return false, fmt.Errorf("RESIZE_FAILURE")
	}

	if detail.TaskStatus != "" {
		return false, nil
	}

	if detail.SubStatus != "NORMAL" {
		return false, nil
	}

	return true, nil
}

// waiting cluster state is available to submit new operate
func checkAndWaitClusterStateAvailable(ctx context.Context, client *golangsdk.ServiceClient, clusterId string,
	isExtendTask bool, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			resp, err := cluster.Get(client, clusterId)
			if err != nil {
				return nil, "failed", err
			}

			if resp.FailedReasons != nil && resp.FailedReasons.ErrorCode != "" {
				return nil, "failed", fmt.Errorf("error_code: %s, error_msg: %s", resp.FailedReasons.ErrorCode,
					resp.FailedReasons.ErrorMsg)
			}

			cState, cErr := parseClusterStatus(resp, isExtendTask)
			if cErr != nil {
				return nil, "failed", cErr
			}
			if cState {
				return resp, "Done", nil
			}
			return resp, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 20 * time.Second,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for DWS (%s) to finish task failed: %s", clusterId, err)
	}
	return nil
}

func parseDwsClusterNotFoundError(respErr error) error {
	var apiErr cluster.FailInfo
	if errCode, ok := respErr.(golangsdk.ErrDefault401); ok && errCode.Body != nil {
		pErr := json.Unmarshal(errCode.Body, &apiErr)
		if pErr == nil && apiErr.ErrorCode == "DWS.0047" {
			return golangsdk.ErrDefault404{}
		}
	}
	return respErr
}
