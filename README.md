# sn

A [Simplenote](https://simplenote.com) downloading CLI client written in Go

## Table of Contents

1. [Introduction](#introduction)
2. [What this is not](#what-this-is-not)
3. [Installation](#installation)
4. [Usage](#usage)
5. [Roadmap](#roadmap)

## Introduction

Syncing textual data across devices is always difficult. I've never found a simple alternative for syncing notes across all my daily devices. I have tried many alternatives such as Evernote, Notion, iNotes, git, and others but none of them can do exactly what I'm looking for.

I chose Simplenote because of its focus on simplicity. It is made for one purpose and one purpose only: to take notes and sync them across devices. However, it is missing one piece: it does not provide command-line or client-side note syncing.

That's where this project comes in. The goal of this project is to readily download simplenote notes as text files and sync with changes to these files.

The program works by using the Simperium websocket API to connect to Simplenote buckets and retrieve data and changes.

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

Run `sn h` for usage details.

## Roadmap

These are things that are "nice to haves" but I haven't or do not plan on implementing in the near future.

* uploading local changes - unfeasible due to the difficult nature of tracking the diff between change versions
