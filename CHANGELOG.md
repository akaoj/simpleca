# 1.2.1 (2018-10-17)

### Buildchain

- Use a Docker-make based workflow, so developers don't have to install anything on their machine to compile or test
  the code. Now you can `git pull` and `make compile` without having installed `go` and everything will work.

### Bug fixes

- Use `SubjectAltName` in addition to `CommonName` (see [RFC 2818](https://tools.ietf.org/html/rfc2818#page-5))



# 1.2 (2018-04-12)

### Improvements

- Add an `init` command ([issue #6](https://github.com/akaoj/simpleca/issues/6)).

  Usage:
  ```
  $ simpleca init
  Folder initialized, please edit the configuration.json file to fit your organization
  ```
- Add a `rm` command ([issue #10](https://github.com/akaoj/simpleca/issues/10)).

  Usage:
  ```
  $ simpleca rm client --name www.domain.com
  client keys and certificates deleted
  ```
- Add this CHANGELOG.md ([issue #15](https://github.com/akaoj/simpleca/issues/15)).

### Bug fixes

- Do not fail tests on master because of version ([issue #13](https://github.com/akaoj/simpleca/issues/13)).

  A warning is displayed in red when building simpleca from a branch which does not match the current simpleca version
  but the tests don't fail anymore if the current branch is master.



# 1.0.1 (2018-03-21)

### Bug fixes

- Fix issue with intermediate certificates

  Intermediate certificates would not behave as expected: certificates signed by intermediates were not valid.



# 1.0 (2018-02-07)

This is the first official release of simpleca.

`simpleca` allows you to:
- create root CAs
- create intermediate CAs
- create client keys
- sign anything with anything (more likely: clients with intermediates and intermediates with root)

All generated keys are by default encrypted.

As of now, only RSA and ECDSA keys are supported (ECDSA being the default generated keys).
