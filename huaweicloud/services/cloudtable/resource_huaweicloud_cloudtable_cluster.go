package cloudtable

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cloudtable/v2/clusters"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

type stateRefresh struct {
	Pending        []string
	Target         []string
	Delay          time.Duration
	Timeout        time.Duration
	PollInterval   time.Duration
	NotFoundChecks int
}

var stateCodes = map[string]string{
	"100": "Creating",
	"200": "Running",
	"300": "Abnormal",
	"303": "Creation failed",
	"800": "Frezon",
}

// @API CloudTable DELETE /v2/{project_id}/clusters/{cluster_id}
// @API CloudTable GET /v2/{project_id}/clusters/{cluster_id}
// @API CloudTable POST /v2/{project_id}/clusters
func ResourceCloudTableCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudTableClusterCreate,
		ReadContext:   resourceCloudTableClusterRead,
		DeleteContext: resourceCloudTableClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

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
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"COMMON", "ULTRAHIGH",
				}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hbase_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"opentsdb_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"rs_num": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  2,
			},
			"iam_auth_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hbase_public_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"open_tsdb_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"opentsdb_public_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"storage_size_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zookeeper_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCloudTableClusterParams(d *schema.ResourceData) clusters.CreateOpts {
	result := clusters.CreateOpts{
		Name:           d.Get("name").(string),
		VpcId:          d.Get("vpc_id").(string),
		StorageType:    d.Get("storage_type").(string),
		IAMAuthEnabled: d.Get("iam_auth_enabled").(bool),
		Datastore: clusters.Datastore{
			Type:    "hbase",
			Version: d.Get("hbase_version").(string),
		},
		Instance: clusters.Instance{
			AvailabilityZone: d.Get("availability_zone").(string),
			CUNum:            d.Get("rs_num").(int),
			Networks: []clusters.Network{
				{
					SubnetId:        d.Get("network_id").(string),
					SecurityGroupId: d.Get("security_group_id").(string),
				},
			},
		},
	}
	if num, ok := d.GetOk("opentsdb_num"); ok {
		result.OpenTSDBEnabled = true
		result.Instance.TSDNum = num.(int)
	}

	return result
}

func resourceCloudTableClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.CloudtableV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CloudTable V2 client: %s", err)
	}

	opt := buildCloudTableClusterParams(d)
	resp, err := clusters.Create(client, opt)
	if err != nil {
		return diag.Errorf("Error creating CloudTable cluster: %s", err)
	}
	d.SetId(resp.ClusterId)
	stateRef := stateRefresh{
		Pending:      []string{"Creating"},
		Target:       []string{"Running"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 15 * time.Second,
		// Due to the design of the API, the cluster cannot be queried for a few time after the create request is sent.
		NotFoundChecks: 1,
	}
	if err := waitForCloudTableClusterStateRefresh(ctx, client, d.Id(), stateRef); err != nil {
		return err
	}

	return resourceCloudTableClusterRead(ctx, d, meta)
}

func setCloudTableClusterStatus(d *schema.ResourceData, strCode string) error {
	if val, ok := stateCodes[strCode]; ok {
		return d.Set("status", val)
	}
	return fmt.Errorf("the status is abnormal, the code is %+v", strCode)
}

// Using this function needs to ensure that the strVal is a string type number.
func setCloudTableIntParam(d *schema.ResourceData, paramName, strVal string) error {
	numVal, err := strconv.Atoi(strVal)
	if err == nil {
		// lintignore:R001
		return d.Set(paramName, numVal)
	}
	return nil
}

func resourceCloudTableClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CloudtableV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CloudTable v2 client: %s", err)
	}

	resp, err := clusters.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error getting the specifies CloudTable cluster form server")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("storage_type", resp.StorageType),
		d.Set("vpc_id", resp.VpcId),
		d.Set("security_group_id", resp.SecurityGroupId),
		d.Set("iam_auth_enabled", resp.IAMAuthEnabled),
		// Computed parameters
		d.Set("created_at", resp.CreateAt),
		d.Set("hbase_public_endpoint", resp.HbasePublicEndpoint),
		d.Set("open_tsdb_link", resp.OpenTSDBLink),
		d.Set("opentsdb_public_endpoint", resp.TSDPublicEndpoint),
		d.Set("zookeeper_link", resp.ZookeeperLink),
		d.Set("hbase_version", resp.Datastore.Version),
		setCloudTableIntParam(d, "opentsdb_num", resp.TSDNum),
		setCloudTableIntParam(d, "rs_num", resp.CUNum),
		setCloudTableIntParam(d, "storage_size", resp.StorageQuota),
		setCloudTableIntParam(d, "storage_size_used", resp.StorageUsed),
		setCloudTableClusterStatus(d, resp.Status),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("Error saving the specifies CloudTable cluster (%s) to state: %s", d.Id(), mErr)
	}

	return nil
}

func resourceCloudTableClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.CloudtableV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating CloudTable v2 client: %s", err)
	}

	err = clusters.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting the specified CloudTable cluster (%s): %s", d.Id(), err)
	}

	stateRef := stateRefresh{
		Pending:      []string{"Running", "Abnormal"},
		Target:       []string{"Deleted"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        6 * time.Second,
		PollInterval: 5 * time.Second,
		// Due to the design of the API, the cluster cannot be queried for a few time after the delete request is sent.
		NotFoundChecks: 1,
	}
	if err := waitForCloudTableClusterStateRefresh(ctx, client, d.Id(), stateRef); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func waitForCloudTableClusterStateRefresh(ctx context.Context, c *golangsdk.ServiceClient, clusterId string,
	s stateRefresh) diag.Diagnostics {
	stateConf := &resource.StateChangeConf{
		Pending:        s.Pending,
		Target:         s.Target,
		Refresh:        clusterStateRefreshFunc(c, clusterId),
		Timeout:        s.Timeout,
		Delay:          s.Delay,
		PollInterval:   s.PollInterval,
		NotFoundChecks: s.NotFoundChecks,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("waiting for status of the cluster (%s) to complete (%s) timeout: %s",
			clusterId, s.Target, err)
	}

	return nil
}

func clusterStateRefreshFunc(c *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := clusters.Get(c, clusterId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return resp, "Deleted", nil
			}
			return resp, "Error", nil
		}
		if val, ok := stateCodes[resp.Status]; ok {
			return resp, val, nil
		}
		log.Printf("[ERROR] The status is abnormal, the code is %+v", resp.Status)
		return resp, "Error", nil
	}
}
