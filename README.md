# Prometheus Exporter for the Twitch.tv API

Proof of Concept, using (deprecated) v5 (Kraken) Twitch API

![Combination of Prometheus logo and Kappa emote from Twitch](images/prometheus-kappa.png)


## Usage

### Local

Copy `example.env` to `.env`, and modify as necessary

```
go run *.go
```

### Docker

Modify `environment` section of `my_metrics` in `docker-compose.yml`.
(see `example.env` for an explanation of each environment variable)

```
docker-compose up
```



## TODO

Migrate to the new (Helix) Twitch API

Helix should already provide most, if not all, of the necessary functionality, so migrate to that at some point

https://github.com/nicklaw5/helix
