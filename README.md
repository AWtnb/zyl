# moko

曖昧模糊 _aimai-**moko**_ means fuzzyness in Japanese.

+ fuzzy file/folder launcher for windows.
+ search target configurable with `.yaml`


```
> moko -h
Usage of moko.exe:
  -all
        switch in order to search including file
  -exclude string
        search exception (comma-separated)
  -filer string
        filer path (default "explorer.exe")
  -src string
        source yaml file path
```

## `launch.yaml`

```yaml
- path: C:\Users\%USERNAME%\Desktop
  alias:
  depth: -1

- path: C:\Users\%USERNAME%\Dropbox
  alias: dbx
  depth: 3
```

+ `alias` : Name displayed during a fuzzy search. Default is the file or folder name at the end of the path.
+ `depth` : Specifies the depth of hierarchy when searching folders.
    + `0` (default) : Open the selected folder in the filer without searching for subfolders.
    + `-1` : Search all folders except hidden folders.