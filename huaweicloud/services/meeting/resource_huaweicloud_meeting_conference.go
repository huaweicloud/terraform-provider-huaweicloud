package meeting

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/meeting/v1/conferences"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type RecordMode int
type EncryptMode int
type RecordType int
type RecordAuthType int
type RoleType int
type MuteType int
type InviteType int
type CallInType int

type MediaType string
type LanguageType string
type ConferenceStatus string
type AccessType string
type CycleType string
type ConferenceRole string

const (
	// Standard format for timestamp conversion.
	standardTimeFormat string = "2006-01-02 15:04"

	// Whether to record conference.
	RecordModeOff RecordMode = 0
	RecordModeOn  RecordMode = 1

	// Media encryption mode.
	EncryptModeAdapt EncryptMode = 0
	EncryptModeForce EncryptMode = 1
	EncryptModeOff   EncryptMode = 2

	// Record type of the conference.
	RecordTypeOff    RecordType = 0
	RecordTypeLive   RecordType = 1
	RecordTypeRc     RecordType = 2
	RecordTypeLiveRc RecordType = 3

	// Authorization type of recording.
	RecordAuthTypeLink       RecordAuthType = 0
	RecordAuthTypeEnterprise RecordAuthType = 1
	RecordAuthTypeAttendee   RecordAuthType = 2

	// Participant role type.
	RoleTypeGuest RoleType = 0
	RoleTypeHost  RoleType = 1

	// Whether to mute the participant after joining.
	MuteTypeOff MuteType = 0
	MuteTypeOn  MuteType = 1

	// Invitation type for participants.
	InviteTypeManual InviteType = 0
	InviteTypeAuto   InviteType = 1

	// Which user can be called.
	CallInTypeAll        CallInType = 0
	CallInTypeEnterprise CallInType = 2
	CallInTypeInvitee    CallInType = 3

	// The media type of the conference.
	MediaTypeVoice   MediaType = "Voice"
	MediaTypeVideo   MediaType = "Video"
	MediaTypeHDVideo MediaType = "HDVideo"
	MediaTypeData    MediaType = "Data"

	// Default language type of the conference.
	LanguageTypeZhCn LanguageType = "zh-CN"
	LanguageTypeEnUs LanguageType = "en-US"

	// Conference status.
	ConferenceStatusSchedule ConferenceStatus = "Schedule"
	ConferenceStatusUsing    ConferenceStatus = "Created"

	// The attendee access type.
	AccessTypeNormal       AccessType = "normal"
	AccessTypeTelepresence AccessType = "telepresence"
	AccessTypeTerminal     AccessType = "terminal"
	AccessTypeOutside      AccessType = "outside"
	AccessTypeMobile       AccessType = "mobile"
	AccessTypeTelephone    AccessType = "telephone"
	AccessTypeIdeahub      AccessType = "ideahub"

	// Cycle type of the cyclical conference.
	CycleTypeDay   CycleType = "Day"
	CycleTypeWeek  CycleType = "Week"
	CycleTypeMonth CycleType = "Month"

	// Participant role title.
	ConferenceRoleChair ConferenceRole = "chair"
	ConferenceRoleGuest ConferenceRole = "general"
)

func subconferenceSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"media_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_auto_record": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"record_auth_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subconfiguration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"callin_restriction": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"audience_callin_restriction": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"allow_guest_start": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"waiting_room_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"show_audience_policies": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"multiple": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"base_audience_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// @API Meeting GET /v1/mmc/management/conferences/history/confDetail
// @API Meeting GET /v1/mmc/management/conferences/history
// @API Meeting DELETE /v1/mmc/management/conferences
// @API Meeting GET /v1/mmc/management/conferences
// @API Meeting POST /v1/mmc/management/conferences
// @API Meeting PUT /v1/mmc/management/conferences
// @API Meeting PUT /v1/mmc/management/conference/duration
// @API Meeting GET /v1/mmc/management/conferences/confDetail
// @API Meeting POST /v1/usg/acs/token/validate
// @API Meeting POST /v1/usg/acs/auth/account
// @API Meeting POST /v2/usg/acs/auth/appauth
// @API Meeting POST /v1/mmc/management/cycleconferences
// @API Meeting PUT /v1/mmc/management/cycleconferences
// @API Meeting DELETE /v1/mmc/management/cycleconferences
func ResourceConference() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConferenceCreate,
		ReadContext:   resourceConferenceRead,
		UpdateContext: resourceConferenceUpdate,
		DeleteContext: resourceConferenceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceConferenceImportState,
		},

		Schema: map[string]*schema.Schema{
			// Authorization arguments
			"account_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"account_password"},
			},
			"account_password": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"app_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"app_key"},
				ExactlyOneOf: []string{"account_name"},
			},
			"app_key": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"corp_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			// Arguments
			"topic": {
				Type:     schema.TypeString,
				Required: true,
			},
			"meeting_room_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$`),
					"The start time format should be 'YYYY-MM-DD hh:mm'."),
			},
			"media_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(MediaTypeVoice), string(MediaTypeVideo), string(MediaTypeHDVideo),
						string(MediaTypeData),
					}, false),
				},
			},
			"is_auto_record": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.IntInSlice([]int{
					int(RecordModeOff), int(RecordModeOn),
				}),
			},
			"encrypt_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.IntInSlice([]int{
					int(EncryptModeAdapt), int(EncryptModeForce), int(EncryptModeOff),
				}),
			},
			"language": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(LanguageTypeZhCn), string(LanguageTypeEnUs),
				}, false),
			},
			"timezone_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"record_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.IntInSlice([]int{
					int(RecordTypeOff), int(RecordTypeLive), int(RecordTypeRc), int(RecordTypeLiveRc),
				}),
			},
			"live_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"aux_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_record_aux_stream": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.IntInSlice([]int{
					int(RecordModeOff), int(RecordModeOn),
				}),
			},
			"record_auth_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.IntInSlice([]int{
					int(RecordAuthTypeLink), int(RecordAuthTypeEnterprise), int(RecordAuthTypeAttendee),
				}),
			},
			"participant_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"participant": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.IntInSlice([]int{
								int(RoleTypeGuest), int(RoleTypeHost),
							}),
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(AccessTypeNormal), string(AccessTypeTelepresence), string(AccessTypeTerminal),
								string(AccessTypeOutside), string(AccessTypeMobile), string(AccessTypeTelephone),
								string(AccessTypeIdeahub),
							}, false),
						},
						"is_mute": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.IntInSlice([]int{
								int(MuteTypeOff), int(MuteTypeOn),
							}),
						},
						"is_auto_invite": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.IntInSlice([]int{
								int(InviteTypeAuto), int(InviteTypeManual),
							}),
						},
						"phone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sms": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"cycle_params": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cycle": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(CycleTypeDay), string(CycleTypeWeek), string(CycleTypeMonth),
							}, false),
						},
						"pre_remind": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"start_date": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`),
								"The date format should be 'YYYY-MM-DD'."),
						},
						"end_date": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`),
								"The date format should be 'YYYY-MM-DD'."),
						},
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"points": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_send_notify": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_send_sms": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_send_calendar": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_auto_mute": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_hard_terminal_auto_mute": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_guest_free_password": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"callin_restriction": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.IntInSlice([]int{
								int(CallInTypeAll), int(CallInTypeEnterprise), int(CallInTypeInvitee),
							}),
						},
						"allow_guest_start": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"guest_password": {
							Type:      schema.TypeString,
							Optional:  true,
							Computed:  true,
							Sensitive: true,
						},
						"prolong_time": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"waiting_room_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"conference_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"conference_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"chair_join_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guest_join_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"audience_join_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subconferences": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     subconferenceSchemaResource(),
			},
			"join_password": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"guest": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildConferenceParticipants(participants []interface{}, appId string) []conferences.Participant {
	if len(participants) < 1 {
		return nil
	}

	result := make([]conferences.Participant, len(participants))
	for i, val := range participants {
		participant := val.(map[string]interface{})
		result[i] = conferences.Participant{
			UserUUID:     participant["user_id"].(string),
			AccountId:    participant["account_id"].(string),
			AppId:        appId,
			Name:         participant["name"].(string),
			Role:         utils.Int(participant["role"].(int)),
			Type:         participant["type"].(string),
			IsMute:       utils.Int(participant["is_mute"].(int)),
			IsAutoInvite: utils.Int(participant["is_auto_invite"].(int)),
			Phone:        participant["phone"].(string),
			Email:        participant["email"].(string),
			SMS:          participant["sms"].(string),
		}
	}
	return result
}

