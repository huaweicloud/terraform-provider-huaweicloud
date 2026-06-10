package geminidb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var commandDisableNonUpdatableParams = []string{
	"instance_id",
	"disabled_type",
}

// @API GeminiDB POST /v3/{project_id}/redis/instances/{instance_id}/disabled-commands
// @API GeminiDB GET /v3/{project_id}/redis/instances/{instance_id}/disabled-commands
// @API GeminiDB DELETE /v3/{project_id}/redis/instances/{instance_id}/disabled-commands
func ResourceCommandDisable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCommandDisableCreate,
		ReadContext:   resourceCommandDisableRead,
		UpdateContext: resourceCommandDisableUpdate,
		DeleteContext: resourceCommandDisableDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCommandDisableImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(commandDisableNonUpdatableParams),

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
			},
			"disabled_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"commands": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"keys"},
			},
			"keys": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				Elem:          disableCommandKeysSchema(),
				ConflictsWith: []string{"commands"},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func disableCommandKeysSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"commands": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func buildCreateCommandDisableBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"disabled_type": d.Get("disabled_type"),
		"commands": utils.ValueIgnoreEmpty(
			utils.ExpandToStringList(d.Get("commands").(*schema.Set).List())),
		"keys": buildKeysInfoBodyParams(d.Get("keys").(*schema.Set).List()),
	}

	return bodyParams
}

func buildKeysInfoBodyParams(keysInfo []interface{}) []map[string]interface{} {
	if len(keysInfo) == 0 {
		return nil
	}

	keyInfos := make([]map[string]interface{}, 0, len(keysInfo))
	for _, v := range keysInfo {
		raw, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"db_id":    raw["db_id"],
			"key":      raw["key"],
			"commands": utils.ExpandToStringList(raw["commands"].(*schema.Set).List()),
		}
		keyInfos = append(keyInfos, params)
	}

	return keyInfos
}

func resourceCommandDisableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		createHttpUrl = "v3/{project_id}/redis/instances/{instance_id}/disabled-commands"
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateCommandDisableBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GeminiDB Redis disable command: %s", err)
	}

	d.SetId(instanceId)

	return resourceCommandDisableRead(ctx, d, meta)
}

func resourceCommandDisableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		commandType = d.Get("disabled_type").(string)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	commandInfos, err := GetCommandDisableInfo(client, d.Id(), commandType)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GeminiDB Redis disabled command")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", d.Id()),
		d.Set("disabled_type", commandType),
	)

	if commandType == "command" {
		mErr = multierror.Append(mErr,
			d.Set("commands", commandInfos),
			d.Set("keys", make([]interface{}, 0)),
		)
	} else {
		mErr = multierror.Append(mErr,
			d.Set("commands", make([]interface{}, 0)),
			d.Set("keys", flattenCommandDisableKeys(commandInfos)),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCommandDisableKeys(keysInfo []interface{}) []map[string]interface{} {
	if len(keysInfo) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(keysInfo))
	for _, item := range keysInfo {
		result = append(result, map[string]interface{}{
			"db_id":    utils.PathSearch("db_id", item, nil),
			"key":      utils.PathSearch("key", item, nil),
			"commands": utils.PathSearch("commands", item, nil),
		})
	}

	return result
}

func GetCommandDisableInfo(client *golangsdk.ServiceClient, instanceId, commandType string) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/redis/instances/{instance_id}/disabled-commands?type={type}&limit={limit}"
		offset  = 0
		limit   = 50
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{type}", commandType)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listOpts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &listOpts)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		if commandType == "command" {
			commands := utils.PathSearch("commands", respBody, make([]interface{}, 0)).([]interface{})
			result = append(result, commands...)
			if len(commands) < limit {
				break
			}

			offset += len(commands)
		}

		if commandType == "key" {
			keys := utils.PathSearch("keys", respBody, make([]interface{}, 0)).([]interface{})
			result = append(result, keys...)
			if len(keys) < limit {
				break
			}

			offset += len(keys)
		}
	}

	if len(result) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return result, nil
}

func resourceCommandDisableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	if d.HasChange("commands") {
		oldCommands, newCommands := d.GetChange("commands")
		addCommands := newCommands.(*schema.Set).Difference(oldCommands.(*schema.Set))
		removeCommands := oldCommands.(*schema.Set).Difference(newCommands.(*schema.Set))

		if removeCommands.Len() > 0 {
			err = removeDisableCommands(client, d, "command", removeCommands.List())
			if err != nil {
				return diag.Errorf("error removing GeminiDB Redis disabled command: %s", err)
			}
		}

		if addCommands.Len() > 0 {
			err = addDisableCommands(client, d, "command", addCommands.List())
			if err != nil {
				return diag.Errorf("error adding GeminiDB Redis disabled command: %s", err)
			}
		}
	}

	if d.HasChange("keys") {
		oldKeys, newKeys := d.GetChange("keys")
		addKeys := newKeys.(*schema.Set).Difference(oldKeys.(*schema.Set))
		removeKeys := oldKeys.(*schema.Set).Difference(newKeys.(*schema.Set))

		if removeKeys.Len() > 0 {
			err = removeDisableCommands(client, d, "key", removeKeys.List())
			if err != nil {
				return diag.Errorf("error removing GeminiDB Redis disabled command: %s", err)
			}
		}

		if addKeys.Len() > 0 {
			err = addDisableCommands(client, d, "key", addKeys.List())
			if err != nil {
				return diag.Errorf("error adding GeminiDB Redis disabled command: %s", err)
			}
		}
	}

	return resourceCommandDisableRead(ctx, d, meta)
}

func addDisableCommands(client *golangsdk.ServiceClient, d *schema.ResourceData, commandType string, commands []interface{}) error {
	addHttpUrl := "v3/{project_id}/redis/instances/{instance_id}/disabled-commands"
	addPath := client.Endpoint + addHttpUrl
	addPath = strings.ReplaceAll(addPath, "{project_id}", client.ProjectID)
	addPath = strings.ReplaceAll(addPath, "{instance_id}", d.Id())
	addIpsOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
		JSONBody:         buildUpdateCommandsBodyParams(commandType, commands),
	}

	_, err := client.Request("POST", addPath, &addIpsOpt)
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateCommandsBodyParams(commandType string, commands []interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"disabled_type": commandType,
	}

	if commandType == "command" {
		bodyParams["commands"] = utils.ExpandToStringList(commands)
	} else {
		bodyParams["keys"] = buildKeysInfoBodyParams(commands)
	}

	return bodyParams
}

func removeDisableCommands(client *golangsdk.ServiceClient, d *schema.ResourceData, commandType string, commands []interface{}) error {
	removeHttpUrl := "v3/{project_id}/redis/instances/{instance_id}/disabled-commands"
	removePath := client.Endpoint + removeHttpUrl
	removePath = strings.ReplaceAll(removePath, "{project_id}", client.ProjectID)
	removePath = strings.ReplaceAll(removePath, "{instance_id}", d.Id())
	removeOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
		JSONBody:         buildUpdateCommandsBodyParams(commandType, commands),
	}

	_, err := client.Request("DELETE", removePath, &removeOpt)
	if err != nil {
		return err
	}

	return nil
}

func resourceCommandDisableDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/redis/instances/{instance_id}/disabled-commands"
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateCommandDisableBodyParams(d)),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// whether the disabled command exist or not, the response HTTP status code of the deletion API is 200.
		return diag.Errorf("error deleting GeminiDB Redis disabled command: %s", err)
	}

	return nil
}

func resourceCommandDisableImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<disabled_type>', but got '%s'", importedId)
	}

	d.SetId(parts[0])

	mErr := multierror.Append(nil,
		d.Set("disabled_type", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
