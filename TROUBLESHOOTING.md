# Troubleshooting Grit

## SSH

If you're having SSH authentication issues similar to:

```
github: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain
```

Here's a few things to try.

### macOS Keychain

If you're using the default macOS SSH setup, make sure the SSH keys in your
keychain are being added to the SSH agent:

```sshconfig
# ~/.ssh/config
Host *
    UseKeychain yes
    AddKeysToAgent yes
```

If the above config doesn't work, you can alternatively add the following
command to `.zshrc`:

```zsh
# ~/.zshrc
ssh-add --apple-load-keychain
```

### 1Password

If you're using 1Password's SSH agent, you'll need to set the `SSH_AGENT_SOCK`
environment variable in your `.zshrc`/`.bashrc` (or similar) to point to the
1Password agent on your OS. Under macOS, this looks like:

```zsh
# ~/.zshrc
export SSH_AUTH_SOCK="$HOME/Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock"
```

## Autocompletion

### ZSH

In order to get completion working under ZSH, you'll need to enable completion,
and also add support for Bash completions. Both are typically configured in
`.zshrc`:

```zsh
# ~/.zshrc
autoload -U compinit && compinit
autoload -U bashcompinit && bashcompinit
```

## Grit doesn't change directory when I `grit clone` / `grit cd`

You haven't enabled shell integration. You're probably missing this line from `.zshrc`/`.bashrc`:

```zsh
# ~/.zshrc
eval "$(grit shell-integration)"
```
