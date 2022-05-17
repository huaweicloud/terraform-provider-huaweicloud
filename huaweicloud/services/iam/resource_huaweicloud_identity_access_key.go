package iam

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/credentials"
	"github.com/chnsz/golangsdk/openstack/identity/v3.0/users"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/encryption"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceIdentityKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityKeyCreate,
		ReadContext:   resourceIdentityKeyRead,
		UpdateContext: resourceIdentityKeyUpdate,
		DeleteContext: resourceIdentityKeyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"pgp_key": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"active", "inactive",
				}, false),
			},
			"key_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"encrypted_secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"user_name": {
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

func resourceIdentityKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud iam client: %s", err)
	}

	userID := d.Get("user_id").(string)
	userInfo, err := users.Get(iamClient, userID).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error fetching iam user %s: %s", userID, err)
	}
	userName := userInfo.Name
	logp.Printf("[DEBUG] Create an access key for user %s", userName)

	opts := credentials.CreateOpts{
		UserID:      userID,
		Description: d.Get("description").(string),
	}
	accessKey, err := credentials.Create(iamClient, opts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating access key: %s", err)
	}
	d.SetId(accessKey.AccessKey)

	var outputFile string
	if v, ok := d.GetOk("secret_file"); ok {
		outputFile = v.(string)
	} else {
		outputFile = fmt.Sprintf("credentials-%s.csv", userName)
	}

	if err := writeToCSVFile(outputFile, accessKey); err != nil {
		// set the SecretKey as it was returned only in creation response
		d.Set("secret", accessKey.SecretKey)
		return fmtp.DiagErrorf("Error saving the access key to %s: %s", outputFile, err)
	}

	if v, ok := d.GetOk("pgp_key"); ok {
		pgpKey := v.(string)
		encryptionKey, err := encryption.RetrieveGPGKey(pgpKey)
		if err != nil {
			return fmtp.DiagErrorf("Error retrieving PGP key: %s", err)
		}
		fingerprint, encrypted, err := encryption.EncryptValue(encryptionKey, accessKey.SecretKey, "IAM Access Key Secret")
		if err != nil {
			return fmtp.DiagErrorf("Error encrypting access key: %s", err)
		}

		mErr := multierror.Append(nil,
			d.Set("key_fingerprint", fingerprint),
			d.Set("encrypted_secret", encrypted),
		)
		if err = mErr.ErrorOrNil(); err != nil {
			return fmtp.DiagErrorf("error setting identity access key fields: %s", err)
		}
	}

	d.Set("user_name", userName)
	return resourceIdentityKeyRead(ctx, d, meta)
}

func resourceIdentityKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud iam client: %s", err)
	}

	accessKey, err := credentials.Get(iamClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "access key")
	}

	mErr := multierror.Append(nil,
		d.Set("status", accessKey.Status),
		d.Set("create_time", accessKey.CreateTime),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting identity access key fields: %s", err)
	}

	return nil
}

func resourceIdentityKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud iam client: %s", err)
	}

	id := d.Id()
	if d.HasChanges("description", "status") {
		opts := credentials.UpdateOpts{
			Description: d.Get("description").(string),
			Status:      d.Get("status").(string),
		}
		_, err := credentials.Update(iamClient, id, opts).Extract()
		if err != nil {
			return fmtp.DiagErrorf("Error updating HuaweiCloud iam access key: %s", err)
		}
	}

	return resourceIdentityKeyRead(ctx, d, meta)
}

func resourceIdentityKeyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud iam client: %s", err)
	}

	if err := credentials.Delete(iamClient, d.Id()).ExtractErr(); err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloud iam access key: %s", err)
	}

	d.SetId("")
	return nil
}

func writeToCSVFile(path string, cred *credentials.Credential) error {
	var csvFile *os.File

	csvFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	var data = make([][]string, 2)
	data[0] = []string{"User ID", "Access Key ID", "Secret Access Key"}
	data[1] = []string{cred.UserID, cred.AccessKey, cred.SecretKey}

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
