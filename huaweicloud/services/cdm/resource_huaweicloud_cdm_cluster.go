package cdm

import (
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
	"github.com/chnsz/golangsdk/openstack/cdm/v1/clusters"
	"github.com/chnsz/golangsdk/openstack/common/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDM DELETE /v1.1/{project_id}/clusters/{clusterId}
// @API CDM GET /v1.1/{project_id}/clusters/{clusterId}
// @API CDM POST /v1.1/{project_id}/clusters
// @API CDM POST /v1.1/{project_id}/cluster/modify/{cluster_id}
func ResourceCdmCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCdmClusterCreate,
		ReadContext:   resourceCdmClusterRead,
		UpdateContext: resourceCdmClusterUpdate,
		DeleteContext: resourceCdmClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

			"flavor_id": {
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

			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_auto_off": {
				Type:          schema.TypeBool,
				Computed:      true,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"schedule_boot_time", "schedule_off_time"},
			},

			"schedule_boot_time": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"is_auto_off"},
			},

			"schedule_off_time": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"is_auto_off"},
			},

			"email": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"phone_num": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instances": {
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

						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"manage_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"traffic_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"flavor_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceCdmClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1 client, err=%s", err)
	}

	opts := clusters.ClusterCreateOpts{}
	buildClusterParamter(d, &opts, cfg.GetEnterpriseProjectID(d))
	buildNotifyParamter(d, &opts)

	rst, createErr := clusters.Create(client, opts)
	if createErr != nil {
		return diag.Errorf("error creating CDM cluster: %s", createErr)
	}

	d.SetId(rst.Id)

	checkCreateErr := waitingforClusterCreated(ctx, client, rst.Id, d.Timeout(schema.TimeoutCreate))
	if checkCreateErr != nil {
		return diag.FromErr(checkCreateErr)
	}
	return resourceCdmClusterRead(ctx, d, meta)
}

func buildClusterParamter(d *schema.ResourceData, opts *clusters.ClusterCreateOpts, enterpriseProjectID string) {
	cluster := clusters.ClusterRequest{
		Name:      d.Get("name").(string),
		VpcId:     d.Get("vpc_id").(string),
		IsAutoOff: utils.Bool(d.Get("is_auto_off").(bool)),
		Instances: []clusters.InstanceReq{
			{
				AvailabilityZone: d.Get("availability_zone").(string),
				FlavorRef:        d.Get("flavor_id").(string),
				Type:             "cdm",
				Nics: []clusters.Nics{
					{
						SecurityGroupId: d.Get("security_group_id").(string),
						NetId:           d.Get("subnet_id").(string),
					},
				},
			},
		},
	}

	if v, ok := d.GetOk("version"); ok {
		cluster.Datastore = &clusters.Datastore{
			Type:    "cdm",
			Version: v.(string),
		}
	}

	// set Schedule boot/off
	bootTime := d.Get("schedule_boot_time").(string)
	offTime := d.Get("schedule_off_time").(string)

	if bootTime != "" || offTime != "" {
		cluster.IsScheduleBootOff = utils.Bool(true)
		cluster.ScheduleBootTime = bootTime
		cluster.ScheduleOffTime = offTime
	}

	if enterpriseProjectID != "" {
		cluster.SysTags = []tags.ResourceTag{
			{
				Key:   "_sys_enterprise_project_id",
				Value: enterpriseProjectID,
			},
		}
	}
	opts.Cluster = cluster
}

func buildNotifyParamter(d *schema.ResourceData, opts *clusters.ClusterCreateOpts) {
	opts.Email = strings.Join(utils.ExpandToStringList(d.Get("email").(*schema.Set).List()), ",")
	opts.PhoneNum = strings.Join(utils.ExpandToStringList(d.Get("phone_num").(*schema.Set).List()), ",")

	if opts.Email != "" || opts.PhoneNum != "" {
		opts.AutoRemind = utils.Bool(true)
	}
}

func resourceCdmClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1 client, err=%s", err)
	}

	detail, err := clusters.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CDM cluster")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", detail.Name),
		d.Set("availability_zone", detail.AzName),
		d.Set("version", detail.Datastore.Version),
		d.Set("flavor_id", detail.Instances[0].Flavor.Id),
		d.Set("vpc_id", detail.VpcId),
		d.Set("subnet_id", detail.SubnetId),
		d.Set("security_group_id", detail.SecurityGroupId),
		d.Set("is_auto_off", detail.IsAutoOff),
		d.Set("name", detail.Name),
		d.Set("instances", flattenInstancs(detail.Instances)),
		d.Set("schedule_boot_time", detail.ScheduleBootTime),
		d.Set("schedule_off_time", detail.ScheduleOffTime),
		d.Set("created", detail.Created),
		d.Set("public_ip", detail.Instances[0].PublicIp),
		d.Set("public_endpoint", detail.PublicEndpoint),
		d.Set("status", detail.StatusDetail),
		d.Set("enterprise_project_id", detail.EnterpriseProjectId),
		d.Set("flavor_name", detail.FlavorName),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstancs(items []clusters.Instance) []map[string]interface{} {
	if len(items) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, instance := range items {
		item := map[string]interface{}{
			"id":         instance.Id,
			"name":       instance.Name,
			"private_ip": instance.PrivateIp,
			"public_ip":  instance.PublicIp,
			"manage_ip":  instance.ManageIp,
			"role":       instance.Role,
			"traffic_ip": instance.TrafficIp,
			"type":       instance.Type,
		}
		result = append(result, item)
	}

	return result
}

func resourceCdmClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	clusterId := d.Id()
	client, err := cfg.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1 client, err=%s", err)
	}

	if d.HasChanges("email", "phone_num") {
		err := updateEmailAndPhoneNum(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   clusterId,
			ResourceType: "cdm-clusters",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCdmClusterRead(ctx, d, meta)
}

func updateEmailAndPhoneNum(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateHttpUrl := "v1.1/{project_id}/cluster/modify/{cluster_id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateClusterBodyParam(d)),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "zh-cn",
		},
	}

	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating cluster notification infos: %s", err)
	}

	return nil
}

func buildUpdateClusterBodyParam(d *schema.ResourceData) map[string]interface{} {
	email := strings.Join(utils.ExpandToStringList(d.Get("email").(*schema.Set).List()), ",")
	phoneNum := strings.Join(utils.ExpandToStringList(d.Get("phone_num").(*schema.Set).List()), ",")

	bodyParams := map[string]interface{}{
		"email":    utils.ValueIgnoreEmpty(email),
		"phoneNum": utils.ValueIgnoreEmpty(phoneNum),
	}
	if email == "" && phoneNum == "" {
		bodyParams["autoRemind"] = utils.Bool(false)
	} else {
		bodyParams["autoRemind"] = utils.Bool(true)
	}

	// Shutdown cluster is no longer supported, to be deprecated,
	// so resource does not support to update those params, input them to avoid being covered.
	bootTime := d.Get("schedule_boot_time").(string)
	offTime := d.Get("schedule_off_time").(string)

	bodyParams["autoOff"] = d.Get("is_auto_off")

	if bootTime != "" || offTime != "" {
		bodyParams["scheduleBootOff"] = utils.Bool(true)
		bodyParams["scheduleBootTime"] = bootTime
		bodyParams["scheduleOffTime"] = offTime
	}

	return bodyParams
}

func resourceCdmClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1 client, err=%s", err)
	}

	dErr := clusters.Delete(client, d.Id(), clusters.ClusterDeleteOpts{})
	if dErr.Err != nil {
		return diag.Errorf("delete CDM cluster failed. %q:%s", d.Id(), dErr)
	}

	err = waitingforClusterDeleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitingforClusterCreated(ctx context.Context, client *golangsdk.ServiceClient, id string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{clusters.StatusCreating, clusters.StatusStarting},
		Target:  []string{clusters.StatusNormal},
		Refresh: func() (interface{}, string, error) {
			cluster, err := clusters.Get(client, id)
			log.Printf("[DEBUG] query CDM cluster in create check func: %#v,%s", cluster, err)
			if err != nil {
				return nil, "", err
			}

			if cluster.Status == clusters.StatusCreationFailed || cluster.Status == clusters.StatusFailed {
				return cluster, "failed", fmt.Errorf("%s:%s", cluster.Status, cluster.StatusDetail)
			}
			return cluster, cluster.Status, nil
		},
		Timeout:      timeout,
		PollInterval: 20 * time.Second,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CDM cluster (%s) to be created: %s", id, err)
	}
	return nil
}

func waitingforClusterDeleted(ctx context.Context, client *golangsdk.ServiceClient, id string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			cluster, err := clusters.Get(client, id)
			log.Printf("[DEBUG] query CDM cluster in delete check func: %#v,%s", cluster, err)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return true, "Done", nil
				}
				return nil, "", err
			}
			return true, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 20 * time.Second,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CDM cluster (%s) to be deleted: %s", id, err)
	}
	return nil
}
