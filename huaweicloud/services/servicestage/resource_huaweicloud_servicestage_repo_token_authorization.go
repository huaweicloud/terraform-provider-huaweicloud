package servicestage

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/servicestage/v1/repositories"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceRepoTokenAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRepoTokenAuthCreate,
		ReadContext:   resourceRepoAuthRead,
		DeleteContext: resourceRepoAuthDelete,

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[\w.-]*$`),
						"The name can only contain letters, digits, underscores (_), hyphens (-) and dots (.)."),
					validation.StringLenBetween(4, 63),
				),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"github", "gitlab", "gitee",
				}, false),
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"token": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceRepoTokenAuthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err error
	config := meta.(*config.Config)
	client, err := config.ServiceStageV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating ServiceStage v1 client: %s", err)
	}

	opt := repositories.PersonalAuthOpts{
		Name:  d.Get("name").(string),
		Host:  d.Get("host").(string),
		Token: d.Get("token").(string),
	}
	auth, err := repositories.CreatePersonalAuth(client, d.Get("type").(string), opt)
	if err != nil {
		return diag.Errorf("error creating the ServiceStage repository token authorization: %s", err)
	}

	d.SetId(auth.Name)

	return resourceRepoAuthRead(ctx, d, meta)
}

func getAuthorizationByName(c *golangsdk.ServiceClient, name string) (*repositories.Authorization, error) {
	resp, err := repositories.List(c)
	if err != nil {
		return nil, err
	}
	for _, auth := range resp {
		if auth.Name == name {
			return &auth, nil
		}
	}
	return nil, fmt.Errorf("unable to find the authorization (%s)", name)
}

func resourceRepoAuthRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ServiceStageV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("error creating ServiceStage v1 client: %s", err)
	}

	auth, err := getAuthorizationByName(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ServiceStage repository authorization")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", auth.Name),
		d.Set("type", auth.RepoType),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRepoAuthDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ServiceStageV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("error creating ServiceStage v1 client: %s", err)
	}

	err = repositories.Delete(client, d.Id())
	if err != nil {
		return fmtp.DiagErrorf("error deleting ServiceStage repository authorization (%s): %s", d.Id(), err)
	}
	return nil
}
