# SVC

SVC is a set of components for building web services based on the concept of containers and middleware.

There are currently two containers.

- `rpc` for remote procedure calls. The preferred container for endpoints.
- `http` for standard http routing.

There will be another for a simple proxy.

## Why?

The core here is to be able to build web services in a consistent way, and to provide documentation generation for the
RPCs in order to drive automation for things like MCP servers.
