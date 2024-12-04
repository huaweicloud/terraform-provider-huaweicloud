package live

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getLiveChannelResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/{project_id}/ott/channels"
		product = "live"
		id      = state.Primary.ID
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Live client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?id=%s", id)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Live channel: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	channelBody := utils.PathSearch("channels|[0]", respBody, nil)
	if channelBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return channelBody, nil
}

func TestAccLiveChannel_RTMP_PUSH(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_live_channel.test"

		dashPackageUrl = fmt.Sprintf("%s/rtmppush/dash/qqqw.mpd", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		hlsPackageUrl  = fmt.Sprintf("%s/rtmppush/hls/adsf.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		mssPackageUrl  = fmt.Sprintf("%s/rtmppush/mss/ffgdf.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)

		hlsPackageUpdateUrl1 = fmt.Sprintf("%s/rtmppush/hls/axe.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		mssPackageUpdateUrl2 = fmt.Sprintf("%s/rtmppush/mss/aax.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)

		hlsPackageUpdateUrl2 = fmt.Sprintf("%s/rtmppush/hls/kk.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLiveChannelResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveStreamingDomainName(t)
			acceptance.TestAccPreCheckLiveIngestRTMPDomainName(t)
			acceptance.TestAccPreCheckLiveTranscodingTemplateID(t)
			acceptance.TestAccPreCheckLiveTranscodingTemplateAnotherID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLiveChannel_RTMP_PUSH(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "name", "test-name"),
					resource.TestCheckResourceAttr(rName, "state", "OFF"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.0.template_id", acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.playlist_window_seconds", "12"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.segment_duration_seconds", "4"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.url", dashPackageUrl),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.encryption.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.encryption.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.encryption.0.level", "profile"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.encryption.0.request_mode", "direct_http"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.encryption.0.resource_id", "test-resource-id"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.encryption.0.speke_version", "1.0"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.encryption.0.url", "http://test-url.cp"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.encryption.0.http_headers.0.key", "key1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.encryption.0.http_headers.0.value", "value1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.encryption.0.http_headers.1.key", "key2"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.encryption.0.http_headers.1.value", "value2"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.request_args.0.record.0.end_time", "end"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.request_args.0.record.0.format", "timestamp"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.request_args.0.record.0.start_time", "begin"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.request_args.0.record.0.unit", "second"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.request_args.0.timeshift.0.back_time", "delay"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.request_args.0.timeshift.0.unit", "second"),

					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.hls_version", "v3"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.playlist_window_seconds", "36"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.segment_duration_seconds", "4"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.url", hlsPackageUrl),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.level", "profile"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.request_mode", "functiongraph_proxy"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.resource_id", "test-resource-id"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.speke_version", "1.0"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.url", "https://test-url.cc"),
					resource.TestCheckResourceAttrPair(rName, "endpoints.0.hls_package.0.encryption.0.urn", "huaweicloud_fgs_function.test", "urn"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.record.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.record.0.end_time", "end"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.record.0.format", "timestamp"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.record.0.start_time", "begin"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.record.0.unit", "second"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.timeshift.0.back_time", "delay"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.timeshift.0.unit", "second"),

					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.delay_segment", "0"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.playlist_window_seconds", "28"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.segment_duration_seconds", "4"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.url", mssPackageUrl),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.encryption.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.encryption.0.level", "content"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.encryption.0.request_mode", "functiongraph_proxy"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.encryption.0.resource_id", "test-resource-id"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.encryption.0.speke_version", "1.0"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.encryption.0.url", "https://test-url.cc"),
					resource.TestCheckResourceAttrPair(rName, "endpoints.0.mss_package.0.encryption.0.urn", "huaweicloud_fgs_function.test", "urn"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.record.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.record.0.end_time", "end"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.record.0.format", "timestamp"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.record.0.start_time", "begin"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.record.0.unit", "second"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.timeshift.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.timeshift.0.back_time", "delay"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.timeshift.0.unit", "second"),

					resource.TestCheckResourceAttr(rName, "input.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.input_protocol", "RTMP_PUSH"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.bitrate", "100"),
					resource.TestCheckResourceAttrSet(rName, "input.0.sources.0.url"),

					resource.TestCheckResourceAttr(rName, "record_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.0.rollingbuffer_duration", "12"),
					resource.TestCheckResourceAttrSet(rName, "channel_id"),
				),
			},
			{
				Config: testLiveChannel_RTMP_PUSH_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "name", "test-name-update"),
					resource.TestCheckResourceAttr(rName, "state", "OFF"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.0.template_id", acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID),
					resource.TestCheckResourceAttr(rName, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.hls_version", "v3"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.playlist_window_seconds", "78"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.segment_duration_seconds", "4"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.url", hlsPackageUpdateUrl1),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.level", "content"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.request_mode", "direct_http"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.resource_id", "test-resource-id"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.speke_version", "1.0"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.url", "https://test-url.cd"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.http_headers.0.key", "aa"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.http_headers.0.value", "dd"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.record.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.record.0.end_time", "end"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.record.0.format", "timestamp"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.record.0.start_time", "begin"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.record.0.unit", "second"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.timeshift.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.timeshift.0.back_time", "delay"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.timeshift.0.unit", "second"),

					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.delay_segment", "0"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.playlist_window_seconds", "42"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.segment_duration_seconds", "4"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.url", mssPackageUpdateUrl2),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.encryption.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.encryption.0.level", "content"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.encryption.0.request_mode", "functiongraph_proxy"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.encryption.0.resource_id", "12weq"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.encryption.0.speke_version", "1.0"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.encryption.0.url", "https://test-url-update.cc"),
					resource.TestCheckResourceAttrPair(rName, "endpoints.0.mss_package.0.encryption.0.urn", "huaweicloud_fgs_function.test", "urn"),

					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.record.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.record.0.end_time", "end"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.record.0.format", "timestamp"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.record.0.start_time", "begin"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.record.0.unit", "second"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.timeshift.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.timeshift.0.back_time", "delay"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.0.timeshift.0.unit", "second"),

					resource.TestCheckResourceAttr(rName, "input.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.input_protocol", "RTMP_PUSH"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.bitrate", "100"),
					resource.TestCheckResourceAttrSet(rName, "input.0.sources.0.url"),

					resource.TestCheckResourceAttr(rName, "record_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.0.rollingbuffer_duration", "20"),

					resource.TestCheckResourceAttrSet(rName, "channel_id"),
				),
			},
			{
				Config: testLiveChannel_RTMP_PUSH_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "state", "ON"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.0.template_id", acceptance.HW_LIVE_TRANSCODING_TEPLATE_ANOTHER_ID),
					resource.TestCheckResourceAttr(rName, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.hls_version", "v3"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.playlist_window_seconds", "40"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.segment_duration_seconds", "4"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.url", hlsPackageUpdateUrl2),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.record.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.0.timeshift.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.input_protocol", "RTMP_PUSH"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.bitrate", "100"),

					resource.TestCheckResourceAttr(rName, "record_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.0.rollingbuffer_duration", "0"),

					resource.TestCheckResourceAttrSet(rName, "input.0.sources.0.url"),
					resource.TestCheckResourceAttrSet(rName, "channel_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testLiveChannel_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  description = "function test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = <<EOT
aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldm
VudCkKICAgIHJldHVybiBvdXRwdXQ=
EOT
}
`, name)
}

func testLiveChannel_RTMP_PUSH(name string) string {
	var (
		ingestDomainName = acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME
		templateID       = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID
		dashPackageUrl   = fmt.Sprintf("%s/rtmppush/dash/qqqw.mpd", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		hlsPackageUrl    = fmt.Sprintf("%s/rtmppush/hls/adsf.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		mssPackageUrl    = fmt.Sprintf("%s/rtmppush/mss/ffgdf.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = "%[2]s"
  name        = "test-name"
  state       = "OFF"

  encoder_settings {
    template_id = "%[3]s"
  }

  endpoints {
    dash_package {
      playlist_window_seconds  = 12
      segment_duration_seconds = 4
      url                      = "%[4]s"

      encryption {
        level         = "profile"
        request_mode  = "direct_http"
        resource_id   = "test-resource-id"
        speke_version = "1.0"
        url           = "http://test-url.cp"

        system_ids = [
          "PlayReady",
        ]

        http_headers {
          key   = "key1"
          value = "value1"
        }
        http_headers {
          key   = "key2"
          value = "value2"
        }
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 36
      segment_duration_seconds = 4
      url                      = "%[5]s"

      encryption {
        level         = "profile"
        request_mode  = "functiongraph_proxy"
        resource_id   = "test-resource-id"
        speke_version = "1.0"
        url           = "https://test-url.cc"
        urn           = huaweicloud_fgs_function.test.urn

        system_ids = [
          "FairPlay",
        ]
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
    mss_package {
      delay_segment            = 0
      playlist_window_seconds  = 28
      segment_duration_seconds = 4
      url                      = "%[6]s"

      encryption {
        level         = "content"
        request_mode  = "functiongraph_proxy"
        resource_id   = "test-resource-id"
        speke_version = "1.0"
        url           = "https://test-url.cc"
        urn           = huaweicloud_fgs_function.test.urn

        system_ids = [
          "PlayReady",
        ]
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "RTMP_PUSH"

    sources {
      bitrate = 100
    }
  }

  record_settings {
    rollingbuffer_duration = 12
  }
}
`, testLiveChannel_base(name), ingestDomainName, templateID, dashPackageUrl, hlsPackageUrl, mssPackageUrl)
}

