# uAdmin the Golang Web Framework

Easy to use, blazing fast and secure.

[![go report card](https://goreportcard.com/badge/github.com/sergeyglazyrindev/go-monolith "go report card")](https://goreportcard.com/report/github.com/sergeyglazyrindev/go-monolith)
[![GoDoc](https://godoc.org/github.com/sergeyglazyrindev/go-monolith?status.svg)](https://godoc.org/github.com/sergeyglazyrindev/go-monolith)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-3%25-brightgreen.svg?longCache=true&style=flat)</a>
![sqlite build status](https://github.com/sergeyglazyrindev/go-monolith/actions/workflows/sqlite.ci.yml/badge.svg)
![postgres build status](https://github.com/sergeyglazyrindev/go-monolith/actions/workflows/postgres.ci.yml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://github.com/sergeyglazyrindev/go-monolith/blob/master/LICENSE)

Originally open sourced by [IntegrityNet Solutions and Services](https://www.integritynet.biz/)
And then fully rewritten by @sergeyglazyrindev.

For Documentation:

- [Please check](https://go-monolithdocs.sergeyg.me/)

Reach us at:

- [gophers.slack.com #uadmin](https://gophers.slack.com/messages/uadmin/)
- [discord #uadmin](https://discord.gg/kADzHWatSj)
- [telegram uadmin_development](https://t.me/joinchat/VzgmokqjF7s4Nzk0)

[join gophers.slack.com](https://join.slack.com/t/gophers/shared_invite/zt-fajz7jh3-2cpkmFU~hQb8d5LmOCnhfQ)

Social Media:

## Screenshots

### Dashboard Menu

![Dashboard](https://github.com/sergeyglazyrindev/go-monolith-docs/raw/master/assets/go-monolithdashboard.png)
&nbsp;

### Log

![Log](https://github.com/sergeyglazyrindev/go-monolith-docs/raw/master/assets/log.png)
&nbsp;

### Login Form

![Login Form](https://github.com/sergeyglazyrindev/go-monolith-docs/raw/master/tutorial/assets/loginform.png)
&nbsp;

## Features

- AB Testing System
- API Configuration
- Approval System
- Authentication and Permissions
- Clean and sharp UI
- Dashboard customization
- Error Handling
- Export to Excel
- Form and List customization
- Image Cropping
- IP address and port configuration
<!-- - Log feature that keeps track of many things in your app
- Metric System -->
<!-- - Multilingual translation -->
- Full Gorm integration
- Pretty good security features (SSL, 2-Factor Authentication, Password Reset, Hash Salt Unique for Each user, Database Encryption)
- Public access to media
- Self relation of foreign key/many2many
- Sending an email from your app by establishing an email configuration
- System settings which can be used system wide to keep application settings
- Tag support for fields
<!-- - Translation files preloading -->
- Validation for user input
- Webcam support on image and file fields

## Minimum requirements

| Operating System                   |                Architectures              |                                Notes                                                |
|------------------------------------|-------------------------------------------|-------------------------------------------------------------------------------------|
| FreeBSD 10.3 or later              |  amd64, 386                               | Debian GNU/kFreeBSD not supported                                                   |
| Linux 2.6.23 or later with glibc   |  amd64, 386, arm, arm64, s390x, ppc64le   | CentOS/RHEL 5.x not supported. Install from source for other libc.                  |
| macOS 10.10 or later               |  amd64                                    | Use the clang or gcc<sup>†</sup> that comes with Xcode<sup>‡</sup> for cgo support. |
| Windows 7, Server 2008 R2 or later |  amd64, 386                               | Use MinGW gcc<sup>†</sup>. No need for cygwin or msys.                              |

- <sup>†</sup> A C compiler is required only if you plan to use cgo.
- <sup>‡</sup> You only need to install the command line tools for Xcode. If you have already installed Xcode 4.3+, you can install it from the Components tab of the Downloads preferences panel.

### Hardware

- RAM - minimum 256MB
- CPU - minimum 2GHz

### Software

- Go Version 1.16 or later

## Installation

```bash
$ go get -u github.com/sergeyglazyrindev/go-monolith/...
```

To test if your installation is fine, run the `go-monolith` command line:

Get full documentation online:
https://go-monolithdocs.sergeyg.me/
```

## Your First App

[Check it out here]https://go-monolithdocs.sergeyg.me/docs/intro
