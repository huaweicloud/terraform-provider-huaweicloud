package bms

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/bms/v1/baremetalservers"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API BMS GET /v1/{project_id}/baremetalservers/{server_id}
// @API BMS PUT /v1/{project_id}/baremetalservers/{server_id}
// @API BMS POST /v1/{project_id}/baremetalservers
// @API BMS POST /v1/{project_id}/baremetalservers/{server_id}/nics
// @API BMS POST /v1/{project_id}/baremetalservers/{server_id}/nics/delete
// @API BMS POST /v1/{project_id}/baremetalservers/{server_id}/tags/action
// @API BMS POST /v1/{project_id}/baremetalservers/{server_id}/metadata
// @API BMS GET /v1/{project_id}/baremetalservers/{server_id}/tags
// @API BMS GET /v1/{project_id}/jobs/{job_id}
// @API VPC GET /v2.0/ports/{port_id}
// @API BSS GET /V2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceBmsInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBmsInstanceCreate,
		ReadContext:   resourceBmsInstanceRead,
		UpdateContext: resourceBmsInstanceUpdate,
		DeleteContext: resourceBmsInstanceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nics": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				// just stash the hash for state & diff comparisons
				StateFunc:        utils.HashAndHexEncode,
				DiffSuppressFunc: utils.SuppressUserData,
			},
			"admin_pass": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"iptype", "eip_charge_mode", "bandwidth_charge_mode", "bandwidth_size", "sharetype",
				},
			},
			"iptype": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"eip_charge_mode", "sharetype", "bandwidth_size",
				},
			},
			"eip_charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid", "postPaid",
				}, true),
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"iptype", "sharetype", "bandwidth_size",
				},
			},
			"sharetype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PER", "WHOLE",
				}, true),
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"iptype", "eip_charge_mode", "bandwidth_size",
				},
			},
			"bandwidth_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_id"},
				RequiredWith: []string{
					"iptype", "eip_charge_mode", "sharetype",
				},
			},
			"bandwidth_charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"traffic", "bandwidth",
				}, true),
				ConflictsWith: []string{"eip_id"},
			},
			"system_disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"system_disk_size",
				},
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"system_disk_type",
				},
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 59,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"charging_mode": common.SchemaChargingMode([]string{}),
			"period_unit":   common.SchemaPeriodUnit([]string{}),
			"period":        common.SchemaPeriod([]string{}),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),

			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// To avoid triggering changes metadata is not backfilled during read.
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceBmsInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	bmsClient, err := cfg.BmsV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating bms client: %s", err)
	}

	createOpts := &baremetalservers.CreateOpts{
		Name:      d.Get("name").(string),
		ImageRef:  d.Get("image_id").(string),
		FlavorRef: d.Get("flavor_id").(string),
		MetaData: baremetalservers.MetaData{
			OpSvcUserId: d.Get("user_id").(string),
			AgencyName:  d.Get("agency_name").(string),
		},
		UserData:         []byte(d.Get("user_data").(string)),
		AdminPass:        d.Get("admin_pass").(string),
		KeyName:          d.Get("key_pair").(string),
		VpcId:            d.Get("vpc_id").(string),
		SecurityGroups:   resourceBmsInstanceSecGroupsV1(d),
		AvailabilityZone: d.Get("availability_zone").(string),
		Nics:             resourceBmsInstanceNicsV1(d),
		DataVolumes:      resourceBmsInstanceDataVolumesV1(d),
		ExtendParam: baremetalservers.ServerExtendParam{
			ChargingMode:        d.Get("charging_mode").(string),
			PeriodType:          d.Get("period_unit").(string),
			PeriodNum:           d.Get("period").(int),
			IsAutoPay:           "true",
			IsAutoRenew:         d.Get("auto_renew").(string),
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		},
	}

	var eipOpts baremetalservers.PublicIp
	var hasEIP bool
	if eipID, ok := d.GetOk("eip_id"); ok {
		hasEIP = true
		eipOpts.Id = eipID.(string)
	} else if eipType, ok := d.GetOk("iptype"); ok {
		hasEIP = true
		eipOpts.Eip = &baremetalservers.Eip{
			IpType: eipType.(string),
			BandWidth: baremetalservers.BandWidth{
				ShareType:  d.Get("sharetype").(string),
				Size:       d.Get("bandwidth_size").(int),
				ChargeMode: d.Get("bandwidth_charge_mode").(string),
			},
			ExtendParam: baremetalservers.EipExtendParam{
				ChargingMode: d.Get("eip_charge_mode").(string),
			},
		}
	}
	if hasEIP {
		createOpts.PublicIp = &eipOpts
	}

	if v, ok := d.GetOk("system_disk_type"); ok {
		volRequest := baremetalservers.RootVolume{
			VolumeType: v.(string),
			Size:       d.Get("system_disk_size").(int),
		}
		createOpts.RootVolume = &volRequest
	}

	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		createOpts.ServerTags = taglist
	}

	n, err := baremetalservers.CreatePrePaid(bmsClient, createOpts).ExtractOrderResponse()
	if err != nil {
		return diag.Errorf("error creating BMS server: %s", err)
	}

	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}
	err = common.WaitOrderComplete(ctx, bssClient, n.OrderID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, n.OrderID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceId)

	// update the user-defined metadata if necessary
	if v, ok := d.GetOk("metadata"); ok {
		metadataOpts := v.(map[string]interface{})
		log.Printf("[DEBUG] BMS metadata options: %v", metadataOpts)

		_, err := baremetalservers.UpdateMetadata(bmsClient, d.Id(), metadataOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating the BMS metadata: %s", err)
		}
	}
	return resourceBmsInstanceRead(ctx, d, meta)
}

func resourceBmsInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	bmsClient, err := cfg.BmsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	server, err := baremetalservers.Get(bmsClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "server")
	}
	if server.Status == "DELETED" {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Retrieved Server %s: %+v", d.Id(), server)

	nics := flattenBmsInstanceNicsV1(d, meta, server.Addresses)

	// Set security groups
	var secGrpIds []string
	for _, sg := range server.SecurityGroups {
		secGrpIds = append(secGrpIds, sg.ID)
	}

	// Set disk ids
	var diskIds []string
	for _, disk := range server.VolumeAttached {
		diskIds = append(diskIds, disk.ID)
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", server.Name),
		d.Set("image_id", server.Image.ID),
		d.Set("flavor_id", server.Flavor.ID),
		d.Set("host_id", server.HostID),
		d.Set("nics", nics),
		d.Set("key_pair", server.KeyName),
		d.Set("security_groups", secGrpIds),
		d.Set("status", server.Status),
		d.Set("user_id", server.Metadata.OpSvcUserId),
		d.Set("image_name", server.Metadata.ImageName),
		d.Set("vpc_id", server.Metadata.VpcID),
		d.Set("agency_name", server.Metadata.AgencyName),
		d.Set("availability_zone", server.AvailabilityZone),
		d.Set("description", server.Description),
		d.Set("user_data", server.UserData),
		d.Set("enterprise_project_id", server.EnterpriseProjectID),
		d.Set("disk_ids", diskIds),
		utils.SetResourceTagsToState(d, bmsClient, "baremetalservers", d.Id()),
		d.Set("tags", d.Get("tags")),
	)

	// Set fixed and floating ip
	if eip := bmsPublicIP(server); eip != "" {
		mErr = multierror.Append(mErr, d.Set("public_ip", eip))
	}
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBmsInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	bmsClient, err := cfg.BmsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	instanceId := d.Id()

	if d.HasChange("name") {
		var updateOpts baremetalservers.UpdateOpts
		updateOpts.Name = d.Get("name").(string)

		_, err = baremetalservers.Update(bmsClient, instanceId, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating bms server: %s", err)
		}
	}

	if d.HasChange("agency_name") {
		metadataOpts := map[string]interface{}{
			"agency_name": d.Get("agency_name").(string),
		}
		_, err := baremetalservers.UpdateMetadata(bmsClient, instanceId, metadataOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating the BMS metadata agency_name: %s", err)
		}
	}

	if d.HasChanges("metadata") {
		_, err := baremetalservers.UpdateMetadata(bmsClient, instanceId, d.Get("metadata").(map[string]interface{})).Extract()
		if err != nil {
			return diag.Errorf("error updating the BMS metadata: %s", err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), instanceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the instance (%s): %s", instanceId, err)
		}
	}

	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(bmsClient, d, "baremetalservers", instanceId)
		if err != nil {
			return diag.Errorf("error updating tags of bms server: %s", err)
		}
	}

	// Security group parammeters are missing in the network card, this is a legacy feature.
	// The first network card can not be deleted, api error message:{"error": {"message": "primary port can not be deleted.", "code":"BMS.0222"}}
	if d.HasChange("nics") {
		err = updateInstanceNics(ctx, d, bmsClient, schema.TimeoutUpdate)
		if err != nil {
			return diag.Errorf("error updating NICs: %s", err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "bms_server",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceBmsInstanceRead(ctx, d, meta)
}

func updateInstanceNics(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string) error {
	addNics, deleteNics := getDiffNics(d)
	if len(deleteNics) > 0 {
		err := deleteInstanceNics(ctx, d, client, timeout, deleteNics)
		if err != nil {
			return err
		}
	}
	if len(addNics) > 0 {
		err := addInstanceNics(ctx, d, client, timeout, addNics)
		if err != nil {
			return err
		}
	}
	return nil
}

func addInstanceNics(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string,
	addNics []interface{}) error {
	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/nics"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{server_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildAddInstanceNicsBodyParams(d, addNics))

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error adding nics to BMS instance: %s", err)
	}
	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, "").(string)
	if jobId == "" {
		return errors.New("error adding nics to BMS instance: job_id is not found in API response")
	}

	return waitForJobComplete(ctx, client, jobId, d.Timeout(timeout))
}

