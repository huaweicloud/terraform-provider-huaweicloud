package helper

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/impl"
)

func loadCredentialFromMetadata(client *impl.DefaultHttpClient, defaultType string) (auth.ICredential, error) {
	if BasicCredentialType == defaultType {
		cred := basic.NewCredentialsBuilder().Build()
		err := cred.UpdateCredential(client)
		if err != nil {
			return nil, err
		}
		return cred, nil
	} else if GlobalCredentialType == defaultType {
		cred := global.NewCredentialsBuilder().Build()
		err := cred.UpdateCredential(client)
		if err != nil {
			return nil, err
		}
		return cred, nil
	}

	return nil, nil
}

func ProcessCredential(client *impl.DefaultHttpClient, defaultType string, cred auth.ICredential) (auth.ICredential, error) {
	if cred == nil {
		return loadCredentialFromMetadata(client, defaultType)
	}

	if basicCred, ok := cred.(basic.Credentials); ok && basicCred.NeedUpdate() {
		err := basicCred.UpdateCredential(client)
		if err != nil {
			return nil, err
		}
	} else if globalCred, ok := cred.(global.Credentials); ok && globalCred.NeedUpdate() {
		err := globalCred.UpdateCredential(client)
		if err != nil {
			return nil, err
		}
	}

	return cred, nil
}
