package dew

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/csms/v1/secrets"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	serviceType = "csms"
)

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
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				StateFunc: utils.HashAndHexEncode,
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
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	// The endpoint of CSMS is the endpoint of KMS.
	client, err := config.KmsV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("failed to create HuaweiCloud CSMS(KMS) client: %s", err)
	}

	name := d.Get("name").(string)
	createOpts := secrets.CreateSecretOpts{
		Name:        name,
		KmsKeyID:    d.Get("kms_key_id").(string),
		Description: d.Get("description").(string),
	}
	logp.Printf("[DEBUG] Create CSMS secret options: %s", createOpts)
	createOpts.SecretString = d.Get("secret_text").(string)

	rst, err := secrets.Create(client, createOpts)
	if err != nil {
		return fmtp.DiagErrorf("failed to create the CSMS secret: %s", err)
	}

	id := fmt.Sprintf("%s/%s", rst.ID, name)
	d.SetId(id)

	// Save tags
	if t, ok := d.GetOk("tags"); ok {
		tMaps := t.(map[string]interface{})
		tagMaps := utils.ExpandResourceTags(tMaps)
		err = tags.Create(client, serviceType, rst.ID, tagMaps).ExtractErr()
		if err != nil {
			logp.Printf("[WARN] Error add tags to CSMS secret: %s, err=%s", rst.ID, err)
		}
	}

	return resourceCsmsSecretRead(ctx, d, meta)
}

func resourceCsmsSecretRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	// The endpoint of CSMS is the endpoint of KMS.
	client, err := config.KmsV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("failed to create HuaweiCloud CSMS(KMS) client: %s", err)
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
	)

	// Query secret version
	version, err := queryLatestVersion(config, region, name)
	if err != nil {
		mErr = multierror.Append(
			mErr,
			err)
	}
	secretTxt := version.SecretString
	encodedSecretTxt := utils.HashAndHexEncode(secretTxt)
	versionID := version.VersionMetadata.ID
	mErr = multierror.Append(
		mErr,
		d.Set("secret_text", encodedSecretTxt),
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
		logp.Printf("[WARN] Error querying CSMS secret tags (%s): %s", id, err)
	}

	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("failed to set attributes for CSMS secret: %s", mErr)
	}
	return nil
}

func queryLatestVersion(config *config.Config, region, name string) (*secrets.Version, error) {
	client, err := config.KmsV1Client(region)
	if err != nil {
		return nil, fmtp.Errorf("failed to create HuaweiCloud CSMS(KMS) client: %s", err)
	}

	// Query the version list
	versions, err := secrets.ListSecretVersions(client, name)
	if err != nil {
		return nil, fmtp.Errorf("failed to query the list of secret versions: %s", err)
	}
	// Sort by created time in descending order.
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].CreateTime > versions[j].CreateTime
	})

	versionID := versions[0].ID

	return queryVersion(config, region, name, versionID)
}

func queryVersion(config *config.Config, region, name, versionID string) (*secrets.Version, error) {
	client, err := config.KmsV1Client(region)
	if err != nil {
		return nil, fmtp.Errorf("failed to create HuaweiCloud CSMS(KMS) client: %s", err)
	}

	// Query version
	version, err := secrets.ShowSecretVersion(client, name, versionID)
	if err != nil {
		return nil, fmtp.Errorf("failed to query secret version: %s", err)
	}
	return version, nil
}

func resourceCsmsSecretUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	// The endpoint of CSMS is the endpoint of KMS.
	client, err := config.KmsV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("failed to create HuaweiCloud CSMS(KMS) client: %s", err)
	}

	var mErr *multierror.Error
	id, name := parseID(d.Id())
	// Update secret basic-info
	if d.HasChanges("kms_key_id", "description") {
		desc := d.Get("description").(string)
		kmsKeyID := d.Get("kms_key_id").(string)
		opts := secrets.UpdateSecretOpts{
			KmsKeyID:    kmsKeyID,
			Description: &desc,
		}
		logp.Printf("[DEBUG] The option to update the basic information of the CSMS secret is: %#v", opts)

		_, err = secrets.Update(client, name, opts)
		if err != nil {
			e := fmtp.Errorf("failed to update the base-info of CSMS secret: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	// Update secret text
	if d.HasChanges("secret_text") {
		opts := secrets.CreateVersionOpts{
			SecretString: d.Get("secret_text").(string),
		}
		_, err = secrets.CreateSecretVersion(client, name, opts)
		if err != nil {
			e := fmtp.Errorf("failed to create a new version of CSMS secret: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	// Update tags
	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(client, d, serviceType, id)
		if err != nil {
			e := fmtp.Errorf("failed to update CSMS secret tags: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("failed to update CSMS secret: %s", mErr)
	}
	return resourceCsmsSecretRead(ctx, d, meta)
}

func resourceCsmsSecretDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	// The endpoint of CSMS is the endpoint of KMS.
	client, err := config.KmsV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("failed to create HuaweiCloud CSMS(KMS) client: %s", err)
	}

	name := d.Get("name").(string)
	err = secrets.Delete(client, name)
	if err != nil {
		return fmtp.DiagErrorf("failed to delete CSMS secret: %s", err)
	}
	d.SetId("")
	return nil
}

func resourceCsmsSecretImport(ctx context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	id, name := parseID(d.Id())
	if id == "" {
		err := fmtp.Errorf("Invalid format specified for the ID of CSMS secret. " +
			"Format must be <id>/<name>")
		return nil, err
	}

	d.Set("secret_id", id)
	d.Set("name", name)
	return []*schema.ResourceData{d}, nil
}

func parseID(id string) (string, string) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
