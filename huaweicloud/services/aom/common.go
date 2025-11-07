package aom

func buildMoreHeaders(epsId string) map[string]string {
	moreHeaders := map[string]string{
		"Content-Type": "application/json",
	}
	if epsId != "" {
		moreHeaders["Enterprise-Project-Id"] = epsId
	}
	return moreHeaders
}
