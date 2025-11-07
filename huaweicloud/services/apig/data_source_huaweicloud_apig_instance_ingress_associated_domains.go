package apig

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports/{ingress_port_id}/domains
func DataSourceInstanceIngressAssociatedDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceIngressAssociatedDomainsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the domains are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the ingress port belongs.`,
			},
			"ingress_port_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the custom ingress port.`,
			},

			// Optional parameters.
			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The domain name that uses the ingress port.`,
			},

			// Attributes.
			"domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of domain information bound to the ingress port that matched the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name bound to the ingress port.`,
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the API group bound to the ingress port.`,
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the API group bound to the ingress port.`,
						},
					},
				},
			},
		},
	}
}

func buildInstanceIngressAssociatedDomainsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("domain_name"); ok {
		res = fmt.Sprintf("%s&domain_name=%v", res, v)
	}

	return res
}

func listInstanceIngressAssociatedDomains(client *golangsdk.ServiceClient, d *schema.ResourceData, instanceId,
	ingressPortId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports/{ingress_port_id}/domains?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{ingress_port_id}", ingressPortId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildInstanceIngressAssociatedDomainsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		domainInfos := utils.PathSearch("domain_infos", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, domainInfos...)
		if len(domainInfos) < limit {
			break
		}
		offset += len(domainInfos)
	}

	return result, nil
}

func flattenInstanceIngressAssociatedDomainInfos(domainInfos []interface{}) []map[string]interface{} {
	if len(domainInfos) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(domainInfos))
	for _, domainInfo := range domainInfos {
		result = append(result, map[string]interface{}{
			"name":       utils.PathSearch("domain_name", domainInfo, nil),
			"group_id":   utils.PathSearch("group_id", domainInfo, nil),
			"group_name": utils.PathSearch("group_name", domainInfo, nil),
		})
	}

	return result
}

func dataSourceInstanceIngressAssociatedDomainsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	ingressPortId := d.Get("ingress_port_id").(string)

	resp, err := listInstanceIngressAssociatedDomains(client, d, instanceId, ingressPortId)
	if err != nil {
		return diag.Errorf("error querying APIG instance ingress associated domains: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("domains", flattenInstanceIngressAssociatedDomainInfos(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
