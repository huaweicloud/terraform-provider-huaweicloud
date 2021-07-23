package huaweicloud

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3.0/credentials"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3.0/users"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func resourceIdentityKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityKeyCreate,
		Read:   resourceIdentityKeyRead,
		Update: resourceIdentityKeyUpdate,
		Delete: resourceIdentityKeyDelete,

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

func resourceIdentityKeyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud iam client: %s", err)
	}

	userID := d.Get("user_id").(string)
	userInfo, err := users.Get(iamClient, userID).Extract()
	if err != nil {
		return fmtp.Errorf("Error fetching iam user %s: %s", userID, err)
	}
	userName := userInfo.Name
	logp.Printf("[DEBUG] Create an access key for user %s", userName)

	opts := credentials.CreateOpts{
		UserID:      userID,
		Description: d.Get("description").(string),
	}
	accessKey, err := credentials.Create(iamClient, opts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating access key: %s", err)
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
		return fmtp.Errorf("Error saving the access key to %s: %s", outputFile, err)
	}

	d.Set("user_name", userName)
	return resourceIdentityKeyRead(d, meta)
}

func resourceIdentityKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud iam client: %s", err)
	}

	accessKey, err := credentials.Get(iamClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "access key")
	}

	d.Set("status", accessKey.Status)
	d.Set("create_time", accessKey.CreateTime)
	return nil
}

func resourceIdentityKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud iam client: %s", err)
	}

	id := d.Id()
	if d.HasChanges("description", "status") {
		opts := credentials.UpdateOpts{
			Description: d.Get("description").(string),
			Status:      d.Get("status").(string),
		}
		_, err := credentials.Update(iamClient, id, opts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud iam access key: %s", err)
		}
	}

	return resourceIdentityKeyRead(d, meta)
}

func resourceIdentityKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud iam client: %s", err)
	}

	if err := credentials.Delete(iamClient, d.Id()).ExtractErr(); err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud iam access key: %s", err)
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

	csvFile.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(csvFile)
	writer.WriteAll(data)
	writer.Flush()

	return nil
}
