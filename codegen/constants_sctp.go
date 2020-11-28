// +build ignore

package sctp_go

//#include <stdint.h>
//#include <linux/types.h>
//#include <sys/socket.h>
//#include <linux/sctp.h>
//#include <netinet/in.h>
//#include <netinet/sctp.h>
//typedef struct in_addr                  InAddr;
//typedef struct in6_addr                 In6Addr;
//typedef struct sockaddr_in6             SockAddrIn6;
//typedef struct sockaddr_in              SockAddrIn;
//typedef struct sockaddr                 SockAddr;
//typedef struct sockaddr_storage         SockAddrStorage;
//typedef struct sctp_sndrcvinfo          SCTPSndRcvInfo;
//typedef struct sctp_initmsg             SCTPInitMsg;
//typedef struct sctp_sndinfo             SCTPSndInfo;
//typedef struct sctp_rcvinfo             SCTPRcvInfo;
//typedef struct sctp_nxtinfo             SCTPNxtInfo;
//typedef struct sctp_prinfo              SCTPPrInfo;
//typedef struct sctp_authinfo            SCTPAuthInfo;
//typedef struct sctp_assoc_change        SCTPAssocChange;
//typedef struct sctp_paddr_change        SCTPPAddrChange;
//typedef struct sctp_remote_error        SCTPRemoteError;
//typedef struct sctp_send_failed         SCTPSendFailed;
//typedef struct sctp_shutdown_event      SCTPShutdownEvent;
//typedef struct sctp_adaptation_event    SCTPAdaptationEvent;
//typedef struct sctp_pdapi_event         SCTPPDApiEvent;
//typedef struct sctp_authkey_event       SCTPAuthKeyEvent;
//typedef struct sctp_sender_dry_event    SCTPSenderDryEvent;
//typedef struct sctp_stream_reset_event  SCTPStreamResetEvent;
//typedef struct sctp_assoc_reset_event   SCTPAssocResetEvent;
//typedef struct sctp_stream_change_event SCTPStreamChangeEvent;
//typedef struct sctp_event_subscribe     SCTPEventSubscribe;
//typedef union sctp_notification         SCTPNotification;
//typedef sctp_cmsg_data_t                SCTPCmsgData;
//typedef struct sn_header {
//	__u16 sn_type;
//	__u16 sn_flags;
//	__u32 sn_length;
//} SCTPNotificationHeader;
//typedef struct sctp_rtoinfo             SCTPRTOInfo;
//typedef struct sctp_assocparams         SCTPAssocParams;
//typedef struct sctp_setpeerprim         SCTPSetPeerPrimary;
//typedef struct sctp_prim                SCTPPrimaryAddr;
//typedef struct sctp_setadaptation       SCTPSetAdaptation;
//typedef struct sctp_paddrparams         SCTPPeerAddrParams;
//typedef struct sctp_authchunk           SCTPAuthChunk;
//typedef struct sctp_hmacalgo            SCTPHmacAlgo;
//typedef struct sctp_authkey             SCTPAuthKey;
//typedef struct sctp_authkeyid           SCTPAuthKeyId;
//typedef struct sctp_sack_info           SCTPSackInfo;
//typedef struct sctp_assoc_value         SCTPAssocValue;
//typedef struct sctp_stream_value        SCTPStreamValue;
//typedef struct sctp_paddrinfo           SCTPPeerAddrInfo;
//typedef struct sctp_status              SCTPStatus;
//typedef struct sctp_authchunks          SCTPAuthChunks;
//typedef struct sctp_assoc_ids           SCTPAssocIds;
//typedef struct sctp_getaddrs_old        SCTPGetAddrsOld;
//typedef struct sctp_getaddrs            SCTPGetAddrs;
//typedef struct sctp_assoc_stats         SCTPAssocStats;
//typedef sctp_peeloff_arg_t              SCTPPeelOffArg;
//typedef sctp_peeloff_flags_arg_t        SCTPPeelOffFlagsArg;
//typedef struct sctp_paddrthlds          SCTPPeerAddrThresholds;
//typedef struct sctp_prstatus            SCTPPRStatus;
//typedef struct sctp_default_prinfo      SCTPDefaultPRInfo;
//typedef struct sctp_info                SCTPInfo;
//typedef struct sctp_reset_streams       SCTPResetStreams;
//typedef struct sctp_add_streams         SCTPAddStreams;
//typedef struct sctp_event               SCTPEvent;
import "C"

