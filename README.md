# sn

A [Simplenote](https://simplenote.com) syncing CLI client written in Go

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [What this is not](#what-this-is-not)
4. [Installation](#installation)
5. [Usage](#usage)

## Introduction

Syncing textual data across devices is always difficult. I've never found a simple alternative for syncing notes across all my daily devices. I have tried many alternatives such as Evernote, Notion, iNotes, git, and others but none of them can do exactly what I'm looking for.

I chose Simplenote because of its focus on simplicity. It is made for one purpose and one purpose only: to take notes and sync them across devices. However, it is missing one piece: it does not provide command-line or client-side note syncing.

That's where this project comes in. The goal of this project is to readily download simplenote notes as text files and sync with edits to these files.

The program works by using the [Simperium websocket API](https://simperium.com/docs/websocket) to connect to Simplenote note buckets and retrieve and upload changes.

## Features

* Authenticate and download a user account's notes in markdown format
* Upload changes to notes
* Automatically open the notes directory with $EDITOR
* Allow refetching for notes
* Clear and delete all notes

## What this is not

It is important to clarify that this Simplenote client is not:

* an official Simplenote CLI
* an extendible Simperium library
* a real time note update system
* a markdown converter
* a websocket library

This project is created entirely by me and for me. I don't intend to make this a general-use program but I'm always open to patches, and I welcome anyone who would like to use this project with me.

## Installation

Compile and install the program.

```sh
make
$ make install
```

If you run NixOS you can use [this derivation from my dotfiles for reference](https://github.com/bossley9/dotfiles/blob/7b7c1d19ba1e1f4cd5fa62a55bc7c553abc1d17c/user/packages/sn.nix).

## Usage

Run `sn h` for more usage details.
