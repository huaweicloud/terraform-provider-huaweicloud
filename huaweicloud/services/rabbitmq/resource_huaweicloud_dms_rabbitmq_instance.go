package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dms/v2/products"
	"github.com/chnsz/golangsdk/openstack/dms/v2/rabbitmq/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/kafka"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const engineRabbitMQ = "rabbitmq"

// @API RabbitMQ POST /v2/{engine}/{project_id}/instances
// @API RabbitMQ POST /v2/{project_id}/instances
// @API RabbitMQ POST /v2/{engine}/{project_id}/instances/{instance_id}/extend
// @API RabbitMQ DELETE /v2/{project_id}/instances/{instance_id}
// @API RabbitMQ GET /v2/{project_id}/instances/{instance_id}
// @API RabbitMQ PUT /v2/{project_id}/instances/{instance_id}
// @API RabbitMQ GET /v2/{project_id}/rabbitmq/{instance_id}/tags
// @API RabbitMQ POST /v2/{project_id}/rabbitmq/{instance_id}/tags/action
// @API RabbitMQ GET /v2/available-zones
// @API RabbitMQ GET /v2/products
// @API RabbitMQ POST /v2/{project_id}/instances/{instance_id}/password
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceDmsRabbitmqInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRabbitmqInstanceCreate,
		ReadContext:   resourceDmsRabbitmqInstanceRead,
		UpdateContext: resourceDmsRabbitmqInstanceUpdate,
		DeleteContext: resourceDmsRabbitmqInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
				Optional: true,
				ForceNew: true,
				Default:  "3.7.17",
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
			"access_user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"availability_zones": {
				// There is a problem with order of elements in Availability Zone list returned by RabbitMQ API.
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"available_zones"},
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				Description:   "schema: Required",
			},
			"flavor_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"product_id"},
				RequiredWith: []string{"storage_space"},
			},
			"broker_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"flavor_id"},
			},
			// The API return format is "HH:mm:ss" for `maintain_begin` and `maintain_end`.
			"maintain_begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
					return regexp.MustCompile(fmt.Sprintf("^%s", n)).MatchString(o)
				},
			},
			"maintain_end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
					return regexp.MustCompile(fmt.Sprintf("^%s", n)).MatchString(o)
				},
			},
			"ssl_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"public_ip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_acl": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disk_encrypted_enable": {
				Type:         schema.TypeBool,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"disk_encrypted_key"},
				Description:  "Whether to enable disk encryption.",
			},
			"disk_encrypted_key": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"disk_encrypted_enable"},
				Description:  "The key ID of the disk encryption.",
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),

			"tags": common.TagsSchema(),
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_public_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"used_storage_space": {
				Type:     schema.TypeInt,
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
			"connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"management_connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extend_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_logical_volume": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"public_ip_address": {
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
			"product_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "product_id has deprecated, please use \"flavor_id\" instead.",
			},
			// Typo, it is only kept in the code, will not be shown in the docs.
			"manegement_connect_address": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "typo in manegement_connect_address, please use \"management_connect_address\" instead.",
			},
		},
	}
}

func getProducts(cfg *config.Config, region, engine string) (*products.GetResponse, error) {
	dmsV2Client, err := cfg.DmsV2Client(region)
	if err != nil {
		return nil, fmt.Errorf("error getting DMS product client V2: %s", err)
	}
	v, err := products.Get(dmsV2Client, engine) // nolint: staticcheck
	return v, err
}