func buildAddInstanceNicsBodyParams(d *schema.ResourceData, addNics []interface{}) map[string]interface{} {
	bodyParams := make([]map[string]interface{}, 0, len(addNics))
	securityGroups := buildAddInstanceNicSecGroupsBodyParams(d)
	for _, nic := range addNics {
		if v, ok := nic.(map[string]interface{}); ok {
			bodyParams = append(bodyParams, map[string]interface{}{
				"subnet_id":       v["subnet_id"].(string),
				"ip_address":      utils.ValueIgnoreEmpty(v["ip_address"]),
				"security_groups": utils.ValueIgnoreEmpty(securityGroups),
			})
		}
	}
	return map[string]interface{}{
		"nics": bodyParams,
	}
}

func buildAddInstanceNicSecGroupsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	segGroups := d.Get("security_groups").(*schema.Set)
	rst := make([]map[string]interface{}, 0, segGroups.Len())
	for _, segGroup := range segGroups.List() {
		rst = append(rst, map[string]interface{}{
			"id": segGroup,
		})
	}
	return rst
}

func deleteInstanceNics(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, timeout string,
	deleteNics []interface{}) error {
	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/nics/delete"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{server_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildDeleteInstanceNicsBodyParams(deleteNics)

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error deleting nics to BMS instance: %s", err)
	}
	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, "").(string)
	if jobId == "" {
		return errors.New("error deleting nics from BMS instance: job_id is not found in API response")
	}

	return waitForJobComplete(ctx, client, jobId, d.Timeout(timeout))
}

func buildDeleteInstanceNicsBodyParams(addNics []interface{}) interface{} {
	bodyParams := make([]interface{}, 0, len(addNics))
	for _, nic := range addNics {
		if v, ok := nic.(map[string]interface{}); ok {
			bodyParams = append(bodyParams, map[string]interface{}{
				"id": v["port_id"].(string),
			})
		}
	}
	return map[string]interface{}{
		"nics": bodyParams,
	}
}

// Get the list of new and to-be-deleted network cards.
func getDiffNics(d *schema.ResourceData) (addList []interface{}, removeList []interface{}) {
	oldNics, newNics := d.GetChange("nics")
	oldList := oldNics.([]interface{})
	newList := newNics.([]interface{})
	for _, ov := range oldList {
		om := ov.(map[string]interface{})
		oSubnetId := om["subnet_id"].(string)
		oIpAddress := om["ip_address"].(string)
		needRemove := true
		for _, nv := range newList {
			nm := nv.(map[string]interface{})
			nSubnetId := nm["subnet_id"].(string)
			nIpAddress := nm["ip_address"].(string)

			if oSubnetId == nSubnetId && oIpAddress == nIpAddress {
				needRemove = false
				break
			}
		}
		if needRemove {
			removeList = append(removeList, ov)
		}
	}

	for _, nv := range newList {
		nm := nv.(map[string]interface{})
		nSubnetId := nm["subnet_id"].(string)
		nIpAddress := nm["ip_address"].(string)
		needAdd := true
		for _, ov := range oldList {
			om := ov.(map[string]interface{})
			oSubnetId := om["subnet_id"].(string)
			oIpAddress := om["ip_address"].(string)
			if nSubnetId == oSubnetId && nIpAddress == oIpAddress {
				needAdd = false
				break
			}
		}
		if needAdd {
			addList = append(addList, nv)
		}
	}
	return
}

func resourceBmsInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	bmsClient, err := cfg.BmsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}
	serverID := d.Id()
	publicIP := d.Get("public_ip").(string)
	diskIds := d.Get("disk_ids").([]interface{})

	resourceIDs := make([]string, 0, 2+len(diskIds))

	resourceIDs = append(resourceIDs, serverID)

	if len(diskIds) > 0 {
		for _, diskId := range diskIds {
			resourceIDs = append(resourceIDs, diskId.(string))
		}
	}

	// unsubscribe the eip if necessary
	if _, ok := d.GetOk("iptype"); ok && publicIP != "" && d.Get("eip_charge_mode").(string) == "prePaid" {
		eipClient, err := cfg.NetworkingV1Client(region)
		if err != nil {
			return diag.Errorf("error creating networking client: %s", err)
		}

		epsID := "all_granted_eps"
		var eipID string
		if eipID, err = common.GetEipIDbyAddress(eipClient, publicIP, epsID); err != nil {
			return diag.Errorf("error fetching EIP ID of BMS server (%s): %s", d.Id(), err)
		}
		resourceIDs = append(resourceIDs, eipID)
	}

	if err := common.UnsubscribePrePaidResource(d, cfg, resourceIDs); err != nil {
		return diag.Errorf("error unsubscribing BMS server: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting", "ACTIVE", "SHUTOFF"},
		Target:       []string{"DELETED"},
		Refresh:      waitForBmsInstanceDelete(bmsClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting BMS instance: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceBmsInstanceNicsV1(d *schema.ResourceData) []baremetalservers.Nic {
	var nicRequests []baremetalservers.Nic

	nics := d.Get("nics").([]interface{})
	for i := range nics {
		nic := nics[i].(map[string]interface{})
		nicRequest := baremetalservers.Nic{
			SubnetId:  nic["subnet_id"].(string),
			IpAddress: nic["ip_address"].(string),
		}

		nicRequests = append(nicRequests, nicRequest)
	}
	return nicRequests
}

func resourceBmsInstanceDataVolumesV1(d *schema.ResourceData) []baremetalservers.DataVolume {
	var volRequests []baremetalservers.DataVolume

	vols := d.Get("data_disks").([]interface{})
	for i := range vols {
		vol := vols[i].(map[string]interface{})
		volRequest := baremetalservers.DataVolume{
			VolumeType: vol["type"].(string),
			Size:       vol["size"].(int),
		}
		volRequests = append(volRequests, volRequest)
	}
	return volRequests
}

func resourceBmsInstanceSecGroupsV1(d *schema.ResourceData) []baremetalservers.SecurityGroup {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secgroups := make([]baremetalservers.SecurityGroup, len(rawSecGroups))
	for i, raw := range rawSecGroups {
		secgroups[i] = baremetalservers.SecurityGroup{
			ID: raw.(string),
		}
	}
	return secgroups
}

func flattenBmsInstanceNicsV1(d *schema.ResourceData, meta interface{},
	addresses map[string][]baremetalservers.Address) []map[string]interface{} {
	cfg := meta.(*config.Config)
	networkingClient, err := cfg.NetworkingV2Client(cfg.GetRegion(d))
	if err != nil {
		log.Printf("Error creating networking client: %s", err)
	}

	var network string
	var nics []map[string]interface{}
	// Loop through all networks and addresses.
	for _, addrs := range addresses {
		for _, addr := range addrs {
			// Skip if not fixed ip
			if addr.Type != "fixed" {
				continue
			}

			p, err := ports.Get(networkingClient, addr.PortID).Extract()
			if err != nil {
				network = ""
				log.Printf("[DEBUG] flattenInstanceNicsV1: failed to fetch port %s", addr.PortID)
			} else {
				network = p.NetworkID
			}

			v := map[string]interface{}{
				"subnet_id":   network,
				"ip_address":  addr.Addr,
				"mac_address": addr.MacAddr,
				"port_id":     addr.PortID,
			}
			nics = append(nics, v)
		}
	}

	log.Printf("[DEBUG] flattenInstanceNicsV1: %#v", nics)
	return nics
}

func bmsPublicIP(server *baremetalservers.CloudServer) string {
	var publicIP string

	for _, addresses := range server.Addresses {
		for _, addr := range addresses {
			if addr.Type == "floating" {
				publicIP = addr.Addr
				break
			}
		}
	}

	return publicIP
}

func waitForBmsInstanceDelete(bmsClient *golangsdk.ServiceClient, serverId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete BMS instance %s", serverId)

		r, err := baremetalservers.Get(bmsClient, serverId).Extract()

		if err != nil {
			return r, "Deleting", err
		}

		return r, r.Status, nil
	}
}
