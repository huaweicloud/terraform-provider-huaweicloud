package iam

import (
	"context"
	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"log"
)

// @IAM GET /v3.0/OS-CREDENTIAL/credentials
func DataSourceIdentityKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityKeyRead,
		Schema: map[string]*schema.Schema{
			"credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityKeyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cfg := m.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	getKeyOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", client.Endpoint+"v3.0/OS-CREDENTIAL/credentials", &getKeyOpts)
	if err != nil {
		return diag.Errorf("failed to fetch credentials: %s", err)
	}

	getKeyRespBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("****************++++++++++++++access body:%v", getKeyRespBody)

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		nil,
		d.Set("credentials", flattenKey(getKeyRespBody)),
		//d.Set("user_id", utils.PathSearch("credentials.user_id", getKeyRespBody, nil)),
		//d.Set("access", utils.PathSearch("credentials.access", getKeyRespBody, nil)),
		//d.Set("status", utils.PathSearch("credentials.status", getKeyRespBody, nil)),
		//d.Set("create_time", utils.PathSearch("credentials.create_time", getKeyRespBody, nil)),
		//d.Set("description", utils.PathSearch("credentials.description", getKeyRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenKey(getKeyRespBody interface{}) []map[string]interface{} {
	keyRaw := utils.PathSearch("credentials", getKeyRespBody, nil)
	if keyRaw == nil {
		return nil
	}
	key := keyRaw.([]interface{})

	res := make([]map[string]interface{}, len(key))
	for i, v := range key {
		res[i] = map[string]interface{}{
			"user_id":     utils.PathSearch("user_id", v, nil),
			"access":      utils.PathSearch("access", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"create_time": utils.PathSearch("create_time", v, nil),
			"description": utils.PathSearch("description", v, nil),
		}
	}
	return res
}
