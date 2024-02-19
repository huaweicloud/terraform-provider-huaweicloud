package apig

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apis"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apis/action
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apis/publish/{api_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/apis/publish/{api_id}
func ResourceApigApiPublishment() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceApiPublishmentCreate,
		ReadContext:   ResourceApiPublishmentRead,
		UpdateContext: ResourceApiPublishmentUpdate,
		DeleteContext: ResourceApiPublishmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApiPublishmentImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to publish API.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the API and the environment belongs.",
			},
			"env_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The ID of the environment to which the current version of the API will be published or " +
					"has been published.",
			},
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the API to be published or already published.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the current publishment.",
			},
			"version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The version ID of the current publishment.",
			},
			// Attributes
			"env_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the environment to which the current version of the API is published.",
			},
			"published_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time when the current version was published.",
			},
			"publish_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The publish ID of the API in current environment.",
			},
			"histories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version ID of the API publishment.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version description of the API publishment.",
						},
					},
				},
				Description: "All publish informations of the API.",
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

// ResourceApiPublishmentCreate is a function that uses terraform configuration to publish new versions or switch
// historical versions.
func ResourceApiPublishmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Get("api_id").(string)
		envId      = d.Get("env_id").(string)
	)

	// If version_id is configured, try to switch to the corresponding version.
	if versionId, hasVer := d.GetOk("version_id"); hasVer {
		histories, err := GetVersionHistories(client, instanceId, envId, apiId)
		if err != nil {
			return diag.Errorf("error finding the publish versions of the API (%s) in the environment (%s): %s",
				apiId, envId, err)
		}
		// Whether the API of the current version has been published, if not, throw an error.
		if ver, ok := isPublished(histories, versionId.(string)); !ok {
			return diag.Errorf("the version (%s) has not published", versionId.(string))
		} else if desc, ok := d.GetOk("description"); ok && desc != ver.Description {
			return diag.Errorf("the description is no correct, want '%s', but '%s', please check your description "+
				"input or API version", ver.Description, desc)
		}
		// If the API have historical versions, check whether the current version is the latest version, if not,
		// switch to this version.
		if _, ok := isLatestVersion(histories, versionId.(string)); !ok {
			_, err := apis.SwitchSpecVersion(client, instanceId, apiId, versionId.(string)).Extract()
			if err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		// If version_id is not configured, try to publish a new version.
		err = publishApiToSpecifiedEnv(client, instanceId, envId, apiId, d.Get("description").(string))
		if err != nil {
			return diag.Errorf("error publishing API: %s", err)
		}
	}

	// The ID is constructed from the instance ID, environment ID, and API ID, separated by slashes.
	d.SetId(fmt.Sprintf("%s/%s/%s", instanceId, envId, apiId))

	return ResourceApiPublishmentRead(ctx, d, meta)
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
		return nil, fmt.Errorf("the API does not have any published information")
	}
	for _, ver := range resp {
		// Status 1 means that this version is a effective version.
		if ver.Status == 1 {
			return &ver, nil
		}
	}
	return nil, fmt.Errorf("unable to find any publish information for the API")
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

func flattenResourceId(id string) (instanceId, envId, apiId string, err error) {
	ids := strings.Split(id, "/")
	if len(ids) != 3 {
		err = fmt.Errorf("invalid ID format, want '<instance_id>/<environment_id>/<api_id>', but '%s'", id)
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
// e.g.
// Publish IDs: {ID A}|{ID B}|{ID C}
// Environment IDs: {Related ID A}|{Related ID B}|{Related ID C}
func getPublishIdByEnvId(client *golangsdk.ServiceClient, instanceId, apiId, envId string) (string, error) {
	resp, err := apis.Get(client, instanceId, apiId).Extract()
	if err != nil {
		return "", err
	}
	var (
		publishIds = strings.Split(resp.PublishId, "|")
		envIds     = strings.Split(resp.RunEnvId, "|")
	)
	log.Printf("[DEBUG] The list of publish IDs is: %#v", publishIds)
	log.Printf("[DEBUG] The list of environment IDs is: %#v", envIds)

	for i, val := range envIds {
		if val == envId {
			if len(publishIds) < i {
				return "", fmt.Errorf("the length of publish ID list is not correct, want '%d', but '%d'",
					len(envIds), len(publishIds))
			}
			return publishIds[i], nil
		}
	}
	return "", fmt.Errorf("the API is not published in this environment (%s)", envId)
}

// ResourceApiPublishmentRead is a method to obtain informations of API publishment and save to the local storage.
func ResourceApiPublishmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	instanceId, envId, apiId, err := flattenResourceId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := GetVersionHistories(client, instanceId, envId, apiId)
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
		d.Set("published_at", publishInfo.PublishTime),
		d.Set("env_name", publishInfo.EnvName),
		setApiPublishHistories(d, resp),
	)

	if publishId, err := getPublishIdByEnvId(client, instanceId, publishInfo.ApiId, publishInfo.EnvId); err != nil {
		mErr = multierror.Append(mErr, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("publish_id", publishId))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving publishment fields: %s", err)
	}
	return nil
}

func ResourceApiPublishmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	versionId := d.Get("version_id").(string)
	instanceId, envId, apiId, err := flattenResourceId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if versionId == "" {
		// If the version of the configuration is empty, whether description is changed or not, publish a new version.
		if err = publishApiToSpecifiedEnv(client, instanceId, envId, apiId, d.Get("description").(string)); err != nil {
			return diag.Errorf("error publishing API: %s", err)
		}
	} else {
		if !d.HasChange("version_id") && d.HasChange("description") {
			return diag.Errorf("only for new API publishment, the description can be updated")
		}
		description := d.Get("description").(string)

		// Obtain the version history of the API from the specified environment and check whether the current version
		// has been published.
		histories, err := GetVersionHistories(client, instanceId, envId, apiId)
		if err != nil {
			return diag.Errorf("error getting version histories of the API (%s): %s", apiId, err)
		}
		if ver, ok := isPublished(histories, versionId); !ok {
			return diag.Errorf("this version (%s) has not published", versionId)
		} else if description != "" && ver.Description != description {
			// If user want to switch an exist version, but the description is not right, throw an error.
			return diag.Errorf("this description is not belongs to version (%s)", versionId)
		}

		if _, err := apis.SwitchSpecVersion(client, instanceId, apiId, versionId).Extract(); err != nil {
			return diag.Errorf("%s", err)
		}
	}

	return ResourceApiPublishmentRead(ctx, d, meta)
}

func ResourceApiPublishmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	instanceId, envId, apiId, err := flattenResourceId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	err = offlineApiFromSpecifiedEnv(client, instanceId, envId, apiId)
	if err != nil {
		return diag.Errorf("error offlining API: %s", err)
	}
	d.SetId("")

	return nil
}

func resourceApiPublishmentImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	instanceId, envId, apiId, err := flattenResourceId(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	// Set all required parameters.
	mErr := multierror.Append(nil,
		d.Set("instance_id", instanceId),
		d.Set("env_id", envId),
		d.Set("api_id", apiId),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
