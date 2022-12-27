{ pkgs, ... }:

{

  # https://devenv.sh/packages/
  packages = [ pkgs.git pkgs.clickhouse pkgs.ngrok ];

  enterShell = ''
    go install github.com/cosmtrek/air@latest
    go install -tags 'clickhouse' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  '';

  # https://devenv.sh/languages/
  languages.go.enable = true;

  processes.clickhouse.exec = "mkdir -p .devenv/state/clickhouse && cd .devenv/state/clickhouse && clickhouse server start";
}
