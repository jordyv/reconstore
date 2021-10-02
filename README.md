# Reconstore

Reconstore is a tool for saving and querying your recon data. You can pipe subdomains and output of your favorite tools to reconstore to build up 
your own recon database. All information is stored in a Sqlite file so you can run your own queries if you like.

## Installation

```
go get -u github.com/jordyv/reconstore/cmd/reconstore
```

## Usage

### Programs

Add your first program:
```
reconstore programs add --name Tesla --platform BugCrowd --bounties
```

Get a list of all your programs:
```
reconstore programs list
ID  Name       Slug       Platform   Private  Has bounties
1   Tesla      tesla      BugCrowd   no       yes
```

### Subdomains
#### Store subdomains
```
echo tesla.com | subfinder | reconstore subdomains save -p tesla
```

#### Tag subdomains
```
cat nuclei | grep wordpress | cut -d ' ' -f 2 | reconstore subdomains tag -t wordpress
```
This will add a tag 'wordpress' to all already existing subdomains (this won't store new subdomains, only existing subdomains are tagged).

#### Query subdomains
Only subdomains from a single program:
```
reconstore subdomains query --slug tesla
```

Only subdomains with a paying program:
```
reconstore subdomains query --bounties
```

Subdomains with a tag:
```
reconstore subdomains query --tag wordpress
```

Subdomains with a search query:
```
reconstore subdomains query --pattern docs
```

## Planned features

- [x] DNS resolution on import of subdomains
- [ ] Query DNS info
- [ ] HTTP info - webserver, status code, title
- [ ] Techs to tags - scrape HTTP, fetch technologies and store in DB
- [ ] Import nmap port scan results
- [ ] GraphQL API
- [ ] Web interface
