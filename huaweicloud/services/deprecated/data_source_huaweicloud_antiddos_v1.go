package deprecated

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/antiddos/v1/antiddos"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceAntiDdosV1() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceAntiDdosV1Read,
		DeprecationMessage: "this is deprecated",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"floating_ip_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"floating_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"period_start": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"bps_attack": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"bps_in": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"total_bps": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"pps_in": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"pps_attack": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"total_pps": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"start_time": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"end_time": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"traffic_cleaning_status": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"trigger_bps": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"trigger_pps": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"trigger_http_pps": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceAntiDdosV1Read(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	antiddosClient, err := cfg.AntiDDosV1Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating antiddos client: %s", err)
	}

	listStatusOpts := antiddos.ListStatusOpts{
		FloatingIpId: d.Get("floating_ip_id").(string),
		Status:       d.Get("status").(string),
		Ip:           d.Get("floating_ip_address").(string),
	}

	refinedAntiddos, err := antiddos.ListStatus(antiddosClient, listStatusOpts)
	if err != nil {
		return fmt.Errorf("unable to retrieve the defense status of EIP, defense is not configured: %s", err)
	}

	if len(refinedAntiddos) < 1 {
		return errors.New("your query returned no results, please change your search criteria and try again")
	}

	if len(refinedAntiddos) > 1 {
		return errors.New("your query returned more than one result, please try a more specific search criteria")
	}

	ddosStatus := refinedAntiddos[0]

	log.Printf("[INFO] Retrieved defense status of EIP %s using given filter", ddosStatus.FloatingIpId)

	d.SetId(ddosStatus.FloatingIpId)

	d.Set("floating_ip_id", ddosStatus.FloatingIpId)
	d.Set("floating_ip_address", ddosStatus.FloatingIpAddress)
	d.Set("network_type", ddosStatus.NetworkType)
	d.Set("status", ddosStatus.Status)

	d.Set("region", cfg.GetRegion(d))

	traffic, err := antiddos.DailyReport(antiddosClient, ddosStatus.FloatingIpId).Extract()
	log.Printf("traffic %#v", traffic)
	if err != nil {
		return fmt.Errorf("unable to retrieve the traffic of a specified EIP, defense is not configured: %s", err)
	}

	periodStart := make([]int, 0)
	for _, param := range traffic {
		periodStart = append(periodStart, param.PeriodStart)
	}
	d.Set("period_start", periodStart)

	bpsIn := make([]int, 0)
	for _, param := range traffic {
		bpsIn = append(bpsIn, param.BpsIn)
	}
	d.Set("bps_in", bpsIn)

	bpsAttack := make([]int, 0)
	for _, param := range traffic {
		bpsAttack = append(bpsAttack, param.BpsAttack)
	}
	d.Set("bps_attack", bpsAttack)

	totalBps := make([]int, 0)
	for _, param := range traffic {
		totalBps = append(totalBps, param.TotalBps)
	}
	d.Set("total_bps", totalBps)

	ppsIn := make([]int, 0)
	for _, param := range traffic {
		ppsIn = append(ppsIn, param.PpsIn)
	}
	d.Set("pps_in", ppsIn)

	ppsAttack := make([]int, 0)
	for _, param := range traffic {
		ppsAttack = append(ppsAttack, param.PpsAttack)
	}
	d.Set("pps_attack", ppsAttack)

	totalPps := make([]int, 0)
	for _, param := range traffic {
		totalPps = append(totalPps, param.TotalPps)
	}
	d.Set("total_pps", totalPps)

	listEventOpts := antiddos.ListLogsOpts{}
	event, err := antiddos.ListLogs(antiddosClient, ddosStatus.FloatingIpId, listEventOpts).Extract()
	log.Printf("event %#v", event)
	if err != nil {
		return fmt.Errorf("unable to retrieve the event of a specified EIP, defense is not configured: %s", err)
	}

	startTime := make([]int, 0)
	for _, param := range event {
		startTime = append(startTime, param.StartTime)
	}
	d.Set("start_time", startTime)

	endTime := make([]int, 0)
	for _, param := range event {
		endTime = append(endTime, param.EndTime)
	}
	d.Set("end_time", endTime)

	cleaningStatus := make([]int, 0)
	for _, param := range event {
		cleaningStatus = append(cleaningStatus, param.Status)
	}
	d.Set("traffic_cleaning_status", cleaningStatus)

	triggerBps := make([]int, 0)
	for _, param := range event {
		triggerBps = append(triggerBps, param.TriggerBps)
	}
	d.Set("trigger_bps", triggerBps)

	triggerPps := make([]int, 0)
	for _, param := range event {
		triggerPps = append(triggerPps, param.TriggerPps)
	}
	d.Set("trigger_pps", triggerPps)

	triggerHttpPps := make([]int, 0)
	for _, param := range event {
		triggerHttpPps = append(triggerHttpPps, param.TriggerHttpPps)
	}
	d.Set("trigger_http_pps", triggerHttpPps)

	return nil
}
