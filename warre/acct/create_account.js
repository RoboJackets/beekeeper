/**
 * @author willstuckey
 * @date 3/6/17.
 */

/**
 * this page does not require authentication to view
 * DATA INTEGRITY SHOULD BE ENFORCED SERVER SIDE
 * @type {boolean}
 */
var AUTH_REQUIRED = false;

/**
 * this page has no navbar
 * @type {boolean}
 */
var NAV_PRESENT = false;

/**
 * this page requires banner injection
 * @type {boolean}
 */
var BANNER_PRESENT = true;

function page_init() {
	if (__is_authenticated()) {
		window.location.replace("/users/preferences.html");
		return;
	}

	add_info("<b>INFO: </b> This page has no functionality. Please use dummy info username: user, password: password");

	$("#create_button").click(function () {
		alert("Please use dummy account info.");
		window.location.replace("/acct/login.html");
	})
}