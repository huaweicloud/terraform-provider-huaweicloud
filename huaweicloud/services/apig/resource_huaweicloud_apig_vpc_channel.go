package apig

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/channels"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

const (
	PROTOCOL_TCP   = "TCP"
	PROTOCOL_HTTP  = "HTTP"
	PROTOCOL_HTTPS = "HTTPS"
)

var (
	balanceStrategy = map[string]int{
		"WRR":         1,
		"WLC":         2,
		"SH":          3,
		"URI hashing": 4,
	}
	channelStatus = map[int]string{
		1: "Normal",
		2: "Abnormal",
	}
	memberType = map[string]string{
		"ECS": "ecs",
		"EIP": "ip",
	}
)

func ResourceApigVpcChannelV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigVpcChannelV2Create,
		Read:   resourceApigVpcChannelV2Read,
		Update: resourceApigVpcChannelV2Update,
		Delete: resourceApigVpcChannelV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceApigVpcChannelResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^([\u4e00-\u9fa5A-Za-z][\u4e00-\u9fa5A-Za-z-_0-9]{2,63})$"),
					"The name consists of 3 to 64 characters and only letters, digits, underscore (_), hyphens (-) "+
						"and chinese characters are allowed. The name must start with a letter or chinese character."),
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"member_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ECS",
				ValidateFunc: validation.StringInSlice([]string{
					"ECS", "EIP",
				}, false),
			},
			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "WRR", // The default value on golangsdk is WLC, but on the console it is WRR.
				ValidateFunc: validation.StringInSlice([]string{
					"WRR", "WLC", "SH", "URI hashing",
				}, false),
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  PROTOCOL_TCP,
				ValidateFunc: validation.StringInSlice([]string{
					PROTOCOL_TCP, PROTOCOL_HTTP, PROTOCOL_HTTPS,
				}, false),
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"healthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(2, 10),
			},
			"unhealthy_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validation.IntBetween(2, 10),
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validation.IntBetween(2, 30),
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				ValidateFunc: validation.IntBetween(5, 300),
			},
			"http_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"members": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 100),
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildApigVpcChannelHealthConfig(d *schema.ResourceData) (channels.VpcHealthConfig, error) {
	conf := channels.VpcHealthConfig{
		Protocol:          d.Get("protocol").(string),
		Path:              d.Get("path").(string),
		Port:              d.Get("port").(int),
		ThresholdNormal:   d.Get("healthy_threshold").(int),
		ThresholdAbnormal: d.Get("unhealthy_threshold").(int),
		Timeout:           d.Get("timeout").(int),
		TimeInterval:      d.Get("interval").(int),
	}
	// The parameter of http codes is required if protocol is set to http or https.
	if val := d.Get("protocol"); val.(string) == PROTOCOL_HTTP || val.(string) == PROTOCOL_HTTPS {
		if codes, ok := d.GetOk("http_code"); ok {
			conf.HttpCodes = codes.(string)
		} else {
			return conf, fmtp.Errorf("The http code cannot be empty if protocol is http or https")
		}
	}
	return conf, nil
}

func buildApigVpcChannelMembers(d *schema.ResourceData, config *config.Config) ([]channels.MemberInfo, error) {
	members := d.Get("members").(*schema.Set)
	mType := d.Get("member_type").(string)
	result := make([]channels.MemberInfo, members.Len())
	// Since the API requires that the name must be supported when the member type is 'ecs'.
	// It is necessary to query the ecs instance information from server to obtain the instance name.
	// We will cancel this unreasonable parameter configuration in the future.
	ecsClient, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return result, fmtp.Errorf("Error creating HuaweiCloud ECS v1 client: %s", err)
	}
	for i, v := range members.List() {
		member := v.(map[string]interface{})
		info := channels.MemberInfo{
			Weight: member["weight"].(int),
		}
		switch mType {
		case "ECS":
			{
				id, ok := member["id"]
				if !ok || id == "" {
					return result, fmtp.Errorf("The instance ID is missing, please check your input of members")
				}
				info.EcsId = id.(string)
				server, err := cloudservers.Get(ecsClient, id.(string)).Extract()
				if err != nil {
					return result, fmtp.Errorf("Error getting ECS instance from server by id: %s", err)
				}
				info.EcsName = server.Name
			}
		case "EIP":
			{
				addr, ok := member["ip_address"]
				if !ok || addr == "" {
					return result, fmtp.Errorf("The ip address of EIP is missing, please check your input of members")
				}
				info.Host = addr.(string)
			}
		default:
			return result, fmtp.Errorf("The member type is wrong, please check your input of members")
		}
		result[i] = info
	}
	return result, nil
}

