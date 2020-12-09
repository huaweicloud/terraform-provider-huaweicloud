package huaweicloud

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/geminidb/v3/instances"
)

func dataSourceGeminiDBInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGeminiDBInstanceRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"backup_strategy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"keep_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"node_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_reduce": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceGeminiDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)
	client, err := config.GeminiDBV3Client(region)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	listOpts := instances.ListGeminiDBInstanceOpts{
		Name:     d.Get("name").(string),
		VpcId:    d.Get("vpc_id").(string),
		SubnetId: d.Get("subnet_id").(string),
	}

	pages, err := instances.List(client, listOpts).AllPages()
	if err != nil {
		return err
	}

	allInstances, err := instances.ExtractGeminiDBInstances(pages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve instances: %s", err)
	}

	if allInstances.TotalCount < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if allInstances.TotalCount > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	instance := allInstances.Instances[0]

	log.Printf("[DEBUG] Retrieved Instance %s: %+v", instance.Id, instance)
	d.SetId(instance.Id)

	d.Set("name", instance.Name)
	d.Set("region", instance.Region)
	d.Set("status", instance.Status)
	d.Set("vpc_id", instance.VpcId)
	d.Set("subnet_id", instance.SubnetId)
	d.Set("security_group_id", instance.SecurityGroupId)
	d.Set("enterprise_project_id", instance.EnterpriseProjectId)
	d.Set("mode", instance.Mode)
	d.Set("db_user_name", instance.DbUserName)

	if dbPort, err := strconv.Atoi(instance.Port); err == nil {
		d.Set("port", dbPort)
	}

	dbList := make([]map[string]interface{}, 0, 1)
	db := map[string]interface{}{
		"engine":         instance.DataStore.Type,
		"version":        instance.DataStore.Version,
		"storage_engine": instance.Engine,
	}
	dbList = append(dbList, db)
	d.Set("datastore", dbList)

	specCode := ""
	wrongFlavor := "Inconsistent Flavor"
	ipsList := []string{}
	azList := []string{}
	nodesList := make([]map[string]interface{}, 0, 1)
	for _, group := range instance.Groups {
		for _, Node := range group.Nodes {
			node := map[string]interface{}{
				"id":                Node.Id,
				"name":              Node.Name,
				"status":            Node.Status,
				"private_ip":        Node.PrivateIp,
				"support_reduce":    Node.SupportReduce,
				"availability_zone": Node.AvailabilityZone,
			}
			if specCode == "" {
				specCode = Node.SpecCode
			} else if specCode != Node.SpecCode && specCode != wrongFlavor {
				specCode = wrongFlavor
			}
			nodesList = append(nodesList, node)
			azList = append(azList, Node.AvailabilityZone)
			// Only return Node private ips which doesn't support reduce
			if !Node.SupportReduce {
				ipsList = append(ipsList, Node.PrivateIp)
			}
		}
		if volSize, err := strconv.Atoi(group.Volume.Size); err == nil {
			d.Set("volume_size", volSize)
		}
		if specCode != "" {
			log.Printf("[DEBUG] Node SpecCode: %s", specCode)
			d.Set("flavor", specCode)
		}
	}
	d.Set("nodes", nodesList)
	d.Set("private_ips", ipsList)

	//remove duplicate az
	azList = removeDuplicateElem(azList)
	sort.Strings(azList)
	d.Set("availability_zone", strings.Join(azList, ","))
	d.Set("node_num", len(nodesList))

	backupStrategyList := make([]map[string]interface{}, 0, 1)
	backupStrategy := map[string]interface{}{
		"start_time": instance.BackupStrategy.StartTime,
		"keep_days":  instance.BackupStrategy.KeepDays,
	}
	backupStrategyList = append(backupStrategyList, backupStrategy)
	d.Set("backup_strategy", backupStrategyList)

	//save geminidb tags
	resourceTags, err := tags.Get(client, "instances", d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching HuaweiCloud geminidb tags: %s", err)
	}

	tagmap := tagsToMap(resourceTags.Tags)
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("Error saving tags for HuaweiCloud geminidb (%s): %s", d.Id(), err)
	}

	return nil
}
