package dcs

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dcs/v2/availablezones"
	"github.com/chnsz/golangsdk/openstack/dcs/v2/instances"
	dcsTags "github.com/chnsz/golangsdk/openstack/dcs/v2/tags"
	"github.com/chnsz/golangsdk/openstack/dcs/v2/whitelists"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	chargeModePostPaid = "postPaid"
	chargeModePrePaid  = "prePaid"
)

var (
	chargingMode = map[int]string{
		0: chargeModePostPaid,
		1: chargeModePrePaid,
	}
)

func ResourceDcsInstance() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceDcsInstancesCreate,
		ReadContext:   resourceDcsInstancesRead,
		UpdateContext: resourceDcsInstancesUpdate,
		DeleteContext: resourceDcsInstancesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
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
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Redis", "Memcached",
				}, true),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"capacity": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"whitelists"},
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"access_user": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				ForceNew:  true,
			},
			"whitelist_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"whitelists": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 4,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_address": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"maintain_begin": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"maintain_end"},
				Default:      "02:00:00",
				ValidateFunc: validation.StringInSlice([]string{
					"22:00:00", "02:00:00", "06:00:00", "10:00:00", "14:00:00", "18:00:00",
				}, false),
			},
			"maintain_end": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "06:00:00",
				ValidateFunc: validation.StringInSlice([]string{
					"22:00:00", "02:00:00", "06:00:00", "10:00:00", "14:00:00", "18:00:00",
				}, false),
			},
			"backup_policy": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"backup_type", "begin_at", "period_type", "backup_at", "save_days"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"save_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 7),
						},
						"backup_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"auto", "manual"}, false),
						},
						"begin_at": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^([0-1]\d|2[0-3]):00-([0-1]\d|2[0-3]):00$`),
								"format must be HH:00-HH:00",
							),
						},
						"period_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "weekly",
							ValidateFunc: validation.StringInSlice([]string{"weekly"}, false),
						},
						"backup_at": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeInt,
								ValidateFunc: validation.IntBetween(1, 7),
							},
						},
					},
				},
			},
			"rename_commands": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenew(nil),
			"auto_pay":      common.SchemaAutoPay(nil),
			"tags":          common.TagsSchema(),
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Deprecated
			"product_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"product_id", "flavor"},
				Deprecated:   "Deprecated, please use `flavor` instead",
			},
			"available_zones": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				AtLeastOneOf: []string{"available_zones", "availability_zones"},
				Deprecated:   "Deprecated, please use `availability_zones` instead",
			},
			"enterprise_project_name": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "Deprecated, this is a non-public attribute.",
			},
			"internal_version": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Deprecated, please us `engine_version` instead.",
			},
			"ip": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Deprecated, please us `private_ip` instead.",
			},
			"user_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Deprecated",
			},
			"user_name": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Deprecated",
			},
			"save_days": {
				Type:       schema.TypeInt,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Deprecated, please use `backup_policy` instead",
			},
			"backup_type": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Deprecated, please use `backup_policy` instead",
			},
			"begin_at": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"period_type", "backup_at", "save_days", "backup_type"},
				Deprecated:   "Deprecated, please use `backup_policy` instead",
			},
			"period_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"begin_at", "backup_at", "save_days", "backup_type"},
				Deprecated:   "Please use `backup_policy` instead",
			},
			"backup_at": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"period_type", "begin_at", "save_days", "backup_type"},
				Deprecated:   "Deprecated, please use `backup_policy` instead",
				Elem:         &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func buildBackupPolicyParams(d *schema.ResourceData) *instances.InstanceBackupPolicyOpts {
	if _, ok := d.GetOk("backup_policy"); !ok { // deprecated branch
		if v, ok := d.GetOk("backup_at"); ok {
			backupAts := v.([]interface{})
			return &instances.InstanceBackupPolicyOpts{
				SaveDays:   d.Get("save_days").(int),
				BackupType: d.Get("backup_type").(string),
				PeriodicalBackupPlan: instances.BackupPlan{
					BeginAt:    d.Get("begin_at").(string),
					PeriodType: d.Get("period_type").(string),
					BackupAt:   utils.ExpandToIntList(backupAts),
				},
			}
		}
		// neither backup_policy nor backup_at is specified
		return nil
	}

	backupPolicyList := d.Get("backup_policy").([]interface{})
	if len(backupPolicyList) == 0 {
		return nil
	}
	backupPolicy := backupPolicyList[0].(map[string]interface{})
	backupType := backupPolicy["backup_type"].(string)
	if len(backupType) == 0 || backupType == "manual" {
		return nil
	}
	// build backup policy options
	backupAt := utils.ExpandToIntList(backupPolicy["backup_at"].([]interface{}))
	backupPlan := instances.BackupPlan{
		BackupAt:   backupAt,
		PeriodType: backupPolicy["period_type"].(string),
		BeginAt:    backupPolicy["begin_at"].(string),
	}
	backupPolicyOpts := &instances.InstanceBackupPolicyOpts{
		BackupType:           backupPolicy["backup_type"].(string),
		SaveDays:             backupPolicy["save_days"].(int),
		PeriodicalBackupPlan: backupPlan,
	}
	return backupPolicyOpts
}

func resourceDcsInstancesCheck(d *schema.ResourceData) error {
	engineVersion := d.Get("engine_version").(string)
	secGroupID := d.Get("security_group_id").(string)

	// check for Redis 4.0 and 5.0
	if engineVersion == "4.0" || engineVersion == "5.0" {
		if secGroupID != "" {
			return fmtp.Errorf("security_group_id is not supported for Redis 4.0 and 5.0. " +
				"please configure the whitelists alternatively")
		}
	} else {
		// check for Memcached and Redis 3.0
		if secGroupID == "" {
			return fmtp.Errorf("security_group_id is mandatory for this DCS instance")
		}
	}

	return nil
}

func buildBssParamParams(d *schema.ResourceData) instances.DcsBssParam {
	bp := instances.DcsBssParam{
		ChargingMode: d.Get("charging_mode").(string),
	}
	if strings.EqualFold(bp.ChargingMode, chargeModePrePaid) {
		bp.IsAutoRenew = d.Get("auto_renew").(string)
		bp.PeriodType = d.Get("period_unit").(string)
		bp.PeriodNum = d.Get("period").(int)
		bp.IsAutoPay = common.GetAutoPay(d)
	}
	return bp
}

func buildDcsTagsParams(tags map[string]interface{}) []dcsTags.ResourceTag {
	tagArr := make([]dcsTags.ResourceTag, 0, len(tags))
	for k, v := range tags {
		tag := dcsTags.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		tagArr = append(tagArr, tag)
	}
	return tagArr
}

func buildWhiteListParams(d *schema.ResourceData) whitelists.WhitelistOpts {
	enable := d.Get("whitelist_enable").(bool)
	groupList := d.Get("whitelists").(*schema.Set).List()

	groups := make([]whitelists.WhitelistGroupOpts, len(groupList))
	for i, v := range groupList {
		item := v.(map[string]interface{})
		groups[i] = whitelists.WhitelistGroupOpts{
			GroupName: item["group_name"].(string),
			IPList:    utils.ExpandToStringList(item["ip_address"].([]interface{})),
		}
	}

	whitelistOpts := whitelists.WhitelistOpts{
		Enable: &enable,
		Groups: groups,
	}
	return whitelistOpts
}

func waitForWhiteListCompleted(ctx context.Context, c *golangsdk.ServiceClient, d *schema.ResourceData) error {
	enable := d.Get("whitelist_enable").(bool)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{strconv.FormatBool(!enable)},
		Target:       []string{strconv.FormatBool(enable)},
		Refresh:      refreshForWhiteListEnableStatus(c, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshForWhiteListEnableStatus(c *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := whitelists.Get(c, id).Extract()
		if err != nil {
			return nil, "Error", err
		}
		return r, strconv.FormatBool(r.Enable), nil
	}
}

func resourceDcsInstancesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DcsV2Client(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DCS Client(v2): %s", err)
	}

	if err = resourceDcsInstancesCheck(d); err != nil {
		return diag.FromErr(err)
	}

	// noPasswordAccess
	noPasswordAccess := true
	if d.Get("access_user").(string) != "" || d.Get("password").(string) != "" {
		noPasswordAccess = false
	}
	// resourceSpecCode
	resourceSpecCode := d.Get("flavor").(string)
	if pid, ok := d.GetOk("product_id"); ok {
		productID := pid.(string)
		resourceSpecCode = productID[0 : len(productID)-2]
	}

	// azCodes
	var azCodes []string
	availabilityZones, ok := d.GetOk("availability_zones")
	if ok {
		azCodes = utils.ExpandToStringList(availabilityZones.([]interface{}))
	} else {
		azCodes, err = getAvailableZoneCodeByID(client, d.Get("available_zones").([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// build a creation options
	createOpts := instances.CreateOpts{
		Name:                d.Get("name").(string),
		Engine:              d.Get("engine").(string),
		EngineVersion:       d.Get("engine_version").(string),
		Capacity:            d.Get("capacity").(float64),
		InstanceNum:         1,
		SpecCode:            resourceSpecCode,
		AzCodes:             azCodes,
		Port:                d.Get("port").(int),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		EnterpriseProjectId: common.GetEnterpriseProjectID(d, conf),
		Description:         d.Get("description").(string),
		PrivateIp:           d.Get("private_ip").(string),
		MaintainBegin:       d.Get("maintain_begin").(string),
		MaintainEnd:         d.Get("maintain_end").(string),
		NoPasswordAccess:    &noPasswordAccess,
		AccessUser:          d.Get("access_user").(string),
		BssParam:            buildBssParamParams(d),
		Tags:                buildDcsTagsParams(d.Get("tags").(map[string]interface{})),
	}

	// build and set rename command if configured.
	renameCmds := d.Get("rename_commands").(map[string]interface{})
	if createOpts.Engine == "Redis" && len(renameCmds) > 0 {
		renameCommands := instances.RedisCommand{
			Command:  renameCmds["command"].(string),
			Keys:     renameCmds["keys"].(string),
			Flushdb:  renameCmds["flushdb"].(string),
			Flushall: renameCmds["flushall"].(string),
			Hgetall:  renameCmds["hgetall"].(string),
		}
		createOpts.RenameCommands = renameCommands
	}

	// build and set backup policy if configured.
	backupPolicy := buildBackupPolicyParams(d)
	if backupPolicy != nil {
		createOpts.BackupPolicy = backupPolicy
	}
	logp.Printf("[DEBUG] Create DCS instance options(hide password) : %#v", createOpts)

	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	// create instance
	r, err := instances.Create(client, createOpts)
	if err != nil || len(r.Instances) == 0 {
		return fmtp.DiagErrorf("error in creating DCS instance : %s", err)
	}
	id := r.Instances[0].InstanceId
	d.SetId(id)

	// If charging mode is PrePaid, wait for the order to be completed.
	if strings.EqualFold(d.Get("charging_mode").(string), chargeModePrePaid) {
		err = common.WaitOrderComplete(ctx, d, conf, r.OrderId)
		if err != nil {
			return fmtp.DiagErrorf("[DEBUG] Error the order is not completed while "+
				"creating DCS instance. %s : %#v", d.Id(), err)
		}
	}

	// wait for the instance to be created successfully and in running state
	err = waitForDcsInstanceCompleted(ctx, client, id, d.Timeout(schema.TimeoutCreate),
		[]string{"CREATING"}, []string{"RUNNING"})
	if err != nil {
		return diag.FromErr(err)
	}

	// create whiteList if configured.
	if d.Get("whitelists").(*schema.Set).Len() > 0 {
		whitelistOpts := buildWhiteListParams(d)
		logp.Printf("[DEBUG] Create whitelist options: %#v", whitelistOpts)

		err = whitelists.Put(client, id, whitelistOpts).ExtractErr()
		if err != nil {
			return fmtp.DiagErrorf("Error creating whitelist for DCS instance (%s): %s", id, err)
		}
		// wait for whitelist created
		err = waitForWhiteListCompleted(ctx, client, d)
		if err != nil {
			return fmtp.DiagErrorf("Error while waiting to create DCS whitelist: %s", err)
		}
	}

	return resourceDcsInstancesRead(ctx, d, meta)
}

func waitForDcsInstanceCompleted(ctx context.Context, c *golangsdk.ServiceClient, id string, timeout time.Duration,
	padding []string, target []string) error {

	stateConf := &resource.StateChangeConf{
		Pending:      padding,
		Target:       target,
		Refresh:      refreshDcsInstanceState(c, id),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("[DEBUG] Error while waiting to create/resize/delete DCS instance. %s : %#v",
			id, err)
	}
	return nil
}

func refreshDcsInstanceState(c *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := instances.Get(c, id)
		if err != nil {
			err404 := golangsdk.ErrDefault404{}
			if errors.As(err, &err404) {
				return &(instances.DcsInstance{}), "DELETED", nil
			}
			return nil, "Error", err
		}
		return r, r.Status, nil
	}
}

func resourceDcsInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DcsV2Client(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DCS Client(v2): %s", err)
	}

	r, err := instances.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DCS instance")
	}
	logp.Printf("[DEBUG] Get DCS instance : %#v", r)

	// capacity
	capacity := r.Capacity
	if capacity == 0 {
		capacity, _ = strconv.ParseFloat(r.CapacityMinor, floatBitSize)
	}

	securityGroupID := r.SecurityGroupId
	// If security_group_id is not set, the default value is returned: securityGroupId. Bad design.
	if securityGroupID == "securityGroupId" {
		securityGroupID = ""
	}

	// batch set attributes
	mErr := multierror.Append(nil,
		d.Set("region", conf.GetRegion(d)),
		d.Set("name", r.Name),
		d.Set("engine", r.Engine),
		d.Set("engine_version", r.EngineVersion),
		d.Set("capacity", capacity),
		d.Set("flavor", r.SpecCode),
		d.Set("availability_zones", r.AzCodes),
		d.Set("vpc_id", r.VpcId),
		d.Set("vpc_name", r.VpcName),
		d.Set("subnet_id", r.SubnetId),
		d.Set("subnet_name", r.SubnetName),
		d.Set("security_group_id", securityGroupID),
		d.Set("security_group_name", r.SecurityGroupName),
		d.Set("enterprise_project_id", r.EnterpriseProjectId),
		d.Set("description", r.Description),
		d.Set("private_ip", r.Ip),
		d.Set("ip", r.Ip),
		d.Set("maintain_begin", r.MaintainBegin),
		d.Set("maintain_end", r.MaintainEnd),
		d.Set("charging_mode", chargingMode[r.ChargingMode]),
		d.Set("port", r.Port),
		d.Set("status", r.Status),
		d.Set("used_memory", r.UsedMemory),
		d.Set("max_memory", r.MaxMemory),
		d.Set("domain_name", r.DomainName),
		d.Set("user_id", r.UserId),
		d.Set("user_name", r.UserName),
		d.Set("access_user", r.AccessUser),
		d.Set("order_id", r.OrderId),
	)

	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("error setting DCS instance attributes: %s", mErr)
	}

	// set backup_policy attribute
	backupPolicy := r.BackupPolicy
	if len(backupPolicy.Policy.BackupType) > 0 {
		bakPolicy := []map[string]interface{}{
			{
				"backup_type": backupPolicy.Policy.BackupType,
				"save_days":   backupPolicy.Policy.SaveDays,
				"begin_at":    backupPolicy.Policy.PeriodicalBackupPlan.BeginAt,
				"period_type": backupPolicy.Policy.PeriodicalBackupPlan.PeriodType,
				"backup_at":   backupPolicy.Policy.PeriodicalBackupPlan.BackupAt,
			},
		}
		d.Set("backup_policy", bakPolicy)
	}

	// set tags
	if resourceTags, err := tags.Get(client, "instances", d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagMap); err != nil {
			return fmtp.DiagErrorf("[DEBUG] Error saving tag to state for DCS instance (%s): %s", d.Id(), err)
		}
	} else {
		logp.Printf("[WARN] fetching tags of DCS instance failed: %s", err)
	}

	// set white list
	// some regions (cn-south-1) will fail to call the API due to the cloud reason
	// ignore the error temporarily.
	wList, err := whitelists.Get(client, d.Id()).Extract()
	if err != nil || wList == nil {
		logp.Printf("[WARN] Error fetching whitelists for DCS instance, error: %s", err)
	}

	logp.Printf("[DEBUG] Find DCS instance white list : %#v", wList.Groups)
	whiteList := make([]map[string]interface{}, len(wList.Groups))
	for i, group := range wList.Groups {
		whiteList[i] = map[string]interface{}{
			"group_name": group.GroupName,
			"ip_address": group.IPList,
		}
	}
	if len(whiteList) > 0 {
		d.Set("whitelists", whiteList)
		d.Set("whitelist_enable", wList.Enable)
	} else {
		// Set to the default value, otherwise it will prompt change after importing resources.
		d.Set("whitelist_enable", true)
	}

	return nil
}

func resourceDcsInstancesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.DcsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DCS Client(v2): %s", err)
	}

	// update basic params
	if d.HasChanges("port", "name", "description", "security_group_id", "backup_policy",
		"maintain_begin", "maintain_end") {
		desc := d.Get("description").(string)
		securityGroupID := d.Get("security_group_id").(string)
		opts := instances.ModifyInstanceOpt{
			Name:            d.Get("name").(string),
			Port:            d.Get("port").(int),
			Description:     &desc,
			MaintainBegin:   d.Get("maintain_begin").(string),
			MaintainEnd:     d.Get("maintain_end").(string),
			SecurityGroupId: &securityGroupID,
			BackupPolicy:    buildBackupPolicyParams(d),
		}
		logp.Printf("[DEBUG] Update DCS instance options : %#v", opts)

		_, err = instances.Update(client, d.Id(), opts)
		if err != nil {
			return diag.FromErr(err)
		}
		if d.HasChange("port") {
			// Modifying the port is asynchronous and needs to wait for completion.
			err = waitForPortUpdated(ctx, client, d)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// resize instance
	err = resizeDcsInstance(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// update tags
	if d.HasChange("tags") {
		oldVal, newVal := d.GetChange("tags")
		err = updateDcsTags(client, d.Id(), oldVal.(map[string]interface{}), newVal.(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update whitelist
	if d.HasChanges("whitelists", "whitelist_enable") {
		whitelistOpts := buildWhiteListParams(d)
		logp.Printf("[DEBUG] Update DCS instance whitelist options: %#v", whitelistOpts)

		err = whitelists.Put(client, d.Id(), whitelistOpts).ExtractErr()
		if err != nil {
			return fmtp.DiagErrorf("Error updating whitelist for instance (%s): %s", d.Id(), err)
		}

		// wait for whitelist updated
		err = waitForWhiteListCompleted(ctx, client, d)
		if err != nil {
			return fmtp.DiagErrorf("Error while waiting to create DCS whitelist: %s", err)
		}
	}

	return resourceDcsInstancesRead(ctx, d, meta)
}

func waitForPortUpdated(ctx context.Context, c *golangsdk.ServiceClient, d *schema.ResourceData) error {
	op, np := d.GetChange("port")
	stateConf := &resource.StateChangeConf{
		Pending:      []string{strconv.Itoa(op.(int))},
		Target:       []string{strconv.Itoa(np.(int))},
		Refresh:      refreshDcsInstancePort(c, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("[DEBUG] Error while waiting to create/resize/delete DCS instance. %s : %#v",
			d.Id(), err)
	}
	return nil
}

func refreshDcsInstancePort(c *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := instances.Get(c, id)
		if err != nil {
			return nil, "Error", err
		}
		return r, strconv.Itoa(r.Port), nil
	}
}

func updateDcsTags(c *golangsdk.ServiceClient, id string, oldVal, newVal map[string]interface{}) error {
	// remove old tags
	if len(oldVal) > 0 {
		tagList := buildDcsTagsParams(oldVal)
		err := dcsTags.Delete(c, id, tagList)
		if err != nil {
			return err
		}
	}

	// add new tags
	if len(newVal) > 0 {
		tagList := buildDcsTagsParams(newVal)
		err := dcsTags.Create(c, id, tagList)
		if err != nil {
			return err
		}
	}
	return nil
}

func resizeDcsInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.DcsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud DCS Client(v2): %s", err)
	}

	if d.HasChanges("flavor", "capacity") {
		specCode := d.Get("flavor").(string)
		opts := instances.ResizeInstanceOpts{
			SpecCode:    specCode,
			NewCapacity: d.Get("capacity").(float64),
		}
		if d.Get("charging_mode").(string) == chargeModePrePaid {
			opts.BssParam = instances.DcsBssParamOpts{
				IsAutoPay: common.GetAutoPay(d),
			}
		}
		logp.Printf("[DEBUG] Resize DCS dcs instance options : %#v", opts)

		r, err := instances.ResizeInstance(client, d.Id(), opts)
		if err != nil {
			return err
		}

		if d.Get("charging_mode").(string) == chargeModePrePaid {
			// wait for order pay
			err = common.WaitOrderComplete(ctx, d, config, r.OrderId)
			if err != nil {
				return err
			}
		}

		// wait for dcs instance change
		err = waitForDcsInstanceCompleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate),
			[]string{"EXTENDING"}, []string{"RUNNING"})
		if err != nil {
			return err
		}

		// check the result of the change
		instance, err := instances.Get(client, d.Id())
		if err != nil {
			return common.CheckDeleted(d, err, "DCS instance")
		}
		if instance.SpecCode != d.Get("flavor").(string) {
			return fmtp.Errorf("[ERROR] Change flavor failed, "+
				"after changed the DCS flavor still is: %s, expected: %s.", instance.SpecCode, specCode)
		}
	}
	return nil
}

func resourceDcsInstancesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DcsV2Client(conf.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DCS Client(v2): %s", err)
	}

	// for prePaid mode, we should unsubscribe the resource
	if d.Get("charging_mode").(string) == chargeModePrePaid {
		err = common.UnsubscribePrePaidResource(d, conf, []string{d.Id()})
		if err != nil {
			return fmtp.DiagErrorf("error unsubscribing HuaweiCloud DCS redis instance : %s", err)
		}
	} else {
		err = instances.Delete(client, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Waiting to delete success
	err = waitForDcsInstanceCompleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete),
		[]string{"RUNNING"}, []string{"DELETED"})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func getAvailableZoneCodeByID(client *golangsdk.ServiceClient, azIds []interface{}) ([]string, error) {
	azCodes := make([]string, 0, len(azIds))
	if len(azIds) == 0 {
		return azCodes, fmtp.Errorf("availability_zones are required")
	}

	list, err := availablezones.List(client)
	if err != nil {
		return azCodes, err
	}

	mapping := make(map[string]string)
	for _, v := range list.AvailableZones {
		mapping[v.ID] = v.Code
	}

	for _, id := range azIds {
		azID := id.(string)
		if v, ok := mapping[azID]; ok {
			azCodes = append(azCodes, v)
		} else {
			return azCodes, fmtp.Errorf("Invalid available zone code: %s", azID)
		}
	}

	return azCodes, nil
}
