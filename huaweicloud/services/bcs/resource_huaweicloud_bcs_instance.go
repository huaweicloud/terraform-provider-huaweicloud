package bcs

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bcs/v2/blockchains"
	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"
	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/kafka"
)

// @API BCS POST /v2/{project_id}/blockchains
// @API BCS DELETE /v2/{project_id}/blockchains/{blockchain_id}
// @API BCS GET /v2/{project_id}/blockchains/{blockchain_id}
// @API CCE GET /api/v3/projects/{project_id}/clusters/{cluster_id}
// @API DMS DELETE /v2/{project_id}/instances/{instance_id}
// @API DMS GET /v2/{project_id}/instances
func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBCSInstanceCreate,
		ReadContext:   resourceBCSInstanceRead,
		UpdateContext: resourceBCSInstanceUpdate,
		DeleteContext: resourceBCSInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
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
			"orderer_node_num": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"consensus": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"edition": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fabric_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"org_disk_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"blockchain_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_mechanism": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"database_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"eip_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"cce_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_orgs": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"org_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"pvc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_detail": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
				Set: resourcePeerOrgsHash,
			},
			"channels": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"org_names": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Set: resourceChannelsHash,
			},
			"couchdb": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							ForceNew:  true,
							Sensitive: true,
						},
					},
				},
			},
			"sfs_turbo": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"flavor": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"share_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"block_info": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transaction_quantity": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"block_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"generation_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"kafka": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"flavor": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"storage_size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tc3_need": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"restful_api_support": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"delete_obs": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delete_storage": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"purchase_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cross_region_support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rollback_support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"old_service_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agent_portal_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceBCSInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var newCluster = false
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.BcsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Blockchain client: %s", err)
	}
	createOpts := blockchains.CreateOpts{
		CreateNewCluster:    &newCluster,
		Name:                d.Get("name").(string),
		VersionType:         d.Get("edition").(int),
		FabricVersion:       d.Get("fabric_version").(string),
		BlockChainType:      d.Get("blockchain_type").(string),
		Consensus:           d.Get("consensus").(string),
		EIPEnable:           d.Get("eip_enable").(bool),
		SignAlgorithm:       d.Get("security_mechanism").(string),
		VolumeType:          d.Get("volume_type").(string),
		OrgDiskSize:         d.Get("org_disk_size").(int),
		DatabaseType:        d.Get("database_type").(string),
		Password:            d.Get("password").(string),
		OrdererNodeNumber:   d.Get("orderer_node_num").(int),
		TC3Need:             d.Get("tc3_need").(bool),
		RestfulAPISupport:   d.Get("restful_api_support").(bool),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		PeerOrgs:            resourceOrgPeer(d),
		Channels:            resourceChannelInfo(d),
		CouchDBInfo:         resourceCouchDBInfo(d),
		SFSTurbo:            resourceTurboInfo(d),
		Block:               resourceBlockInfo(d),
		Kafka:               resourceKafkaCreateInfo(d),
	}

	v, err := resourceClusterInfo(cfg, d.Get("cce_cluster_id").(string), region)
	if err != nil {
		return diag.Errorf("get cluster information failed: %s ", err)
	}
	createOpts.ClusterType = "cce"
	createOpts.CCEClusterInfo = v

	res, err := blockchains.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating Blockchain instance: %s ", err)
	}

	d.SetId(res.ID)
	instanceID := d.Id()
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"IsCreating"},
		Target:       []string{"Normal"},
		Refresh:      blockchainStateRefreshFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        150 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for Blockchain instance (%s) status to normal: %s ", instanceID, err)
	}

	return resourceBCSInstanceRead(ctx, d, meta)
}

func resourceBCSInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.BcsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Blockchain client: %s", err)
	}

	instanceID := d.Id()
	instance, err := blockchains.Get(client, instanceID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting Blockchain instance")
	}
	log.Printf("[DEBUG] retrieved Blockchain instance %s: %#v", instanceID, instance)

	mErr := multierror.Append(
		d.Set("name", instance.Basic.Name),
		d.Set("edition", instance.Basic.VersionType),
		d.Set("blockchain_type", instance.Basic.ServiceType),
		d.Set("consensus", instance.Basic.Consensus),
		d.Set("security_mechanism", instance.Basic.SignAlgorithm),
		d.Set("database_type", instance.Basic.DatabaseType),
		d.Set("restful_api_support", instance.Basic.IsSupportRestful),
		d.Set("status", instance.Basic.Status),
		d.Set("tc3_need", instance.CouchDB != (blockchains.CouchDB{})),
		d.Set("cluster_type", instance.Basic.ClusterType),
		d.Set("cce_cluster_id", instance.Basic.ClusterID),
		d.Set("version", instance.Basic.Version),
		d.Set("purchase_type", instance.Basic.PurchaseType),
		d.Set("cross_region_support", instance.Basic.IsCrossRegion),
		d.Set("rollback_support", instance.Basic.IsSupportRollback),
		d.Set("old_service_version", instance.Basic.OldServiceVersion),
		d.Set("agent_portal_address", instance.Basic.AgentPortalAddress),
	)

	channelList := make([]map[string]interface{}, len(instance.Channels))
	for i, v := range instance.Channels {
		channel := map[string]interface{}{
			"name":      v.Name,
			"org_names": v.OrgNames,
		}
		channelList[i] = channel
	}
	mErr = multierror.Append(mErr, d.Set("channels", channelList))

	peerList := make([]map[string]interface{}, len(instance.Peer))
	for i, org := range instance.Peer {
		address := make([]map[string]interface{}, len(org.Address))
		for j, v := range org.Address {
			addr := map[string]interface{}{
				"domain_port": v.DomainPort,
				"ip_port":     v.IPPort,
			}
			address[j] = addr
		}
		peerList[i] = map[string]interface{}{
			"org_name":      org.Name,
			"count":         org.NodeCount,
			"status":        org.Status,
			"status_detail": org.StatusDetail,
			"pvc_name":      org.PVCName,
			"address":       address,
		}
	}
	mErr = multierror.Append(mErr, d.Set("peer_orgs", peerList))

	if instance.CouchDB != (blockchains.CouchDB{}) {
		couchDBList := make([]map[string]interface{}, 1)
		info := map[string]interface{}{
			"user_name": instance.CouchDB.User,
			"password":  d.Get("couchdb.0.password"),
		}
		couchDBList[0] = info
		mErr = multierror.Append(mErr, d.Set("couchdb", couchDBList))
	}

	if instance.Basic.Consensus == "kafka" {
		kafkaList := make([]map[string]interface{}, 1)
		info := map[string]interface{}{
			"name":              instance.DMSKafka.Name,
			"flavor":            d.Get("kafka.0.flavor"),
			"storage_size":      d.Get("kafka.0.storage_size"),
			"availability_zone": d.Get("kafka.0.availability_zone"),
		}
		kafkaList[0] = info
		mErr = multierror.Append(mErr, d.Set("kafka", kafkaList))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBCSInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	instanceId := d.Id()

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "bcs",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceBCSInstanceRead(ctx, d, meta)
}

func resourceBCSInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	bcsClient, err := cfg.BcsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Blockchain client: %s", err)
	}
	if d.Get("consensus").(string) == "kafka" {
		dmsClient, err := cfg.DmsV2Client(region)
		if err != nil {
			return diag.Errorf("error creating kafka client: %s", err)
		}
		kafkaName := d.Get("kafka.0.name").(string)
		listOpts := instances.ListOpts{
			Engine: "kafka",
			Name:   kafkaName,
		}
		pages, err := instances.List(dmsClient, listOpts).AllPages()
		if err != nil {
			return diag.Errorf("error getting kafka instance in queue (%s): %s", kafkaName, err)
		}
		res, err := instances.ExtractInstances(pages)
		if err != nil {
			return diag.Errorf("error quering kafka instances: %s", err)
		}
		if len(res.Instances) < 1 {
			return diag.Errorf("error quering kafka, returned no results")
		}
		if len(res.Instances) > 1 {
			return diag.Errorf("error quering kafka, returned more than one result")
		}
		kafkaID := res.Instances[0].InstanceID
		if r := instances.Delete(dmsClient, kafkaID); r.Result.Err != nil {
			return diag.Errorf("error deleting kafka instance (%s): %s ", kafkaID, r.Result.Err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"DELETING", "RUNNING"},
			Target:     []string{"DELETED"},
			Refresh:    kafka.KafkaInstanceStateRefreshFunc(dmsClient, kafkaID),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return diag.Errorf("error waiting for instance (%s) to delete: %s", kafkaID, err)
		}
	}

	blockchainID := d.Id()
	deleteOpts := blockchains.DeleteOpts{}
	if v, ok := d.GetOk("delete_obs"); ok {
		deleteOpts.IsDeleteOBS = v.(bool)
	}
	if v, ok := d.GetOk("delete_storage"); ok {
		deleteOpts.IsDeleteStorage = v.(bool)
	}

	if err := blockchains.Delete(bcsClient, deleteOpts, blockchainID).Extract(); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting Blockchain instance")
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"IsDeleting"},
		Target:       []string{"IsDeleted"},
		Refresh:      blockchainStateRefreshFunc(bcsClient, blockchainID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        15 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for Blockchain instance (%s) status to deleted: %s ", blockchainID, err)
	}

	return nil
}

func blockchainStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := blockchains.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault400); ok {
				return instance, "IsDeleted", nil
			}
			if _, ok := err.(golangsdk.ErrDefault401); ok {
				return instance, "IsDeleted", nil
			}
			return nil, "FOUND ERROR", err
		}
		if instance.Basic.ID == "" {
			return instance, "IsDeleted", nil
		}
		if instance.Basic.ProcessStatus != "" {
			return instance, instance.Basic.ProcessStatus, nil
		}

		return instance, instance.Basic.Status, nil
	}
}

func resourceClusterInfo(cfg *config.Config, clusterID, region string) (*blockchains.CCEClusterInfo, error) {
	clusterInfo := new(blockchains.CCEClusterInfo)

	client, err := cfg.CceV3Client(region)
	if err != nil {
		return clusterInfo, fmt.Errorf("error creating CCE client: %s", err)
	}
	n, err := clusters.Get(client, clusterID).Extract()
	if err != nil {
		return clusterInfo, fmt.Errorf("error retrieving CCE: %s", err)
	}
	clusterInfo.ID = clusterID
	clusterInfo.Name = n.Metadata.Name

	return clusterInfo, nil
}

func resourceOrgPeer(d *schema.ResourceData) []blockchains.PeerOrg {
	infoRaw := d.Get("peer_orgs")
	nodeList := make([]blockchains.PeerOrg, infoRaw.(*schema.Set).Len())
	for i, v := range infoRaw.(*schema.Set).List() {
		var peerOrgInfo blockchains.PeerOrg
		peerOrgInfo.Name = v.(map[string]interface{})["org_name"].(string)
		peerOrgInfo.NodeCount = v.(map[string]interface{})["count"].(int)
		nodeList[i] = peerOrgInfo
	}
	return nodeList
}

