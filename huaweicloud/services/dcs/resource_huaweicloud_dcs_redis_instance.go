package dcs

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	dcsTags "github.com/chnsz/golangsdk/openstack/dcs/v2/tags"
	"github.com/chnsz/golangsdk/openstack/dcs/v2/whitelists"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dcs/v2/instances"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	floatBitSize = 64

	chargeModePostPaid = "postPaid"
	chargeModePrePaid  = "prePaid"
)

var (
	chargingMode = map[int]string{
		0: chargeModePostPaid,
		1: chargeModePrePaid,
	}
)

func ResourceDcsRedisInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsRedisInstanceCreate,
		ReadContext:   resourceDcsRedisInstanceRead,
		UpdateContext: resourceDcsRedisInstanceUpdate,
		DeleteContext: resourceDcsRedisInstanceDelete,
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
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[A-Za-z][\u4e00-\u9fa5\\w]{3,63}$"),
					"name must be 4 to 64 characters in length and start with a letter. "+
						"Only Chinese characters, letters (case-insensitive), digits, underscores (_), "+
						"and hyphens (-) are allowed."),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"capacity": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"available_zones": {
				Type:     schema.TypeList,
				Required: true,
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
							MaxItems: 20,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"maintain_begin": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "02:00:00",
				ValidateFunc: validation.StringInSlice([]string{
					"22:00:00", "02:00:00", "06:00:00", "10:00:00", "14:00:00", "18:00:00",
				}, false),
			},
			"maintain_end": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "06:00:00",
				RequiredWith: []string{"maintain_begin"},
				ValidateFunc: validation.StringInSlice([]string{
					"22:00:00", "02:00:00", "06:00:00", "10:00:00", "14:00:00", "18:00:00",
				}, false),
			},
			"backup_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "auto",
							ValidateFunc: validation.StringInSlice([]string{"auto", "manual"}, false),
						},
						"save_days": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 7),
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
						"begin_at": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^([0-1]\d|2[0-3]):00-([0-1]\d|2[0-3]):00$`),
								"format must be HH:00-HH:00",
							),
						},
						"timezone_offset": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^([+-](0[1-9]|1[0-2])00|\+0000)$`),
								"must be between -1200 and +1200",
							),
						},
					},
				},
			},
			"tags":          common.TagsSchema(),
			"charging_mode": common.SchemeChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenew(nil),
			"rename_commands": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
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
		},
	}
}

func buildCreateOptsParams(d *schema.ResourceData, config *config.Config) *instances.CreateOpts {
	password, pwdOk := d.GetOk("password")
	noPasswordAccess := !pwdOk

	// build a creation options
	opts := instances.CreateOpts{
		Name:                d.Get("name").(string),
		Engine:              "Redis",
		EngineVersion:       d.Get("engine_version").(string),
		Capacity:            d.Get("capacity").(float64),
		InstanceNum:         1,
		SpecCode:            d.Get("resource_spec_code").(string),
		AzCodes:             utils.ExpandToStringList(d.Get("available_zones").([]interface{})),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		EnterpriseProjectId: common.GetEnterpriseProjectID(d, config),
		Description:         d.Get("description").(string),
		PrivateIp:           d.Get("ip").(string),
		MaintainBegin:       d.Get("maintain_begin").(string),
		MaintainEnd:         d.Get("maintain_end").(string),
		NoPasswordAccess:    &noPasswordAccess,
		BssParam:            buildBssParamParams(d),
		Tags:                buildDcsTagsParams(d.Get("tags").(map[string]interface{})),
		Port:                d.Get("port").(int),
	}

	// build and set rename command
	renameCmds := d.Get("rename_commands").(map[string]interface{})
	if len(renameCmds) > 0 {
		renameCommands := instances.RedisCommand{
			Command:  renameCmds["command"].(string),
			Keys:     renameCmds["keys"].(string),
			Flushdb:  renameCmds["flushdb"].(string),
			Flushall: renameCmds["flushall"].(string),
			Hgetall:  renameCmds["hgetall"].(string),
		}
		opts.RenameCommands = renameCommands
	}

	// build and set backup policy if configured.
	backupPolicy := buildBackupPolicyParams(d.Get("backup_policy").([]interface{}))
	if backupPolicy != nil {
		opts.BackupPolicy = backupPolicy
	}
	logp.Printf("[DEBUG] Create DCS redis options(hide password) : %#v", opts)

	// set password after printing log
	if pwdOk {
		opts.Password = password.(string)
	}

	return &opts
}