const (
	SOL_SCTP                       = C.SOL_SCTP
	IPPROTO_SCTP                   = C.IPPROTO_SCTP
	InAddrSize                     = C.sizeof_InAddr
	In6AddrSize                    = C.sizeof_In6Addr
	SockAddrInSize                 = C.sizeof_SockAddrIn
	SockAddrIn6Size                = C.sizeof_SockAddrIn6
	SockAddrSize                   = C.sizeof_SockAddr
	SockAddrStorageSize            = C.sizeof_SockAddrStorage
	SCTPSndRcvInfoSize             = C.sizeof_SCTPSndRcvInfo
	SCTPInitMsgSize                = C.sizeof_SCTPInitMsg
	SCTPSndInfoSize                = C.sizeof_SCTPSndInfo
	SCTPRcvInfoSize                = C.sizeof_SCTPRcvInfo
	SCTPNxtInfoSize                = C.sizeof_SCTPNxtInfo
	SCTPPrInfoSize                 = C.sizeof_SCTPPrInfo
	SCTPAuthInfoSize               = C.sizeof_SCTPAuthInfo
	SCTPAssocChangeSize            = C.sizeof_SCTPAssocChange
	SCTPPAddrChangeSize            = C.sizeof_SCTPPAddrChange
	SCTPRemoteErrorSize            = C.sizeof_SCTPRemoteError
	SCTPSendFailedSize             = C.sizeof_SCTPSendFailed
	SCTPShutdownEventSize          = C.sizeof_SCTPShutdownEvent
	SCTPAdaptationEventSize        = C.sizeof_SCTPAdaptationEvent
	SCTPPDApiEventSize             = C.sizeof_SCTPPDApiEvent
	SCTPAuthKeyEventSize           = C.sizeof_SCTPAuthKeyEvent
	SCTPSenderDryEventSize         = C.sizeof_SCTPSenderDryEvent
	SCTPStreamResetEventSize       = C.sizeof_SCTPStreamResetEvent
	SCTPAssocResetEventSize        = C.sizeof_SCTPAssocResetEvent
	SCTPStreamChangeEventSize      = C.sizeof_SCTPStreamChangeEvent
	SCTPEventSubscribeSize         = C.sizeof_SCTPEventSubscribe
	SCTPNotificationSize           = C.sizeof_SCTPNotification
	SCTPCmsgDataSize               = C.sizeof_SCTPCmsgData
	SCTPNotificationHeaderSize     = C.sizeof_SCTPNotificationHeader
	SCTPRTOInfoSize                = C.sizeof_SCTPRTOInfo
	SCTPAssocParamsSize            = C.sizeof_SCTPAssocParams
	SCTPSetPeerPrimarySize         = C.sizeof_SCTPSetPeerPrimary
	SCTPPrimaryAddrSize            = C.sizeof_SCTPPrimaryAddr
	SCTPSetAdaptationSize          = C.sizeof_SCTPSetAdaptation
	SCTPPeerAddrParamsSize         = C.sizeof_SCTPPeerAddrParams
	SCTPAuthChunkSize              = C.sizeof_SCTPAuthChunk
	SCTPHmacAlgoSize               = C.sizeof_SCTPHmacAlgo
	SCTPAuthKeySize                = C.sizeof_SCTPAuthKey
	SCTPAuthKeyIdSize              = C.sizeof_SCTPAuthKeyId
	SCTPSackInfoSize               = C.sizeof_SCTPSackInfo
	SCTPAssocValueSize             = C.sizeof_SCTPAssocValue
	SCTPStreamValueSize            = C.sizeof_SCTPStreamValue
	SCTPPeerAddrInfoSize           = C.sizeof_SCTPPeerAddrInfo
	SCTPStatusSize                 = C.sizeof_SCTPStatus
	SCTPAuthChunksSize             = C.sizeof_SCTPAuthChunks
	SCTPAssocIdsSize               = C.sizeof_SCTPAssocIds
	SCTPGetAddrsOldSize            = C.sizeof_SCTPGetAddrsOld
	SCTPGetAddrsSize               = C.sizeof_SCTPGetAddrs
	SCTPAssocStatsSize             = C.sizeof_SCTPAssocStats
	SCTPPeelOffArgSize             = C.sizeof_SCTPPeelOffArg
	SCTPPeelOffFlagsArgSize        = C.sizeof_SCTPPeelOffFlagsArg
	SCTPPeerAddrThresholdsSize     = C.sizeof_SCTPPeerAddrThresholds
	SCTPPRStatusSize               = C.sizeof_SCTPPRStatus
	SCTPDefaultPRInfoSize          = C.sizeof_SCTPDefaultPRInfo
	SCTPInfoSize                   = C.sizeof_SCTPInfo
	SCTPResetStreamsSize           = C.sizeof_SCTPResetStreams
	SCTPAddStreamsSize             = C.sizeof_SCTPAddStreams
	SCTPEventSize                  = C.sizeof_SCTPEvent
	SCTP_FUTURE_ASSOC              = C.SCTP_FUTURE_ASSOC
	SCTP_CURRENT_ASSOC             = C.SCTP_CURRENT_ASSOC
	SCTP_ALL_ASSOC                 = C.SCTP_ALL_ASSOC
	SCTP_RTOINFO                   = C.SCTP_RTOINFO
	SCTP_ASSOCINFO                 = C.SCTP_ASSOCINFO
	SCTP_INITMSG                   = C.SCTP_INITMSG
	SCTP_NODELAY                   = C.SCTP_NODELAY
	SCTP_AUTOCLOSE                 = C.SCTP_AUTOCLOSE
	SCTP_SET_PEER_PRIMARY_ADDR     = C.SCTP_SET_PEER_PRIMARY_ADDR
	SCTP_PRIMARY_ADDR              = C.SCTP_PRIMARY_ADDR
	SCTP_ADAPTATION_LAYER          = C.SCTP_ADAPTATION_LAYER
	SCTP_DISABLE_FRAGMENTS         = C.SCTP_DISABLE_FRAGMENTS
	SCTP_PEER_ADDR_PARAMS          = C.SCTP_PEER_ADDR_PARAMS
	SCTP_DEFAULT_SEND_PARAM        = C.SCTP_DEFAULT_SEND_PARAM
	SCTP_EVENTS                    = C.SCTP_EVENTS
	SCTP_I_WANT_MAPPED_V4_ADDR     = C.SCTP_I_WANT_MAPPED_V4_ADDR
	SCTP_MAXSEG                    = C.SCTP_MAXSEG
	SCTP_STATUS                    = C.SCTP_STATUS
	SCTP_GET_PEER_ADDR_INFO        = C.SCTP_GET_PEER_ADDR_INFO
	SCTP_DELAYED_ACK_TIME          = C.SCTP_DELAYED_ACK_TIME
	SCTP_DELAYED_ACK               = C.SCTP_DELAYED_ACK
	SCTP_DELAYED_SACK              = C.SCTP_DELAYED_SACK
	SCTP_CONTEXT                   = C.SCTP_CONTEXT
	SCTP_FRAGMENT_INTERLEAVE       = C.SCTP_FRAGMENT_INTERLEAVE
	SCTP_PARTIAL_DELIVERY_POINT    = C.SCTP_PARTIAL_DELIVERY_POINT
	SCTP_MAX_BURST                 = C.SCTP_MAX_BURST
	SCTP_AUTH_CHUNK                = C.SCTP_AUTH_CHUNK
	SCTP_HMAC_IDENT                = C.SCTP_HMAC_IDENT
	SCTP_AUTH_KEY                  = C.SCTP_AUTH_KEY
	SCTP_AUTH_ACTIVE_KEY           = C.SCTP_AUTH_ACTIVE_KEY
	SCTP_AUTH_DELETE_KEY           = C.SCTP_AUTH_DELETE_KEY
	SCTP_PEER_AUTH_CHUNKS          = C.SCTP_PEER_AUTH_CHUNKS
	SCTP_LOCAL_AUTH_CHUNKS         = C.SCTP_LOCAL_AUTH_CHUNKS
	SCTP_GET_ASSOC_NUMBER          = C.SCTP_GET_ASSOC_NUMBER
	SCTP_GET_ASSOC_ID_LIST         = C.SCTP_GET_ASSOC_ID_LIST
	SCTP_AUTO_ASCONF               = C.SCTP_AUTO_ASCONF
	SCTP_PEER_ADDR_THLDS           = C.SCTP_PEER_ADDR_THLDS
	SCTP_RECVRCVINFO               = C.SCTP_RECVRCVINFO
	SCTP_RECVNXTINFO               = C.SCTP_RECVNXTINFO
	SCTP_DEFAULT_SNDINFO           = C.SCTP_DEFAULT_SNDINFO
	SCTP_AUTH_DEACTIVATE_KEY       = C.SCTP_AUTH_DEACTIVATE_KEY
	SCTP_REUSE_PORT                = C.SCTP_REUSE_PORT
	SCTP_SOCKOPT_BINDX_ADD         = C.SCTP_SOCKOPT_BINDX_ADD
	SCTP_SOCKOPT_BINDX_REM         = C.SCTP_SOCKOPT_BINDX_REM
	SCTP_SOCKOPT_PEELOFF           = C.SCTP_SOCKOPT_PEELOFF
	SCTP_SOCKOPT_CONNECTX_OLD      = C.SCTP_SOCKOPT_CONNECTX_OLD
	SCTP_GET_PEER_ADDRS            = C.SCTP_GET_PEER_ADDRS
	SCTP_GET_LOCAL_ADDRS           = C.SCTP_GET_LOCAL_ADDRS
	SCTP_SOCKOPT_CONNECTX          = C.SCTP_SOCKOPT_CONNECTX
	SCTP_SOCKOPT_CONNECTX3         = C.SCTP_SOCKOPT_CONNECTX3
	SCTP_GET_ASSOC_STATS           = C.SCTP_GET_ASSOC_STATS
	SCTP_PR_SUPPORTED              = C.SCTP_PR_SUPPORTED
	SCTP_DEFAULT_PRINFO            = C.SCTP_DEFAULT_PRINFO
	SCTP_PR_ASSOC_STATUS           = C.SCTP_PR_ASSOC_STATUS
	SCTP_PR_STREAM_STATUS          = C.SCTP_PR_STREAM_STATUS
	SCTP_RECONFIG_SUPPORTED        = C.SCTP_RECONFIG_SUPPORTED
	SCTP_ENABLE_STREAM_RESET       = C.SCTP_ENABLE_STREAM_RESET
	SCTP_RESET_STREAMS             = C.SCTP_RESET_STREAMS
	SCTP_RESET_ASSOC               = C.SCTP_RESET_ASSOC
	SCTP_ADD_STREAMS               = C.SCTP_ADD_STREAMS
	SCTP_SOCKOPT_PEELOFF_FLAGS     = C.SCTP_SOCKOPT_PEELOFF_FLAGS
	SCTP_STREAM_SCHEDULER          = C.SCTP_STREAM_SCHEDULER
	SCTP_STREAM_SCHEDULER_VALUE    = C.SCTP_STREAM_SCHEDULER_VALUE
	SCTP_INTERLEAVING_SUPPORTED    = C.SCTP_INTERLEAVING_SUPPORTED
	SCTP_SENDMSG_CONNECT           = C.SCTP_SENDMSG_CONNECT
	SCTP_EVENT                     = C.SCTP_EVENT
	SCTP_ASCONF_SUPPORTED          = C.SCTP_ASCONF_SUPPORTED
	SCTP_AUTH_SUPPORTED            = C.SCTP_AUTH_SUPPORTED
	SCTP_ECN_SUPPORTED             = C.SCTP_ECN_SUPPORTED
	SCTP_PR_SCTP_NONE              = C.SCTP_PR_SCTP_NONE
	SCTP_PR_SCTP_TTL               = C.SCTP_PR_SCTP_TTL
	SCTP_PR_SCTP_RTX               = C.SCTP_PR_SCTP_RTX
	SCTP_PR_SCTP_PRIO              = C.SCTP_PR_SCTP_PRIO
	SCTP_PR_SCTP_MAX               = C.SCTP_PR_SCTP_MAX
	SCTP_PR_SCTP_MASK              = C.SCTP_PR_SCTP_MASK
	SCTP_ENABLE_RESET_STREAM_REQ   = C.SCTP_ENABLE_RESET_STREAM_REQ
	SCTP_ENABLE_RESET_ASSOC_REQ    = C.SCTP_ENABLE_RESET_ASSOC_REQ
	SCTP_ENABLE_CHANGE_ASSOC_REQ   = C.SCTP_ENABLE_CHANGE_ASSOC_REQ
	SCTP_ENABLE_STRRESET_MASK      = C.SCTP_ENABLE_STRRESET_MASK
	SCTP_STREAM_RESET_INCOMING     = C.SCTP_STREAM_RESET_INCOMING
	SCTP_STREAM_RESET_OUTGOING     = C.SCTP_STREAM_RESET_OUTGOING
	SCTP_MSG_NOTIFICATION          = C.MSG_NOTIFICATION
	SCTP_UNORDERED                 = C.SCTP_UNORDERED
	SCTP_ADDR_OVER                 = C.SCTP_ADDR_OVER
	SCTP_ABORT                     = C.SCTP_ABORT
	SCTP_SACK_IMMEDIATELY          = C.SCTP_SACK_IMMEDIATELY
	SCTP_SENDALL                   = C.SCTP_SENDALL
	SCTP_PR_SCTP_ALL               = C.SCTP_PR_SCTP_ALL
	SCTP_NOTIFICATION              = C.SCTP_NOTIFICATION
	SCTP_EOF                       = C.SCTP_EOF
	SCTP_INIT                      = C.SCTP_INIT
	SCTP_SNDRCV                    = C.SCTP_SNDRCV
	SCTP_SNDINFO                   = C.SCTP_SNDINFO
	SCTP_RCVINFO                   = C.SCTP_RCVINFO
	SCTP_NXTINFO                   = C.SCTP_NXTINFO
	SCTP_PRINFO                    = C.SCTP_PRINFO
	SCTP_AUTHINFO                  = C.SCTP_AUTHINFO
	SCTP_DSTADDRV4                 = C.SCTP_DSTADDRV4
	SCTP_DSTADDRV6                 = C.SCTP_DSTADDRV6
	SCTP_COMM_UP                   = C.SCTP_COMM_UP
	SCTP_COMM_LOST                 = C.SCTP_COMM_LOST
	SCTP_RESTART                   = C.SCTP_RESTART
	SCTP_SHUTDOWN_COMP             = C.SCTP_SHUTDOWN_COMP
	SCTP_CANT_STR_ASSOC            = C.SCTP_CANT_STR_ASSOC
	SCTP_ADDR_AVAILABLE            = C.SCTP_ADDR_AVAILABLE
	SCTP_ADDR_UNREACHABLE          = C.SCTP_ADDR_UNREACHABLE
	SCTP_ADDR_REMOVED              = C.SCTP_ADDR_REMOVED
	SCTP_ADDR_ADDED                = C.SCTP_ADDR_ADDED
	SCTP_ADDR_MADE_PRIM            = C.SCTP_ADDR_MADE_PRIM
	SCTP_ADDR_CONFIRMED            = C.SCTP_ADDR_CONFIRMED
	SCTP_DATA_UNSENT               = C.SCTP_DATA_UNSENT
	SCTP_DATA_SENT                 = C.SCTP_DATA_SENT
	SCTP_PARTIAL_DELIVERY_ABORTED  = C.SCTP_PARTIAL_DELIVERY_ABORTED
	SCTP_AUTH_NEW_KEY              = C.SCTP_AUTH_NEW_KEY
	SCTP_AUTH_FREE_KEY             = C.SCTP_AUTH_FREE_KEY
	SCTP_AUTH_NO_AUTH              = C.SCTP_AUTH_NO_AUTH
	SCTP_STREAM_RESET_INCOMING_SSN = C.SCTP_STREAM_RESET_INCOMING_SSN
	SCTP_STREAM_RESET_OUTGOING_SSN = C.SCTP_STREAM_RESET_OUTGOING_SSN
	SCTP_STREAM_RESET_DENIED       = C.SCTP_STREAM_RESET_DENIED
	SCTP_STREAM_RESET_FAILED       = C.SCTP_STREAM_RESET_FAILED
	SCTP_ASSOC_RESET_DENIED        = C.SCTP_ASSOC_RESET_DENIED
	SCTP_ASSOC_RESET_FAILED        = C.SCTP_ASSOC_RESET_FAILED
	SCTP_ASSOC_CHANGE_DENIED       = C.SCTP_ASSOC_CHANGE_DENIED
	SCTP_ASSOC_CHANGE_FAILED       = C.SCTP_ASSOC_CHANGE_FAILED
	SCTP_STREAM_CHANGE_DENIED      = C.SCTP_STREAM_CHANGE_DENIED
	SCTP_STREAM_CHANGE_FAILED      = C.SCTP_STREAM_CHANGE_FAILED
	SCTP_SN_TYPE_BASE              = C.SCTP_SN_TYPE_BASE
	SCTP_DATA_IO_EVENT             = C.SCTP_DATA_IO_EVENT
	SCTP_ASSOC_CHANGE              = C.SCTP_ASSOC_CHANGE
	SCTP_PEER_ADDR_CHANGE          = C.SCTP_PEER_ADDR_CHANGE
	SCTP_SEND_FAILED               = C.SCTP_SEND_FAILED
	SCTP_REMOTE_ERROR              = C.SCTP_REMOTE_ERROR
	SCTP_SHUTDOWN_EVENT            = C.SCTP_SHUTDOWN_EVENT
	SCTP_PARTIAL_DELIVERY_EVENT    = C.SCTP_PARTIAL_DELIVERY_EVENT
	SCTP_ADAPTATION_INDICATION     = C.SCTP_ADAPTATION_INDICATION
	SCTP_AUTHENTICATION_EVENT      = C.SCTP_AUTHENTICATION_EVENT
	SCTP_SENDER_DRY_EVENT          = C.SCTP_SENDER_DRY_EVENT
	SCTP_STREAM_RESET_EVENT        = C.SCTP_STREAM_RESET_EVENT
	SCTP_ASSOC_RESET_EVENT         = C.SCTP_ASSOC_RESET_EVENT
	SCTP_STREAM_CHANGE_EVENT       = C.SCTP_STREAM_CHANGE_EVENT
	SCTP_SN_TYPE_MAX               = C.SCTP_SN_TYPE_MAX
	SCTP_FAILED_THRESHOLD          = C.SCTP_FAILED_THRESHOLD
	SCTP_RECEIVED_SACK             = C.SCTP_RECEIVED_SACK
	SCTP_HEARTBEAT_SUCCESS         = C.SCTP_HEARTBEAT_SUCCESS
	SCTP_RESPONSE_TO_USER_REQ      = C.SCTP_RESPONSE_TO_USER_REQ
	SCTP_INTERNAL_ERROR            = C.SCTP_INTERNAL_ERROR
	SCTP_SHUTDOWN_GUARD_EXPIRES    = C.SCTP_SHUTDOWN_GUARD_EXPIRES
	SCTP_PEER_FAULTY               = C.SCTP_PEER_FAULTY
	SPP_HB_ENABLE                  = C.SPP_HB_ENABLE
	SPP_HB_DISABLE                 = C.SPP_HB_DISABLE
	SPP_HB                         = C.SPP_HB
	SPP_HB_DEMAND                  = C.SPP_HB_DEMAND
	SPP_PMTUD_ENABLE               = C.SPP_PMTUD_ENABLE
	SPP_PMTUD_DISABLE              = C.SPP_PMTUD_DISABLE
	SPP_PMTUD                      = C.SPP_PMTUD
	SPP_SACKDELAY_ENABLE           = C.SPP_SACKDELAY_ENABLE
	SPP_SACKDELAY_DISABLE          = C.SPP_SACKDELAY_DISABLE
	SPP_SACKDELAY                  = C.SPP_SACKDELAY
	SPP_HB_TIME_IS_ZERO            = C.SPP_HB_TIME_IS_ZERO
	SPP_IPV6_FLOWLABEL             = C.SPP_IPV6_FLOWLABEL
	SPP_DSCP                       = C.SPP_DSCP
	SCTP_AUTH_HMAC_ID_SHA1         = C.SCTP_AUTH_HMAC_ID_SHA1
	SCTP_AUTH_HMAC_ID_SHA256       = C.SCTP_AUTH_HMAC_ID_SHA256
	SCTP_INACTIVE                  = C.SCTP_INACTIVE
	SCTP_PF                        = C.SCTP_PF
	SCTP_ACTIVE                    = C.SCTP_ACTIVE
	SCTP_UNCONFIRMED               = C.SCTP_UNCONFIRMED
	SCTP_UNKNOWN                   = C.SCTP_UNKNOWN
	SCTP_EMPTY                     = C.SCTP_EMPTY
	SCTP_CLOSED                    = C.SCTP_CLOSED
	SCTP_COOKIE_WAIT               = C.SCTP_COOKIE_WAIT
	SCTP_COOKIE_ECHOED             = C.SCTP_COOKIE_ECHOED
	SCTP_ESTABLISHED               = C.SCTP_ESTABLISHED
	SCTP_SHUTDOWN_PENDING          = C.SCTP_SHUTDOWN_PENDING
	SCTP_SHUTDOWN_SENT             = C.SCTP_SHUTDOWN_SENT
	SCTP_SHUTDOWN_RECEIVED         = C.SCTP_SHUTDOWN_RECEIVED
	SCTP_SHUTDOWN_ACK_SENT         = C.SCTP_SHUTDOWN_ACK_SENT
	SCTP_BINDX_ADD_ADDR            = C.SCTP_BINDX_ADD_ADDR
	SCTP_BINDX_REM_ADDR            = C.SCTP_BINDX_REM_ADDR
	SCTP_SS_FCFS                   = C.SCTP_SS_FCFS
	SCTP_SS_DEFAULT                = C.SCTP_SS_DEFAULT
	SCTP_SS_PRIO                   = C.SCTP_SS_PRIO
	SCTP_SS_RR                     = C.SCTP_SS_RR
	SCTP_SS_MAX                    = C.SCTP_SS_RR
)
