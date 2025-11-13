package apig

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports
func DataSourceInstanceIngressPorts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceIngressPortsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the ingress ports are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the ingress ports belong.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The protocol of the ingress port to be queried.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The port number of the ingress port to be queried.`,
			},
			"ingress_ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the ingress port.`,
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The protocol of the ingress port.`,
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The port number of the ingress port.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the ingress port.`,
						},
					},
				},
				Description: `The list of the ingress ports that matched filter parameters.`,
			},
		},
	}
}

func flattenInstanceIngressPorts(ingressPorts []interface{}) []map[string]interface{} {
	if len(ingressPorts) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(ingressPorts))
	for _, port := range ingressPorts {
		result = append(result, map[string]interface{}{
			"id":       utils.PathSearch("ingress_port_id", port, nil),
			"protocol": utils.PathSearch("protocol", port, nil),
			"port":     utils.PathSearch("ingress_port", port, nil),
			"status":   utils.PathSearch("status", port, nil),
		})
	}

	return result
}

func dataSourceInstanceIngressPortsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	ingressPorts, err := listInstanceIngressPorts(client, instanceId, d)
	if err != nil {
		return diag.Errorf("error querying APIG instance ingress ports: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("ingress_ports", flattenInstanceIngressPorts(ingressPorts)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
