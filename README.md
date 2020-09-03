check_sophos_central
====================

Check the status of alerts and endpoints over the API of the Sophos Central cloud service.

The plugin currently checks the state of all alerts and endpoints within a tenant, you need to supply API Token
(ID and secret) for a single tenant.

## Usage

```
Arguments:
      --client-id string       API Client ID (env:SOPHOS_CLIENT_ID)
      --client-secret string   API Client Secret (env:SOPHOS_CLIENT_SECRET)
      --show-all               List all non-ok endpoints
      --api string             API Base URL (default "https://api.central.sophos.com")
  -t, --timeout int            Abort the check after n seconds (default 30)
  -d, --debug                  Enable debug mode
  -v, --verbose                Enable verbose mode
  -V, --version                Print version and exit
```

## Example

```
$ ./check_sophos_central --client-id efce870a-6c53-4a6b-8c49-864894b9d8ee --client-secret thatwouldbeagoodjoke
CRITICAL - no alerts - endpoints: 2 good, 3 bad, 6 suspicious

## Endpoints
bad: HOST1, HOST2, HOST6
suspicious: HOST11, HOST12, HOST13, HOST14, HOST15, ...
| 'alerts'=0 'alerts_high'=0 'alerts_medium'=0 'alerts_low'=0 'endpoints_total'=11 'endpoints_good'=2 'endpoints_bad'=3 'endpoints_suspicious'=6 'endpoints_unknown'=0
```

## API Documentation

Full API documentation is available at [developer.sophos.com](https://developer.sophos.com/intro).

## License

Copyright (c) 2020 [NETWAYS GmbH](mailto:info@netways.de) \
Copyright (c) 2020 [Markus Frosch](mailto:markus.frosch@netways.de)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see [gnu.org/licenses](https://www.gnu.org/licenses/).
