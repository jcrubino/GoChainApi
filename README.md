## Synopsis

This is an unofficial golang client api for the bitcoin blockchain service from Chain.com
It is a work in progress that needs polishing as well as the webhooks api implemented.
Currently all data is returned as a string but the json deserialization structs are included.
Next on my to do is json deserialization, then finish the web hooks and finally reduce the dependecies to zero by creating a built in config processing.

Dependencies: code.google.com/p/gcfg

Anybody and everybody is welcome to submit fixes and features and a very straight forward Go project to work on for tire kickers.

## Code Example

For now the code is single file main package with examples at the end.
Will change over to proper package convention and seperate json structs into seperate files when ready.

## Motivation

Chain.com is a solid bitcoin blockchain service.
That said anyone choosing to use golang should look at Conformal LLC code at github.com/conformal where bitcoin node and wallet code can be found in pure golang.  Ideally this package becomes a local and remote api capable for btcd and chain.

## Installation

git clone

create a config.cfg file with your chain credentials and settings
example

```
[auth]
key="yourchainkey"
secret="yourchainsecret"

[mode]
network="bitcoin" # or testnet3
verbose=true # or false; no quotes! Outputs log info from individual api calls, yet to be prettyfied
```

<b>Send txn has not been tested and all tests use bitcoin network.  Will change over to testnet3 soon<b>

## API Reference

Naming conventions follow the node.js api conventions mostly

## Tests

see end of file

## Contributors

jcrubino, YourUserNameHere

## License

MIT liscense