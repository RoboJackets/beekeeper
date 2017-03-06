/**
 * @author willstuckey
 * @date 3/4/17.
 */

function on_page_ready() {
	// if the user is not authenticated, take them to the login page
	if (AUTH_REQUIRED && !(__is_authenticated())) {
		window.location.replace("/acct/login.html");
		return;
	}

	if (NAV_PRESENT) {
		nav_init();
	}

	if (BANNER_PRESENT) {
		banner_init();
	}

	// calls page init function from corresponding page js file
	page_init();
}