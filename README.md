# Gizo 
## Contribution
We welcome and appreciate contributions from anyone regardless of magnitude of fixes. Just make sure to have fun while at it!

If you'd like to contribute to gizo, start by forking the repo. Make sure to have [Go](https://golang.org/doc/install) and [Dep](https://golang.github.io/dep/docs/installation.html) installed. Clone this repo into `$GOPATH/src/github.com/gizo-network/` by running the command `git clone https://github.com/gizo-network/gizo.git`. Enter the directory and install the dependencies by running `cd $GOPATH/src/github.com/gizo-network/gizo && dep ensure`. Now that you've successfully done that, it's time to setup git origin and upstream remote url's so you can push straight to your forks rather than into `github.com/gizo-network/gizo`. To achieve this, run the following commands:
* `git remote set-url origin https://github.com/USERNAME/gizo.git` - make sure to replace USERNAME with your github username
* `git remote add upstream https://github.com/gizo-network/gizo.git`

Please make sure your contributions follow our coding guidelines:
* Stick to [git branching model](http://nvie.com/posts/a-successful-git-branching-model/)
    * Make pull requests from your feature or hotfix branches into the develop branch
* Push to the origin of your forked repo and make pull requests to gizo 
* Commit messages should be prefixed with the package(s) they modify.
    * E.g `core, helpers: example message`