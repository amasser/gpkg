**WARNING**: this is currently still in the _prototyping_ stage, do not use.

![go](https://github.com/git-pkg/gpkg/workflows/go/badge.svg)
[![codecov](https://codecov.io/gh/git-pkg/gpkg/branch/main/graph/badge.svg)](https://codecov.io/gh/git-pkg/gpkg)
[![Gitter](https://badges.gitter.im/git-pkg/community.svg)](https://gitter.im/git-pkg/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

# Git Package Manager

Package management using [Github][github], [Gitlab][gitlab] and vanilla [Git][git] repositories.

**TODO**: In it's current state this tool only works with Github git repositories as that was the highest priority.

# About

This tool started because I had several tools that I was manually installing from [Github][github] repositories and I wanted to automate that process in a simple and generic way requiring no special permissions.

For example: I regularly use the [Github CLI](https://github.com/cli/cli). I use to download tarballs from the their [Releases][ghreleases] but I wanted instead to be able to run something like this to install:

```go
gpkg install github.com/cli/cli
```

And then when I need to update to the latest patch release:

```go
gpkg update github.com/cli/cli
```

This is the kind of functionality this tool provides in a nutshell for a variety of repositories.

# Installation

```shell
GO111MODULE="on" go get github.com/git-pkg/gpkg/cmd/gpkg
export PATH="$HOME/.local/gpkg/bin:$PATH"
```

# Goals

To track current progress on check the [Organization's Projects][projects].

## Immediate Goals

* Support [Linux][linux] based operating systems on [64bit x86 Architecture][x86_64]
* Support [Github][github] based Git repositories using [Releases][githubreleases]

## Upcoming

* Support [Gitlab][gitlab] based Git repositories using [Release][gitlabreleases]
* Support generic [Git][git] repositories utilizing [Tags][tags]

# Contributing

Currently this is being prototyped. Feel free to contribute during this time, or if you're waiting for something a little more stable keep a lookout for our [first release][milestone1]

[git]:https://git-scm.com
[tags]:https://git-scm.com/book/en/Git-Basics-Tagging
[linux]:https://kernel.org
[x86_64]:https://en.wikipedia.org/wiki/X86-64
[github]:https://github.com
[githubreleases]:https://developer.github.com/v3/repos/releases/
[gitlab]:https://gitlab.com
[gitlabreleases]:https://docs.gitlab.com/ee/user/project/releases/
[projects]:https://github.com/orgs/git-pkg/projects
[milestone1]:https://github.com/git-pkg/gpkg/milestone/1
