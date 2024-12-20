## OGIMG

A simple OG image capture service(cache with redis).

## Usage

**Basic usage**

`https://ogimg.peterroe.me?url=<encoded_url>`

A few examples:

* GitHub: https://ogimg.peterroe.me?url=https%3A%2F%2Fgithub.com
* YouTube: https://ogimg.peterroe.me/?url=https%3A%2F%2Fyoutube.com
* Instagram: https://ogimg.peterroe.me/?url=https%3A%2F%2Finstagram.com

**Description**

`https://ogimg.peterroe.me/desc?url=<encoded_url>`

Return the description of the website. Include `title`, `description` and `logo`.

A few examples:

* GitHub: https://ogimg.peterroe.me/desc?url=https%3A%2F%2Fgithub.com
* YouTube: https://ogimg.peterroe.me/desc?url=https%3A%2F%2Fyoutube.com
* Instagram: https://ogimg.peterroe.me/desc?url=https%3A%2F%2Finstagram.com


## Self-hosted

Easy to self-host, just run the following command:

```bash
$ wget -O docker-compose.yml https://raw.githubusercontent.com/peterroe/ogimg/main/docker-compose.yml
$ docker-compose up -d
```

Then visit http://localhost:8888?url=https%3A%2F%2Fgithub.com
