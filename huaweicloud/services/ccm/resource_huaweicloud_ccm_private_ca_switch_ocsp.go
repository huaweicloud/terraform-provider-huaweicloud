package ccm

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCM POST /v1/private-certificate-authorities/{ca_id}/ocsp/switch
func ResourcePrivateCaSwitchOcsp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateCaSwitchOcspCreate,
		UpdateContext: resourcePrivateCaSwitchOcspUpdate,
		ReadContext:   resourcePrivateCaSwitchOcspRead,
		DeleteContext: resourcePrivateCaSwitchOcspDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"ca_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"ca_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ocsp_switch": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func switchPrivateCaOcsp(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/private-certificate-authorities/{ca_id}/ocsp/switch"
	requestPath = strings.ReplaceAll(requestPath, "{ca_id}", d.Get("ca_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
			202,
			204,
		},
		JSONBody: map[string]interface{}{
			"ocsp_switch": d.Get("ocsp_switch"),
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err == nil {
		return nil
	}

	var errCode403 golangsdk.ErrDefault403
	if errors.As(err, &errCode403) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode403.Body, &apiError); jsonErr != nil {
			log.Printf("[ERROR] unable to unmarshal switch private CA OCSP error response: %s", jsonErr)
			return err
		}

		// Error code "PCA.10010071" and "PCA.10010072" mean a duplicate operation status, which can be considered as
		// a successful operation status.
		errCode := utils.PathSearch("error_code", apiError, "").(string)
		if errCode == "PCA.10010072" || errCode == "PCA.10010071" {
			return nil
		}
	}

	return err
}

func resourcePrivateCaSwitchOcspCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ccm"
		caId    = d.Get("ca_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	if err := switchPrivateCaOcsp(client, d); err != nil {
		return diag.Errorf("error switching CCM private CA OCSP in create operation: %s", err)
	}

	d.SetId(caId)

	return resourcePrivateCaSwitchOcspRead(ctx, d, meta)
}

func resourcePrivateCaSwitchOcspRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePrivateCaSwitchOcspUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "ccm"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	if err := switchPrivateCaOcsp(client, d); err != nil {
		return diag.Errorf("error switching CCM private CA OCSP in update operation: %s", err)
	}

	return resourcePrivateCaSwitchOcspRead(ctx, d, meta)
}

func resourcePrivateCaSwitchOcspDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `Deleting this resource will not recover the private CA OCSP switch status, but will only remove the
	resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
