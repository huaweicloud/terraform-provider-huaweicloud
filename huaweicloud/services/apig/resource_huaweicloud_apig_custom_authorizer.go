package apig

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/authorizers"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

const (
	frontAuth   = "FRONTEND"
	backendAuth = "BACKEND"
)

func ResourceApigCustomAuthorizerV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigCustomAuthorizerV2Create,
		Read:   resourceApigCustomAuthorizerV2Read,
		Update: resourceApigCustomAuthorizerV2Update,
		Delete: resourceApigCustomAuthorizerV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceApigCustomAuthorizerResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[A-Za-z][\\w]{2,63}$"),
					"The name consists of 3 to 64 characters, starting with a letter. "+
						"Only letters, digits and underscores (_) are allowed."),
			},
			"function_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  frontAuth,
				ValidateFunc: validation.StringInSlice([]string{
					frontAuth, backendAuth,
				}, false),
			},
			"is_body_send": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cache_age": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 3600),
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The parameter identity only required if type is 'FRONTEND'.
			"identity": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"location": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HEADER", "QUERY",
							}, false),
						},
						"validation": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 2048),
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildApigIdentities(identities *schema.Set) []authorizers.AuthCreateIdentitiesReq {
	result := make([]authorizers.AuthCreateIdentitiesReq, identities.Len())
	for i, val := range identities.List() {
		identity := val.(map[string]interface{})
		validContent := identity["validation"].(string)
		result[i] = authorizers.AuthCreateIdentitiesReq{
			Name:       identity["name"].(string),
			Location:   identity["location"].(string),
			Validation: &validContent,
		}
	}
	return result
}

func buildApigCustomAuthorizerParameters(d *schema.ResourceData) (authorizers.CustomAuthOpts, error) {
	isBodySend := d.Get("is_body_send").(bool)
	userData := d.Get("user_data").(string)
	t := d.Get("type").(string) // The 'authType' is easily confused with 'AuthorizerType', and 'type' is a keyword.
	identities := d.Get("identity").(*schema.Set)
	if identities.Len() > 0 && t != frontAuth {
		return authorizers.CustomAuthOpts{}, fmt.Errorf("The identities can only be set when the type is 'FRONTEND'")
	}
	return authorizers.CustomAuthOpts{
		Name:           d.Get("name").(string),
		Type:           t,
		AuthorizerType: "FUNC", // The custom authorizer only support 'FUNC'.
		AuthorizerURI:  d.Get("function_urn").(string),
		IsBodySend:     &isBodySend,
		TTL:            golangsdk.IntToPointer(d.Get("cache_age").(int)),
		UserData:       &userData,
		Identities:     buildApigIdentities(identities),
	}, nil
}

func resourceApigCustomAuthorizerV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opt, err := buildApigCustomAuthorizerParameters(d)
	if err != nil {
		return err
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := authorizers.Create(client, instanceId, opt).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "Error creating HuaweiCloud APIG custom authorizer")
	}
	d.SetId(resp.ID)
	return resourceApigCustomAuthorizerV2Read(d, meta)
}

func setApigCustomAuthorizerIdentities(d *schema.ResourceData,
	identities []authorizers.Identity) error {
	result := make([]map[string]interface{}, len(identities))
	for i, val := range identities {
		result[i] = map[string]interface{}{
			"name":       val.Name,
			"location":   val.Location,
			"validation": val.Validation,
		}
	}
	return d.Set("identity", result)
}

func setApigCustomAuthorizerParamters(d *schema.ResourceData, config *config.Config,
	resp *authorizers.CustomAuthorizer) error {
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("function_urn", resp.AuthorizerURI),
		d.Set("type", resp.Type),
		d.Set("is_body_send", resp.IsBodySend),
		d.Set("cache_age", resp.TTL),
		d.Set("user_data", resp.UserData),
		d.Set("create_time", resp.CreateTime),
		setApigCustomAuthorizerIdentities(d, resp.Identities),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func resourceApigCustomAuthorizerV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := authorizers.Get(client, instanceId, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, fmt.Sprintf("Unable to get the custom authorizer (%s) form server", d.Id()))
	}
	return setApigCustomAuthorizerParamters(d, config, resp)
}

func resourceApigCustomAuthorizerV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	opt, err := buildApigCustomAuthorizerParameters(d)
	if err != nil {
		return err
	}
	instanceId := d.Get("instance_id").(string)
	_, err = authorizers.Update(client, instanceId, d.Id(), opt).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating HuaweiCloud APIG custom authorizer (%s): %s", d.Id(), err)
	}

	return resourceApigCustomAuthorizerV2Read(d, meta)
}

func resourceApigCustomAuthorizerV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	err = authorizers.Delete(client, instanceId, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud APIG custom authorizer from the instance (%s): %s",
			instanceId, err)
	}
	d.SetId("")
	return nil
}

// The ID cannot find on the console, so we need to import by authorizer name.
func resourceApigCustomAuthorizerResourceImportState(d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmtp.Errorf("Invalid format specified for import id, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	name := parts[1]
	opt := authorizers.ListOpts{
		Name: name,
	}
	pages, err := authorizers.List(client, instanceId, opt).AllPages()
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error retrieving custom authorizer: %s", err)
	}
	resp, err := authorizers.ExtractCustomAuthorizers(pages)
	if len(resp) < 1 {
		return []*schema.ResourceData{d}, fmtp.Errorf("Unable to find the custom authorizer (%s) form server: %s",
			name, err)
	}
	d.SetId(resp[0].ID)
	d.Set("instance_id", instanceId)
	return []*schema.ResourceData{d}, nil
}
