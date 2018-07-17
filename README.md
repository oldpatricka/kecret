# kecret - The Kubernetes Secret Editor

This little app helps you with editing secrets when working with Kubernetes.
It decodes all base64 encoded values in your secret and automatically opens
$EDITOR, then automatically re-encodes them when you save your file and exit
the editor. Handy!

## Installation

```
go get github.com/oldpatricka/kecret
```

## Usage

Simply use kecret as your editor to edit your secret file like so:

```
kecret mysecret.yaml
```

Your editor will display the decoded secret values which you can edit with ease.
When you save, your secret file will automatically be re-encoded.

## TODO

* Add support for using kecret with `kubectl edit secret/whatever`
* Some sanity checks
* Passthrough when used as KUBE_EDITOR

## Motivation

Why? Well, my regular workflow when working with secrets in Kubernetes is to
keep [SOPS](https://github.com/mozilla/sops) encrypted files in my source
repository, and decrypt them and apply them every time I make a change to the
secrets. This means my full workflow whenever I need to edit a secret is:

```
sops --decrypt mysecret.enc.yaml > mysecret.yaml
$EDITOR mysecret.yaml
<copy secret>
pbpaste | base64 -D
<copy decoded secret into another editor>
<edit secret>
<copy edited secret>
pbpaste | base64 | pbcopy
$EDITOR mysecret.yaml
<paste secret in the right location>
<save>
sops --encrypt mysecret.yaml > mysecret.enc.yaml
```

This is pretty annoying, and kecret helps with this a bit, by changing this
workflow to:

```
sops --decrypt mysecret.enc.yaml > mysecret.yaml
kecret mysecret.yaml
<edit secret>
<save>
sops --encrypt mysecret.yaml > mysecret.enc.yaml
```

I think this is much better.
