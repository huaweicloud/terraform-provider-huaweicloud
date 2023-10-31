package antiddos

import (
	"context"
	"fmt"
	"log"
	"math"
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
				ForceNew: true,
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

func waitForAntiDDoSAvailable(ctx context.Context, client *golangsdk.ServiceClient, antiDDoSId string,
	timeout time.Duration, isDeleteCheck bool) (*antiddossdk.GetResponse, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := antiddossdk.Get(client, antiDDoSId).Extract()
			if err != nil {
				analyzedErr := parseAntiDDoSQueryError(err)
				if _, ok := analyzedErr.(golangsdk.ErrDefault404); !ok {
					return resp, "ERROR", analyzedErr
				}
				// For deletion operations, the 404 error returned by the query is considered to be the completion of
				// the deletion, while for other operations, this error is considered to require continued waiting.
				if !isDeleteCheck {
					return resp, "PENDING", nil
				}
			}
			return resp, "COMPLETED", ResourceCloudNativeAntiDdos().Importer.InternalValidate()
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	resp, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return nil, fmt.Errorf("error waiting for AntiDDoS (%s) to become expected status: %s", antiDDoSId, stateErr)
	}
	preProtection, ok := resp.(*antiddossdk.GetResponse)
	if !ok || resp == nil {
		return preProtection, fmt.Errorf("invalid result type of the AntiDDoS query, want '*antiddossdk.GetResponse', but got '%T'", resp)
	}
	return preProtection, nil
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

	if err := updateAntiDdosTrafficThreshold(ctx, d, client, eipID, thresholdID, false); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(eipID)
	return resourceCloudNativeAntiDdosRead(ctx, d, meta)
}

func resourceCloudNativeAntiDdosRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if err := updateAntiDdosTrafficThreshold(ctx, d, client, d.Id(), 99, true); err != nil {
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
	antiDDoSId string, threshold int, check bool) error {
	preProtection, err := waitForAntiDDoSAvailable(ctx, client, antiDDoSId, d.Timeout(schema.TimeoutUpdate), check)
	if err != nil {
		return err
	}

	if preProtection.TrafficPosId != threshold {
		updateOpts := antiddossdk.UpdateOpts{
			EnableL7:         preProtection.EnableL7,
			HttpRequestPosId: preProtection.HttpRequestPosId,
			// make sure the CleaningAccessPosId not larger than 8
			// CleaningAccessPosId has no practical meaning in the request
			// this will avoid error in partners cloud
			CleaningAccessPosId: int(math.Min(float64(preProtection.CleaningAccessPosId), 8)),
			AppTypeId:           preProtection.AppTypeId,
			TrafficPosId:        threshold,
		}

		log.Printf("[DEBUG] AntiDdos updating options: %#v", updateOpts)
		if _, err := antiddossdk.Update(client, antiDDoSId, updateOpts).Extract(); err != nil {
			return fmt.Errorf("error updating AntiDdos: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"configging"},
			Target:       []string{"normal"},
			Refresh:      getAntiDdosStatus(client, antiDDoSId),
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

func parseAntiDDoSQueryError(respErr error) error {
	if errCode, ok := respErr.(golangsdk.ErrDefault403); ok {
		resp, err := common.ParseErrorMsg(errCode.Body)
		if err == nil && (resp.ErrorCode == "10001020" && resp.ErrorMsg == "IPID is invalid") {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return respErr
}
