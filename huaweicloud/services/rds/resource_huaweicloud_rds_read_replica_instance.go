package rds

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourceRdsReadReplicaInstance is the impl for huaweicloud_rds_read_replica_instance resource
// @API RDS POST /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/jobs
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/alias
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/port
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/security-group
// @API RDS POST /v3/{project_id}/instances/{instance_id}/action
// @API RDS POST /v3/{project_id}/instances/{id}/tags/action
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/ip
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/ssl
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/disk-auto-expansion
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/configurations
// @API RDS PUT /v3.1/{project_id}/instances/{instance_id}/configurations
// @API RDS GET /v3/{project_id}/instances/{instance_id}/disk-auto-expansion
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/name
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceRdsReadReplicaInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsReadReplicaInstanceCreate,
		ReadContext:   resourceRdsReadReplicaInstanceRead,
		UpdateContext: resourceRdsReadReplicaInstanceUpdate,
		DeleteContext: resourceRdsInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"primary_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"limit_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"trigger_threshold": {
							Type:         schema.TypeInt,
							Optional:     true,
							RequiredWith: []string{"volume.0.limit_size"},
						},
						"disk_encryption_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "schema: Computed",
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fixed_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
			},
			"maintain_begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_end": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"maintain_begin"},
			},
			"db": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			// charge info: charging_mode, period_unit, period, auto_renew
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceRdsReadReplicaInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.RdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating rds client: %s ", err)
	}

	primaryInstanceID := d.Get("primary_instance_id").(string)
	createOpts := instances.CreateReplicaOpts{
		Name:                d.Get("name").(string),
		ReplicaOfId:         primaryInstanceID,
		FlavorRef:           d.Get("flavor").(string),
		Region:              region,
		AvailabilityZone:    d.Get("availability_zone").(string),
		Volume:              buildRdsReplicaInstanceVolume(d),
		DiskEncryptionId:    d.Get("volume.0.disk_encryption_id").(string),
		EnterpriseProjectId: config.GetEnterpriseProjectID(d),
	}

	// PrePaid
	if d.Get("charging_mode") == "prePaid" {
		if err = common.ValidatePrePaidChargeInfo(d); err != nil {
			return diag.FromErr(err)
		}

		chargeInfo := &instances.ChargeInfo{
			ChargeMode: d.Get("charging_mode").(string),
			PeriodType: d.Get("period_unit").(string),
			PeriodNum:  d.Get("period").(int),
			IsAutoPay:  true,
		}
		if d.Get("auto_renew").(string) == "true" {
			chargeInfo.IsAutoRenew = true
		}
		createOpts.ChargeInfo = chargeInfo
	}

	log.Printf("[DEBUG] Create replica instance Options: %#v", createOpts)
	retryFunc := func() (interface{}, bool, error) {
		resp, err := instances.CreateReplica(client, createOpts).Extract()
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, primaryInstanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating replica instance: %s ", err)
	}

	resp := r.(*instances.CreateResponse)

	instance := resp.Instance
	d.SetId(instance.Id)
	instanceID := d.Id()
	// wait for order success
	if resp.OrderId != "" {
		bssClient, err := config.BssV2Client(config.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, resp.OrderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, resp.OrderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for replica order resource %s complete: %s", resp.OrderId, err)
		}
		d.SetId(resourceId)
	} else {
		if err := checkRDSInstanceJobFinish(client, resp.JobId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error creating replica instance (%s): %s", instanceID, err)
		}
	}

	res, err := GetRdsInstanceByID(client, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceDescription(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceMaintainWindow(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	port := utils.PathSearch("port", res, float64(0)).(float64)
	if v, ok := d.GetOk("db.0.port"); ok && v.(int) != int(port) {
		if err = updateRdsInstanceDBPort(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	securityGroupId := utils.PathSearch("security_group_id", res, "").(string)
	if v, ok := d.GetOk("security_group_id"); ok && v.(string) != securityGroupId {
		if err = updateRdsInstanceSecurityGroup(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	volumeSize := utils.PathSearch("volume.size", res, float64(0)).(float64)
	if v, ok := d.GetOk("volume.0.size"); ok && v.(int) != int(volumeSize) {
		if err = updateRdsInstanceVolumeSize(ctx, d, config, client); err != nil {
			return diag.FromErr(err)
		}
	}

	fixedIp := utils.PathSearch("private_ips[0]", res, "").(string)
	if v, ok := d.GetOk("fixed_ip"); ok && v.(string) != fixedIp {
		if err = updateRdsInstanceFixedIp(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	enableSsl := utils.PathSearch("enable_ssl", res, false).(bool)
	if v, ok := d.GetOk("ssl_enable"); ok && v.(bool) != enableSsl {
		dataStoreType := utils.PathSearch("datastore.type", res, "").(string)
		if strings.ToLower(dataStoreType) != "mysql" {
			return diag.Errorf("only MySQL database support SSL enable and disable")
		}
		err = configRdsInstanceSSL(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("volume.0.limit_size"); ok {
		if v.(int) > 0 {
			if err = enableVolumeAutoExpand(ctx, d, client, v.(int)); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err = disableVolumeAutoExpand(ctx, schema.TimeoutCreate, client, d); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// Set Parameters
	if parametersRaw := d.Get("parameters").(*schema.Set); parametersRaw.Len() > 0 {
		if err = initializeParameters(ctx, d, client, parametersRaw); err != nil {
			return diag.FromErr(err)
		}
	}

	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := utils.ExpandResourceTags(tagRaw)
		err := tags.Create(client, "instances", instanceID, tagList).ExtractErr()
		if err != nil {
			return diag.Errorf("error setting tags of RDS read replica instance %s: %s", instanceID, err)
		}
	}

	return resourceRdsReadReplicaInstanceRead(ctx, d, meta)
}

func resourceRdsReadReplicaInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating rds client: %s", err)
	}

	instanceID := d.Id()
	instance, err := GetRdsInstanceByID(client, instanceID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting RDS replice instance")
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", instance, nil)),
		d.Set("description", utils.PathSearch("alias", instance, nil)),
		d.Set("flavor", utils.PathSearch("flavor_ref", instance, nil)),
		d.Set("region", utils.PathSearch("region", instance, nil)),
		d.Set("private_ips", utils.PathSearch("private_ips", instance, nil)),
		d.Set("public_ips", utils.PathSearch("public_ips", instance, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", instance, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", instance, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", instance, nil)),
		d.Set("type", utils.PathSearch("type", instance, nil)),
		d.Set("status", utils.PathSearch("status", instance, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", instance, nil)),
		d.Set("ssl_enable", utils.PathSearch("enable_ssl", instance, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", instance, make([]interface{}, 0)))),
		d.Set("fixed_ip", utils.PathSearch("private_ips[0]", instance, nil)),
		d.Set("availability_zone", utils.PathSearch("nodes[0].availability_zone", instance, nil)),
		d.Set("primary_instance_id", utils.PathSearch("related_instance[?type=='replica_of']|[0].id", instance, nil)),
		d.Set("volume", flattenInstanceVolume(client, instance, instanceID)),
		d.Set("db", flattenReplicaInstanceDb(instance)),
	)

	if v := utils.PathSearch("maintenance_window", instance, "").(string); v != "" {
		maintainWindow := strings.Split(v, "-")
		mErr = multierror.Append(mErr, d.Set("maintain_begin", maintainWindow[0]))
		mErr = multierror.Append(mErr, d.Set("maintain_end", maintainWindow[1]))
	}

	diagErr := setRdsInstanceParameters(ctx, d, client)
	resErr := append(diag.FromErr(mErr.ErrorOrNil()), diagErr...)

	return resErr
}

func flattenReplicaInstanceDb(instance interface{}) []interface{} {
	database := map[string]interface{}{
		"type":      utils.PathSearch("datastore.type", instance, nil),
		"version":   utils.PathSearch("datastore.version", instance, nil),
		"port":      utils.PathSearch("port", instance, nil),
		"user_name": utils.PathSearch("db_user_name", instance, nil),
	}
	return []interface{}{database}
}

func resourceRdsReadReplicaInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.RdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating rds v3 client: %s ", err)
	}

	instanceID := d.Id()

	if err = updateRdsInstanceName(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceDescription(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceVolumeSize(ctx, d, cfg, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceMaintainWindow(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceDBPort(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceFixedIp(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceSecurityGroup(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceSSLConfig(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceFlavor(ctx, d, cfg, client, false); err != nil {
		return diag.FromErr(err)
	}

	if err = updateRdsInstanceAutoRenew(d, cfg); err != nil {
		return diag.FromErr(err)
	}

	if ctx, err = updateRdsParameters(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if err = updateVolumeAutoExpand(ctx, d, client); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "instances", instanceID)
		if tagErr != nil {
			return diag.Errorf("error updating tags of RDS read replica instance: %s, err: %s", instanceID, tagErr)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceID,
			ResourceType: "rds",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceRdsReadReplicaInstanceRead(ctx, d, meta)
}

func updateRdsInstanceAutoRenew(d *schema.ResourceData, config *config.Config) error {
	if d.HasChange("auto_renew") {
		bssClient, err := config.BssV2Client(config.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return fmt.Errorf("error updating the auto-renew of the instance (%s): %s", d.Id(), err)
		}
	}
	return nil
}

func buildRdsReplicaInstanceVolume(d *schema.ResourceData) *instances.Volume {
	var volume *instances.Volume
	volumeRaw := d.Get("volume").([]interface{})

	if len(volumeRaw) == 1 {
		volume = new(instances.Volume)
		volume.Type = volumeRaw[0].(map[string]interface{})["type"].(string)
		volume.Size = volumeRaw[0].(map[string]interface{})["size"].(int)
		// the size is optional and invalid for replica, but it's required in sdk
		// so just set 100 if not specified
		if volume.Size == 0 {
			volume.Size = 100
		}
	}
	log.Printf("[DEBUG] volume: %+v", volume)
	return volume
}
