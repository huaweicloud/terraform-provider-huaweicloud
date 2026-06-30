package cce

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE POST /api/v3/projects/{project_id}/clusters/{cluster_id}/assume-agency-for-pod-identity
func DataSourceCCEClusterPodIdentityAssumeAgency() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEClusterPodIdentityAssumeAgencyRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the cluster ID.",
			},
			"token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"assumed_agency": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"audience": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_access_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"security_token": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"expiration": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"pod_identity_association_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_account": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCCEClusterPodIdentityAssumeAgencyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)

	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	requestURL := "api/v3/projects/{project_id}/clusters/{cluster_id}/assume-agency-for-pod-identity"
	requestPath := client.Endpoint + requestURL
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", clusterID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"token": d.Get("token").(string),
		},
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error getting assumed agency for pod identity: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error retrieving assumed agency for pod identity: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("assumed_agency", flattenAssumedAgency(respBody)),
		d.Set("audience", utils.PathSearch("audience", respBody, nil)),
		d.Set("credentials", flattenCredentials(respBody)),
		d.Set("pod_identity_association_id", utils.PathSearch("podIdentityAssociationId", respBody, nil)),
		d.Set("subject", flattenSubject(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAssumedAgency(resp interface{}) []map[string]interface{} {
	assumedAgency := utils.PathSearch("assumedAgency", resp, nil)
	if assumedAgency == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"urn": utils.PathSearch("urn", assumedAgency, nil),
			"id":  utils.PathSearch("id", assumedAgency, nil),
		},
	}
}

func flattenCredentials(resp interface{}) []map[string]interface{} {
	credentials := utils.PathSearch("credentials", resp, nil)
	if credentials == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"access_key_id":     utils.PathSearch("accessKeyId", credentials, nil),
			"secret_access_key": utils.PathSearch("secretAccessKey", credentials, nil),
			"security_token":    utils.PathSearch("securityToken", credentials, nil),
			"expiration":        utils.PathSearch("expiration", credentials, nil),
		},
	}
}

func flattenSubject(resp interface{}) []map[string]interface{} {
	subject := utils.PathSearch("subject", resp, nil)
	if subject == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"namespace":       utils.PathSearch("namespace", subject, nil),
			"service_account": utils.PathSearch("serviceAccount", subject, nil),
		},
	}
}
