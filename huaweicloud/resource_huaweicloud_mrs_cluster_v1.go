package huaweicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/mrs/v1/cluster"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/subnets"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/vpcs"
)

func resourceMRSClusterV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterV1Create,
		Read:   resourceClusterV1Read,
		Update: resourceClusterV1Update,
		Delete: resourceClusterV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"available_zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"billing_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"master_node_num": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 2),
			},
			"master_node_size": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"core_node_num": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 500),
			},
			"core_node_size": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

			"volume_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SATA", "SAS", "SSD",
				}, false),
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"safe_mode": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"cluster_admin_secret": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(8, 26),
			},
			"node_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				ExactlyOneOf: []string{"node_public_cert_name"},
				ValidateFunc: validation.StringLenBetween(8, 26),
			},
			"node_public_cert_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"node_password"},
			},
			"log_collection": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"component_list": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"component_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"component_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"add_jobs": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_type": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"job_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"jar_path": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"arguments": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"input": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"output": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"job_log": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"shutdown_cluster": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"file_action": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"submit_job_once_cluster_run": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
						"hql": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"hive_script_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"tags": tagsSchema(),
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hadoop_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_node_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_first": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slave_security_groups_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_groups_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_alternate_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_node_spec_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"core_node_spec_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_node_product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"core_node_product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vnc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fee": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charging_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getAllClusterComponents(d *schema.ResourceData) []cluster.ComponentOpts {
	var componentOpts []cluster.ComponentOpts

	components := d.Get("component_list").([]interface{})
	for _, v := range components {
		component := v.(map[string]interface{})
		component_name := component["component_name"].(string)

		v := cluster.ComponentOpts{
			ComponentName: component_name,
		}
		componentOpts = append(componentOpts, v)
	}

	log.Printf("[DEBUG] getAllClusterComponents: %#v", componentOpts)
	return componentOpts
}

func getAllClusterJobs(d *schema.ResourceData) []cluster.JobOpts {
	var jobOpts []cluster.JobOpts

	jobs := d.Get("add_jobs").([]interface{})
	for _, v := range jobs {
		job := v.(map[string]interface{})

		v := cluster.JobOpts{
			JobType:                 job["job_type"].(int),
			JobName:                 job["job_name"].(string),
			JarPath:                 job["jar_path"].(string),
			Arguments:               job["arguments"].(string),
			Input:                   job["input"].(string),
			Output:                  job["output"].(string),
			JobLog:                  job["job_log"].(string),
			ShutdownCluster:         job["shutdown_cluster"].(bool),
			FileAction:              job["file_action"].(string),
			SubmitJobOnceClusterRun: job["submit_job_once_cluster_run"].(bool),
			Hql:                     job["hql"].(string),
			HiveScriptPath:          job["hive_script_path"].(string),
		}
		jobOpts = append(jobOpts, v)
	}

	log.Printf("[DEBUG] getAllClusterJobs: %#v", jobOpts)
	return jobOpts
}

func ClusterStateRefreshFunc(client *golangsdk.ServiceClient, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		clusterGet, err := cluster.Get(client, clusterID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return clusterGet, "DELETED", nil
			}
			return nil, "", err
		}
		log.Printf("[DEBUG] ClusterStateRefreshFunc: %#v", clusterGet)
		return clusterGet, clusterGet.Clusterstate, nil
	}
}

