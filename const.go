package zrule

// Redis Keys
const CACHE_ZRULE_AUTH_ATTEMPT = "zrule::auth::attempt::%s"
const CACHE_ZRULE_AUTH_TOKEN = "zrule::auth::attempt::%s::TOKEN"
const CACHE_CCP_JWKS = "ccp::jwks"

// ESI Timestamp Format
const ESI_EXPIRES_HEADER_FORMAT = "Mon, 02 Jan 2006 15:04:05 MST"

// REDIS KEY
const CACHE_ESI_ERROR_COUNT = "esi::error::count"
const CACHE_ESI_ERROR_RESET = "esi::error::reset"
const CACHE_ESI_TRACKING_STATUS = "esi::tracking::status"

const CACHE_ESI_TRACKING_OK = "zrule::esi::tracking::ok"                     // 200
const CACHE_ESI_TRACKING_NOT_MODIFIED = "zrule::esi::tracking::not_modified" // 304
const CACHE_ESI_TRACKING_CALM_DOWN = "zrule::esi::tracking::calm_down"       // 420
const CACHE_ESI_TRACKING_4XX = "zrule::esi::tracking::4xx"                   // Does not include 420s. Those are in the calm down set
const CACHE_ESI_TRACKING_5XX = "zrule::esi::tracking::5xx"

const CACHE_ALLIANCE = "zrule::alliance::%d"
const CACHE_CORPORATION = "zrule::corporation::%d"
const CACHE_CHARACTER = "zrule::character::%d"

const CACHE_REGION = "zrule::region::%d"
const CACHE_CONSTELLATION = "zrule::constellation::%d"
const CACHE_SOLARSYSTEM = "zrule::solarsystem::%d"

const CACHE_ITEM = "zrule::item::%d"
const CACHE_ITEMGROUP = "zrule::itemgroup::%d"

const QUEUES_KILLMAIL_PROCESSING = "zrule::killmail::processing"
const QUEUE_STOP = "zrule::queue::stop"
const QUEUE_RESTART_TRACKER = "zrule::tracker::restart"
const QUEUES_KILLMAIL_MATCHED = "zrule::killmail::matched"