func buildConferenceCycleParams(cycleParams []interface{}) *conferences.CycleParams {
	if len(cycleParams) < 1 {
		return nil
	}

	cycleParam := cycleParams[0].(map[string]interface{})
	return &conferences.CycleParams{
		Cycle:         cycleParam["cycle"].(string),
		PreRemindDays: utils.Int(cycleParam["pre_remind"].(int)),
		StartDate:     cycleParam["start_date"].(string),
		EndDate:       cycleParam["end_date"].(string),
		Interval:      cycleParam["interval"].(int),
		Points:        utils.ExpandToIntList(cycleParam["points"].([]interface{})),
	}
}

func buildConferenceConfiguration(configurations []interface{}) *conferences.Configuration {
	if len(configurations) < 1 {
		return nil
	}

	configuration := configurations[0].(map[string]interface{})
	return &conferences.Configuration{
		IsSendNotify:           utils.Bool(configuration["is_send_notify"].(bool)),
		IsSendSms:              utils.Bool(configuration["is_send_sms"].(bool)),
		IsSendCalendar:         utils.Bool(configuration["is_send_calendar"].(bool)),
		IsAutoMute:             utils.Bool(configuration["is_auto_mute"].(bool)),
		IsHardTerminalAutoMute: utils.Bool(configuration["is_hard_terminal_auto_mute"].(bool)),
		IsGuestFreePwd:         utils.Bool(configuration["is_guest_free_password"].(bool)),
		CallInRestriction:      utils.Int(configuration["callin_restriction"].(int)),
		AllowGuestStartConf:    utils.Bool(configuration["allow_guest_start"].(bool)),
		GuestPwd:               configuration["guest_password"].(string),
		ProlongLength:          utils.Int(configuration["prolong_time"].(int)),
		EnableWaitingRoom:      utils.Bool(configuration["waiting_room_enabled"].(bool)),
	}
}

func buildConferenceCreateOpts(d *schema.ResourceData, token string) conferences.CreateOpts {
	opt := conferences.CreateOpts{
		Subject:                d.Get("topic").(string),
		StartTime:              d.Get("start_time").(string),
		Length:                 d.Get("duration").(int),
		MediaTypes:             strings.Join(utils.ExpandToStringListBySet(d.Get("media_types").(*schema.Set)), ","),
		IsAutoRecord:           utils.Int(d.Get("is_auto_record").(int)),
		EncryptMode:            utils.Int(d.Get("encrypt_mode").(int)),
		Language:               d.Get("language").(string),
		TimeZoneID:             strconv.Itoa(d.Get("timezone_id").(int)),
		RecordType:             utils.Int(d.Get("record_type").(int)),
		LiveAddress:            d.Get("live_address").(string),
		AuxAddress:             d.Get("aux_address").(string),
		RecordAuxStream:        utils.Int(d.Get("is_record_aux_stream").(int)),
		RecordAuthType:         utils.Int(d.Get("record_auth_type").(int)),
		ConcurrentParticipants: utils.Int(d.Get("participant_number").(int)),
		Participants:           buildConferenceParticipants(d.Get("participant").([]interface{}), d.Get("app_id").(string)),
		CycleParams:            buildConferenceCycleParams(d.Get("cycle_params").([]interface{})),
		Configuration:          buildConferenceConfiguration(d.Get("configuration").([]interface{})),
		// The authorization token
		Token: token,
	}
	if warRoomId, ok := d.GetOk("meeting_room_id"); ok {
		opt.VmrFlag = utils.Int(1)
		opt.VmrId = warRoomId.(string)
	}

	return opt
}

func resourceConferenceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	token, err := NewMeetingToken(conf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	opt := buildConferenceCreateOpts(d, token)
	log.Printf("[DEBUG] The createOpts of the cloud conference is: %v", opt)
	resp, err := conferences.Create(NewMeetingV1Client(conf), opt)
	if err != nil {
		return diag.Errorf("error creating cloud conference: %s", err)
	}
	if len(resp) < 1 {
		return diag.Errorf("the corresponding cloud conference is not found after the creation is complete, please " +
			"contact customer service to solve this problem")
	}
	d.SetId(resp[0].ID)

	return resourceConferenceRead(ctx, d, meta)
}

func buildConferenceGetOpts(id, userId, token string) conferences.GetOpts {
	return conferences.GetOpts{
		ConferenceId: id,
		UserId:       userId,
		Limit:        500,
		Token:        token,
	}
}

func buildHistoryConferenceGetOpts(uuid, userId, token string) conferences.GetHistoryOpts {
	return conferences.GetHistoryOpts{
		ConferenceUuid: uuid,
		UserId:         userId,
		Limit:          500,
		Token:          token,
	}
}

func buildHistoryConferenceListOpts(userId, startTime, token string, duration int) (conferences.ListHistoryOpts,
	error) {
	tm, err := time.Parse(standardTimeFormat, startTime)
	if err != nil {
		return conferences.ListHistoryOpts{}, nil
	}
	return conferences.ListHistoryOpts{
		UserUUID: userId,
		// The time range here refers to the time range when the meeting starts, the actual start time depends on when
		// the meeting is used.
		StartDate: int(tm.Unix()) * 1000, //ms
		EndDate:   (int(tm.Unix()) + duration*60) * 1000,
		Limit:     500,
		Token:     token,
	}, nil
}

func parseConferenceNotFoundError(respErr error) error {
	var apiErr conferences.ErrResponse
	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok && errCode.Body != nil {
		pErr := json.Unmarshal(errCode.Body, &apiErr)
		if pErr == nil && apiErr.Code == "MMC.111070005" {
			return golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Body: []byte("the conference is outdate"),
				},
			}
		}
	}
	return respErr
}

func GetConferenceDetail(client *golangsdk.ServiceClient, d *schema.ResourceData, token string) (*conferences.GetResp, error) {
	userId := d.Get("user_id").(string)
	opt := buildConferenceGetOpts(d.Id(), userId, token)
	log.Printf("[DEBUG] The GetOpts of the cloud conference is: %v", opt)
	// Get scheduled or online conference.
	resp, err := conferences.Get(client, opt)
	if err == nil {
		return resp, nil
	}
	if _, ok := parseConferenceNotFoundError(err).(golangsdk.ErrDefault404); ok {
		var uuid interface{}
		if uuid, ok = d.GetOk("conference_uuid"); !ok {
			startTime := d.Get("start_time").(string)
			duration := d.Get("duration").(int)
			hisListOpt, err := buildHistoryConferenceListOpts(userId, startTime, token, duration)
			if err != nil {
				return nil, err
			}
			resp, err := conferences.ListHistory(client, hisListOpt)
			if err != nil {
				return nil, err
			}
			if len(resp) > 0 {
				for _, val := range resp {
					if val.ID == d.Id() {
						uuid = val.ConfUUID
					}
				}
			}
			if uuid == "" {
				return nil, fmt.Errorf("no history conference found")
			}
		}
		// Get history conference.
		hisOpt := buildHistoryConferenceGetOpts(uuid.(string), userId, token)
		resp, err = conferences.GetHistory(client, hisOpt)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, err
}

func flattenConferenceParticipants(participants []conferences.ParticipantDetail) []map[string]interface{} {
	if len(participants) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(participants))
	for i, participant := range participants {
		result[i] = map[string]interface{}{
			"user_id":        participant.UserUUID,
			"account_id":     participant.AccountId,
			"name":           participant.Name,
			"role":           participant.Role,
			"type":           participant.Type,
			"is_auto_invite": participant.IsAutoInvite,
			"phone":          participant.ID, // The ID is the phone number of the attendee.
			"email":          participant.Email,
			"sms":            participant.SMS,
		}
	}

	log.Printf("[DEBUG] The participant list result of the conference is %#v", result)
	return result
}

