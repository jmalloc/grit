# Grit

> Not working? Try the [troubleshooting guide].

[troubleshooting guide]: TROUBLESHOOTING.md

## What?

Grit keeps track of your local clones of Git repositories, allowing you to
quickly switch to clone directories in your terminal based on the repository
name or a portion thereof.

## Why?

I spend most of my day working with Git. Many of the repositories are hosted on
GitHub.com, but many more are in my employer's private GitHub Enterprise and
BitBucket installations.

Keeping track of hundreds of clones can be a little tedious, so some time back
I adopted a basic directory naming convention and wrote some shell scripts to
handle cloning in a consistent way.

This worked well for a while, until the list of places I needed to clone from
increased further, and I started working more heavily in [Go](http://golang.org),
which, at the time, placed it's [own requirements](https://github.com/golang/go/wiki/GOPATH)
on the location of your Git clones.

Grit is the logical evolution of those original scripts into a standalone
project that clones from multiple Git sources. It's hacked together, there are
no tests, and there are other more general solutions for navigating your
filesystem; but this works for me. I've published it to GitHub in case it works
for you too.

## How?

If you want to get a sense of what Grit does without having to read or
understand anything, try this:

1. Download `grit` from the [releases page](https://github.com/jmalloc/grit/releases) or [install with Homebrew](https://github.com/jmalloc/homebrew-grit)
1. Clone a repository with `grit clone <repo>` maybe try `grit clone jmalloc/grit`
1. Find a repository with `grit index find <repo>`

Joke's on you! Grit only knows about that one repository. Now you have to read
a bunch of stuff anyway ...

### Configuration

Grit is useful because it knows where you want to clone repositories from. But
it's not a magician, it only knows after you tell it, which is less impressive.

By default, Grit looks for a configuration file at `~/.config/grit.toml`, which
at its most basic is a list of named clone sources, such as:

```toml
[clone.sources]
my-company = "git@git.example.com:{{slug}}.git"
```

If you only use GitHub you don't need to define any additional sources, and so
the configuration file can be omitted entirely. For a complete list of the
available configuration directives, see the [example](etc/complete-example.toml)
configuration file.

From now on, I'll assume you're using the configuration shown above, which is a
bit silly, but here we are. Rest assured, all paths and filenames used by Grit
are probably 100% configurable, more or less.

### Cloning a repository

There are two parts to Grit. One of them is the `clone` command which clones
repositories. It's like `git clone` but without any of the flexibility or power.
That's minimalism.

The clone command accepts a single argument, the repository "slug". The slug
is the part of the repository URL that a standard-issue human would use to
identify a repository. For GitHub, this is the familiar `username/repository`
syntax. For sources configured in `grit.toml`, it's the part represented
by the string `{{slug}}` in the URL.

Try it out:

    grit clone jmalloc/grit

Grit will print the absolute path of the local clone, which should look
something like:

    /Users/james/grit/github.com/jmalloc/grit

I'm assuming a lot about your system configuration and your given name, just
roll with it.

By the way, you can pass a complete Git URL instead of just the slug and Grit
will clone it into the correct location, even if the URL does not match any of
your configured sources.

### Querying the index

Whenever you clone a repository with Grit, the repository is added to the index.
The index is a database mapping repository slugs to directories. The index is
also that second part of Grit that I mentioned back when we were talking about
cloning. I haven't forgotten.

The index can be queried to find a repository by slug:

    grit index find jmalloc/grit

Or by the part after the slash, let's call this the repository name:

    grit index find grit

Just like the `clone` command, `index find` prints the clone path. If there is
more than one matching path, it just prints them all. It's relentless.

If nothing is found, Grit exits with a non-zero exit code. The universal signal
that you should re-think your actions.

### Selecting between multiple clones

"What happens if I have two repositories with the same name?", you ask me in a
contrived way perfectly suited to this next example.

Well, try this...

First, clone the `mojombo/grit` repository; an unsuspecting pawn in our example
that has nothing to do with `jmalloc/grit`. Then, change to the directory of
your "grit" clone.

    grit clone mojombo/grit
    grit cd grit

If you've followed all the steps until now, and I haven't messed up the examples
too badly, you should be presented with a list of matching repositories:

    Use the arrow keys to navigate: ↓ ↑ → ←
    ? :
      ▸ github.com/jmalloc/grit
        [current] github.com/mojombo/grit

This time, Grit isn't content to blindly print paths to the terminal. Oh no!
This time Grit wants to know what's on your mind. It prints a list of all
matching directories and eagerly awaits your decision. Navigate to your
selection using the arrow keys and press enter.

Grit diligently prints the absolute path to your selection, just like before.
Interminable!

## Ok, so Grit prints paths to the terminal, I get it.

> Anybody can print to a terminal, even a programmer!
> Hell, Grit's `cd` command doesn't even change the current directory!

That's what you sound like right now. Shh bby is ok.

What we need is some glue that makes Grit's `cd` command behave like the _real_
`cd` command, and that glue is:

    eval "$(grit shell-integration)"

Add this line to your `.bash_profile` file and many of the following things will
happen:

1. You'll get auto-completion of all command names and indexed repository slugs
1. `grit cd` will start working the way you expect
1. `grit clone` will also change the current directory after cloning

To achieve the latter two, Grit needs to execute shell commands in the context
of it's parent shell (the one you ran Grit from). It does this by writing the
commands to a separate file which is then sourced by the parent shell. The
details can be seen by running `grit shell-integration` directly from your
terminal.

## What next?

My colleagues are helping me iron out some of the kinks. Undoubtedly there are
still some changes to come, but Grit is feature complete insofar as it already
does everything that my original shell scripts could do.

If you find Grit useful, or have a feature request or bug-report, please don't
hesitate to create a new [issue](https://github.com/jmalloc/grit/issues).
