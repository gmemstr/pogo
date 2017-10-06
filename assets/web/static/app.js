Vue.component('episode-publish-form', {
	template: '<form enctype="multipart/form-data" action="/admin/publish" method="post"><label for="title">Episode Title</label><input type="text" id="title" name="title"><label for="description">Episode Description</label><textarea name="description" id="description" cols="100" rows="20" style="resize: none;"></textarea><label for="file">Media File</label><input type="file" id="file" name="file"><label for="date">Publish Date</label><input type="date" id="date" name="date"><input type="submit" value="Publish"></form>'
})

Vue.component('custom-css', {
	template: '<form action="/admin/css" method="post" enctype="multipart/form-data"><label for="css">Custom CSS</label><textarea name="css" id="css" cols="120" rows="20"></textarea><br /><input type="submit" value="Submit"></form>'
})

var app = new Vue({
  el: '#app',
  data: {
    header: 'Pogo Admin',
    current_page: "Page Title",
    page: false,
  }
})

window.onhashchange = setpagecontents
window.onload = setpagecontents
// I know I'm probably not using 
// vue.js properly here but it's the
// best I can do right now
function setpagecontents(){
	page = window.location.href.split('#')[1]
	app.page = page
	
	if (page == "publish") {
		app.current_page = "Publish Episode"
	} 
	else if (page == "customcss") {
		app.current_page = "Edit Theme"
		getcss()
	} 
	else {
		app.current_page = "404 Not found"
	}
}

function getcss(){
	var xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function() { 
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
            document.getElementById("css").innerHTML=xmlHttp.responseText;
    }
    xmlHttp.open("GET", "/admin/css", true);
    xmlHttp.send(null);
}