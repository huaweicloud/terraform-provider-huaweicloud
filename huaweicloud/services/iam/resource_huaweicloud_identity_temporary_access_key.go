package iam

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourceIdentitySecurityTokens
// @API IAM POST /v3.0/OS-CREDENTIAL/securitytokens
func ResourceIdentityTemporaryAccessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityTemporaryCreate,
		ReadContext:   resourceIdentityTemporaryRead,
		DeleteContext: resourceIdentityTemporaryDelete,

		Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"methods": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Authentication method, the content of this field is either `token` or `assume_role`.",
			},
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"duration_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  900,
				ForceNew: true,
			},
			"session_user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"securitytoken": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityTemporaryCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient := common.NewCustomClient(true, "https://iam.{region_id}.myhuaweicloud.com")
	basePath := iamClient.ResourceBase + "v3.0/OS-CREDENTIAL/securitytokens"
	basePath = strings.ReplaceAll(basePath, "{region_id}", cfg.GetRegion(d))
	createTokenOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"X-Auth-Token": d.Get("token").(string),
		},
		JSONBody: map[string]interface{}{
			"auth": buildCreateSecurityTokensBodyParams(d),
		},
	}
	response, err := iamClient.Request("POST", basePath, &createTokenOpt)
	if err != nil {
		return diag.Errorf("error getting identity security token: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.Errorf("error flattening response: %s", err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)
	mErr := multierror.Append(
		d.Set("expires_at", utils.PathSearch("credential.expires_at", respBody, "").(string)),
		d.Set("access", utils.PathSearch("credential.access", respBody, "").(string)),
		d.Set("secret", utils.PathSearch("credential.secret", respBody, "").(string)),
		d.Set("securitytoken", utils.PathSearch("credential.securitytoken", respBody, "").(string)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCreateSecurityTokensBodyParams(d *schema.ResourceData) map[string]interface{} {
	methods := d.Get("methods").(string)
	bodyParams := map[string]interface{}{}
	if methods == "token" {
		bodyParams = map[string]interface{}{
			"methods": []string{methods},
			"token": map[string]interface{}{
				"duration_seconds": d.Get("duration_seconds").(int),
			},
		}
	} else if methods == "assume_role" {
		bodyParams = map[string]interface{}{
			"methods": []string{methods},
			"assume_role": map[string]interface{}{
				"agency_name":      d.Get("agency_name"),
				"domain_name":      d.Get("domain_name"),
				"domain_id":        d.Get("domain_id"),
				"duration_seconds": d.Get("duration_seconds").(int),
				"session_user": map[string]interface{}{
					"name": d.Get("session_user_name"),
				},
			},
		}
	}
	if policy, ok := d.GetOk("policy"); ok {
		bodyParams["policy"] = policy
	}
	bodyParams = map[string]interface{}{"identity": bodyParams}
	return bodyParams
}

func resourceIdentityTemporaryRead(c context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	expiresAt, err := time.ParseInLocation(`2006-01-02T15:04:05Z`, d.Get("expires_at").(string), time.UTC)
	if err != nil {
		diag.Errorf("error parsing expires at: %s", err)
	}
	if time.Now().After(expiresAt) {
		return resourceIdentityTemporaryCreate(c, d, meta)
	}
	return nil
}

func resourceIdentityTemporaryDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting token is not supported. The token is only removed from the state, but it remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
