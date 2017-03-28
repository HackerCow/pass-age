#   pass-age
## Prerequisites
Requires ZX2C4's [pass](https://www.passwordstore.org/) and an initialized and git-tracked password store in `~/.password-store`
## Building
```
$ go build
```
## Installation
```
# go install
```
## Usage
An invocation with no parameters lists all passwords in the password store and their last-modified date.
```
$ pass-age
github.com/HackerCow was last changed 1 day 1 hour 24 minutes 31 seconds ago
```

An invocation with one parameter fetches the last-modified time of a single password entry
```
$ pass-age twitter.com/HackerCow
twitter.com/HackerCow was last changed 23 hours 40 minutes 4 seconds ago
```