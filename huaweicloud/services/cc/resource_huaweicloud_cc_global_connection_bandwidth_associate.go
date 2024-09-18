package cc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CC GET /v3/{domain_id}/gcb/gcbandwidths/{id}
// @API CC POST /v3/{domain_id}/gcb/gcbandwidths/{id}/disassociate-instance
// @API CC POST /v3/{domain_id}/gcb/gcbandwidths/{id}/associate-instance
func ResourceGlobalConnectionBandwidthAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlobalConnectionBandwidthAssociateCreate,
		UpdateContext: resourceGlobalConnectionBandwidthAssociateUpdate,
		ReadContext:   resourceGlobalConnectionBandwidthAssociateRead,
		DeleteContext: resourceGlobalConnectionBandwidthAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"gcb_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The global connection bandwidth ID.`,
			},
			"gcb_binding_resources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the resource to associate with the global connection bandwidth.`,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the resource to associate with the global connection bandwidth.`,
						},
						"region_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The region ID of the resource to associate with the global connection bandwidth.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The project ID of the resource to associate with the global connection bandwidth.`,
						},
					},
				},
				Description: `The resources to associate with the global connection bandwidth.`,
			},
		},
	}
}

func resourceGlobalConnectionBandwidthAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cc"
		gcbID   = d.Get("gcb_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	err = associateGlobalConnectionBandwidth(client, gcbID, cfg.DomainID, d.Get("gcb_binding_resources").(*schema.Set))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(gcbID)

	return resourceGlobalConnectionBandwidthAssociateRead(ctx, d, meta)
}

func associateGlobalConnectionBandwidth(client *golangsdk.ServiceClient, gcbID, domainID string, resources *schema.Set) error {
	httpUrl := "v3/{domain_id}/gcb/gcbandwidths/{id}/associate-instance"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", domainID)
	path = strings.ReplaceAll(path, "{id}", gcbID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	opt.JSONBody = utils.RemoveNil(buildGlobalConnectionBandwidthAssociateBodyParams(resources))
	resp, err := client.Request("POST", path, &opt)
	if err != nil {
		return fmt.Errorf("error associating the resource instance to the global connection bandwidth: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	results := utils.PathSearch("gcbandwidths[?result!='success']", respBody, make([]interface{}, 0)).([]interface{})
	if len(results) > 0 {
		messages := utils.PathSearch("[*].message", results, make([]interface{}, 0)).([]interface{})
		errMessage := "unable to associate the resource instance to the global connection bandwidth:"
		for _, message := range messages {
			errMessage += fmt.Sprintf(" %s", message)
		}
		return fmt.Errorf(errMessage)
	}

	return nil
}

func buildGlobalConnectionBandwidthAssociateBodyParams(resources *schema.Set) map[string]interface{} {
	gcbBindingResources := make([]map[string]interface{}, 0, resources.Len())
	for _, r := range resources.List() {
		resourceMap := r.(map[string]interface{})
		gcbBindingResources = append(gcbBindingResources, map[string]interface{}{
			"resource_id":   resourceMap["resource_id"],
			"resource_type": resourceMap["resource_type"],
			"region_id":     utils.ValueIgnoreEmpty(resourceMap["region_id"]),
			"project_id":    utils.ValueIgnoreEmpty(resourceMap["project_id"]),
		})
	}

	return map[string]interface{}{
		"gcbandwidths": gcbBindingResources,
	}
}

func resourceGlobalConnectionBandwidthAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{domain_id}/gcb/gcbandwidths/{id}"
		product = "cc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", cfg.DomainID)
	path = strings.ReplaceAll(path, "{id}", d.Id())

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", path, &opt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the global connection bandwidth")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	resources := utils.PathSearch("globalconnection_bandwidth.instances", respBody, make([]interface{}, 0))
	if v, ok := resources.([]interface{}); ok && len(v) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("gcb_id", utils.PathSearch("globalconnection_bandwidth.id", respBody, nil)),
		d.Set("gcb_binding_resources", flattenGlobalConnectionBandwidthAssociate(resources)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGlobalConnectionBandwidthAssociate(resources interface{}) []interface{} {
	res, ok := resources.([]interface{})
	if !ok || len(res) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(res))
	for _, raw := range res {
		rawMap := raw.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"resource_id":   utils.PathSearch("id", rawMap, nil),
			"resource_type": utils.PathSearch("type", rawMap, nil),
			"region_id":     utils.PathSearch("region_id", rawMap, nil),
			"project_id":    utils.PathSearch("project_id", rawMap, nil),
		})
	}

	return result
}

func resourceGlobalConnectionBandwidthAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	gcbID := d.Id()

	client, err := cfg.NewServiceClient("cc", region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	oldResources, newResources := d.GetChange("gcb_binding_resources")
	associateResources := newResources.(*schema.Set).Difference(oldResources.(*schema.Set))
	disassociateResources := oldResources.(*schema.Set).Difference(newResources.(*schema.Set))

	if disassociateResources.Len() > 0 {
		err := disassociateGlobalConnectionBandwidth(client, d, cfg.DomainID, disassociateResources)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if associateResources.Len() > 0 {
		err := associateGlobalConnectionBandwidth(client, gcbID, cfg.DomainID, associateResources)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGlobalConnectionBandwidthAssociateRead(ctx, d, meta)
}

func disassociateGlobalConnectionBandwidth(client *golangsdk.ServiceClient, d *schema.ResourceData, domainID string, resources *schema.Set) error {
	httpUrl := "v3/{domain_id}/gcb/gcbandwidths/{id}/disassociate-instance"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", domainID)
	path = strings.ReplaceAll(path, "{id}", d.Id())

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	opt.JSONBody = utils.RemoveNil(buildGlobalConnectionBandwidthAssociateBodyParams(resources))
	resp, err := client.Request("POST", path, &opt)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	results := utils.PathSearch("gcbandwidths[?result!='success']", respBody, make([]interface{}, 0)).([]interface{})
	if len(results) > 0 {
		messages := utils.PathSearch("[*].message", results, make([]interface{}, 0)).([]interface{})
		errMessage := "unable to disassociate the resource instance from the global connection bandwidth:"
		for _, message := range messages {
			errMessage += fmt.Sprintf(" %s", message)
		}
		return fmt.Errorf(errMessage)
	}

	return nil
}

func resourceGlobalConnectionBandwidthAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("cc", region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	disassociateResources, _ := d.GetChange("gcb_binding_resources")
	err = disassociateGlobalConnectionBandwidth(client, d, cfg.DomainID, disassociateResources.(*schema.Set))
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "GCB.0001"),
			"error disassociating the resource instance from the global connection bandwidth")
	}
	return nil
}
