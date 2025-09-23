package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var domainAssociateNonUpdatableParams = []string{
	"instance_id",
	"group_id",
	"url_domain",
}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}/domains
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}/domains/{domain_id}
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}/domains/{domain_id}
func ResourceGroupDomainAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupDomainAssociateCreate,
		ReadContext:   resourceGroupDomainAssociateRead,
		UpdateContext: resourceGroupDomainAssociateUpdate,
		DeleteContext: resourceGroupDomainAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGroupDomainAssociateImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(domainAssociateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the APIG group and domain are located.",
			},
			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the dedicated instance to which the group belongs.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the group to associate with the domain name.",
			},
			"url_domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The associated domain name.",
			},
			// Optional parameters.
			"min_ssl_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The minimum SSL protocol version.",
			},
			"ingress_http_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The HTTP protocol inbound access port bound to the domain name.",
			},
			"ingress_https_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The HTTPS protocol inbound access port bound to the domain name.",
			},
			"is_http_redirect_to_https": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable redirection from HTTP to HTTPS.`,
			},
			// Internal parameters/attributes.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The associated domain ID.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func resourceGroupDomainAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}/domains"
		instanceId = d.Get("instance_id").(string)
		groupId    = d.Get("group_id").(string)
		urlDomain  = d.Get("url_domain").(string)
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{group_id}", groupId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"url_domain":                utils.ValueIgnoreEmpty(urlDomain),
			"min_ssl_version":           utils.ValueIgnoreEmpty(d.Get("min_ssl_version")),
			"is_http_redirect_to_https": d.Get("is_http_redirect_to_https"),
			"ingress_http_port":         utils.ValueIgnoreEmpty(d.Get("ingress_http_port")),
			"ingress_https_port":        utils.ValueIgnoreEmpty(d.Get("ingress_https_port")),
		}),
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error associating the domain name with the group (%s): %s", groupId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceId, groupId, urlDomain))

	return resourceGroupDomainAssociateRead(ctx, d, meta)
}

func getGroupAssociatedDomains(client *golangsdk.ServiceClient, instanceId, groupId string) ([]interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{group_id}", groupId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("url_domains", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func GetGroupAssociatedDomainByUrl(client *golangsdk.ServiceClient, instanceId, groupId, urlDomain string) (interface{}, error) {
	urlDomains, err := getGroupAssociatedDomains(client, instanceId, groupId)
	if err != nil {
		return nil, err
	}

	result := utils.PathSearch(fmt.Sprintf("[?domain=='%s']|[0]", urlDomain), urlDomains, nil)
	if result == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the domain (%s) does not bound to the API group (%s)", urlDomain, groupId)),
			},
		}
	}
	return result, nil
}

func resourceGroupDomainAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		groupId    = d.Get("group_id").(string)
		urlDomain  = d.Get("url_domain").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	resp, err := GetGroupAssociatedDomainByUrl(client, instanceId, groupId, urlDomain)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying domain associate info for the API group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("min_ssl_version", utils.PathSearch("min_ssl_version", resp, nil)),
		d.Set("ingress_http_port", utils.PathSearch("ingress_http_port", resp, nil)),
		d.Set("ingress_https_port", utils.PathSearch("ingress_https_port", resp, nil)),
		d.Set("url_domain", utils.PathSearch("domain", resp, nil)),
		// Attributes
		d.Set("domain_id", utils.PathSearch("id", resp, nil)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving associate resource for the domain and API group: %s", err)
	}
	return nil
}

func buildDomainAssociateUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"min_ssl_version":           utils.ValueIgnoreEmpty(d.Get("min_ssl_version")),
		"is_http_redirect_to_https": d.Get("is_http_redirect_to_https"),
		"ingress_http_port":         utils.ValueIgnoreEmpty(d.Get("ingress_http_port")),
		"ingress_https_port":        utils.ValueIgnoreEmpty(d.Get("ingress_https_port")),
	}
}

func resourceGroupDomainAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}/domains/{domain_id}"
		instanceId = d.Get("instance_id").(string)
		groupId    = d.Get("group_id").(string)
		domainId   = d.Get("domain_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	// Get domain ID if the value of attribute 'domain_id' is empty.
	if domainId == "" {
		urlDomain := d.Get("url_domain").(string)
		associateInfo, err := GetGroupAssociatedDomainByUrl(client, instanceId, groupId, urlDomain)
		if err != nil {
			return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying the associated domain information (%s)", urlDomain))
		}
		domainId = utils.PathSearch("id", associateInfo, "").(string)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)
	updatePath = strings.ReplaceAll(updatePath, "{group_id}", groupId)
	updatePath = strings.ReplaceAll(updatePath, "{domain_id}", domainId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildDomainAssociateUpdateBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &opt)
	if err != nil {
		return diag.Errorf("error updating associate configuration for the domain (%s) and the API group (%s): %s",
			domainId, groupId, err)
	}
	return resourceGroupDomainAssociateRead(ctx, d, meta)
}

func resourceGroupDomainAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}/domains/{domain_id}"
		instanceId = d.Get("instance_id").(string)
		groupId    = d.Get("group_id").(string)
		domainId   = d.Get("domain_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	// Get domain ID if the value of attribute 'domain_id' is empty.
	if domainId == "" {
		urlDomain := d.Get("url_domain").(string)
		associateInfo, err := GetGroupAssociatedDomainByUrl(client, instanceId, groupId, urlDomain)
		if err != nil {
			return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying the associated domain information (%s)", urlDomain))
		}
		domainId = utils.PathSearch("id", associateInfo, "").(string)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{group_id}", groupId)
	deletePath = strings.ReplaceAll(deletePath, "{domain_id}", domainId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error dissociate the domain name with the group (%s)", groupId))
	}
	return nil
}

func resourceGroupDomainAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid resource ID format, want '<instance_id>/<group_id>/<url_domain>', but got '%s'",
			importedId)
	}

	instanceId := parts[0]
	groupId := parts[1]
	urlDomain := parts[2]

	mrr := multierror.Append(nil,
		d.Set("instance_id", instanceId),
		d.Set("group_id", groupId),
		d.Set("url_domain", urlDomain),
	)
	return []*schema.ResourceData{d}, mrr.ErrorOrNil()
}