func getRabbitMQProductDetail(cfg *config.Config, d *schema.ResourceData) (*products.ProductInfo, error) {
	productRsp, err := getProducts(cfg, cfg.GetRegion(d), "rabbitmq")
	if err != nil {
		return nil, fmt.Errorf("error querying product detail, please check product_id, error: %s", err)
	}

	productID := d.Get("product_id").(string)
	engineVersion := d.Get("engine_version").(string)

	for _, ps := range productRsp.Hourly {
		if ps.Version != engineVersion {
			continue
		}
		for _, v := range ps.Values {
			for _, detail := range v.Details {
				// All informations of product for single instance type and the kafka engine type are stored in the
				// detail structure.
				if v.Name == "single" {
					if detail.ProductID == productID {
						return &products.ProductInfo{
							Storage:          detail.Storage,
							ProductID:        detail.ProductID,
							SpecCode:         detail.SpecCode,
							IOs:              detail.IOs,
							AvailableZones:   detail.AvailableZones,
							UnavailableZones: detail.UnavailableZones,
						}, nil
					}
				} else {
					for _, product := range detail.ProductInfos {
						if product.ProductID == productID {
							return &product, nil
						}
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("can not found product detail base on product_id: %s", productID)
}

func resourceDmsRabbitmqInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var dErr diag.Diagnostics
	if _, ok := d.GetOk("flavor_id"); ok {
		dErr = createRabbitMQInstanceWithFlavor(ctx, d, meta)
	} else {
		dErr = createRabbitMQInstanceWithProductID(ctx, d, meta)
	}
	if dErr != nil {
		return dErr
	}

	return resourceDmsRabbitmqInstanceRead(ctx, d, meta)
}

func createRabbitMQInstanceWithFlavor(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DmsV2Client(region)
	if err != nil {
		return diag.Errorf("error initializing DMS RabbitMQ(v2) client: %s", err)
	}

	var availableZones []string // Available zones IDs
	azIDs, ok := d.GetOk("availability_zones")
	if ok {
		// convert the codes of the availability zone into ids
		azCodes := azIDs.(*schema.Set)
		availableZones, err = kafka.GetAvailableZoneIDByCode(cfg, region, azCodes.List())
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		availableZones = utils.ExpandToStringList(d.Get("available_zones").([]interface{}))
	}

	createOpts := &instances.CreateOps{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Engine:              engineRabbitMQ,
		EngineVersion:       d.Get("engine_version").(string),
		StorageSpace:        d.Get("storage_space").(int),
		AccessUser:          d.Get("access_user").(string),
		VPCID:               d.Get("vpc_id").(string),
		SecurityGroupID:     d.Get("security_group_id").(string),
		SubnetID:            d.Get("network_id").(string),
		AvailableZones:      availableZones,
		ProductID:           d.Get("flavor_id").(string),
		BrokerNum:           d.Get("broker_num").(int),
		MaintainBegin:       d.Get("maintain_begin").(string),
		MaintainEnd:         d.Get("maintain_end").(string),
		SslEnable:           d.Get("ssl_enable").(bool),
		StorageSpecCode:     d.Get("storage_spec_code").(string),
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
		EnableAcl:           d.Get("enable_acl").(bool),
		DiskEncryptedEnable: d.Get("disk_encrypted_enable").(bool),
		DiskEncryptedKey:    d.Get("disk_encrypted_key").(string),
	}

	if chargingMode, ok := d.GetOk("charging_mode"); ok && chargingMode == "prePaid" {
		var autoRenew bool
		if d.Get("auto_renew").(string) == "true" {
			autoRenew = true
		}
		isAutoPay := true
		createOpts.BssParam = &instances.BssParam{
			ChargingMode: d.Get("charging_mode").(string),
			PeriodType:   d.Get("period_unit").(string),
			PeriodNum:    d.Get("period").(int),
			IsAutoRenew:  &autoRenew,
			IsAutoPay:    &isAutoPay,
		}
	}

	if pubIpID, ok := d.GetOk("public_ip_id"); ok {
		createOpts.EnablePublicIP = true
		createOpts.PublicIpID = pubIpID.(string)
	}

	// set tags
	if tagRaw := d.Get("tags").(map[string]interface{}); len(tagRaw) > 0 {
		createOpts.Tags = utils.ExpandResourceTags(tagRaw)
	}

	log.Printf("[DEBUG] Create DMS RabbitMQ instance Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	v, err := instances.CreateWithEngine(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating DMS RabbitMQ instance: %s", err)
	}
	instanceID := v.InstanceID
	// Store the instance ID now
	d.SetId(instanceID)
	log.Printf("[INFO] Creating RabbitMQ instance, ID: %s", instanceID)

	var delayTime time.Duration = 300
	if chargingMode, ok := d.GetOk("charging_mode"); ok && chargingMode == "prePaid" {
		err = waitForRabbitMQOrderComplete(ctx, d, cfg, client, instanceID)
		if err != nil {
			return diag.FromErr(err)
		}
		delayTime = 5
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      rabbitmqInstanceStateRefreshFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        delayTime * time.Second,
		PollInterval: 15 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RabbitMQ instance (%s) to be ready: %s", instanceID, err)
	}

	return nil
}

func waitForRabbitMQOrderComplete(ctx context.Context, d *schema.ResourceData, conf *config.Config,
	client *golangsdk.ServiceClient, instanceID string) error {
	region := conf.GetRegion(d)
	orderId, err := getRabbitMQInstanceOrderId(ctx, d, client, instanceID)
	if err != nil {
		return err
	}
	if orderId == "" {
		log.Printf("[WARN] error get order id by instance ID: %s", instanceID)
		return nil
	}

	bssClient, err := conf.BssV2Client(region)
	if err != nil {
		return fmt.Errorf("error creating BSS v2 client: %s", err)
	}
	// wait for order success
	err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error waiting for RabbitMQ order resource %s complete: %s", orderId, err)
	}
	return nil
}

func getRabbitMQInstanceOrderId(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"EMPTY"},
		Target:       []string{"CREATING"},
		Refresh:      rabbitMQInstanceCreatingFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        500 * time.Millisecond,
		PollInterval: 500 * time.Millisecond,
	}
	orderId, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", fmt.Errorf("error waiting for RabbitMQ instance (%s) to creating: %s", instanceID, err)
	}
	return orderId.(string), nil
}

func rabbitMQInstanceCreatingFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return instance, "EMPTY", nil
			}
			return nil, "", err
		}
		return instance.OrderID, "CREATING", nil
	}
}

func createRabbitMQInstanceWithProductID(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DmsV2Client(region)
	if err != nil {
		return diag.Errorf("error initializing DMS RabbitMQ(v2) client: %s", err)
	}
	var availableZones []string // Available zones IDs
	azIDs, ok := d.GetOk("availability_zones")
	if ok {
		// convert the codes of the availability zone into ids
		azCodes := azIDs.(*schema.Set)
		availableZones, err = kafka.GetAvailableZoneIDByCode(cfg, region, azCodes.List())
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		availableZones = utils.ExpandToStringList(d.Get("available_zones").([]interface{}))
	}

	storageSpace := d.Get("storage_space").(int)
	if storageSpace == 0 {
		product, err := getRabbitMQProductDetail(cfg, d)
		if err != nil || product == nil {
			return diag.Errorf("failed to query DMS RabbitMQ product details: %s", err)
		}
		defaultStorageSpace, err := strconv.ParseInt(product.Storage, 10, 32)
		if err != nil {
			return diag.Errorf("failed to create RabbitMQ instance, error parsing storage_space to int %v: %s",
				product.Storage, err)
		}
		storageSpace = int(defaultStorageSpace)
	}

	createOpts := &instances.CreateOps{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Engine:              engineRabbitMQ,
		EngineVersion:       d.Get("engine_version").(string),
		StorageSpace:        storageSpace,
		AccessUser:          d.Get("access_user").(string),
		VPCID:               d.Get("vpc_id").(string),
		SecurityGroupID:     d.Get("security_group_id").(string),
		SubnetID:            d.Get("network_id").(string),
		AvailableZones:      availableZones,
		ProductID:           d.Get("product_id").(string),
		MaintainBegin:       d.Get("maintain_begin").(string),
		MaintainEnd:         d.Get("maintain_end").(string),
		SslEnable:           d.Get("ssl_enable").(bool),
		StorageSpecCode:     d.Get("storage_spec_code").(string),
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
		DiskEncryptedEnable: d.Get("disk_encrypted_enable").(bool),
		DiskEncryptedKey:    d.Get("disk_encrypted_key").(string),
	}

	if pubIpID, ok := d.GetOk("public_ip_id"); ok {
		createOpts.EnablePublicIP = true
		createOpts.PublicIpID = pubIpID.(string)
	}

	// set tags
	if tagRaw := d.Get("tags").(map[string]interface{}); len(tagRaw) > 0 {
		createOpts.Tags = utils.ExpandResourceTags(tagRaw)
	}

	log.Printf("[DEBUG] Create DMS RabbitMQ instance Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	v, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating DMS RabbitMQ instance: %s", err)
	}
	// Store the instance ID now
	d.SetId(v.InstanceID)

	log.Printf("[INFO] Creating RabbitMQ instance, ID: %s", v.InstanceID)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      rabbitmqInstanceStateRefreshFunc(client, v.InstanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        300 * time.Second,
		PollInterval: 15 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RabbitMQ instance (%s) to be ready: %s", v.InstanceID, err)
	}

	return nil
}

func resourceDmsRabbitmqInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.DmsV2Client(region)
	if err != nil {
		return diag.Errorf("error initializing DMS RabbitMQ(v2) client: %s", err)
	}

	v, err := instances.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS RabbitMQ instance")
	}

	log.Printf("[DEBUG] DMS RabbitMQ instance %+v", v)

	azIDs := v.AvailableZones
	asCodes, err := kafka.GetAvailableZoneCodeByID(cfg, region, azIDs)
	mErr := multierror.Append(nil, err)

	var chargingMode = "postPaid"
	if v.ChargingMode == 0 {
		chargingMode = "prePaid"
	}

	d.SetId(v.InstanceID)

	createdAt, _ := strconv.ParseInt(v.CreatedAt, 10, 64)
	mErr = multierror.Append(mErr,
		d.Set("region", region),
		d.Set("name", v.Name),
		d.Set("description", v.Description),
		d.Set("engine", v.Engine),
		d.Set("engine_version", v.EngineVersion),
		d.Set("specification", v.Specification),
		// storage_space indicates total_storage_space while creating
		// set value of total_storage_space to storage_space to keep consistent
		d.Set("storage_space", v.TotalStorageSpace),

		d.Set("vpc_id", v.VPCID),
		d.Set("security_group_id", v.SecurityGroupID),
		d.Set("network_id", v.SubnetID),
		d.Set("available_zones", v.AvailableZones),
		d.Set("availability_zones", asCodes),
		setRabbitMQFlavorId(d, v.ProductID),
		d.Set("maintain_begin", v.MaintainBegin),
		d.Set("maintain_end", v.MaintainEnd),
		d.Set("enable_public_ip", v.EnablePublicIP),
		d.Set("public_ip_id", v.PublicIPID),
		d.Set("ssl_enable", v.SslEnable),
		d.Set("storage_spec_code", v.StorageSpecCode),
		d.Set("broker_num", v.BrokerNum),
		d.Set("enterprise_project_id", v.EnterpriseProjectID),
		d.Set("used_storage_space", v.UsedStorageSpace),
		d.Set("connect_address", v.ConnectAddress),
		d.Set("management_connect_address", v.ManagementConnectAddress),
		d.Set("manegement_connect_address", v.ManagementConnectAddress),
		d.Set("port", v.Port),
		d.Set("status", v.Status),
		d.Set("resource_spec_code", v.ResourceSpecCode),
		d.Set("user_id", v.UserID),
		d.Set("user_name", v.UserName),
		d.Set("type", v.Type),
		d.Set("access_user", v.AccessUser),
		d.Set("charging_mode", chargingMode),
		d.Set("created_at", utils.FormatTimeStampRFC3339(createdAt/1000, false)),
		d.Set("extend_times", v.ExtendTimes),
		d.Set("is_logical_volume", v.IsLogicalVolume),
		d.Set("public_ip_address", v.PublicIPAddress),
		d.Set("enable_acl", v.EnableAcl),
		d.Set("disk_encrypted_enable", v.DiskEncrypted),
		d.Set("disk_encrypted_key", v.DiskEncryptedKey),
	)

	// set tags
	if resourceTags, err := tags.Get(client, engineRabbitMQ, d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		err = d.Set("tags", tagMap)
		if err != nil {
			mErr = multierror.Append(mErr,
				fmt.Errorf("error saving tags to state for DMS RabbitMQ instance (%s): %s", d.Id(), err))
		}
	} else {
		log.Printf("[WARN] error fetching tags of DMS RabbitMQ instance (%s): %s", d.Id(), err)
		e := fmt.Errorf("error fetching tags of DMS RabbitMQ instance (%s): %s", d.Id(), err)
		mErr = multierror.Append(mErr, e)
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("failed to set attributes for DMS RabbitMQ instance: %s", mErr)
	}

	return nil
}

func setRabbitMQFlavorId(d *schema.ResourceData, flavorId string) error {
	re := regexp.MustCompile(`^\d(\d|-)*\d$`)
	if re.MatchString(flavorId) {
		return d.Set("product_id", flavorId)
	}
	return d.Set("flavor_id", flavorId)
}

func resourceDmsRabbitmqInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error initializing DMS RabbitMQ(v2) client: %s", err)
	}

	var mErr *multierror.Error
	if d.HasChanges("name", "description", "maintain_begin", "maintain_end",
		"security_group_id", "enterprise_project_id", "enable_acl") {
		description := d.Get("description").(string)
		updateOpts := instances.UpdateOpts{
			Description:         &description,
			MaintainBegin:       d.Get("maintain_begin").(string),
			MaintainEnd:         d.Get("maintain_end").(string),
			SecurityGroupID:     d.Get("security_group_id").(string),
			EnterpriseProjectID: d.Get("enterprise_project_id").(string),
		}

		if d.HasChange("enable_acl") {
			enableAcl := d.Get("enable_acl").(bool)
			updateOpts.EnableAcl = &enableAcl
		}

		if d.HasChange("name") {
			updateOpts.Name = d.Get("name").(string)
		}

		retryFunc := func() (interface{}, bool, error) {
			err = instances.Update(client, d.Id(), updateOpts).Err
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rabbitmqInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			e := fmt.Errorf("error updating DMS RabbitMQ Instance: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	if d.HasChange("public_ip_id") {
		oldEIP, newEIP := d.GetChange("public_ip_id")
		if oldEIP.(string) != "" {
			// unbind the EIP
			enablePublicIP := false
			updateOpts := instances.UpdateOpts{
				EnablePublicIP: &enablePublicIP,
			}
			err := rabbitmqBindOrUnbindEIP(ctx, client, d.Timeout(schema.TimeoutUpdate), updateOpts, d.Id(), "unbindInstancePublicIp")
			if err != nil {
				mErr = multierror.Append(mErr, err)
			}
		}
		if newEIP.(string) != "" {
			// bind the new EIP
			enablePublicIP := true
			updateOpts := instances.UpdateOpts{
				EnablePublicIP: &enablePublicIP,
				PublicIpID:     newEIP.(string),
			}
			err := rabbitmqBindOrUnbindEIP(ctx, client, d.Timeout(schema.TimeoutUpdate), updateOpts, d.Id(), "bindInstancePublicIp")
			if err != nil {
				mErr = multierror.Append(mErr, err)
			}
		}
	}

	if d.HasChange("tags") {
		// update tags
		tagErr := utils.UpdateResourceTags(client, d, engineRabbitMQ, d.Id())
		if tagErr != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("error updating tags of DMS RabbitMQ instance: %s, err: %s",
				d.Id(), tagErr))
		}
	}

	if d.HasChanges("product_id", "flavor_id", "broker_num", "storage_space") {
		err = resizeRabbitMQInstance(ctx, d, meta, engineRabbitMQ)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if d.HasChange("password") {
		resetPasswordOpts := instances.ResetPasswordOpts{
			NewPassword: d.Get("password").(string),
		}
		retryFunc := func() (interface{}, bool, error) {
			err = instances.ResetPassword(client, d.Id(), resetPasswordOpts).Err
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rabbitmqInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			e := fmt.Errorf("error resetting password: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error while updating DMS RabbitMQ instances, %s", mErr)
	}

	return resourceDmsRabbitmqInstanceRead(ctx, d, meta)
}

func resizeRabbitMQInstance(ctx context.Context, d *schema.ResourceData, meta interface{}, engineType string) error {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error initializing DMS(v2) client: %s", err)
	}

	if d.HasChanges("product_id") {
		product, err := getRabbitMQProductDetail(cfg, d)
		if err != nil {
			return fmt.Errorf("failed to resize RabbitMQ instance, query product details error: %s", err)
		}
		storage, err := strconv.Atoi(product.Storage)
		if err != nil {
			return fmt.Errorf("failed to resize RabbitMQ instance, error parsing storage_space to int %v: %s",
				product.Storage, err)
		}

		resizeOpts := instances.ResizeInstanceOpts{
			NewSpecCode:     &product.SpecCode,
			NewStorageSpace: &storage,
		}
		log.Printf("[DEBUG] Resize DMS RabbitMQ instance option : %#v", resizeOpts)

		if err = doRabbitMQInstanceResize(ctx, d, client, resizeOpts); err != nil {
			return err
		}
	}

	if d.HasChanges("flavor_id") {
		flavorID := d.Get("flavor_id").(string)
		operType := "vertical"
		resizeOpts := instances.ResizeInstanceOpts{
			OperType:     &operType,
			NewProductID: &flavorID,
		}
		log.Printf("[DEBUG] Resize RabbitMQ instance flavor ID options: %s", utils.MarshalValue(resizeOpts))

		if err = doRabbitMQInstanceResize(ctx, d, client, resizeOpts); err != nil {
			return err
		}
	}

	if d.HasChanges("broker_num") {
		brokerNum := d.Get("broker_num").(int)
		operType := "horizontal"
		resizeOpts := instances.ResizeInstanceOpts{
			OperType:     &operType,
			NewBrokerNum: &brokerNum,
		}
		log.Printf("[DEBUG] Resize RabbitMQ instance broker num options: %s", utils.MarshalValue(resizeOpts))

		if err = doRabbitMQInstanceResize(ctx, d, client, resizeOpts); err != nil {
			return err
		}
	}

	if d.HasChanges("storage_space") {
		storageSpace := d.Get("storage_space").(int)
		operType := "storage"
		resizeOpts := instances.ResizeInstanceOpts{
			OperType:        &operType,
			NewStorageSpace: &storageSpace,
		}
		log.Printf("[DEBUG] Resize RabbitMQ instance storage space options: %s", utils.MarshalValue(resizeOpts))

		if err = doRabbitMQInstanceResize(ctx, d, client, resizeOpts); err != nil {
			return err
		}
	}

	return nil
}

