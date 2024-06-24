FROM golang:alpine
#renombrar
COPY ./api.go ./api-tf.go
#port
EXPOSE 9000
# run api-svc001.go
CMD ["go", "run", "api-tf.go"]
 
 #FROM golang:alpine
 #RUN mkdir /app
 #ADD . /app/
 #WORKDIR /app
 #RUN go build -o svc .
 #CMD ["./svc"]