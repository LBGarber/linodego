module github.com/linode/linodego

require (
	github.com/go-resty/resty/v2 v2.1.1-0.20191201195748-d7b97669fe48
	github.com/google/go-cmp v0.5.8
	golang.org/x/oauth2 v0.8.0
)

go 1.16

retract v1.0.0 // Accidental branch push
