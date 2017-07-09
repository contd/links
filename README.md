# Links REST Service

This is a rest api service for managing saved links for personal use and to replace the storage currently used by the front end web app.

To run this in docker and persist the sqlite database, use the following once you've created an image from the `Dockerfile`:

```shell
docker run --name golinks -d -p 5555:5555 -v $GOPATH/src/githib.com/contd/links:/data contd/links
```

Then use the `docker-links.service` file to make the service autostart on boot:

```shell
sudo cp $GOPATH/src/github.com/contd/links/docker-links.service /etc/systemd/system/
sudo systemctl enable docker-links.service
sudo systemctl start docker-links.service
```