func buildBackupPolicyParams(p []interface{}) *instances.InstanceBackupPolicyOpts {
	if len(p) == 0 {
		return nil
	}

	backupPolicy := p[0].(map[string]interface{})
	backupType := backupPolicy["backup_type"].(string)
	if len(backupType) == 0 || backupType == "manual" {
		return nil
	}
	// build backup policy options
	backupAt := utils.ExpandToIntList(backupPolicy["backup_at"].([]interface{}))
	backupPlan := instances.BackupPlan{
		TimezoneOffset: backupPolicy["timezone_offset"].(string),
		BackupAt:       backupAt,
		PeriodType:     backupPolicy["period_type"].(string),
		BeginAt:        backupPolicy["begin_at"].(string),
	}
	backupPolicyOpts := &instances.InstanceBackupPolicyOpts{
		BackupType:           backupPolicy["backup_type"].(string),
		SaveDays:             backupPolicy["save_days"].(int),
		PeriodicalBackupPlan: backupPlan,
	}
	return backupPolicyOpts
}

func buildBssParamParams(d *schema.ResourceData) instances.DcsBssParam {
	bp := instances.DcsBssParam{
		ChargingMode: d.Get("charging_mode").(string),
	}
	if strings.EqualFold(bp.ChargingMode, chargeModePrePaid) {
		bp.IsAutoRenew = d.Get("auto_renew").(string)
		bp.IsAutoPay = "true"
		bp.PeriodType = d.Get("period_unit").(string)
		bp.PeriodNum = d.Get("period").(int)
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

func resourceDcsRedisInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.DcsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DCS Client(v2): %s", err)
	}

	// check for security_group_id and version
	version := d.Get("engine_version").(string)
	if _, ok := d.GetOk("security_group_id"); ok {
		if version != "3.0" {
			return fmtp.DiagErrorf("Only Redis 3.0 supports security_group_id, " +
				"please configure whitelist for other versions.")
		}
	} else {
		if version == "3.0" {
			return fmtp.DiagErrorf("security_group_id can not be empty for Redis 3.0.")
		}
	}

	// build a creation options
	opts := buildCreateOptsParams(d, config)

	// create redis instance
	r, err := instances.Create(client, *opts)
	if err != nil || len(r.Instances) == 0 {
		return fmtp.DiagErrorf("error in creating DCS redis : %s", err)
	}
	id := r.Instances[0].InstanceId
	d.SetId(id)

	// If charging mode is PrePaid, wait for the order to be completed.
	if strings.EqualFold(d.Get("charging_mode").(string), chargeModePrePaid) {
		err = common.WaitOrderComplete(ctx, d, config, r.OrderId)
		if err != nil {
			return fmtp.DiagErrorf("[DEBUG] Error the order is not completed while "+
				"creating DCS redis instance. %s : %#v", d.Id(), err)
		}
	}

	// wait for the instance to be created successfully and in running state
	err = waitForRedisInstanceCompleted(ctx, client, id, d.Timeout(schema.TimeoutCreate),
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
			return fmtp.DiagErrorf("Error creating whitelist for redis instance (%s): %s", id, err)
		}
		// wait for whitelist created
		err = waitForWhiteListCompleted(ctx, client, d)
		if err != nil {
			return fmtp.DiagErrorf("Error while waiting to create redis whitelist: %s", err)
		}
	}

	return resourceDcsRedisInstanceRead(ctx, d, meta)
}

