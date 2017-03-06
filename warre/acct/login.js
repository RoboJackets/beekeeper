/**
 * @author willstuckey
 * @date 3/6/17
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

/**
 * initialize login page
 */
function page_init() {
	//redirect to index if already authenticated
	//TODO there should probably be a redirect source cookie
	if (__is_authenticated()) {
		window.location.replace("/users/preferences.html");
		return;
	}

	add_info("<b>INFO: </b> You may login with the dummy info username: user, password: password");

	$("#login_button").click(function () {
		if (__authenticate($("#email_input").val(), $("#passwd_input").val())) {
			window.location.replace("/user/preferences.html");
		} else {
			alert("invalid login creds");
		}
	})
}