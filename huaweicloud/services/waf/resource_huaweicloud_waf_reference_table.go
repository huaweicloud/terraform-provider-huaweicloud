package waf

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/valuelists"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF DELETE /v1/{project_id}/waf/valuelist/{valuelistid}
// @API WAF GET /v1/{project_id}/waf/valuelist/{valuelistid}
// @API WAF PUT /v1/{project_id}/waf/valuelist/{valuelistid}
// @API WAF POST /v1/{project_id}/waf/valuelist
func ResourceWafReferenceTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafReferenceTableCreate,
		ReadContext:   resourceWafReferenceTableRead,
		UpdateContext: resourceWafReferenceTableUpdate,
		DeleteContext: resourceWafReferenceTableDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"conditions": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 30,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "schema: Required",
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceWafReferenceTableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	opt := valuelists.CreateOpts{
		Name:                d.Get("name").(string),
		Type:                d.Get("type").(string),
		Values:              utils.ExpandToStringList(d.Get("conditions").([]interface{})),
		Description:         d.Get("description").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}

	r, err := valuelists.Create(client, opt)
	if err != nil {
		return diag.Errorf("error creating WAF reference table: %s", err)
	}
	d.SetId(r.Id)

	return resourceWafReferenceTableRead(ctx, d, meta)
}

func resourceWafReferenceTableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.WafV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	r, err := valuelists.GetWithEpsID(client, d.Id(), cfg.GetEnterpriseProjectID(d))
	if err != nil {
		// If the reference table does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF reference table")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", r.Name),
		d.Set("type", r.Type),
		d.Set("conditions", r.Values),
		d.Set("description", r.Description),
		d.Set("creation_time", time.Unix(r.CreationTime/1000, 0).Format("2006-01-02 15:04:05")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceWafReferenceTableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	desc := d.Get("description").(string)
	opt := valuelists.UpdateValueListOpts{
		Name: d.Get("name").(string),
		// Type is required, but it cannot be changed.
		Type:                d.Get("type").(string),
		Values:              utils.ExpandToStringList(d.Get("conditions").([]interface{})),
		Description:         &desc,
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}

	_, err = valuelists.Update(client, d.Id(), opt)
	if err != nil {
		return diag.Errorf("error updating WAF reference table: %s", err)
	}

	return resourceWafReferenceTableRead(ctx, d, meta)
}

func resourceWafReferenceTableDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	_, err = valuelists.DeleteWithEpsID(client, d.Id(), cfg.GetEnterpriseProjectID(d))
	if err != nil {
		// If the reference table does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF reference table")
	}

	return nil
}
