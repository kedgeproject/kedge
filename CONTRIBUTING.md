# Contributing guidelines

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

The following is a set of guidelines (not rules) for contributing to Kedge.

These are just guidelines, not rules, use your best judgment and feel free to
propose changes to this document in a pull request.

## Submitting a Pull Request
Before you submit your pull request consider the following guidelines:

* Make your changes in a new git branch:

     ```shell
     git checkout -b bug/my-fix-branch master
     ```

* Create your patch, **including appropriate test cases**.
* Include documentation that either describe a change to a behavior of kedge or the changed capability to an end user of kedge.
* Commit your changes using **a descriptive commit message**. If you are fixing an issue please include something like 'this closes issue #xyz'.
* Make sure your tests pass! 

    ```shell
    make test
    ```

* Push your branch to GitHub:

    ```shell
    git push origin bug/my-fix-branch
    ```

* In GitHub, send a pull request to `kedge:master`.
* If we suggest changes then:
  * Make the required updates.
  * Rebase your branch and force push to your GitHub repository (this will update your Pull Request):

    ```shell
    git rebase master -i
    git push origin bug/my-fix-branch -f
    ```

That's it! Thank you for your contribution!

### Merge Rules

* Include unit or integration tests for the capability you have implemented
* Include documentation for the capability you have implemented
* If you are fixing an issue within Kedge, include the issue number you are fixing

### After your pull request is merged

After your pull request is merged, you can safely delete your branch and pull the changes
from the upstream repository:

* Delete the remote branch on GitHub either through the GitHub web UI or your local shell as follows:

    ```shell
    git push origin --delete bug/my-fix-branch
    ```

* Check out the master branch:

    ```shell
    git checkout master -f
    ```

* Delete the local branch:

    ```shell
    git branch -D bug/my-fix-branch
    ```

* Update your master with the latest upstream version:

    ```shell
    git pull --ff upstream master

## Git Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Reference issues and pull requests liberally
