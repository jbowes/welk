<!--
  Attractive html formatting for rendering in github. sorry text editor
  readers! Besides the header and section links, everything should be clean and
  readable.
-->
<h1 align="center">welk</h1>
<p align="center"><i>What's inside the shell? welk manages `curl | sh` style package installs 🐌</i></p>

<div align="center">
  <img alt="Alpha Quality" src="https://img.shields.io/badge/status-ALPHA-orange.svg" >
  <a href="https://github.com/jbowes/welk/releases/latest"><img alt="GitHub tag" src="https://img.shields.io/github/tag/jbowes/welk.svg"></a>
  <a href="https://github.com/jbowes/welk/actions/workflows/go.yml"><img alt="Build Status" src="https://github.com/jbowes/welk/actions/workflows/go.yml/badge.svg?branch=main"></a>
  <a href="./LICENSE"><img alt="BSD license" src="https://img.shields.io/badge/license-BSD-blue.svg"></a>
  <a href="https://codecov.io/gh/jbowes/welk"><img alt="codecov" src="https://img.shields.io/codecov/c/github/jbowes/welk.svg"></a>
  <a href="https://goreportcard.com/report/github.com/jbowes/welk"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/jbowes/welk"></a>
</div><br /><br />

---

🚧 ___Disclaimer___: _`welk` is alpha quality software. The API may change
without warning between revisions._ 🚧

`welk` is the package manager for software installed with [`curl`][curl] and
[`sh`][sh]. `welk` gives you:
- A sandboxed environment for script execution during install
- Tracking of installed software. No more `$HOME` littered with random
  [Kubernetes][k8s] tools!
- Deletion of installed software.
- Automatic `$PATH` management. Your [`.zshrc`][zsh] has never looked better!

When documentation tells you to run `curl <some url> | /bin/sh`, use `welk`.

## Getting `welk`

Use `curl | sh`, of course! Run:
```sh
curl <TODO ADD FLAGS> <TODO URL> | sh
```

After `welk` is installed it will manage itself.

Make sure `.local/bin` is in your path, as defined by [XDG].

## Quick start

- Install new software with `welk install $URL`
- List installed software with `welk list`
- See the files installed by a script with `welk info $URL`
- Remove everything a script installed with `welk delete $URL`

## Seasons

Open Source is free to play. Like all good free to play experiences, `welk`
releases content and feature drops in seasons. Content from past seasons
remains available to new users; we won't put it in a vault!

The current season is the [*Season of discovery*](./milestone/2).

The next season is the [*Season of the herald*](./milestone/3).

Find more upcoming and past seasons under [milestones][milestones]

## Contributing

We would love your help!

`welk` is still a work in progress. You can help by:

- Opening a pull request to resolve an [open issue][issues].
- Adding a feature or enhancement of your own! If it might be big, please
  [open an issue][enhancement] first so we can discuss it.
- Improving this `README` or adding other documentation to `welk`.
- Letting [me] know if you're using `welk`.

[curl]: https://curl.se/
[sh]: https://en.wikipedia.org/wiki/Unix_shell
[k8s]: https://kubernetes.io/
[zsh]: https://www.zsh.org/
[xdg]: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html

[miletones]: ./milestones
[issues]: ./issues
[bug]: ./issues/new?labels=bug
[enhancement]: ./issues/new?labels=enhancement

[me]: https://twitter.com/jrbowes
