package cse

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v4/services"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
)

func ResourceMicroservice() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMicroserviceCreate,
		ReadContext:   resourceMicroserviceRead,
		DeleteContext: resourceMicroserviceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMicroserviceImportState,
		},

		Schema: map[string]*schema.Schema{
			"connect_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9]([\w-.]*[A-Za-z0-9])?$`),
						"The name must start and end with a letter or a digit, and can only contain letters, digits, "+
							"underscore (_), hyphens (-) and dots (.)."),
					validation.StringLenBetween(1, 128),
				),
			},
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"development", "testing", "acceptance", "production",
				}, false),
			},
			"level": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"FRONT", "MIDDLE", "BACK",
				}, false),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"admin_user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"admin_pass": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				RequiredWith: []string{"admin_user"},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildMicroserviceCreateOpts(d *schema.ResourceData) services.CreateOpts {
	env := d.Get("environment").(string)
	return services.CreateOpts{
		Services: services.Service{
			Name:        d.Get("name").(string),
			AppId:       d.Get("app_name").(string),
			Environment: &env,
			Version:     d.Get("version").(string),
			Level:       d.Get("level").(string),
			Description: d.Get("description").(string),
		},
	}
}

func resourceMicroserviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := GetAuthorizationToken(d.Get("connect_address").(string), d.Get("admin_user").(string),
		d.Get("admin_pass").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	client := common.NewCustomClient(true, d.Get("connect_address").(string), "v4", "default")
	createOpts := buildMicroserviceCreateOpts(d)
	log.Printf("[DEBUG] The createOpts of the Microservice is: %v", createOpts)
	resp, err := services.Create(client, createOpts, token)
	if err != nil {
		return diag.Errorf("error creating dedicated microservice: %s", err)
	}
	d.SetId(resp.ID)

	return resourceMicroserviceRead(ctx, d, meta)
}

func resourceMicroserviceRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	token, err := GetAuthorizationToken(d.Get("connect_address").(string), d.Get("admin_user").(string),
		d.Get("admin_pass").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	client := common.NewCustomClient(true, d.Get("connect_address").(string), "v4", "default")
	resp, err := services.Get(client, d.Id(), token)
	if err != nil {
		return diag.Errorf("error getting dedicated microservice (%s): %s", d.Id(), err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", resp.Name),
		d.Set("app_name", resp.AppId),
		d.Set("environment", resp.Environment),
		d.Set("version", resp.Version),
		d.Set("level", resp.Level),
		d.Set("description", resp.Description),
		d.Set("status", resp.Status),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceMicroserviceDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	token, err := GetAuthorizationToken(d.Get("connect_address").(string), d.Get("admin_user").(string),
		d.Get("admin_pass").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// The current configuration is force deletion that delete microservices and related configuration and binding
	// instances
	deleteOpts := services.DeleteOpts{
		Force: true,
	}
	client := common.NewCustomClient(true, d.Get("connect_address").(string), "v4", "default")
	err = services.Delete(client, deleteOpts, d.Id(), token)
	if err != nil {
		return diag.Errorf("error deleting dedicated microservice (%s): %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceMicroserviceImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	re := regexp.MustCompile(`^(https://\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5})/(.*)$`)
	if !re.MatchString(d.Id()) {
		return nil, fmt.Errorf("The imported microservice ID specifies an invalid format, must start with the " +
			"connection address of the service registry center for the dedicated CSE engine.")
	}

	var mErr *multierror.Error
	formatErr := fmt.Errorf("The imported microservice ID specifies an invalid format, must be " +
		"<cnnect_address>/<microservice_id> or <cnnect_address>/<microservice_id>/<admin_user>/<admin_pass>.")

	resp := re.FindAllStringSubmatch(d.Id(), -1)
	if len(resp) >= 1 && len(resp[0]) == 3 {
		mErr = multierror.Append(mErr, d.Set("connect_address", resp[0][1]))
		parts := strings.SplitN(resp[0][2], "/", 3)
		switch len(parts) {
		case 1:
			d.SetId(parts[0])
		case 3:
			d.SetId(parts[0])
			mErr = multierror.Append(mErr,
				d.Set("admin_user", parts[1]),
				d.Set("admin_pass", parts[2]),
			)
		default:
			return nil, formatErr
		}
		return []*schema.ResourceData{d}, mErr.ErrorOrNil()
	}

	return nil, formatErr
}