func resourceClusterV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := GetRegion(d, config)

	client, err := config.MrsV1Client(region)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud MRS client: %s", err)
	}
	vpcClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud Vpc client: %s", err)
	}

	// Get vpc name
	vpc, err := vpcs.Get(vpcClient, d.Get("vpc_id").(string)).Extract()
	if err != nil {
		return fmt.Errorf("Error retrieving HuaweiCloud Vpc: %s", err)
	}
	// Get subnet name
	subnet, err := subnets.Get(vpcClient, d.Get("subnet_id").(string)).Extract()
	if err != nil {
		return fmt.Errorf("Error retrieving HuaweiCloud Subnet: %s", err)
	}

	loginMode := 0
	if _, ok := d.GetOk("node_public_cert_name"); ok {
		loginMode = 1
	}

	createOpts := &cluster.CreateOpts{
		DataCenter:          region,
		BillingType:         d.Get("billing_type").(int),
		MasterNodeNum:       d.Get("master_node_num").(int),
		MasterNodeSize:      d.Get("master_node_size").(string),
		CoreNodeNum:         d.Get("core_node_num").(int),
		CoreNodeSize:        d.Get("core_node_size").(string),
		AvailableZoneID:     d.Get("available_zone_id").(string),
		ClusterName:         d.Get("cluster_name").(string),
		ClusterVersion:      d.Get("cluster_version").(string),
		ClusterType:         d.Get("cluster_type").(int),
		VpcID:               d.Get("vpc_id").(string),
		SubnetID:            d.Get("subnet_id").(string),
		Vpc:                 vpc.Name,
		SubnetName:          subnet.Name,
		VolumeType:          d.Get("volume_type").(string),
		VolumeSize:          d.Get("volume_size").(int),
		SafeMode:            d.Get("safe_mode").(int),
		LoginMode:           loginMode,
		NodePublicCertName:  d.Get("node_public_cert_name").(string),
		ClusterMasterSecret: d.Get("node_password").(string),
		ClusterAdminSecret:  d.Get("cluster_admin_secret").(string),
		LogCollection:       d.Get("log_collection").(int),
		ComponentList:       getAllClusterComponents(d),
		AddJobs:             getAllClusterJobs(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	clusterCreate, err := cluster.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Cluster: %s", err)
	}

	d.SetId(clusterCreate.ClusterID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"starting"},
		Target:     []string{"running"},
		Refresh:    ClusterStateRefreshFunc(client, clusterCreate.ClusterID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for cluster (%s) to become ready: %s ",
			clusterCreate.ClusterID, err)
	}

	// create tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := expandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "clusters", d.Id(), taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of MRS cluster %s: %s", d.Id(), tagErr)
		}
	}

	return resourceClusterV1Read(d, meta)
}

