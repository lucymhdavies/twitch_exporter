FROM grafana/grafana

ADD ./provisioning/ /etc/grafana/provisioning/

RUN grafana-cli plugins install grafana-clock-panel
