/**
 * @author willstuckey
 * @date 3/6/17
 */

/**
 * this is our auth state until backend is online
 * @type {boolean}
 * @private
 */
var __authenticated = false;

/**
 * check if a user is "authenticated"
 * @returns {boolean}
 * @private
 */
function __is_authenticated() {
	return __authenticated;
}

/**
 * "authenticate" a user
 * @param un
 * @param pw
 * @returns {boolean}
 * @private
 */
function __authenticate(un, pw) {
	if (un == "user" && pw == "password") {
		__authenticated = true;
	}

	return __authenticated;
}

/**
 * "deauthenticate" a user
 * @private
 */
function __deauthenticate() {
	__authenticated = false;
}