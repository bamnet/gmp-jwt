# gmp-jwt (experimental)

GMP-JWT generates JWTs suitable for use authenticating to Google Maps Platform APIs.

Supports:

* [Routes API](https://developers.google.com/maps/documentation/routes)

------

Tokens are generated using an available Service Account via [Application Default Credentials](https://cloud.google.com/docs/authentication/provide-credentials-adc).
Practically, this means the either provided automatically when run on Cloud Run, GCE, etc
or manually provided via the `GOOGLE_APPLICATION_CREDENTIALS` ENV variable.