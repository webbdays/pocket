## client Fisherman Stake

Stake a node in the network. Custodial stake uses the same address as operator/output for rewards/return of staked funds.

### Synopsis

Stake the node into the network, making it available for service.

Will prompt the user for the *fromAddr* account passphrase. If the node is already staked, this transaction acts as an *update* transaction.

A node can update relayChainIDs, serviceURI, and raise the stake amount with this transaction.

If the node is currently staked at X and you submit an update with new stake Y. Only Y-X will be subtracted from an account.

If no changes are desired for the parameter, just enter the current param value just as before.

```
client Fisherman Stake <fromAddr> <amount> <relayChainIDs> <serviceURI> [flags]
```

### Options

```
  -h, --help         help for Stake
      --pwd string   passphrase used by the cmd, non empty usage bypass interactive prompt
```

### Options inherited from parent commands

```
      --path_to_private_key_file string   Path to private key to use when signing (default "./pk.json")
      --remote_cli_url string             takes a remote endpoint in the form of <protocol>://<host> (uses RPC Port) (default "http://localhost:50832")
```

### SEE ALSO

* [client Fisherman](client_Fisherman.md)	 - Fisherman actor specific commands

###### Auto generated by spf13/cobra on 9-Nov-2022
