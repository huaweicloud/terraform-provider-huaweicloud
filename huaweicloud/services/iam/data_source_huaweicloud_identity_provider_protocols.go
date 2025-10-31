package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceIamIdentityProviderProtocols
// @API IAM GET /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols
// @API IAM GET /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}
func DataSourceIamIdentityProviderProtocols() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityProviderProtocolsRead,

		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of an identity provider.",
			},
			"protocol_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The protocol id.",
			},

			"protocols": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Protocol Information List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol Id",
						},
						"mapping_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mapping Id",
						},
						"links": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The links of protocol.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"self": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `resource link.`,
									},
									"identity_provider": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `identity provider resource link.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityProviderProtocolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamV3Client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	protocolId := d.Get("protocol_id").(string)
	if protocolId == "" {
		return listProviderProtocols(iamV3Client, d)
	}
	return showProviderProtocol(iamV3Client, d)
}

func listProviderProtocols(iamV3Client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	protocolPath := iamV3Client.Endpoint + "v3/OS-FEDERATION/identity_providers/{idp_id}/protocols"
	protocolPath = strings.ReplaceAll(protocolPath, "{idp_id}", d.Get("provider_id").(string))
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := iamV3Client.Request("GET", protocolPath, &options)
	if err != nil {
		return diag.Errorf("ListProtocol error : %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	protocolsBody := utils.PathSearch("protocols", respBody, make([]interface{}, 0)).([]interface{})
	protocols := make([]interface{}, 0, len(protocolsBody))
	for _, protocol := range protocolsBody {
		protocols = append(protocols, flattenProtocol(protocol))
	}
	if err = d.Set("protocols", protocols); err != nil {
		return diag.Errorf("error setting protocols fields: %s", err)
	}
	return nil
}

func showProviderProtocol(iamV3Client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	protocolPath := getProtocolPath(iamV3Client.Endpoint, d.Get("provider_id").(string), d.Get("protocol_id").(string))
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := iamV3Client.Request("GET", protocolPath, &options)
	if err != nil {
		return diag.Errorf("ShowProtocol error : %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	protocolBody := utils.PathSearch("protocol", respBody, make([]interface{}, 0))
	protocols := append(make([]interface{}, 0, 1), flattenProtocol(protocolBody))
	if err = d.Set("protocols", protocols); err != nil {
		return diag.Errorf("error setting protocols fields: %s", err)
	}
	return nil
}

func flattenProtocol(protocolModel interface{}) map[string]interface{} {
	if protocolModel == nil {
		return nil
	}
	protocol := make(map[string]interface{})
	protocol["id"] = utils.PathSearch("id", protocolModel, "")
	protocol["mapping_id"] = utils.PathSearch("mapping_id", protocolModel, "")
	links := append(make([]interface{}, 0, 1), map[string]string{
		"self":              utils.PathSearch("links.self", protocolModel, "").(string),
		"identity_provider": utils.PathSearch("links.identity_provider", protocolModel, "").(string),
	})
	protocol["links"] = links
	return protocol
}
