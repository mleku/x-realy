# database

## index key scheme

this scheme is designed to be stable and semi-human-readable, and use two bytes as the key prefix in most cases.

all keys further contain the database serial (sequence number) as the last 8 bytes

- `ev` - the event itself, encoded in binary format


- `cf` - a free form configuration JSON


- `id` - event id - truncated 8 bytes hash
  
  these override any other filter


- `fi` - full index: full event id, pubkey truncated hash, kind and created_at, enabling identifying and filtering search results to return only the event id of a match while enabling filtering by timestamp and allowing the exclusion of matches based on a user's mute list


- `pk` - public key - truncated 8 byte hash of public key


- `pc` - public key, created at - varint encoded (ltr encoder)

  these index all events associated to a pubkey, easy to pick by timestamp


- `ca` - created_at timestamp - varint encoded (ltr encoder)

  these timestamps are not entirely reliable but a since/until filter these are sequential


- `fs` - index that stores the timestamp when the event was received

  this enables search by first-seen


- `ki` - kind, created_at - 2 bytes kind, varint encoded created_at

  kind and timestamp - to catch events by time window and kind


- `ta` - kind, pubkey, hash of d tag (a tag value)

  these are a reference used by parameterized replaceable events


- `te` - event id - truncated 8 bytes hash

  these are events that refer to another event (root, reply, etc)


- `tp` - public key - truncated 8 bytes hash

  these are references to another user


- `tt` - hashtag - 8 bytes hash of full hashtag

  this enables fast hashtag searches


- `t` - tag for other letters (literally the letter), 8 bytes truncated hash of value

  all other tags, with a distinguishable value compactly encoded


- `t-` - 8 bytes hash of pubkey

  these are protected events that cannot be saved unless the author has authed, they can't be broadcast by the relay either


- `t?` - 8 bytes hash of key field - 8 bytes hash of value field

  this in fact enables search by other tags but this is not exposed in filter syntax


- `fw` - fulltext search index - the whole word follows, serial is last 8 bytes

  when searching, whole match has no prefix, * for contains ^ for prefix $ suffix


- `la` - serial, value is last accessed timestamp


- `ac` - serial, value is incremented counter of accesses

  increment at each time this event is matched by other indexes in a result
