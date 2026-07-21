FROM texlive/texlive:latest
WORKDIR /app
COPY --from=golang:1.24-alpine /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:/home/rlevidev/.mimocode/bin:/home/rlevidev/.local/bin:/home/rlevidev/.opencode/bin:/home/rlevidev/.local/bin:/home/linuxbrew/.linuxbrew/bin:/home/linuxbrew/.linuxbrew/sbin:/home/rlevidev/.local/share/zinit/polaris/bin:/usr/lib/jvm/java-25-openjdk/bin:/usr/local/bin:/usr/bin:/var/lib/snapd/snap/bin:/usr/local/go/bin:/home/rlevidev/go/bin"
COPY . .
RUN go build -o main .
CMD ["./main"]
