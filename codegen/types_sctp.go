// +build ignore

package sctp_go

//#include <stdint.h>
//#include <linux/types.h>
//#include <sys/socket.h>
//#include <linux/sctp.h>
//#include <netinet/in.h>
//typedef struct iovec                    IoVector;
//typedef struct msghdr                   MsgHeader;
//typedef struct cmsghdr                  CMsgHeader;
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

type IoVector               C.IoVector
type MsgHeader              C.MsgHeader
type CMsgHeader             C.CMsgHeader
type InAddr                 C.InAddr
type In6Addr                C.In6Addr
type SockAddrIn6            C.SockAddrIn6
type SockAddrIn             C.SockAddrIn
type SockAddr               C.SockAddr
type SockAddrStorage        C.SockAddrStorage
type SCTPAssocId            C.sctp_assoc_t
type SCTPInitMsg            C.SCTPInitMsg
type SCTPSndRcvInfo         C.SCTPSndRcvInfo
type SCTPSndInfo            C.SCTPSndInfo
type SCTPRcvInfo            C.SCTPRcvInfo
type SCTPNxtInfo            C.SCTPNxtInfo
type SCTPPrInfo             C.SCTPPrInfo
type SCTPAuthInfo           C.SCTPAuthInfo
type SCTPCmsgData           C.SCTPCmsgData
type SCTPAssocChange        C.SCTPAssocChange
type SCTPPAddrChange        C.SCTPPAddrChange
type SCTPRemoteError        C.SCTPRemoteError
type SCTPSendFailed         C.SCTPSendFailed
type SCTPShutdownEvent      C.SCTPShutdownEvent
type SCTPAdaptationEvent    C.SCTPAdaptationEvent
type SCTPPDApiEvent         C.SCTPPDApiEvent
type SCTPAuthKeyEvent       C.SCTPAuthKeyEvent
type SCTPSenderDryEvent     C.SCTPSenderDryEvent
type SCTPStreamResetEvent   C.SCTPStreamResetEvent
type SCTPAssocResetEvent    C.SCTPAssocResetEvent
type SCTPStreamChangeEvent  C.SCTPStreamChangeEvent
type SCTPEventSubscribe     C.SCTPEventSubscribe
type SCTPNotification       C.SCTPNotification
type SCTPNotificationHeader C.SCTPNotificationHeader
type SCTPRTOInfo            C.SCTPRTOInfo
type SCTPAssocParams        C.SCTPAssocParams
type SCTPSetPeerPrimary     C.SCTPSetPeerPrimary
type SCTPPrimaryAddr        C.SCTPPrimaryAddr
type SCTPSetAdaptation      C.SCTPSetAdaptation
type SCTPPeerAddrParams     C.SCTPPeerAddrParams
type SCTPAuthChunk          C.SCTPAuthChunk
type SCTPHmacAlgo           C.SCTPHmacAlgo
type SCTPAuthKey            C.SCTPAuthKey
type SCTPAuthKeyId          C.SCTPAuthKeyId
type SCTPSackInfo           C.SCTPSackInfo
type SCTPAssocValue         C.SCTPAssocValue
type SCTPStreamValue        C.SCTPStreamValue
type SCTPPeerAddrInfo       C.SCTPPeerAddrInfo
type SCTPStatus             C.SCTPStatus
type SCTPAuthChunks         C.SCTPAuthChunks
type SCTPAssocIds           C.SCTPAssocIds
type SCTPGetAddrsOld        C.SCTPGetAddrsOld
type SCTPGetAddrs           C.SCTPGetAddrs
type SCTPAssocStats         C.SCTPAssocStats
type SCTPPeelOffArg         C.SCTPPeelOffArg
type SCTPPeelOffFlagsArg    C.SCTPPeelOffFlagsArg
type SCTPPeerAddrThresholds C.SCTPPeerAddrThresholds
type SCTPPRStatus           C.SCTPPRStatus
type SCTPDefaultPRInfo      C.SCTPDefaultPRInfo
type SCTPInfo               C.SCTPInfo
type SCTPResetStreams       C.SCTPResetStreams
type SCTPAddStreams         C.SCTPAddStreams
type SCTPEvent              C.SCTPEvent