func resourceChannelInfo(d *schema.ResourceData) []blockchains.ChannelInfo {
	chRaw := d.Get("channels")
	channelList := make([]blockchains.ChannelInfo, 0, chRaw.(*schema.Set).Len())

	for _, v := range chRaw.(*schema.Set).List() {
		var channelInfo blockchains.ChannelInfo
		channelInfo.Name = v.(map[string]interface{})["name"].(string)
		orgNameList := make([]string, len(v.(map[string]interface{})["org_names"].([]interface{})))
		for j, org := range v.(map[string]interface{})["org_names"].([]interface{}) {
			orgNameList[j] = org.(string)
		}
		channelInfo.OrgNames = orgNameList
		channelList = append(channelList, channelInfo)
	}
	return channelList
}

func resourceCouchDBInfo(d *schema.ResourceData) *blockchains.CouchDBInfo {
	var couchDBInfo *blockchains.CouchDBInfo
	var infoRaw []interface{}
	if v, ok := d.GetOk("couchdb"); ok {
		infoRaw = v.([]interface{})
	}
	if len(infoRaw) == 1 {
		couchDBInfo = new(blockchains.CouchDBInfo)
		couchDBInfo.UserName = infoRaw[0].(map[string]interface{})["user_name"].(string)
		couchDBInfo.Password = infoRaw[0].(map[string]interface{})["password"].(string)
	}
	return couchDBInfo
}

func resourceTurboInfo(d *schema.ResourceData) *blockchains.SFSTurbo {
	var turboInfo *blockchains.SFSTurbo
	var infoRaw []interface{}
	if v, ok := d.GetOk("sfs_turbo"); ok {
		infoRaw = v.([]interface{})
	}
	if len(infoRaw) == 1 {
		turboInfo = new(blockchains.SFSTurbo)
		turboInfo.ShareType = infoRaw[0].(map[string]interface{})["share_type"].(string)
		turboInfo.Type = infoRaw[0].(map[string]interface{})["type"].(string)
		turboInfo.AvailabilityZone = infoRaw[0].(map[string]interface{})["availability_zone"].(string)
		turboInfo.Flavor = infoRaw[0].(map[string]interface{})["flavor"].(string)
	}
	return turboInfo
}

func resourceBlockInfo(d *schema.ResourceData) *blockchains.BlockInfo {
	var blockInfo *blockchains.BlockInfo
	var infoRaw []interface{}
	if v, ok := d.GetOk("block_info"); ok {
		infoRaw = v.([]interface{})
	}
	if len(infoRaw) == 1 {
		blockInfo = new(blockchains.BlockInfo)
		blockInfo.BatchTimeout = infoRaw[0].(map[string]interface{})["generation_interval"].(int)
		blockInfo.MaxMessageCount = infoRaw[0].(map[string]interface{})["transaction_quantity"].(int)
		blockInfo.PreferredMaxbytes = infoRaw[0].(map[string]interface{})["block_size"].(int)
	}
	return blockInfo
}

func resourceKafkaCreateInfo(d *schema.ResourceData) *blockchains.KafkaInfo {
	var createInfo *blockchains.KafkaInfo
	var buffer bytes.Buffer
	var infoRaw []interface{}

	buffer.WriteString("dms.instance.kafka.cluster.")
	if v, ok := d.GetOk("kafka"); ok {
		infoRaw = v.([]interface{})
	}
	if len(infoRaw) == 1 {
		createInfo = new(blockchains.KafkaInfo)
		createInfo.Storage = infoRaw[0].(map[string]interface{})["storage_size"].(int)
		buffer.WriteString(infoRaw[0].(map[string]interface{})["flavor"].(string))
		createInfo.Flavor = buffer.String()
		sliceAZ := make([]string, len(infoRaw[0].(map[string]interface{})["availability_zone"].([]interface{})))
		for i, v := range infoRaw[0].(map[string]interface{})["availability_zone"].([]interface{}) {
			sliceAZ[i] = v.(string)
		}
		createInfo.AvailabilityZone = strings.Join(sliceAZ, ",")
	}
	return createInfo
}

func resourceChannelsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["name"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourcePeerOrgsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["org_name"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["org_name"].(string)))
	}

	return hashcode.String(buf.String())
}
