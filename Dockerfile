#buiding minimal docker image from scratch for running GO Application
FROM scratch
ADD main /
ADD config/config.json /config/config.json
ADD config/config-dev.json /config/config-dev.json
CMD ["/main"]