func flattenShowAudiencePolicy(subConference conferences.ShowAudiencePolicy) []map[string]interface{} {
	if reflect.DeepEqual(subConference, conferences.ShowAudiencePolicy{}) {
		return nil
	}

	result := []map[string]interface{}{
		{
			"mode":                subConference.ShowAudienceMode,
			"base_audience_count": subConference.BaseAudienceCount,
			"multiple":            subConference.Multiple,
		},
	}

	log.Printf("[DEBUG] The audience display policy result of the subconference is %#v", result)
	return result
}

func flattenSubConferenceConfiguration(subConference conferences.SubConfiguration) []map[string]interface{} {
	if reflect.DeepEqual(subConference, conferences.SubConfiguration{}) {
		return nil
	}

	result := []map[string]interface{}{
		{
			"callin_restriction":          subConference.CallInRestriction,
			"audience_callin_restriction": subConference.AudienceCallInRestriction,
			"allow_guest_start":           subConference.AllowGuestStartConf,
			"waiting_room_enabled":        subConference.EnableWaitingRoom,
			"show_audience_policies":      flattenShowAudiencePolicy(subConference.ShowAudiencePolicy),
		},
	}

	log.Printf("[DEBUG] The configuration result of the subconference is %#v", result)
	return result
}

func flattenSubConferences(subConferences []conferences.SubConference) []map[string]interface{} {
	if len(subConferences) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(subConferences))
	for i, subConference := range subConferences {
		mediaType := strings.Split(subConference.MediaType, ",")
		sort.Strings(mediaType)
		result[i] = map[string]interface{}{
			"id":               subConference.ID,
			"media_types":      mediaType,
			"start_time":       subConference.StartTime,
			"end_time":         subConference.EndTime,
			"is_auto_record":   subConference.IsAutoRecord,
			"record_auth_type": subConference.RecordAuthType,
			"subconfiguration": flattenSubConferenceConfiguration(subConference.SubConfiguration),
		}
	}

	log.Printf("[DEBUG] The result of the subconference is %#v", result)
	return result
}

func flattenConferenceCycleParams(cycleParams conferences.CycleParams) []map[string]interface{} {
	if reflect.DeepEqual(cycleParams, conferences.CycleParams{}) {
		return nil
	}

	result := []map[string]interface{}{
		{
			"cycle":      cycleParams.Cycle,
			"pre_remind": cycleParams.PreRemindDays,
			"start_date": cycleParams.StartDate,
			"end_date":   cycleParams.EndDate,
			"interval":   cycleParams.Interval,
			"points":     cycleParams.Points,
		},
	}

	log.Printf("[DEBUG] The cycle parameters result of the conference is %#v", result)
	return result
}

func flattenConferenceConfiguration(configuration conferences.Configuration) []map[string]interface{} {
	if reflect.DeepEqual(configuration, conferences.Configuration{}) {
		return nil
	}

	// No password return, it is a sensitive parameter.
	result := []map[string]interface{}{
		{
			"is_send_notify":             configuration.IsSendNotify,
			"is_send_sms":                configuration.IsSendSms,
			"is_send_calendar":           configuration.IsSendCalendar,
			"is_auto_mute":               configuration.IsAutoMute,
			"is_hard_terminal_auto_mute": configuration.IsHardTerminalAutoMute,
			"is_guest_free_password":     configuration.IsGuestFreePwd,
			"callin_restriction":         configuration.CallInRestriction,
			"allow_guest_start":          configuration.AllowGuestStartConf,
			"prolong_time":               configuration.ProlongLength,
			"waiting_room_enabled":       configuration.EnableWaitingRoom,
		},
	}

	log.Printf("[DEBUG] The configuration result of the conference is %#v", result)
	return result
}

func flattenJoinPasswords(passwords []conferences.PasswordEntry) []map[string]interface{} {
	if len(passwords) < 1 {
		return nil
	}

	result := make(map[string]interface{})
	for _, pwd := range passwords {
		// Only set passwords of chair role and guest role.
		if pwd.ConferenceRole == string(ConferenceRoleChair) {
			result["host"] = pwd.Password
		}
		if pwd.ConferenceRole == string(ConferenceRoleGuest) {
			result["guest"] = pwd.Password
		}
	}
	return []map[string]interface{}{result}
}

func resourceConferenceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	token, err := NewMeetingToken(conf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := GetConferenceDetail(NewMeetingV1Client(conf), d, token)
	if err != nil {
		return diag.Errorf("error getting conference (%s) detail: %v", d.Id(), err)
	}
	mErr := multierror.Append(nil,
		d.Set("conference_type", resp.Conference.ConfType),
		d.Set("topic", resp.Conference.Subject),
		d.Set("meeting_room_id", resp.Conference.VmrId),
		d.Set("media_types", strings.Split(resp.Conference.MediaTypes, ",")),
		d.Set("is_auto_record", resp.Conference.IsAutoRecord),
		d.Set("language", resp.Conference.Language),
		d.Set("record_type", resp.Conference.RecordType),
		d.Set("live_address", resp.Conference.LiveAddress),
		d.Set("aux_address", resp.Conference.AuxAddress),
		d.Set("is_record_aux_stream", resp.Conference.RecordAuxStream),
		d.Set("record_auth_type", resp.Conference.RecordAuthType),
		d.Set("participant_number", resp.Conference.ConcurrentParticipants),
		// Parameters of the structrue list
		d.Set("participant", flattenConferenceParticipants(resp.Data.Participants)),
		d.Set("cycle_params", flattenConferenceCycleParams(resp.Conference.CycleParams)),
		d.Set("configuration", flattenConferenceConfiguration(resp.Conference.Configuration)),
		// Computed parameters
		d.Set("conference_uuid", resp.Conference.ConfUUID),
		d.Set("access_number", resp.Conference.AccessNumber),
		d.Set("status", resp.Conference.ConferenceState),
		d.Set("chair_join_uri", resp.Conference.ChairJoinUri),
		d.Set("guest_join_uri", resp.Conference.GuestJoinUri),
		d.Set("audience_join_uri", resp.Conference.AudienceJoinUri),
		d.Set("subconferences", flattenSubConferences(resp.Conference.Subconferences)),
		d.Set("join_password", flattenJoinPasswords(resp.Conference.PasswordEntry)),
	)
	if timezoneId, err := strconv.Atoi(resp.Conference.TimeZoneId); err != nil {
		log.Printf("[ERROR] The format of timezone ID (%#v) is wrong: %v", resp.Conference.TimeZoneId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("timezone_id", timezoneId))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func getConferenceControlToken(client *golangsdk.ServiceClient, conferenceId string, chairPwd string) (string, error) {
	opt := conferences.AuthOpts{
		ConferenceId: conferenceId,
		Password:     chairPwd,
	}
	resp, err := conferences.GetControlToken(client, opt)
	if err != nil {
		return "", err
	}
	return resp.TokenDetail.Token, nil
}

func ExtendConference(client *golangsdk.ServiceClient, conferenceId string, chairPwd string, duration int) error {
	token, err := getConferenceControlToken(client, conferenceId, chairPwd)
	if err != nil {
		return err
	}

	opt := conferences.ExtendOpts{
		IsAuto:   0,
		Token:    token,
		Duration: duration,
	}
	return conferences.ExtendTime(client, opt)
}

func getChairPassword(d *schema.ResourceData) string {
	pwdSet := d.Get("join_password").([]interface{})
	if len(pwdSet) < 1 {
		return ""
	}
	pwdInfo := pwdSet[0].(map[string]interface{})
	return pwdInfo["chair"].(string)
}

func buildConferenceUpdateOpts(d *schema.ResourceData, conferenceId, token string) conferences.UpdateOpts {
	opt := conferences.UpdateOpts{
		// Query parameter
		ConferenceID: conferenceId,
		// Arguments
		Length:                 d.Get("duration").(int),
		Subject:                d.Get("topic").(string),
		StartTime:              d.Get("start_time").(string),
		MediaTypes:             strings.Join(utils.ExpandToStringListBySet(d.Get("media_types").(*schema.Set)), ","),
		IsAutoRecord:           utils.Int(d.Get("is_auto_record").(int)),
		EncryptMode:            utils.Int(d.Get("encrypt_mode").(int)),
		Language:               d.Get("language").(string),
		TimeZoneID:             strconv.Itoa(d.Get("timezone_id").(int)),
		RecordType:             utils.Int(d.Get("record_type").(int)),
		LiveAddress:            d.Get("live_address").(string),
		AuxAddress:             d.Get("aux_address").(string),
		RecordAuxStream:        utils.Int(d.Get("is_record_aux_stream").(int)),
		RecordAuthType:         utils.Int(d.Get("record_auth_type").(int)),
		ConcurrentParticipants: utils.Int(d.Get("participant_number").(int)),
		Participants:           buildConferenceParticipants(d.Get("participant").([]interface{}), d.Get("app_id").(string)),
		CycleParams:            buildConferenceCycleParams(d.Get("cycle_params").([]interface{})),
		Configuration:          buildConferenceConfiguration(d.Get("configuration").([]interface{})),
		// The authorization token
		Token: token,
	}
	if warRoomId, ok := d.GetOk("meeting_room_id"); ok {
		opt.VmrFlag = utils.Int(1)
		opt.VmrId = warRoomId.(string)
	}

	return opt
}

func resourceConferenceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	token, err := NewMeetingToken(conf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	opt := buildConferenceGetOpts(d.Id(), d.Get("user_id").(string), token)
	log.Printf("[DEBUG][Before Update] The GetOpts of the cloud conference is: %v", opt)
	client := NewMeetingV1Client(conf)

	switch d.Get("status").(string) {
	case string(ConferenceStatusUsing):
		if d.HasChangeExcept("length") {
			return diag.Errorf("when the conference is in progress, the parameters except the conference length do " +
				"not support changing.")
		}
		oldLen, newLen := d.GetChange("length")
		if newLen.(int) < oldLen.(int) {
			return diag.Errorf("invalid duration, only extension is supported when the conference is in progress.")
		}
		newDuration := newLen.(int) - oldLen.(int)
		err = ExtendConference(client, d.Id(), getChairPassword(d), newDuration)
		if err != nil {
			return diag.FromErr(err)
		}
	case string(ConferenceStatusSchedule):
		opt := buildConferenceUpdateOpts(d, d.Id(), token)
		_, err = conferences.Update(client, opt)
		if err != nil {
			return diag.Errorf("error creating cloud conference: %s", err)
		}
	default:
		return diag.Errorf("update is not supported because of the conference is out of date.")
	}

	return resourceConferenceRead(ctx, d, meta)
}

func buildConferenceDeleteOpts(d *schema.ResourceData, token string) conferences.DeleteOpts {
	result := conferences.DeleteOpts{
		ConferenceId: d.Id(),
		UserId:       d.Get("user_id").(string),
		Type:         1,
		Token:        token,
	}
	if _, ok := d.GetOk("cycle_params"); ok {
		result.IsCycle = true
	}
	return result
}

func resourceConferenceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	token, err := NewMeetingToken(conf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	opt := buildConferenceDeleteOpts(d, token)
	log.Printf("[DEBUG] The DeleteOpts of the cloud conference is: %v", opt)
	if conferences.Delete(NewMeetingV1Client(conf), opt) != nil {
		if _, ok := parseConferenceNotFoundError(err).(golangsdk.ErrDefault404); !ok {
			return diag.Errorf("error deleting conference (%s): %s", d.Id(), err)
		}
	}

	return nil
}

func resourceConferenceImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	var mErr *multierror.Error
	parts := strings.Split(d.Id(), "/")
	switch len(parts) {
	case 3:
		d.SetId(parts[0])
		mErr = multierror.Append(mErr,
			d.Set("account_name", parts[1]),
			d.Set("account_password", parts[2]),
		)
	case 5:
		d.SetId(parts[0])
		mErr = multierror.Append(mErr,
			d.Set("app_id", parts[1]),
			d.Set("app_key", parts[2]),
			d.Set("corp_id", parts[3]),
			d.Set("user_id", parts[4]),
		)
	default:
		return nil, fmt.Errorf("the imported ID specifies an invalid format, must be " +
			"<id>/<account_name>/<account_password> or <id>/<app_id>/<app_key>/<corp_id>/<user_id>.")
	}
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
