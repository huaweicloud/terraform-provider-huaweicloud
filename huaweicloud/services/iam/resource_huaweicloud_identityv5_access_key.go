package iam

import (
	"context"
	"fmt"
	"log"
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

var v5AccessKeyNonUpdatableParams = []string{"user_id"}

// @API IAM POST /v5/users/{user_id}/access-keys
// @API IAM GET /v5/users/{user_id}/access-keys
// @API IAM GET /v5/users/{user_id}/access-keys/{access_key_id}/last-used
// @API IAM PUT /v5/users/{user_id}/access-keys/{access_key_id}
// @API IAM DELETE /v5/users/{user_id}/access-keys/{access_key_id}
func ResourceV5AccessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5AccessKeyCreate,
		ReadContext:   resourceV5AccessKeyRead,
		UpdateContext: resourceV5AccessKeyUpdate,
		DeleteContext: resourceV5AccessKeyDelete,

		CustomizeDiff: config.FlexibleForceNew(v5AccessKeyNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceV5AccessKeyImportState,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the user.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The status of the access key.`,
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the generated access key.`,
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `The generated secret access key.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the access key.`,
			},
			"last_used_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the access key was last used.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceV5AccessKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId := d.Get("user_id").(string)
	httpUrl := "v5/users/{user_id}/access-keys"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{user_id}", userId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createAccessKeyResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IAM access key for user (%s): %s", userId, err)
	}

	resp, err := utils.FlattenResponse(createAccessKeyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	accessKeyId := utils.PathSearch("access_key.access_key_id", resp, "").(string)
	if accessKeyId == "" {
		return diag.Errorf("unable to find the IAM access key ID from the API response: %s", err)
	}

	d.SetId(accessKeyId)

	mErr := multierror.Append(
		d.Set("secret_access_key", utils.PathSearch("access_key.secret_access_key", resp, "").(string)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving secret access key field to state: %s", err)
	}

	if v, ok := d.GetOk("status"); ok && v.(string) != "active" {
		err = updateV5AccessKeyStatus(client, userId, accessKeyId, v.(string))
		if err != nil {
			return diag.Errorf("error updating status of access key for user (%s): %s", userId, err)
		}
	}

	return resourceV5AccessKeyRead(ctx, d, meta)
}

func updateV5AccessKeyStatus(client *golangsdk.ServiceClient, userId string, accessKeyId string, status string) error {
	httpUrl := "v5/users/{user_id}/access-keys/{access_key_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{user_id}", userId)
	updatePath = strings.ReplaceAll(updatePath, "{access_key_id}", accessKeyId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"status": status,
		},
	}
	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func listV5AccessKeys(client *golangsdk.ServiceClient, userId string) ([]interface{}, error) {
	var (
		httpUrl = "v5/users/{user_id}/access-keys"
		result  = make([]interface{}, 0)
		marker  = ""
		// The default limit is 200, maximum is 200.
		limit = 200
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{user_id}", userId)
	listPath += fmt.Sprintf("?limit=%d", limit)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &getOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		accessKeys := utils.PathSearch("access_keys", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, accessKeys...)
		if len(accessKeys) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func GetV5AccessKeyById(client *golangsdk.ServiceClient, userId string, accessKeyId string) (interface{}, error) {
	accessKeys, err := listV5AccessKeys(client, userId)
	if err != nil {
		return nil, err
	}

	accessKey := utils.PathSearch(fmt.Sprintf("[?access_key_id=='%s']|[0]", accessKeyId), accessKeys, nil)
	if accessKey == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v5/users/{user_id}/access-keys",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the IAM access key (%s) for user (%s) does not exist", accessKeyId, userId)),
			},
		}
	}

	return accessKey, nil
}

func resourceV5AccessKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId := d.Get("user_id").(string)
	accessKey, err := GetV5AccessKeyById(client, userId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving IAM access key of the user (%s)", userId))
	}

	mErr := multierror.Append(
		d.Set("access_key_id", utils.PathSearch("access_key_id", accessKey, nil)),
		d.Set("created_at", utils.PathSearch("created_at", accessKey, "").(string)),
		d.Set("status", utils.PathSearch("status", accessKey, "").(string)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM access key fields: %s", err)
	}

	accessKeyId := utils.PathSearch("access_key_id", accessKey, "").(string)
	if accessKeyId == "" {
		return nil
	}

	lastUsedAt, err := getV5AccessKeyLastUsedAt(client, userId, accessKeyId)
	if err != nil {
		log.Printf("[ERROR] error retrieving last used time of the access key (%s): %s", accessKeyId, err)
	}

	return diag.FromErr(d.Set("last_used_at", lastUsedAt))
}

func getV5AccessKeyLastUsedAt(client *golangsdk.ServiceClient, userId, accessKeyId string) (string, error) {
	path := client.Endpoint + "v5/users/" + userId + "/access-keys/" + accessKeyId + "/last-used"
	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	r, err := client.Request("GET", path, &reqOpt)
	if err != nil {
		return "", err
	}

	resp, err := utils.FlattenResponse(r)
	if err != nil {
		return "", err
	}

	return utils.PathSearch("access_key_last_used.last_used_at", resp, "").(string), nil
}

func resourceV5AccessKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		userId      = d.Get("user_id").(string)
		accessKeyId = d.Id()
	)

	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	err = updateV5AccessKeyStatus(client, userId, accessKeyId, d.Get("status").(string))
	if err != nil {
		return diag.Errorf("error updating status of access key for user (%s): %s", userId, err)
	}

	return resourceV5AccessKeyRead(ctx, d, meta)
}

func resourceV5AccessKeyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		userId = d.Get("user_id").(string)
	)

	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	httpUrl := "v5/users/{user_id}/access-keys/{access_key_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{user_id}", userId)
	deletePath = strings.ReplaceAll(deletePath, "{access_key_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting access key for user (%s): %s", userId, err)
	}

	return nil
}

func resourceV5AccessKeyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<user_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("user_id", parts[0])
}
