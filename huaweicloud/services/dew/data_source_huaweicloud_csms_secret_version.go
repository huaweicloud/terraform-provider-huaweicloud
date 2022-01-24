package dew

import (
	"context"
	"fmt"
	"time"

	"github.com/chnsz/golangsdk/openstack/csms/v1/secrets"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func DataSourceDewCsmsSecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDewCsmsSecretRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secret_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secret_text": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDewCsmsSecretRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var version *secrets.Version
	var err error

	secretName := d.Get("secret_name").(string)
	if ver, ok := d.GetOk("version"); ok {
		version, err = queryVersion(config, region, secretName, ver.(string))
	} else {
		version, err = queryLatestVersion(config, region, secretName)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	vMetadata := version.VersionMetadata
	id := fmt.Sprintf("%s/%s", vMetadata.ID, vMetadata.SecretName)
	d.SetId(id)

	createTime := time.Unix(int64(vMetadata.CreateTime)/1000, 0).UTC().Format("2006-01-02 15:04:05 MST")
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("secret_text", version.SecretString),
		d.Set("secret_name", vMetadata.SecretName),
		d.Set("kms_key_id", vMetadata.KmsKeyID),
		d.Set("version", vMetadata.ID),
		d.Set("status", vMetadata.VersionStages),
		d.Set("created_at", createTime),
	)

	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("error setting CSMS secret attributes: %s", mErr)
	}
	return nil
}
