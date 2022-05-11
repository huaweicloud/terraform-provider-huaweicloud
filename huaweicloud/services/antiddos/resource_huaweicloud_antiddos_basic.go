package antiddos

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	antiddossdk "github.com/chnsz/golangsdk/openstack/antiddos/v1/antiddos"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var (
	trafficThresholds  = []int{10, 30, 50, 70, 100, 120, 150, 200, 250, 300, 1000}
	trafficThresholdID = map[int]int{
		10:   1,
		30:   2,
		50:   3,
		70:   4,
		100:  5,
		120:  99,
		150:  6,
		200:  7,
		250:  8,
		300:  9,
		1000: 88,
	}
	trafficThresholdBandwidth = map[int]int{
		1:  10,
		2:  30,
		3:  50,
		4:  70,
		5:  100,
		6:  150,
		7:  200,
		8:  250,
		9:  300,
		88: 1000,
		99: 120,
	}
)

// ResourceCloudNativeAntiDdos is the imple of huaweicloud_antiddos_basic
func ResourceCloudNativeAntiDdos() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudNativeAntiDdosUpdate,
		ReadContext:   resourceCloudNativeAntiDdosRead,
		UpdateContext: resourceCloudNativeAntiDdosUpdate,
		DeleteContext: resourceCloudNativeAntiDdosDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"traffic_threshold": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice(trafficThresholds),
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudNativeAntiDdosUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.AntiDDosV1Client(region)
	if err != nil {
		return diag.Errorf("error creating antiddos client: %s", err)
	}

	eipID := d.Get("eip_id").(string)
	thresholdID := getTrafficThresholdID(d.Get("traffic_threshold").(int))

	if err := updateAntiDdosTrafficThreshold(ctx, d, client, eipID, thresholdID, true); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(eipID)
	return resourceCloudNativeAntiDdosRead(ctx, d, meta)
}

func resourceCloudNativeAntiDdosRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.AntiDDosV1Client(region)
	if err != nil {
		return diag.Errorf("error creating antiddos client: %s", err)
	}
	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating VPC client: %s", err)
	}

	eIP, err := eips.Get(vpcClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving cloud native AntiDdos")
	}

	listStatusOpts := antiddossdk.ListStatusOpts{
		Ip: eIP.PublicAddress,
	}
	results, err := antiddossdk.ListStatus(client, listStatusOpts)
	if err != nil {
		return diag.Errorf("error retrieving cloud native AntiDdos: %s", err)
	}

	if len(results) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving cloud native AntiDdos")
	}

	ddosStatus := results[0]
	log.Printf("[DEBUG] Retrieved cloud native AntiDdos %s: %#v", d.Id(), ddosStatus)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("eip_id", ddosStatus.FloatingIpId),
		d.Set("public_ip", ddosStatus.FloatingIpAddress),
		d.Set("traffic_threshold", getTrafficThresholdBandwidth(ddosStatus.TrafficThreshold)),
		d.Set("status", ddosStatus.Status),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting resource: %s", mErr)
	}
	return nil
}

func resourceCloudNativeAntiDdosDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.AntiDDosV1Client(region)
	if err != nil {
		return diag.Errorf("error creating antiddos client: %s", err)
	}

	if err := updateAntiDdosTrafficThreshold(ctx, d, client, d.Id(), 99, false); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getTrafficThresholdID(bandwidth int) int {
	return trafficThresholdID[bandwidth]
}

func getTrafficThresholdBandwidth(id int) int {
	bandwidth, ok := trafficThresholdBandwidth[id]
	if !ok {
		bandwidth = id
	}
	return bandwidth
}

func updateAntiDdosTrafficThreshold(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	antiddosID string, threshold int, check bool) error {
	// calling antiddossdk.Get method to get other requied parameters
	preProtection, err := antiddossdk.Get(client, antiddosID).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok && !check {
			return nil
		}

		return fmt.Errorf("error retrieving cloud native AntiDdos: %s", err)
	}

	if preProtection.TrafficPosId != threshold {
		updateOpts := antiddossdk.UpdateOpts{
			EnableL7:            preProtection.EnableL7,
			HttpRequestPosId:    preProtection.HttpRequestPosId,
			CleaningAccessPosId: preProtection.CleaningAccessPosId,
			AppTypeId:           preProtection.AppTypeId,
			TrafficPosId:        threshold,
		}

		log.Printf("[DEBUG] AntiDdos updating options: %#v", updateOpts)
		if _, err := antiddossdk.Update(client, antiddosID, updateOpts).Extract(); err != nil {
			return fmt.Errorf("error updating AntiDdos: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"configging"},
			Target:       []string{"normal"},
			Refresh:      getAntiDdosStatus(client, antiddosID),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        5 * time.Second,
			PollInterval: 5 * time.Second,
		}

		_, stateErr := stateConf.WaitForStateContext(ctx)
		if stateErr != nil {
			return fmt.Errorf("error waiting for AntiDdos to become normal: %s", stateErr)
		}
	}

	return nil
}

func getAntiDdosStatus(client *golangsdk.ServiceClient, antiddosID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := antiddossdk.GetStatus(client, antiddosID).Extract()
		if err != nil {
			return nil, "", err
		}

		return s, s.Status, nil
	}
}
