package iam

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/credentials"
	"github.com/chnsz/golangsdk/openstack/identity/v3.0/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/encryption"
)

// @API IAM POST /v3.0/OS-CREDENTIAL/credentials
// @API IAM PUT /v3.0/OS-CREDENTIAL/credentials/{access_key}
// @API IAM DELETE /v3.0/OS-CREDENTIAL/credentials/{access_key}
// @API IAM GET /v3.0/OS-CREDENTIAL/credentials/{access_key}
// @API IAM GET /v3.0/OS-USER/users/{user_id}
func ResourceAccessKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessKeyCreate,
		ReadContext:   resourceAccessKeyRead,
		UpdateContext: resourceAccessKeyUpdate,
		DeleteContext: resourceAccessKeyDelete,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The IAM user ID for which access key (AK/SK) to be created.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the access key.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The status of the access key.`,
			},
			"secret_file": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The file name of the credentials (CSV) that can save access key and access secret key.`,
			},
			"pgp_key": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: `The PGP public key (base64 encoded) used to encrypt the storaged secret key.`,
			},
			"secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `The access secret key.`,
			},
			"key_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fingerprint of the PGP key used to encrypt the secret.`,
			},
			"encrypted_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The encrypted secret, which encoded in base64.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the IAM user.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the access key, in ISO-8601 UTC format.`,
			},
		},
	}
}

func storeAccessKeyToCsvFile(path string, cred *credentials.Credential) error {
	var (
		csvFile *os.File
		data    = [][]string{
			{"User ID", "Access Key ID", "Secret Access Key"},
			{cred.UserID, cred.AccessKey, cred.SecretKey},
		}
	)

	csvFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	_, err = csvFile.WriteString("\xEF\xBB\xBF")
	if err != nil {
		return err
	}

	writer := csv.NewWriter(csvFile)
	err = writer.WriteAll(data)
	if err != nil {
		return err
	}

	writer.Flush()

	return nil
}

func resourceAccessKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		mErr   *multierror.Error
		userId = d.Get("user_id").(string)
		diags  diag.Diagnostics
	)
	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	// Precheck the user information and ensure the user exists.
	userInfo, err := users.Get(client, userId).Extract()
	if err != nil {
		return diag.Errorf("error fetching IAM user %s: %s", userId, err)
	}

	opts := credentials.CreateOpts{
		UserID:      userId,
		Description: d.Get("description").(string),
	}
	accessKey, err := credentials.Create(client, opts).Extract()
	if err != nil {
		return diag.Errorf("error creating permanent access key: %s", err)
	}

	d.SetId(accessKey.AccessKey)
	userName := userInfo.Name
	mErr = multierror.Append(mErr, d.Set("user_name", userName))

	// Default storage in the current execution directory, can also be specified by 'secret_file'.
	outputFile := fmt.Sprintf("credentials-%s.csv", userName)
	if customStoragePath, ok := d.GetOk("secret_file"); ok {
		outputFile = customStoragePath.(string)
	}

	if err := storeAccessKeyToCsvFile(outputFile, accessKey); err != nil {
		// When the CSV file fails to be saved, the secret key is stored in tfstate to prevent the value from being lost.
		mErr = multierror.Append(mErr, d.Set("secret", accessKey.SecretKey))
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  `Unable to save the secret key to the specified path, which is configured by 'secret_file'`,
			Detail:   fmt.Sprintf("unable to save the secret key to %s: %s", outputFile, err),
		})
	}

	if v, ok := d.GetOk("pgp_key"); ok {
		encryptionKey, err := encryption.RetrieveGPGKey(v.(string))
		if err != nil {
			return diag.Errorf("error retrieving PGP key: %s", err)
		}
		fingerprint, encrypted, err := encryption.EncryptValue(encryptionKey, accessKey.SecretKey, "IAM Access Key Secret")
		if err != nil {
			return diag.Errorf("error encrypting access key using PGP key: %s", err)
		}

		mErr = multierror.Append(mErr,
			d.Set("key_fingerprint", fingerprint),
			d.Set("encrypted_secret", encrypted),
		)
		if err = mErr.ErrorOrNil(); err != nil {
			return diag.Errorf("error setting PGP encryption fields of permanent access key: %s", err)
		}
	}

	return append(diags, resourceAccessKeyRead(ctx, d, meta)...)
}

func resourceAccessKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		accessKeyId = d.Id()
	)
	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	accessKey, err := credentials.Get(client, accessKeyId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "access key")
	}

	mErr := multierror.Append(nil,
		d.Set("status", accessKey.Status),
		d.Set("create_time", accessKey.CreateTime),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting resource fields of permanent access key: %s", err)
	}

	return nil
}

func resourceAccessKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		accessKeyId = d.Id()
	)
	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if d.HasChanges("description", "status") {
		opts := credentials.UpdateOpts{
			Description: d.Get("description").(string),
			Status:      d.Get("status").(string),
		}
		_, err := credentials.Update(client, accessKeyId, opts).Extract()
		if err != nil {
			return diag.Errorf("error updating permanent access key: %s", err)
		}
	}

	return resourceAccessKeyRead(ctx, d, meta)
}

func resourceAccessKeyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		accessKeyId = d.Id()
	)
	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if err := credentials.Delete(client, accessKeyId).ExtractErr(); err != nil {
		return diag.Errorf("error deleting permanent access key: %s", err)
	}

	return nil
}
