package dms

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dms/v2/availablezones"
	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"
	"github.com/chnsz/golangsdk/openstack/dms/v2/products"
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

func ResourceDmsKafkaInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaInstanceCreate,
		ReadContext:   resourceDmsKafkaInstanceRead,
		UpdateContext: resourceDmsKafkaInstanceUpdate,
		DeleteContext: resourceDmsKafkaInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
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
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_space": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"manager_user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"manager_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"access_user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"password",
				},
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				ForceNew:  true,
				RequiredWith: []string{
					"access_user",
				},
			},
			"maintain_begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_ip_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"retention_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"produce_reject", "time_base",
				}, false),
			},
			"dumping": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_auto_topic": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"partition_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable_public_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ssl_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"used_storage_space": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manegement_connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_zones": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				AtLeastOneOf: []string{"available_zones", "availability_zones"},
				Deprecated:   "available_zones has deprecated, please use \"availability_zones\" instead.",
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Deprecated: "The bandwidth has been deprecated. " +
					"If you need to change the bandwidth, please update the product_id.",
			},
		},
	}
}

func resourceDmsKafkaPublicIpIDs(d *schema.ResourceData, bandwidth string) (string, error) {
	publicIpIDsRaw := d.Get("public_ip_ids").([]interface{})

	IdNumMap := map[string]int{
		"100MB":  3,
		"300MB":  3,
		"600MB":  4,
		"1200MB": 8,
	}
	if IdNumMap[bandwidth] != len(publicIpIDsRaw) {
		return "", fmtp.Errorf("error creating HuaweiCloud DMS kafka instance: "+
			"%d public ip IDs needed when bandwidth is set to %s, but got %d",
			IdNumMap[bandwidth], bandwidth, len(publicIpIDsRaw))
	}

	publicIpIDs := utils.ExpandToStringList(publicIpIDsRaw)
	return strings.Join(publicIpIDs, ","), nil
}

func getProductDetail(config *config.Config, d *schema.ResourceData) (*products.Detail, error) {
	products, err := getProducts(config, config.GetRegion(d), "kafka")
	if err != nil {
		return nil, fmtp.Errorf("error querying product detail, please check product_id, error: %s", err)
	}

	productID := d.Get("product_id").(string)
	engineVersion := d.Get("engine_version").(string)

	for _, ps := range products.Hourly {
		if ps.Version != engineVersion {
			continue
		}
		for _, v := range ps.Values {
			for _, p := range v.Details {
				if p.ProductID == productID {
					return &p, nil
				}
			}
		}
	}
	return nil, fmtp.Errorf("can not found product detail base on product_id: %s", productID)
}

func resourceDmsKafkaInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	dmsV2Client, err := config.DmsV2Client(region)
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DMS instance client: %s", err)
	}

	product, err := getProductDetail(config, d)
	if err != nil {
		return fmtp.DiagErrorf("Error querying product detail: %s", err)
	}

	bandwidth := product.Bandwidth
	defaultPartitionNum, _ := strconv.ParseInt(product.PartitionNum, 10, 64)
	defaultStorageSpace, _ := strconv.ParseInt(product.Storage, 10, 64)

	// check storage
	storageSpace, ok := d.GetOk("storage_space")
	if ok && storageSpace.(int) < int(defaultStorageSpace) {
		return fmtp.DiagErrorf("The storage capacity is less than the default capacity of the product. "+
			"The default storage capacity of product is %v, storage_space is %v.", defaultStorageSpace, storageSpace)
	}

	sslEnable := false
	if d.Get("access_user").(string) != "" && d.Get("password").(string) != "" {
		sslEnable = true
	}

	var availableZones []string
	zoneIDs, ok := d.GetOk("available_zones")
	if ok {
		availableZones = utils.ExpandToStringList(zoneIDs.([]interface{}))
	} else {
		// convert the codes of the availability zone into ids
		azCodes := d.Get("availability_zones").([]interface{})
		availableZones, err = getAvailableZoneIDByCode(config, region, azCodes)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	createOpts := &instances.CreateOps{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Engine:              "kafka",
		EngineVersion:       d.Get("engine_version").(string),
		Specification:       bandwidth,
		StorageSpace:        int(defaultStorageSpace),
		PartitionNum:        int(defaultPartitionNum),
		AccessUser:          d.Get("access_user").(string),
		VPCID:               d.Get("vpc_id").(string),
		SecurityGroupID:     d.Get("security_group_id").(string),
		SubnetID:            d.Get("network_id").(string),
		AvailableZones:      availableZones,
		ProductID:           d.Get("product_id").(string),
		KafkaManagerUser:    d.Get("manager_user").(string),
		MaintainBegin:       d.Get("maintain_begin").(string),
		MaintainEnd:         d.Get("maintain_end").(string),
		SslEnable:           sslEnable,
		RetentionPolicy:     d.Get("retention_policy").(string),
		ConnectorEnalbe:     d.Get("dumping").(bool),
		EnableAutoTopic:     d.Get("enable_auto_topic").(bool),
		StorageSpecCode:     d.Get("storage_spec_code").(string),
		EnterpriseProjectID: common.GetEnterpriseProjectID(d, config),
	}

	if _, ok := d.GetOk("public_ip_ids"); ok {
		publicIpIDs, err := resourceDmsKafkaPublicIpIDs(d, bandwidth)
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.EnablePublicIP = true
		createOpts.PublicIpID = publicIpIDs
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		createOpts.Tags = taglist
	}

	logp.Printf("[DEBUG] Create DMS Kafka instance options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)
	createOpts.KafkaManagerPassword = d.Get("manager_password").(string)

	v, err := instances.Create(dmsV2Client, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DMS kafka instance: %s", err)
	}
	logp.Printf("[INFO] instance ID: %s", v.InstanceID)

	// Store the instance ID now
	d.SetId(v.InstanceID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      DmsKafkaInstanceStateRefreshFunc(dmsV2Client, v.InstanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        300 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("error waiting for instance (%s) to become ready: %s", v.InstanceID, err)
	}

	// resize storage capacity of the instance
	if ok && storageSpace.(int) != int(defaultStorageSpace) {
		err = resizeInstance(ctx, d, meta)
		if err != nil {
			dErrs := fmtp.DiagErrorf("Kafka instance has created, "+
				"but an error occurred while resizing the storage capacity. "+
				"Current storage capacity are %vGB, expected storage_space=%vGB, error message: %s ",
				defaultStorageSpace, storageSpace.(int), err)
			dErrs[0].Severity = diag.Warning
			return dErrs
		}
	}

	return resourceDmsKafkaInstanceRead(ctx, d, meta)
}

func resourceDmsKafkaInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	dmsV2Client, err := config.DmsV2Client(region)
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DMS instance client: %s", err)
	}
	v, err := instances.Get(dmsV2Client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS instance")
	}
	logp.Printf("[DEBUG] DMS kafka instance created success %s: %+v", d.Id(), v)
	d.SetId(v.InstanceID)

	partitionNum, _ := strconv.ParseInt(v.PartitionNum, 10, 64)
	// convert the ids of the availability zone into codes
	availableZoneIDs := v.AvailableZones
	availableZoneCodes, err := getAvailableZoneCodeByID(config, region, availableZoneIDs)
	mErr := multierror.Append(nil, err)

	mErr = multierror.Append(mErr,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", v.Name),
		d.Set("description", v.Description),
		d.Set("engine", v.Engine),
		d.Set("engine_version", v.EngineVersion),
		d.Set("bandwidth", v.Specification),
		// storage_space indicates total_storage_space while creating
		// set value of total_storage_space to storage_space to keep consistent
		d.Set("storage_space", v.TotalStorageSpace),
		d.Set("partition_num", partitionNum),
		d.Set("vpc_id", v.VPCID),
		d.Set("security_group_id", v.SecurityGroupID),
		d.Set("network_id", v.SubnetID),
		d.Set("available_zones", availableZoneIDs),
		d.Set("availability_zones", availableZoneCodes),
		d.Set("product_id", v.ProductID),
		d.Set("manager_user", v.KafkaManagerUser),
		d.Set("maintain_begin", v.MaintainBegin),
		d.Set("maintain_end", v.MaintainEnd),
		d.Set("enable_public_ip", v.EnablePublicIP),
		d.Set("ssl_enable", v.SslEnable),
		d.Set("retention_policy", v.RetentionPolicy),
		d.Set("dumping", v.ConnectorEnalbe),
		d.Set("enable_auto_topic", v.EnableAutoTopic),
		d.Set("storage_spec_code", v.StorageSpecCode),
		d.Set("enterprise_project_id", v.EnterpriseProjectID),
		d.Set("used_storage_space", v.UsedStorageSpace),
		d.Set("connect_address", v.ConnectAddress),
		d.Set("port", v.Port),
		d.Set("status", v.Status),
		d.Set("resource_spec_code", v.ResourceSpecCode),
		d.Set("user_id", v.UserID),
		d.Set("user_name", v.UserName),
		d.Set("manegement_connect_address", v.ManagementConnectAddress),
		d.Set("type", v.Type),
		d.Set("access_user", v.AccessUser),
	)
	// set tags
	engine := "kafka"
	if resourceTags, err := tags.Get(dmsV2Client, engine, d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			e := fmtp.Errorf("error saving tags to state for DMS kafka instance (%s): %s", d.Id(), err)
			mErr = multierror.Append(mErr, e)
		}
	} else {
		logp.Printf("[WARN] error fetching tags of DMS kafka instance (%s): %s", d.Id(), err)
	}

	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("Error setting attributes for DMS kafka instance: %s", mErr)
	}

	return nil
}

func resourceDmsKafkaInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DMS instance client: %s", err)
	}

	var mErr *multierror.Error
	if d.HasChanges("name", "description", "maintain_begin", "maintain_end",
		"security_group_id", "retention_policy", "enterprise_project_id") {
		description := d.Get("description").(string)
		updateOpts := instances.UpdateOpts{
			Description:         &description,
			MaintainBegin:       d.Get("maintain_begin").(string),
			MaintainEnd:         d.Get("maintain_end").(string),
			SecurityGroupID:     d.Get("security_group_id").(string),
			RetentionPolicy:     d.Get("retention_policy").(string),
			EnterpriseProjectID: d.Get("enterprise_project_id").(string),
		}

		if d.HasChange("name") {
			updateOpts.Name = d.Get("name").(string)
		}

		err = instances.Update(dmsV2Client, d.Id(), updateOpts).Err
		if err != nil {
			e := fmtp.Errorf("error updating HuaweiCloud DMS kafka Instance: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	if d.HasChanges("storage_space", "product_id") {
		err = resizeInstance(ctx, d, meta)
		if err != nil {
			e := fmtp.Errorf("error resizing HuaweiCloud DMS kafka Instance: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	if d.HasChange("tags") {
		// update tags
		engine := "kafka"
		tagErr := utils.UpdateResourceTags(dmsV2Client, d, engine, d.Id())
		if tagErr != nil {
			e := fmtp.Errorf("error updating tags of DMS kafka instance:%s, err:%s", d.Id(), tagErr)
			mErr = multierror.Append(mErr, e)
		}
	}
	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("error while updating DMS Kafka instances, there %s", mErr)
	}
	return resourceDmsKafkaInstanceRead(ctx, d, meta)
}

func resizeInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud DMS instance client: %s", err)
	}

	product, err := getProductDetail(config, d)
	if err != nil || product == nil {
		return fmtp.Errorf("change storage_space failed, error querying product detail: %s", err)
	}

	storage := d.Get("storage_space").(int)
	opts := instances.ResizeInstanceOpts{
		NewSpecCode:     product.SpecCode,
		NewStorageSpace: storage,
	}
	_, err = instances.Resize(dmsV2Client, d.Id(), opts)
	if err != nil {
		return fmtp.Errorf("resize failed, error: %s", err)
	}

	productID := d.Get("product_id").(string)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"EXTENDING", "REFRESHING"},
		Target:       []string{"RUNNING"},
		Refresh:      refreshResizeProductIDFunc(dmsV2Client, d.Id(), productID),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        300 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("error waiting for instance (%s) to resized: %v", d.Id(), err)
	}
	return nil
}

func resourceDmsKafkaInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DMS instance client: %s", err)
	}

	err = instances.Delete(dmsV2Client, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("error deleting HuaweiCloud instance: %s", err)
	}

	// Wait for the instance to delete before moving on.
	logp.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETING", "RUNNING"},
		Target:       []string{"DELETED"},
		Refresh:      DmsKafkaInstanceStateRefreshFunc(dmsV2Client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        120 * time.Second,
		PollInterval: 15 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf(
			"error waiting for instance (%s) to delete: %s",
			d.Id(), err)
	}

	logp.Printf("[DEBUG] DMS instance %s deactivated", d.Id())
	d.SetId("")
	return nil
}

func DmsKafkaInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}

		return v, v.Status, nil
	}
}

func refreshResizeProductIDFunc(client *golangsdk.ServiceClient, instanceID,
	productID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}
		if v.Status == "RUNNING" && v.ProductID != productID {
			return v, "REFRESHING", nil
		}
		return v, v.Status, nil
	}
}

func getAvailableZoneIDByCode(config *config.Config, region string, azCodes []interface{}) ([]string, error) {
	if len(azCodes) == 0 {
		return nil, fmtp.Errorf("availability_zones is required")
	}

	availableZones, err := getAvailableZones(config, region)
	if err != nil {
		return nil, err
	}

	mappingData := make(map[string]availablezones.AvailableZone)
	for _, v := range availableZones {
		mappingData[v.Code] = v
	}

	azIDs := make([]string, 0, len(azCodes))
	for _, code := range azCodes {
		if az, ok := mappingData[code.(string)]; ok {
			azIDs = append(azIDs, az.ID)
		}
	}
	logp.Printf("[DEBUG] DMS convert the codes of the availability zone into ids: \n%#v => \n%#v",
		azCodes, azIDs)
	return azIDs, nil
}

func getAvailableZoneCodeByID(config *config.Config, region string, azIDs []string) ([]string, error) {
	if len(azIDs) == 0 {
		return nil, fmtp.Errorf("availability_zones is required")
	}

	availableZones, err := getAvailableZones(config, region)
	if err != nil {
		return nil, err
	}

	mappingData := make(map[string]availablezones.AvailableZone)
	for _, v := range availableZones {
		mappingData[v.ID] = v
	}

	azCodes := make([]string, 0, len(mappingData))
	for _, id := range azIDs {
		if az, ok := mappingData[id]; ok {
			azCodes = append(azCodes, az.Code)
		}
	}
	logp.Printf("[DEBUG] DMS convert the ids of the availability zone into codes: \n%#v => \n%#v",
		azIDs, azCodes)
	return azCodes, nil
}

func getAvailableZones(config *config.Config, region string) ([]availablezones.AvailableZone, error) {
	dmsV2Client, err := config.DmsV2Client(region)
	if err != nil {
		return nil, fmtp.Errorf("Error creating HuaweiCloud DMS client V2 : %s", err)
	}

	r, err := availablezones.Get(dmsV2Client)
	if err != nil {
		return nil, fmtp.Errorf("Error querying available Zones: %s", err)
	}

	return r.AvailableZones, nil
}
