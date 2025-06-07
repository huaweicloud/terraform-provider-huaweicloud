package dew

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DecryptPasswordWithDefaultKmsKey(_ context.Context, meta interface{}, d *schema.ResourceData, cipherText string) (string, error) {
	var (
		cfg                = meta.(*config.Config)
		region             = cfg.GetRegion(d)
		dataDecryptHttpUrl = "v1.0/{project_id}/kms/decrypt-data"
		dataDecryptProduct = "kms"
	)

	client, err := cfg.NewServiceClient(dataDecryptProduct, region)
	if err != nil {
		return "", fmt.Errorf("error creating KMS client: %s", err)
	}

	dataDecryptPath := client.Endpoint + dataDecryptHttpUrl
	dataDecryptPath = strings.ReplaceAll(dataDecryptPath, "{project_id}", client.ProjectID)

	dataDecryptOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	bodyParams := map[string]interface{}{
		"cipher_text": cipherText,
	}

	dataDecryptOpt.JSONBody = utils.RemoveNil(bodyParams)
	dataDecryptResp, err := client.Request("POST", dataDecryptPath, &dataDecryptOpt)
	if err != nil {
		return "", fmt.Errorf("error decrypting data with KMS service: %s", err)
	}

	dataDecryptRespBody, err := utils.FlattenResponse(dataDecryptResp)
	if err != nil {
		return "", fmt.Errorf("error flatting response: %s", err)
	}

	plainText := utils.PathSearch("plain_text", dataDecryptRespBody, "").(string)
	if plainText == "" {
		return "", errors.New("unable to find the plain text from the API response")
	}

	return plainText, nil
}
