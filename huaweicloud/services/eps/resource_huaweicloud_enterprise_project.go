package eps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EPS POST /v1.0/enterprise-projects
// @API EPS GET /v1.0/enterprise-projects/{enterprise_project_id}
// @API EPS PUT /v1.0/enterprise-projects/{enterprise_project_id}
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/action
// @API EPS DELETE /v1.0/enterprise-projects/{enterprise_project_id}
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
			"delete_flag": {
				Type:     schema.TypeBool,
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

	if !d.Get("enable").(bool) {
		if err := updateEnterpriseProjectEnable(epsClient, d); err != nil {
			return diag.Errorf("error disabling Enterprise Project in create: %s", err)
		}
	}

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

func updateEnterpriseProjectEnable(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var actionOpts enterpriseprojects.ActionOpts
	if d.Get("enable").(bool) {
		actionOpts.Action = "enable"
	} else {
		actionOpts.Action = "disable"
	}

	_, err := enterpriseprojects.Action(client, actionOpts, d.Id()).Extract()
	return err
}

func resourceEnterpriseProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	epsClient, err := conf.EnterpriseProjectClient(conf.GetRegion(d))

	if err != nil {
		return diag.Errorf("unable to create EPS client: %s", err)
	}

	if d.HasChange("enable") {
		if err := updateEnterpriseProjectEnable(epsClient, d); err != nil {
			return diag.Errorf("error enabling/disabling Enterprise Project in update: %s", err)
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

func resourceEnterpriseProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if d.Get("delete_flag").(bool) {
		return resourceEnterpriseProjectDeleteFlag(ctx, d, meta)
	}

	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "The project is only disabled and removed from the state, but it remains in the cloud.",
		},
	}
}

func handleProjectDeleteError(err error) (interface{}, string, error) {
	if err == nil {
		return "success", "COMPLETED", nil
	}

	var errCode404 golangsdk.ErrDefault404
	if errors.As(err, &errCode404) {
		return "deleted", "COMPLETED", nil
	}

	var errCode400 golangsdk.ErrDefault400
	if errors.As(err, &errCode400) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode400.Body, &apiError); jsonErr != nil {
			return nil, "ERROR", fmt.Errorf("error unmarshaling the deletion enterprise project response body: %s", jsonErr)
		}

		// EPS.0086: The enterprise project is checking resources. Please try again later.
		errCode := utils.PathSearch("error.error_code", apiError, "").(string)
		if errCode == "EPS.0086" {
			return "retry", "PENDING", nil
		}
	}

	return nil, "ERROR", err
}

func resourceEnterpriseProjectDeleteFlag(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("eps", cfg.Region)
	if err != nil {
		return diag.Errorf("error creating EPS client: %s", err)
	}
	httpUrl := "v1.0/enterprise-projects/{enterprise_project_id}"
	deletePath := client.Endpoint + httpUrl

	deletePath = strings.ReplaceAll(deletePath, "{enterprise_project_id}", d.Id())
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			_, err := client.Request("DELETE", deletePath, &opt)
			return handleProjectDeleteError(err)
		},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for Enterprise Project (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}
