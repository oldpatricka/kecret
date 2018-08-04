# kecret - The Kubernetes Secret Editor

This little app helps you with editing secrets when working with Kubernetes.
It decodes all base64 encoded values in your secret and automatically opens
$EDITOR, then automatically re-encodes them when you save your file and exit
the editor. Handy!

## Installation

```
go get github.com/oldpatricka/kecret
```

Or [download](https://github.com/oldpatricka/kecret/releases) latest release, and install to /usr/local/bin:

```
export PLATFORM="`uname`_`uname -m`"
curl -L https://github.com/oldpatricka/kecret/releases/download/v0.1.1/kecret_0.1.1_$PLATFORM.tar.gz -o kecret.tar.gz
tar xzf kecret.tar.gz
mv kecret /usr/local/bin
```

## Usage

Simply use kecret as your editor to edit your secret file like so:

```
kecret mysecret.yaml
```

Your editor will display the decoded secret values which you can edit with ease.
When you save, your secret file will automatically be re-encoded.

If you would like to edit your secret directly on a live Kubernetes system,
use:

```
KUBE_EDITOR=kecret kubectl edit secrets/mysecret
```

You can also set KUBE_EDITOR in your shell to kecret, and if you attempt to edit
something that isn't a secret, it will simply pass through to your EDITOR.

## TODO

* Some more sanity checks

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
