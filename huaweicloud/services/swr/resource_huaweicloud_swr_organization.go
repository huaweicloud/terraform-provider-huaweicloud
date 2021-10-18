package swr

import (
	"context"
	"fmt"
	"time"

	"github.com/chnsz/golangsdk/openstack/swr/v2/namespaces"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

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

		//request and response parameters
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
	config := meta.(*config.Config)
	swrClient, err := config.SwrV2Client(config.GetRegion(d))

	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloud SWR client : %s", err)
	}

	name := d.Get("name").(string)
	createOpts := namespaces.CreateOpts{
		Namespace: name,
	}

	err = namespaces.Create(swrClient, createOpts).ExtractErr()

	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud SWR Organization: %s", err)
	}

	d.SetId(name)

	return resourceSWROrganizationRead(ctx, d, meta)
}

func resourceSWROrganizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	swrClient, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud SWR client: %s", err)
	}

	n, err := namespaces.Get(swrClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving HuaweiCloud SWR")
	}

	permission := resourceSWRAuthToPermission(n.Auth)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", n.Name),
		d.Set("creator", n.CreatorName),
		d.Set("permission", permission),
	)

	login := fmt.Sprintf("swr.%s.%s", config.GetRegion(d), config.Cloud)
	mErr = multierror.Append(mErr, d.Set("login_server", login))

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceSWROrganizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	swrClient, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud SWR Client: %s", err)
	}

	err = namespaces.Delete(swrClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloud SWR Organization: %s", err)
	}

	d.SetId("")
	return nil
}
