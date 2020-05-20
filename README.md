# DEPRECATED

Use the git resource now that this is merged: https://github.com/concourse/git-resource/pull/225

# Git Tags Resource

Tracks github git tags (Annotated or Lightweight) regardless of the branch.

## IMPORTANT

**Be aware that this resource only supports tags in line with [Semver](https://semver.org/)**
**All non-semver tags will be dropped/ignored**

```yaml
resource_types:
- name: git-tags-resource
  type: docker-image
  source:
    repository: adgear/git-tags-resource
```

## Source Configuration

* `repository_name`: *Required.* The repository name.

* `uri`: *Optional.* The git URI. Defaults to git@github.com:`repository_name`.git

* `tag_filter`: *Optional.* The glob pattern to match tags against. Defaults to "*.*.*"

* `private_key`: *Optional.* Required if uri starts with `git@`. The private key string to clone the repository.

* `PrivateKeyPassword`: *Optional.* Required if `private_key` is password protected. The password for the private key.

* `LatestOnly`: *Optional.* Get only the latest tag. Defaults to `true`.

### Example

Resource configuration for incubator repository

``` yaml
resources:
- name: concourse-tags
  type: git-tags-resource
  source:
    repository_name: concourse/concourse
    uri: "https://github.com/concourse/concourse.git"
    tag_filter: "*.*.*"
```

Resource configuration for private repository

``` yaml
resources:
- name: concourse-tags
  type: git-tags-resource
  source:
    repository_name: concourse/concourse
    uri: "git@github.com/concourse/concourse.git"
    tag_filter: "*.*.*"
    private_key: "----RSA----\nKEY\n----RSA----"
```

## Behavior

### `check`: Check for new git tags

Search for the latest version of `source.uri`.

### `in`: Get the latest tag ref.

Output `tag`, `ref`, `shortref`, `tag_type`, `committer`, `author` in `.git` folder.

### `out`: Nothing

No uses.

## Development

### Prerequisites

* Common sense.

### Running the tests

To be implemented.

### License

[MIT](LICENSE)

### Contributing

TBD.
