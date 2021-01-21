Simpson
=======

Simpson is a

**B**uild  
**A**nd  
**R**elease  
**T**ool

for Go on Github.

Considerations
--------------

### Fetching Tags ###

The tag `latest` will change with every push. Fetching/pulling tags will thus fail with the error:
```
! [rejected]        latest     -> latest  (would clobber existing tag)
```

On the commandline you must add the option `-f`:
```
git fetch --tags -f
git pull --tags -f
```

In VSCode you should change the setting `git.pullTags` to `false` and fetch tags
manually when required.
