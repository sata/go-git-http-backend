## About

`go-git-http-backend` provides HTTP interface for the
[`go-git`](https://github.com/go-git/go-git) server side
implementation in a similar way as what the utility [git
http-backend](https://git-scm.com/docs/git-http-backend) does for
`git`. However, in it's current version it only supports clone and
push operations.

## Motivation

To provide an easy way for writing tests in Go for any tooling that
interacts with a remote Git repository.

In my case, it was to test a GitOps promotion tool behaviour was
correct. As in, the actual artefact of running the promotion tool,
resulted in producing commits and Git blobs that matched our
expectations of both path and content, while still testing end to end
of the tool, without any mocks or test interfaces over client or
storage.

## Development

This project uses [Mage](magefile.org) as a build tool. Install it
following [installation steps](https://github.com/magefile/mage#installation).

List available targets by running `mage`. 

## Limitations

- The project only supports Smart protocol, see
[http-protocol](https://github.com/git/git/blob/master/Documentation/technical/http-protocol.txt)
- Only support clone & push operations so far.

## Resources
useful documentation to understand the Git protocol and the transfer
protocol over HTTP.

- https://github.com/git/git/blob/master/Documentation/technical/http-protocol.txt
- https://github.com/git/git/blob/master/Documentation/technical/pack-protocol.txt
- https://github.com/git/git/blob/master/Documentation/technical/protocol-common.txt
