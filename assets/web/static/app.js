const episodepublishform = {
    template: '<div><h3>Publish Episode</h3><form enctype="multipart/form-data" action="/admin/publish" method="post"><label for="title">Episode Title</label><input type="text" id="title" name="title"><label for="description">Episode Description</label><textarea name="description" id="description" cols="100" rows="20" style="resize: none;"></textarea><label for="file">Media File</label><input type="file" id="file" name="file"><label for="date">Publish Date</label><input type="date" id="date" name="date"><input type="submit" value="Publish"></form></div>',
    current_page: "Publish Episode"
}

const customcss = {
    template: '<div><h3>Edit CSS</h3><form action="/admin/css" method="post" enctype="multipart/form-data"><label for="css">Custom CSS</label><textarea name="css" id="css" cols="120" rows="20">{{ css }}</textarea><br /><input type="submit" value="Submit"></form></div>',
    current_page: "Publish Episode",
    data() {
        return {
            loading: false,
            css: null,
            error: null
        }
    },
    created() {
        // fetch the data when the view is created and the data is
        // already being observed
        this.fetchData()
    },
    watch: {
        // call again the method if the route changes
        '$route': 'fetchData'
    },
    methods: {
        fetchData() {
            this.error = this.css = null
            this.loading = true

            getCss( (err, css) => {
                this.loading = false
                if (err) {
                    this.error = err.toString()
                } else {
                    this.css = css
                }
            })
        }
    }
}

const routes = [
    { path: '/publish', component: episodepublishform },
    { path: '/theme', component: customcss }
]

const router = new VueRouter({
    routes
})

const app = new Vue({
    router,
    data: {
        header: 'Pogo Admin',
    }
}).$mount('#app')

function getCss(callback) {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function() { 
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
            callback(null, xmlHttp.responseText)
    }
    xmlHttp.open("GET", "/admin/css", true);
    xmlHttp.send(null);
}