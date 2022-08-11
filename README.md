
# Bottle

The friendliest blogging engine


## Used By

- [My blog](https://u8.quest)


## Installation

If you're on a *nix system:
```sh
$ curl https://bottle.quest/install.sh | bash
```
    
## Features

- Built-in production webserver
- Live edit your blog posts. No more restarting your server.
- Easy to customize theme system.
- Works on any platform Go is supported on


## Deployment

To create a new blog using `bottle`, first create your site's directory

```bash
mkdir blog
cd blog
```

Use `bottle` to create your blog's structure. 
```bash
bottle init
```

Customize your blog's metadata by editing `bottle.yaml`
```yaml
title: 99 Bottles of beer
author: Brewmaster
description: A blog by Brewmaster
keywords:
    - blog
    - bottle
    - seo
    - brewing
theme: default
url: "http://localhost:8080"
locale: en_us
fb_user_id: ""
fb_app_id: ""
tw_pub_handle: ""
tw_author_handle: ""
```

Create a new post by creating a new markdown file in ./posts or by running 
```bash
bottle new post "your_post_slug"
```

Everything is YAML! Customize your post's metadata.  
```yaml
---
title: My new post
subtitle: A blank canvas--full of adventure
author: Brewmaster
description: ""
keywords: []
site_img: ""
body: ""
slug: test
publish_date: 2022-07-20T22:49:35.876486611-04:00
og_video: ""
section: ""
modifiedTime: 0001-01-01T00:00:00Z
expirationTime: 0001-01-01T00:00:00Z
tw_author_handle: ""
tw_preview_image: ""
tw_vid_aud_player: ""
---
```
Don't worry about filling all those fields. Bottle's smart enough to fill in sensible defaults.
All you need to publish a blog post is a Title! 
```yaml
---
title: Hello World!
---
```

## Serving your Website

`bottle` comes pre-packaged with a [production level webserver](https://gofiber.io).
The server comes with logging, compression, and some other optimizations on by default.
```bash
bottle serve
```

If you want to run your server across multiple processes, start it with:
```bash
bottle serve -j
```

More options are coming as bottle gets developed to allows pros to customize their site as heavily as needed.
## Contributing

Contributions are always welcome!

See `contributing.md` for ways to get started.

Please adhere to this project's `code of conduct`.


## License

[MIT](https://choosealicense.com/licenses/mit/)

