FROM golang:1.14 as build

RUN apt-get update && apt-get install -y ninja-build git

# TODO: Змініть на власну реалізацію системи збірки
WORKDIR /go/src
RUN git clone https://github.com/AnastasiaYarema/design-practice-2
WORKDIR design-practice-2/
RUN go get -u github.com/roman-mazur/bood/cmd/bood
WORKDIR build/
RUN bood

WORKDIR /go/src/practice-3
COPY . .

# RUN CGO_ENABLED=0 bood
RUN ls
RUN CGO_ENABLED=0 /go/src/design-practice-2/build/out/bin/bood

# ==== Final image ====
FROM alpine:3.11
WORKDIR /opt/practice-3
COPY entry.sh ./
COPY --from=build /go/src/practice-3/out/bin/* ./
ENTRYPOINT ["/opt/practice-3/entry.sh"]
CMD ["server"]