func waitForRedisInstanceCompleted(ctx context.Context, c *golangsdk.ServiceClient,
	id string, timeout time.Duration, padding []string, target []string) error {

	stateConf := &resource.StateChangeConf{
		Pending:      padding,
		Target:       target,
		Refresh:      refreshDcsInstanceState(c, id),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("[DEBUG] Error while waiting to create/resize/delete DCS redis instance. %s : %#v",
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

func waitForWhiteListCompleted(ctx context.Context, c *golangsdk.ServiceClient, d *schema.ResourceData) error {
	enable := d.Get("whitelist_enable").(bool)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{strconv.FormatBool(!enable)},
		Target:       []string{strconv.FormatBool(enable)},
		Refresh:      refreshForWhiteListEnableStatus(c, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
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

func resourceDcsRedisInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.DcsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DCS Client(v2): %s", err)
	}

	r, err := instances.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DCS instance")
	}
	logp.Printf("[DEBUG] Get DCS redis instance : %#v", r)

	// capacity
	capacity := r.Capacity
	if capacity == 0 {
		capacity, _ = strconv.ParseFloat(r.CapacityMinor, floatBitSize)
	}
	// batch set attributes
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", r.Name),
		d.Set("engine_version", r.EngineVersion),
		d.Set("capacity", capacity),
		d.Set("resource_spec_code", r.SpecCode),
		d.Set("available_zones", r.AzCodes),
		d.Set("vpc_id", r.VpcId),
		d.Set("subnet_id", r.SubnetId),
		d.Set("security_group_id", r.SecurityGroupId),
		d.Set("enterprise_project_id", r.EnterpriseProjectId),
		d.Set("description", r.Description),
		d.Set("ip", r.Ip),
		d.Set("maintain_begin", r.MaintainBegin),
		d.Set("maintain_end", r.MaintainEnd),
		d.Set("charging_mode", chargingMode[r.ChargingMode]),
		d.Set("port", r.Port),
		d.Set("status", r.Status),
		d.Set("used_memory", r.UsedMemory),
		d.Set("max_memory", r.MaxMemory),
		d.Set("domain_name", r.DomainName),
	)
	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("error setting DCS redis instance attributes: %s", mErr)
	}

	// set backup_policy attribute
	backupPolicy := r.BackupPolicy
	if len(backupPolicy.Policy.BackupType) > 0 {
		bakPolicy := make([]map[string]interface{}, 0, 1)
		bp := make(map[string]interface{})
		bp["backup_type"] = backupPolicy.Policy.BackupType
		bp["save_days"] = backupPolicy.Policy.SaveDays
		bp["begin_at"] = backupPolicy.Policy.PeriodicalBackupPlan.BeginAt
		bp["period_type"] = backupPolicy.Policy.PeriodicalBackupPlan.PeriodType
		bp["backup_at"] = backupPolicy.Policy.PeriodicalBackupPlan.BackupAt
		bp["timezone_offset"] = backupPolicy.Policy.PeriodicalBackupPlan.TimezoneOffset
		bakPolicy = append(bakPolicy, bp)
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
	wList, err := whitelists.Get(client, d.Id()).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error setting whitelists for DCS instance, error: %s", err)
	}
	logp.Printf("[DEBUG] Find DCS instance white list : %#v", r)

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
	}

	return nil
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

//lintignore:R019
func resourceDcsRedisInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.DcsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DCS Client(v2): %s", err)
	}

	// update basic params
	if d.HasChanges("name", "description", "port", "security_group_id", "backup_policy",
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
			BackupPolicy:    buildBackupPolicyParams(d.Get("backup_policy").([]interface{})),
		}
		logp.Printf("[DEBUG] Update DCS redis instance options : %#v", opts)

		_, err = instances.Update(client, d.Id(), opts)
		if err != nil {
			return diag.FromErr(err)
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

	// update password
	if d.HasChanges("password") {
		old, nw := d.GetChange("password")

		opts := instances.UpdatePasswordOpts{
			OldPassword: old.(string),
			NewPassword: nw.(string),
		}
		logp.Printf("[DEBUG] Update password of DCS redis instance options(%s)", d.Id())

		_, err = instances.UpdatePassword(client, d.Id(), opts)
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

	return resourceDcsRedisInstanceRead(ctx, d, meta)
}

func resizeDcsInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.DcsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud DCS Client(v2): %s", err)
	}

	if d.HasChanges("resource_spec_code", "capacity") {
		specCode := d.Get("resource_spec_code").(string)
		opts := instances.ResizeInstanceOpts{
			SpecCode:    specCode,
			NewCapacity: d.Get("capacity").(float64),
		}
		if d.Get("charging_mode").(string) == chargeModePrePaid {
			opts.BssParam = instances.DcsBssParamOpts{
				IsAutoPay: "true",
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
		err = waitForRedisInstanceCompleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate),
			[]string{"EXTENDING"}, []string{"RUNNING"})
		if err != nil {
			return err
		}

		// check the result of the change
		inst, err := instances.Get(client, d.Id())
		if err != nil {
			return common.CheckDeleted(d, err, "DCS instance")
		}
		if inst.SpecCode != d.Get("resource_spec_code").(string) {
			return fmtp.Errorf("[ERROR] Change specification failed, "+
				"the specification code of DCS instance still is: %s, expected: %s.", inst.SpecCode, specCode)
		}
	}
	return nil
}

func resourceDcsRedisInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.DcsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DCS Client(v2): %s", err)
	}

	// for prePaid mode, we should unsubscribe the resource
	if d.Get("charging_mode").(string) == chargeModePrePaid {
		err = common.UnsubscribePrePaidResource(d, config, []string{d.Id()})
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
	err = waitForRedisInstanceCompleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete),
		[]string{"RUNNING"}, []string{"DELETED"})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
