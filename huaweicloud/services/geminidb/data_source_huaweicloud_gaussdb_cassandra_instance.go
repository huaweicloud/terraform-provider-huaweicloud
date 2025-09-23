package geminidb

import (
	"context"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/geminidb/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforNoSQL GET /v3/{project_id}/instances
// @API GaussDBforNoSQL GET /v3/{project_id}/instances/{instance_id}/tags
func DataSourceGeminiDBInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBInstanceRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func dataSourceGeminiDBInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.GeminiDBV3Client(region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listOpts := instances.ListGeminiDBInstanceOpts{
		Name:     d.Get("name").(string),
		VpcId:    d.Get("vpc_id").(string),
		SubnetId: d.Get("subnet_id").(string),
	}

	pages, err := instances.List(client, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("error getting GaussDB cassandra instance list: %s", err)
	}

	allInstances, err := instances.ExtractGeminiDBInstances(pages)
	if err != nil {
		return diag.Errorf("unable to retrieve instances: %s", err)
	}

	if allInstances.TotalCount < 1 {
		return diag.Errorf("your query returned no results. " +
			"please change your search criteria and try again.")
	}

	if allInstances.TotalCount > 1 {
		return diag.Errorf("your query returned more than one result." +
			" please try a more specific search criteria")
	}

	instance := allInstances.Instances[0]

	log.Printf("[DEBUG] retrieved instance %s: %+v", instance.Id, instance)
	d.SetId(instance.Id)

	mErr := multierror.Append(
		d.Set("name", instance.Name),
		d.Set("region", instance.Region),
		d.Set("status", instance.Status),
		d.Set("vpc_id", instance.VpcId),
		d.Set("subnet_id", instance.SubnetId),
		d.Set("security_group_id", instance.SecurityGroupId),
		d.Set("enterprise_project_id", instance.EnterpriseProjectId),
		d.Set("mode", instance.Mode),
		d.Set("db_user_name", instance.DbUserName),
	)
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
			log.Printf("[DEBUG] node specCode: %s", specCode)
			d.Set("flavor", specCode)
		}
	}
	d.Set("nodes", nodesList)
	d.Set("private_ips", ipsList)

	// remove duplicate az
	azList = utils.RemoveDuplicateElem(azList)
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

	// save geminidb tags
	if resourceTags, err := tags.Get(client, "instances", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return diag.Errorf("error saving tags to state for geminidb (%s): %s", d.Id(), err)
		}
	} else {
		log.Printf("[WARN] error fetching tags of geminidb (%s): %s", d.Id(), err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}
