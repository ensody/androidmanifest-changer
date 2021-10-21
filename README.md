# androidmanifest-changer

Change Android AAB/APK attributes like the versionCode and versionName. This tool modified the binary AndroidManifest.xml within AAB (Bundles) and APK files.

## Supported attributes

* versionCode
* versionName
* package

## Usage

```
# Change only versionCode
androidmanifest-changer --versionCode 4 app.aab

# Change multiple values
androidmanifest-changer \
  --versionCode 4 \
  --versionName 1.0.2 \
  --package com.some.app \
  app.aab
```

This will rewrite the given aab/apk with the new values.

## Requirements

These tools must be installed and reachable on your PATH:

* zip (Go's built-in zip library produces invalid aab/apk files; TODO: find a workaround)
* aapt2 (only if you want to manipulate APKs)
