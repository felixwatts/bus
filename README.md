This is a very simple message bus server. It will eventually support PUBLISH/SUBSCRIBE and REQUEST/RESPONSE. It uses a very simple text-based TCP protocol. 
subscription keys may contain wildcards:

* `SB 1 stocks/NASDAQ/GOOG/latest` subscribe to key 'stocks/GOOG/latest'
* `SB 2 stocks/NASDAQ/*/latest` subscribe to all keys starting with 'stocks/NASDAQ' followed by a single field followed by 'latest'
* `PB 3 stocks/**/state closed` publish 'closed' as the value for all keys starting with 'stocks' and ending with 'state'

Implemented Features

- Subscribe to a key
- Wildcards in subscription keys
- Double wildcards in subscription keys
- Publish to a key
- Wildcards in publish keys
- Double wildcards in publish keys
- Publisher can CLAIM an area of the keyspace. This means no one else can claim it or publish to it and the publisher gets informed when subscribers subscribe within its space

Future Features

- Register as a METHOD provider
- Request a method and get a response

Possibly later

- Decentralisation



