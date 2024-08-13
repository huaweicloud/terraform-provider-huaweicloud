package dew

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/csms/v1/secrets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	serviceType = "csms"
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
func ResourceCsmsSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCsmsSecretCreate,
		ReadContext:   resourceCsmsSecretRead,
		UpdateContext: resourceCsmsSecretUpdate,
		DeleteContext: resourceCsmsSecretDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCsmsSecretImport,
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
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[\w-\.]{1,64}$`),
					"The maximum length is 64 characters. "+
						"Only letters, digits, underscores (_) hyphens (-) and dots (.) are allowed."),
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
				ForceNew: true,
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

func resourceCsmsSecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// The endpoint of CSMS is the endpoint of KMS.
	client, err := cfg.KmsV1Client(region)
	if err != nil {
		return diag.Errorf("failed to create CSMS(KMS) client: %s", err)
	}

	name := d.Get("name").(string)
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
		return diag.Errorf("failed to create the CSMS secret: %s", err)
	}

	id := fmt.Sprintf("%s/%s", rst.ID, name)
	d.SetId(id)

	// Save tags
	if t, ok := d.GetOk("tags"); ok {
		tMaps := t.(map[string]interface{})
		tagMaps := utils.ExpandResourceTags(tMaps)
		err = tags.Create(client, serviceType, rst.ID, tagMaps).ExtractErr()
		if err != nil {
			log.Printf("[WARN] Error add tags to CSMS secret: %s, err=%s", rst.ID, err)
		}
	}

	return resourceCsmsSecretRead(ctx, d, meta)
}

// The method can be used to remove empty strings from an array.
// When the sercet no associated event, the query secret API return the attribute(event_subscriptions) as follow:
// "event_subscriptions": [null]
// so need the function dealing the situation.
func removeNullValues(s []interface{}) []interface{} {
	result := make([]interface{}, 0, len(s))
	for _, elem := range s {
		if v, ok := elem.(string); ok && v != "" {
			result = append(result, v)
		}
	}
	return result
}

func resourceCsmsSecretRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// The endpoint of CSMS is the endpoint of KMS.
	client, err := cfg.KmsV1Client(region)
	if err != nil {
		return diag.Errorf("failed to create CSMS(KMS) client: %s", err)
	}

	id, name := parseID(d.Id())
	// Query secret details
	secret, err := secrets.Get(client, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "failed to query CSMS secret details")
	}

	createTime := time.Unix(int64(secret.CreateTime)/1000, 0).UTC().Format("2006-01-02 15:04:05 MST")
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("secret_id", secret.ID),
		d.Set("name", secret.Name),
		d.Set("kms_key_id", secret.KmsKeyID),
		d.Set("description", secret.Description),
		d.Set("status", secret.State),
		d.Set("create_time", createTime),
		d.Set("secret_type", secret.SecretType),
		d.Set("enterprise_project_id", secret.EnterpriseProjectID),
		d.Set("event_subscriptions", removeNullValues(secret.EventSubscriptions)),
	)

	// Query secret version
	version, err := queryLatestVersion(cfg, region, name)
	if err != nil {
		mErr = multierror.Append(
			mErr,
			err)
	}
	var encodedSecretTxt string
	var encodedSecretBinary string
	if version.SecretString != "" {
		encodedSecretTxt = utils.HashAndHexEncode(version.SecretString)
		mErr = multierror.Append(
			mErr,
			d.Set("secret_text", encodedSecretTxt),
		)
	} else {
		encodedSecretBinary = utils.HashAndHexEncode(version.SecretBinary)
		mErr = multierror.Append(
			mErr,
			d.Set("secret_binary", encodedSecretBinary),
		)
	}

	if version.VersionMetadata.ExpireTime > 0 {
		mErr = multierror.Append(
			mErr,
			d.Set("expire_time", version.VersionMetadata.ExpireTime),
		)
	}
	versionID := version.VersionMetadata.ID
	mErr = multierror.Append(
		mErr,
		d.Set("latest_version", versionID),
	)

	// Query secret tags
	if resourceTags, err := tags.Get(client, serviceType, id).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(
			mErr,
			d.Set("tags", tagMap),
		)
	} else {
		log.Printf("[WARN] error querying CSMS secret tags (%s): %s", id, err)
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("failed to set attributes for CSMS secret: %s", mErr)
	}
	return nil
}

func queryLatestVersion(cfg *config.Config, region, name string) (*secrets.Version, error) {
	client, err := cfg.KmsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("failed to create CSMS(KMS) client: %s", err)
	}

	// Query the version list
	versions, err := secrets.ListSecretVersions(client, name)
	if err != nil {
		return nil, fmt.Errorf("failed to query the list of secret versions: %s", err)
	}
	// Sort by created time in descending order.
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].CreateTime > versions[j].CreateTime
	})

	versionID := versions[0].ID

	return queryVersion(cfg, region, name, versionID)
}

func queryVersion(cfg *config.Config, region, name, versionID string) (*secrets.Version, error) {
	client, err := cfg.KmsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("failed to create CSMS(KMS) client: %s", err)
	}

	// Query version
	version, err := secrets.ShowSecretVersion(client, name, versionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query secret version: %s", err)
	}
	return version, nil
}

func resourceCsmsSecretUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// The endpoint of CSMS is the endpoint of KMS.
	client, err := cfg.KmsV1Client(region)
	if err != nil {
		return diag.Errorf("failed to create CSMS(KMS) client: %s", err)
	}

	id, name := parseID(d.Id())
	// Update secret basic-info
	if d.HasChanges("kms_key_id", "description", "event_subscriptions") {
		desc := d.Get("description").(string)
		kmsKeyID := d.Get("kms_key_id").(string)
		events := d.Get("event_subscriptions").(*schema.Set)
		opts := secrets.UpdateSecretOpts{
			KmsKeyID:           kmsKeyID,
			Description:        &desc,
			EventSubscriptions: utils.ExpandToStringListBySet(events),
		}
		log.Printf("[DEBUG] The option to update the basic information of the CSMS secret is: %#v", opts)

		_, err = secrets.Update(client, name, opts)
		if err != nil {
			return diag.Errorf("failed to update the base-info of CSMS secret: %s", err)
		}
	}

	// Update secret text
	if d.HasChanges("secret_text", "secret_binary") {
		opts := secrets.CreateVersionOpts{
			SecretString: d.Get("secret_text").(string),
			SecretBinary: d.Get("secret_binary").(string),
		}

		_, err = secrets.CreateSecretVersion(client, name, opts)
		if err != nil {
			return diag.Errorf("failed to create a new version of CSMS secret: %s", err)
		}
	}

	// The expire_time can be update only when the secret key state is ENABLED.
	if d.HasChanges("expire_time") && d.Get("status") == "ENABLED" {
		err = updateSecretVersion(client, d, name)
		if err != nil {
			return diag.Errorf("failed to update the CSMS secret version: %s", err)
		}
	}

	// Update tags
	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(client, d, serviceType, id)
		if err != nil {
			return diag.Errorf("failed to update CSMS secret tags: %s", err)
		}
	}
	return resourceCsmsSecretRead(ctx, d, meta)
}

func updateSecretVersion(client *golangsdk.ServiceClient, d *schema.ResourceData, name string) error {
	opts := secrets.UpdateVersionOpts{
		ExpireTime: d.Get("expire_time").(int),
	}

	// Query the version list
	versions, err := secrets.ListSecretVersions(client, name)
	if err != nil {
		return fmt.Errorf("failed to query the list of secret versions: %s", err)
	}
	if len(versions) == 0 {
		return fmt.Errorf("the list of secret versions is empty")
	}

	// Sort by created time in descending order.
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].CreateTime > versions[j].CreateTime
	})

	versionID := versions[0].ID
	_, err = secrets.UpdateSecretVersion(client, name, versionID, opts)
	if err != nil {
		return fmt.Errorf("failed to update the expire time of CSMS secret: %s", err)
	}

	return nil
}

func resourceCsmsSecretDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// The endpoint of CSMS is the endpoint of KMS.
	client, err := cfg.KmsV1Client(region)
	if err != nil {
		return diag.Errorf("failed to create CSMS(KMS) client: %s", err)
	}

	name := d.Get("name").(string)
	err = secrets.Delete(client, name)
	if err != nil {
		return diag.Errorf("failed to delete CSMS secret: %s", err)
	}
	d.SetId("")
	return nil
}

func resourceCsmsSecretImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	id, name := parseID(d.Id())
	if id == "" {
		err := fmt.Errorf("invalid format specified for the ID of CSMS secret. " +
			"Format must be <id>/<name>")
		return nil, err
	}

	d.Set("secret_id", id)
	d.Set("name", name)
	return []*schema.ResourceData{d}, nil
}

func parseID(id string) (string, string) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
