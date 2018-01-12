# Contributing

Generally, any contributions to Pogo are more than welcome, but we'd like it if you follow a couple guidelines. We'll also point out a couple of tricks for ease of use.

First, fork the repository and clone it locally


```
git clone git@github.com:your-username/pogo.git
```

Set up Go and install the dependencies

```
cd pogo
go get github.com/tools/godep
godep restore
```

Then make your changes. If you use Sublime Text 3, please check out [our snippets](https://gist.github.com/gmemstr/60831109f0ae6c40861c1751a367524e) to add some shortcuts to make your life easier. 

The platform is divided into two main parts: The main Go app, which does everything from generating RSS to serving up webpages, and the "DeV" web interface, which implements a basic frontend for testing - if you want to contribute to the frontend, check out [gmemstr/pogo-vue](https://github.com/gmemstr/pogo-vue). 

Once you've made your changes, please make sure the app can build

```
go build
```

Once you've verified your addition works, push to your repository and [create a pull request](https://github.com/gmemstr/pogo/compare). During the review your PR will also pass through our TravisCI testing.
