package dms

import (
	"context"
	"strconv"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dms/v2/products"
	"github.com/chnsz/golangsdk/openstack/dms/v2/rabbitmq/instances"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

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
				Optional: true,
				ForceNew: true,
				Default:  "3.7.17",
			},
			"storage_space": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
				ForceNew:  true,
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
				// There is a problem with order of elements in Availability Zone list returned by RabbitMQ API.
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"available_zones"},
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
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
		},
	}
}

func getRabbitMQProductDetail(config *config.Config, d *schema.ResourceData) (*products.Detail, error) {
	productRsp, err := getProducts(config, config.GetRegion(d), "rabbitmq")
	if err != nil {
		return nil, fmtp.Errorf("error querying product detail, please check product_id, error: %s", err)
	}

	productID := d.Get("product_id").(string)
	engineVersion := d.Get("engine_version").(string)

	for _, ps := range productRsp.Hourly {
		if ps.Version != engineVersion {
			continue
		}
		for _, v := range ps.Values {
			for _, p := range v.Details {
				// All informations of product for single instance type and the kafka engine type are stored in the
				// detail structure.
				if v.Name == "single" {
					if p.ProductID == productID {
						return &p, nil
					}
				} else {
					for _, pi := range p.ProductInfos {
						if pi.ProductID == productID {
							p.ProductInfos = []products.ProductInfo{pi}
							return &p, nil
						}
					}
				}

			}
		}
	}
	return nil, fmtp.Errorf("can not found product detail base on product_id: %s", productID)
}

func resourceDmsRabbitmqInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	dmsV2Client, err := config.DmsV2Client(region)
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DMS instance client: %s", err)
	}

	var availableZones []string
	zoneIDs, ok := d.GetOk("available_zones")
	if ok {
		availableZones = utils.ExpandToStringList(zoneIDs.([]interface{}))
	} else {
		// convert the codes of the availability zone into ids
		azCodes := d.Get("availability_zones").(*schema.Set)
		availableZones, err = getAvailableZoneIDByCode(config, region, azCodes.List())
		if err != nil {
			return diag.FromErr(err)
		}
	}

	storageSpace := d.Get("storage_space").(int)
	if storageSpace == 0 {
		product, err := getRabbitMQProductDetail(config, d)
		if err != nil || product == nil {
			return fmtp.DiagErrorf("query DMS RabbimtMQ product failed: %s", err)
		}
		space := product.ProductInfos[0].Storage
		defaultStorageSpace, err := strconv.ParseInt(space, 10, 32)
		if err != nil {
			return fmtp.DiagErrorf("Parse storage capacity to int error, %v: %s", space, err)
		}
		storageSpace = int(defaultStorageSpace)
	}

	createOpts := &instances.CreateOps{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Engine:              "rabbitmq",
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
		EnterpriseProjectID: common.GetEnterpriseProjectID(d, config),
	}

	if v, ok := d.GetOk("public_ip_id"); ok {
		createOpts.EnablePublicIP = true
		createOpts.PublicIpID = v.(string)
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		createOpts.Tags = taglist
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	v, err := instances.Create(dmsV2Client, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DMS rabbitmq instance: %s", err)
	}
	logp.Printf("[INFO] instance ID: %s", v.InstanceID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      DmsRabbitmqInstanceStateRefreshFunc(dmsV2Client, v.InstanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        500 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf(
			"error waiting for instance (%s) to become ready: %s",
			v.InstanceID, err)
	}

	// Store the instance ID now
	d.SetId(v.InstanceID)

	return resourceDmsRabbitmqInstanceRead(ctx, d, meta)
}

func resourceDmsRabbitmqInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	logp.Printf("[DEBUG] DMS rabbitmq instance %s: %+v", d.Id(), v)

	availableZoneIDs := v.AvailableZones
	availableZoneCodes, err := getAvailableZoneCodeByID(config, region, availableZoneIDs)
	mErr := multierror.Append(nil, err)

	d.SetId(v.InstanceID)
	mErr = multierror.Append(mErr,
		d.Set("region", config.GetRegion(d)),
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
		d.Set("availability_zones", availableZoneCodes),
		d.Set("product_id", v.ProductID),
		d.Set("maintain_begin", v.MaintainBegin),
		d.Set("maintain_end", v.MaintainEnd),
		d.Set("enable_public_ip", v.EnablePublicIP),
		d.Set("public_ip_id", v.PublicIPID),
		d.Set("ssl_enable", v.SslEnable),
		d.Set("storage_spec_code", v.StorageSpecCode),
		d.Set("enterprise_project_id", v.EnterpriseProjectID),
		d.Set("used_storage_space", v.UsedStorageSpace),
		d.Set("connect_address", v.ConnectAddress),
		d.Set("manegement_connect_address", v.ManagementConnectAddress),
		d.Set("port", v.Port),
		d.Set("status", v.Status),
		d.Set("resource_spec_code", v.ResourceSpecCode),
		d.Set("user_id", v.UserID),
		d.Set("user_name", v.UserName),
		d.Set("type", v.Type),
		d.Set("access_user", v.AccessUser),
	)

	// set tags
	engine := "rabbitmq"
	if resourceTags, err := tags.Get(dmsV2Client, engine, d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		err = d.Set("tags", tagmap)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	} else {
		logp.Printf("[WARN] error fetching tags of DMS rabbitmq instance (%s): %s", d.Id(), err)
		e := fmtp.Errorf("error fetching tags of DMS rabbitmq instance (%s): %s", d.Id(), err)
		mErr = multierror.Append(mErr, e)
	}

	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("Error setting attributes for DMS rabbitmq instance: %s", mErr)
	}

	return nil
}

func resourceDmsRabbitmqInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DMS instance client: %s", err)
	}

	var mErr *multierror.Error
	if d.HasChanges("name", "description", "maintain_begin", "maintain_end",
		"security_group_id", "public_ip_id", "enterprise_project_id") {
		description := d.Get("description").(string)
		updateOpts := instances.UpdateOpts{
			Description:         &description,
			MaintainBegin:       d.Get("maintain_begin").(string),
			MaintainEnd:         d.Get("maintain_end").(string),
			SecurityGroupID:     d.Get("security_group_id").(string),
			EnterpriseProjectID: d.Get("enterprise_project_id").(string),
		}

		if d.HasChange("name") {
			updateOpts.Name = d.Get("name").(string)
		}

		if d.HasChange("public_ip_id") {
			if v, ok := d.GetOk("public_ip_id"); ok {
				enablePublicIP := true
				updateOpts.EnablePublicIP = &enablePublicIP
				updateOpts.PublicIpID = v.(string)
			} else {
				enablePublicIP := false
				updateOpts.EnablePublicIP = &enablePublicIP
			}
		}

		err = instances.Update(dmsV2Client, d.Id(), updateOpts).Err
		if err != nil {
			e := fmtp.Errorf("error updating HuaweiCloud DMS rabbitMQ Instance: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	if d.HasChange("tags") {
		// update tags
		engine := "rabbitmq"
		tagErr := utils.UpdateResourceTags(dmsV2Client, d, engine, d.Id())
		if tagErr != nil {
			e := fmtp.Errorf("error updating tags of DMS rabbitmq instance:%s, err:%s", d.Id(), tagErr)
			mErr = multierror.Append(mErr, e)
		}
	}

	if d.HasChange("product_id") {
		err = resizeInstance(ctx, d, meta, "rabbitmq")
		if err != nil {
			e := fmtp.Errorf("error resizing HuaweiCloud DMS rabbitMQ Instance: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}
	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("error while updating DMS rabbitMQ instances, there %s", mErr)
	}

	return resourceDmsRabbitmqInstanceRead(ctx, d, meta)
}

func resourceDmsRabbitmqInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		Refresh:      DmsRabbitmqInstanceStateRefreshFunc(dmsV2Client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        90 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf(
			"error waiting for instance (%s) to delete: %s", d.Id(), err)
	}

	logp.Printf("[DEBUG] DMS instance %s deactivated", d.Id())
	d.SetId("")
	return nil
}

func DmsRabbitmqInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
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
