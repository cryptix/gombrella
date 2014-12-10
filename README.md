# extract bookmarks from raindrop.io

the api documentation is declared as [`soon`](https://raindrop.io/pages/dev) since they launched and apperantly they [want to charge](https://raindrop.io/static/love) now for backing up your content.

I respect that they want to cover their costs with such a model and not mine out their users data but I don't want to pay to get my data back. So I hacked this together, you mainly just have to do 2 requests to get the list of collection ids. you can do so with curl as well.


## Requests

The system is build with express and uses JSON for everything

### Login
`https://raindrop.io/api/auth/login` 

Body:
```
{
"email": "...",
"password": "..."
}
```

### Collections and Bookmarks
`https://raindrop.io/api/collections` for the List of IDs of collections

See [types.go](https://github.com/cryptix/gombrella/blob/master/types.go) for a description of the fields in the response.

## Usage
```
gumbrealla <email>
```


## Installation

if you have go installed you can just run:

`go get github.com/cryptix/gombrella`

Pre-build binaries are hosted at [bintray.com](https://bintray.com/cryptix/golang/gombrella)
