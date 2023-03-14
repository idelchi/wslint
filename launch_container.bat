cd /D "%~dp0"
set "PWD=.."
docker rm -f devenv
docker compose run --name devenv --rm --service-ports devenv
