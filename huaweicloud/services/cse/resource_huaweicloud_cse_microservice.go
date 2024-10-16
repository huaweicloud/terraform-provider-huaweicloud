package cse

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v4/services"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var microserviceNotFoundCodes = []string{
	"400012",
}

// @API CSE DELETE /v2/{project_id}/registry/microservices/{serviceId}
// @API CSE GET /v2/{project_id}/registry/microservices/{serviceId}
// @API CSE POST /v2/{project_id}/registry/microservices
func ResourceMicroservice() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMicroserviceCreate,
		ReadContext:   resourceMicroserviceRead,
		DeleteContext: resourceMicroserviceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMicroserviceImportState,
		},

		Schema: map[string]*schema.Schema{
			// Authentication and request parameters.
			"auth_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					`The address that used to request the access token.`,
					utils.SchemaDescInput{
						Required: true,
					}),
			},
			"connect_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The address that used to send requests and manage configuration.`,
			},
			"admin_user": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The user name that used to pass the RBAC control.",
			},
			"admin_pass": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				RequiredWith: []string{"admin_user"},
				Description:  "The user password that used to pass the RBAC control.",
			},
			// Resource parameters.
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getAuthAddress(d *schema.ResourceData) string {
	if v, ok := d.GetOk("auth_address"); ok {
		return v.(string)
	}
	// Using the connect address as the auth address if its empty.
	// The behavior of the connect address is required.
	return d.Get("connect_address").(string)
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
	token, err := GetAuthorizationToken(getAuthAddress(d), d.Get("admin_user").(string),
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
	token, err := GetAuthorizationToken(getAuthAddress(d), d.Get("admin_user").(string),
		d.Get("admin_pass").(string))
	if err != nil {
		// When the engine does not exist, obtaining a token will cause a request connection exception.
		// To ensure that the resource is available on RFS platform, this situation is specially handled as a 404 error.
		log.Printf("[ERROR] %s", err)
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	client := common.NewCustomClient(true, d.Get("connect_address").(string), "v4", "default")
	resp, err := services.Get(client, d.Id(), token)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "errorCode", microserviceNotFoundCodes...),
			"CSE Microservice")
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
	token, err := GetAuthorizationToken(getAuthAddress(d), d.Get("admin_user").(string),
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
	var (
		authAddr, connectAddr, importedIdWithoutAddrs, microserviceId, adminUser, adminPwd string
		mErr                                                                               *multierror.Error

		importedId   = d.Id()
		addressRegex = `https://\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}`
		re           = regexp.MustCompile(fmt.Sprintf(`^(%[1]s)?/?(%[1]s)/(.*)$`, addressRegex))
		formatErr    = fmt.Errorf("the imported microservice ID specifies an invalid format, want "+
			"'<auth_address>/<connect_address>/<id>' or '<auth_address>/<connect_address>/<id>/<admin_user>/<admin_pass>', but got '%s'",
			importedId)
	)

	if !re.MatchString(importedId) {
		return nil, formatErr
	}
	resp := re.FindAllStringSubmatch(importedId, -1)
	// If the imported ID matches the address regular expression, the length of the response result must be greater than 1.
	switch len(resp[0]) {
	case 4:
		authAddr = resp[0][1]
		connectAddr = resp[0][2]
		importedIdWithoutAddrs = resp[0][3]
		if authAddr == "" {
			authAddr = connectAddr // Using the connect address as the auth address if the auth address input is omitted.
		}
	default:
		return nil, formatErr
	}

	mErr = multierror.Append(mErr,
		d.Set("auth_address", authAddr),
		d.Set("connect_address", connectAddr),
	)

	parts := strings.Split(importedIdWithoutAddrs, "/")
	switch len(parts) {
	case 1:
		microserviceId = parts[0]
	case 3:
		microserviceId = parts[0]
		adminUser = parts[1]
		adminPwd = parts[2]

		mErr = multierror.Append(mErr,
			d.Set("admin_user", adminUser),
			d.Set("admin_pass", adminPwd),
		)
	default:
		return nil, formatErr
	}

	d.SetId(microserviceId)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
