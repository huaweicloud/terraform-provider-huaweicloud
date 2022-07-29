package apig

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apis"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceApigApiPublishment() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceApigApiPublishmentCreate,
		ReadContext:   ResourceApigApiPublishmentRead,
		UpdateContext: ResourceApigApiPublishmentUpdate,
		DeleteContext: ResourceApigApiPublishmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApiPublishmentImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"env_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"env_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publish_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publish_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"histories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_id": {
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

func isPublished(versionList []apis.ApiVersionInfo, versionId string) (*apis.ApiVersionInfo, bool) {
	for _, ver := range versionList {
		if ver.VersionId == versionId {
			return &ver, true
		}
	}
	return nil, false
}

func isLatestVersion(versionList []apis.ApiVersionInfo, versionId string) (*apis.ApiVersionInfo, bool) {
	if versionList[0].VersionId == versionId {
		return &versionList[0], true
	}
	return nil, false
}

func publishApiToSpecifiedEnv(c *golangsdk.ServiceClient, instanceId, envId, apiId, description string) error {
	opts := apis.PublishOpts{
		Action:      "online",
		EnvId:       envId,
		ApiId:       apiId,
		Description: description,
	}
	_, err := apis.Publish(c, instanceId, opts).Extract()
	if err != nil {
		return err
	}
	return nil
}

func offlineApiFromSpecifiedEnv(c *golangsdk.ServiceClient, instanceId, envId, apiId string) error {
	opts := apis.PublishOpts{
		Action: "offline",
		EnvId:  envId,
		ApiId:  apiId,
	}
	_, err := apis.Publish(c, instanceId, opts).Extract()
	if err != nil {
		return err
	}
	return nil
}

// ResourceApigApiPublishmentCreate is a function that uses terraform configuration to publish new versions or switch
// historical versions.
func ResourceApigApiPublishmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	apiId := d.Get("api_id").(string)
	envId := d.Get("env_id").(string)
	// If version_id is configured, try to switch to the corresponding version.
	if versionId, hasVer := d.GetOk("version_id"); hasVer {
		histories, err := GetVersionHistories(c, instanceId, envId, apiId)
		if err != nil {
			return fmtp.DiagErrorf("Error finding the publish versions of the API (%s) in the environment(%s): %s",
				apiId, envId, err)
		}
		// Check if the current version has been published, if not, throw an error.
		if ver, ok := isPublished(histories, versionId.(string)); !ok {
			return fmtp.DiagErrorf("The version (%s) has not published.", versionId.(string))
		} else if desc, ok := d.GetOk("description"); ok && desc != ver.Description {
			return fmtp.DiagErrorf("The description of version (%s) is not %s, but %s.", versionId, desc, ver.Description)
		}
		// If the API have historical versions, check whether the current version is the latest version, if not,
		// switch to this version.
		if _, ok := isLatestVersion(histories, versionId.(string)); !ok {
			_, err := apis.SwitchSpecVersion(c, instanceId, apiId, versionId.(string)).Extract()
			if err != nil {
				return fmtp.DiagErrorf("%s", err)
			}
		}
	} else {
		// If version_id is not configured, try to publish a new version.
		err = publishApiToSpecifiedEnv(c, instanceId, envId, apiId, d.Get("description").(string))
		if err != nil {
			return fmtp.DiagErrorf("Error publishing API: %s", err)
		}
	}

	// The ID is constructed from the instance ID, environment ID, and API ID, separated by slashes.
	d.SetId(fmt.Sprintf("%s/%s/%s", instanceId, envId, apiId))

	return ResourceApigApiPublishmentRead(ctx, d, meta)
}

// GetVersionHistories is a function that obtains version histories by APIG IDs.
func GetVersionHistories(c *golangsdk.ServiceClient, instanceId, envId, apiId string) ([]apis.ApiVersionInfo, error) {
	// Get version histories of the API from the specified environment.
	pages, err := apis.ListPublishHistories(c, instanceId, apiId, apis.ListPublishHistoriesOpts{
		EnvId: envId,
	}).AllPages()
	if err != nil {
		return nil, err
	}
	histories, err := apis.ExtractHistories(pages)
	if err != nil {
		return nil, err
	}

	return histories, nil
}

func getCertainPublishInfo(resp []apis.ApiVersionInfo) (*apis.ApiVersionInfo, error) {
	if len(resp) == 0 {
		return nil, fmtp.Errorf("The API does not have any published information.")
	}
	for _, ver := range resp {
		// Status 1 means that this version is a effective version.
		if ver.Status == 1 {
			return &ver, nil
		}
	}
	return nil, fmtp.Errorf("Unable to find any publish information for the API.")
}

func setApiPublishHistories(d *schema.ResourceData, resp []apis.ApiVersionInfo) error {
	result := make([]map[string]interface{}, len(resp))
	for i, ver := range resp {
		result[i] = map[string]interface{}{
			"version_id":  ver.VersionId,
			"description": ver.Description,
		}
	}
	return d.Set("histories", result)
}

