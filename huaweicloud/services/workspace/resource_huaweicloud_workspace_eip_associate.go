package workspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/workspace/v2/desktops"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API Workspace POST /v2/{project_id}/eips/binding
// @API Workspace POST /v2/{project_id}/eips/unbinding
// @API Workspace GET /v2/{project_id}/eips
func ResourceEipAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEipAssociateCreate,
		ReadContext:   resourceEipAssociateRead,
		DeleteContext: resourceEipAssociateDelete,

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
			"desktop_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEipAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	desktopId := d.Get("desktop_id").(string)
	eipId := d.Get("eip_id").(string)
	createOpts := desktops.BindEipOpts{
		DesktopId: desktopId,
		ID:        eipId,
	}
	err = desktops.BindEip(client, createOpts)
	if err != nil {
		return diag.Errorf("error binding desktop EIP: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", desktopId, eipId))
	return resourceEipAssociateRead(ctx, d, meta)
}

func resourceEipAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.WorkspaceV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	desktopId := d.Get("desktop_id").(string)
	eips, err := desktops.ListEips(client, desktopId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Workspace EIP associate")
	}

	if len(eips) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "Workspace EIP associate")
	}

	associatedEipInfo := eips[0]
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("desktop_id", associatedEipInfo.AttachedDesktopId),
		d.Set("eip_id", associatedEipInfo.ID),
		d.Set("enterprise_project_id", associatedEipInfo.EnterpriseProjectId),
		d.Set("public_ip", associatedEipInfo.Address),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceEipAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	desktopId := d.Get("desktop_id").(string)
	opts := desktops.UnbindEipOpt{
		DesktopIds: []string{desktopId},
	}
	err = desktops.UnbindEip(client, opts)
	if err != nil {
		return diag.Errorf("error unbinding desktop EIP: %s", err)
	}

	return resourceEipAssociateRead(ctx, d, meta)
}
