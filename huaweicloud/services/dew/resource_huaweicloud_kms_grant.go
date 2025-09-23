// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product KMS
// ---------------------------------------------------------------

package dew

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1.0/{project_id}/kms/create-grant
// @API DEW POST /v1.0/{project_id}/kms/list-grants
// @API DEW POST /v1.0/{project_id}/kms/revoke-grant
func ResourceKmsGrant() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsGrantCreate,
		ReadContext:   resourceKmsGrantRead,
		DeleteContext: resourceKmsGrantDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceKmsGrantImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Key ID.`,
			},
			"grantee_principal": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the authorized user or account.`,
			},
			"operations": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
				Description: `List of granted operations.
`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Grant name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "user",
				ForceNew:    true,
				Description: `Authorization type. The options are: **user**, **domain**. The default value is **user**.`,
				ValidateFunc: validation.StringInSlice([]string{
					"user", "domain",
				}, false),
			},
			"retiring_principal": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The ID of the retiring user.`,
			},
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the creator.`,
			},
		},
	}
}

func resourceKmsGrantCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createGrant: create a KMS Grant.
	var (
		createGrantHttpUrl = "v1.0/{project_id}/kms/create-grant"
		createGrantProduct = "kms"
	)
	createGrantClient, err := cfg.NewServiceClient(createGrantProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS Client: %s", err)
	}

	createGrantPath := createGrantClient.Endpoint + createGrantHttpUrl
	createGrantPath = strings.ReplaceAll(createGrantPath, "{project_id}", createGrantClient.ProjectID)

	createGrantOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createGrantOpt.JSONBody = utils.RemoveNil(buildCreateGrantBodyParams(d, cfg))
	createGrantResp, err := createGrantClient.Request("POST", createGrantPath, &createGrantOpt)
	if err != nil {
		return diag.Errorf("error creating KMS: %s", err)
	}

	createGrantRespBody, err := utils.FlattenResponse(createGrantResp)
	if err != nil {
		return diag.FromErr(err)
	}

	grantId := utils.PathSearch("grant_id", createGrantRespBody, "").(string)
	if grantId == "" {
		return diag.Errorf("unable to find the KMS grant ID from the API response")
	}
	d.SetId(grantId)

	return resourceKmsGrantRead(ctx, d, meta)
}

func buildCreateGrantBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                   utils.ValueIgnoreEmpty(d.Get("name")),
		"key_id":                 utils.ValueIgnoreEmpty(d.Get("key_id")),
		"grantee_principal_type": utils.ValueIgnoreEmpty(d.Get("type")),
		"grantee_principal":      utils.ValueIgnoreEmpty(d.Get("grantee_principal")),
		"operations":             utils.ValueIgnoreEmpty(d.Get("operations").(*schema.Set).List()),
		"retiring_principal":     utils.ValueIgnoreEmpty(d.Get("retiring_principal")),
	}
	return bodyParams
}

func resourceKmsGrantRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getGrant: Query the KMS manual Grant
	var (
		getGrantHttpUrl = "v1.0/{project_id}/kms/list-grants"
		getGrantProduct = "kms"
	)
	getGrantClient, err := cfg.NewServiceClient(getGrantProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS Client: %s", err)
	}

	getGrantPath := getGrantClient.Endpoint + getGrantHttpUrl
	getGrantPath = strings.ReplaceAll(getGrantPath, "{project_id}", getGrantClient.ProjectID)

	getGrantOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	allGrants := make([]interface{}, 0)
	var nextMarker string
	getGrantOpt.JSONBody = utils.RemoveNil(buildReadGrantBodyParams(d, cfg))
	getGrantJSONBody := getGrantOpt.JSONBody.(map[string]interface{})
	for {
		getGrantResp, err := getGrantClient.Request("POST", getGrantPath, &getGrantOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "KMS grant")
		}
		getGrantRespBody, err := utils.FlattenResponse(getGrantResp)
		if err != nil {
			return diag.FromErr(err)
		}
		grants := utils.PathSearch("grants", getGrantRespBody, make([]interface{}, 0)).([]interface{})
		if len(grants) > 0 {
			allGrants = append(allGrants, grants...)
		}
		nextMarker = utils.PathSearch("next_marker", getGrantRespBody, "").(string)
		if nextMarker == "" {
			break
		}
		getGrantJSONBody["marker"] = nextMarker
	}

	searchPath := fmt.Sprintf("[?grant_id=='%s']|[0]", d.Id())
	grantDetail := utils.PathSearch(searchPath, allGrants, nil)
	if grantDetail == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "KMS grant")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", grantDetail, nil)),
		d.Set("key_id", utils.PathSearch("key_id", grantDetail, nil)),
		d.Set("type", utils.PathSearch("grantee_principal_type", grantDetail, nil)),
		d.Set("grantee_principal", utils.PathSearch("grantee_principal", grantDetail, nil)),
		d.Set("operations", utils.PathSearch("operations", grantDetail, nil)),
		d.Set("creator", utils.PathSearch("issuing_principal", grantDetail, nil)),
		d.Set("retiring_principal", utils.PathSearch("retiring_principal", grantDetail, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildReadGrantBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id": d.Get("key_id"),
		"limit":  100,
	}
	return bodyParams
}

func resourceKmsGrantDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteGrant: missing operation notes
	var (
		deleteGrantHttpUrl = "v1.0/{project_id}/kms/revoke-grant"
		deleteGrantProduct = "kms"
	)
	deleteGrantClient, err := cfg.NewServiceClient(deleteGrantProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS Client: %s", err)
	}

	deleteGrantPath := deleteGrantClient.Endpoint + deleteGrantHttpUrl
	deleteGrantPath = strings.ReplaceAll(deleteGrantPath, "{project_id}", deleteGrantClient.ProjectID)

	deleteGrantOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	deleteGrantOpt.JSONBody = utils.RemoveNil(buildDeleteGrantBodyParams(d, cfg))
	_, err = deleteGrantClient.Request("POST", deleteGrantPath, &deleteGrantOpt)
	if err != nil {
		return diag.Errorf("error deleting KMS grant: %s", err)
	}

	return nil
}

func buildDeleteGrantBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id":   d.Get("key_id"),
		"grant_id": d.Id(),
	}
	return bodyParams
}

func resourceKmsGrantImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <key_id>/<grant_id>")
	}

	d.Set("key_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
