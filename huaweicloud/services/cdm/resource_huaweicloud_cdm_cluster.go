package cdm

import (
	"context"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cdm/v1/clusters"
	"github.com/chnsz/golangsdk/openstack/common/tags"
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

func ResourceCdmCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCdmClusterCreate,
		ReadContext:   resourceCdmClusterRead,
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

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "2.8.6.2",
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

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"phone_num": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 5,
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
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceCdmClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1 client, err=%s", err)
	}

	opts := clusters.ClusterCreateOpts{}
	buildClusterParamter(d, &opts, config.GetEnterpriseProjectID(d))
	buildNotifyParamter(d, &opts)

	logp.Printf("[DEBUG] Creating CDM cluster opts: %#v", opts)

	rst, createErr := clusters.Create(client, opts)
	if createErr != nil {
		return fmtp.DiagErrorf("Error creating CDM cluster: %s", createErr)
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
		Datastore: clusters.Datastore{Type: "cdm", Version: d.Get("version").(string)},
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
	opts.Email = strings.Join(utils.ExpandToStringList(d.Get("email").([]interface{})), ",")
	opts.PhoneNum = strings.Join(utils.ExpandToStringList(d.Get("phone_num").([]interface{})), ",")

	if opts.Email != "" || opts.PhoneNum != "" {
		opts.AutoRemind = utils.Bool(true)
	}
}

func resourceCdmClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1 client, err=%s", err)
	}

	detail, err := clusters.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving CDM cluster")
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
		setInstancesToState(d, detail.Instances),
		d.Set("schedule_boot_time", detail.ScheduleBootTime),
		d.Set("schedule_off_time", detail.ScheduleOffTime),
		d.Set("created", detail.Created),
		d.Set("public_ip", detail.Instances[0].PublicIp),
		d.Set("public_endpoint", detail.PublicEndpoint),
		d.Set("status", detail.StatusDetail),
	)

	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("Error setting CDM fields: %s", mErr)
	}

	return nil
}

func setInstancesToState(d *schema.ResourceData, items []clusters.Instance) error {
	if len(items) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(items))
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

	return d.Set("instances", result)
}

func resourceCdmClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1 client, err=%s", err)
	}

	dErr := clusters.Delete(client, d.Id(), clusters.ClusterDeleteOpts{})
	if dErr.Err != nil {
		return fmtp.DiagErrorf("delete CDM cluster failed. %q:%s", d.Id(), dErr)
	}

	err = waitingforClusterDeleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func waitingforClusterCreated(ctx context.Context, client *golangsdk.ServiceClient, id string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{clusters.StatusCreating, clusters.StatusStarting},
		Target:  []string{clusters.StatusNormal},
		Refresh: func() (interface{}, string, error) {
			cluster, err := clusters.Get(client, id)
			logp.Printf("[DEBUG] query CDM cluster in create check func: %#v,%s", cluster, err)
			if err != nil {
				return nil, "", err
			}

			if cluster.Status == clusters.StatusCreationFailed || cluster.Status == clusters.StatusFailed {
				return cluster, "failed", fmtp.Errorf("%s:%s", cluster.Status, cluster.StatusDetail)
			}
			return cluster, cluster.Status, nil
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("error waiting for CDM cluster (%s) to be created: %s", id, err)
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
			logp.Printf("[DEBUG] query CDM cluster in delete check func: %#v,%s", cluster, err)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return true, "Done", nil
				}
				return nil, "", nil
			}
			return true, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("error waiting for CDM cluster (%s) to be deleted: %s", id, err)
	}
	return nil
}
