package codeartspipeline

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

var publisherNonUpdatableParams = []string{
	"name", "en_name", "support_url", "description", "logo_url", "website", "source_url",
}

// @API CodeArtsPipeline POST /v1/{domain_id}/publisher/create
// @API CodeArtsPipeline POST /v1/{domain_id}/publisher/detail
// @API CodeArtsPipeline GET /v1/{domain_id}/publisher/query-all
// @API CodeArtsPipeline DELETE /v1/{domain_id}/publisher/delete
func ResourceCodeArtsPipelinePublisher() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelinePublisherCreate,
		ReadContext:   resourcePipelinePublisherRead,
		UpdateContext: resourcePipelinePublisherUpdate,
		DeleteContext: resourcePipelinePublisherDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(publisherNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the publisher name.`,
			},
			"en_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the publisher English name.`,
			},
			"support_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the support URL.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description.`,
			},
			"logo_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the logo URL.`,
			},
			"website": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the website URL.`,
			},
			"source_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the source URL.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"auth_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the authorization status.`,
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the user ID.`,
			},
			"last_update_user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the updater ID.`,
			},
			"last_update_user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the updater name.`,
			},
			"last_update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the update time.`,
			},
		},
	}
}

func resourcePipelinePublisherCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v1/{domain_id}/publisher/create"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{domain_id}", cfg.DomainID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePipelinePublisherBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline publisher: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody, ""); err != nil {
		return diag.Errorf("error creating CodeArts Pipeline publisher: %s", err)
	}

	publisher, err := getPipelinePublisherByList(client, d, cfg.DomainID)
	if err != nil {
		return diag.Errorf("error getting CodeArts Pipeline publishers: %s", err)
	}

	id := utils.PathSearch("publisher_unique_id", publisher, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the CodeArts Pipeline publisher ID from the API response")
	}

	d.SetId(id)

	return resourcePipelinePublisherRead(ctx, d, meta)
}

func buildCreatePipelinePublisherBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"en_name":     d.Get("en_name"),
		"support_url": d.Get("support_url"),
		"description": d.Get("description"),
		"logo_url":    d.Get("logo_url"),
		"website":     d.Get("website"),
		"source_url":  d.Get("source_url"),
	}

	return bodyParams
}

func getPipelinePublisherByList(client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string) (interface{}, error) {
	getHttpUrl := "v1/{domain_id}/publisher/query-all?limit=10&name={name}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", domainId)
	getPath = strings.ReplaceAll(getPath, "{name}", d.Get("name").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// id is not in return, and there is no unique value except id, use all params to filter
	var path []string
	searchKey := []string{"name", "en_name", "support_url", "description", "logo_url", "website", "source_url"}
	for _, k := range searchKey {
		path = append(path, fmt.Sprintf("%s=='%v'", k, d.Get(k)))
	}
	searchPath := fmt.Sprintf("data[?%s]|[0]", strings.Join(path, "&&"))

	offset := 0
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return nil, err
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, fmt.Errorf("error flatten response: %s", err)
		}

		publishers := utils.PathSearch("data", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(publishers) == 0 {
			return nil, golangsdk.ErrDefault404{}
		}

		publisher := utils.PathSearch(searchPath, getRespBody, nil)
		if publisher != nil {
			return publisher, nil
		}

		offset += 10
	}
}

func resourcePipelinePublisherRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	publisher, err := GetPipelinePublisher(client, cfg.DomainID, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DEVPIPE.30011001"),
			"error retrieving CodeArts Pipeline publisher")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", publisher, nil)),
		d.Set("en_name", utils.PathSearch("en_name", publisher, nil)),
		d.Set("support_url", utils.PathSearch("support_url", publisher, nil)),
		d.Set("description", utils.PathSearch("description", publisher, nil)),
		d.Set("logo_url", utils.PathSearch("logo_url", publisher, nil)),
		d.Set("website", utils.PathSearch("website", publisher, nil)),
		d.Set("source_url", utils.PathSearch("source_url", publisher, nil)),
		d.Set("auth_status", utils.PathSearch("auth_status", publisher, nil)),
		d.Set("user_id", utils.PathSearch("user_id", publisher, nil)),
		d.Set("last_update_user_id", utils.PathSearch("last_update_user_id", publisher, nil)),
		d.Set("last_update_user_name", utils.PathSearch("last_update_user_name", publisher, nil)),
		d.Set("last_update_time", utils.PathSearch("last_update_time", publisher, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetPipelinePublisher(client *golangsdk.ServiceClient, domainId, id string) (interface{}, error) {
	httpUrl := "v1/{domain_id}/publisher/detail"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", domainId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         []string{id},
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	searchPath := fmt.Sprintf(`"%s"`, id)
	publisher := utils.PathSearch(searchPath, getRespBody, nil)
	if publisher == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return publisher, nil
}

func resourcePipelinePublisherUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePipelinePublisherDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v1/{domain_id}/publisher/delete?publisher_unique_id={id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{domain_id}", cfg.DomainID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DEVPIPE.30011001"),
			"error deleting CodeArts Pipeline publisher")
	}

	return nil
}
