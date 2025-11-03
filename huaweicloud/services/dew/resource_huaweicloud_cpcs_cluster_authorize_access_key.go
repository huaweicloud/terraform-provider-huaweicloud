package dew

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1/{project_id}/dew/cpcs/cluster/{cluster_id}/authorize-access-keys
// @API DEW POST /v1/{project_id}/dew/cpcs/cluster/{cluster_id}/de-authorize-access-keys
// @API DEW GET /v1/{project_id}/dew/cpcs/cluster/{cluster_id}/access-keys
func ResourceClusterAuthorizeAccessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterAuthorizeAccessKeyCreate,
		ReadContext:   resourceClusterAuthorizeAccessKeyRead,
		UpdateContext: resourceClusterAuthorizeAccessKeyUpdate,
		DeleteContext: resourceClusterAuthorizeAccessKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceClusterAuthorizeAccessKeyImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"cluster_id",
			"app_id",
			"access_key_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			// no response value
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the cluster.`,
			},
			// no response value
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the application.`,
			},
			// The official documentation describes this field as allowing "all" to be entered, signifying authorization
			// for all applications. However, in actual testing, entering "all" results in an error,
			// so this value is not mentioned in the provider documentation.
			"access_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the access key.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the access key.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the application.`,
			},
			"access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `The access key.`,
			},
			"key_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the access key.`,
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The creation time of the access key, UNIX timestamp in milliseconds.`,
			},
		},
	}
}

func resourceClusterAuthorizeAccessKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dew/cpcs/cluster/{cluster_id}/authorize-access-keys"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", d.Get("cluster_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"app_id":         d.Get("app_id").(string),
			"access_key_ids": []string{d.Get("access_key_id").(string)},
		},
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error authorizing access key for DEW CPCS cluster: %s", err)
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(generateId)

	return resourceClusterAuthorizeAccessKeyRead(ctx, d, meta)
}

func QueryCpcsClusterAuthorizeAccessKey(client *golangsdk.ServiceClient, clusterId, accessKeyId string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/dew/cpcs/cluster/{cluster_id}/access-keys"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", clusterId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var (
		pageNum       = 1
		allAccessKeys []interface{}
		expression    = fmt.Sprintf("[?access_key_id=='%s']|[0]", accessKeyId)
	)

	for {
		requestPathWithPageNum := requestPath + fmt.Sprintf("?page_num=%d", pageNum)
		resp, err := client.Request("GET", requestPathWithPageNum, &requestOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		accessKeys := utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{})
		if len(accessKeys) == 0 {
			break
		}

		allAccessKeys = append(allAccessKeys, accessKeys...)
		pageNum++
	}

	targetAccessKey := utils.PathSearch(expression, allAccessKeys, nil)
	if targetAccessKey == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return targetAccessKey, nil
}

func resourceClusterAuthorizeAccessKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "kms"
		clusterId = d.Get("cluster_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	accessKeyDetail, err := QueryCpcsClusterAuthorizeAccessKey(client, clusterId, d.Get("access_key_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DEW CPCS cluster authorize access key")
	}

	mErr := multierror.Append(
		d.Set("cluster_id", clusterId),
		d.Set("access_key_id", utils.PathSearch("access_key_id", accessKeyDetail, nil)),
		d.Set("status", utils.PathSearch("status", accessKeyDetail, nil)),
		d.Set("app_name", utils.PathSearch("app_name", accessKeyDetail, nil)),
		d.Set("access_key", utils.PathSearch("access_key", accessKeyDetail, nil)),
		d.Set("key_name", utils.PathSearch("key_name", accessKeyDetail, nil)),
		d.Set("create_time", utils.PathSearch("create_time", accessKeyDetail, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceClusterAuthorizeAccessKeyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterAuthorizeAccessKeyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dew/cpcs/cluster/{cluster_id}/de-authorize-access-keys"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", d.Get("cluster_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"app_id":         d.Get("app_id").(string),
			"access_key_ids": []string{d.Get("access_key_id").(string)},
		},
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting DEW CPCS cluster authorize access key: %s", err)
	}

	return nil
}

func resourceClusterAuthorizeAccessKeyImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <cluster_id>/<access_key_id>, but got %s", d.Id())
	}

	mErr := multierror.Append(
		d.Set("cluster_id", parts[0]),
		d.Set("access_key_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
