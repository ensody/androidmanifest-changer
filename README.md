# androidmanifest-changer

Change Android AAB/APK attributes like the versionCode and versionName. This tool modified the binary AndroidManifest.xml within AAB (Bundles) and APK files.

## Supported attributes

* minSdkVersion
* versionCode
* versionName
* package

## Usage

```
# Change only versionCode
androidmanifest-changer --versionCode 4 app.aab

# Change multiple values
androidmanifest-changer \
  --minSdkVersion 33 \
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


## License

```
Copyright 2023 Ensody GmbH, Waldemar Kornewald

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
