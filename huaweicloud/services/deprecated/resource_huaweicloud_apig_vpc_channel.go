package deprecated

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/apigw/deprecated/dedicated/v2/channels"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

type (
	MemberType    string
	AlgorithmType string
	ProtocolType  string
	ChannelStatus int
)

const (
	MemberTypeEcs MemberType = "ECS"
	MemberTypeEip MemberType = "EIP"

	AlgorithmTypeWrr AlgorithmType = "WRR"
	AlgorithmTypeWlc AlgorithmType = "WLC"
	AlgorithmTypeSh  AlgorithmType = "SH"
	AlgorithmTypeUri AlgorithmType = "URI hashing"

	ProtocolTypeTcp   ProtocolType = "TCP"
	ProtocolTypeHttp  ProtocolType = "HTTP"
	ProtocolTypeHttps ProtocolType = "HTTPS"
	ProtocolTypeBoth  ProtocolType = "BOTH"

	ChannelStatusNormal   ChannelStatus = 1
	ChannelStatusAbnormal ChannelStatus = 2

	ProtocolTypeTCP   ProtocolType = "TCP"
	ProtocolTypeHTTP  ProtocolType = "HTTP"
	ProtocolTypeHTTPS ProtocolType = "HTTPS"
)

var (
	memberType = map[MemberType]string{
		MemberTypeEcs: "ecs",
		MemberTypeEip: "ip",
	}

	balanceStrategy = map[AlgorithmType]int{
		AlgorithmTypeWrr: 1,
		AlgorithmTypeWlc: 2,
		AlgorithmTypeSh:  3,
		AlgorithmTypeUri: 4,
	}

	channelStatus = map[int]string{
		1: "Normal",
		2: "Abnormal",
	}
)

func ResourceApigVpcChannelV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcChannelCreate,
		ReadContext:   resourceVpcChannelRead,
		UpdateContext: resourceVpcChannelUpdate,
		DeleteContext: resourceVpcChannelDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceVpcChannelResourceImportState,
		},
		DeprecationMessage: "VPC channel has been deprecated.",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the dedicated instance is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the VPC channel belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the VPC channel.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The host port of the VPC channel.",
			},
			"members": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The IP address of the backend server.",
						},
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the backend server.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "The weight of current backend server.",
						},
					},
				},
				Description: "The configuration of the backend servers that bind the VPC channel.",
			},
			"member_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(MemberTypeEcs),
				ValidateFunc: validation.StringInSlice([]string{
					string(MemberTypeEcs),
					string(MemberTypeEip),
				}, false),
				Description: "The member type of the VPC channel.",
			},
			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				// The default value on golangsdk is WLC, but on the console it is WRR.
				Default: string(AlgorithmTypeWrr),
				ValidateFunc: validation.StringInSlice([]string{
					string(AlgorithmTypeWrr),
					string(AlgorithmTypeWlc),
					string(AlgorithmTypeSh),
					string(AlgorithmTypeUri),
				}, false),
				Description: "The distribution algorithm.",
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(ProtocolTypeTCP),
				ValidateFunc: validation.StringInSlice([]string{
					string(ProtocolTypeTCP),
					string(ProtocolTypeHTTP),
					string(ProtocolTypeHTTPS),
				}, false),
				Description: "The rotocol for performing health checks on backend servers in the VPC channel.",
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The destination path for health checks.",
			},
			"healthy_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
				Description: "The the healthy threshold, which refers to the number of consecutive successful " +
					"checks required for a backend server to be considered healthy.",
			},
			"unhealthy_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
				Description: "The unhealthy threshold, which refers to the number of consecutive failed checks " +
					"required for a backend server to be considered unhealthy.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The timeout for determining whether a health check fails, in second.",
			},
			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "The interval between consecutive checks, in second.",
			},
			"http_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The response codes for determining a successful HTTP response.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the VPC channel was created.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the VPC channel.`,
			},
		},
	}
}

