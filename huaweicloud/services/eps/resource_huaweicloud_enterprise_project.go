package eps

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API EPS POST /v1.0/enterprise-projects
// @API EPS GET /v1.0/enterprise-projects/{enterprise_project_id}
// @API EPS PUT /v1.0/enterprise-projects/{enterprise_project_id}
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/action
func ResourceEnterpriseProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnterpriseProjectCreate,
		ReadContext:   resourceEnterpriseProjectRead,
		UpdateContext: resourceEnterpriseProjectUpdate,
		DeleteContext: resourceEnterpriseProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		// Request and response parameters.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"poc", "prod"}, false),
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"skip_disable_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeInt,
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

func resourceEnterpriseProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	epsClient, err := conf.EnterpriseProjectClient(conf.GetRegion(d))

	if err != nil {
		return diag.Errorf("unable to create EPS client: %s", err)
	}

	createOpts := enterpriseprojects.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
	}

	project, err := enterpriseprojects.Create(epsClient, createOpts).Extract()

	if err != nil {
		return diag.Errorf("error creating Enterprise Project: %s", err)
	}

	d.SetId(project.ID)

	return resourceEnterpriseProjectRead(ctx, d, meta)
}

func resourceEnterpriseProjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	epsClient, err := conf.EnterpriseProjectClient(conf.GetRegion(d))

	if err != nil {
		return diag.Errorf("unable to create EPS client: %s", err)
	}

	project, err := enterpriseprojects.Get(epsClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving HuaweiCloud Enterprise Project")
	}

	var enable bool
	if project.Status == 1 {
		enable = true
	}

	mErr := multierror.Append(nil,
		d.Set("name", project.Name),
		d.Set("description", project.Description),
		d.Set("type", project.Type),
		d.Set("status", project.Status),
		d.Set("enable", enable),
		d.Set("created_at", project.CreatedAt),
		d.Set("updated_at", project.UpdatedAt),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting enterprise project fields: %s", err)
	}

	return nil
}

func resourceEnterpriseProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	epsClient, err := conf.EnterpriseProjectClient(conf.GetRegion(d))

	if err != nil {
		return diag.Errorf("unable to create EPS client: %s", err)
	}

	if d.HasChange("enable") {
		var actionOpts enterpriseprojects.ActionOpts
		if d.Get("enable").(bool) {
			actionOpts.Action = "enable"
		} else {
			actionOpts.Action = "disable"
		}

		_, err = enterpriseprojects.Action(epsClient, actionOpts, d.Id()).Extract()
		if err != nil {
			return diag.Errorf("error enabling/disabling Enterprise Project: %s", err)
		}
	}

	if d.HasChanges("name", "description", "type") {
		updateOpts := enterpriseprojects.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Type:        d.Get("type").(string),
		}

		_, err = enterpriseprojects.Update(epsClient, updateOpts, d.Id()).Extract()

		if err != nil {
			return diag.Errorf("error updating Enterprise Project: %s", err)
		}
	}

	return resourceEnterpriseProjectRead(ctx, d, meta)
}

func resourceEnterpriseProjectDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	epsClient, err := conf.EnterpriseProjectClient(conf.GetRegion(d))

	if err != nil {
		return diag.Errorf("unable to create EPS client: %s", err)
	}

	if d.Get("skip_disable_on_destroy").(bool) {
		log.Printf("[DEBUG] skip disable on destroy for %s", d.Id())
		return nil
	}

	actionOpts := enterpriseprojects.ActionOpts{
		Action: "disable",
	}

	_, err = enterpriseprojects.Action(epsClient, actionOpts, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("error disabling Enterprise Project: %s", err)
	}

	d.SetId("")

	errorMsg := "Deleting enterprise projects is not supported. The project is only disabled and " +
		"removed from the state, but it remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
