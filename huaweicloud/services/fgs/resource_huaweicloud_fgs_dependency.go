package fgs

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/dependencies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API FunctionGraph POST /v2/{project_id}/fgs/dependencies
// @API FunctionGraph GET /v2/{project_id}/fgs/dependencies/{depend_id}
// @API FunctionGraph PUT /v2/{project_id}/fgs/dependencies/{depend_id}
// @API FunctionGraph DELETE /v2/{project_id}/fgs/dependencies/{depend_id}
func ResourceFgsDependency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFgsDependencyCreate,
		ReadContext:   resourceFgsDependencyRead,
		UpdateContext: resourceFgsDependencyUpdate,
		DeleteContext: resourceFgsDependencyDelete,
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
			"runtime": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"link": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildFgsDependencyOpts(d *schema.ResourceData) dependencies.DependOpts {
	desc := d.Get("description").(string)
	// Since the zip file upload is limited in size and requires encoding, only the OBS type is supported.
	// The zip file uploading can also be achieved by uploading OBS objects and is more secure.
	return dependencies.DependOpts{
		Name:        d.Get("name").(string),
		Runtime:     d.Get("runtime").(string),
		Description: &desc,
		Type:        "obs",
		Link:        d.Get("link").(string),
	}
}

func resourceFgsDependencyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	opts := buildFgsDependencyOpts(d)
	resp, err := dependencies.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating custom dependency: %s", err)
	}
	d.SetId(resp.ID)

	return resourceFgsDependencyRead(ctx, d, meta)
}

func resourceFgsDependencyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	resp, err := dependencies.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "FunctionGraph dependency")
	}

	log.Printf("[DEBUG] Retrieved custom dependency %s: %+v", d.Id(), resp)
	mErr := multierror.Append(
		d.Set("runtime", resp.Runtime),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("link", resp.Link),
		d.Set("etag", resp.Etag),
		d.Set("size", resp.Size),
		d.Set("owner", resp.Owner),
		d.Set("version", resp.Version),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting resource fields of custom dependency (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceFgsDependencyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	opts := buildFgsDependencyOpts(d)
	_, err = dependencies.Update(client, d.Id(), opts)
	if err != nil {
		return diag.Errorf("error updating custom dependency: %s", err)
	}

	return resourceFgsDependencyRead(ctx, d, meta)
}

func resourceFgsDependencyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	fgsClient, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	err = dependencies.Delete(fgsClient, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting custom dependency")
	}
	return nil
}
