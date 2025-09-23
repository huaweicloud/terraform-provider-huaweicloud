package servicestage

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/servicestage/v1/repositories"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

// @API ServiceStage POST /v1/{project_id}/git/auths/{repo_type}/password
// @API ServiceStage GET /v1/{project_id}/git/auths
// @API ServiceStage DELETE /v1/{project_id}/git/auths/{name}
func ResourceRepoPwdAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRepoPwdAuthCreate,
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
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"devcloud", "bitbucket",
				}, false),
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceRepoPwdAuthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var err error
	config := meta.(*config.Config)
	client, err := config.ServiceStageV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating ServiceStage v1 client: %s", err)
	}

	s := fmt.Sprintf("%s:%s", d.Get("user_name").(string), d.Get("password").(string))
	opt := repositories.PwdAuthOpts{
		Name:  d.Get("name").(string),
		Token: base64.StdEncoding.EncodeToString([]byte(s)),
	}
	auth, err := repositories.CreatePwdAuth(client, d.Get("type").(string), opt)
	if err != nil {
		return diag.Errorf("error creating the ServiceStage password authorization: %s", err)
	}

	d.SetId(auth.Name)

	return resourceRepoAuthRead(ctx, d, meta)
}