func getApigIdsFromPublishmentId(id string) (instanceId, envId, apiId string, err error) {
	ids := strings.Split(id, "/")
	if len(ids) != 3 {
		err = fmtp.Errorf("the format of the publishment ID should be 'instance_id/environment_id/api_id'.")
		return
	}
	instanceId = ids[0]
	envId = ids[1]
	apiId = ids[2]
	return
}

// The publishIdsStr string representing one or more API publication IDs, in descending order from left to right by
// publication time. The publish ID corresponds to the environment one by one.
// If the first matching environment ID is found, it means the publish ID under the current index we need.
func analysisPublishIds(publishIdsStr, envIdsStr, envId string) (string, error) {
	publishIds := strings.Split(publishIdsStr, "|")
	envIds := strings.Split(envIdsStr, "|")
	for i, val := range envIds {
		if val == envId {
			if len(publishIds) < i {
				log.Printf("[ERROR] The length of publish ID list is not right, the envIds is '%s', but the "+
					"publishIds is '%s'.", publishIdsStr, envIdsStr)
				return "", fmtp.Errorf("The publish ID is lost, please contact customer service")
			} else {
				return publishIds[i], nil
			}
		}
	}
	return "", fmtp.Errorf("the API is not published in this environment (%s)", envId)
}

func getPublishIdByApi(client *golangsdk.ServiceClient, instanceId, envId, apiId string) (string, error) {
	resp, err := apis.Get(client, instanceId, apiId).Extract()
	if err != nil {
		return "", err
	}
	publishId, err := analysisPublishIds(resp.PublishId, resp.RunEnvId, envId)
	if err != nil {
		return "", err
	}
	return publishId, nil
}

// ResourceApigApiPublishmentRead is a method to obtain informations of API publishment and save to the local storage.
func ResourceApigApiPublishmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}

	instanceId, envId, apiId, err := getApigIdsFromPublishmentId(d.Id())
	if err != nil {
		return fmtp.DiagErrorf("Wrong ID format: %s", err)
	}
	resp, err := GetVersionHistories(c, instanceId, envId, apiId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error getting the publish versions of the API (%s) in the environment (%s)", apiId, envId))
	}

	publishInfo, err := getCertainPublishInfo(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	mErr := multierror.Append(nil,
		d.Set("env_id", publishInfo.EnvId),
		d.Set("api_id", publishInfo.ApiId),
		d.Set("description", publishInfo.Description),
		d.Set("publish_time", publishInfo.PublishTime),
		d.Set("env_name", publishInfo.EnvName),
		setApiPublishHistories(d, resp),
	)

	if publishId, err := getPublishIdByApi(c, instanceId, publishInfo.EnvId, publishInfo.ApiId); err != nil {
		mErr = multierror.Append(mErr, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("publish_id", publishId))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func ResourceApigApiPublishmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}

	versionId := d.Get("version_id").(string)
	instanceId, envId, apiId, err := getApigIdsFromPublishmentId(d.Id())
	if err != nil {
		return fmtp.DiagErrorf("Wrong ID format: %s", err)
	}
	if versionId == "" {
		// If the version of the configuration is empty, whether description is changed or not, publish a new version.
		if err = publishApiToSpecifiedEnv(c, instanceId, envId, apiId, d.Get("description").(string)); err != nil {
			return fmtp.DiagErrorf("Error publishing API: %s", err)
		}
	} else {
		if !d.HasChange("version_id") && d.HasChange("description") {
			return fmtp.DiagErrorf("Only for new API publishment, the description can be updated.")
		}
		description := d.Get("description").(string)

		// Obtain the version history of the API from the specified environment and check whether the current version
		// has been published.
		histories, err := GetVersionHistories(c, instanceId, envId, apiId)
		if err != nil {
			return fmtp.DiagErrorf("Error getting version histories of the API (%s): %s", apiId, err)
		}
		if ver, ok := isPublished(histories, versionId); !ok {
			return fmtp.DiagErrorf("This version (%s) has not published.", versionId)
		} else if description != "" && ver.Description != description {
			// If user want to switch an exist version, but the description is not right, throw an error.
			return fmtp.DiagErrorf("This description is not belongs to version (%s).", versionId)
		}

		if _, err := apis.SwitchSpecVersion(c, instanceId, apiId, versionId).Extract(); err != nil {
			return fmtp.DiagErrorf("%s", err)
		}
	}

	return ResourceApigApiPublishmentRead(ctx, d, meta)
}

func ResourceApigApiPublishmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.ApigV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}

	instanceId, envId, apiId, err := getApigIdsFromPublishmentId(d.Id())
	if err != nil {
		return fmtp.DiagErrorf("Wrong ID format: %s", err)
	}
	err = offlineApiFromSpecifiedEnv(c, instanceId, envId, apiId)
	if err != nil {
		return fmtp.DiagErrorf("Error offlining API: %s", err)
	}
	d.SetId("")

	return nil
}

func resourceApiPublishmentImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	instanceId, envId, apiId, err := getApigIdsFromPublishmentId(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Wrong ID format: %s", err)
	}
	// Set all required parameters.
	d.Set("instance_id", instanceId)
	d.Set("env_id", envId)
	d.Set("api_id", apiId)

	return []*schema.ResourceData{d}, nil
}