func resourceClusterV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.MrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud MRS client: %s", err)
	}

	clusterGet, err := cluster.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "MRS Cluster")
	}
	if clusterGet.Clusterstate == "terminated" {
		log.Printf("[WARN] The Cluster %s has been terminated.", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Retrieved Cluster %s: %#v", d.Id(), clusterGet)
	d.SetId(clusterGet.Clusterid)
	d.Set("region", GetRegion(d, config))
	d.Set("order_id", clusterGet.Orderid)
	d.Set("available_zone_name", clusterGet.Azname)
	d.Set("available_zone_id", clusterGet.Azid)
	d.Set("cluster_name", clusterGet.Clustername)
	d.Set("cluster_version", clusterGet.Clusterversion)

	if clusterGet.Masternodenum != "" {
		masterNodeNum, err := strconv.Atoi(clusterGet.Masternodenum)
		if err != nil {
			return fmt.Errorf("Error converting Masternodenum: %s", err)
		}
		d.Set("master_node_num", masterNodeNum)
	}

	if clusterGet.Corenodenum != "" {
		coreNodeNum, err := strconv.Atoi(clusterGet.Corenodenum)
		if err != nil {
			return fmt.Errorf("Error converting Corenodenum: %s", err)
		}
		d.Set("core_node_num", coreNodeNum)
	}

	// the following attributes are empty during to the API backend
	// d.Set("master_node_size", clusterGet.Masternodesize)
	// d.Set("core_node_size", clusterGet.Corenodesize)
	// d.Set("volume_size", clusterGet.Volumesize)

	d.Set("node_public_cert_name", clusterGet.Nodepubliccertname)
	d.Set("safe_mode", clusterGet.Safemode)
	d.Set("instance_id", clusterGet.Instanceid)
	d.Set("hadoop_version", clusterGet.Hadoopversion)
	d.Set("master_node_ip", clusterGet.Masternodeip)
	d.Set("external_ip", clusterGet.Externalip)
	d.Set("private_ip_first", clusterGet.Privateipfirst)
	d.Set("internal_ip", clusterGet.Internalip)
	d.Set("slave_security_groups_id", clusterGet.Slavesecuritygroupsid)
	d.Set("security_groups_id", clusterGet.Securitygroupsid)
	d.Set("external_alternate_ip", clusterGet.Externalalternateip)
	d.Set("master_node_spec_id", clusterGet.Masternodespecid)
	d.Set("core_node_spec_id", clusterGet.Corenodespecid)
	d.Set("master_node_product_id", clusterGet.Masternodeproductid)
	d.Set("core_node_product_id", clusterGet.Corenodeproductid)
	d.Set("duration", clusterGet.Duration)
	d.Set("vnc", clusterGet.Vnc)
	d.Set("fee", clusterGet.Fee)
	d.Set("deployment_id", clusterGet.Deploymentid)
	d.Set("cluster_state", clusterGet.Clusterstate)
	d.Set("error_info", clusterGet.Errorinfo)
	d.Set("remark", clusterGet.Remark)

	updateAt, err := strconv.ParseInt(clusterGet.Updateat, 10, 64)
	if err != nil {
		return fmt.Errorf("Error converting Updateat: %s", err)
	}
	updateAtTm := time.Unix(updateAt, 0)

	createAt, err := strconv.ParseInt(clusterGet.Createat, 10, 64)
	if err != nil {
		return fmt.Errorf("Error converting Createat: %s", err)
	}
	createAtTm := time.Unix(createAt, 0)

	chargingStartTime, err := strconv.ParseInt(clusterGet.Chargingstarttime, 10, 64)
	if err != nil {
		return fmt.Errorf("Error converting chargingStartTime: %s", err)
	}
	chargingStartTimeTm := time.Unix(chargingStartTime, 0)

	d.Set("update_at", updateAtTm.Format(time.RFC3339))
	d.Set("create_at", createAtTm.Format(time.RFC3339))
	d.Set("charging_start_time", chargingStartTimeTm.Format(time.RFC3339))
	d.Set("component_list", clusterGet.Duration)

	components := make([]map[string]interface{}, len(clusterGet.Componentlist))
	for i, attachment := range clusterGet.Componentlist {
		components[i] = make(map[string]interface{})
		components[i]["component_id"] = attachment.Componentid
		components[i]["component_name"] = attachment.Componentname
		components[i]["component_version"] = attachment.Componentversion
		components[i]["component_desc"] = attachment.Componentdesc
		log.Printf("[DEBUG] components: %v", components)
	}

	d.Set("component_list", components)

	// set tags
	if resourceTags, err := tags.Get(client, "clusters", d.Id()).Extract(); err == nil {
		tagmap := tagsToMap(resourceTags.Tags)
		d.Set("tags", tagmap)
	} else {
		log.Printf("[WARN] fetching tags of MRS cluster failed: %s", err)
	}

	return nil
}

func resourceClusterV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.MrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud MRS client: %s", err)
	}

	// update tags
	tagErr := UpdateResourceTags(client, d, "clusters", d.Id())
	if tagErr != nil {
		return fmt.Errorf("Error updating tags of MRS cluster:%s, err:%s", d.Id(), tagErr)
	}

	return resourceClusterV1Read(d, meta)
}

func resourceClusterV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.MrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud MRS client: %s", err)
	}

	rId := d.Id()
	clusterGet, err := cluster.Get(client, d.Id()).Extract()
	if err != nil {
		if isResourceNotFound(err) {
			log.Printf("[INFO] getting an unavailable Cluster: %s", rId)
			return nil
		}
		return fmt.Errorf("Error getting Cluster %s: %s", rId, err)
	}

	if clusterGet.Clusterstate == "terminated" {
		log.Printf("[DEBUG] The Cluster %s has been terminated.", rId)
		return nil
	}

	log.Printf("[DEBUG] Deleting Cluster %s", rId)

	err = cluster.Delete(client, rId).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud Cluster: %s", err)
	}

	log.Printf("[DEBUG] Waiting for Cluster (%s) to be terminated", rId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"running", "terminating"},
		Target:     []string{"terminated"},
		Refresh:    ClusterStateRefreshFunc(client, rId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Cluster (%s) to be terminated: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}
