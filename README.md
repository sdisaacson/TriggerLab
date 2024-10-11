<p align="center">
  <img src="https://hsto.org/webt/mn/fz/q-/mnfzq-lgdnbmv-3xv-1qm6gn82e.png" alt="Logo" width=128" />
</p>

# TriggerLab

[![Release version][badge_release_version]][link_releases]
![Project language][badge_language]
[![Build Status][badge_build]][link_build]
[![Release Status][badge_release]][link_build]
[![Image size][badge_size_latest]][link_docker_hub]
[![License][badge_license]][link_license]

This application allows you to test and debug Webhooks and HTTP requests using unique (random) URLs. You can customize the response code, `content-type` HTTP header, response content and set some delay for the HTTP responses. The main idea is viewed [here](https://github.com/fredsted/webhook.site).

<p align="center">
  <img src="https://user-images.githubusercontent.com/7326800/201918441-ec70a826-48dc-4bb0-af27-d194ea89a4ad.gif" alt="screencast" />
</p>

This application is written in GoLang and works very fast. It comes with a tiny UI (written in `Vue.js`), which is built in the binary file, so you don't need any additional assets for the application using. Websockets are also used for incoming webhook notifications in the UI - you don't need any 3rd party solutions (like `pusher.com`) for this!

### 🔥 Features list

- Liveness/readiness probes (routes `/live` and `/ready` respectively)
- Can be started without any 3rd party dependencies
- Metrics in prometheus format (route `/metrics`)
- Built-in tiny and fast UI, based on `vue.js`
- Multi-arch docker image, based on `scratch`
- Unprivileged user in docker image is used
- Well-tested and documented source code
- Built-in CLI health check sub-command
- Recorded request binary view using UI
- JSON/human-readable logging formats
- Customizable webhook responses
- Built-in Websockets support
- Low memory/cpu usage
- Free and open-source
- Ready to scale


### 🗃 Storage

At the moment 2 types of data storage are supported - **memory** and **redis server** (flag `--storage-driver`).

The **memory** driver is useful for fast local debugging when recorded requests will not be needed after the app stops. The **Redis driver**, on the contrary, stores all the data on the redis server, and the data will not be lost after the app restarts. When running multiple app instances (behind the load balancer), it is also necessary to use the redis driver.

### 📢 Pub/sub

Publishing/subscribing are used to send notifications using WebSockets, and it also supports 2 types of driver - **memory** and **redis server** (flag `--pubsub-driver`).

For multiple app instances redis driver must be used.
