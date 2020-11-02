package zrule

// Cookie Keys
const COOKIE_zrule_AUTH_ATTEMPT = "zrule-auth-attempt"

// Redis Keys
const REDIS_zrule_AUTH_ATTEMPT = "zrule::auth::attempt::%s"
const REDIS_zrule_AUTH_TOKEN = "zrule::auth::attempt::%s::TOKEN"
const REDIS_CCP_JWKS = "ccp::jwks"

// ESI Timestamp Format
const ESI_EXPIRES_HEADER_FORMAT = "Mon, 02 Jan 2006 15:04:05 MST"

// REDIS KEY
const REDIS_ESI_ERROR_COUNT = "esi::error::count"
const REDIS_ESI_ERROR_RESET = "esi::error::reset"
const REDIS_ESI_TRACKING_STATUS = "esi::tracking::status"

const REDIS_ESI_TRACKING_OK = "zrule::esi::tracking::ok"                     // 200
const REDIS_ESI_TRACKING_NOT_MODIFIED = "zrule::esi::tracking::not_modified" // 304
const REDIS_ESI_TRACKING_CALM_DOWN = "zrule::esi::tracking::calm_down"       // 420
const REDIS_ESI_TRACKING_4XX = "zrule::esi::tracking::4xx"                   // Does not include 420s. Those are in the calm down set
const REDIS_ESI_TRACKING_5XX = "zrule::esi::tracking::5xx"

const REDIS_CHARACTER = "zrule::character::%d"

const QUEUES_KILLMAIL_PROCESSING = "zrule::killmail::processing"
const QUEUE_STOP = "zrule::queue::stop"
const QUEUE_RESTART_TRACKER = "zrule::tracker::restart"
const QUEUES_KILLMAIL_MATCHED = "zrule::killmail::matched"
