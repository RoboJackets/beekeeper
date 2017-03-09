/**
 * @author willstuckey
 * @date 3/4/17.
 *
 * This file contains utility functions to manage cookies.
 *
 * Any functions regarding specific cookies should go
 * elsewhere.
 */

/**
 * You should probably use a function with a more coarse expiration.
 *
 * creates/sets cookie
 * @param n name
 * @param v value
 * @param t_ms expiration in milliseconds
 * @private
 */
function _set_cookie(n, v, t_ms) {
	var dt = new Date();
	dt.setTime(dt.getTime() + t_ms);
	var exp = ("expires=" + dt.toUTCString());
	document.cookie = (n + "=" + v + ";" + exp + ";path=/");
}

/**
 * creates a cookie
 * @param cname cookie name
 * @param cvalue cookie value
 * @param min_til_exp number of minutes till expiry
 */
function set_cookie_m(cname, cvalue, min_til_exp) {
	_set_cookie(cname, cvalue, min_til_exp * 60 * 1000);
}

/**
 * creates a cookie
 * @param cname cookie name
 * @param cvalue cookie value
 * @param hr_til_exp number of hours till expiry
 */
function set_cookie_h(cname, cvalue, hr_til_exp) {
	set_cookie_m(cname, cvalue, hr_til_exp * 60);
}

/**
 * creates a cookie
 * @param cname cookie name
 * @param cvalue cookie value
 * @param days_til_exp number of hours till expiry
 */
function set_cookie_d(cname, cvalue, days_til_exp) {
	set_cookie_h(cname, cvalue, days_til_exp * 24);
}

/**
 * deletes a cookie
 * @param cname cookie name
 */
function delete_cookie(cname) {
	document.cookie = (cname + '=; Path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;');
}

/**
 * gets a cookie
 * @param name cookie name
 * @returns {*} cookie, null if not found
 */
function get_cookie(name) {
	var nameEQ = name + "=";
	var ca = document.cookie.split(';');
	for (var i = 0; i < ca.length; i++) {
		var c = ca[i];
		while (c.charAt(0) == ' ') {
			c = c.substring(1, c.length);
		}

		if (c.indexOf(nameEQ) == 0) {
			return c.substring(nameEQ.length, c.length);
		}
	}

	return null;
}
