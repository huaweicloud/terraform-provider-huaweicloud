package apig

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/environments"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceApigEnvironmentV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceEnvironmentResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the dedicated instance is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the environment belongs.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z]\w*$`),
						"Only letters, digits and underscores (_) are allowed, and must start with a letter."),
					validation.StringLenBetween(3, 64),
				),
				Description: "The environment name.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[^<>]*$`),
						"The angle brackets (< and >) are not allowed."),
					validation.StringLenBetween(0, 255),
				),
				Description: "The environment description.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the environment was created.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "Use 'created_at' instead",
				Description: `schema: Deprecated; The time when the environment was created.`,
			},
		},
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)

		opts = environments.EnvironmentOpts{
			Name:        d.Get("name").(string),
			Description: utils.String(d.Get("description").(string)),
		}
	)

	resp, err := environments.Create(client, instanceId, opts).Extract()
	if err != nil {
		return diag.Errorf("error creating dedicated environment: %s", err)
	}
	d.SetId(resp.Id)

	return resourceEnvironmentRead(ctx, d, meta)
}

// GetEnvironmentFormServer is a method to get dedicated environment details form server using IDs.
func GetEnvironmentFormServer(client *golangsdk.ServiceClient, instanceId,
	envId string) (*environments.Environment, error) {
	allPages, err := environments.List(client, instanceId, environments.ListOpts{}).AllPages()
	if err != nil {
		return nil, err
	}
	envs, err := environments.ExtractEnvironments(allPages)
	if err != nil {
		return nil, err
	}
	for _, v := range envs {
		if v.Id == envId {
			return &v, nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func resourceEnvironmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := GetEnvironmentFormServer(client, instanceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "dedicated environment")
	}
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("created_at", resp.CreateTime),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving dedicated environment fields: %s", err)
	}
	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId    = d.Get("instance_id").(string)
		environmentId = d.Id()
	)

	opt := environments.EnvironmentOpts{
		Name:        d.Get("name").(string), // Due to API restrictions, the name must be provided.
		Description: utils.String(d.Get("description").(string)),
	}
	_, err = environments.Update(client, instanceId, environmentId, opt).Extract()
	if err != nil {
		return diag.Errorf("error updating dedicated environment (%s): %s", environmentId, err)
	}

	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	err = environments.Delete(client, instanceId, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting dedicated environment from the instance (%s): %s", instanceId, err)
	}

	return nil
}

// The ID cannot find on console, so we need to import by environment name.
func resourceEnvironmentResourceImportState(_ context.Context, d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<name>")
	}
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = parts[0]
		name       = parts[1]

		opt = environments.ListOpts{
			Name: name,
		}
	)
	pages, err := environments.List(client, instanceId, opt).AllPages()
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error retrieving environment: %s", err)
	}
	resp, err := environments.ExtractEnvironments(pages)
	if len(resp) < 1 {
		return []*schema.ResourceData{d}, fmt.Errorf("unable to find the environment (%s) form server: %s", name, err)
	}

	d.SetId(resp[0].Id)
	return []*schema.ResourceData{d}, d.Set("instance_id", instanceId)
}
