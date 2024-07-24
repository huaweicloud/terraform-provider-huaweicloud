package swr

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/swr/v2/namespaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API SWR GET /v2/manage/namespaces/{name}
// @API SWR DELETE /v2/manage/namespaces/{name}
// @API SWR POST /v2/manage/namespaces
func ResourceSWROrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSWROrganizationCreate,
		ReadContext:   resourceSWROrganizationRead,
		DeleteContext: resourceSWROrganizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		// Request and response parameters.
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login_server": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSWROrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	swrClient, err := cfg.SwrV2Client(cfg.GetRegion(d))

	if err != nil {
		return diag.Errorf("unable to create SWR client: %s", err)
	}

	name := d.Get("name").(string)
	createOpts := namespaces.CreateOpts{
		Namespace: name,
	}

	err = namespaces.Create(swrClient, createOpts).ExtractErr()

	if err != nil {
		return diag.Errorf("error creating SWR organization: %s", err)
	}

	d.SetId(name)

	return resourceSWROrganizationRead(ctx, d, meta)
}

func resourceSWROrganizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	swrClient, err := cfg.SwrV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	n, err := namespaces.Get(swrClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR")
	}

	permission := resourceSWRAuthToPermission(n.Auth)

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", n.Name),
		d.Set("creator", n.CreatorName),
		d.Set("permission", permission),
	)

	login := fmt.Sprintf("swr.%s.%s", cfg.GetRegion(d), cfg.Cloud)
	mErr = multierror.Append(mErr, d.Set("login_server", login))

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceSWROrganizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	swrClient, err := cfg.SwrV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	err = namespaces.Delete(swrClient, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR organization")
	}

	return nil
}
