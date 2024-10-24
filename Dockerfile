FROM ubuntu:22.04 as build

WORKDIR /cm-explorer
COPY . .

WORKDIR /build
RUN mkdir -p bin
RUN cp -r /cm-explorer/scripts/* bin/
RUN cp -r /cm-explorer/configs/ .

FROM ubuntu:22.04

WORKDIR /chainmaker-explorer-backend
COPY --from=build /build/ .

EXPOSE 9999

WORKDIR /chainmaker-explorer-backend/bin
ENTRYPOINT ["./chainmaker-browser.bin", "-config" ]
CMD ["../configs/"]
