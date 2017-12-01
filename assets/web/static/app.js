const episodepublishform = {
    template: `<div>
    <h3>Publish Episode</h3>
    <form enctype="multipart/form-data" action="/admin/publish" method="post" class="publish">
        <label for="title">Episode Title</label>
        <input type="text" id="title" name="title">
        <label for="description">Episode Description</label>
        <textarea name="description" id="description" style="resize: none;" class="epdesc"></textarea>
        <label for="file">Media File</label>
        <input type="file" id="file" name="file">
        <label for="date">Publish Date</label>
        <input type="date" id="date" name="date"><br /><br />
        <input type="submit" value="Publish" class="button">
    </form>
</div>`
}

const message = {
    template: `<div><h3>{{ this.$route.params.message }}</h3></div>`
}

const userlist = {
    template: `<div>
    <router-link :to="\'users/new\'" tag="button">New</router-link>
    <table>
        <tr>
            <th>Username</th>
            <th>Email</th>
            <th></th>
        </tr>
        <tr v-for="item in items">
            <td>{{ item.username }}</td>
            <td>{{ item.email }}</td>
            <td>
                <router-link :to="\'user/\' + item.id" class="button">Edit</router-link>
            </td>
        </tr>
    </table>
</div>`,
    data() {
        return {
            loading: false,
            items: null,
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
            this.error = this.items = []
            this.loading = true

            get("/admin/listusers", (err, items) => {
                this.loading = false
                if (err) {
                    this.error = err.toString()
                } else {
                    var t = JSON.parse(items).reverse();
                    for (var i = t.length - 1; i >= 0; i--) {
                        this.items.push({
                            id: t[i].id,
                            username: t[i].username,
                            email: t[i].email,
                        })
                    }
                }
            })
        }
    }
}

const usernew = {
    template: `<div>
    <div>
        <h3>New User</h3>
        <form enctype="multipart/form-data" action="/admin/adduser" method="post">
        <label for="username">Username</label>
        <input type="text" id="username" name="username">
        <label for="email">Email</label>
        <input type="text" id="email" name="email">
        <label for="realname">Real Name</label>
        <input type="text" id="realname" name="realname">

        <label for="password">New Password</label>
        <input type="password" id="password" name="password">
        <label for="permissions">Permission Level</label>
        <select name="permissions">
            <option value="0">Publishing only</option>
            <option value="1">Publishing and Episode Management</option>
            <option value="2">Publishing, Episode and User management</option>
        </select>
        <br /><br />
        <input type="submit" class="button" value="Save"></form>
    </div>
</div>`
}

const useredit = {
    template: `<div>
    <div>
        <h3>Edit User</h3>
        <form enctype="multipart/form-data" action="/admin/edituser" method="post">
        <label for="username">Username</label>
        <input type="text" id="username" name="username" :value="user.username">
        <label for="email">Email</label>
        <input type="text" id="email" name="email" :value="user.email">
        <label for="realname">Real Name</label>
        <input type="text" id="realname" name="realname" :value="user.realname">

        <label for="newpw1">New Password</label>
        <input type="password" id="newpw1" name="newpw1">
        <label for="newpw2">Repeat New Password</label>
        <input type="password" id="newpw2" name="newpw2">
        <label for="oldpw">Old Password</label>
        <input type="password" id="oldpw" name="oldpw">
        <input name="id" id="id" :value="user.id" type="hidden">
        <br /><br />
        <input type="submit" class="button" value="Save" class="button"></form>
        <a v-bind:href="'/admin/deleteuser/'+user.id+''" class="button">Delete User</a>
    </div>
</div>`,
    data() {
        return {
            loading: false,
            user: null,
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
            this.error = this.user = []
            this.loading = true

            get("/admin/listusers", (err, items) => {
                this.loading = false
                if (err) {
                    this.error = err.toString()
                } else {
                    var t = JSON.parse(items)
                    for (var i = t.length - 1; i >= 0; i--) {
                        if (t[i].id == this.$route.params.id) {
                            this.user = {
                                id: t[i].id,
                                username: t[i].username,
                                email: t[i].email,
                                realname: t[i].realname
                            }
                        }
                    }
                }
            })
        }
    }
}

