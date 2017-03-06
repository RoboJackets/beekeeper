/**
 * @author willstuckey
 * @date 3/6/17
 */

function add_info(content) {
	var html_content = "";
	html_content += "<div class='alert alert-info alert-dismissable' role='alert'>";
	html_content += "<a href='#' class='close' data-dismiss='alert' aria-label='close'>×</a>";
	html_content += content;
	html_content += "</div>";
	$("#banner_inject").append(html_content);
}

/**
 * banner init
 */
function banner_init() {
	$("#banner_inject").empty();

	//TODO will fetch global notifications server side
	add_info("<b>INFO: </b> This site is under development. Links, actions, and data are all volatile.");

	/*
	 * TODO
	 * read notification cookies
	 * display notifications
	 * clear notification cookies
	 */
}