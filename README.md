# wsmux

A mux for websockets, utilizing pub/sub.

In a modern architecture, the backend for a web app consists of dozens,
hundreds, or even more microservices, glued together by something like
GraphQL. That doesn't play well with websockets, which stay attached to
just one server over over the course of their lifetime. We obviously
don't want just one monolithic server for our backend, so we need a way
to let microservices send messages through a websocket that the
microservices themselves do not control.

`wsmux` is that missing piece. If you spin up a `wsmux`-based service,
it will keep track of the websockets and which users they belong to.
When a user establishes a websocket with `wsmux`, `wsmux` will subscribe
to a topic based on the user's unique identifier. Microservices can then
publish messages to those topics to asynchronously notify a user via
your web's front-end.

You can spin up multiple `wsmux` instances for load balancing. Since
each instance subscribes to the user topics for its websockets, all of
them will respond to published messages for those topics.
