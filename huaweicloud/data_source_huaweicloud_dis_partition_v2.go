package huaweicloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
)

func dataSourceDisPartitionV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDisPartitionV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"stream_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"partitions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hash_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sequence_number_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDisPartitionV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.disV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	url, err := replaceVars(d, "streams/{stream_name}", nil)
	if err != nil {
		return err
	}
	url = client.ServiceURL(url)
	url1 := url

	result := make([]interface{}, 0, 50)
	for {
		r := golangsdk.Result{}
		_, r.Err = client.Get(url1, &r.Body, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{"Content-Type": "application/json"}})
		if r.Err != nil {
			return fmt.Errorf("Error running api(read) for resource(DisStreamV2), err=%s", r.Err)
		}

		v, err := navigateValue(r.Body, []string{"partitions"}, nil)
		if err != nil {
			return err
		}
		partitions, ok := v.([]interface{})
		if !ok {
			break
		}

		hasMore, err := navigateValue(r.Body, []string{"has_more_partitions"}, nil)
		if err != nil {
			return err
		}

		for _, i := range partitions {
			val := i.(map[string]interface{})
			result = append(result, map[string]interface{}{
				"id":                    val["partition_id"],
				"status":                val["status"],
				"hash_range":            val["hash_range"],
				"sequence_number_range": val["sequence_number_range"],
			})
		}

		if m, ok := hasMore.(bool); ok && m {
			url1 = url + "?start_partitionId=" + result[len(result)-1].(map[string]interface{})["id"].(string)
		} else {
			break
		}
	}

	d.SetId(d.Get("stream_name").(string))
	d.Set("partitions", result)
	return nil
}
