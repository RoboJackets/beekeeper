/**
 * @author willstuckey
 * @date 3/6/17
 *
 * TODO this feature will require a more spohisticated set
 * of cookie management functions in cookie.js
 */

/**
 * cookie string prefix for frontend banner notifications
 * @type {string}
 * @private
 */
var _BANNER_LOCAL = "c_banner_local";

/**
 * cookie string for local info banners
 * @type {string}
 */
var C_BANNER_LOCAL_INFO = _BANNER_LOCAL + "_info";

/**
 * cookie string for local success banners
 * @type {string}
 */
var C_BANNER_LOCAL_SUCCESS = _BANNER_LOCAL + "_suc";

/**
 * cookie string for local warning banners
 * @type {string}
 */
var C_BANNER_LOCAL_WARNING = _BANNER_LOCAL + "_warn";

/**
 * cookie string for local error banners
 * @type {string}
 */
var C_BANNER_LOCAL_ERROR = _BANNER_LOCAL + "_err";

/**
 * cookie string for cookie the indicates the frontend page
 * that issued a redirect to another frontend page
 * @type {string}
 */
var C_SOURCE_OF_REDIRECT = "c_source_of_redirect";

/**
 * cookie string for cookie that indicates why a frontend page
 * issued a redirect to another frontend page
 * @type {string}
 */
var C_REASON_FOR_REDIRECT = "c_reason_for_redirect";