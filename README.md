# androidmanifest-changer

Change Android AAB/APK attributes like the versionCode and versionName. This tool modifies the binary AndroidManifest.xml within AAB (Bundles) and APK files.


## Supported Attributes

**minSdkVersion**
   - An integer designating the minimum API level required for the application to run. 
   - The Android system prevents the user from installing the application if the system's API level is lower than the value specified in this attribute. 
   - Always declare this attribute.
   - Caution: If you don't declare this attribute, the system assumes a default value of "1", which indicates that your application is compatible with all versions of Android. 
   - If it isn't, and you didn't declare the proper minSdkVersion, then when installed on a system with an incompatible API level, the application crashes during runtime when attempting to access the unavailable APIs.
   - For this reason, be certain to declare the appropriate API level in the minSdkVersion attribute.
   - See https://developer.android.com/guide/topics/manifest/uses-sdk-element for more information.

**versionCode**
   - A positive integer used as an internal version number.
   - This number helps determine whether one version is more recent than another, with higher numbers indicating more recent versions.
   - This is not the version number shown to users; that number is set by the versionName setting.
   - The Android system uses the versionCode value to protect against downgrades by preventing users from installing an APK with a lower versionCode than the version currently installed on their device.
   - The value is a positive integer so that other apps can programmatically evaluate it—to check an upgrade or downgrade relationship, for instance.
   - You can set the value to any positive integer. However, make sure that each successive release of your app uses a greater value.
   - Note: The greatest value Google Play allows for versionCode is 2100000000.
   - You can't upload an APK to the Play Store with a versionCode you have already used for a previous version.
   - Note: In some situations, you might want to upload a version of your app with a lower versionCode than the most recent version.
   - For example, if you are publishing multiple APKs, you might have pre-set versionCode ranges for specific APKs.
   - For more about assigning versionCode values for multiple APKs, see Assigning version codes.
   - Typically, you release the first version of your app with versionCode set to 1, then monotonically increase the value with each release, regardless of whether the release constitutes a major or minor release.
   - This means that the versionCode value doesn't necessarily resemble the app release version that is visible to the user.
   - Apps and publishing services shouldn't display this version value to users.
   - See https://developer.android.com/studio/publish/versioning for more information.

**versionName**
   - A string used as the version number shown to users. This setting can be specified as a raw string or as a reference to a string resource.
   - The value is a string so that you can describe the app version as a <major>.<minor>.<point> string or as any other type of absolute or relative version identifier.
   - The versionName is the only value displayed to users.
   - See https://developer.android.com/studio/publish/versioning for more information.

**package**
   - The value of the package attribute in the APK's manifest file represents your app's universally unique application ID.
   - It is formatted as a full Java-language-style package name for the Android app.
   - The name can contain uppercase or lowercase letters, numbers, and underscores ('_'). However, individual package name parts can only start with letters.
   - Be careful not to change the package value, since that essentially creates a new app.
   - Users of the previous version of your app don't receive an update and can't transfer their data between the old and new versions.
   - In the Gradle-based build system, starting with AGP 7.3, don't set the package value in the source manifest file directly.
   - For more information, see Set the application ID.
   - See https://developer.android.com/guide/topics/manifest/manifest-element for more information.


## Coming soon

**targetSdkVersion**
   - Similar to Java, the targetSdkVersion of your app can change the runtime behavior of native code.
   - Behavior changes in the system are, when feasible, only applied to apps with a targetSdkVersion greater than or equal to the OS version that introduced the change.
   - For a new app, choose the newest version available. For an existing app, update this to the latest version when convenient (after updating compileSdkVersion).
   - While application developers generally know their app's targetSdkVersion, this API is useful for library developers that cannot know which targetSdkVersion their users will choose.
   - At runtime, you can get the targetSdkVersion used by an application by calling android_get_application_target_sdk_version().
   - This API is available in API level 24 and later. This function has the following signature:
```
/**
 * Returns the `targetSdkVersion` of the caller, or `__ANDROID_API_FUTURE__` if
 * there is no known target SDK version (for code not running in the context of
 * an app).
 *
 * The returned values correspond to the named constants in `<android/api-level.h>`,
 * and is equivalent to the AndroidManifest.xml `targetSdkVersion`.
 *
 * See also android_get_device_api_level().
 *
 * Available since API level 24.
 */
int android_get_application_target_sdk_version() __INTRODUCED_IN(24);
```
   - Other behavior changes might depend on the device API level. You can get the API level of the device your application is running on by calling android_get_device_api_level(). This function has the following signature:
```
/**
 * Returns the API level of the device we're actually running on, or -1 on failure.
 * The returned values correspond to the named constants in `<android/api-level.h>`,
 * and is equivalent to the Java `Build.VERSION.SDK_INT` API.
 *
 * See also android_get_application_target_sdk_version().
 */
int android_get_device_api_level();
```
   - See https://developer.android.com/ndk/guides/sdk-versions for more information.

**compileSdkVersion**
   - This property has no effect on NDK builds. API availability for the NDK is instead governed by minSdkVersion.
   - This is because C++ symbols are eagerly resolved at library load time rather than lazily resolved when first called (as they are in Java).
   - Using any symbols that are not available in the minSdkVersion will cause the library to fail to load on OS versions that do not have the newer API, regardless of whether or not those APIs will be called.
   - For a new app, choose the newest version available. For an existing app, update this to the latest version when convenient.
   - See https://developer.android.com/ndk/guides/sdk-versions for more information.

**compileSdkVersionCodename**
   - The development codename (ex. "S", "REL") of the framework against which the application claims to have been compiled, or null if not specified.
   - This property is the compile-time equivalent of Build.VERSION.CODENAME.
   - See https://developer.android.com/reference/android/content/pm/ApplicationInfo#compileSdkVersionCodename for more information.

**platformBuildVersionCode**


**platformBuildVersionName**


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
