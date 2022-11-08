# gmp-jwt (experimental)

GMP-JWT generates JWTs suitable for use authenticating to Google Maps Platform APIs.

Supports:

* [Routes API](https://developers.google.com/maps/documentation/routes)

The server (server/main.go) can be customized with the following flags:

|flag||
|------|-|
| `--cors_origins` | Value to set for the 'Access-Control-Allow-Origin' header.  Use * for wildcard, which is dangerous in production. |
| `--enable_appcheck`| If set, requests must a valid token from app check in the `X-Firebase-AppCheck` header. |
| `--token_duration` | Duration a generated token is valid for (default 30m0s) |

[![Run on Google Cloud](https://deploy.cloud.run/button.svg)](https://deploy.cloud.run)

------

Tokens are generated using an available Service Account via [Application Default Credentials](https://cloud.google.com/docs/authentication/provide-credentials-adc).
Practically, this means the either provided automatically when run on Cloud Run, GCE, etc
or manually provided via the `GOOGLE_APPLICATION_CREDENTIALS` ENV variable.