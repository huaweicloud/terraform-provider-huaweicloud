package dws

import (
	"context"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dws/v1/cluster"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
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

			"tags": common.TagsForceNewSchema(),

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
		return fmtp.DiagErrorf("error creating DWS v1 client, err=%s", err)
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
		Tags:                utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}

	if obj, ok := d.GetOk("number_of_cn"); ok {
		numberOfCn := obj.(int)
		opts.NumberOfCn = &numberOfCn
	}

	publicIPProp, err := expandDwsClusterPublicIP(d)
	if err != nil {
		return diag.FromErr(err)
	}
	opts.PublicIp = publicIPProp

	rst, createErr := cluster.Create(client, opts)
	if createErr != nil {
		return fmtp.DiagErrorf("Error creating DWS Cluster: %s", createErr)
	}

	clusterId := rst.Cluster.Id

	d.SetId(clusterId)

	checkCreateErr := checkClusterCreateResult(ctx, client, clusterId, d.Timeout(schema.TimeoutCreate))
	if checkCreateErr != nil {
		return diag.FromErr(checkCreateErr)
	}

	return resourceDwsClusterRead(ctx, d, meta)
}

func resourceDwsClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DwsV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("error creating DWS v1 client, err=%s", err)
	}

	clusterDetail, dErr := cluster.Get(client, d.Id())
	if dErr != nil {
		return fmtp.DiagErrorf("Error query DwsCluster %q:%s", d.Id(), dErr)
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
		setPublicIpToState(d, clusterDetail.PublicIp),
		d.Set("enterprise_project_id", clusterDetail.EnterpriseProjectId),
		d.Set("tags", utils.TagsToMap(clusterDetail.Tags)),
		d.Set("created", clusterDetail.Created),
		setEndpointsToState(d, clusterDetail.Endpoints),
		setPublicEndpointsToState(d, clusterDetail.PublicEndpoints),
		d.Set("recent_event", clusterDetail.RecentEvent),
		d.Set("status", clusterDetail.Status),
		d.Set("sub_status", clusterDetail.SubStatus),
		d.Set("task_status", clusterDetail.TaskStatus),
		d.Set("updated", clusterDetail.Updated),
		d.Set("version", clusterDetail.Version),
		d.Set("private_ip", clusterDetail.PrivateIp),
	)
	if setSdErr := mErr.ErrorOrNil(); setSdErr != nil {
		return fmtp.DiagErrorf("Error setting vault fields: %s", setSdErr)
	}

	return nil
}

func setPublicIpToState(d *schema.ResourceData, publicIp *cluster.PublicIp) error {
	if publicIp == nil {
		return nil
	}
	value := []interface{}{map[string]string{
		"eip_id":           publicIp.EipID,
		"public_bind_type": publicIp.PublicBindType,
	}}

	return d.Set("public_ip", value)
}

func resourceDwsClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DwsV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("error creating DWS v1 client, err=%s", err)
	}

	clusterId := d.Id()
	errResult := cluster.Delete(client, clusterId)
	if errResult.Err != nil {
		return fmtp.DiagErrorf("Delete DWS Cluster failed. %s", errResult.Err)
	}

	errCheckRt := checkClusterDeleteResult(ctx, client, clusterId, d.Timeout(schema.TimeoutDelete))
	if errCheckRt != nil {
		return fmtp.DiagErrorf("Failed to check the result of deletion %s", errCheckRt)
	}
	d.SetId("")
	return nil
}

func resourceDwsClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DwsV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("error creating DWS v1 client, err=%s", err)
	}

	clusterId := d.Id()
	//check cluster state is available before update
	checkErr := checkAndWaitClusterStateAvailable(ctx, client, clusterId, true, d.Timeout(schema.TimeoutUpdate))
	if checkErr != nil {
		return fmtp.DiagErrorf("cluster state is not available to update. cluster_id:%s,error:%s", clusterId, checkErr)
	}

	// extend cluster
	if d.HasChange("number_of_node") {
		oldValue, newValue := d.GetChange("number_of_node")
		num := newValue.(int) - oldValue.(int)
		_, extendErr := cluster.Resize(client, clusterId, num)
		if extendErr != nil {
			return fmtp.DiagErrorf("Extend DWS cluster failed.cluster_id=%s,error=%s", clusterId, extendErr)
		}
		checkErr = checkAndWaitClusterStateAvailable(ctx, client, clusterId, true, d.Timeout(schema.TimeoutUpdate))
		if checkErr != nil {
			return fmtp.DiagErrorf("Extends node failed. cluster_id:%s,error:%s", clusterId, checkErr)
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
			return fmtp.DiagErrorf("reset password of DWS cluster failed. cluster_id=%s,error:%s", clusterId, rErr)
		}

		checkErr = checkAndWaitClusterStateAvailable(ctx, client, clusterId, false, d.Timeout(schema.TimeoutUpdate))
		if checkErr != nil {
			return fmtp.DiagErrorf("reset password of dws cluster failed. cluster_id:%s,error:%s", clusterId, checkErr)
		}

	}

	return resourceDwsClusterRead(ctx, d, meta)
}

func setEndpointsToState(d *schema.ResourceData, endpoints []cluster.Endpoints) error {
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

	return d.Set("endpoints", result)
}

func setPublicEndpointsToState(d *schema.ResourceData, endpoints []cluster.PublicEndpoints) error {
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

	return d.Set("public_endpoints", result)
}

func expandDwsClusterPublicIP(d *schema.ResourceData) (*cluster.PublicIpOpts, error) {
	var rst cluster.PublicIpOpts
	if obj, ok := d.GetOk("public_ip.0.public_bind_type"); ok {
		publicBindType := obj.(string)

		switch publicBindType {
		case cluster.PublicBindTypeBindExisting:
			if obj, ok := d.GetOk("public_ip.0.eip_id"); ok {
				rst.EipID = obj.(string)
				rst.PublicBindType = publicBindType
			} else {
				return nil, fmtp.Errorf("Illegal parameter:When public_bind_type is equal '%s', eip_id is required",
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
		Pending: []string{"CREATING", "Pending"},
		Target:  []string{"AVAILABLE"},
		Refresh: func() (interface{}, string, error) {
			resp, err := cluster.Get(client, clusterId)
			if err != nil {
				return nil, "failed", err
			}

			if resp.FailedReasons != nil && resp.FailedReasons.ErrorCode != "" {
				return nil, "failed", fmtp.Errorf("create DWS failed. error_code: %s, error_msg: %s",
					resp.FailedReasons.ErrorCode, resp.FailedReasons.ErrorMsg)
			}
			return resp, resp.Status, err
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("error waiting for DWS (%s) to be created: %s", clusterId, err)
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
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return true, "Done", nil
				}
				return nil, "", nil
			}
			return true, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("error waiting for DWS (%s) to be delete: %s", clusterId, err)
	}
	return nil
}

// when extend=true: if TaskStatus = RESIZE_FAILURE ,return error; else just check cluster is no task running
func parseClusterStatus(detail *cluster.ClusterDetail, extend bool) (bool, error) {
	//actions --- the behaviors on a cluster
	if len(detail.ActionProgress) > 0 {
		return false, nil
	}

	if detail.Status != "AVAILABLE" {
		return false, nil
	}

	if extend && detail.TaskStatus == "RESIZE_FAILURE" {
		return false, fmtp.Errorf("RESIZE_FAILURE")
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
				return nil, "failed", fmtp.Errorf("error_code: %s, error_msg: %s", resp.FailedReasons.ErrorCode,
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
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("waiting for DWS (%s) to finish task failed: %s", clusterId, err)
	}
	return nil
}
