package dew

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/csms/v1/secrets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1/{project_id}/secrets
// @API DEW POST /v1/{project_id}/{resourceType}/{id}/tags/action
// @API DEW GET /v1/{project_id}/secrets/{secret_name}
// @API DEW GET /v1/{project_id}/secrets/{secret_name}/versions
// @API DEW GET /v1/{project_id}/secrets/{secret_name}/versions/{version_id}
// @API DEW PUT /v1/{project_id}/secrets/{secret_name}/versions/{version_id}
// @API DEW GET /v1/{project_id}/{resourceType}/{id}/tags
// @API DEW PUT /v1/{project_id}/secrets/{secret_name}
// @API DEW POST /v1/{project_id}/secrets/{secret_name}/versions
// @API DEW DELETE /v1/{project_id}/secrets/{secret_name}
func ResourceSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecretCreate,
		ReadContext:   resourceSecretRead,
		UpdateContext: resourceSecretUpdate,
		DeleteContext: resourceSecretDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
			"secret_text": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				StateFunc:    utils.HashAndHexEncode,
				ExactlyOneOf: []string{"secret_text", "secret_binary"},
			},
			"secret_binary": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				StateFunc: utils.HashAndHexEncode,
			},
			"expire_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"event_subscriptions": {
				// the field can be left blank, no need add Computed attribute
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": common.TagsSchema(),
			"secret_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_stages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		name   = d.Get("name").(string)
	)

	client, err := cfg.KmsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	createOpts := secrets.CreateSecretOpts{
		Name:                name,
		KmsKeyID:            d.Get("kms_key_id").(string),
		Description:         d.Get("description").(string),
		SecretString:        d.Get("secret_text").(string),
		SecretBinary:        d.Get("secret_binary").(string),
		SecretType:          d.Get("secret_type").(string),
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
		EventSubscriptions:  utils.ExpandToStringListBySet(d.Get("event_subscriptions").(*schema.Set)),
	}

	rst, err := secrets.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating CSMS secret: %s", err)
	}
	log.Printf("[DEBUG] The response body information for creating CSMS secret: %#v", rst)

	if rst.ID == "" {
		return diag.Errorf("error creating CSMS secret: ID is not found in API response")
	}
	d.SetId(fmt.Sprintf("%s/%s", rst.ID, name))

	if err := utils.CreateResourceTags(client, d, "csms", rst.ID); err != nil {
		return diag.Errorf("error setting tags of CSMS secret (%s): %s", d.Id(), err)
	}

	if _, ok := d.GetOk("expire_time"); ok && rst.State == "ENABLED" {
		if err := updateSecretVersion(client, d, name); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSecretRead(ctx, d, meta)
}

func parseSecretResourceID(id string) (secretID, secretName string, err error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid format for CSMS secret resource ID, want '<secret_id>/<name>', but got '%s'", id)
		return
	}
	secretID = parts[0]
	secretName = parts[1]
	return
}

// Due to API reasons, the response array needs to be filtered by empty strings.
func removeNullValues(s []interface{}) []interface{} {
	result := make([]interface{}, 0, len(s))
	for _, elem := range s {
		if v, ok := elem.(string); ok && v != "" {
			result = append(result, v)
		}
	}
	return result
}

func flattenSecretText(version *secrets.Version) string {
	if version.SecretString == "" {
		return ""
	}
	return utils.HashAndHexEncode(version.SecretString)
}

func flattenSecretBinary(version *secrets.Version) string {
	if version.SecretBinary == "" {
		return ""
	}
	return utils.HashAndHexEncode(version.SecretBinary)
}

func resourceSecretRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.KmsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	id, name, err := parseSecretResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	secret, err := secrets.Get(client, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CSMS secret")
	}
	log.Printf("[DEBUG] The response body information for getting CSMS secret: %#v", secret)

	version, err := queryLatestVersion(cfg, region, name)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("secret_id", secret.ID),
		d.Set("name", secret.Name),
		d.Set("kms_key_id", secret.KmsKeyID),
		d.Set("description", secret.Description),
		d.Set("status", secret.State),
		d.Set("create_time", utils.FormatTimeStampRFC3339(int64(secret.CreateTime)/1000, true, "2006-01-02 15:04:05 MST")),
		d.Set("secret_type", secret.SecretType),
		d.Set("enterprise_project_id", secret.EnterpriseProjectID),
		d.Set("event_subscriptions", removeNullValues(secret.EventSubscriptions)),
		d.Set("secret_text", flattenSecretText(version)),
		d.Set("secret_binary", flattenSecretBinary(version)),
		d.Set("expire_time", version.VersionMetadata.ExpireTime),
		d.Set("latest_version", version.VersionMetadata.ID),
		d.Set("version_stages", version.VersionMetadata.VersionStages),
		utils.SetResourceTagsToState(d, client, "csms", id),
		d.Set("tags", d.Get("tags")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func queryLatestVersion(cfg *config.Config, region, name string) (*secrets.Version, error) {
	client, err := cfg.KmsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS client: %s", err)
	}

	versions, err := secrets.ListSecretVersions(client, name)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CSMS secret versions: %s", err)
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("error retrieving CSMS secret versions: The versions in API response is empty")
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].CreateTime > versions[j].CreateTime
	})

	versionID := versions[0].ID

	return queryVersion(cfg, region, name, versionID)
}

func queryVersion(cfg *config.Config, region, name, versionID string) (*secrets.Version, error) {
	client, err := cfg.KmsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS client: %s", err)
	}

	version, err := secrets.ShowSecretVersion(client, name, versionID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CSMS secret version (%s): %s", versionID, err)
	}
	return version, nil
}

func updateSecret(client *golangsdk.ServiceClient, d *schema.ResourceData, name string) error {
	opts := secrets.UpdateSecretOpts{
		KmsKeyID:           d.Get("kms_key_id").(string),
		Description:        utils.String(d.Get("description").(string)),
		EventSubscriptions: utils.ExpandToStringListBySet(d.Get("event_subscriptions").(*schema.Set)),
	}

	_, err := secrets.Update(client, name, opts)
	if err != nil {
		return fmt.Errorf("error updating CSMS secret (%s): %s", name, err)
	}
	return nil
}

// Credential values need to be updated by creating a new credential version.
func createSecretVersion(client *golangsdk.ServiceClient, d *schema.ResourceData, name string) error {
	opts := secrets.CreateVersionOpts{
		SecretString: d.Get("secret_text").(string),
		SecretBinary: d.Get("secret_binary").(string),
		ExpireTime:   d.Get("expire_time").(int),
	}

	_, err := secrets.CreateSecretVersion(client, name, opts)
	if err != nil {
		return fmt.Errorf("error creating a new CSMS secret (%s) version: %s", name, err)
	}
	return nil
}

// updateSecretVersion using to update the secret version expiration time.
func updateSecretVersion(client *golangsdk.ServiceClient, d *schema.ResourceData, name string) error {
	opts := secrets.UpdateVersionOpts{
		ExpireTime: d.Get("expire_time").(int),
	}

	versions, err := secrets.ListSecretVersions(client, name)
	if err != nil {
		return fmt.Errorf("error retrieving CSMS secret versions: %s", err)
	}

	if len(versions) == 0 {
		return fmt.Errorf("error retrieving CSMS secret versions: The versions in API response is empty")
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].CreateTime > versions[j].CreateTime
	})

	versionID := versions[0].ID
	_, err = secrets.UpdateSecretVersion(client, name, versionID, opts)
	if err != nil {
		return fmt.Errorf("error updating CSMS secret (%s) version: %s", name, err)
	}

	return nil
}

func resourceSecretUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.KmsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	id, name, err := parseSecretResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("kms_key_id", "description", "event_subscriptions") {
		if err := updateSecret(client, d, name); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("secret_text", "secret_binary") {
		if err := createSecretVersion(client, d, name); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("expire_time") && d.Get("status") == "ENABLED" {
		if err := updateSecretVersion(client, d, name); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		if err := utils.UpdateResourceTags(client, d, "csms", id); err != nil {
			return diag.Errorf("error updating tags of CSMS secret (%s): %s", id, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   id,
			ResourceType: "csms",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSecretRead(ctx, d, meta)
}

func resourceSecretDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		name   = d.Get("name").(string)
	)

	client, err := cfg.KmsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	if err := secrets.Delete(client, name); err != nil {
		return diag.Errorf("error deleting CSMS secret (%s): %s", name, err)
	}
	return nil
}
