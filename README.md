# Telegraf http output

## Instalacja
1. [Instalujemy go](https://golang.org/doc/install).

1. Należy pamiętać o [modyfikacji zmiennej $PATH](https://golang.org/doc/code.html#GOPATH). W przeciwnym przypadku podczas kompilacji będziemy otrzymywać komunikat, że program `gdm` nie został odnaleziony.

1. Pobieramy źródła:
    ```
    go get -d github.com/influxdata/telegraf
    go get -d github.com/gilek/telegraf-output-http/plugins/outputs/http
    ```
1. Przechodzimy do katalogu:
    ```
    cd ~/go/src/github.com/influxdata/telegraf
    ```
1. Modyfikujemy plik `plugins/outputs/all/all.go` dodając nowy wpis o adresie:
    ```
    github.com/gilek/telegraf-output-http/plugins/outputs/http
    ```
1. Wykonujemy `make`.

1. Przykładowa konfiguracja jest w pliku `config.toml`.
