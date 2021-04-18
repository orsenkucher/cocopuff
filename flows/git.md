### Rewriting history
It is possible to revrite git history  
https://stackoverflow.com/questions/17577409/git-remove-merge-commit-from-history

### Fixing history

After you screwed up all the history with `git rebase -i ...` and want to reset timestamps as they were before, run

```bash
git filter-branch --env-filter 'export GIT_COMMITTER_DATE="$GIT_AUTHOR_DATE"'
```
> `Explanation`: What is going on: when your rebase, the committer's timestamp changes, but not the author's timestamp, which suddenly all makes sense.

### Sign all the commits
But now all the newly timestamped history lost GPG signes...  
```bash
git filter-branch --commit-filter 'git commit-tree -S "$@";' -- --all
```
Fixes it.

### Merging unrelated histories
Now you want to merge one repository with another without losing file history  
https://saintgimp.org/2013/01/22/merging-two-git-repositories-into-one-repository-without-losing-file-history/
```bash
# in lama
git checkout -b cocopuff
git remote add -f cocopuff <Repo URL>
git merge cocopuff/lama --allow-unrelated-histories
```
And then Merge PR to keep history.

### Bring over a feature branch from one of the old repos
If we have in-progress feature branches in the old repositories that also need to come over to the new repository, that’s also quite easy:
```bash
git checkout -b feature-in-progress
git merge -s recursive -Xsubtree=old_a old_a/feature-in-progress
```
This is the only non-obvious part of the whole operation.  We’re doing a normal recursive merge here (the “-s recursive” part isn’t strictly necessary because that’s the default) but we’re passing an argument to the recursive merge that tells Git that we’ve renamed the target and that helps Git line them up correctly.

### Stackoverflow
https://stackoverflow.com/questions/2973996/git-rebase-without-changing-commit-timestamps
>If you've already screwed up the commit dates (perhaps with a rebase) and want to reset them to their corresponding author dates, you can run:
>```bash
>git filter-branch --env-filter 'GIT_COMMITTER_DATE=$GIT_AUTHOR_DATE; export GIT_COMMITTER_DATE'
>```


>I'd like to add another approach if you've already screwed up but don't want to iterate through the whole history: git rebase --committer-date-is-author-date <base_branch> This way, git will reset the commit date only for the commits applied upon <base_branch> (which is probably the same branch name you used when you screwed up).

> A crucial question of Von C helped me understand what is going on: when your rebase, the committer's timestamp changes, but not the author's timestamp, which suddenly all makes sense. So my question was actually not precise enough.  
>The answer is that rebase actually doesn't change the author's timestamps (you don't need to do anything for that), which suits me perfectly.

https://stackoverflow.com/questions/30790645/how-to-make-a-git-rebase-and-keep-the-commit-timestamp/30819930

https://stackoverflow.com/questions/17577409/git-remove-merge-commit-from-history
