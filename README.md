# random_scripts

This repo is a place to hold random scripts i have created for fun or to help me with day-to-day tasks and list useful packages and extensions that i've come across

## Useful packages

- [pwgen](https://formulae.brew.sh/formula/pwgen): useful for generating passwords
- [Kubectx and Kubens](https://github.com/ahmetb/kubectx): tools that allow you to switch K8s contexts and namespaces faster. remember to add a alias for this command: `alias kx="kubectx"` and `alias kn="kubens"`
- [Kube-ps1](https://github.com/jonmosco/kube-ps1): displays your current k8s context and namespace on the command prompt
- [terraform-docs](https://github.com/terraform-docs/terraform-docs): useful for generating documentation for terraform modules

- [ksd](https://github.com/mfuentesg/ksd): quick way to decode k8s secrets
- [iterm2](https://iterm2.com/): good replacement terminal for macOS
- [ohmyzsh](https://github.com/ohmyzsh/ohmyzsh): framework for managing your zsh configuration (requires zsh to be installed)

## Useful VSCode extensions

- Git Graph: used to view the Git graph of a repository
- GitLens: useful for checking when a line of code was last edited, running git blame against a file, what changes where made in a comment etc.
- HashiCorp Terraform: syntax highlighting and (some) autocompletion for terraform code
- markdownlint: useful for linting markdown documentation
- vscode-icons: more icons for different file types
- YAML: adds YAML language support
- Docker: adds Dockerfile support

## .zshrc useful functions, alias' and autocomplete

### kubectl autocomplete and alias

```bash
source <(kubectl completion zsh)
alias k="kubectl"
complete -F __start_kubectl k

# use if you have installed kubectx and kubens
alias kx="kubectx"
alias kn="kubens"
```

### prune

removes all git local git branches apart from the current branch and the main branch

```bash
prune(){
    main=$(git remote show origin | grep 'HEAD branch' | awk '{print $3}')
    git branch | grep -v '*' | grep -v $main | xargs git branch -D
}
```

### rebase

checkouts the specified branch (uses main/master by default), pulls any upstream changes, checkouts current branch and rebases from specified branch

```bash
rebase(){
    my_branch=$(git branch | grep '*' | awk '{print $2}')
    main=$(git remote show origin | grep 'HEAD branch' | awk '{print $3}')

    if [[ -z $1 ]]
    then
        echo "no branch specified. using main/master"
        base_branch=$main
    else
        echo "using branch: $1"
        base_branch=$1
    fi

    git checkout $base_branch 
    git pull
    git checkout $my_branch
    git rebase $base_branch
}
```