func testLiveChannel_RTMP_PUSH_update1(name string) string {
	var (
		ingestDomainName = acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME
		templateID       = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID
		hlsPackageUrl    = fmt.Sprintf("%s/rtmppush/hls/axe.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		mssPackageUrl    = fmt.Sprintf("%s/rtmppush/mss/aax.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = "%[2]s"
  name        = "test-name-update"
  state       = "OFF"

  encoder_settings {
    template_id = "%[3]s"
  }

  endpoints {
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 78
      segment_duration_seconds = 4
      url                      = "%[4]s"

      encryption {
        level         = "content"
        request_mode  = "direct_http"
        resource_id   = "test-resource-id"
        speke_version = "1.0"
        url           = "https://test-url.cd"

        system_ids = [
          "FairPlay",
        ]

        http_headers {
          key   = "aa"
          value = "dd"
        }
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }

    mss_package {
      delay_segment            = 0
      playlist_window_seconds  = 42
      segment_duration_seconds = 4
      url                      = "%[5]s"

      encryption {
        level         = "content"
        request_mode  = "functiongraph_proxy"
        resource_id   = "12weq"
        speke_version = "1.0"
        url           = "https://test-url-update.cc"
        urn           = huaweicloud_fgs_function.test.urn

        system_ids = [
          "PlayReady",
        ]
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "RTMP_PUSH"

    sources {
      bitrate = 100
    }
  }

  record_settings {
    rollingbuffer_duration = 20
  }
}
`, testLiveChannel_base(name), ingestDomainName, templateID, hlsPackageUrl, mssPackageUrl)
}

func testLiveChannel_RTMP_PUSH_update2(name string) string {
	var (
		ingestDomainName = acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME
		templateID       = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ANOTHER_ID
		hlsPackageUrl    = fmt.Sprintf("%s/rtmppush/hls/kk.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = "%[2]s"
  state       = "ON"

  encoder_settings {
    template_id = "%[3]s"
  }

  endpoints {
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 40
      segment_duration_seconds = 4
      url                      = "%[4]s"

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "RTMP_PUSH"

    sources {
      bitrate = 100
    }
  }

  record_settings {
    rollingbuffer_duration = 0
  }
}
`, testLiveChannel_base(name), ingestDomainName, templateID, hlsPackageUrl)
}

func TestAccLiveChannel_FLV_PULL(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_live_channel.test"

		dashPackageUrl       = fmt.Sprintf("%s/flvpull/dash/sa.mpd", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		hlsPackageUrl        = fmt.Sprintf("%s/flvpull/hls/ss.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		dashPackageUrlUpdate = fmt.Sprintf("%s/flvpull/dash/yy.mpd", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		hlsPackageUrlUpdate  = fmt.Sprintf("%s/flvpull/hls/ss.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		mssPackageUrlUpdate  = fmt.Sprintf("%s/flvpull/mss/ss.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLiveChannelResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveStreamingDomainName(t)
			acceptance.TestAccPreCheckLiveIngestRTMPDomainName(t)
			acceptance.TestAccPreCheckLiveTranscodingTemplateID(t)
			acceptance.TestAccPreCheckLiveTranscodingTemplateAnotherID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLiveChannel_FLV_PULL(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "name", "test-name"),
					resource.TestCheckResourceAttr(rName, "state", "ON"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.0.template_id", "test-template"),
					resource.TestCheckResourceAttr(rName, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.playlist_window_seconds", "12"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.segment_duration_seconds", "2"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.url", dashPackageUrl),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.request_args.#", "1"),

					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.hls_version", "v3"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.playlist_window_seconds", "4"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.segment_duration_seconds", "2"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.url", hlsPackageUrl),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.#", "1"),

					resource.TestCheckResourceAttr(rName, "input.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.input_protocol", "FLV_PULL"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.0.input_loss_threshold_msec", "3000"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.0.input_preference", "PRIMARY"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.0.bitrate", "100"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.0.url", "https://sss.vv"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.bitrate", "100"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.url", "https://sss.cc"),

					resource.TestCheckResourceAttr(rName, "record_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.0.rollingbuffer_duration", "10"),

					resource.TestCheckResourceAttrSet(rName, "channel_id"),
				),
			},
			{
				Config: testLiveChannel_FLV_PULL_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "name", "test-name-update"),
					resource.TestCheckResourceAttr(rName, "state", "ON"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.#", "2"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.0.template_id", acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID),
					resource.TestCheckResourceAttr(rName, "encoder_settings.1.template_id", acceptance.HW_LIVE_TRANSCODING_TEPLATE_ANOTHER_ID),
					resource.TestCheckResourceAttr(rName, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.playlist_window_seconds", "24"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.segment_duration_seconds", "4"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.url", dashPackageUrlUpdate),
					resource.TestCheckResourceAttr(rName, "endpoints.0.dash_package.0.request_args.#", "1"),

					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.hls_version", "v3"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.playlist_window_seconds", "8"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.segment_duration_seconds", "4"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.url", hlsPackageUrlUpdate),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.#", "1"),

					resource.TestCheckResourceAttr(rName, "input.0.input_protocol", "FLV_PULL"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.0.input_loss_threshold_msec", "4000"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.0.input_preference", "EQUAL"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.0.bitrate", "100"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.0.url", "http://hgf.vv"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.bitrate", "100"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.url", "http://qwe.cc"),

					resource.TestCheckResourceAttr(rName, "record_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.0.rollingbuffer_duration", "12"),

					resource.TestCheckResourceAttrSet(rName, "channel_id"),
				),
			},
			{
				Config: testLiveChannel_FLV_PULL_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "name", ""),
					resource.TestCheckResourceAttr(rName, "state", "OFF"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.0.template_id", acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID),
					resource.TestCheckResourceAttr(rName, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.delay_segment", "0"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.playlist_window_seconds", "54"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.segment_duration_seconds", "4"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.url", mssPackageUrlUpdate),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.input_protocol", "FLV_PULL"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.0.input_loss_threshold_msec", "4000"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.0.input_preference", "EQUAL"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.bitrate", "100"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.url", "http://ssd.cc"),

					resource.TestCheckResourceAttr(rName, "record_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.0.rollingbuffer_duration", "0"),

					resource.TestCheckResourceAttrSet(rName, "channel_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testLiveChannel_FLV_PULL(name string) string {
	var (
		ingestDomainName = acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME
		templateID       = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID
		dashPackageUrl   = fmt.Sprintf("%s/flvpull/dash/sa.mpd", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		hlsPackageUrl    = fmt.Sprintf("%s/flvpull/hls/ss.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = "%[2]s"
  name        = "test-name"
  state       = "ON"

  encoder_settings {
    template_id = "%[3]s"
  }

  endpoints {
    dash_package {
      playlist_window_seconds  = 12
      segment_duration_seconds = 2
      url                      = "%[4]s"

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 4
      segment_duration_seconds = 2
      url                      = "%[5]s"

      encryption {
        level         = "content"
        request_mode  = "direct_http"
        resource_id   = "qqq"
        speke_version = "1.0"
        url           = "http://sssddd.vf"

        system_ids = [
          "FairPlay",
        ]

        http_headers {
          key   = "aaa"
          value = "sss"
        }
        http_headers {
          key   = "www"
          value = "qqq"
        }
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "FLV_PULL"

    failover_conditions {
      input_loss_threshold_msec = 3000
      input_preference          = "PRIMARY"
    }

    secondary_sources {
      bitrate = 100
      url     = "https://sss.vv"
    }

    sources {
      bitrate = 100
      url     = "https://sss.cc"
    }
  }

  record_settings {
    rollingbuffer_duration = 10
  }
}
`, testLiveChannel_base(name), ingestDomainName, templateID, dashPackageUrl, hlsPackageUrl)
}

func testLiveChannel_FLV_PULL_update1(name string) string {
	var (
		ingestDomainName  = acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME
		templateID        = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID
		templateAnotherID = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ANOTHER_ID
		dashPackageUrl    = fmt.Sprintf("%s/flvpull/dash/yy.mpd", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		hlsPackageUrl     = fmt.Sprintf("%s/flvpull/hls/ss.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = "%[2]s"
  name        = "test-name-update"
  state       = "ON"

  encoder_settings {
    template_id = "%[3]s"
  }
  encoder_settings {
    template_id = "%[4]s"
  }


  endpoints {
    dash_package {
      playlist_window_seconds  = 24
      segment_duration_seconds = 4
      url                      = "%[5]s"

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }

    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 8
      segment_duration_seconds = 4
      url                      = "%[6]s"

      encryption {
        level         = "content"
        request_mode  = "direct_http"
        resource_id   = "qqq"
        speke_version = "1.0"
        url           = "http://xxx.sp"

        system_ids = [
          "FairPlay",
        ]

        http_headers {
          key   = "aaa"
          value = "sss"
        }
        http_headers {
          key   = "www"
          value = "qqq"
        }
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "FLV_PULL"

    failover_conditions {
      input_loss_threshold_msec = 4000
      input_preference          = "EQUAL"
    }

    secondary_sources {
      bitrate = 100
      url     = "http://hgf.vv"
    }

    sources {
      bitrate = 100
      url     = "http://qwe.cc"
    }
  }

  record_settings {
    rollingbuffer_duration = 12
  }
}
`, testLiveChannel_base(name), ingestDomainName, templateID, templateAnotherID, dashPackageUrl, hlsPackageUrl)
}

func testLiveChannel_FLV_PULL_update2(name string) string {
	var (
		ingestDomainName = acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME
		templateID       = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID
		mssPackageUrl    = fmt.Sprintf("%s/flvpull/mss/ss.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = "%[2]s"
  state       = "OFF"

  encoder_settings {
    template_id = "%[3]s"
  }

  endpoints {
    mss_package {
      playlist_window_seconds  = 54
      segment_duration_seconds = 4
      url                      = "%[4]s"

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "FLV_PULL"

    sources {
      bitrate = 100
      url     = "http://ssd.cc"
    }
  }

  record_settings {
    rollingbuffer_duration = 0
  }
}
`, testLiveChannel_base(name), ingestDomainName, templateID, mssPackageUrl)
}

func TestAccLiveChannel_HLS_PULL(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_live_channel.test"

		hlsPackageUrl       = fmt.Sprintf("%s/hlspull/hls/ss.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		mssPackageUrl       = fmt.Sprintf("%s/hlspull/mss/ure.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		hlsPackageUrlUpdate = fmt.Sprintf("%s/hlspull/hls/ss.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLiveChannelResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveStreamingDomainName(t)
			acceptance.TestAccPreCheckLiveIngestRTMPDomainName(t)
			acceptance.TestAccPreCheckLiveTranscodingTemplateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLiveChannel_HLS_PULL(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "name", "test-name"),
					resource.TestCheckResourceAttr(rName, "state", "OFF"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.0.template_id", acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID),
					resource.TestCheckResourceAttr(rName, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.hls_version", "v3"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.playlist_window_seconds", "4"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.segment_duration_seconds", "2"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.url", hlsPackageUrl),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.level", "content"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.request_mode", "functiongraph_proxy"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.resource_id", "qwerwq"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.speke_version", "1.0"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.0.url", "https://test-url.sd"),
					resource.TestCheckResourceAttrPair(rName, "endpoints.0.hls_package.0.encryption.0.urn", "huaweicloud_fgs_function.test", "urn"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.#", "1"),

					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.delay_segment", "0"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.playlist_window_seconds", "8"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.segment_duration_seconds", "2"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.url", mssPackageUrl),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.encryption.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.#", "1"),

					resource.TestCheckResourceAttr(rName, "input.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.input_protocol", "HLS_PULL"),
					resource.TestCheckResourceAttr(rName, "input.0.max_bandwidth_limit", "200"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.0.input_loss_threshold_msec", "2000"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.0.input_preference", "PRIMARY"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.0.bitrate", "100"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.0.url", "http://qqwe.dd"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.bitrate", "100"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.url", "http://ssa.qw"),

					resource.TestCheckResourceAttr(rName, "record_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.0.rollingbuffer_duration", "3"),

					resource.TestCheckResourceAttrSet(rName, "channel_id"),
				),
			},
			{
				Config: testLiveChannel_HLS_PULL_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "name", "test-name-update"),
					resource.TestCheckResourceAttr(rName, "state", "ON"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.0.template_id", acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID),
					resource.TestCheckResourceAttr(rName, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.hls_version", "v3"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.playlist_window_seconds", "12"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.segment_duration_seconds", "2"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.url", hlsPackageUrlUpdate),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.input_protocol", "HLS_PULL"),
					resource.TestCheckResourceAttr(rName, "input.0.max_bandwidth_limit", "300"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.0.input_loss_threshold_msec", "3000"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.0.input_preference", "EQUAL"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.0.bitrate", "100"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.0.url", "http://vvgj.dd"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.bitrate", "100"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.url", "http://eeds.qw"),

					resource.TestCheckResourceAttr(rName, "record_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.0.rollingbuffer_duration", "9"),

					resource.TestCheckResourceAttrSet(rName, "channel_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testLiveChannel_HLS_PULL(name string) string {
	var (
		ingestDomainName = acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME
		templateID       = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID
		hlsPackageUrl    = fmt.Sprintf("%s/hlspull/hls/ss.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		mssPackageUrl    = fmt.Sprintf("%s/hlspull/mss/ure.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = "%[2]s"
  name        = "test-name"
  state       = "OFF"

  encoder_settings {
    template_id = "%[3]s"
  }

  endpoints {
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 4
      segment_duration_seconds = 2
      url                      = "%[4]s"

      encryption {
        level         = "content"
        request_mode  = "functiongraph_proxy"
        resource_id   = "qwerwq"
        speke_version = "1.0"
        url           = "https://test-url.sd"
        urn           = huaweicloud_fgs_function.test.urn

        system_ids = [
          "FairPlay",
        ]
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
    mss_package {
      playlist_window_seconds  = 8
      segment_duration_seconds = 2
      url                      = "%[5]s"

      encryption {
        level         = "content"
        request_mode  = "direct_http"
        resource_id   = "dfge"
        speke_version = "1.0"
        url           = "https://ssc.cd"

        system_ids = [
          "PlayReady",
        ]

        http_headers {
          key   = "aa"
          value = "ss"
        }
        http_headers {
          key   = "gg"
          value = "ff"
        }
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol      = "HLS_PULL"
    max_bandwidth_limit = 200

    failover_conditions {
      input_loss_threshold_msec = 2000
      input_preference          = "PRIMARY"
    }

    secondary_sources {
      bitrate = 100
      url     = "http://qqwe.dd"
    }

    sources {
      bitrate = 100
      url     = "http://ssa.qw"
    }
  }

  record_settings {
    rollingbuffer_duration = 3
  }
}
`, testLiveChannel_base(name), ingestDomainName, templateID, hlsPackageUrl, mssPackageUrl)
}

func testLiveChannel_HLS_PULL_update(name string) string {
	var (
		ingestDomainName = acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME
		templateID       = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID
		hlsPackageUrl    = fmt.Sprintf("%s/hlspull/hls/ss.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = "%[2]s"
  name        = "test-name-update"
  state       = "ON"

  encoder_settings {
    template_id = "%[3]s"
  }

  endpoints {
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 12
      segment_duration_seconds = 2
      url                      = "%[4]s"

      encryption {
        level         = "content"
        request_mode  = "functiongraph_proxy"
        resource_id   = "qwerwq"
        speke_version = "1.0"
        url           = "https://test-url.fg"
        urn           = huaweicloud_fgs_function.test.urn

        system_ids = [
          "FairPlay",
        ]
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol      = "HLS_PULL"
    max_bandwidth_limit = 300

    failover_conditions {
      input_loss_threshold_msec = 3000
      input_preference          = "EQUAL"
    }

    secondary_sources {
      bitrate = 100
      url     = "http://vvgj.dd"
    }

    sources {
      bitrate = 100
      url     = "http://eeds.qw"
    }
  }

  record_settings {
    rollingbuffer_duration = 9
  }
}
`, testLiveChannel_base(name), ingestDomainName, templateID, hlsPackageUrl)
}

func TestAccLiveChannel_SRT_PUSH(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_live_channel.test"

		hlsPackageUrl       = fmt.Sprintf("%s/srtpush/hls/aa.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		mssPackageUrl       = fmt.Sprintf("%s/srtpush/mss/ggf.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		mssPackageUrlUpdate = fmt.Sprintf("%s/srtpush/mss/ggf.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLiveChannelResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveStreamingDomainName(t)
			acceptance.TestAccPreCheckLiveIngestSRTDomainName(t)
			acceptance.TestAccPreCheckLiveTranscodingTemplateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLiveChannel_SRT_PUSH(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_SRT_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "name", "test-name"),
					resource.TestCheckResourceAttr(rName, "state", "OFF"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.0.template_id", acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID),
					resource.TestCheckResourceAttr(rName, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.url", hlsPackageUrl),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.#", "1"),

					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.url", mssPackageUrl),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.#", "1"),

					resource.TestCheckResourceAttr(rName, "input.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.input_protocol", "SRT_PUSH"),
					resource.TestCheckResourceAttr(rName, "input.0.ip_whitelist", "192.168.0.1/16,192.168.1.1/16,192.168.2.1/16"),
					resource.TestCheckResourceAttr(rName, "input.0.audio_selectors.#", "3"),

					resource.TestCheckResourceAttr(rName, "input.0.sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.0.rollingbuffer_duration", "2"),

					resource.TestCheckResourceAttrSet(rName, "channel_id"),
				),
			},
			{
				Config: testLiveChannel_SRT_PUSH_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_SRT_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "name", "test-name-update"),
					resource.TestCheckResourceAttr(rName, "state", "ON"),
					resource.TestCheckResourceAttr(rName, "encoder_settings_expand.#", "1"),
					resource.TestCheckResourceAttr(rName, "encoder_settings_expand.0.audio_descriptions.#", "2"),
					resource.TestCheckResourceAttr(rName, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.url", mssPackageUrlUpdate),
					resource.TestCheckResourceAttr(rName, "input.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.input_protocol", "SRT_PUSH"),
					resource.TestCheckResourceAttr(rName, "input.0.ip_whitelist", "192.168.0.1/16,192.168.1.1/16"),
					resource.TestCheckResourceAttr(rName, "input.0.audio_selectors.#", "2"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.0.rollingbuffer_duration", "4"),

					resource.TestCheckResourceAttrSet(rName, "channel_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testLiveChannel_SRT_PUSH(name string) string {
	var (
		ingestDomainName = acceptance.HW_LIVE_INGEST_SRT_DOMAIN_NAME
		templateID       = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID
		hlsPackageUrl    = fmt.Sprintf("%s/srtpush/hls/aa.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		mssPackageUrl    = fmt.Sprintf("%s/srtpush/mss/ggf.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = "%[2]s"
  name        = "test-name"
  state       = "OFF"

  encoder_settings {
    template_id = "%[3]s"
  }

  endpoints {
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 4
      segment_duration_seconds = 2
      url                      = "%[4]s"

      encryption {
        encryption_method             = "SAMPLE-AES"
        key_rotation_interval_seconds = 0
        level                         = "content"
        request_mode                  = "direct_http"
        resource_id                   = "asdfwee"
        speke_version                 = "1.0"
        url                           = "http://qqq.co"

        system_ids = [
          "FairPlay",
        ]

        http_headers {
          key   = "aa"
          value = "sss"
        }
        http_headers {
          key   = "dd"
          value = "sss"
        }
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
    mss_package {
      playlist_window_seconds  = 8
      segment_duration_seconds = 2
      url                      = "%[5]s"

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "SRT_PUSH"
    ip_whitelist   = "192.168.0.1/16,192.168.1.1/16,192.168.2.1/16"

    audio_selectors {
      name = "qweweqrd"

      selector_settings {
        audio_pid_selection {
          pid = 2
        }
      }
    }
    audio_selectors {
      name = "ferfddfag"

      selector_settings {
        audio_language_selection {
          language_code             = "ch"
          language_selection_policy = "LOOSE"
        }
      }
    }
    audio_selectors {
      name = "wqwerwed"

      selector_settings {
        audio_pid_selection {
          pid = 0
        }
      }
    }

    sources {
      bitrate = 100
    }
  }

  record_settings {
    rollingbuffer_duration = 2
  }
}
`, testLiveChannel_base(name), ingestDomainName, templateID, hlsPackageUrl, mssPackageUrl)
}

func testLiveChannel_SRT_PUSH_update(name string) string {
	var (
		ingestDomainName = acceptance.HW_LIVE_INGEST_SRT_DOMAIN_NAME
		templateID       = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID
		mssPackageUrl    = fmt.Sprintf("%s/srtpush/mss/ggf.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = "%[2]s"
  name        = "test-name-update"
  state       = "ON"

  encoder_settings {
    template_id = "%[3]s"
  }

  encoder_settings_expand {
    audio_descriptions {
      audio_selector_name   = "verfgger"
      language_code         = "ch"
      language_code_control = "FOLLOW_INPUT"
      name                  = "fgdfvs"
      stream_name           = "hhee"
    }
    audio_descriptions {
      audio_selector_name   = "gfdsge"
      language_code         = "en"
      language_code_control = "USE_CONFIGURED"
      name                  = "vsdfg5e"
      stream_name           = "ssaw"
    }
  }

  endpoints {
    mss_package {
      playlist_window_seconds  = 24
      segment_duration_seconds = 2
      url                      = "%[4]s"

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "SRT_PUSH"
    ip_whitelist   = "192.168.0.1/16,192.168.1.1/16"

    audio_selectors {
      name = "verfgger"

      selector_settings {
        audio_pid_selection {
          pid = 5
        }
      }
    }
    audio_selectors {
      name = "gfdsge"

      selector_settings {
        audio_language_selection {
          language_code             = "ee"
          language_selection_policy = "STRICT"
        }
      }
    }

    sources {
      bitrate = 100
    }
  }

  record_settings {
    rollingbuffer_duration = 4
  }
}
`, testLiveChannel_base(name), ingestDomainName, templateID, mssPackageUrl)
}

func TestAccLiveChannel_SRT_PULL(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_live_channel.test"

		hlsPackageUrl       = fmt.Sprintf("%s/srtpull/hls/ss.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		mssPackageUrl       = fmt.Sprintf("%s/srtpull/mss/ss.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		mssPackageUrlUpdate = fmt.Sprintf("%s/srtpull/mss/ee.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLiveChannelResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveStreamingDomainName(t)
			acceptance.TestAccPreCheckLiveIngestSRTDomainName(t)
			acceptance.TestAccPreCheckLiveTranscodingTemplateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLiveChannel_SRT_PULL(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_SRT_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "name", "test-name"),
					resource.TestCheckResourceAttr(rName, "state", "OFF"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.0.template_id", acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID),
					resource.TestCheckResourceAttr(rName, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.url", hlsPackageUrl),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.encryption.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.url", mssPackageUrl),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.#", "1"),

					resource.TestCheckResourceAttr(rName, "input.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.input_protocol", "SRT_PULL"),
					resource.TestCheckResourceAttr(rName, "input.0.audio_selectors.#", "3"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.#", "1"),

					resource.TestCheckResourceAttr(rName, "record_settings.#", "1"),

					resource.TestCheckResourceAttrSet(rName, "channel_id"),
				),
			},
			{
				Config: testLiveChannel_SRT_PULL_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_SRT_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "name", "dfgr"),
					resource.TestCheckResourceAttr(rName, "state", "ON"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.0.template_id", acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID),

					resource.TestCheckResourceAttr(rName, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.url", mssPackageUrlUpdate),
					resource.TestCheckResourceAttr(rName, "endpoints.0.mss_package.0.request_args.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.input_protocol", "SRT_PULL"),
					resource.TestCheckResourceAttr(rName, "input.0.audio_selectors.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.failover_conditions.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.secondary_sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.#", "1"),

					resource.TestCheckResourceAttr(rName, "record_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.0.rollingbuffer_duration", "6"),

					resource.TestCheckResourceAttrSet(rName, "channel_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testLiveChannel_SRT_PULL(name string) string {
	var (
		ingestDomainName = acceptance.HW_LIVE_INGEST_SRT_DOMAIN_NAME
		templateID       = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID
		hlsPackageUrl    = fmt.Sprintf("%s/srtpull/hls/ss.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
		mssPackageUrl    = fmt.Sprintf("%s/srtpull/mss/ss.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = "%[2]s"
  name        = "test-name"
  state       = "OFF"

  encoder_settings {
    template_id = "%[3]s"
  }

  endpoints {
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 16
      segment_duration_seconds = 4
      url                      = "%[4]s"

      encryption {
        encryption_method = "SAMPLE-AES"
        level             = "content"
        request_mode      = "direct_http"
        resource_id       = "vdfgdfg"
        speke_version     = "1.0"
        url               = "http://sss.cc"

        system_ids = [
          "FairPlay",
        ]

        http_headers {
          key   = "aa"
          value = "ss"
        }
        http_headers {
          key   = "ff"
          value = "dd"
        }
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
    mss_package {
      playlist_window_seconds  = 24
      segment_duration_seconds = 4
      url                      = "%[5]s"

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "SRT_PULL"

    audio_selectors {
      name = "dvcvbdfg"

      selector_settings {
        audio_language_selection {
          language_code             = "dfg"
          language_selection_policy = "LOOSE"
        }
      }
    }
    audio_selectors {
      name = "fgvsdfre"

      selector_settings {
        audio_pid_selection {
          pid = 13
        }
      }
    }
    audio_selectors {
      name = "vbfgh5t"

      selector_settings {
        audio_pid_selection {
          pid = 0
        }
      }
    }

    failover_conditions {
      input_loss_threshold_msec = 2000
      input_preference          = "EQUAL"
    }

    secondary_sources {
      bitrate   = 100
      latency   = 1000
      stream_id = "vcbeer"
      url       = "srt://192.168.1.215:9001"
    }

    sources {
      bitrate   = 100
      latency   = 2000
      stream_id = "dfawerw"
      url       = "srt://192.168.1.216:9001"
    }
  }

  record_settings {
    rollingbuffer_duration = 4
  }
}
`, testLiveChannel_base(name), ingestDomainName, templateID, hlsPackageUrl, mssPackageUrl)
}

func testLiveChannel_SRT_PULL_update(name string) string {
	var (
		ingestDomainName = acceptance.HW_LIVE_INGEST_SRT_DOMAIN_NAME
		templateID       = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID
		mssPackageUrl    = fmt.Sprintf("%s/srtpull/mss/ee.ism/Manifest", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = "%[2]s"
  name        = "dfgr"
  state       = "ON"

  encoder_settings {
    template_id = "%[3]s"
  }

  endpoints {
    mss_package {
      playlist_window_seconds  = 42
      segment_duration_seconds = 4
      url                      = "%[4]s"

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }
        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "SRT_PULL"

    audio_selectors {
      name = "vbrt"

      selector_settings {
        audio_pid_selection {
          pid = 16
        }
      }
    }

    failover_conditions {
      input_loss_threshold_msec = 3000
      input_preference          = "PRIMARY"
    }

    secondary_sources {
      bitrate   = 100
      latency   = 2000
      stream_id = "cvbf"
      url       = "srt://192.168.1.215:9006"
    }

    sources {
      bitrate   = 100
      latency   = 3000
      stream_id = "dfgv"
      url       = "srt://192.168.1.216:9005"
    }
  }

  record_settings {
    rollingbuffer_duration = 6
  }
}
`, testLiveChannel_base(name), ingestDomainName, templateID, mssPackageUrl)
}

func TestAccLiveChannel_CUSTOM_CHANNEL_ID(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_live_channel.test"

		channelID, _  = uuid.GenerateUUID()
		hlsPackageUrl = fmt.Sprintf("%s/customchannel/hls/ss.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getLiveChannelResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveStreamingDomainName(t)
			acceptance.TestAccPreCheckLiveIngestRTMPDomainName(t)
			acceptance.TestAccPreCheckLiveTranscodingTemplateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testLiveChannel_CUSTOM_CHANNEL(name, channelID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "app_name", "live"),
					resource.TestCheckResourceAttr(rName, "channel_id", channelID),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME),
					resource.TestCheckResourceAttr(rName, "name", ""),
					resource.TestCheckResourceAttr(rName, "state", "OFF"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "encoder_settings.0.template_id", acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID),
					resource.TestCheckResourceAttr(rName, "endpoints.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.#", "1"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.hls_version", "v3"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.playlist_window_seconds", "40"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.segment_duration_seconds", "4"),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.url", hlsPackageUrl),
					resource.TestCheckResourceAttr(rName, "endpoints.0.hls_package.0.request_args.#", "1"),

					resource.TestCheckResourceAttr(rName, "input.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.input_protocol", "RTMP_PUSH"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.#", "1"),
					resource.TestCheckResourceAttr(rName, "input.0.sources.0.bitrate", "100"),
					resource.TestCheckResourceAttrSet(rName, "input.0.sources.0.url"),

					resource.TestCheckResourceAttr(rName, "record_settings.#", "1"),
					resource.TestCheckResourceAttr(rName, "record_settings.0.rollingbuffer_duration", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testLiveChannel_CUSTOM_CHANNEL(name, channelID string) string {
	var (
		ingestDomainName = acceptance.HW_LIVE_INGEST_RTMP_DOMAIN_NAME
		templateID       = acceptance.HW_LIVE_TRANSCODING_TEPLATE_ID
		hlsPackageUrl    = fmt.Sprintf("%s/customchannel/hls/ss.m3u8", acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
	)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = "%[2]s"
  state       = "OFF"
  channel_id  = "%[3]s"

  encoder_settings {
    template_id = "%[4]s"
  }

  endpoints {
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 40
      segment_duration_seconds = 4
      url                      = "%[5]s"

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }

        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "RTMP_PUSH"

    sources {
      bitrate = 100
    }
  }

  record_settings {
    rollingbuffer_duration = 0
  }
}
`, testLiveChannel_base(name), ingestDomainName, channelID, templateID, hlsPackageUrl)
}
