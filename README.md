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

## .zshrc useful functions, alias' and autocomplete

kubectl autocomplete and alias:

```bash
alias k="kubectl"
complete -F __start_kubectl k
```

prune: removes all git local git branches apart from the current branch and the main branch

```bash
prune(){
    main=$(git remote show origin | grep 'HEAD branch' | awk '{print $3}')
    git branch | grep -v '*' | grep -v $main | xargs git branch -D
}
```

rebase: checkouts the main branch, pulls any upstream changes, checkout current branch and rebases from main branch

```bash
rebase(){
    branch=$(git branch | grep '*' | awk '{print $2}')
    main=$(git remote show origin | grep 'HEAD branch' | awk '{print $3}')
    git checkout $main 
    git pull
    git checkout $branch
    git rebase $main
}
```
