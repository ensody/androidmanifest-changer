# androidmanifest-changer

Change Android AAB/APK attributes like the versionCode and versionName. This tool modifies the binary AndroidManifest.xml within AAB (Bundles) and APK files.


## Supported Attributes

**minSdkVersion**
   - Integer specifying the minimum Android API level required for the app to run.
   - Android system blocks installation on devices with API levels lower than this value.
   - Default assumed value is "1" if not specified, implying compatibility with all Android versions.
   - Not declaring this attribute can lead to app crashes on incompatible systems due to unavailable APIs.
   - Critical to specify the correct minSdkVersion for app stability.
   - See https://developer.android.com/guide/topics/manifest/uses-sdk-element for more information.


**versionCode**
   - A positive integer that serves as the internal version number of the app.
   - Used to differentiate between newer and older versions, with higher numbers indicating more recent versions.
   - Not visible to users, as the versionName attribute is used for display.
   - Prevents downgrading by blocking the installation of APKs with lower versionCodes than the one installed.
   - Important to increment for each app update.
   - Maximum allowable value on Google Play: 2100000000.
   - Reuse of versionCode for Play Store uploads is not permitted.
   - See https://developer.android.com/studio/publish/versioning for more information.


**versionName**
   - The version number shown to users, specified as a string.
   - Flexible format, commonly used as <major>.<minor>.<point> or other version identifiers.
   - The primary version identifier for end users.
   - See https://developer.android.com/studio/publish/versioning for more information.


**package**
   - Represents the app's unique application ID, formatted as a Java package name.
   - Can include uppercase and lowercase letters, numbers, and underscores, but must start with a letter.
   - Modifying this value essentially creates a new application, impacting updates and data transfer.
   - In AGP 7.3+, avoid setting this directly in the source manifest.
   - See https://developer.android.com/guide/topics/manifest/manifest-element for more information.


## Coming soon

**targetSdkVersion**
   - Influences the runtime behavior of the app's native code.
   - System applies behavior changes to apps with targetSdkVersion at or above the OS version introducing these changes.
   - New apps should target the latest version; existing apps should update when feasible.
   - Retrieve the targetSdkVersion at runtime with android_get_application_target_sdk_version() in API level 24 and later.
   - See https://developer.android.com/ndk/guides/sdk-versions for more information.


**compileSdkVersion**
   - Determines the API availability for NDK builds, independent of the compileSdkVersion property.
   - Governed by minSdkVersion, as C++ symbols are resolved at library load time.
   - Recommended to use the newest version for new apps and update existing apps as needed.
   - See https://developer.android.com/ndk/guides/sdk-versions for more information.


**compileSdkVersionCodename**
   - Reflects the development codename of the Android framework used for compiling the app.
   - Compile-time equivalent of Build.VERSION.CODENAME.
   - See https://developer.android.com/reference/android/content/pm/ApplicationInfo#compileSdkVersionCodename for more information.


**platformBuildVersionCode**
   - Description unavailable.


**platformBuildVersionName**
   - Description unavailable.


## Usage

```bash
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
