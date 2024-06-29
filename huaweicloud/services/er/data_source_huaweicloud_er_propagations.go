package er

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/er/v3/propagations"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ER GET /v3/{project_id}/enterprise-router/{er_id}/route-tables/{route_table_id}/propagations
func DataSourcePropagations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePropagationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attachment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attachment_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"propagations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     propagationSchema(),
			},
		},
	}
}

func propagationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildPropagationListOpts(d *schema.ResourceData) propagations.ListOpts {
	opts := propagations.ListOpts{}
	if attachmentId, ok := d.GetOk("attachment_id"); ok {
		opts.AttachmentIds = []string{attachmentId.(string)}
	}

	if attachmentType, ok := d.GetOk("attachment_type"); ok {
		opts.ResourceTypes = []string{attachmentType.(string)}
	}

	if status, ok := d.GetOk("status"); ok {
		opts.Statuses = []string{status.(string)}
	}

	return opts
}

func flattenPropagations(all []propagations.Propagation) []map[string]interface{} {
	if len(all) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(all))
	for i, propagation := range all {
		result[i] = map[string]interface{}{
			"id":              propagation.ID,
			"instance_id":     propagation.InstanceId,
			"route_table_id":  propagation.RouteTableId,
			"attachment_id":   propagation.AttachmentId,
			"attachment_type": propagation.ResourceType,
			"resource_id":     propagation.ResourceId,
			"route_policy_id": propagation.RoutePolicy.ImportPoilicyId,
			"status":          propagation.Status,
			// The time results are not the time in RF3339 format without milliseconds.
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(propagation.CreatedAt)/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(propagation.UpdatedAt)/1000, false),
		}
	}
	return result
}

func dataSourcePropagationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	opts := buildPropagationListOpts(d)
	resp, err := propagations.List(client, d.Get("instance_id").(string), d.Get("route_table_id").(string), opts)
	if err != nil {
		return diag.Errorf("error retrieving propagations: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("propagations", flattenPropagations(resp)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving data source fields of ER propagations: %s", mErr)
	}
	return nil
}
