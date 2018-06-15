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

Prometheus will be available on http://localhost:9090

For example, to view all live viewers:

http://localhost:9090/graph?g0.range_input=30m&g0.expr=lmhd_twitch_stream_viewers&g0.tab=0



## TODO

Migrate to the new (Helix) Twitch API

Helix should already provide most, if not all, of the necessary functionality, so migrate to that at some point

https://github.com/nicklaw5/helix
