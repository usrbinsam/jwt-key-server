# JWT Key Server
Software key licensing server based on JWTs. Licenses are just permissions, so they fit the JWT model. This project
is a reboot of the original Python-Flask text-based licensing server, which can be found on the **legacy** branch.

[![Unit Tests](https://github.com/usrbinsam/mini-key-server/actions/workflows/main.yml/badge.svg?branch=beta)](https://github.com/usrbinsam/jwt-key-server/actions/workflows/main.yml)


## JWT as a License Principle

| JWT Standard Claim | Application License Use                 |
|--------------------|-----------------------------------------|
| Subject            | The key itself                          |
| Issuer             | The application the key activates       |
| Audience           | The developer who wrote the application |

Additionally, private claims and scopes can be used to enable sub features of licenses, such as trial periods.

## Activation Protocols

The license model should be accessible through a variety of transports, the most widely accessible being HTTP. JWT can
be re-used within the same project as an activation protocol for clients to remotely authenticate. The audience claim
allows for multi-tenancy developers with parameters per-developer as well as per-key and per-application.