func buildApigVpcChannelParameters(d *schema.ResourceData, config *config.Config) (channels.ChannelOpts, error) {
	opt := channels.ChannelOpts{
		Name: d.Get("name").(string),
		Port: d.Get("port").(int),
		Type: 2, // The type is required and type 1 (private network ELB channel) is to be deprecated.
	}
	// Member type, use 'ECS' and 'EIP' to save the stored value for better understanding.
	mType := d.Get("member_type").(string)
	if val, ok := memberType[mType]; ok {
		opt.MemberType = val
	} else {
		return opt, fmtp.Errorf("Wrong member type: %s", mType)
	}
	// Backend servers
	members, err := buildApigVpcChannelMembers(d, config)
	if err != nil {
		return opt, err
	}
	opt.Members = members
	// Healthy check config
	conf, err := buildApigVpcChannelHealthConfig(d)
	if err != nil {
		return opt, err
	}
	opt.VpcHealthConfig = conf
	// algorithm
	v, ok := balanceStrategy[d.Get("algorithm").(string)]
	if ok {
		opt.BalanceStrategy = v
	} else {
		return opt, fmtp.Errorf("The value of algorithm is invalid")
	}
	return opt, nil
}

func resourceApigVpcChannelV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	opts, err := buildApigVpcChannelParameters(d, config)
	if err != nil {
		return fmtp.Errorf("Unable to get the create option of the vpc channel: %s", err)
	}

	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	v, err := channels.Create(client, instanceId, opts).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "Error creating HuaweiCloud vpc channel")
	}
	d.SetId(v.Id)
	return resourceApigVpcChannelV2Read(d, meta)
}

func setApigVpcChannelMembers(d *schema.ResourceData, mType string, members []channels.MemberInfo) error {
	result := make([]map[string]interface{}, len(members))
	for i, v := range members {
		memberMap := map[string]interface{}{
			"weight": v.Weight,
		}
		switch mType {
		case "ecs":
			memberMap["id"] = v.EcsId
		case "ip":
			memberMap["ip_address"] = v.Host
		default:
			return fmtp.Errorf("Wrong members (%+v) with the type (%s)", v, mType)
		}
		result[i] = memberMap
	}
	return d.Set("members", result)
}

func setApigVpcChannelMemberType(d *schema.ResourceData, mType string) error {
	for k, v := range memberType {
		if v == mType {
			return d.Set("member_type", k)
		}
	}
	return fmtp.Errorf("The member type (%s) is not supported", mType)
}

func setApigVpcChannelAlgorithm(d *schema.ResourceData, algorithm int) error {
	for k, v := range balanceStrategy {
		if v == algorithm {
			return d.Set("algorithm", k)
		}
	}
	return fmtp.Errorf("The algorithm (%d) is not supported", algorithm)
}

func setApigVpcChannelParameters(d *schema.ResourceData, config *config.Config, resp channels.VpcChannel) error {
	status, ok := channelStatus[resp.Status]
	if !ok {
		return fmtp.Errorf("The response status is invalid")
	}
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("port", resp.Port),
		d.Set("protocol", strings.ToUpper(resp.VpcHealthConfig.Protocol)),
		d.Set("path", resp.VpcHealthConfig.Path),
		d.Set("healthy_threshold", resp.VpcHealthConfig.ThresholdNormal),
		d.Set("unhealthy_threshold", resp.VpcHealthConfig.ThresholdAbnormal),
		d.Set("timeout", resp.VpcHealthConfig.Timeout),
		d.Set("interval", resp.VpcHealthConfig.TimeInterval),
		d.Set("http_code", resp.VpcHealthConfig.HttpCodes),
		d.Set("create_time", resp.CreateTime),
		d.Set("status", status),
		setApigVpcChannelMemberType(d, resp.MemberType),
		setApigVpcChannelAlgorithm(d, resp.BalanceStrategy),
		setApigVpcChannelMembers(d, resp.MemberType, resp.Members),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func resourceApigVpcChannelV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := channels.Get(client, instanceId, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, fmt.Sprintf("error getting vpc channel (%s) form server", d.Id()))
	}
	return setApigVpcChannelParameters(d, config, *resp)
}

func resourceApigVpcChannelV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}

	opt, err := buildApigVpcChannelParameters(d, config)
	if err != nil {
		return fmtp.Errorf("Unable to get the update option of the vpc channel: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	_, err = channels.Update(client, instanceId, d.Id(), opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating vpc channel: %s", err)
	}
	return resourceApigVpcChannelV2Read(d, meta)
}

func resourceApigVpcChannelV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	if err = channels.Delete(client, instanceId, d.Id()).ExtractErr(); err != nil {
		return fmtp.Errorf("Unable to delete the vpc channel (%s): %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

// The ID cannot find on console, so we need to import by channel name.
func resourceApigVpcChannelResourceImportState(d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmtp.Errorf("Invalid format specified for import id, must be <instance_id>/<channel name>")
	}
	instanceId := parts[0]
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	name := parts[1]
	opt := channels.ListOpts{
		Name: name,
	}
	pages, err := channels.List(client, instanceId, opt).AllPages()
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error retrieving channels: %s", err)
	}
	resp, err := channels.ExtractChannels(pages)
	if len(resp) < 1 {
		return []*schema.ResourceData{d}, fmtp.Errorf("Unable to find the channel (%s) form server: %s", name, err)
	}
	d.SetId(resp[0].Id)
	d.Set("instance_id", instanceId)
	return []*schema.ResourceData{d}, nil
}