func doRabbitMQInstanceResize(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	opts instances.ResizeInstanceOpts) error {
	retryFunc := func() (interface{}, bool, error) {
		_, err := instances.Resize(client, d.Id(), opts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rabbitmqInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("resize RabbitMQ instance failed: resizeInstanceOpts: %#v, err: %s", opts, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"RUNNING"},
		Refresh:      rabbitMQResizeStateRefresh(client, d, opts.OperType),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        60 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for instance (%s) to resize: %v", d.Id(), err)
	}
	return nil
}

func rabbitMQResizeStateRefresh(client *golangsdk.ServiceClient, d *schema.ResourceData, operType *string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.Get(client, d.Id()).Extract()
		if err != nil {
			return nil, "failed", err
		}

		if v.Task.Status != "" && v.Task.Status != "SUCCESS" {
			return v, "PENDING", nil
		}

		return v, v.Status, nil
	}
}

func resourceDmsRabbitmqInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error initializing DMS RabbitMQ(v2) client: %s", err)
	}

	if d.Get("charging_mode") == "prePaid" {
		retryFunc := func() (interface{}, bool, error) {
			err = common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()})
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rabbitmqInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error unsubscribe RabbitMQ instance: %s", err)
		}
	} else {
		retryFunc := func() (interface{}, bool, error) {
			err = instances.Delete(client, d.Id()).ExtractErr()
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rabbitmqInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutDelete),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return common.CheckDeletedDiag(d, err, "failed to delete RabbitMQ instance")
		}
	}

	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for DMS RabbitMQ instance (%s) to be deleted", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETING", "RUNNING", "ERROR"}, // Status may change to ERROR on deletion.
		Target:       []string{"DELETED"},
		Refresh:      rabbitmqInstanceStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        90 * time.Second,
		PollInterval: 15 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for DMS RabbitMQ instance (%s) to be deleted: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] DMS RabbitMQ instance %s has been deleted", d.Id())
	d.SetId("")
	return nil
}

func rabbitmqInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
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

func rabbitmqBindOrUnbindEIP(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	updateOpts instances.UpdateOpts, id, action string) error {
	retryFunc := func() (interface{}, bool, error) {
		err := instances.Update(client, id, updateOpts).Err
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rabbitmqInstanceStateRefreshFunc(client, id),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      timeout,
		DelayTimeout: 5 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating DMS RabbitMQ Instance with action(%s): %s", action, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATED"},
		Target:       []string{"SUCCESS"},
		Refresh:      kafka.FilterTaskRefreshFunc(client, id, action),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 15 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for job(%s) success: %s", action, err)
	}

	return nil
}
