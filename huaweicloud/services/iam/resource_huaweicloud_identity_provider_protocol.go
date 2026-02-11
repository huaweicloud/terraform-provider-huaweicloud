package iam

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourceV3ProviderProtocol
// @API IAM PUT /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}
// @API IAM GET /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}
// @API IAM PATCH /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}
// @API IAM DELETE /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}
func ResourceV3ProviderProtocol() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3ProviderProtocolCreate,
		ReadContext:   resourceV3ProviderProtocolRead,
		UpdateContext: resourceV3ProviderProtocolUpdate,
		DeleteContext: resourceV3ProviderProtocolDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV3ProviderProtocolImportState,
		},

		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the identity provider used to manage the protocol.`,
			},
			"protocol_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The identity protocol of the identity provider.`,
			},
			"mapping_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The mapping_id for the protocol.`,
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
							Description: `The resource link.`,
						},
						"identity_provider": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The identity provider resource link.`,
						},
					},
				},
			},
		},
	}
}

func resourceV3ProviderProtocolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	idpId := d.Get("provider_id").(string)
	protocolId := d.Get("protocol_id").(string)
	mappingId := d.Get("mapping_id").(string)

	protocolPath := getProtocolPath(client.Endpoint, idpId, protocolId)
	options := getProtocolRequestOptsWithBody(mappingId)
	_, err = client.Request("PUT", protocolPath, &options)
	conflictMsg := "Conflict occurred attempting to store federation_protocol"
	if err != nil {
		if strings.Contains(err.Error(), "got 409") && strings.Contains(err.Error(), conflictMsg) {
			log.Printf("protocol `%s` of identity provider `%s` has already existed. Now Update it with mapping `%s`",
				idpId, protocolId, mappingId)
			return resourceV3ProviderProtocolUpdate(ctx, d, meta)
		}
		return diag.Errorf("CreateProtocol error : %s", err)
	}
	d.SetId(fmt.Sprintf("%s:%s", idpId, protocolId))
	return resourceV3ProviderProtocolRead(ctx, d, meta)
}

func resourceV3ProviderProtocolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	protocolPath := getProtocolPath(client.Endpoint, d.Get("provider_id").(string), d.Get("protocol_id").(string))
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := client.Request("GET", protocolPath, &options)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving identity provider protocol")
	}

	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}

	links := append(make([]interface{}, 0, 1), map[string]interface{}{
		"self":              utils.PathSearch("protocol.links.self", respBody, ""),
		"identity_provider": utils.PathSearch("protocol.links.identity_provider", respBody, ""),
	})

	mErr := multierror.Append(nil,
		d.Set("protocol_id", utils.PathSearch("protocol.id", respBody, "")),
		d.Set("mapping_id", utils.PathSearch("protocol.mapping_id", respBody, "")),
		d.Set("links", links),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting identity provider protocol: %s", err)
	}
	return nil
}

func resourceV3ProviderProtocolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	idpId := d.Get("provider_id").(string)
	protocolId := d.Get("protocol_id").(string)
	mappingId := d.Get("mapping_id").(string)
	protocolPath := getProtocolPath(client.Endpoint, idpId, protocolId)
	options := getProtocolRequestOptsWithBody(mappingId)
	_, err = client.Request("PATCH", protocolPath, &options)
	if err != nil {
		return diag.Errorf("UpdateProtocol error : %s", err)
	}
	return resourceV3ProviderProtocolRead(ctx, d, meta)
}

func resourceV3ProviderProtocolDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IAMV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	protocolPath := getProtocolPath(client.Endpoint, d.Get("provider_id").(string), d.Get("protocol_id").(string))
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	_, err = client.Request("DELETE", protocolPath, &options)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting identity provider protocol")
	}
	return nil
}

func resourceV3ProviderProtocolImportState(
	_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id: %s, id must be {provider_id}:{protocol_id}", d.Id())
	}
	mErr := multierror.Append(nil,
		d.Set("provider_id", parts[0]),
		d.Set("protocol_id", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import identity provider protocol, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

func getProtocolPath(endpoint string, idpId string, protocolId string) string {
	protocolPath := endpoint + "v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}"
	protocolPath = strings.ReplaceAll(protocolPath, "{idp_id}", idpId)
	protocolPath = strings.ReplaceAll(protocolPath, "{protocol_id}", protocolId)
	return protocolPath
}

func getProtocolRequestOptsWithBody(mappingId string) golangsdk.RequestOpts {
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	if mappingId != "" {
		options.JSONBody = map[string]interface{}{
			"protocol": map[string]string{
				"mapping_id": mappingId},
		}
	} else {
		options.JSONBody = map[string]interface{}{
			"protocol": map[string]string{},
		}
	}
	return options
}
