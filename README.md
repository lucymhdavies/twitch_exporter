# Prometheus Exporter for the Twitch.tv API

Proof of Concept, using (deprecated) v5 (Kraken) Twitch API

![Combination of Prometheus logo and Kappa emote from Twitch](images/prometheus-kappa.png)


## Usage

Copy `example.env` to `.env`, and modify as necessary

### Local

Run locally with:

```
go run *.go
```

### Docker

Launch the exporter, as well as a prometheus instance with:

```
docker-compose up
```

Prometheus will be available, on [http://localhost:9090](http://localhost:9090)

The docker-compose file also comes with Grafana, with a pre-configured dashboard, on [http://localhost:3000](http://localhost:3000)

![screenshot of Grafana dashboard](images/grafana.png)




## TODO

Migrate to the new (Helix) Twitch API

Helix should already provide most, if not all, of the necessary functionality, so migrate to that at some point

https://github.com/nicklaw5/helix
