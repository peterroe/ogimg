## OGIMG

A simple OG image capture service.

## Usage

`https://ogimg.peterroe.me?url=<encoded_url>`

Sample:

https://ogimg.peterroe.me?url=https%3A%2F%2Fgithub.com

## Self-host

```bash
$ wget -O docker-compose.yml https://raw.githubusercontent.com/peterroe/ogimg/main/docker-compose.yml
$ docker-compose up -d
```

Visit http://localhost:8888?url=https%3A%2F%2Fgithub.com