const episodemanagement = {
    template: `<div>
                <table style="width:100%">
                <tr>
                    <th>Title</th>
                    <th>URL</th>
                    <th></th>
                </tr>
                <tr v-for="item in items">
                    <td>{{ item.id }}: {{ item.title }}</td><td>{{ item.url }}</td><td><router-link class="button" :to="\'edit/\' + item.id">Edit</router-link></td>
                </tr>
                </table>
                </div>`,
    data() {
        return {
            loading: false,
            items: null,
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
            this.error = this.items = []
            this.loading = true

            get("/json", (err, items) => {
                this.loading = false
                if (err) {
                    this.error = err.toString()
                } else {
                    var t = JSON.parse(items).items
                    for (var i = t.length - 1; i >= 0; i--) {
                        console.log(i) 
                        this.items.push({
                            title: t[i].title,
                            url: t[i].url,
                            id: t[i].id
                        })
                    }
                }
            })
        }
    }
}

const episodeedit = {
    template: `<div>
    <div>
        <h3>Edit Episode</h3>
        <form enctype="multipart/form-data" action="/admin/edit" method="post">
        <label for="title">Episode Title</label>
        <input type="text" id="title" name="title" :value="episode.title">
        <label for="description">Episode Description</label>
        <textarea name="description" id="description" cols="100" rows="20" style="resize: none;">{{ episode.description }}</textarea>
        <label for="date">Publish Date</label>
        <input type="date" id="date" name="date" :value="episode.time">
        <input name="previousfilename" id="previousfilename" :value="episode.previousfilename" type="hidden">
        <input type="submit" class="button" value="Publish"></form>
    </div>
</div>`,
    data() {
        return {
            loading: false,
            episode: null,
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
            this.error = this.episode = {}
            this.loading = true

            get("/json", (err, items) => {
                this.loading = false
                if (err) {
                    this.error = err.toString()
                } else {
                    var t = JSON.parse(items).items
                    for (var i = t.length - 1; i >= 0; i--) {
                        if (t[i].id == this.$route.params.id) {
                        	time = t[i].date_published.split("T")
                        	var prev = time[0] + "_" + t[i].title
                            this.episode = {
                                title: t[i].title,
                                description: t[i].summary,
                                url: t[i].url,
                                id: t[i].id,
                                time: time[0],
                                previousfilename: prev
                            }
                        }
                    }
                }
                console.log(this.episode)
            })
        }
    }
}

const customcss = {
    template: `<div>
    <h3>Theme</h3>
    <form action="/admin/css" method="post" enctype="multipart/form-data">
        <textarea spellcheck="false" name="css" id="css" cols="120" rows="20" class="css">{{ css }}</textarea>
        <br /><br />
        <input type="submit" class="button" value="Submit" class="button">
    </form>
</div>`,
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

            get("/admin/css", (err, css) => {
                this.loading = false
                if (css == "{}") {
                    this.css = "You aren't allowed to edit this CSS!"
                } else {
                    if (err) {
                        this.error = err.toString()
                    } else {
                        this.css = css
                    }
                }
            })
        }
    }
}

const routes = [
    {path: '/', redirect: '/publish'},
    { path: '/publish', component: episodepublishform },
    { path: '/manage', component: episodemanagement },
    { path: '/theme', component: customcss },
    { path: '/edit/:id', component: episodeedit },
    { path: '/users/', component: userlist },
    { path: '/msg/:message', component: message },
    { path: '/user/:id', component: useredit },
    { path: '/users/new', component: usernew }
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

function get(url,callback) {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function() {
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
            if (xmlHttp.responseText == "Unauthorized") {
                callback(null, "{}")
            } else {
                callback(null, xmlHttp.responseText)
            }
        }
    }
    xmlHttp.open("GET", url, true);
    xmlHttp.send(null);
}