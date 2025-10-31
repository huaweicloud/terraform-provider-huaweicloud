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

// ResourceIdentityProviderProtocol
// @API IAM PUT /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}
// @API IAM GET /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}
// @API IAM PATCH /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}
// @API IAM DELETE /v3/OS-FEDERATION/identity_providers/{idp_id}/protocols/{protocol_id}
func ResourceIdentityProviderProtocol() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityProviderProtocolCreate,
		ReadContext:   resourceIdentityProviderProtocolRead,
		UpdateContext: resourceIdentityProviderProtocolUpdate,
		DeleteContext: resourceIdentityProviderProtocolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityProviderProtocolImportState,
		},

		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of an identity provider.",
			},
			"protocol_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The protocol ID to be registered. The content of this field is `saml` or `oidc`.",
			},
			"mapping_id": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "When the identity provider type is `iam_user_sso`, there is no need to bind a " +
					"mapping ID, and this field does not need to be passed; otherwise, this field is mandatory.",
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
	}
}

func resourceIdentityProviderProtocolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamV3Client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	idpId := d.Get("provider_id").(string)
	protocolId := d.Get("protocol_id").(string)
	mappingId := d.Get("mapping_id").(string)

	protocolPath := getProtocolPath(iamV3Client.Endpoint, idpId, protocolId)
	options := getProtocolRequestOptsWithBody(mappingId)
	_, err = iamV3Client.Request("PUT", protocolPath, &options)
	conflictMsg := "Conflict occurred attempting to store federation_protocol"
	if err != nil {
		if strings.Contains(err.Error(), "got 409") && strings.Contains(err.Error(), conflictMsg) {
			log.Printf("protocol `%s` of identity provider `%s` has already existed. Now Update it with mapping `%s`",
				idpId, protocolId, mappingId)
			return resourceIdentityProviderProtocolUpdate(ctx, d, meta)
		}
		return diag.Errorf("CreateProtocol error : %s", err)
	}
	d.SetId(fmt.Sprintf("%s:%s", idpId, protocolId))
	return resourceIdentityProviderProtocolRead(ctx, d, meta)
}

func resourceIdentityProviderProtocolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iamV3Client, err := conf.IAMV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	protocolPath := getProtocolPath(iamV3Client.Endpoint, d.Get("provider_id").(string), d.Get("protocol_id").(string))
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := iamV3Client.Request("GET", protocolPath, &options)
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

func resourceIdentityProviderProtocolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iamV3Client, err := conf.IAMV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	idpId := d.Get("provider_id").(string)
	protocolId := d.Get("protocol_id").(string)
	mappingId := d.Get("mapping_id").(string)
	protocolPath := getProtocolPath(iamV3Client.Endpoint, idpId, protocolId)
	options := getProtocolRequestOptsWithBody(mappingId)
	_, err = iamV3Client.Request("PATCH", protocolPath, &options)
	if err != nil {
		return diag.Errorf("UpdateProtocol error : %s", err)
	}
	return resourceIdentityProviderProtocolRead(ctx, d, meta)
}

func resourceIdentityProviderProtocolDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iamV3Client, err := conf.IAMV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	protocolPath := getProtocolPath(iamV3Client.Endpoint, d.Get("provider_id").(string), d.Get("protocol_id").(string))
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	_, err = iamV3Client.Request("DELETE", protocolPath, &options)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting identity provider protocol")
	}
	return nil
}

func resourceIdentityProviderProtocolImportState(
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
