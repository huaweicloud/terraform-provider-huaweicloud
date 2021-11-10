package dis

import (
	"context"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/dis/v2/streams"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceDisPartitions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDisPartitionRead,
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

func dataSourceDisPartitionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DisV2Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud DIS client: %s", err)
	}
	name := d.Get("stream_name").(string)
	var result []map[string]interface{}
	opts := streams.GetOpts{}
	for {
		rst, gErr := streams.Get(client, name, opts)
		if gErr != nil {
			return fmtp.DiagErrorf("Error query the partitions of DIS stream, err=%s", gErr)
		}

		if len(rst.Partitions) < 1 {
			break
		}

		for _, i := range rst.Partitions {
			result = append(result, map[string]interface{}{
				"id":                    i.PartitionId,
				"status":                i.Status,
				"hash_range":            i.HashRange,
				"sequence_number_range": i.SequenceNumberRange,
			})
		}

		if rst.HasMorePartitions {
			opts.StartPartitionId = rst.Partitions[len(rst.Partitions)-1].PartitionId
		} else {
			break
		}
	}

	dErr := d.Set("partitions", result)
	if dErr != nil {
		return fmtp.DiagErrorf("Error set partitions, err=%s", dErr)
	}

	d.SetId(d.Get("stream_name").(string))
	return nil
}
