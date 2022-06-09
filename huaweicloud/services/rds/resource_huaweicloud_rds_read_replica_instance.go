package rds

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourceRdsReadReplicaInstance is the impl for huaweicloud_rds_read_replica_instance resource
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
				ForceNew: true,
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
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"disk_encryption_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

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

			"security_group_id": {
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

			"db": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"tags": common.TagsSchema(),
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

	createOpts := instances.CreateReplicaOpts{
		Name:                d.Get("name").(string),
		ReplicaOfId:         d.Get("primary_instance_id").(string),
		FlavorRef:           d.Get("flavor").(string),
		Region:              region,
		AvailabilityZone:    d.Get("availability_zone").(string),
		Volume:              buildRdsReplicaInstanceVolume(d),
		DiskEncryptionId:    d.Get("volume.0.disk_encryption_id").(string),
		EnterpriseProjectId: config.GetEnterpriseProjectID(d),
	}

	log.Printf("[DEBUG] Create replica instance Options: %#v", createOpts)
	resp, err := instances.CreateReplica(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating replica instance: %s ", err)
	}

	instance := resp.Instance
	d.SetId(instance.Id)
	instanceID := d.Id()
	if err := checkRDSInstanceJobFinish(client, resp.JobId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error creating replica instance (%s): %s", instanceID, err)
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
		return diag.FromErr(err)
	}
	if instance.Id == "" {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Retrieved rds read replica instance %s: %#v", instanceID, instance)
	d.Set("name", instance.Name)
	d.Set("flavor", instance.FlavorRef)
	d.Set("region", instance.Region)
	d.Set("private_ips", instance.PrivateIps)
	d.Set("public_ips", instance.PublicIps)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("type", instance.Type)
	d.Set("status", instance.Status)
	d.Set("enterprise_project_id", instance.EnterpriseProjectId)
	d.Set("tags", utils.TagsToMap(instance.Tags))

	az := expandAvailabilityZone(instance)
	d.Set("availability_zone", az)

	if primaryInstanceID, err := expandPrimaryInstanceID(instance); err == nil {
		d.Set("primary_instance_id", primaryInstanceID)
	} else {
		return diag.FromErr(err)
	}

	volumeList := make([]map[string]interface{}, 0, 1)
	volume := map[string]interface{}{
		"type":               instance.Volume.Type,
		"size":               instance.Volume.Size,
		"disk_encryption_id": instance.DiskEncryptionId,
	}
	volumeList = append(volumeList, volume)
	if err := d.Set("volume", volumeList); err != nil {
		return diag.Errorf("error saving volume to RDS read replica instance (%s): %s", instanceID, err)
	}

	dbList := make([]map[string]interface{}, 0, 1)
	database := map[string]interface{}{
		"type":      instance.DataStore.Type,
		"version":   instance.DataStore.Version,
		"port":      instance.Port,
		"user_name": instance.DbUserName,
	}
	dbList = append(dbList, database)
	if err := d.Set("db", dbList); err != nil {
		return diag.Errorf("error saving data base to RDS read replica instance (%s): %s", instanceID, err)
	}

	return nil
}

func resourceRdsReadReplicaInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating rds v3 client: %s ", err)
	}

	instanceID := d.Id()
	if err := updateRdsInstanceFlavor(d, config, client, instanceID, false); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "instances", instanceID)
		if tagErr != nil {
			return diag.Errorf("error updating tags of RDS read replica instance: %s, err: %s", instanceID, tagErr)
		}
	}

	return resourceRdsReadReplicaInstanceRead(ctx, d, meta)
}

func expandAvailabilityZone(resp *instances.RdsInstanceResponse) string {
	node := resp.Nodes[0]
	return node.AvailabilityZone
}

func expandPrimaryInstanceID(resp *instances.RdsInstanceResponse) (string, error) {
	relatedInst := resp.RelatedInstance
	for _, relate := range relatedInst {
		if relate.Type == "replica_of" {
			return relate.Id, nil
		}
	}
	return "", fmt.Errorf("error when get primary instance id for replica %s", resp.Id)
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
