# sn

A [Simplenote](https://simplenote.com) static-syncing CLI client written in Go

## Table of Contents

1. [Introduction](#introduction)
2. [What this is not](#what-this-is-not)
3. [Installation](#installation)
4. [Usage](#usage)

## Introduction

Syncing textual data across devices is always difficult. I've never found a simple alternative for syncing notes across all my daily devices. I have tried many alternatives such as Evernote, Notion, iNotes, git, and others but none of them can do exactly what I'm looking for.

I chose Simplenote because of its focus on simplicity. It is made for one purpose and one purpose only: to take notes and sync them across devices. However, it is missing one piece: it does not provide command-line or client-side note syncing.

That's where this project comes in. The goal of this project is to readily download simplenote notes as text files and sync with edits to these files.

The program works by using the Simperium websocket API to connect to Simplenote buckets and retrieve and upload data and changes.

## What this is not

It is important to clarify that this syncing client is not:

* a fully-fledged Simplenote CLI
* an extendible library
* a markdown converter
* a websocket library

This project is entirely made by me and for me. I don't plan on making this a general-purpose program or library, but I'm always open to suggestions and patches.

## Installation

Compile and install the program.

```sh
make
$ make install
```

## Usage

Sn has two main functions: to download note changes from the server and upload local modifications to the server.

Downloading new notes or new changes is as simple as running `sn` to sync local changes with the Simperium server.

Local modifications can be uploaded to the server with `sn u`.

Run `sn h` for more usage details.