func buildVpcChannelHealthConfig(d *schema.ResourceData) (channels.VpcHealthConfig, error) {
	protocol := d.Get("protocol").(string)
	result := channels.VpcHealthConfig{
		Protocol:          protocol,
		Path:              d.Get("path").(string),
		Port:              d.Get("port").(int),
		ThresholdNormal:   d.Get("healthy_threshold").(int),
		ThresholdAbnormal: d.Get("unhealthy_threshold").(int),
		Timeout:           d.Get("timeout").(int),
		TimeInterval:      d.Get("interval").(int),
	}
	// The parameter of HTTP codes is required if protocol is set to 'HTTP' or 'HTTPS'.
	if ProtocolType(protocol) == ProtocolTypeHTTP || ProtocolType(protocol) == ProtocolTypeHTTPS {
		codes, ok := d.GetOk("http_code")
		if !ok {
			return result, fmt.Errorf("the HTTP code cannot be empty if the protocol is 'HTTP' or 'HTTPS'")
		}
		result.HttpCodes = codes.(string)
	}
	return result, nil
}

func buildVpcChannelMembers(d *schema.ResourceData, config *config.Config) ([]channels.MemberInfo, error) {
	var (
		members = d.Get("members").(*schema.Set)
		mType   = MemberType(d.Get("member_type").(string))
		result  = make([]channels.MemberInfo, members.Len())
	)

	// Since the API requires that the name must be supported when the member type is 'ecs'.
	// It is necessary to query the ecs instance information from server to obtain the instance name.
	// We will cancel this unreasonable parameter configuration in the future.
	ecsClient, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return result, fmt.Errorf("error creating ECS v1 client: %s", err)
	}
	for i, v := range members.List() {
		member := v.(map[string]interface{})
		info := channels.MemberInfo{
			Weight: member["weight"].(int),
		}
		switch mType {
		case MemberTypeEcs:
			id, ok := member["id"]
			if !ok || id == "" {
				return result, fmt.Errorf("the instance ID is missing, please check your input of members")
			}
			info.EcsId = id.(string)
			server, err := cloudservers.Get(ecsClient, id.(string)).Extract()
			if err != nil {
				return result, fmt.Errorf("error getting ECS instance from server by ID: %s", err)
			}
			info.EcsName = server.Name
		case MemberTypeEip:
			addr, ok := member["ip_address"]
			if !ok || addr == "" {
				return result, fmt.Errorf("the IP address of EIP is missing, please check your input of members")
			}
			info.Host = addr.(string)
		default:
			return result, fmt.Errorf("the member type is wrong, please check your input of members")
		}
		result[i] = info
	}
	return result, nil
}

func buildVpcChannelParameters(d *schema.ResourceData, config *config.Config) (channels.ChannelOpts, error) {
	var (
		mType         = MemberType(d.Get("member_type").(string))
		algorithmType = AlgorithmType(d.Get("algorithm").(string))

		opt = channels.ChannelOpts{
			Name: d.Get("name").(string),
			Port: d.Get("port").(int),
			Type: 2, // The type is required and type 1 (private network ELB channel) is to be deprecated.
		}
	)

	if val, ok := memberType[mType]; ok {
		opt.MemberType = val
	} else {
		return opt, fmt.Errorf("wrong member type: %v", mType)
	}
	// Backend servers
	members, err := buildVpcChannelMembers(d, config)
	if err != nil {
		return opt, err
	}
	opt.Members = members
	// Healthy check config
	conf, err := buildVpcChannelHealthConfig(d)
	if err != nil {
		return opt, err
	}
	opt.VpcHealthConfig = conf
	// algorithm
	v, ok := balanceStrategy[algorithmType]
	if ok {
		opt.BalanceStrategy = v
	} else {
		return opt, fmt.Errorf("the value of algorithm parameter is invalid")
	}
	return opt, nil
}

func resourceVpcChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opts, err := buildVpcChannelParameters(d, cfg)
	if err != nil {
		return diag.Errorf("unable to get the create option of the VPC channel: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	v, err := channels.Create(client, instanceId, opts).Extract()
	if err != nil {
		return diag.Errorf("error creating VPC channel: %s", err)
	}
	d.SetId(v.Id)
	return resourceVpcChannelRead(ctx, d, meta)
}

func flattenVpcChannelMembers(mType string, members []channels.MemberInfo) ([]map[string]interface{}, error) {
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
			return nil, fmt.Errorf("wrong members (%+v) with the type (%v)", v, mType)
		}
		result[i] = memberMap
	}
	return result, nil
}

func extractVpcChannelMemberType(mType string) *MemberType {
	for k, v := range memberType {
		if v == mType {
			return &k
		}
	}
	return nil
}

func extractVpcChannelAlgorithm(algorithm int) *AlgorithmType {
	for k, v := range balanceStrategy {
		if v == algorithm {
			return &k
		}
	}
	return nil
}

func resourceVpcChannelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := channels.Get(client, instanceId, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC channel")
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("port", resp.Port),
		d.Set("protocol", strings.ToUpper(resp.VpcHealthConfig.Protocol)),
		d.Set("path", resp.VpcHealthConfig.Path),
		d.Set("healthy_threshold", resp.VpcHealthConfig.ThresholdNormal),
		d.Set("unhealthy_threshold", resp.VpcHealthConfig.ThresholdAbnormal),
		d.Set("timeout", resp.VpcHealthConfig.Timeout),
		d.Set("interval", resp.VpcHealthConfig.TimeInterval),
		d.Set("http_code", resp.VpcHealthConfig.HttpCodes),
		d.Set("created_at", resp.CreateTime),
	)

	status, ok := channelStatus[resp.Status]
	if !ok {
		return diag.Errorf("the response status is invalid")
	} else {
		mErr = multierror.Append(mErr, d.Set("status", status))
	}
	// Extract the member type.
	memberType := extractVpcChannelMemberType(resp.MemberType)
	if memberType != nil {
		mErr = multierror.Append(mErr, d.Set("member_type", string(*memberType)))
	}
	// Extract the algorithm.
	algorithm := extractVpcChannelAlgorithm(resp.BalanceStrategy)
	if algorithm != nil {
		mErr = multierror.Append(mErr, d.Set("algorithm", string(*algorithm)))
	}

	if members, err := flattenVpcChannelMembers(resp.MemberType, resp.Members); err != nil {
		mErr = multierror.Append(mErr, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("members", members))
	}

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}
	return nil
}

func resourceVpcChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opt, err := buildVpcChannelParameters(d, cfg)
	if err != nil {
		return diag.Errorf("unable to get the update option of the VPC channel: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	_, err = channels.Update(client, instanceId, d.Id(), opt).Extract()
	if err != nil {
		return diag.Errorf("error updating VPC channel: %s", err)
	}
	return resourceVpcChannelRead(ctx, d, meta)
}

func resourceVpcChannelDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	if err = channels.Delete(client, instanceId, d.Id()).ExtractErr(); err != nil {
		return diag.Errorf("unable to delete the VPC channel (%s): %s", d.Id(), err)
	}

	return nil
}

// The ID cannot find on console, so we need to import by VPC channel name.
func resourceVpcChannelResourceImportState(_ context.Context, d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	name := parts[1]
	opt := channels.ListOpts{
		Name: name,
	}
	pages, err := channels.List(client, instanceId, opt).AllPages()
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error retrieving VPC channels: %s", err)
	}
	resp, err := channels.ExtractChannels(pages)
	if len(resp) < 1 {
		return []*schema.ResourceData{d}, fmt.Errorf("unable to find the VPC channel (%s) form server: %s", name, err)
	}
	d.SetId(resp[0].Id)

	return []*schema.ResourceData{d}, d.Set("instance_id", instanceId)
}
