# moko

曖昧模糊 あいまいもこ _aimai-**moko**_ means fuzzyness in Japanese.

- fuzzy file/folder launcher for windows.
- search target configurable with `.yaml`


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
- path: C:\Users\${USERNAME}\Desktop
  alias:
  depth: -1

- path: C:\Users\${USERNAME}\Dropbox
  alias: dbx
  depth: 3
```

- `alias` : Name displayed during a fuzzy search. Default is the file or folder name at the end of the path.
- `depth` : Specifies the depth of hierarchy when searching folders.
    - `0` (default) : Open the selected folder in the filer without searching for subfolders.
    - `-1` : Search all folders except hidden folders.
- File search on network directory is slow. In that case, [Everything](https://www.voidtools.com) can be used.
    - Requirement: Everything is running on PC and `Everything64.dll` exists on the same directory with `moko.exe` .
    - `Everything64.dll` is released on [official site](https://www.voidtools.com/support/everything/sdk/) .

Great thanks for:

- https://www.voidtools.com/support/everything/sdk/
- https://github.com/jof4002/Everything/blob/master/everything_windows_amd64.go