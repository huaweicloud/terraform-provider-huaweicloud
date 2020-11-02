package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/antiddos/v1/antiddos"
)

func dataSourceAntiDdosV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAntiDdosV1Read,

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
	config := meta.(*Config)
	antiddosClient, err := config.antiddosV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating antiddos client: %s", err)
	}

	listStatusOpts := antiddos.ListStatusOpts{
		FloatingIpId: d.Get("floating_ip_id").(string),
		Status:       d.Get("status").(string),
		Ip:           d.Get("floating_ip_address").(string),
	}

	refinedAntiddos, err := antiddos.ListStatus(antiddosClient, listStatusOpts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve the defense status of  EIP, defense is not configured.: %s", err)
	}

	if len(refinedAntiddos) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedAntiddos) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	ddosStatus := refinedAntiddos[0]

	log.Printf("[INFO] Retrieved defense status of  EIP %s using given filter", ddosStatus.FloatingIpId)

	d.SetId(ddosStatus.FloatingIpId)

	d.Set("floating_ip_id", ddosStatus.FloatingIpId)
	d.Set("floating_ip_address", ddosStatus.FloatingIpAddress)
	d.Set("network_type", ddosStatus.NetworkType)
	d.Set("status", ddosStatus.Status)

	d.Set("region", GetRegion(d, config))

	traffic, err := antiddos.DailyReport(antiddosClient, ddosStatus.FloatingIpId).Extract()
	log.Printf("traffic %#v", traffic)
	if err != nil {
		return fmt.Errorf("Unable to retrieve the traffic of a specified EIP, defense is not configured: %s", err)
	}

	period_start := make([]int, 0)
	for _, param := range traffic {
		period_start = append(period_start, param.PeriodStart)
	}
	d.Set("period_start", period_start)

	bps_in := make([]int, 0)
	for _, param := range traffic {
		bps_in = append(bps_in, param.BpsIn)
	}
	d.Set("bps_in", bps_in)

	bps_attack := make([]int, 0)
	for _, param := range traffic {
		bps_attack = append(bps_attack, param.BpsAttack)
	}
	d.Set("bps_attack", bps_attack)

	total_bps := make([]int, 0)
	for _, param := range traffic {
		total_bps = append(total_bps, param.TotalBps)
	}
	d.Set("total_bps", total_bps)

	pps_in := make([]int, 0)
	for _, param := range traffic {
		pps_in = append(pps_in, param.PpsIn)
	}
	d.Set("pps_in", pps_in)

	pps_attack := make([]int, 0)
	for _, param := range traffic {
		pps_attack = append(pps_attack, param.PpsAttack)
	}
	d.Set("pps_attack", pps_attack)

	total_pps := make([]int, 0)
	for _, param := range traffic {
		total_pps = append(total_pps, param.TotalPps)
	}
	d.Set("total_pps", total_pps)

	listEventOpts := antiddos.ListLogsOpts{}
	event, err := antiddos.ListLogs(antiddosClient, ddosStatus.FloatingIpId, listEventOpts).Extract()
	log.Printf("event %#v", event)
	if err != nil {
		return fmt.Errorf("Unable to retrieve the event of a specified EIP, defense is not configured: %s", err)
	}

	start_time := make([]int, 0)
	for _, param := range event {
		start_time = append(start_time, param.StartTime)
	}
	d.Set("start_time", start_time)

	end_time := make([]int, 0)
	for _, param := range event {
		end_time = append(end_time, param.EndTime)
	}
	d.Set("end_time", end_time)

	cleaning_status := make([]int, 0)
	for _, param := range event {
		cleaning_status = append(cleaning_status, param.Status)
	}
	d.Set("traffic_cleaning_status", cleaning_status)

	trigger_bps := make([]int, 0)
	for _, param := range event {
		trigger_bps = append(trigger_bps, param.TriggerBps)
	}
	d.Set("trigger_bps", trigger_bps)

	trigger_pps := make([]int, 0)
	for _, param := range event {
		trigger_pps = append(trigger_pps, param.TriggerPps)
	}
	d.Set("trigger_pps", trigger_pps)

	trigger_http_pps := make([]int, 0)
	for _, param := range event {
		trigger_http_pps = append(trigger_http_pps, param.TriggerHttpPps)
	}
	d.Set("trigger_http_pps", trigger_http_pps)

	return nil

}
